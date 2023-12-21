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

package resource

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/relationship"
	"github.com/KusionStack/karbour/pkg/util/cache"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	topologyutil "github.com/KusionStack/karbour/pkg/util/topology"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8syaml "sigs.k8s.io/yaml"
)

type ResourceManager struct {
	config *ResourceConfig
	cache  *cache.Cache[core.Locator, map[string]ResourceTopology]
}

// NewResourceManager returns a new ResourceManager
func NewResourceManager(config *ResourceConfig) *ResourceManager {
	return &ResourceManager{
		config: config,
		cache:  cache.NewCache[core.Locator, map[string]ResourceTopology](10 * time.Minute),
	}
}

// GetResource returns the unstructured cluster object for a given cluster
func (r *ResourceManager) GetResource(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) (*unstructured.Unstructured, error) {
	resourceGVR, err := topologyutil.GetGVRFromGVK(loc.APIVersion, loc.Kind)
	if err != nil {
		return nil, err
	}
	return client.DynamicClient.Resource(resourceGVR).Namespace(loc.Namespace).Get(ctx, loc.Name, metav1.GetOptions{})
}

// GetResourceSummary returns the unstructured cluster object summary for a given cluster. Possibly will add more metrics to it in the future.
func (r *ResourceManager) GetResourceSummary(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) (*ResourceSummary, error) {
	obj, err := r.GetResource(ctx, client, loc)
	if err != nil {
		return nil, err
	}

	return &ResourceSummary{
		Resource: Resource{
			Name:       obj.GetName(),
			Namespace:  obj.GetNamespace(),
			APIVersion: obj.GetAPIVersion(),
			Cluster:    loc.Cluster,
			Kind:       obj.GetKind(),
		},
		CreationTimestamp: obj.GetCreationTimestamp(),
		ResourceVersion:   obj.GetResourceVersion(),
		UID:               obj.GetUID(),
	}, nil
}

// GetResourceSummary returns the unstructured cluster object summary for a given cluster. Possibly will add more metrics to it in the future.
func (r *ResourceManager) GetResourceEvents(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) ([]unstructured.Unstructured, error) {
	eventGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}
	var eventList *unstructured.UnstructuredList
	filteredList := make([]unstructured.Unstructured, 0)
	// Retrieve the list of events for the specific resource
	eventList, err := client.DynamicClient.Resource(eventGVR).Namespace(loc.Namespace).List(context.TODO(), metav1.ListOptions{
		// FieldSelector is case-sensitive so this would depend on user input. Safer way is to list all events within namespace and compare afterwards
		// FieldSelector: fmt.Sprintf("involvedObject.apiVersion=%s,involvedObject.kind=%s,involvedObject.name=%s", res.APIVersion, res.Kind, res.Name),
	})
	if err != nil {
		return nil, err
	}
	// Iterate over the list and filter events for the specific resource
	for _, event := range eventList.Items {
		involvedObjectName, foundName, _ := unstructured.NestedString(event.Object, "involvedObject", "name")
		involvedObjectKind, foundKind, _ := unstructured.NestedString(event.Object, "involvedObject", "kind")
		// case-insensitive comparison
		if foundName && foundKind && strings.EqualFold(involvedObjectName, loc.Name) && strings.EqualFold(involvedObjectKind, loc.Kind) {
			filteredList = append(filteredList, event)
		}
	}

	return filteredList, nil
}

// GetYAMLForResource returns the yaml byte array for a given cluster
func (r *ResourceManager) GetYAMLForResource(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) ([]byte, error) {
	obj, err := r.GetResource(ctx, client, loc)
	if err != nil {
		return nil, err
	}
	return k8syaml.Marshal(obj.Object)
}

// GetTopologyForResource returns a map that describes topology for a given cluster
func (r *ResourceManager) GetTopologyForResource(ctx context.Context, client *multicluster.MultiClusterClient, locator *core.Locator) (map[string]ResourceTopology, error) {
	log := ctxutil.GetLogger(ctx)

	if topologyData, exist := r.cache.Get(*locator); exist {
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
	g, err = r.GetResourceRelationship(ctx, client, *unObj, rg, g)
	if err != nil {
		return nil, err
	}

	topologyMap := r.ConvertResourceGraphToMap(g)
	r.cache.Set(*locator, topologyMap)
	log.Info("Added to resource topology cache for locator", "locator", locator)

	// Draw graph
	// TODO: This is drawn on the server side, not needed eventually
	file, _ := os.Create("./resource.gv")
	_ = draw.DOT(g, file)

	return topologyMap, nil
}

// GetResourceRelationship returns a full graph that contains all the resources that are related to obj
func (r *ResourceManager) GetResourceRelationship(ctx context.Context, client *multicluster.MultiClusterClient, obj unstructured.Unstructured, relationshipGraph graph.Graph[string, relationship.RelationshipGraphNode], resourceGraph graph.Graph[string, relationship.ResourceGraphNode]) (graph.Graph[string, relationship.ResourceGraphNode], error) {
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

func (r *ResourceManager) ConvertResourceGraphToMap(g graph.Graph[string, relationship.ResourceGraphNode]) map[string]ResourceTopology {
	am, _ := g.AdjacencyMap()
	m := make(map[string]ResourceTopology)
	for key, edgeMap := range am {
		childList := make([]string, 0)
		for edgeTarget := range edgeMap {
			childList = append(childList, edgeTarget)
		}
		m[key] = ResourceTopology{
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
		if node, ok := m[key]; ok {
			node.Parents = parentList
			m[key] = node
		}
	}
	return m
}

func (r *ResourceManager) ConvertResourceToLocator(res *Resource) *core.Locator {
	if res != nil {
		return &core.Locator{
			Cluster:    res.Cluster,
			APIVersion: res.APIVersion,
			Kind:       res.Kind,
			Namespace:  res.Namespace,
			Name:       res.Name,
		}
	}
	return nil
}
