/*
Copyright The Karpor Authors.

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
	"github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster"

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

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a
// resource.
func (r *REST) ShortNames() []string {
	return []string{"cl"}
}
