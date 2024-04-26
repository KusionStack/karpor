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

package entity

import (
	"fmt"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	CreatedAt *metav1.Time `yaml:"createdAt,omitempty" json:"createdAt,omitempty"`
	// CreatedAt is the timestamp of the updated for the resourceGroupRule.
	UpdatedAt *metav1.Time `yaml:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	// DeletedAt is the timestamp of the deleted for the resourceGroupRule.
	DeletedAt *metav1.Time `yaml:"deletedAt,omitempty" json:"deletedAt,omitempty"`
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

// UUID() returns a randomly generated UUID string.
func UUID() string {
	uuid := uuid.New()
	return uuid.String()
}
