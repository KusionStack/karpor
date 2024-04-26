// Copyright The Karpor Authors.
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

package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

// Selectable defines the interface for objects that can be selected by labels and fields.
type Selectable interface {
	GetLabels() labels.Labels
	GetFields() fields.Fields
}

// Selector represents a selection based on labels and fields.
type Selector struct {
	Label labels.Selector
	Field FieldsSelector
}

// ServerSupported checks if the server supports the fields and labels specified in the Selector.
func (s Selector) ServerSupported() bool {
	if !s.Field.Empty() && !s.Field.ServerSupported {
		return false
	}
	return true
}

// FieldsSelector represents a selection based on fields.
type FieldsSelector struct {
	fields.Selector
	ServerSupported bool
}

// Empty checks if the FieldsSelector is empty, meaning no fields are specified for selection.
func (fs FieldsSelector) Empty() bool {
	if fs.Selector == nil {
		return true
	}
	return fs.Selector.Empty()
}

// Matches checks if the provided Selectable object matches the selection criteria specified by the Selector.
func (s Selector) Matches(obj Selectable) bool {
	if s.Label != nil &&
		!s.Label.Matches(obj.GetLabels()) {
		return false
	}
	if !s.Field.Empty() &&
		!s.Field.Matches(obj.GetFields()) {
		return false
	}
	return true
}

// MultiSelectors is a slice of Selectors, allowing for multiple selection criteria to be combined.
type MultiSelectors []Selector

// Matches checks if the provided Selectable object matches any of the selection criteria specified in the MultiSelectors.
func (m MultiSelectors) Matches(obj Selectable) bool {
	for _, s := range m {
		if s.Matches(obj) {
			return true
		}
	}
	return false
}

// ApplyToList applies the MultiSelectors to the provided metav1.ListOptions for use in Kubernetes API list requests.
func (m MultiSelectors) ApplyToList(options *metav1.ListOptions) {
	if len(m) == 0 || !m.canApplyToListWatch() {
		return
	}

	selector := m[0]
	if selector.Label != nil && !selector.Label.Empty() {
		options.LabelSelector = selector.Label.String()
	}
	if !selector.Field.Empty() {
		options.FieldSelector = selector.Field.String()
	}
}

// canApplyToListWatch checks if the MultiSelectors can be applied to list and watch Kubernetes API requests.
func (m MultiSelectors) canApplyToListWatch() bool {
	if len(m) > 1 {
		return false
	}
	if len(m) == 1 {
		return m[0].ServerSupported()
	}
	return true
}

// Predicate returns a function that implements the Predicate interface, allowing the MultiSelectors to be used as a go-restful Predicate.
func (m MultiSelectors) Predicate(obj interface{}) bool {
	if len(m) == 0 || m.canApplyToListWatch() {
		// no filtering after list/watch
		return true
	}

	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return false
	}
	return m.Matches(selectableUnstructured{Unstructured: u, parser: DefaultJSONPathParser})
}

// selectableUnstructured is a wrapper around unstructured.Unstructured for implementing the Selectable interface.
type selectableUnstructured struct {
	*unstructured.Unstructured
	parser *JSONPathParser
}

// GetLabels retrieves the labels from the wrapped unstructured.Unstructured object.
func (w selectableUnstructured) GetLabels() labels.Labels {
	return labels.Set(w.Unstructured.GetLabels())
}

// GetFields retrieves the fields from the wrapped unstructured.Unstructured object using the embedded JSONPathParser.
func (w selectableUnstructured) GetFields() fields.Fields {
	return NewJSONPathFields(w.parser, w.UnstructuredContent())
}
