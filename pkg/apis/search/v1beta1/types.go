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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SyncClustersResources struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// +optional
	Spec SyncClustersResourcesSpec `json:"spec,omitempty"`

	// +optional
	Status SyncClustersResourcesStatus `json:"status,omitempty"`
}

type SyncClustersResourcesSpec struct {
	// ClusterSelector is used to filter the target clusters that need to be synced from.
	// +optional
	ClusterSelector Selector `json:"clusterSelector,omitempty"`

	// ClusterNames is the list of the target clusters to be be synced from.
	// +optional
	ClusterNames []string `json:"clusterNames,omitempty"`

	// +optional
	SyncResources []ResourceSyncRule `json:"syncResources,omitempty"`

	// +optional
	SyncResourcesRefName string `json:"SyncResourcesRefName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SyncClustersResourcesList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SyncClustersResources `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SyncResources struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec SyncResourcesSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SyncResourcesList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SyncResources `json:"items"`
}

type SyncResourcesSpec struct {
	// +optional
	SyncResources []ResourceSyncRule `json:"syncResources,omitempty"`
}

// ResourceSyncRule is used to specify the way to sync the specified resource
type ResourceSyncRule struct {
	// APIVersion represents the group version of the target resource.
	// +required
	APIVersion string `json:"apiVersion"`

	// Kind represents the kind of the target resource.
	// +required
	Kind string `json:"kind"`

	// Namespace specifies the namespace in which the ListWatch of the target resources is limited to.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// selectors are used to filter the target resources to sync. Multiple selectors are ORed.
	// +optional
	Selectors []Selector `json:"selectors,omitempty"`

	// Transform is the rule applied to the original resource to transform it to the desired target resource.
	// +optional
	Transform TransformRuleSpec `json:"transform,omitempty"`

	// TransformRefName is the name of the TransformRule
	// +optional
	TransformRefName string `json:"transformRefName,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// TransformRule is used to define the rule to transform the original resource into the desired target resource.
type TransformRule struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// +optional
	Spec TransformRuleSpec `json:"spec,omitempty"`
}

type TransformRuleSpec struct {
	// Type is the type of transformer.
	// +required
	Type string `json:"type"`

	// ValueTemplate is the template of the input data to be paased to the transformer
	// +required
	ValueTemplate string `json:"valueTemplate"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TransformRuleList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []TransformRule `json:"items"`
}

// Selector represents a resource filter
type Selector struct {
	// LabelSelector is a filter to select resources by labels.
	// If non-nil and non-empty, only the resource match this filter will be selected.
	// +optional
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`

	// FieldSelector is a filter to select resources by fields.
	// If non-nil and non-empty, only the resource match this filter will be selected.
	// +optional
	FieldSelector *FieldSelector `json:"fieldSelector,omitempty"`
}

// FieldSelector is a field filter.
type FieldSelector struct {
	// matchFields is a map of {field,value} pairs. A single {field,value} in the matchFields
	// map means that the specified field should have an exact match with the specified value. Multiple entries are ANDed.
	// +optional
	MatchFields map[string]string `json:"matchFields,omitempty"`
}

type SyncClustersResourcesStatus struct {
	// +optional
	Clusters []ClusterSyncResourcesCondition `json:"clusters"`

	// +required
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}

type ClusterSyncResourcesCondition struct {
	// +required
	Cluster string `json:"cluster"`

	// +required
	Status string `json:"status"`

	// optional
	Resources []ResourceSyncCondition `json:"resources"`
}

type ResourceSyncCondition struct {
	// +required
	APIVersion string

	// +required
	Kind string `json:"kind"`

	// +required
	Status string `json:"status"`

	// +optional
	Reason string `json:"reason,omitempty"`

	// +optional
	Message string `json:"message,omitempty"`

	// +required
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}