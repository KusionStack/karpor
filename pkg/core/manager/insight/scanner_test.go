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
	"testing"

	"kusionstack.io/karpor/pkg/core/entity"
	"github.com/stretchr/testify/require"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

func TestInsightManager_Audit(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Test cases
	testCases := []struct {
		name          string
		resourceGroup entity.ResourceGroup
		noCache       bool
		expectError   bool
		expectedCount int
	}{
		{
			name: "Audit with Cache Enabled",
			resourceGroup: entity.ResourceGroup{
				Cluster:    "test-cluster",
				Namespace:  "test-namespace",
				APIVersion: "v1",
				Kind:       "Pod",
				Name:       "test-pod",
			},
			noCache:       false,
			expectError:   false,
			expectedCount: 9,
		},
		{
			name: "Audit with Cache Disabled",
			resourceGroup: entity.ResourceGroup{
				Cluster:    "test-cluster",
				Namespace:  "test-namespace",
				APIVersion: "v1",
				Kind:       "Deployment",
				Name:       "test-deployment",
			},
			noCache:       true,
			expectError:   false,
			expectedCount: 9,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call Audit method
			result, err := manager.Audit(context.Background(), tc.resourceGroup, tc.noCache)

			// Check for errors
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, result, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, result, "Expected non-nil result")

				// Check the number of issues
				require.Equal(t, tc.expectedCount, result.IssueTotal(), "Unexpected number of issues")
			}
		})
	}
}

func TestInsightManager_Score(t *testing.T) {
	// Initialize InsightManager
	manager, err := NewInsightManager(&mockSearchStorage{}, &mockResourceStorage{}, &mockResourceGroupRuleStorage{}, &genericapiserver.CompletedConfig{})
	require.NoError(t, err, "Unexpected error initializing InsightManager")

	// Test cases
	testCases := []struct {
		name                      string
		resourceGroup             entity.ResourceGroup
		noCache                   bool
		expectScore               int
		expectResourceTotal       int
		expectIssuesTotal         int
		expectHighSeverityCount   int
		expectMediumSeverityCount int
		expectLowSeverityCount    int
		expectError               bool
	}{
		{
			name:                      "Score Calculation with Cache Enabled",
			resourceGroup:             entity.ResourceGroup{Cluster: "existing-cluster", APIVersion: "v1", Kind: "Pod", Namespace: "default", Name: "existing-pod"},
			noCache:                   false,
			expectScore:               17,
			expectResourceTotal:       1,
			expectIssuesTotal:         9,
			expectHighSeverityCount:   7,
			expectMediumSeverityCount: 0,
			expectLowSeverityCount:    2,
			expectError:               false,
		},
		{
			name:                      "Score Calculation with Cache Disabled",
			resourceGroup:             entity.ResourceGroup{Cluster: "existing-cluster", APIVersion: "v1", Kind: "Pod", Namespace: "default", Name: "existing-pod"},
			noCache:                   true,
			expectScore:               17,
			expectResourceTotal:       1,
			expectIssuesTotal:         9,
			expectHighSeverityCount:   7,
			expectMediumSeverityCount: 0,
			expectLowSeverityCount:    2,
			expectError:               false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call Score method
			scoreData, err := manager.Score(context.Background(), tc.resourceGroup, tc.noCache)

			// Check for errors
			if tc.expectError {
				require.Error(t, err, "Expected an error")
				require.Nil(t, scoreData, "Expected nil score data on error")
			} else {
				require.NoError(t, err, "Did not expect error")
				require.NotNil(t, scoreData, "Expected non-nil score data")

				// Check the score data
				require.NotNil(t, scoreData.Score, "Expected non-nil score")
				require.Equal(t, tc.expectScore, int(scoreData.Score))
				require.Equal(t, tc.expectResourceTotal, scoreData.ResourceTotal)
				require.Equal(t, tc.expectIssuesTotal, scoreData.IssuesTotal)
				require.Equal(t, tc.expectHighSeverityCount, scoreData.SeverityStatistic["High"])
				require.Equal(t, tc.expectMediumSeverityCount, scoreData.SeverityStatistic["Medium"])
				require.Equal(t, tc.expectLowSeverityCount, scoreData.SeverityStatistic["Low"])
			}
		})
	}
}
