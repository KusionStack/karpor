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
	"errors"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
)

type ResourceGroupManager struct {
	rgrStorage storage.ResourceGroupRuleStorage
}

// NewResourceGroupManager creates a new instance of ResourceGroupManager with the given storage.
func NewResourceGroupManager(rgrStorage storage.ResourceGroupRuleStorage) (*ResourceGroupManager, error) {
	return &ResourceGroupManager{
		rgrStorage: rgrStorage,
	}, nil
}

// GetResourceGroupRule retrieves a specific resource group rule by its name.
func (m *ResourceGroupManager) GetResourceGroupRule(ctx context.Context, name string) (*entity.ResourceGroupRule, error) {
	if len(name) == 0 {
		return nil, ErrMissingResourceGroupRuleName
	}
	return m.rgrStorage.GetResourceGroupRule(ctx, name)
}

// ListResourceGroupRules returns a list of all resource group rules.
func (m *ResourceGroupManager) ListResourceGroupRules(ctx context.Context) ([]*entity.ResourceGroupRule, error) {
	return m.rgrStorage.ListResourceGroupRules(ctx)
}

// CreateResourceGroupRule creates a new resource group rule.
func (m *ResourceGroupManager) CreateResourceGroupRule(ctx context.Context, rgr *entity.ResourceGroupRule) error {
	if rgr == nil {
		return ErrNilResourceGroupRule
	}
	if len(rgr.Name) == 0 {
		return ErrMissingResourceGroupRuleName
	}

	// Check if the rule already exists to prevent duplicates.
	_, err := m.GetResourceGroupRule(ctx, rgr.Name)
	if err == nil {
		return ErrResourceGroupRuleAlreadyExists
	}
	if !errors.Is(err, elasticsearch.ErrResourceGroupRuleNotFound) {
		return err
	}
	// Save the new rule to the storage.
	return m.rgrStorage.SaveResourceGroupRule(ctx, rgr)
}

// UpdateResourceGroupRule updates an existing resource group rule.
func (m *ResourceGroupManager) UpdateResourceGroupRule(ctx context.Context, name string, rgr *entity.ResourceGroupRule) error {
	if len(name) == 0 {
		return ErrMissingResourceGroupRuleName
	}
	if rgr == nil {
		return ErrNilResourceGroupRule
	}
	if name != rgr.Name {
		return ErrResourceGroupRuleNameCannotModify
	}

	// Get the existing rule.
	existingRGR, err := m.GetResourceGroupRule(ctx, name)
	if err != nil {
		return err
	}
	if existingRGR == nil {
		return ErrResourceGroupRuleNotFound
	}

	// Update the fields of the existing rule with the new values.
	newRGR := &entity.ResourceGroupRule{
		ID:          existingRGR.ID,
		Name:        existingRGR.Name,
		Description: rgr.Description,
		Fields:      rgr.Fields,
		CreatedAt:   existingRGR.CreatedAt,
		UpdatedAt:   rgr.UpdatedAt,
		DeletedAt:   nil,
	}
	// Save the updated rule to the storage.
	return m.rgrStorage.SaveResourceGroupRule(ctx, newRGR)
}

// DeleteResourceGroupRule deletes a resource group rule by name.
func (m *ResourceGroupManager) DeleteResourceGroupRule(ctx context.Context, name string) error {
	if len(name) == 0 {
		return ErrMissingResourceGroupRuleName
	}

	// Delete the rule from the storage.
	return m.rgrStorage.DeleteResourceGroupRule(ctx, name)
}

// ListResourceGroupsBy lists all resource groups by specified resource group
// rule name.
func (m *ResourceGroupManager) ListResourceGroupsBy(ctx context.Context, ruleName string) (*storage.ResourceGroupResult, error) {
	if len(ruleName) == 0 {
		return nil, ErrMissingResourceGroupRuleName
	}

	// List the resource groups by specified rule.
	return m.rgrStorage.ListResourceGroupsBy(ctx, ruleName)
}
