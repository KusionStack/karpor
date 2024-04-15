package repository

import (
	"context"

	"github.com/KusionStack/karbour/pkg/core/entity"
)

// ResourceGroupRuleRepository is an interface that defines the repository
// operations for resource group rule.
// It follows the principles of domain-driven design (DDD).
type ResourceGroupRuleRepository interface {
	// Create creates a new resourceGroupRule.
	Create(ctx context.Context, resourceGroupRule *entity.ResourceGroupRule) error
	// Delete deletes a resourceGroupRule by its ID.
	Delete(ctx context.Context, id string) error
	// Update updates an existing resourceGroupRule.
	Update(ctx context.Context, resourceGroupRule *entity.ResourceGroupRule) error
	// Get retrieves a resourceGroupRule by its ID.
	Get(ctx context.Context, id string) (*entity.ResourceGroupRule, error)
}
