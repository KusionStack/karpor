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

package insight

import (
	"context"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/infra/multicluster"
	"github.com/KusionStack/karbour/pkg/infra/topology"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	topologyutil "github.com/KusionStack/karbour/pkg/util/topology"
	"github.com/dominikbraun/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GetTopologyForCluster returns a map that describes topology for a given cluster
func (i *InsightManager) GetTopologyForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	locator := core.Locator{
		Cluster: name,
	}
	if topologyData, exist := i.clusterTopologyCache.Get(locator); exist {
		log.Info("Cache hit for cluster topology", "locator", locator)
		return topologyData, nil
	}

	// Build relationship graph based on GVK
	_, rg, _ := topology.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Count resources in all namespaces
	log.Info("Retrieving topology for cluster", "clusterName", name)
	rg, err := rg.CountRelationshipGraph(ctx, client.DynamicClient, client.ClientSet.DiscoveryClient, "")
	if err != nil {
		return nil, err
	}

	clusterTopologyMap := i.ConvertGraphToMap(rg, locator)
	i.clusterTopologyCache.Set(locator, clusterTopologyMap)
	log.Info("Added to cluster topology cache for locator", "locator", locator)

	return clusterTopologyMap, nil
}

// GetTopologyForResource returns a map that describes topology for a given cluster
func (i *InsightManager) GetTopologyForResource(ctx context.Context, client *multicluster.MultiClusterClient, locator *core.Locator) (map[string]ResourceTopology, error) {
	log := ctxutil.GetLogger(ctx)

	if topologyData, exist := i.resourceTopologyCache.Get(*locator); exist {
		log.Info("Cache hit for resource topology", "locator", locator)
		return topologyData, nil
	}

	log.Info("Cache miss for locator", "locator", locator)
	// Build relationship graph based on GVK
	rg, _, err := topology.BuildRelationshipGraph(ctx, client.DynamicClient)
	if err != nil {
		return nil, err
	}
	log.Info("Retrieving topology for resource", "resourceName", locator.Name)

	ResourceGraphNodeHash := func(rgn topology.ResourceGraphNode) string {
		return rgn.Group + "/" + rgn.Version + "." + rgn.Kind + ":" + rgn.Namespace + "." + rgn.Name
	}
	g := graph.New(ResourceGraphNodeHash, graph.Directed(), graph.PreventCycles())

	// Get target resource
	resourceGVR, err := topologyutil.GetGVRFromGVK(locator.APIVersion, locator.Kind)
	if err != nil {
		return nil, err
	}
	resObj, err := client.DynamicClient.Resource(resourceGVR).Namespace(locator.Namespace).Get(ctx, locator.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	unObj := &unstructured.Unstructured{}
	unObj.SetUnstructuredContent(resObj.Object)

	// Build resource graph for target resource
	g, err = i.GetResourceRelationship(ctx, client, *unObj, rg, g)
	if err != nil {
		return nil, err
	}

	topologyMap := i.ConvertResourceGraphToMap(g, *locator)
	i.resourceTopologyCache.Set(*locator, topologyMap)
	log.Info("Added to resource topology cache for locator", "locator", locator)

	return topologyMap, nil
}

// GetResourceRelationship returns a full graph that contains all the resources that are related to obj
func (i *InsightManager) GetResourceRelationship(ctx context.Context, client *multicluster.MultiClusterClient, obj unstructured.Unstructured, relationshipGraph graph.Graph[string, topology.RelationshipGraphNode], resourceGraph graph.Graph[string, topology.ResourceGraphNode]) (graph.Graph[string, topology.ResourceGraphNode], error) {
	var err error
	namespace := obj.GetNamespace()
	objName := obj.GetName()
	gv, _ := schema.ParseGroupVersion(obj.GetAPIVersion())
	objResourceNode := topology.ResourceGraphNode{
		Group:     gv.Group,
		Version:   gv.Version,
		Kind:      obj.GetKind(),
		Name:      objName,
		Namespace: namespace,
	}
	resourceGraph.AddVertex(objResourceNode)

	objGVKOnGraph, err := topology.FindNodeOnGraph(relationshipGraph, gv.Group, gv.Version, obj.GetKind())
	// When obj GVK is not found on relationship graph, return an empty graph with no error
	if err != nil {
		return nil, nil //nolint:nilerr
	}

	// Recursively find parents
	for _, objParent := range objGVKOnGraph.Parent {
		resourceGraph, err = topology.GetParents(ctx, client.DynamicClient, obj, objParent, namespace, objName, objResourceNode, relationshipGraph, resourceGraph)
		if err != nil {
			return nil, err
		}
	}

	// Recursively find children
	for _, objChild := range objGVKOnGraph.Children {
		resourceGraph, err = topology.GetChildren(ctx, client.DynamicClient, obj, objChild, namespace, objName, objResourceNode, relationshipGraph, resourceGraph)
		if err != nil {
			return nil, err
		}
	}

	return resourceGraph, nil
}

func (i *InsightManager) ConvertResourceGraphToMap(g graph.Graph[string, topology.ResourceGraphNode], loc core.Locator) map[string]ResourceTopology {
	rtm := make(map[string]ResourceTopology)
	if g == nil {
		return rtm
	}
	am, _ := g.AdjacencyMap()
	for key, edgeMap := range am {
		childList := make([]string, 0)
		for edgeTarget := range edgeMap {
			childList = append(childList, edgeTarget)
		}
		vertex, _ := g.Vertex(key)
		locator := core.Locator{
			Cluster:    loc.Cluster,
			APIVersion: schema.GroupVersion{Group: vertex.Group, Version: vertex.Version}.String(),
			Kind:       vertex.Kind,
			Namespace:  vertex.Namespace,
			Name:       vertex.Name,
		}
		rtm[key] = ResourceTopology{
			Locator:  locator,
			Children: childList,
		}
	}

	pm, _ := g.PredecessorMap()
	for key, edgeMap := range pm {
		parentList := make([]string, 0)
		for edgeSource := range edgeMap {
			parentList = append(parentList, edgeSource)
		}
		if node, ok := rtm[key]; ok {
			node.Parents = parentList
			rtm[key] = node
		}
	}
	return rtm
}

// GetTopologyForClusterNamespace returns a map that describes topology for a given namespace in a given cluster
func (i *InsightManager) GetTopologyForClusterNamespace(ctx context.Context, client *multicluster.MultiClusterClient, cluster, namespace string) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	locator := core.Locator{
		Cluster:   cluster,
		Namespace: namespace,
	}
	if topologyData, exist := i.clusterTopologyCache.Get(locator); exist {
		log.Info("Cache hit for cluster topology", "locator", locator)
		return topologyData, nil
	}

	// Build relationship graph based on GVK
	_, rg, _ := topology.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Only count resources that belong to a specific namespace
	log.Info("Retrieving topology", "namespace", namespace, "cluster", cluster)
	rg, err := rg.CountRelationshipGraph(ctx, client.DynamicClient, client.ClientSet.DiscoveryClient, namespace)
	if err != nil {
		return nil, err
	}

	namespaceTopologyMap := i.ConvertGraphToMap(rg, locator)
	i.clusterTopologyCache.Set(locator, namespaceTopologyMap)
	log.Info("Added to cluster topology cache for locator", "locator", locator)

	return namespaceTopologyMap, nil
}

// ConvertGraphToMap returns a map[string]ClusterTopology for a given relationship.RelationshipGraph
func (i *InsightManager) ConvertGraphToMap(rg *topology.RelationshipGraph, loc core.Locator) map[string]ClusterTopology {
	ctm := make(map[string]ClusterTopology)
	for _, rgn := range rg.RelationshipNodes {
		rgnMap := rgn.ConvertToMap()
		locator := core.Locator{
			Cluster:    loc.Cluster,
			APIVersion: schema.GroupVersion{Group: rgn.Group, Version: rgn.Version}.String(),
			Kind:       rgn.Kind,
			Namespace:  loc.Namespace,
		}
		ctm[rgn.GetHash()] = ClusterTopology{
			Locator:      locator,
			Count:        rgn.ResourceCount,
			Relationship: rgnMap,
		}
	}
	return ctm
}
