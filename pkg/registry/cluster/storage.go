/*
Copyright The Karbour Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"context"
	"fmt"

	"github.com/KusionStack/karbour/pkg/apis/cluster"
	clustermgr "github.com/KusionStack/karbour/pkg/core/manager/cluster"

	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	_ rest.Lister = &REST{}
	_ rest.Getter = &REST{}
)

type Storage struct {
	Cluster *REST
	Status  *StatusREST
	Proxy   *ProxyREST
}

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(optsGetter generic.RESTOptionsGetter) (*Storage, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &cluster.Cluster{} },
		NewListFunc:              func() runtime.Object { return &cluster.ClusterList{} },
		DefaultQualifiedResource: cluster.Resource("clusters"),
		PredicateFunc:            MatchCluster,
		// SingularQualifiedResource: cluster.Resource("cluster"),

		CreateStrategy: Strategy,
		UpdateStrategy: Strategy,
		DeleteStrategy: Strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: rest.NewDefaultTableConvertor(cluster.Resource("clusters")),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = StatusStartegy

	return &Storage{
		Cluster: &REST{store},
		Status:  &StatusREST{&statusStore},
		Proxy:   &ProxyREST{store},
	}, nil
}

type REST struct {
	*genericregistry.Store
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"cl"}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *REST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	sanitizedClusterList := &cluster.ClusterList{}
	clusterList, err := r.Store.List(ctx, options)
	if err != nil {
		return sanitizedClusterList, err
	}
	for _, c := range clusterList.(*cluster.ClusterList).Items {
		// Convert to unstructured
		clusterMap, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(&c)
		clusterUnstructured := &unstructured.Unstructured{}
		clusterUnstructured.SetUnstructuredContent(clusterMap)
		// Sanitize credentials
		sanitized, _ := clustermgr.SanitizeUnstructuredCluster(ctx, clusterUnstructured)
		clusterObj := &cluster.Cluster{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(sanitized.Object, clusterObj)
		if err != nil {
			fmt.Printf("Error converting unstructured item to Cluster: %v\n", err)
		}
		sanitizedClusterList.Items = append(sanitizedClusterList.Items, *clusterObj)
	}
	return sanitizedClusterList, nil
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	sanitized := &unstructured.Unstructured{}
	cluster, err := r.Store.Get(ctx, name, options)
	if err != nil {
		return sanitized, nil
	}
	clusterMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cluster)
	if err != nil {
		return sanitized, nil
	}
	clusterUnstructured := &unstructured.Unstructured{}
	clusterUnstructured.SetUnstructuredContent(clusterMap)
	sanitized, _ = clustermgr.SanitizeUnstructuredCluster(ctx, clusterUnstructured)
	return sanitized, nil
}
