package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ResourceGroupRule struct {
	// ID is the id of the resourceGroupRule.
	ID string `yaml:"id" json:"id"`
	// Name is the name of the resourceGroupRule.
	Name string `yaml:"name" json:"name"`
	// Description is a human-readable description of the resourceGroupRule.
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Fields      []string `yaml:"fields,omitempty" json:"fields,omitempty"`
	// CreatedAt is the timestamp of the created for the resourceGroupRule.
	CreatedAt time.Time `yaml:"createdAt,omitempty" json:"createdAt,omitempty"`
	// CreatedAt is the timestamp of the updated for the resourceGroupRule.
	UpdatedAt time.Time `yaml:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	// DeletedAt is the timestamp of the deleted for the resourceGroupRule.
	DeletedAt time.Time `yaml:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// Validate checks if the resourceGroupRule is valid.
// It returns an error if the resourceGroupRule is not valid.
func (e *ResourceGroupRule) Validate() error {
	if e == nil {
		return fmt.Errorf("resource group rule is nil")
	}

	if e.Name == "" {
		return fmt.Errorf("resource group rule must have a name")
	}

	return nil
}

func UUID() string {
	uuid := uuid.New()
	return uuid.String()
}
