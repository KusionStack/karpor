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
	"os"

	"github.com/KusionStack/karbour/pkg/apis/search"
	"github.com/KusionStack/karbour/pkg/registry/cluster"
	"github.com/KusionStack/karbour/pkg/search/storage"
	filtersutil "github.com/KusionStack/karbour/pkg/util/filters"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/klog/v2"
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

type REST struct {
	Storage storage.SearchStorage
	Cluster *cluster.Storage
}

type Scope struct {
	Scope meta.RESTScopeName
}

func (s Scope) Name() meta.RESTScopeName {
	return s.Scope
}

func NewREST(searchStorageGetter storage.SearchStorageGetter, clusterOptsGetter generic.RESTOptionsGetter) (rest.Storage, error) {
	searchStorage, err := searchStorageGetter.GetSearchStorage()
	if err != nil {
		return nil, err
	}

	clusterStorage, err := cluster.NewREST(clusterOptsGetter)
	if err != nil {
		return nil, err
	}

	return &REST{
		Storage: searchStorage,
		Cluster: clusterStorage,
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

// Get retrieves the uniresource information from storage. Current supports topology calculation for a single uniresource.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	_, err := r.BuildDynamicClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	rt := &search.UniResourceList{}
	if name == "topology" {
		resource, ok := filtersutil.ResourceDetailFrom(ctx)
		if !ok {
			return nil, fmt.Errorf("name, namespace, cluster, apiVersion and kind are used to locate a unique resource so they can't be empty")
		}
		queryString := fmt.Sprintf("%s where name = '%s' AND namespace = '%s' AND cluster = '%s' AND apiVersion = '%s' AND kind = '%s'", SQLQueryDefault, resource.Name, resource.Namespace, resource.Cluster, resource.APIVersion, resource.Kind)
		// TODO: Should we enforce all fields to be present? Or do we allow topology graph for multiple (fuzzy search) resources at a time?
		// if resource.Namespace != "" {
		// 	queryString += fmt.Sprintf(" AND namespace = '%s'", resource.Namespace)
		// }
		// ...

		klog.Infof("Query string: %s", queryString)

		rg, _ := BuildResourceRelationshipGraph()
		res, err := r.Storage.Search(ctx, queryString, storage.SQLPatternType)
		if err != nil {
			return nil, err
		}

		ResourceGraphNodeHash := func(rgn ResourceGraphNode) string {
			return rgn.Group + "/" + rgn.Version + "." + rgn.Kind + ":" + rgn.Namespace + "." + rgn.Name
		}
		g := graph.New(ResourceGraphNodeHash, graph.Directed(), graph.PreventCycles())
		for _, resource := range res.Resources {
			unObj := &unstructured.Unstructured{}
			unObj.SetUnstructuredContent(resource.Object)
			g, err = r.GetResourceRelationship(ctx, *unObj, rg, g)
			if err != nil {
				return rt, err
			}
			rt.Items = append(rt.Items, unObj)
		}
		// Draw graph
		file, _ := os.Create("./resource.gv")
		_ = draw.DOT(g, file)

		// am, _ := g.AdjacencyMap()
		// spew.Dump(am)

		return rt, nil
	} else {
		return nil, fmt.Errorf("only support getting topology for uniresource at the moment")
	}
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

func (r *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: add real logic of convert to table when the storage layer is implemented
	return rest.NewDefaultTableConvertor(search.Resource("uniresources")).ConvertToTable(ctx, object, tableOptions)
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"ur"}
}
