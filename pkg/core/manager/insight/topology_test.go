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
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/topology"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/dynamic"
)

func TestInsightManager_GetTopologyForCluster(t *testing.T) {
	// Set up environment variable for relationship file
	setRelationshipFilePath()
	defer os.Unsetenv("KARPOR_RELATIONSHIP_FILE")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	mockey.Mock(topology.GVRNamespaced).Return(true).Build()
	defer mockey.UnPatchAll()

	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Test cases
	testCases := []struct {
		name        string
		cluster     string
		noCache     bool
		expectedMap map[string]ClusterTopology
		expectError bool
	}{
		{
			name:        "Success - Cache Hit",
			cluster:     "existing-cluster",
			noCache:     false,
			expectedMap: mockClusterTopologyMapForCluster(),
			expectError: false,
		},
		{
			name:        "Success - Cache Miss",
			cluster:     "existing-cluster",
			noCache:     true,
			expectedMap: mockClusterTopologyMapForCluster(),
			expectError: false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetTopologyForCluster method
			topologyMap, err := manager.GetTopologyForCluster(context.Background(), mockMultiClusterClient(), tc.cluster, tc.noCache)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, topologyMap, "Expected nil topology map on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, topologyMap, "Expected non-nil topology map")

				// Compare topology map
				require.Equal(t, tc.expectedMap, topologyMap, "Topology map does not match expected")
			}
		})
	}
}

func TestInsightManager_GetTopologyForResource(t *testing.T) {
	// Set up environment variable for relationship file
	setRelationshipFilePath()
	defer os.Unsetenv("KARPOR_RELATIONSHIP_FILE")

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
		noCache       bool
		expectedMap   map[string]ResourceTopology
		expectError   bool
	}{
		{
			name: "Success - Existing Pod",
			resourceGroup: &entity.ResourceGroup{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
				Name:       "existing-pod",
			},
			noCache: true,
			expectedMap: map[string]ResourceTopology{
				"/v1.Pod:default.existing-pod": {
					ResourceGroup: entity.ResourceGroup{
						Cluster:    "existing-cluster",
						APIVersion: "v1",
						Kind:       "Pod",
						Namespace:  "default",
						Name:       "existing-pod",
					},
					Parents:  []string{},
					Children: []string{},
				},
			},
			expectError: false,
		},
		{
			name:          "Error - Non-existing cluster",
			resourceGroup: &entity.ResourceGroup{},
			noCache:       true,
			expectError:   true,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetTopologyForResource method
			topologyMap, err := manager.GetTopologyForResource(context.Background(), mockMultiClusterClient(), tc.resourceGroup, tc.noCache)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, topologyMap, "Expected nil topology map on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, topologyMap, "Expected non-nil topology map")

				// Compare topology map
				require.True(t, reflect.DeepEqual(tc.expectedMap, topologyMap), "Topology map does not match expected")
			}
		})
	}
}

func TestInsightManager_GetTopologyForClusterNamespace(t *testing.T) {
	// Set up environment variable for relationship file
	setRelationshipFilePath()
	defer os.Unsetenv("KARPOR_RELATIONSHIP_FILE")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	mockey.Mock(topology.GVRNamespaced).Return(true).Build()
	defer mockey.UnPatchAll()

	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Test cases
	testCases := []struct {
		name        string
		cluster     string
		namespace   string
		noCache     bool
		expectedMap map[string]ClusterTopology
		expectError bool
	}{
		{
			name:        "Success - Cache Hit",
			cluster:     "existing-cluster",
			namespace:   "default",
			noCache:     false,
			expectedMap: mockClusterTopologyMapForClusterNamespace(),
			expectError: false,
		},
		{
			name:        "Success - Cache Miss",
			cluster:     "existing-cluster",
			namespace:   "default",
			noCache:     true,
			expectedMap: mockClusterTopologyMapForClusterNamespace(),
			expectError: false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetTopologyForClusterNamespace method
			topologyMap, err := manager.GetTopologyForClusterNamespace(context.Background(), mockMultiClusterClient(), tc.cluster, tc.namespace, tc.noCache)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, topologyMap, "Expected nil topology map on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, topologyMap, "Expected non-nil topology map")

				// Compare topology map
				require.True(t, reflect.DeepEqual(tc.expectedMap, topologyMap), "Topology map does not match expected")
			}
		})
	}
}

func setRelationshipFilePath() {
	filePath := os.Getenv("KARPOR_RELATIONSHIP_FILE")
	if filePath == "" {
		// Default file path
		filePath = "relationship.yaml"
		// Find the file in each parent directory
		curdir, err := os.Getwd()
		if err != nil {
			panic("Unable to get current directory")
		}
		dir := curdir
		for {
			configPath := filepath.Join(dir, "config", "relationship.yaml")
			if _, err := os.Stat(configPath); err == nil {
				filePath, err = filepath.Rel(curdir, configPath)
				if err != nil {
					panic("Failed to get relative path of relationship.yaml as " + err.Error())
				}
				break
			}
			// Move up to the parent directory
			parent := filepath.Dir(dir)
			if parent == dir {
				// Reached the root directory, break the loop
				break
			}
			dir = parent
		}
		// Set the environment variable
		if err := os.Setenv("KARPOR_RELATIONSHIP_FILE", filePath); err != nil {
			panic("Failed to set environment variable")
		}
	}
}
