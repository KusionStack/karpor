package clusterextension

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"code.alipay.com/ant-iac/karbour/pkg/apis/cluster"
	proxyutil "code.alipay.com/ant-iac/karbour/pkg/util/proxy"
	"k8s.io/apimachinery/pkg/util/proxy"
)

var proxyMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

var _ rest.Connecter = &ProxyREST{}

type ProxyREST struct {
	store *genericregistry.Store
}

func (r *ProxyREST) Destroy() {
}

// New returns an empty cluster proxy subresource.
func (r *ProxyREST) New() runtime.Object {
	return &cluster.ClusterExtensionProxyOptions{}
}

func (r *ProxyREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &cluster.ClusterExtensionProxyOptions{}, true, "path"
}

func (r *ProxyREST) ConnectMethods() []string {
	return proxyMethods
}

func (r *ProxyREST) Connect(ctx context.Context, id string, options runtime.Object, responder rest.Responder) (http.Handler, error) {
	startTime := time.Now()
	proxyOpts, ok := options.(*cluster.ClusterExtensionProxyOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", options)
	}

	parentObj, err := r.store.Get(ctx, id, &metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("no such cluster %v", id)
	}
	clusterExtension := parentObj.(*cluster.ClusterExtension)

	reqInfo, _ := request.RequestInfoFrom(ctx)

	factory := request.RequestInfoFactory{
		APIPrefixes:          sets.NewString("api", "apis"),
		GrouplessAPIPrefixes: sets.NewString("api"),
	}
	proxyReqInfo, _ := factory.NewRequestInfo(&http.Request{
		URL: &url.URL{
			Path: proxyOpts.Path,
		},
		Method: strings.ToUpper(reqInfo.Verb),
	})
	proxyReqInfo.Verb = reqInfo.Verb

	return &proxyHandler{
		clusterName:      id,
		resource:         proxyReqInfo.Resource,
		verb:             proxyReqInfo.Verb,
		path:             proxyOpts.Path,
		startConnect:     startTime,
		responder:        responder,
		clusterExtension: clusterExtension,
	}, nil
}

type proxyHandler struct {
	clusterName      string
	resource         string
	verb             string
	path             string
	startConnect     time.Time
	responder        rest.Responder
	clusterExtension *cluster.ClusterExtension
}

func (p *proxyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	setErr := func(err error) {
		klog.ErrorS(err, "request failed", "clusterName", p.clusterName, "verb", p.verb, "resource", p.resource, "userAgent", request.UserAgent())
		responsewriters.InternalError(writer, request, err)
	}

	cluster := p.clusterExtension
	cfg, err := NewConfigFromCluster(cluster)
	if err != nil {
		setErr(errors.Wrapf(err, "failed to create cluster proxy client config %s", cluster.Name))
		return
	}

	userAgent := request.UserAgent()
	if userAgent != "" {
		cfg.UserAgent = userAgent
	}
	impersonateUser(cfg, request)

	urlAddr, err := GetEndpointURL(cluster)
	if err != nil {
		setErr(errors.Wrapf(err, "failed to parsing endpoint for cluster %s", cluster.Name))
		return
	}

	location := &url.URL{
		Scheme:   urlAddr.Scheme,
		Path:     filepath.Join(urlAddr.Path, p.path),
		Host:     urlAddr.Host,
		RawQuery: request.URL.RawQuery,
	}

	rt, err := restclient.TransportFor(cfg)
	if err != nil {
		setErr(errors.Wrapf(err, "failed to create round tripper from rest config"))
		return
	}

	handler := proxy.NewUpgradeAwareHandler(location, rt, false, false, proxy.NewErrorResponder(p.responder))
	handler.UseLocationHost = true
	proxyutil.WithLogs(handler).ServeHTTP(writer, request)
}

func impersonateUser(cfg *restclient.Config, req *http.Request) {
	user, _ := request.UserFrom(req.Context())
	isSA := false
	for _, g := range user.GetGroups() {
		if g == "system:serviceaccounts" {
			isSA = true
		}
	}
	// only impersonate serviceaccounts to compatible with old certificates
	if isSA {
		cfg.Impersonate = restclient.ImpersonationConfig{
			UserName: user.GetName(),
			Groups:   user.GetGroups(),
			Extra:    user.GetExtra(),
		}
	}
}

func GetEndpointURL(c *cluster.ClusterExtension) (*url.URL, error) {
	if c.Spec.Access.Credential == nil {
		return nil, fmt.Errorf("proxying cluster %s not support due to lacking credentials", c.Name)
	}

	urlAddr, err := url.Parse(c.Spec.Access.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parsing url from cluster %s invalid value %s", c.Name, c.Spec.Access.Endpoint)
	}
	return urlAddr, nil
}

func NewConfigFromCluster(c *cluster.ClusterExtension) (*restclient.Config, error) {
	cfg := &restclient.Config{}
	cfg.Host = c.Spec.Access.Endpoint
	cfg.CAData = c.Spec.Access.CABundle
	if c.Spec.Access.Insecure != nil && *c.Spec.Access.Insecure {
		cfg.TLSClientConfig = restclient.TLSClientConfig{Insecure: true}
	}
	switch c.Spec.Access.Credential.Type {
	case cluster.CredentialTypeServiceAccountToken:
		cfg.BearerToken = c.Spec.Access.Credential.ServiceAccountToken
	case cluster.CredentialTypeX509Certificate:
		cfg.CertData = c.Spec.Access.Credential.X509.Certificate
		cfg.KeyData = c.Spec.Access.Credential.X509.PrivateKey
	case cluster.CredentialTypeUnifiedIdentity:
		cfg.CAData = c.Spec.Access.CABundle
		cfg.CertFile = "/var/run/secrets/kubernetes.io/serviceaccount/app.crt"
		cfg.KeyFile = "/var/run/secrets/kubernetes.io/serviceaccount/app.key"
	}
	u, err := url.Parse(c.Spec.Access.Endpoint)
	if err != nil {
		return nil, err
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}
	cfg.ServerName = host // apiserver may listen on SNI cert
	return cfg, nil
}
