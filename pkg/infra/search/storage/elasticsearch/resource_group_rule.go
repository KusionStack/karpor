package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/KusionStack/karbour/pkg/core/entity"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/aquasecurity/esquery"
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

var ErrResourceGroupRuleNotFound = fmt.Errorf("resource group rule not found")

func (s *Storage) DeleteResourceGroupRule(ctx context.Context, name string) error {
	if rgr, err := s.GetResourceGroupRule(ctx, name); err != nil {
		return err
	} else {
		return s.client.DeleteDocument(ctx, s.resourceGroupRuleIndexName, rgr.ID)
	}
}

func (s *Storage) GetResourceGroupRule(ctx context.Context, name string) (*entity.ResourceGroupRule, error) {
	query := generateResourceGroupRuleQuery(name)
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}
	resp, err := s.client.SearchDocument(ctx, s.resourceGroupRuleIndexName, buf)
	if err != nil {
		return nil, err
	}

	if resp.Hits.Total.Value == 0 {
		return nil, ErrResourceGroupRuleNotFound
	}

	res, err := storage.Map2ResourceGroupRule(resp.Hits.Hits[0].Source)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ListResourceGroupRules lists all resource group rules by searching the entire
// index.
func (s *Storage) ListResourceGroupRules(ctx context.Context) ([]*entity.ResourceGroupRule, error) {
	// Create a query to search for all resource group rules.
	query := generateResourceGroupRuleQueryForAll()

	// Buffer to hold the query JSON.
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}

	// Execute the search document call to the storage.
	resp, err := s.client.SearchDocument(ctx, s.resourceGroupRuleIndexName, buf)
	if err != nil {
		return nil, err
	}

	// Check if the search found any resource group rules.
	if resp.Hits.Total.Value == 0 {
		return nil, ErrResourceGroupRuleNotFound
	}

	// Initialize a slice to hold the resource group rules.
	var rgrList []*entity.ResourceGroupRule

	// Iterate over the search hits and map each hit to a ResourceGroupRule entity.
	for _, hit := range resp.Hits.Hits {
		// Map the source of the hit to a ResourceGroupRule entity.
		rgr, err := storage.Map2ResourceGroupRule(hit.Source)
		if err != nil {
			return nil, err
		}
		rgrList = append(rgrList, rgr)
	}

	return rgrList, nil
}

func (s *Storage) SaveResourceGroupRule(ctx context.Context, data *entity.ResourceGroupRule) error {
	id, body, err := s.generateResourceGroupRuleDocument(data)
	if err != nil {
		return err
	}
	return s.client.SaveDocument(ctx, s.resourceGroupRuleIndexName, id, bytes.NewReader(body))
}

// generateResourceGroupRuleDocument creates an resource group rule document for
// Elasticsearch with the specified name, description etc.
func (s *Storage) generateResourceGroupRuleDocument(data *entity.ResourceGroupRule) (id string, body []byte, err error) {
	if len(data.ID) == 0 {
		id = entity.UUID()
	} else {
		id = data.ID
	}
	body, err = json.Marshal(map[string]interface{}{
		resourceGroupRuleKeyID:          id,
		resourceGroupRuleKeyName:        data.Name,
		resourceGroupRuleKeyDescription: data.Description,
		resourceGroupRuleKeyFields:      data.Fields,
		resourceGroupRuleKeyCreatedAt:   data.CreatedAt,
		resourceGroupRuleKeyUpdatedAt:   data.UpdatedAt,
	})
	if err != nil {
		return
	}
	return
}

// generateResourceGroupRuleQuery creates a query to search for an object in
// Elasticsearch based on resource group rule's name.
func generateResourceGroupRuleQuery(name string) map[string]interface{} {
	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term(resourceKeyName, name),
	).Map()
	return query
}

// generateResourceGroupRuleQueryForAll creates a query to search for all
// resource group rules.
func generateResourceGroupRuleQueryForAll() map[string]interface{} {
	query := make(map[string]interface{})
	// This query will match all documents in the index.
	query["query"] = esquery.MatchAll().Map()
	return query
}
