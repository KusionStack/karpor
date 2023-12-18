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

	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/relationship"
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
}

// NewResourceManager returns a new ResourceManager
func NewResourceManager(config *ResourceConfig) *ResourceManager {
	return &ResourceManager{
		config: config,
	}
}

// GetResource returns the unstructured cluster object for a given cluster
func (r *ResourceManager) GetResource(ctx context.Context, client *multicluster.MultiClusterClient, res *Resource) (*unstructured.Unstructured, error) {
	resourceGVR, err := topologyutil.GetGVRFromGVK(res.APIVersion, res.Kind)
	if err != nil {
		return nil, err
	}
	return client.DynamicClient.Resource(resourceGVR).Namespace(res.Namespace).Get(ctx, res.Name, metav1.GetOptions{})
}

// GetResourceSummary returns the unstructured cluster object summary for a given cluster. Possibly will add more metrics to it in the future.
func (r *ResourceManager) GetResourceSummary(ctx context.Context, client *multicluster.MultiClusterClient, res *Resource) (*ResourceSummary, error) {
	obj, err := r.GetResource(ctx, client, res)
	if err != nil {
		return nil, err
	}

	return &ResourceSummary{
		Resource: Resource{
			Name:       obj.GetName(),
			Namespace:  obj.GetNamespace(),
			APIVersion: obj.GetAPIVersion(),
			Cluster:    res.Cluster,
			Kind:       obj.GetKind(),
		},
		CreationTimestamp: obj.GetCreationTimestamp(),
		ResourceVersion:   obj.GetResourceVersion(),
		UID:               obj.GetUID(),
	}, nil
}

// GetYAMLForResource returns the yaml byte array for a given cluster
func (r *ResourceManager) GetYAMLForResource(ctx context.Context, client *multicluster.MultiClusterClient, res *Resource) ([]byte, error) {
	obj, err := r.GetResource(ctx, client, res)
	if err != nil {
		return nil, err
	}
	return k8syaml.Marshal(obj.Object)
}

// GetTopologyForResource returns a map that describes topology for a given cluster
func (r *ResourceManager) GetTopologyForResource(ctx context.Context, client *multicluster.MultiClusterClient, res *Resource) (map[string]ResourceTopology, error) {
	log := ctxutil.GetLogger(ctx)

	// Build relationship graph based on GVK
	rg, _, err := relationship.BuildRelationshipGraph(ctx, client.DynamicClient)
	if err != nil {
		return nil, err
	}
	log.Info("Retrieving topology for resource", "resourceName", res.Name)

	ResourceGraphNodeHash := func(rgn relationship.ResourceGraphNode) string {
		return rgn.Group + "/" + rgn.Version + "." + rgn.Kind + ":" + rgn.Namespace + "." + rgn.Name
	}
	g := graph.New(ResourceGraphNodeHash, graph.Directed(), graph.PreventCycles())

	// Get target resource
	resourceGVR, err := topologyutil.GetGVRFromGVK(res.APIVersion, res.Kind)
	if err != nil {
		return nil, err
	}
	resObj, _ := client.DynamicClient.Resource(resourceGVR).Namespace(res.Namespace).Get(ctx, res.Name, metav1.GetOptions{})
	unObj := &unstructured.Unstructured{}
	unObj.SetUnstructuredContent(resObj.Object)

	// Build resource graph for target resource
	g, err = r.GetResourceRelationship(ctx, client, *unObj, rg, g)
	if err != nil {
		return nil, err
	}

	// Draw graph
	// TODO: This is drawn on the server side, not needed eventually
	file, _ := os.Create("./resource.gv")
	_ = draw.DOT(g, file)

	return r.ConvertResourceGraphToMap(g), nil
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
