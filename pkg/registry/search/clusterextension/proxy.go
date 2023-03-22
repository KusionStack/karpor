package searchextension

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"code.alipay.com/multi-cluster/karbour/pkg/apis/search"
	"code.alipay.com/multi-cluster/karbour/pkg/identity"
	proxyutil "code.alipay.com/multi-cluster/karbour/pkg/util/proxy"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	restclient "k8s.io/client-go/rest"
)

var proxyMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

var _ rest.Connecter = &ProxyREST{}

type ProxyREST struct {
	Store *genericregistry.Store
}

func (r *ProxyREST) Destroy() {
}

// New returns an empty search proxy subresource.
func (r *ProxyREST) New() runtime.Object {
	return &search.SearchExtensionProxyOptions{}
}

func (r *ProxyREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &search.SearchExtensionProxyOptions{}, true, "path"
}

func (r *ProxyREST) ConnectMethods() []string {
	return proxyMethods
}

func (r *ProxyREST) Connect(ctx context.Context, id string, options runtime.Object, responder rest.Responder) (http.Handler, error) {
	proxyOpts, ok := options.(*search.SearchExtensionProxyOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", options)
	}

	obj, err := r.Store.Get(ctx, id, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	searchExtension := obj.(*search.SearchExtension)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		location, transport, err := resourceLocation(searchExtension, proxyOpts.Path, r)
		if err != nil {
			responsewriters.InternalError(w, r, err)
			return
		}

		proxyHandler := proxy.NewUpgradeAwareHandler(location, transport, false, false, proxy.NewErrorResponder(responder))
		proxyHandler.UseLocationHost = true
		handler := proxyutil.WithLogs(proxyHandler)
		handler.ServeHTTP(w, r)
	}), nil
}

func resourceLocation(searchExtension *search.SearchExtension, path string, request *http.Request) (location *url.URL, transport http.RoundTripper, err error) {
	location, err = getEndpointURL(searchExtension)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to parsing endpoint for search %s", searchExtension.Name)
	}
	location.Path = path
	location.RawQuery = request.URL.RawQuery

	cfg, err := NewConfigFromSearch(searchExtension)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create search proxy client config %s", searchExtension.Name)
	}

	userAgent := request.UserAgent()
	if userAgent != "" {
		cfg.UserAgent = userAgent
	}

	impersonateUser(cfg, request)

	transport, err = restclient.TransportFor(cfg)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create round tripper from rest config")
	}
	return location, transport, nil
}

func NewConfigFromSearch(c *search.SearchExtension) (*restclient.Config, error) {
	cfg := &restclient.Config{}
	cfg.Host = c.Spec.Access.Endpoint
	cfg.CAData = c.Spec.Access.CABundle
	if c.Spec.Access.Insecure != nil && *c.Spec.Access.Insecure {
		cfg.TLSClientConfig = restclient.TLSClientConfig{Insecure: true}
	}
	switch c.Spec.Access.Credential.Type {
	case search.CredentialTypeServiceAccountToken:
		cfg.BearerToken = c.Spec.Access.Credential.ServiceAccountToken
	case search.CredentialTypeX509Certificate:
		cfg.CertData = c.Spec.Access.Credential.X509.Certificate
		cfg.KeyData = c.Spec.Access.Credential.X509.PrivateKey
	case search.CredentialTypeUnifiedIdentity:
		cfg.CAData = c.Spec.Access.CABundle
		cfg.CertFile = identity.GetCertFile()
		cfg.KeyFile = identity.GetKeyFile()
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

func getEndpointURL(c *search.SearchExtension) (*url.URL, error) {
	if c.Spec.Access.Credential == nil {
		return nil, fmt.Errorf("proxying search %s not support due to lacking credentials", c.Name)
	}

	urlAddr, err := url.Parse(c.Spec.Access.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parsing url from search %s invalid value %s", c.Name, c.Spec.Access.Endpoint)
	}
	return urlAddr, nil
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
