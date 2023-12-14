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

package search

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SyncRegistry struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   SyncRegistrySpec
	Status SyncRegistryStatus
}

type SyncRegistrySpec struct {
	// ClusterLabelSelector is used to filter the target clusters that need to be synced from.
	ClusterLabelSelector *metav1.LabelSelector

	// Clusters is the list of the target clusters to be be synced from.
	Clusters []string

	SyncResources []ResourceSyncRule

	SyncResourcesRefName string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SyncRegistryList struct {
	metav1.TypeMeta

	metav1.ListMeta

	Items []SyncRegistry
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SyncResources struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec SyncResourcesSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SyncResourcesList struct {
	metav1.TypeMeta

	metav1.ListMeta

	Items []SyncResources
}

type SyncResourcesSpec struct {
	SyncResources []ResourceSyncRule
}

// ResourceSyncRule is used to specify the way to sync the specified resource
type ResourceSyncRule struct {
	// APIVersion represents the group version of the target resource.
	APIVersion string

	// Resource is the the target resource.
	Resource string

	// Namespace specifies the namespace in which the ListWatch of the target resources is limited to.
	Namespace string

	// ResynPeriod is the period to resync
	ResyncPeriod *metav1.Duration

	// MaxConcurrent is the maximum number of workers (default: 10)
	MaxConcurrent int

	// Selectors are used to filter the target resources to sync. Multiple selectors are ORed.
	Selectors []Selector

	// Transform is the rule applied to the original resource to transform it to the desired target resource.
	Transform *TransformRuleSpec

	// TransformRefName is the name of the TransformRule
	TransformRefName string
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TransformRule is used to define the rule to transform the original resource into the desired target resource.
type TransformRule struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec TransformRuleSpec
}

type TransformRuleSpec struct {
	// Type is the type of transformer.
	Type string

	// ValueTemplate is the template of the input data to be paased to the transformer
	ValueTemplate string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TransformRuleList struct {
	metav1.TypeMeta

	metav1.ListMeta

	Items []TransformRule
}

// Selector represents a resource filter
type Selector struct {
	// LabelSelector is a filter to select resources by labels.
	// If non-nil and non-empty, only the resource match this filter will be selected.
	LabelSelector *metav1.LabelSelector

	// FieldSelector is a filter to select resources by fields.
	// If non-nil and non-empty, only the resource match this filter will be selected.
	FieldSelector *FieldSelector
}

// FieldSelector is a field filter.
type FieldSelector struct {
	// MatchFields is a map of {field,value} pairs. A single {field,value} in the matchFields
	// map means that the specified field should have an exact match with the specified value. Multiple entries are ANDed.
	MatchFields map[string]string
	// ServerSupported specifies whether the field selection is supported in api server side
	SeverSupported bool
}

type SyncRegistryStatus struct {
	Clusters []ClusterResourcesSyncCondition

	LastTransitionTime metav1.Time
}

type ClusterResourcesSyncCondition struct {
	Cluster string

	Status string

	Resources []ResourceSyncCondition
}

type ResourceSyncCondition struct {
	APIVersion string

	Kind string

	Status string

	Reason string

	Message string

	LastTransitionTime metav1.Time
}
