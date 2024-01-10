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
	"strings"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/infra/multicluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GetResourceEvents returns the list of events specified by core.Locator
func (i *InsightManager) GetResourceEvents(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) ([]unstructured.Unstructured, error) {
	// Retrieve the list of events for the specific resource
	eventGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}
	eventList, err := client.DynamicClient.Resource(eventGVR).Namespace(loc.Namespace).List(ctx, metav1.ListOptions{
		// FieldSelector is case-sensitive so this would depend on user input. Safer way is to list all events within namespace and compare afterwards
		// FieldSelector: fmt.Sprintf("involvedObject.apiVersion=%s,involvedObject.kind=%s,involvedObject.name=%s", res.APIVersion, res.Kind, res.Name),
	})
	if err != nil {
		return nil, err
	}
	filteredList := make([]unstructured.Unstructured, 0)
	// Iterate over the list and filter events for the specific resource
	for _, event := range eventList.Items {
		involvedObjectName, foundName, _ := unstructured.NestedString(event.Object, "involvedObject", "name")
		involvedObjectAPIVersion, foundAPIVersion, _ := unstructured.NestedString(event.Object, "involvedObject", "apiVersion")
		involvedObjectKind, foundKind, _ := unstructured.NestedString(event.Object, "involvedObject", "kind")
		// case-insensitive comparison
		if foundName && foundKind && foundAPIVersion &&
			strings.EqualFold(involvedObjectName, loc.Name) &&
			strings.EqualFold(involvedObjectAPIVersion, loc.APIVersion) &&
			strings.EqualFold(involvedObjectKind, loc.Kind) {
			filteredList = append(filteredList, event)
		}
	}

	return filteredList, nil
}

// GetNamespaceEvents returns the complete list of events in a namespace
func (i *InsightManager) GetNamespaceGVKEvents(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) ([]unstructured.Unstructured, error) {
	eventGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}
	eventList, err := client.DynamicClient.Resource(eventGVR).Namespace(loc.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	filteredList := make([]unstructured.Unstructured, 0)
	for _, event := range eventList.Items {
		involvedObjectAPIVersion, foundAPIVersion, _ := unstructured.NestedString(event.Object, "involvedObject", "apiVersion")
		involvedObjectKind, foundKind, _ := unstructured.NestedString(event.Object, "involvedObject", "kind")
		if foundAPIVersion && foundKind && strings.EqualFold(involvedObjectAPIVersion, loc.APIVersion) && strings.EqualFold(involvedObjectKind, loc.Kind) {
			filteredList = append(filteredList, event)
		}
	}
	return filteredList, nil
}

// GetNamespaceEvents returns the complete list of events in a namespace
func (i *InsightManager) GetNamespaceEvents(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) ([]unstructured.Unstructured, error) {
	eventGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}
	// Another option is to add .Namespace(loc.Namespace) here
	// It is quicker but it will not include the events that are related to the namespace resource itself
	eventList, err := client.DynamicClient.Resource(eventGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	filteredList := make([]unstructured.Unstructured, 0)
	for _, event := range eventList.Items {
		involvedObjectName, foundName, _ := unstructured.NestedString(event.Object, "involvedObject", "name")
		involvedObjectKind, foundKind, _ := unstructured.NestedString(event.Object, "involvedObject", "kind")
		involvedObjectNamespace, foundNamespace, _ := unstructured.NestedString(event.Object, "involvedObject", "namespace")
		// Either the event is for a resource whose namespace == locator's namespace, or
		// the event is for the namespace itself, in which case we are checking if
		// involvedObjectKind == "Namespace" AND involvedObjectName == locator's namespace
		if (foundNamespace && strings.EqualFold(involvedObjectNamespace, loc.Namespace)) ||
			(foundName && foundKind && strings.EqualFold(involvedObjectKind, "Namespace") && strings.EqualFold(involvedObjectName, loc.Namespace)) {
			filteredList = append(filteredList, event)
		}
	}
	return filteredList, nil
}

// GetGVKEvents returns the complete list of events for a GVK
func (i *InsightManager) GetGVKEvents(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) ([]unstructured.Unstructured, error) {
	eventGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}
	eventList, err := client.DynamicClient.Resource(eventGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	filteredList := make([]unstructured.Unstructured, 0)
	for _, event := range eventList.Items {
		involvedObjectAPIVersion, foundAPIVersion, _ := unstructured.NestedString(event.Object, "involvedObject", "apiVersion")
		involvedObjectKind, foundKind, _ := unstructured.NestedString(event.Object, "involvedObject", "kind")
		if foundAPIVersion && foundKind && strings.EqualFold(involvedObjectAPIVersion, loc.APIVersion) && strings.EqualFold(involvedObjectKind, loc.Kind) {
			filteredList = append(filteredList, event)
		}
	}
	return filteredList, nil
}

// GetClusterEvents returns the complete list of events in a cluster
func (i *InsightManager) GetClusterEvents(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) ([]unstructured.Unstructured, error) {
	eventGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}
	eventList, err := client.DynamicClient.Resource(eventGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return eventList.Items, nil
}
