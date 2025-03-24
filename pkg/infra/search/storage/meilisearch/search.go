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
	"reflect"
	"strings"

	"github.com/KusionStack/karpor/pkg/infra/persistence/meilisearch"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/pkg/errors"
)

// Pagination defines the struct for pagination which contains page number and page size.
type Pagination struct {
	Page     int
	PageSize int
}

// Search performs a search operation with the given query string, pattern type, and pagination settings.
func (s *Storage) Search(ctx context.Context, queryStr, patternType string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	var sr *storage.SearchResult
	var err error

	switch patternType {
	case storage.DSLPatternType:
		sr, err = s.searchByDSL(ctx, queryStr, pagination)
		if err != nil {
			return nil, errors.Wrap(err, "search by DSL failed")
		}
	case storage.SQLPatternType, storage.NLPatternType:
		sr, err = s.searchBySQL(ctx, queryStr, pagination)
		if err != nil {
			return nil, errors.Wrap(err, "search by SQL failed")
		}
	default:
		return nil, fmt.Errorf("invalid type %s", patternType)
	}

	return sr, nil
}

// searchByDSL performs a search operation using a DSL (Domain Specific Language) string and pagination settings.
func (s *Storage) searchByDSL(ctx context.Context, dslStr string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	request := &meilisearch.SearchRequest{}
	if pagination != nil {
		request.Limit = int64(pagination.PageSize)
		request.Offset = int64((pagination.Page - 1) * pagination.PageSize)
	}

	filter, err := Parse(dslStr)
	if err != nil {
		return nil, err
	}
	request.Filter = filter
	return s.search(ctx, request)
}

// searchBySQL performs a search operation using an SQL string and pagination settings.
func (s *Storage) searchBySQL(ctx context.Context, sqlStr string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	searchRequest, _, err := Convert(sqlStr)
	if err != nil {
		return nil, err
	}
	if pagination != nil {
		searchRequest.Limit = int64(pagination.PageSize)
		searchRequest.Offset = int64((pagination.Page - 1) * pagination.PageSize)
	}
	return s.search(ctx, searchRequest)
}

// search performs a search operation using an io.Reader as the query body and pagination settings.
func (s *Storage) search(ctx context.Context, request *meilisearch.SearchRequest) (*storage.SearchResult, error) {
	resp, err := s.client.SearchDocument(ctx, s.resourceIndexName, request)
	if err != nil {
		return nil, err
	}

	return convertSearchResult(resp)
}

// SearchByTerms performs a search operation with a map of keys and values and pagination information.
func (s *Storage) SearchByTerms(ctx context.Context, keysAndValues map[string]any, pagination *storage.Pagination) (*storage.SearchResult, error) {
	req := &meilisearch.SearchRequest{}
	if pagination != nil {
		req.Limit = int64(pagination.PageSize)
		req.Offset = int64((pagination.Page - 1) * pagination.PageSize)
	} else {
		req.Limit = 1000
		req.Offset = 0
	}
	if len(keysAndValues) != 0 {
		filter, err := ConvertToFilter(keysAndValues)
		if err != nil {
			return nil, err
		}
		req.Filter = filter
	}

	resp, err := s.client.SearchDocument(ctx, s.resourceIndexName, req)
	if err != nil {
		return nil, err
	}
	return convertSearchResult(resp)
}

// convertSearchResult converts an elasticsearch.SearchResponse to a storage.SearchResult.
func convertSearchResult(in *meilisearch.SearchResponse) (*storage.SearchResult, error) {
	out := &storage.SearchResult{
		Total:     int(in.EstimatedTotalHits),
		Resources: make([]*storage.Resource, len(in.Hits)),
	}

	for i, hit := range in.Hits {
		var err error
		obj := hit.(map[string]interface{})
		if out.Resources[i], err = storage.Map2Resource(obj); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// convertAggregationResult converts an elasticsearch.AggResults to a storage.AggregateResults.
func convertAggregationResult(in *meilisearch.AggResults) *storage.AggregateResults {
	buckets := make([]storage.Bucket, len(in.Buckets))
	for i := range in.Buckets {
		buckets[i] = storage.Bucket{
			Keys:  in.Buckets[i].Keys,
			Count: in.Buckets[i].Count,
		}
	}
	return &storage.AggregateResults{
		Buckets: buckets,
		Total:   in.Total,
	}
}

// AggregateByTerms performs an aggregation operation using the provided list of keys and returns the results.
func (s *Storage) AggregateByTerms(ctx context.Context, keys []string) (*storage.AggregateResults, error) {
	resp, err := s.client.AggregateDocumentByTerms(ctx, s.resourceIndexName, keys)
	if err != nil {
		return nil, err
	}
	return convertAggregationResult(resp), nil
}

// ConvertResourceGroup2Map converts a ResourceGroup to a map[string]any.
func ConvertResourceGroup2Map(rg *entity.ResourceGroup) (map[string]any, error) {
	result := make(map[string]interface{})
	v := reflect.ValueOf(rg).Elem()

	// Iterate through the fields of the struct.
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i).Interface()
		s, ok := field.Tag.Lookup("json")
		if !ok {
			return nil, fmt.Errorf("the JSON tag for the %s field in the ResourceGroup does not exist", field.Name)
		}
		ss := strings.Split(s, ",")
		if len(ss) == 0 {
			return nil, fmt.Errorf("invalid json tag: %s", s)
		}
		tag := strings.TrimSpace(ss[0])

		switch fieldValue := value.(type) {
		case map[string]string:
			// Handle the map field by iterating its keys and values.
			for key, val := range fieldValue {
				result[tag+"."+key] = val
			}
		case string:
			// For non-map fields, directly add them to the result map.
			// TODO: use pointers instead of null values to determine if a field exists
			if fieldValue != "" {
				result[tag] = value
			}
		default:
			return nil, fmt.Errorf("type %T not supported", fieldValue)
		}
	}

	return result, nil
}
