// Copyright The Karbour Authors.
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

package uniresource

import (
	"context"
	"fmt"

	"github.com/KusionStack/karbour/pkg/apis/search"
	clusterstorage "github.com/KusionStack/karbour/pkg/registry/cluster"
	"github.com/KusionStack/karbour/pkg/search/storage"
	filtersutil "github.com/KusionStack/karbour/pkg/util/filters"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"

	"k8s.io/apiserver/pkg/registry/rest"

	"github.com/pkg/errors"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"

	cluster "github.com/KusionStack/karbour/pkg/apis/cluster"
)

var (
	_ rest.Storage            = &REST{}
	_ rest.Scoper             = &REST{}
	_ rest.Lister             = &REST{}
	_ rest.ShortNamesProvider = &REST{}
)

const (
	SQLQueryDefault = "select * from resources"
)

type Storage struct {
	Uniresource *REST
	Topology    *TopologyREST
	YAML        *YAMLREST
}

type REST struct {
	Storage storage.SearchStorage
	Cluster *clusterstorage.Storage
}

func NewREST(searchStorageGetter storage.SearchStorageGetter, clusterOptsGetter generic.RESTOptionsGetter) (*Storage, error) {
	searchStorage, err := searchStorageGetter.GetSearchStorage()
	if err != nil {
		return nil, err
	}

	clusterStorage, err := clusterstorage.NewREST(clusterOptsGetter)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Uniresource: &REST{
			Storage: searchStorage,
			Cluster: clusterStorage,
		},
		Topology: &TopologyREST{
			Storage: searchStorage,
			Cluster: clusterStorage,
		},
		YAML: &YAMLREST{
			Storage: searchStorage,
			Cluster: clusterStorage,
		},
	}, nil
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
	queryString, ok := filtersutil.SearchQueryFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("query can't be empty")
	}

	patternType, ok := filtersutil.PatternTypeFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("pattern type can't be empty")
	}

	res, err := r.Storage.Search(ctx, queryString, patternType)
	if err != nil {
		return nil, err
	}

	rt := &search.UniResourceList{}
	for _, resource := range res.Resources {
		unObj := &unstructured.Unstructured{}
		unObj.SetUnstructuredContent(resource.Object)
		rt.Items = append(rt.Items, unObj)
	}
	return rt, nil
}

// BuildDynamicClient returns a dynamic client based on the cluster name in the request
func (r *REST) BuildDynamicClient(ctx context.Context) (*dynamic.DynamicClient, error) {
	// Extract the cluster name from context
	resourceDetail, ok := filtersutil.ResourceDetailFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("name, namespace, cluster, apiVersion and kind are used to locate a unique resource so they can't be empty")
	}

	// Locate the cluster resource and build config with it
	obj, err := r.Cluster.Status.Store.Get(ctx, resourceDetail.Cluster, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clusterFromContext := obj.(*cluster.Cluster)
	klog.Infof("Cluster found: %s", clusterFromContext.Name)
	config, err := clusterstorage.NewConfigFromCluster(clusterFromContext)
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

func (r *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: add real logic of convert to table when the storage layer is implemented
	return rest.NewDefaultTableConvertor(search.Resource("uniresources")).ConvertToTable(ctx, object, tableOptions)
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"ur"}
}
