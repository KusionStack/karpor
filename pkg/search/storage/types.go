package storage

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

type Storage interface {
	Get(ctx context.Context, cluster, namespace, apiVerson, kind, name string) (runtime.Object, error)
	List(ctx context.Context, options *ListOptions) (runtime.Object, error)
	Delete(ctx context.Context, cluster, namespace, apiVerson, kind, name string) error
	Create(ctx context.Context, cluster string, obj runtime.Object) error
	Update(ctx context.Context, cluster string, obj runtime.Object) error
}

type ListOptions struct {
	Clusters    []string
	APIVersions []string
	Groups      []string
	Kinds       []string
	Namespaces  []string
	Names       []string
}
