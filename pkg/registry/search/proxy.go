package search

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KusionStack/karbour/pkg/apis/search"
	"github.com/KusionStack/karbour/pkg/registry/search/esserver"
	"github.com/elastic/go-elasticsearch/v8"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

var proxyMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

var (
	_ rest.Connecter = &ProxyREST{}
	_ rest.Storage   = &ProxyREST{}
	_ rest.Scoper    = &ProxyREST{}
)

type ProxyREST struct{}

// NamespaceScoped returns false because Resources is not namespaced
func (r *ProxyREST) NamespaceScoped() bool {
	return false
}

func (r *ProxyREST) Destroy() {
}

// New returns an empty search proxy subresource.
func (r *ProxyREST) New() runtime.Object {
	return &search.Search{}
}

func (r *ProxyREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &search.Search{}, false, ""
}

func (r *ProxyREST) ConnectMethods() []string {
	return proxyMethods
}

func (r *ProxyREST) Connect(ctx context.Context, id string, options runtime.Object, responder rest.Responder) (http.Handler, error) {
	opt, ok := options.(*search.Search)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", opt)
	}
	es := esserver.NewElasticServerOrDie(elasticsearch.Config{
		Addresses: []string{"http://100.88.101.58:9200"},
	})
	return es, nil
}
