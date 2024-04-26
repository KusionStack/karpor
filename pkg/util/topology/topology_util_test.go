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

package topology

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestIsMapSubset(t *testing.T) {
	tests := []struct {
		m    map[string]string
		sub  map[string]string
		want bool
	}{
		// subset with equal maps
		{map[string]string{"a": "1"}, map[string]string{"a": "1"}, true},
		// subset with sub being smaller
		{map[string]string{"a": "1", "b": "2"}, map[string]string{"a": "1"}, true},
		// sub is not a subset
		{map[string]string{"a": "1"}, map[string]string{"b": "2"}, false},
		// empty subset
		{map[string]string{"a": "1"}, map[string]string{}, true},
	}

	for _, tt := range tests {
		t.Run("IsMapSubset", func(t *testing.T) {
			result := IsMapSubset(tt.m, tt.sub)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestJSONPathMatch(t *testing.T) {
	source := &unstructured.Unstructured{}
	source.SetName("testname")
	source.SetNamespace("testnamespace")

	target := &unstructured.Unstructured{}
	target.SetName("testname")
	target.SetNamespace("testnamespace-2")

	tests := []struct {
		source      *unstructured.Unstructured
		target      *unstructured.Unstructured
		criteriaSet []map[string]string
		wantMatch   bool
		wantErr     error
	}{
		// names match
		{
			source,
			target,
			[]map[string]string{
				{"name": "$.metadata.name"},
			},
			true,
			nil,
		},
		// namespaces don't match
		{
			source,
			target,
			[]map[string]string{
				{"namespace": "$.metadata.namespace"},
			},
			false,
			nil,
		},
		// invalid criteria key
		{
			source,
			target,
			[]map[string]string{
				{"invalid-key": "$.metadata.name"},
			},
			false,
			errors.New("invalid criteria key"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("JSONPathMatch-%d", i), func(t *testing.T) {
			match, err := JSONPathMatch(*tt.source, *tt.target, tt.criteriaSet)
			require.Equal(t, tt.wantMatch, match)
			if tt.wantErr != nil {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestGetGVRFromGVK(t *testing.T) {
	tests := []struct {
		apiVersion string
		kind       string
		wantGVR    schema.GroupVersionResource
		wantErr    bool
	}{
		// correct GroupVersionKind, should not error
		{
			"v1",
			"Pod",
			schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run("GetGVRFromGVK", func(t *testing.T) {
			gvr, err := GetGVRFromGVK(tt.apiVersion, tt.kind)
			require.Equal(t, tt.wantGVR, gvr)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestOwnerReferencesMatch(t *testing.T) {
	tests := []struct {
		name      string
		parent    unstructured.Unstructured
		child     unstructured.Unstructured
		wantMatch bool
	}{
		{
			name: "Parent in child's OwnerReferences",
			parent: func() unstructured.Unstructured {
				obj := unstructured.Unstructured{}
				obj.SetUID("parent-uid")
				return obj
			}(),
			child: func() unstructured.Unstructured {
				obj := unstructured.Unstructured{}
				obj.SetOwnerReferences([]metav1.OwnerReference{{
					UID: "parent-uid",
				}})
				return obj
			}(),
			wantMatch: true,
		},
		{
			name: "Parent not in child's OwnerReferences",
			parent: func() unstructured.Unstructured {
				obj := unstructured.Unstructured{}
				obj.SetUID("parent-uid")
				return obj
			}(),
			child: func() unstructured.Unstructured {
				obj := unstructured.Unstructured{}
				obj.SetOwnerReferences([]metav1.OwnerReference{{
					UID: "other-uid",
				}})
				return obj
			}(),
			wantMatch: false,
		},
		{
			name: "Child with no OwnerReferences",
			parent: func() unstructured.Unstructured {
				obj := unstructured.Unstructured{}
				obj.SetUID("parent-uid")
				return obj
			}(),
			child: func() unstructured.Unstructured {
				obj := unstructured.Unstructured{}
				return obj
			}(),
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := OwnerReferencesMatch(tt.parent, tt.child)
			require.NoError(t, err)
			require.Equal(t, tt.wantMatch, match)
		})
	}
}

func TestLabelSelectorsMatch(t *testing.T) {
	tests := []struct {
		name         string
		selectorPath string
		selectorObj  map[string]interface{}
		selectedObj  map[string]string
		wantMatch    bool
		wantErr      bool
	}{
		{
			name:         "Matching labels",
			selectorPath: "spec.selector.matchLabels",
			selectorObj:  map[string]interface{}{"spec": map[string]interface{}{"selector": map[string]interface{}{"matchLabels": map[string]interface{}{"app": "myapp"}}}},
			selectedObj:  map[string]string{"app": "myapp"},
			wantMatch:    true,
			wantErr:      false,
		},
		{
			name:         "Unmatching labels",
			selectorPath: "spec.selector.matchLabels",
			selectorObj:  map[string]interface{}{"spec": map[string]interface{}{"selector": map[string]interface{}{"matchLabels": map[string]interface{}{"app": "myapp"}}}},
			selectedObj:  map[string]string{"app": "wrongapp"},
			wantMatch:    false,
			wantErr:      false,
		},
		{
			name:         "Selector path not found",
			selectorPath: "invalid.path",
			selectorObj:  map[string]interface{}{"spec": map[string]interface{}{"selector": map[string]interface{}{"matchLabels": map[string]interface{}{"app": "myapp"}}}},
			selectedObj:  map[string]string{"app": "myapp"},
			wantMatch:    false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selectorUnstructured := &unstructured.Unstructured{
				Object: tt.selectorObj,
			}
			selectedUnstructured := &unstructured.Unstructured{}
			selectedUnstructured.SetLabels(tt.selectedObj)

			match, err := LabelSelectorsMatch(*selectorUnstructured, *selectedUnstructured, tt.selectorPath)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantMatch, match)
		})
	}
}
