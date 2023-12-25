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
	"os"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/relationship"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	topologyutil "github.com/KusionStack/karbour/pkg/util/topology"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
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
	graph, rg, _ := relationship.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Count resources in all namespaces
	log.Info("Retrieving topology for cluster", "clusterName", name)
	rg, err := rg.CountRelationshipGraph(ctx, client.DynamicClient, client.ClientSet.DiscoveryClient, "")
	if err != nil {
		return nil, err
	}

	clusterTopologyMap := i.ConvertGraphToMap(rg)
	i.clusterTopologyCache.Set(locator, clusterTopologyMap)
	log.Info("Added to cluster topology cache for locator", "locator", locator)

	// Draw graph
	// TODO: This is drawn on the server side, not needed eventually
	file, _ := os.Create("./relationship.gv")
	_ = draw.DOT(graph, file)

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
	rg, _, err := relationship.BuildRelationshipGraph(ctx, client.DynamicClient)
	if err != nil {
		return nil, err
	}
	log.Info("Retrieving topology for resource", "resourceName", locator.Name)

	ResourceGraphNodeHash := func(rgn relationship.ResourceGraphNode) string {
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

	topologyMap := i.ConvertResourceGraphToMap(g)
	i.resourceTopologyCache.Set(*locator, topologyMap)
	log.Info("Added to resource topology cache for locator", "locator", locator)

	// Draw graph
	// TODO: This is drawn on the server side, not needed eventually
	file, _ := os.Create("./resource.gv")
	_ = draw.DOT(g, file)

	return topologyMap, nil
}

// GetResourceRelationship returns a full graph that contains all the resources that are related to obj
func (i *InsightManager) GetResourceRelationship(ctx context.Context, client *multicluster.MultiClusterClient, obj unstructured.Unstructured, relationshipGraph graph.Graph[string, relationship.RelationshipGraphNode], resourceGraph graph.Graph[string, relationship.ResourceGraphNode]) (graph.Graph[string, relationship.ResourceGraphNode], error) {
	var err error
	namespace := obj.GetNamespace()
	objName := obj.GetName()
	gv, _ := schema.ParseGroupVersion(obj.GetAPIVersion())
	objResourceNode := relationship.ResourceGraphNode{
		Group:     gv.Group,
		Version:   gv.Version,
		Kind:      obj.GetKind(),
		Name:      objName,
		Namespace: namespace,
	}
	resourceGraph.AddVertex(objResourceNode)

	objGVKOnGraph, _ := relationship.FindNodeOnGraph(relationshipGraph, gv.Group, gv.Version, obj.GetKind())
	// TODO: process error
	// Recursively find parents
	for _, objParent := range objGVKOnGraph.Parent {
		resourceGraph, err = relationship.GetParents(ctx, client.DynamicClient, obj, objParent, namespace, objName, objResourceNode, relationshipGraph, resourceGraph)
		if err != nil {
			return nil, err
		}
	}

	// Recursively find children
	for _, objChild := range objGVKOnGraph.Children {
		resourceGraph, err = relationship.GetChildren(ctx, client.DynamicClient, obj, objChild, namespace, objName, objResourceNode, relationshipGraph, resourceGraph)
		if err != nil {
			return nil, err
		}
	}

	return resourceGraph, nil
}

func (i *InsightManager) ConvertResourceGraphToMap(g graph.Graph[string, relationship.ResourceGraphNode]) map[string]ResourceTopology {
	am, _ := g.AdjacencyMap()
	rtm := make(map[string]ResourceTopology)
	for key, edgeMap := range am {
		childList := make([]string, 0)
		for edgeTarget := range edgeMap {
			childList = append(childList, edgeTarget)
		}
		rtm[key] = ResourceTopology{
			Identifier: key,
			Children:   childList,
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
	graph, rg, _ := relationship.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Only count resources that belong to a specific namespace
	log.Info("Retrieving topology", "namespace", namespace, "cluster", cluster)
	rg, err := rg.CountRelationshipGraph(ctx, client.DynamicClient, client.ClientSet.DiscoveryClient, namespace)
	if err != nil {
		return nil, err
	}

	namespaceTopologyMap := i.ConvertGraphToMap(rg)
	i.clusterTopologyCache.Set(locator, namespaceTopologyMap)
	log.Info("Added to cluster topology cache for locator", "locator", locator)

	// Draw graph
	// TODO: This is drawn on the server side, not needed eventually
	file, _ := os.Create("./relationship.gv")
	_ = draw.DOT(graph, file)

	return namespaceTopologyMap, nil
}

// ConvertGraphToMap returns a map[string]ClusterTopology for a given relationship.RelationshipGraph
func (i *InsightManager) ConvertGraphToMap(rg *relationship.RelationshipGraph) map[string]ClusterTopology {
	ctm := make(map[string]ClusterTopology)
	for _, rgn := range rg.RelationshipNodes {
		rgnMap := rgn.ConvertToMap()
		ctm[rgn.GetHash()] = ClusterTopology{
			GroupVersionKind: rgn.GetHash(),
			Count:            rgn.ResourceCount,
			Relationship:     rgnMap,
		}
	}
	return ctm
}
