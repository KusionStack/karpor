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

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestESListerGetter_ListKeys(t *testing.T) {
	// Set up mocks for dynamic client
	mockey.Mock((*elasticsearch.Storage).SearchByQuery).Return(&storage.SearchResult{
		Resources: []*storage.Resource{
			{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "default",
						"namespace": "existing-namespace",
					},
				},
			},
		},
	}, nil).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		gvr            schema.GroupVersionResource
		expectedLength int
		expectError    bool
	}{
		{
			name: "Success - ListKeys",
			gvr: schema.GroupVersionResource{
				Version:  "v1",
				Resource: "Pod",
			},
			expectedLength: 1,
			expectError:    false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call ListKeys method
			esLister := NewESListerGetter("defalut", &elasticsearch.Storage{}, tc.gvr)
			keys := esLister.ListKeys()

			// Check error expectation
			if tc.expectError {
				require.Len(t, keys, tc.expectedLength, "Expected nil result on error")
			} else {
				require.Len(t, keys, tc.expectedLength, "Expected nil result on error")
			}
		})
	}

}

func TestESListerGetter_GetByKey(t *testing.T) {
	// Set up mocks for dynamic client
	mockey.Mock((*elasticsearch.Storage).SearchByQuery).Return(&storage.SearchResult{
		Resources: []*storage.Resource{
			{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "default",
						"namespace": "existing-namespace",
					},
				},
			},
		},
	}, nil).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		gvr            schema.GroupVersionResource
		expectedResult *unstructured.Unstructured
		expectError    bool
	}{
		{
			name: "Success - GetByKey",
			gvr: schema.GroupVersionResource{
				Version:  "v1",
				Resource: "Pod",
			},
			expectedResult: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "default",
						"namespace": "existing-namespace",
					},
				},
			},
			expectError: false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call ListKeys method
			esLister := NewESListerGetter("defalut", &elasticsearch.Storage{}, tc.gvr)
			value, exists, err := esLister.GetByKey("existing-namespace/default")

			// Check error expectation
			if tc.expectError {
				require.Error(t, err)
				require.False(t, exists)
			} else {
				require.True(t, exists)
				require.Equal(t, tc.expectedResult, value)
			}
		})
	}
}
