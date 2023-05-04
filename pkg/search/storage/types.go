package storage

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

type Storage interface {
	Get(ctx context.Context, cluster string, obj runtime.Object) error
	Create(ctx context.Context, cluster string, obj runtime.Object) error
	Update(ctx context.Context, cluster string, obj runtime.Object) error
	Delete(ctx context.Context, cluster string, obj runtime.Object) error
}

type Searcher interface {
	Search(ctx context.Context, queries []Query) (*SearchResult, error)
}

type Resource struct {
	Cluster    string                 `json:"cluster"`
	Namespace  string                 `json:"namespace"`
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Name       string                 `json:"name"`
	Object     map[string]interface{} `json:"object"`
}

type SearchResult struct {
	Total     int
	Resources []*Resource
}
