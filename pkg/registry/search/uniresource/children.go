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

package uniresource

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"

	"github.com/dominikbraun/graph"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// GetChildResourcesList returns an *unstructured.UnstructuredList representing all resources that matches the child GVK in the current namespace
func GetChildResourcesList(ctx context.Context, client *dynamic.DynamicClient, childRelation *Relationship, namespace string) (*unstructured.UnstructuredList, error) {
	gv, _ := schema.ParseGroupVersion(childRelation.Group + "/" + childRelation.Version)
	gvk := gv.WithKind(childRelation.Kind)
	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{})
	scope := Scope{"namespace"}
	mapper.Add(gvk, scope)
	mapping, err := mapper.RESTMapping(schema.GroupKind{Group: gv.Group, Kind: childRelation.Kind}, gv.Version)
	if err != nil {
		return nil, err
	}
	childrenRes := mapping.Resource
	klog.Infof("Listing child resource %s in namespace %s: \n", childRelation.Kind, namespace)
	var childrenResList *unstructured.UnstructuredList
	// Depends on whether child object is namespaced or not
	// TODO-think: Can this be derived from discovery.ServerResourcesForGroupVersion(version)?
	// TODO-think: Can this be retrieved from storage instead directly from cluster?
	if childRelation.ClusterScoped {
		childrenResList, err = client.Resource(childrenRes).List(ctx, metav1.ListOptions{})
	} else {
		childrenResList, err = client.Resource(childrenRes).Namespace(namespace).List(ctx, metav1.ListOptions{})
	}
	if err != nil {
		return nil, err
	}
	klog.Infof("List return size: %d\n", len(childrenResList.Items))
	return childrenResList, nil
}

// GetChildren returns a graph that includes all of the child resources for the current obj that are described by the childRelation
func GetChildren(ctx context.Context, client *dynamic.DynamicClient, obj unstructured.Unstructured, childRelation *Relationship, namespace, objName string, objResourceNode ResourceGraphNode, relationshipGraph graph.Graph[string, RelationshipGraphNode], resourceGraph graph.Graph[string, ResourceGraphNode]) (graph.Graph[string, ResourceGraphNode], error) {
	if childRelation.Type == "OwnerReference" {
		// If relationship type is ownerreference, honor that instead of relationship graph
		gv, _ := schema.ParseGroupVersion(childRelation.Group + "/" + childRelation.Version)
		gvk := gv.WithKind(childRelation.Kind)
		childrenResList, err := GetChildResourcesList(ctx, client, childRelation, namespace)
		if err != nil {
			return nil, err
		}
		resourceGraph, err = GetChildrenByOwnerReference(childrenResList, ctx, client, obj, gvk, relationshipGraph, resourceGraph)
		if err != nil {
			return nil, err
		}
	} else {
		// otherwise, use the children GVK on relationship graph to get a list of resources that match the children kind. Only proceed if the result size > 0.
		gv, _ := schema.ParseGroupVersion(childRelation.Group + "/" + childRelation.Version)
		gvk := gv.WithKind(childRelation.Kind)
		mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{})
		scope := Scope{"namespace"}
		mapper.Add(gvk, scope)
		mapping, err := mapper.RESTMapping(schema.GroupKind{Group: gv.Group, Kind: childRelation.Kind}, gv.Version)
		if err != nil {
			return nil, err
		}
		childRes := mapping.Resource
		klog.Infof("Listing child resource %s in namespace %s: \n", childRelation.Kind, namespace)
		var childrenResList *unstructured.UnstructuredList
		// Depends on whether child object is namespaced or not
		// TODO-think: Can this be derived from discovery.ServerResourcesForGroupVersion(version)?
		// TODO-think: Can this be retrieved from storage instead directly from cluster?
		if childRelation.ClusterScoped {
			childrenResList, err = client.Resource(childRes).List(ctx, metav1.ListOptions{})
		} else {
			childrenResList, err = client.Resource(childRes).Namespace(namespace).List(ctx, metav1.ListOptions{})
		}
		klog.Infof("List return size: %d\n", len(childrenResList.Items))
		if errors.IsNotFound(err) {
			klog.Infof("Obj %s in namespace %s not found\n", objName, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			klog.Infof("Error getting obj %s in namespace %s: %v\n", objName, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			return nil, err
		} else if len(childrenResList.Items) > 0 {
			if childRelation.Type == "JSONPath" {
				resourceGraph, err = GetByJSONPath(childrenResList, "child", ctx, client, obj, childRelation, gvk, objResourceNode, relationshipGraph, resourceGraph)
				if err != nil {
					return nil, err
				}
			} else if childRelation.Type == "Selector" {
				resourceGraph, err = GetByLabelSelector(childrenResList, "child", ctx, client, obj, childRelation, gvk, objResourceNode, relationshipGraph, resourceGraph)
				if err != nil {
					return nil, err
				}
			} else {
				klog.Infof("Something went wrong. Type should be either OwnerReference, Selector, or JSONPath")
			}
		}
	}
	return resourceGraph, nil
}

// GetChildrenByOwnerReference returns a graph that includes all of the child resources for the current obj described by their children's OwnerReferences field
func GetChildrenByOwnerReference(childrenResList *unstructured.UnstructuredList, ctx context.Context, client *dynamic.DynamicClient, obj unstructured.Unstructured, childGVK schema.GroupVersionKind, relationshipGraph graph.Graph[string, RelationshipGraphNode], resourceGraph graph.Graph[string, ResourceGraphNode]) (graph.Graph[string, ResourceGraphNode], error) {
	// For ownerreference-identified children, look up all instances of the child GVK and filter by ownerreference
	klog.Infof("Using OwnerReferences to find children...\n")
	gv, _ := schema.ParseGroupVersion(obj.GetAPIVersion())
	objResourceNode := ResourceGraphNode{
		Group:     gv.Group,
		Version:   gv.Version,
		Kind:      obj.GetKind(),
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}

	for _, childRes := range childrenResList.Items {
		if orMatch, err := OwnerReferencesMatch(obj, childRes); orMatch && err == nil {
			klog.Infof("Child resource found for kind %s, name %s based on OwnerReference.\n", obj.GetKind(), obj.GetName())
			klog.Infof("Child resource is: kind %s, name %s.\n", childRes.GetKind(), childRes.GetName())
			klog.Infof("---------------------------------------------------------------------------\n")
			cgv, _ := schema.ParseGroupVersion(childRes.GetAPIVersion())
			childResourceNode := ResourceGraphNode{
				Group:     cgv.Group,
				Version:   cgv.Version,
				Kind:      childRes.GetKind(),
				Name:      childRes.GetName(),
				Namespace: childRes.GetNamespace(),
			}
			resourceGraph.AddVertex(childResourceNode)
			resourceGraph.AddEdge(objResourceNode.GetHash(), childResourceNode.GetHash())
			childGVKOnGraph, _ := FindNodeOnGraph(relationshipGraph, childGVK.Group, childGVK.Version, childRes.GetKind())
			if len(childGVKOnGraph.Children) > 0 {
				// repeat for child resources
				// shorten call stack
				for _, childRelation := range childGVKOnGraph.Children {
					resourceGraph, _ = GetChildren(ctx, client, childRes, childRelation, childRes.GetNamespace(), childRes.GetName(), childResourceNode, relationshipGraph, resourceGraph)
				}
			}
		}
	}
	return resourceGraph, nil
}
