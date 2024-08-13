package insight

import (
	"context"
	"testing"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/dynamic"
)

func TestInsightManager_GetResourceEvents(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockEventResource{}).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		resourceGroup  *entity.ResourceGroup
		expectedLength int
		expectError    bool
	}{
		{
			name: "Success - GetResourceEvents",
			resourceGroup: &entity.ResourceGroup{
				Name:       "default-name",
				APIVersion: "v1",
				Kind:       "Pod",
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "Failed - GetResourceEvents",
			resourceGroup: &entity.ResourceGroup{
				Name:       "default-test",
				APIVersion: "v1",
				Kind:       "Pod",
			},
			expectedLength: 0,
			expectError:    false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetResourceGroupSummary method
			events, err := manager.GetResourceEvents(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			}
		})
	}
}

func TestInsightManager_GetNamespaceGVKEvents(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockEventResource{}).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		resourceGroup  *entity.ResourceGroup
		expectedLength int
		expectError    bool
	}{
		{
			name: "Success - GetNamespaceGVKEvents",
			resourceGroup: &entity.ResourceGroup{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "Failed - GetNamespaceGVKEvents",
			resourceGroup: &entity.ResourceGroup{
				APIVersion: "v2",
				Kind:       "Pod",
			},
			expectedLength: 0,
			expectError:    false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetResourceGroupSummary method
			events, err := manager.GetNamespaceGVKEvents(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			}
		})
	}
}

func TestInsightManager_GetNamespaceEvents(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockEventResource{}).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		resourceGroup  *entity.ResourceGroup
		expectedLength int
		expectError    bool
	}{
		{
			name: "Success - GetNamespaceEvents",
			resourceGroup: &entity.ResourceGroup{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "Failed - GetNamespaceEvents",
			resourceGroup: &entity.ResourceGroup{
				APIVersion: "v2",
				Kind:       "Pod",
				Namespace:  "non-default",
			},
			expectedLength: 0,
			expectError:    false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetResourceGroupSummary method
			events, err := manager.GetNamespaceEvents(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			}
		})
	}
}

func TestInsightManager_GetGVKEvents(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockEventResource{}).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		resourceGroup  *entity.ResourceGroup
		expectedLength int
		expectError    bool
	}{
		{
			name: "Success - GetGVKEvents",
			resourceGroup: &entity.ResourceGroup{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "Failed - GetGVKEvents",
			resourceGroup: &entity.ResourceGroup{
				APIVersion: "v2",
				Kind:       "Pod",
				Namespace:  "non-default",
			},
			expectedLength: 0,
			expectError:    false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetResourceGroupSummary method
			events, err := manager.GetGVKEvents(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			}
		})
	}
}

func TestInsightManager_GetClusterEvents(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Set up mocks for dynamic client
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockEventResource{}).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		resourceGroup  *entity.ResourceGroup
		expectedLength int
		expectError    bool
	}{
		{
			name: "Success - GetGVKEvents",
			resourceGroup: &entity.ResourceGroup{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
			},
			expectedLength: 1,
			expectError:    false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetResourceGroupSummary method
			events, err := manager.GetClusterEvents(context.Background(), mockMultiClusterClient(), tc.resourceGroup)

			// Check error expectation
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.Len(t, events, tc.expectedLength, "Expected nil result on error")
			}
		})
	}
}
