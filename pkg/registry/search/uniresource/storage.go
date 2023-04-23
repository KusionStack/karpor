package uniresource

import (
	"context"

	"github.com/KusionStack/karbour/pkg/apis/search"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	_ rest.Storage = &REST{}
	_ rest.Scoper  = &REST{}
	_ rest.Lister  = &REST{}
)

type REST struct{}

func NewREST() rest.Storage {
	return &REST{}
}

func (r *REST) New() runtime.Object {
	return &search.UniResource{}
}

func (r *REST) Destroy() {
}

func (r *REST) NamespaceScoped() bool {
	return false
}

func (r *REST) NewList() runtime.Object {
	return &search.UniResourceList{}
}

func (r *REST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	return &search.UniResourceList{}, nil
}

func (r *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return rest.NewDefaultTableConvertor(search.Resource("uniresources")).ConvertToTable(ctx, object, tableOptions)
}
