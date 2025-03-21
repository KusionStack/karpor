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

package meilisearch

import (
	"context"
	"fmt"
	"github.com/KusionStack/karpor/pkg/infra/persistence/meilisearch"
	"strings"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
)

const (
	resourceGroupRuleKeyID          = "id"
	resourceGroupRuleKeyName        = "name"
	resourceGroupRuleKeyDescription = "description"
	resourceGroupRuleKeyFields      = "fields"
	resourceGroupRuleKeyCreatedAt   = "createdAt"
	resourceGroupRuleKeyUpdatedAt   = "updatedAt"
	resourceGroupRuleKeyDeletedAt   = "deletedAt"
)

// DeleteResourceGroupRule deletes a resource group rule based on the given name.
func (s *Storage) DeleteResourceGroupRule(ctx context.Context, name string) error {

	if rgr, err := s.GetResourceGroupRule(ctx, name); err != nil {
		return err
	} else {
		return s.client.DeleteDocument(ctx, s.resourceGroupRuleIndexName, rgr.ID)
	}
}

// GetResourceGroupRule retrieves a resource group rule based on the given name.
func (s *Storage) GetResourceGroupRule(ctx context.Context, name string) (*entity.ResourceGroupRule, error) {

	filter := generateFilter(resourceKeyName, name)

	resp, err := s.client.SearchDocument(ctx, s.resourceGroupRuleIndexName, &meilisearch.SearchRequest{Filter: filter})
	if err != nil {
		return nil, err
	}

	if resp.TotalHits == 0 {
		return nil, storage.ErrResourceGroupRuleNotFound
	}

	res, err := storage.Map2ResourceGroupRule(resp.Hits[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ListResourceGroupRules lists all resource group rules by searching the entire
// index.
func (s *Storage) ListResourceGroupRules(ctx context.Context) ([]*entity.ResourceGroupRule, error) {

	// Execute the search document call to the storage.
	resp, err := s.client.SearchDocument(ctx, s.resourceGroupRuleIndexName, &meilisearch.SearchRequest{})
	if err != nil {
		return nil, err
	}

	// Check if the search found any resource group rules.
	if resp.TotalHits == 0 {
		return nil, storage.ErrResourceGroupRuleNotFound
	}

	// Initialize a slice to hold the resource group rules.
	rgrList := make([]*entity.ResourceGroupRule, 0, len(resp.Hits))

	// Iterate over the search hits and map each hit to a ResourceGroupRule entity.
	for _, hit := range resp.Hits {
		// Map the source of the hit to a ResourceGroupRule entity.
		rgr, err := storage.Map2ResourceGroupRule(hit.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		rgrList = append(rgrList, rgr)
	}

	return rgrList, nil
}

// ListResourceGroupsBy lists all resource groups by specified resource group
// rule name.
func (s *Storage) ListResourceGroupsBy(ctx context.Context, ruleName string) (*storage.ResourceGroupResult, error) {

	rgr, err := s.GetResourceGroupRule(ctx, ruleName)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.AggregateDocumentByTerms(ctx, s.resourceIndexName, rgr.Fields)
	if err != nil {
		return nil, err
	}

	// Check if the search found any resource groups.
	if resp.Total == 0 {
		return nil, storage.ErrResourceGroupNotFound
	}

	// Initialize a slice to hold the resource group rules.
	rgList := make([]*entity.ResourceGroup, 0, len(resp.Buckets))

	// Iterate over the search hits and map each hit to a ResourceGroupRule entity.
	for _, bucket := range resp.Buckets {
		if len(rgr.Fields) != len(bucket.Keys) {
			return nil, fmt.Errorf("mismatched number of fields: expected %d, got %d", len(rgr.Fields), len(bucket.Keys))
		}
		// Convert the current bucket to a resource group.
		rg := &entity.ResourceGroup{}
		for i, v := range bucket.Keys {
			field := rgr.Fields[i]
			switch field {
			case "cluster":
				rg.Cluster = v
			case "apiVersion":
				rg.APIVersion = v
			case "kind":
				rg.Kind = v
			case "namespace":
				rg.Namespace = v
			case "name":
				rg.Name = v
			default:
				if strings.HasPrefix(field, "annotations.") {
					annoKey := strings.TrimPrefix(field, "annotations.")
					if rg.Annotations == nil {
						rg.Annotations = map[string]string{annoKey: v}
					} else {
						rg.Annotations[annoKey] = v
					}
				} else if strings.HasPrefix(field, "labels.") {
					labelKey := strings.TrimPrefix(field, "labels.")
					if rg.Labels == nil {
						rg.Labels = map[string]string{labelKey: v}
					} else {
						rg.Labels[labelKey] = v
					}
				}
			}
		}
		rgList = append(rgList, rg)
	}

	return &storage.ResourceGroupResult{
		Groups: rgList,
		Fields: rgr.Fields,
	}, nil
}

// SaveResourceGroupRule saves a resource group rule to the storage.
func (s *Storage) SaveResourceGroupRule(ctx context.Context, data *entity.ResourceGroupRule) error {
	obj := s.generateResourceGroupRuleDocument(data)
	if err := s.client.UpdateIndexFacets(ctx, s.resourceIndexName, data.Fields); err != nil {
		return err
	}
	return s.client.SaveDocument(ctx, s.resourceGroupRuleIndexName, obj)
}

// CountResourceGroupRules return a count of resource group rules in the
// Elasticsearch storage.
func (s *Storage) CountResourceGroupRules(ctx context.Context) (int, error) {
	if resp, err := s.client.Count(ctx, s.resourceGroupRuleIndexName); err != nil {
		return 0, err
	} else {
		return int(resp.Count), nil
	}
}

// generateResourceGroupRuleDocument creates a resource group rule document for
// Elasticsearch with the specified name, description etc.
func (s *Storage) generateResourceGroupRuleDocument(data *entity.ResourceGroupRule) map[string]any {
	var id string
	if len(data.ID) == 0 {
		id = entity.UUID()
	} else {
		id = data.ID
	}
	return map[string]any{
		resourceGroupRuleKeyID:          id,
		resourceGroupRuleKeyName:        data.Name,
		resourceGroupRuleKeyDescription: data.Description,
		resourceGroupRuleKeyFields:      data.Fields,
		resourceGroupRuleKeyCreatedAt:   data.CreatedAt,
		resourceGroupRuleKeyUpdatedAt:   data.UpdatedAt,
	}
}

// generateResourceGroupRuleFilter creates a query to search for an object in
// Elasticsearch based on resource group rule's name.
func generateFilter(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}
