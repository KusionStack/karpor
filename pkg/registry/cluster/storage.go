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

	"github.com/KusionStack/karbour/pkg/apis/cluster"

	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/KusionStack/karbour/pkg/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
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

// Destroy cleans up resources on shutdown.
func (r *StatusREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.Store.Get(ctx, name, options)
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

var _ registry.RESTStorageProvider = &RESTStorageProvider{}

type RESTStorageProvider struct{}

func (p RESTStorageProvider) GroupName() string {
	return cluster.GroupName
}

func (p RESTStorageProvider) NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(cluster.GroupName, scheme.Scheme, scheme.ParameterCodec, scheme.Codecs)

	v1beta1 := map[string]rest.Storage{}
	clusterStorage, err := NewREST(restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}

	v1beta1["clusters"] = clusterStorage.Cluster
	v1beta1["clusters/status"] = clusterStorage.Status
	v1beta1["clusters/proxy"] = clusterStorage.Proxy

	apiGroupInfo.VersionedResourcesStorageMap["v1beta1"] = v1beta1
	return apiGroupInfo, nil
}
