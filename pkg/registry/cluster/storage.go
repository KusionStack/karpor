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
	"os"

	"github.com/KusionStack/karbour/pkg/apis/cluster"
	"github.com/KusionStack/karbour/pkg/registry/search/relationship"
	"github.com/dominikbraun/graph/draw"
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
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

type StatusREST struct {
	Store *genericregistry.Store
}

// New returns empty Cluster object.
func (r *StatusREST) New() runtime.Object {
	return &cluster.Cluster{}
}

func (r *StatusREST) NewList() runtime.Object {
	return &cluster.ClusterList{}
}

// Destroy cleans up resources on shutdown.
func (r *StatusREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.Store.Get(ctx, name, options)
}

func (r *StatusREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	rt := &cluster.Cluster{}
	client, err := r.BuildDynamicClient(ctx)
	if err != nil {
		return rt, err
	}
	graph, rg, _ := relationship.BuildRelationshipGraph(ctx, client, true)

	// Draw graph
	file, _ := os.Create("./relationship.gv")
	_ = draw.DOT(graph, file)

	m := make(map[string]cluster.ClusterTopology)
	for _, rgn := range rg.RelationshipNodes {
		rgnMap := rgn.ConvertToMap()
		m[rgn.GetHash()] = cluster.ClusterTopology{
			GroupVersionKind: rgn.GetHash(),
			Count:            rgn.ResourceCount,
			Relationship:     rgnMap,
		}
	}
	rt.Status.Graph = m
	return rt, nil
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.Store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// GetResetFields implements rest.ResetFieldsStrategy
func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.Store.GetResetFields()
}

func (r *StatusREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.Store.ConvertToTable(ctx, object, tableOptions)
}

// BuildDynamicClient returns a dynamic client based on the cluster name in the request
func (r *StatusREST) BuildDynamicClient(ctx context.Context) (*dynamic.DynamicClient, error) {
	// Extract the cluster name from context
	info, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve request info from context")
	}
	klog.Infof("Getting topology for cluster %s...", info.Name)

	// Locate the cluster resource and build config with it
	obj, err := r.Store.Get(ctx, info.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clusterFromContext := obj.(*cluster.Cluster)
	klog.Infof("Cluster found: %s", clusterFromContext.Name)
	config, err := NewConfigFromCluster(clusterFromContext)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster client config %s", clusterFromContext.Name)
	}

	// Create the dynamic client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
