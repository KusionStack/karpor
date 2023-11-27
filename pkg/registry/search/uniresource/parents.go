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

	topologyutil "github.com/KusionStack/karbour/pkg/util/topology"
	"github.com/dominikbraun/graph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// GetParentResourcesList returns an *unstructured.UnstructuredList representing all resources that matches the parent GVK in the current namespace
func GetParentResourcesList(ctx context.Context, client *dynamic.DynamicClient, parentRelation *Relationship, namespace string) (*unstructured.UnstructuredList, error) {
	parentAPIVersion := parentRelation.Group + "/" + parentRelation.Version
	parentRes, err := topologyutil.GetGVRFromGVK(parentAPIVersion, parentRelation.Kind)
	if err != nil {
		return nil, err
	}
	klog.Infof("Listing parent resource %s in namespace %s: \n", parentRelation.Kind, namespace)
	var parentResList *unstructured.UnstructuredList
	// Depends on whether parent object is namespaced or not
	// TODO-think: Can this be derived from discovery.ServerResourcesForGroupVersion(version)?
	// TODO-think: Can this be retrieved from storage instead directly from cluster?
	if parentRelation.ClusterScoped {
		parentResList, err = client.Resource(parentRes).List(ctx, metav1.ListOptions{})
	} else {
		parentResList, err = client.Resource(parentRes).Namespace(namespace).List(ctx, metav1.ListOptions{})
	}
	if err != nil {
		return nil, err
	}
	if parentResList != nil {
		klog.Infof("List return size: %d\n", len(parentResList.Items))
	}
	return parentResList, nil
}

// GetParents returns a graph that includes all of the parent resources for the current obj that are described by the parentRelation
func GetParents(ctx context.Context, client *dynamic.DynamicClient, obj unstructured.Unstructured, parentRelation *Relationship, namespace, objName string, objResourceNode ResourceGraphNode, relationshipGraph graph.Graph[string, RelationshipGraphNode], resourceGraph graph.Graph[string, ResourceGraphNode]) (graph.Graph[string, ResourceGraphNode], error) {
	var err error
	if parentRelation.Type == "OwnerReference" {
		// If relationship type is ownerreference, honor that instead of relationship graph
		resourceGraph, err = GetParentsByOwnerReference(ctx, client, obj, objResourceNode, relationshipGraph, resourceGraph)
		if err != nil {
			return nil, err
		}
	} else {
		gv, _ := schema.ParseGroupVersion(parentRelation.Group + "/" + parentRelation.Version)
		gvk := gv.WithKind(parentRelation.Kind)
		parentResList, err := GetParentResourcesList(ctx, client, parentRelation, namespace)
		if errors.IsNotFound(err) {
			klog.Infof("Obj %s in namespace %s not found\n", objName, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			klog.Infof("Error getting obj %s in namespace %s: %v\n", objName, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			return nil, err
		} else if len(parentResList.Items) > 0 {
			if parentRelation.Type == "JSONPath" {
				resourceGraph, err = GetByJSONPath(parentResList, "parent", ctx, client, obj, parentRelation, gvk, objResourceNode, relationshipGraph, resourceGraph)
				if err != nil {
					return nil, err
				}
			} else if parentRelation.Type == "Selector" {
				resourceGraph, err = GetByLabelSelector(parentResList, "parent", ctx, client, obj, parentRelation, gvk, objResourceNode, relationshipGraph, resourceGraph)
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

// GetParentsByOwnerReference returns a graph that includes all of the parent resources for the current obj described by its OwnerReferences field
func GetParentsByOwnerReference(ctx context.Context, client *dynamic.DynamicClient, obj unstructured.Unstructured, objResourceNode ResourceGraphNode, relationshipGraph graph.Graph[string, RelationshipGraphNode], resourceGraph graph.Graph[string, ResourceGraphNode]) (graph.Graph[string, ResourceGraphNode], error) {
	klog.Infof("Using OwnerReferences to find parents...\n")
	objName := obj.GetName()
	namespace := obj.GetNamespace()
	objOwnerList := obj.GetOwnerReferences()
	var objOwner metav1.OwnerReference
	if len(objOwnerList) == 1 {
		objOwner = objOwnerList[0]
	} else if len(objOwnerList) == 0 {
		klog.Infof("No owner found for kind %s, name %s\n", obj.GetKind(), obj.GetName())
		return resourceGraph, nil
	} else {
		klog.Infof("Found more than 1 owner.")
		return resourceGraph, nil
	}

	parentRes, err := topologyutil.GetGVRFromGVK(objOwner.APIVersion, objOwner.Kind)
	if err != nil {
		return nil, err
	}

	klog.Infof("Listing parent resource %s in namespace %s: \n", objOwner.Kind, namespace)
	parentResList, err := client.Resource(parentRes).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if errors.IsNotFound(err) {
		klog.Infof("Obj %s in namespace %s not found\n", objName, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		klog.Infof("Error getting obj %s in namespace %s: %v\n", objName, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, err
	} else if len(parentResList.Items) > 0 {
		for _, parentRes := range parentResList.Items {
			if parentRes.GetUID() == objOwner.UID {
				klog.Infof("Parent resource found for kind %s, name %s based on OwnerReference. \n", obj.GetKind(), objName)
				klog.Infof("Parent resource is: kind %s, name %s.\n", parentRes.GetKind(), parentRes.GetName())
				klog.Infof("---------------------------------------------------------------------------\n")
				pgv, _ := schema.ParseGroupVersion(parentRes.GetAPIVersion())
				parentResourceNode := ResourceGraphNode{
					Group:     pgv.Group,
					Version:   pgv.Version,
					Kind:      parentRes.GetKind(),
					Name:      parentRes.GetName(),
					Namespace: parentRes.GetNamespace(),
				}
				resourceGraph.AddVertex(parentResourceNode)
				resourceGraph.AddEdge(parentResourceNode.GetHash(), objResourceNode.GetHash())
				if len(parentRes.GetOwnerReferences()) > 0 {
					resourceGraph, _ = GetParentsByOwnerReference(ctx, client, parentRes, parentResourceNode, relationshipGraph, resourceGraph)
				}
			}
		}
	}
	return resourceGraph, nil
}
