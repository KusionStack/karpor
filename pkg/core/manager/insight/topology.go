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

	"github.com/KusionStack/karbour/pkg/core/entity"
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
func (i *InsightManager) GetTopologyForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string, noCache bool) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	resourceGroup := entity.ResourceGroup{
		Cluster: name,
	}

	// If noCache is set to false, attempt to retrieve the result from cache first
	if !noCache {
		if topologyData, exist := i.clusterTopologyCache.Get(resourceGroup); exist {
			log.Info("Cache hit for cluster topology", "resourceGroup", resourceGroup)
			return topologyData, nil
		}
		log.Info("Cache miss for resourceGroup", "resourceGroup", resourceGroup)
	}

	log.Info("Calculating topology for cluster...", "cluster", name)
	// Build relationship graph based on GVK
	_, rg, _ := topology.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Count resources in all namespaces
	log.Info("Retrieving topology for cluster", "clusterName", name)
	rg, err := rg.CountRelationshipGraph(ctx, client.DynamicClient, client.ClientSet.DiscoveryClient, "")
	if err != nil {
		return nil, err
	}

	clusterTopologyMap := i.ConvertGraphToMap(rg, resourceGroup)
	i.clusterTopologyCache.Set(resourceGroup, clusterTopologyMap)
	log.Info("Added to cluster topology cache for resourceGroup", "resourceGroup", resourceGroup)

	return clusterTopologyMap, nil
}

// GetTopologyForResource returns a map that describes topology for a given cluster
func (i *InsightManager) GetTopologyForResource(ctx context.Context, client *multicluster.MultiClusterClient, resourceGroup *entity.ResourceGroup, noCache bool) (map[string]ResourceTopology, error) {
	log := ctxutil.GetLogger(ctx)

	// If noCache is set to false, attempt to retrieve the result from cache first
	if !noCache {
		if topologyData, exist := i.resourceTopologyCache.Get(*resourceGroup); exist {
			log.Info("Cache hit for resource topology", "resourceGroup", resourceGroup)
			return topologyData, nil
		}
		log.Info("Cache miss for resourceGroup", "resourceGroup", resourceGroup)
	}

	log.Info("Calculating topology for resource...", "resourceGroup", resourceGroup)
	// Build relationship graph based on GVK
	rg, _, err := topology.BuildRelationshipGraph(ctx, client.DynamicClient)
	if err != nil {
		return nil, err
	}
	log.Info("Retrieving topology for resource", "resourceName", resourceGroup.Name)

	ResourceGraphNodeHash := func(rgn topology.ResourceGraphNode) string {
		return rgn.Group + "/" + rgn.Version + "." + rgn.Kind + ":" + rgn.Namespace + "." + rgn.Name
	}
	g := graph.New(ResourceGraphNodeHash, graph.Directed(), graph.PreventCycles())

	// Get target resource
	resourceGVR, err := topologyutil.GetGVRFromGVK(resourceGroup.APIVersion, resourceGroup.Kind)
	if err != nil {
		return nil, err
	}
	resObj, err := client.DynamicClient.Resource(resourceGVR).Namespace(resourceGroup.Namespace).Get(ctx, resourceGroup.Name, metav1.GetOptions{})
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

	topologyMap := i.ConvertResourceGraphToMap(g, *resourceGroup)
	i.resourceTopologyCache.Set(*resourceGroup, topologyMap)
	log.Info("Added to resource topology cache for resourceGroup", "resourceGroup", resourceGroup)

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

// ConvertResourceGraphToMap converts a resource graph to a map of ResourceTopology based on the given graph and resourceGroup.
func (i *InsightManager) ConvertResourceGraphToMap(g graph.Graph[string, topology.ResourceGraphNode], resourceGroup entity.ResourceGroup) map[string]ResourceTopology {
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
		resourceGroup := entity.ResourceGroup{
			Cluster:    resourceGroup.Cluster,
			APIVersion: schema.GroupVersion{Group: vertex.Group, Version: vertex.Version}.String(),
			Kind:       vertex.Kind,
			Namespace:  vertex.Namespace,
			Name:       vertex.Name,
		}
		rtm[key] = ResourceTopology{
			ResourceGroup: resourceGroup,
			Children:      childList,
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
func (i *InsightManager) GetTopologyForClusterNamespace(ctx context.Context, client *multicluster.MultiClusterClient, cluster, namespace string, noCache bool) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	resourceGroup := entity.ResourceGroup{
		Cluster:   cluster,
		Namespace: namespace,
	}

	if !noCache {
		if topologyData, exist := i.clusterTopologyCache.Get(resourceGroup); exist {
			log.Info("Cache hit for cluster topology", "resourceGroup", resourceGroup)
			return topologyData, nil
		}
		log.Info("Cache miss for resourceGroup", "resourceGroup", resourceGroup)
	}

	log.Info("Calculating topology for namespace...", "cluster", cluster, "namespace", namespace)
	// Build relationship graph based on GVK
	_, rg, _ := topology.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Only count resources that belong to a specific namespace
	log.Info("Retrieving topology", "namespace", namespace, "cluster", cluster)
	rg, err := rg.CountRelationshipGraph(ctx, client.DynamicClient, client.ClientSet.DiscoveryClient, namespace)
	if err != nil {
		return nil, err
	}

	namespaceTopologyMap := i.ConvertGraphToMap(rg, resourceGroup)
	i.clusterTopologyCache.Set(resourceGroup, namespaceTopologyMap)
	log.Info("Added to cluster topology cache for resourceGroup", "resourceGroup", resourceGroup)

	return namespaceTopologyMap, nil
}

// GetTopologyForCustomResourceGroup returns a map that describes topology for custom resource group
func (i *InsightManager) GetTopologyForCustomResourceGroup(ctx context.Context, client *multicluster.MultiClusterClient, customResourceGroup string, clusters []string, noCache bool) (map[string]map[string]ClusterTopology, error) {
	result := map[string]map[string]ClusterTopology{}
	for _, cluster := range clusters {
		m, err := i.GetTopologyForCustomResourceGroupSingleCluster(ctx, client, customResourceGroup, cluster, noCache)
		if err != nil {
			return nil, err
		}
		result[cluster] = m
	}
	return result, nil
}

// GetTopologyForCustomResourceGroupSingleCluster returns a map that describes topology for single cluster custom resource group
func (i *InsightManager) GetTopologyForCustomResourceGroupSingleCluster(ctx context.Context, client *multicluster.MultiClusterClient, customResourceGroup string, cluster string, noCache bool) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	resourceGroup := entity.ResourceGroup{CustomResourceGroup: customResourceGroup, Cluster: cluster}

	// If noCache is set to false, attempt to retrieve the result from cache first
	if !noCache {
		if topologyData, exist := i.clusterTopologyCache.Get(resourceGroup); exist {
			log.Info("Cache hit for cluster topology", "resourceGroup", resourceGroup)
			return topologyData, nil
		}
		log.Info("Cache miss for resourceGroup", "resourceGroup", resourceGroup)
	}

	log.Info("Calculating topology for cluster...", "cluster", cluster)
	// Build relationship graph based on GVK
	_, rg, _ := topology.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Count resources in all namespaces
	log.Info("Retrieving topology for cluster", "clusterName", cluster)
	rg, err := rg.CountRelationshipGraphByCustomResourceGroup(ctx, i.search, customResourceGroup, cluster)
	if err != nil {
		return nil, err
	}

	clusterTopologyMap := i.ConvertGraphToMap(rg, resourceGroup)
	i.clusterTopologyCache.Set(resourceGroup, clusterTopologyMap)
	log.Info("Added to cluster topology cache for resourceGroup", "resourceGroup", resourceGroup)

	return clusterTopologyMap, nil
}

// ConvertGraphToMap returns a map[string]ClusterTopology for a given relationship.RelationshipGraph
func (i *InsightManager) ConvertGraphToMap(rg *topology.RelationshipGraph, resourceGroup entity.ResourceGroup) map[string]ClusterTopology {
	ctm := make(map[string]ClusterTopology)
	for _, rgn := range rg.RelationshipNodes {
		rgnMap := rgn.ConvertToMap()
		resourceGroup := entity.ResourceGroup{
			Cluster:    resourceGroup.Cluster,
			APIVersion: schema.GroupVersion{Group: rgn.Group, Version: rgn.Version}.String(),
			Kind:       rgn.Kind,
			Namespace:  resourceGroup.Namespace,
		}
		ctm[rgn.GetHash()] = ClusterTopology{
			ResourceGroup: resourceGroup,
			Count:         rgn.ResourceCount,
			Relationship:  rgnMap,
		}
	}
	return ctm
}
