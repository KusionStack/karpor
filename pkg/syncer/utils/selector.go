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

package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

type Selectable interface {
	GetLabels() labels.Labels
	GetFields() fields.Fields
}

type Selector struct {
	Label labels.Selector
	Field FieldsSelector
}

func (s Selector) ServerSupported() bool {
	if !s.Field.Empty() && !s.Field.ServerSupported {
		return false
	}
	return true
}

type FieldsSelector struct {
	fields.Selector
	ServerSupported bool
}

func (fs FieldsSelector) Empty() bool {
	if fs.Selector == nil {
		return true
	}
	return fs.Selector.Empty()
}

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

type MultiSelectors []Selector

func (m MultiSelectors) Matches(obj Selectable) bool {
	for _, s := range m {
		if s.Matches(obj) {
			return true
		}
	}
	return false
}

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

func (m MultiSelectors) canApplyToListWatch() bool {
	if len(m) > 1 {
		return false
	}
	if len(m) == 1 {
		return m[0].ServerSupported()
	}
	return true
}

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

type selectableUnstructured struct {
	*unstructured.Unstructured
	parser *JSONPathParser
}

func (w selectableUnstructured) GetLabels() labels.Labels {
	return labels.Set(w.Unstructured.GetLabels())
}

func (w selectableUnstructured) GetFields() fields.Fields {
	return NewJSONPathFields(w.parser, w.UnstructuredContent())
}
