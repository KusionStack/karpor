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

package resourcegroup

import (
	"context"
	"testing"
	"time"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// mockResourceGroupRuleStorage is a mock implementation of storage.ResourceGroupRuleStorage
type mockResourceGroupRuleStorage struct {
	rules map[string]*entity.ResourceGroupRule
}

func newMockResourceGroupRuleStorage() *mockResourceGroupRuleStorage {
	return &mockResourceGroupRuleStorage{
		rules: make(map[string]*entity.ResourceGroupRule),
	}
}

func (m *mockResourceGroupRuleStorage) GetResourceGroupRule(ctx context.Context, name string) (*entity.ResourceGroupRule, error) {
	if rule, exists := m.rules[name]; exists {
		return rule, nil
	}
	return nil, elasticsearch.ErrResourceGroupRuleNotFound
}

func (m *mockResourceGroupRuleStorage) ListResourceGroupRules(ctx context.Context) ([]*entity.ResourceGroupRule, error) {
	rules := make([]*entity.ResourceGroupRule, 0, len(m.rules))
	for _, rule := range m.rules {
		rules = append(rules, rule)
	}
	return rules, nil
}

func (m *mockResourceGroupRuleStorage) SaveResourceGroupRule(ctx context.Context, rule *entity.ResourceGroupRule) error {
	m.rules[rule.Name] = rule
	return nil
}

func (m *mockResourceGroupRuleStorage) DeleteResourceGroupRule(ctx context.Context, name string) error {
	if _, exists := m.rules[name]; !exists {
		return elasticsearch.ErrResourceGroupRuleNotFound
	}
	delete(m.rules, name)
	return nil
}

func (m *mockResourceGroupRuleStorage) CountResourceGroupRules(ctx context.Context) (int, error) {
	return len(m.rules), nil
}

func (m *mockResourceGroupRuleStorage) ListResourceGroupsBy(ctx context.Context, ruleName string) (*storage.ResourceGroupResult, error) {
	if _, exists := m.rules[ruleName]; !exists {
		return nil, elasticsearch.ErrResourceGroupRuleNotFound
	}
	return &storage.ResourceGroupResult{
		Groups: []*entity.ResourceGroup{
			{
				Name: "test-group",
			},
		},
	}, nil
}

func TestNewResourceGroupManager(t *testing.T) {
	tests := []struct {
		name        string
		storage     storage.ResourceGroupRuleStorage
		expectError bool
	}{
		{
			name:        "Success",
			storage:     newMockResourceGroupRuleStorage(),
			expectError: false,
		},
		{
			name:        "Success with nil storage",
			storage:     nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewResourceGroupManager(tt.storage)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, manager)
			} else {
				require.NoError(t, err)
				require.NotNil(t, manager)
			}
		})
	}
}

func TestResourceGroupManager_GetResourceGroupRule(t *testing.T) {
	mockStorage := newMockResourceGroupRuleStorage()
	manager, err := NewResourceGroupManager(mockStorage)
	require.NoError(t, err)

	// Create a test rule
	testRule := &entity.ResourceGroupRule{
		Name:        "test-rule",
		Description: "Test Rule",
		Fields:      []string{"field1", "field2"},
		CreatedAt:   &metav1.Time{Time: time.Now()},
		UpdatedAt:   &metav1.Time{Time: time.Now()},
	}
	err = mockStorage.SaveResourceGroupRule(context.Background(), testRule)
	require.NoError(t, err)

	tests := []struct {
		name        string
		ruleName    string
		expectError bool
		expectRule  *entity.ResourceGroupRule
	}{
		{
			name:        "Success",
			ruleName:    "test-rule",
			expectError: false,
			expectRule:  testRule,
		},
		{
			name:        "Empty name",
			ruleName:    "",
			expectError: true,
			expectRule:  nil,
		},
		{
			name:        "Non-existent rule",
			ruleName:    "non-existent",
			expectError: true,
			expectRule:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := manager.GetResourceGroupRule(context.Background(), tt.ruleName)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, rule)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectRule, rule)
			}
		})
	}
}

func TestResourceGroupManager_ListResourceGroupRules(t *testing.T) {
	mockStorage := newMockResourceGroupRuleStorage()
	manager, err := NewResourceGroupManager(mockStorage)
	require.NoError(t, err)

	// Create test rules
	testRules := []*entity.ResourceGroupRule{
		{
			Name:        "test-rule-1",
			Description: "Test Rule 1",
			Fields:      []string{"field1"},
			CreatedAt:   &metav1.Time{Time: time.Now()},
			UpdatedAt:   &metav1.Time{Time: time.Now()},
		},
		{
			Name:        "test-rule-2",
			Description: "Test Rule 2",
			Fields:      []string{"field2"},
			CreatedAt:   &metav1.Time{Time: time.Now()},
			UpdatedAt:   &metav1.Time{Time: time.Now()},
		},
	}

	for _, rule := range testRules {
		err = mockStorage.SaveResourceGroupRule(context.Background(), rule)
		require.NoError(t, err)
	}

	t.Run("Success", func(t *testing.T) {
		rules, err := manager.ListResourceGroupRules(context.Background())
		require.NoError(t, err)
		require.Len(t, rules, len(testRules))
		require.ElementsMatch(t, testRules, rules)
	})
}

func TestResourceGroupManager_CreateResourceGroupRule(t *testing.T) {
	mockStorage := newMockResourceGroupRuleStorage()
	manager, err := NewResourceGroupManager(mockStorage)
	require.NoError(t, err)

	testRule := &entity.ResourceGroupRule{
		Name:        "test-rule",
		Description: "Test Rule",
		Fields:      []string{"field1", "field2"},
		CreatedAt:   &metav1.Time{Time: time.Now()},
		UpdatedAt:   &metav1.Time{Time: time.Now()},
	}

	tests := []struct {
		name        string
		rule        *entity.ResourceGroupRule
		expectError bool
	}{
		{
			name:        "Success",
			rule:        testRule,
			expectError: false,
		},
		{
			name:        "Nil rule",
			rule:        nil,
			expectError: true,
		},
		{
			name: "Empty name",
			rule: &entity.ResourceGroupRule{
				Description: "Test Rule",
				Fields:      []string{"field1"},
			},
			expectError: true,
		},
		{
			name:        "Duplicate rule",
			rule:        testRule,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.CreateResourceGroupRule(context.Background(), tt.rule)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify the rule was created
				savedRule, err := manager.GetResourceGroupRule(context.Background(), tt.rule.Name)
				require.NoError(t, err)
				require.Equal(t, tt.rule, savedRule)
			}
		})
	}
}

func TestResourceGroupManager_UpdateResourceGroupRule(t *testing.T) {
	mockStorage := newMockResourceGroupRuleStorage()
	manager, err := NewResourceGroupManager(mockStorage)
	require.NoError(t, err)

	// Create initial rule
	initialRule := &entity.ResourceGroupRule{
		Name:        "test-rule",
		Description: "Initial Description",
		Fields:      []string{"field1"},
		CreatedAt:   &metav1.Time{Time: time.Now()},
		UpdatedAt:   &metav1.Time{Time: time.Now()},
	}
	err = manager.CreateResourceGroupRule(context.Background(), initialRule)
	require.NoError(t, err)

	tests := []struct {
		name        string
		ruleName    string
		updateRule  *entity.ResourceGroupRule
		expectError bool
	}{
		{
			name:     "Success",
			ruleName: "test-rule",
			updateRule: &entity.ResourceGroupRule{
				Name:        "test-rule",
				Description: "Updated Description",
				Fields:      []string{"field1", "field2"},
				UpdatedAt:   &metav1.Time{Time: time.Now()},
			},
			expectError: false,
		},
		{
			name:        "Empty name",
			ruleName:    "",
			updateRule:  initialRule,
			expectError: true,
		},
		{
			name:        "Nil rule",
			ruleName:    "test-rule",
			updateRule:  nil,
			expectError: true,
		},
		{
			name:     "Name mismatch",
			ruleName: "test-rule",
			updateRule: &entity.ResourceGroupRule{
				Name: "different-name",
			},
			expectError: true,
		},
		{
			name:     "Non-existent rule",
			ruleName: "non-existent",
			updateRule: &entity.ResourceGroupRule{
				Name: "non-existent",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.UpdateResourceGroupRule(context.Background(), tt.ruleName, tt.updateRule)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify the update
				updatedRule, err := manager.GetResourceGroupRule(context.Background(), tt.ruleName)
				require.NoError(t, err)
				require.Equal(t, tt.updateRule.Description, updatedRule.Description)
				require.Equal(t, tt.updateRule.Fields, updatedRule.Fields)
			}
		})
	}
}

func TestResourceGroupManager_DeleteResourceGroupRule(t *testing.T) {
	mockStorage := newMockResourceGroupRuleStorage()
	manager, err := NewResourceGroupManager(mockStorage)
	require.NoError(t, err)

	// Create a test rule
	testRule := &entity.ResourceGroupRule{
		Name:        "test-rule",
		Description: "Test Rule",
		Fields:      []string{"field1"},
		CreatedAt:   &metav1.Time{Time: time.Now()},
		UpdatedAt:   &metav1.Time{Time: time.Now()},
	}
	err = manager.CreateResourceGroupRule(context.Background(), testRule)
	require.NoError(t, err)

	tests := []struct {
		name        string
		ruleName    string
		expectError bool
	}{
		{
			name:        "Success",
			ruleName:    "test-rule",
			expectError: false,
		},
		{
			name:        "Empty name",
			ruleName:    "",
			expectError: true,
		},
		{
			name:        "Non-existent rule",
			ruleName:    "non-existent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.DeleteResourceGroupRule(context.Background(), tt.ruleName)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify the rule was deleted
				_, err := manager.GetResourceGroupRule(context.Background(), tt.ruleName)
				require.Error(t, err)
				require.ErrorIs(t, err, elasticsearch.ErrResourceGroupRuleNotFound)
			}
		})
	}
}

func TestResourceGroupManager_ListResourceGroupsBy(t *testing.T) {
	mockStorage := newMockResourceGroupRuleStorage()
	manager, err := NewResourceGroupManager(mockStorage)
	require.NoError(t, err)

	// Create a test rule
	testRule := &entity.ResourceGroupRule{
		Name:        "test-rule",
		Description: "Test Rule",
		Fields:      []string{"field1"},
		CreatedAt:   &metav1.Time{Time: time.Now()},
		UpdatedAt:   &metav1.Time{Time: time.Now()},
	}
	err = manager.CreateResourceGroupRule(context.Background(), testRule)
	require.NoError(t, err)

	tests := []struct {
		name        string
		ruleName    string
		expectError bool
	}{
		{
			name:        "Success",
			ruleName:    "test-rule",
			expectError: false,
		},
		{
			name:        "Empty name",
			ruleName:    "",
			expectError: true,
		},
		{
			name:        "Non-existent rule",
			ruleName:    "non-existent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.ListResourceGroupsBy(context.Background(), tt.ruleName)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotEmpty(t, result.Groups)
			}
		})
	}
}
