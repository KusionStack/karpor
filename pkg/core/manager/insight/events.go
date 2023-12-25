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

// GetResourceSummary returns the unstructured cluster object summary for a given cluster. Possibly will add more metrics to it in the future.
func (m *InsightManager) GetResourceEvents(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) ([]unstructured.Unstructured, error) {
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
