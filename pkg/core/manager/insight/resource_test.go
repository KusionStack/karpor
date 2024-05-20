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

package insight

import (
	"context"
	"reflect"
	"testing"

	"kusionstack.io/karpor/pkg/core/entity"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/dynamic"
)

// TestGetResource tests the TestGetResource method of the InsightManager for
// various scenarios.
func TestGetResource(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		resourceGroup  *entity.ResourceGroup
		expectError    bool
		expectSanitize bool
	}{
		{
			name: "Success - Existing ConfigMap",
			resourceGroup: &entity.ResourceGroup{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Namespace:  "default",
				Name:       "existing-configmap",
			},
			expectError:    false,
			expectSanitize: false, // Not a Secret kind, so no sanitization expected
		},
		{
			name:           "Error - Non-existing cluster",
			resourceGroup:  &entity.ResourceGroup{},
			expectError:    true,
			expectSanitize: false, // Not applicable as there is an error
		},
		{
			name: "Success - Existing Secret",
			resourceGroup: &entity.ResourceGroup{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Secret",
				Namespace:  "default",
				Name:       "existing-secret",
			},
			expectError:    false,
			expectSanitize: true, // Expecting data sanitization for Secret kind
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetResource method
			resource, err := manager.GetResource(context.TODO(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, resource, "Expected nil resource on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, resource, "Expected non-nil resource")

				// Check if data is sanitized for Secret kind
				if tc.expectSanitize {
					data, found, err := unstructured.NestedString(resource.Object, "data")
					require.NoError(t, err, "Did not expect error")
					require.True(t, found, "Expected 'data' field in Secret")
					require.Equal(t, "[redacted]", data, "Expected data to be redacted in Secret")
				}
			}
		})
	}
}

func TestInsightManager_GetYAMLForResource(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name          string
		resourceGroup *entity.ResourceGroup
		expectedYAML  []byte
		expectError   bool
	}{
		{
			name: "Success - Existing ConfigMap",
			resourceGroup: &entity.ResourceGroup{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Namespace:  "default",
				Name:       "existing-configmap",
			},
			expectedYAML: []byte(`apiVersion: v1
data:
  key1: value1
  key2: value2
kind: ConfigMap
metadata:
  name: existing-configmap
  namespace: default
`),
			expectError: false,
		},
		{
			name:          "Error - Non-existing cluster",
			resourceGroup: &entity.ResourceGroup{},
			expectError:   true,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetYAMLForResource method
			yamlData, err := manager.GetYAMLForResource(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, yamlData, "Expected nil YAML data on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, yamlData, "Expected non-nil YAML data")

				// Compare YAML data
				require.True(t, reflect.DeepEqual(tc.expectedYAML, yamlData), "YAML data does not match expected")
			}
		})
	}
}
