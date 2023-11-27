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

	cluster "github.com/KusionStack/karbour/pkg/apis/cluster"
	"github.com/KusionStack/karbour/pkg/apis/search"
	clusterstorage "github.com/KusionStack/karbour/pkg/registry/cluster"
	"github.com/KusionStack/karbour/pkg/search/storage"
	filtersutil "github.com/KusionStack/karbour/pkg/util/filters"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
)

type TopologyREST struct {
	Storage storage.SearchStorage
	Cluster *clusterstorage.Storage
}

// New returns an empty cluster proxy subresource.
func (r *TopologyREST) New() runtime.Object {
	return &search.UniresourceTopology{}
}

func (r *TopologyREST) NewList() runtime.Object {
	return &search.UniresourceTopologyList{}
}

// Destroy cleans up resources on shutdown.
func (r *TopologyREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

func (r *TopologyREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: add real logic of convert to table when the storage layer is implemented
	return rest.NewDefaultTableConvertor(search.Resource("uniresources")).ConvertToTable(ctx, object, tableOptions)
}

func (r *TopologyREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	rt := &search.UniResourceList{}
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
}

// BuildDynamicClient returns a dynamic client based on the cluster name in the request
func (r *TopologyREST) BuildDynamicClient(ctx context.Context) (*dynamic.DynamicClient, error) {
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

// GetResourceRelationship returns a full graph that contains all the resources that are related to obj
func (r *TopologyREST) GetResourceRelationship(ctx context.Context, obj unstructured.Unstructured, relationshipGraph graph.Graph[string, RelationshipGraphNode], resourceGraph graph.Graph[string, ResourceGraphNode]) (graph.Graph[string, ResourceGraphNode], error) {
	namespace := obj.GetNamespace()
	objName := obj.GetName()
	gv, _ := schema.ParseGroupVersion(obj.GetAPIVersion())
	objResourceNode := ResourceGraphNode{
		Group:     gv.Group,
		Version:   gv.Version,
		Kind:      obj.GetKind(),
		Name:      objName,
		Namespace: namespace,
	}
	resourceGraph.AddVertex(objResourceNode)
	client, err := r.BuildDynamicClient(ctx)
	if err != nil {
		return resourceGraph, err
	}

	objGVKOnGraph, _ := FindNodeOnGraph(relationshipGraph, gv.Group, gv.Version, obj.GetKind())
	// TODO: process error
	// Recursively find parents
	for _, objParent := range objGVKOnGraph.Parent {
		resourceGraph, err = GetParents(ctx, client, obj, objParent, namespace, objName, objResourceNode, relationshipGraph, resourceGraph)
		if err != nil {
			return nil, err
		}
	}

	// Recursively find children
	for _, objChild := range objGVKOnGraph.Children {
		resourceGraph, err = GetChildren(ctx, client, obj, objChild, namespace, objName, objResourceNode, relationshipGraph, resourceGraph)
		if err != nil {
			return nil, err
		}
	}

	return resourceGraph, nil
}
