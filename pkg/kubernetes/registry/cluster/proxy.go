// Copyright The Karpor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"kusionstack.io/karpor/pkg/kubernetes/apis/cluster"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
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

// New returns an empty cluster proxy subresource.
func (r *ProxyREST) New() runtime.Object {
	return &cluster.ClusterProxyOptions{}
}

func (r *ProxyREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &cluster.ClusterProxyOptions{}, true, "path"
}

func (r *ProxyREST) ConnectMethods() []string {
	return proxyMethods
}

func (r *ProxyREST) Connect(
	ctx context.Context,
	id string,
	options runtime.Object,
	responder rest.Responder,
) (http.Handler, error) {
	proxyOpts, ok := options.(*cluster.ClusterProxyOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", options)
	}

	obj, err := r.Store.Get(ctx, id, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clusterExtension := obj.(*cluster.Cluster)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		location, transport, err := resourceLocation(clusterExtension, proxyOpts.Path, r)
		if err != nil {
			responsewriters.InternalError(w, r, err)
			return
		}

		proxyHandler := proxy.NewUpgradeAwareHandler(
			location,
			transport,
			false,
			false,
			proxy.NewErrorResponder(responder),
		)
		proxyHandler.UseLocationHost = true
		proxyHandler.ServeHTTP(w, r)
	}), nil
}

func resourceLocation(
	clusterExtension *cluster.Cluster,
	path string,
	request *http.Request,
) (location *url.URL, transport http.RoundTripper, err error) {
	location, err = getEndpointURL(clusterExtension)
	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"failed to parsing endpoint for cluster %s",
			clusterExtension.Name,
		)
	}
	location.Path = path
	location.RawQuery = request.URL.RawQuery

	cfg, err := NewConfigFromCluster(clusterExtension)
	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"failed to create cluster proxy client config %s",
			clusterExtension.Name,
		)
	}

	userAgent := request.UserAgent()
	if userAgent != "" {
		cfg.UserAgent = userAgent
	}

	transport, err = restclient.TransportFor(cfg)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create round tripper from rest config")
	}
	return location, transport, nil
}

func NewConfigFromCluster(c *cluster.Cluster) (*restclient.Config, error) {
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

func getEndpointURL(c *cluster.Cluster) (*url.URL, error) {
	if c.Spec.Access.Credential == nil {
		return nil, fmt.Errorf("proxying cluster %s not support due to lacking credentials", c.Name)
	}

	urlAddr, err := url.Parse(c.Spec.Access.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to parsing url from cluster %s invalid value %s",
			c.Name,
			c.Spec.Access.Endpoint,
		)
	}
	return urlAddr, nil
}
