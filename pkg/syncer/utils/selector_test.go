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
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

// 测试数据和辅助函数
func createTestSelectable(labels map[string]interface{}, fields map[string]interface{}) Selectable {
	return &selectableUnstructured{
		Unstructured: &unstructured.Unstructured{
			Object: map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": labels,
				},
				"spec": fields,
			}},
		parser: DefaultJSONPathParser,
	}
}

func TestGetLabels(t *testing.T) {
	expectedLabels := labels.Set{"app": "test"}
	selectable := createTestSelectable(map[string]interface{}{"app": "test"}, nil)
	require.Equal(t, expectedLabels, selectable.GetLabels())
}

func TestGetFields(t *testing.T) {
	u := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"spec.replicas": "2",
		},
	}
	expectedFields := NewJSONPathFields(NewJSONPathParser(), u.Object)
	label := selectableUnstructured{Unstructured: u, parser: DefaultJSONPathParser}
	require.Equal(t, expectedFields, label.GetFields())
}

func TestSelector_Matches(t *testing.T) {
	labelSelector := labels.SelectorFromSet(labels.Set{"app": "test"})
	fieldSelector := fields.SelectorFromSet(fields.Set{"spec.replicas": "2"})
	selector := Selector{Label: labelSelector, Field: FieldsSelector{Selector: fieldSelector, ServerSupported: true}}

	testCases := []struct {
		name     string
		labels   map[string]interface{}
		fields   map[string]interface{}
		expected bool
	}{
		{"MatchBoth", map[string]interface{}{"app": "test"}, map[string]interface{}{"replicas": "2"}, true},
		{"MismatchLabel", map[string]interface{}{"app": "wrong"}, map[string]interface{}{"replicas": "2"}, false},
		{"MismatchField", map[string]interface{}{"app": "test"}, map[string]interface{}{"replicas": "3"}, false},
		{"NoFields", map[string]interface{}{"app": "test"}, nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			selectable := createTestSelectable(tc.labels, tc.fields)
			require.Equal(t, tc.expected, selector.Matches(selectable))
		})
	}
}

func TestMultiSelectors_Matches(t *testing.T) {
	labelSelector := labels.SelectorFromSet(labels.Set{"app": "test"})
	fieldSelector := fields.SelectorFromSet(fields.Set{"spec.replicas": "2"})
	selector := Selector{Label: labelSelector, Field: FieldsSelector{Selector: fieldSelector, ServerSupported: true}}
	selectors := MultiSelectors{selector}

	testCases := []struct {
		name     string
		labels   map[string]interface{}
		fields   map[string]interface{}
		expected bool
	}{
		{"MatchBoth", map[string]interface{}{"app": "test"}, map[string]interface{}{"replicas": "2"}, true},
		{"MismatchLabel", map[string]interface{}{"app": "wrong"}, map[string]interface{}{"replicas": "2"}, false},
		{"MismatchField", map[string]interface{}{"app": "test"}, map[string]interface{}{"replicas": "3"}, false},
		{"NoFields", map[string]interface{}{"app": "test"}, nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			selectable := createTestSelectable(tc.labels, tc.fields)
			require.Equal(t, tc.expected, selectors.Matches(selectable))
		})
	}
}

func TestApplyToList(t *testing.T) {
	labelSelector := labels.SelectorFromSet(labels.Set{"app": "test"})
	fieldSelector := fields.SelectorFromSet(fields.Set{"spec.replicas": "2"})
	selector := Selector{Label: labelSelector, Field: FieldsSelector{Selector: fieldSelector, ServerSupported: true}}
	selectors := MultiSelectors{selector}

	meta := &metav1.ListOptions{}
	selectors.ApplyToList(meta)

	require.Equal(t, "app=test", meta.LabelSelector)
	require.Equal(t, "spec.replicas=2", meta.FieldSelector)
}

func TestPredicate(t *testing.T) {
	labelSelector := labels.SelectorFromSet(labels.Set{"app": "test"})
	fieldSelector := fields.SelectorFromSet(fields.Set{"spec.replicas": "2"})
	s1 := Selector{Label: labelSelector, Field: FieldsSelector{Selector: fieldSelector, ServerSupported: true}}
	labelSelector = labels.SelectorFromSet(labels.Set{"app": "test"})
	fieldSelector = fields.SelectorFromSet(fields.Set{"spec.replicas": "4"})
	s2 := Selector{Label: labelSelector, Field: FieldsSelector{Selector: fieldSelector, ServerSupported: true}}
	selectors := MultiSelectors{s1, s2}

	testCases := []struct {
		name     string
		labels   map[string]interface{}
		fields   map[string]interface{}
		expected bool
	}{
		{"MatchBoth", map[string]interface{}{"app": "test"}, map[string]interface{}{"replicas": "2"}, true},
		{"MismatchLabel", map[string]interface{}{"app": "wrong"}, map[string]interface{}{"replicas": "2"}, false},
		{"MismatchField", map[string]interface{}{"app": "test"}, map[string]interface{}{"replicas": "3"}, false},
		{"NoFields", map[string]interface{}{"app": "test"}, nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, selectors.Predicate(&unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": tc.labels,
					},
					"spec": tc.fields,
				}}))
		})
	}
}
