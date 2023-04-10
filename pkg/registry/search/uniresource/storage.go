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
	// _ rest.Getter  = &REST{}
	_ rest.Lister = &REST{}
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

// func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
// 	u := search.UniResource{}
// 	u.Data = name
// 	return &u, nil
// }

func (r *REST) NewList() runtime.Object {
	return &search.UniResourceList{}
}

func (r *REST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	u1, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(r1)
	return &search.UniResourceList{
		Items: []search.UniResource{
			{Data: r1},
			{Data: r2},
			{Data: r3},
			{Data: r4},
			{Data: r5},
			{Data: r6},
		},
	}, nil
}

func (r *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return rest.NewDefaultTableConvertor(search.Resource("uniresources")).ConvertToTable(ctx, object, tableOptions)
}

var (
	r1 = `{
    "apiVersion": "v1",
    "kind": "Namespace",
    "metadata": {
        "name": "test1",
        "namespace": ""
    }
}`
	r2 = `{
    "apiVersion": "v1",
    "kind": "Namespace",
    "metadata": {
        "name": "test2",
        "namespace": ""
    }
}`
	r3 = `{
    "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
    "kind": "CafeDeployment",
    "metadata": {
        "name": "test1",
        "namespace": "test1"
    }
}`
	r4 = `{
    "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
    "kind": "CafeDeployment",
    "metadata": {
        "name": "test2",
        "namespace": "test2"
    }
}`
	r5 = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "name": "test1",
        "namespace": "test1"
    }
},`
	r6 = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "name": "test2",
        "namespace": "test2"
    }
}`
)
