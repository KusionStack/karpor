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
	"testing"

	"github.com/KusionStack/karbour/pkg/core/entity"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/dynamic"
)

func TestInsightManager_GetResourceSummary(t *testing.T) {
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
		expectedResult *ResourceSummary
		expectError    bool
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
			expectedResult: &ResourceSummary{
				Resource: entity.ResourceGroup{
					Name:       "existing-configmap",
					Namespace:  "default",
					APIVersion: "v1",
					Cluster:    "existing-cluster",
					Kind:       "ConfigMap",
				},
			},
			expectError: false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetResourceSummary method
			result, err := manager.GetResourceSummary(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, result, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, result, "Expected non-nil result")

				// Compare results
				require.Equal(t, tc.expectedResult, result, "Result does not match expected")
			}
		})
	}
}

func TestInsightManager_GetGVKSummary(t *testing.T) {
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
		expectedResult *GVKSummary
		expectError    bool
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
			expectedResult: &GVKSummary{
				Cluster: "existing-cluster",
				Group:   "",
				Version: "v1",
				Kind:    "ConfigMap",
				Count:   1,
			},
			expectError: false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetGVKSummary method
			result, err := manager.GetGVKSummary(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, result, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, result, "Expected non-nil result")

				// Compare results
				require.Equal(t, tc.expectedResult, result, "Result does not match expected")
			}
		})
	}
}

func TestInsightManager_GetNamespaceSummary(t *testing.T) {
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
		expectedResult *NamespaceSummary
		expectError    bool
	}{
		{
			name: "Success - Existing Namespace",
			resourceGroup: &entity.ResourceGroup{
				Cluster:   "existing-cluster",
				Namespace: "default",
			},
			expectedResult: &NamespaceSummary{
				Cluster:   "existing-cluster",
				Namespace: "default",
				CountByGVK: map[string]int{
					"Pod.v1": 1,
				},
			},
			expectError: false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetNamespaceSummary method
			result, err := manager.GetNamespaceSummary(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, result, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, result, "Expected non-nil result")

				// Compare results
				require.Equal(t, tc.expectedResult, result, "Result does not match expected")
			}
		})
	}
}
