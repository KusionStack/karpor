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

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/persistence/elasticsearch"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/util/sql2es"
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
	case storage.SQLPatternType:
		sr, err = s.searchBySQL(ctx, queryStr, pagination)
		if err != nil {
			return nil, errors.Wrap(err, "search by SQL failed")
		}
	default:
		return nil, fmt.Errorf("invalid type %s", patternType)
	}

	return sr, nil
}

// SearchByQuery performs a search operation using a query map and pagination settings.
func (s *Storage) SearchByQuery(ctx context.Context, query map[string]interface{}, pagination *storage.Pagination) (*storage.SearchResult, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}
	return s.search(ctx, buf, pagination)
}

// searchByDSL performs a search operation using a DSL (Domain Specific Language) string and pagination settings.
func (s *Storage) searchByDSL(ctx context.Context, dslStr string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	queries, err := Parse(dslStr)
	if err != nil {
		return nil, err
	}
	esQuery, err := ParseQueries(queries)
	if err != nil {
		return nil, err
	}
	res, err := s.SearchByQuery(ctx, esQuery, pagination)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// searchBySQL performs a search operation using an SQL string and pagination settings.
func (s *Storage) searchBySQL(ctx context.Context, sqlStr string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	dsl, _, err := sql2es.ConvertWithDefaultFilter(sqlStr, &sql2es.DeletedFilter)
	if err != nil {
		return nil, err
	}
	return s.search(ctx, strings.NewReader(dsl), pagination)
}

// search performs a search operation using an io.Reader as the query body and pagination settings.
func (s *Storage) search(ctx context.Context, body io.Reader, pagination *storage.Pagination) (*storage.SearchResult, error) {
	var opts []elasticsearch.Option
	if pagination != nil {
		opts = append(opts, elasticsearch.Pagination(pagination.Page, pagination.PageSize))
	}
	resp, err := s.client.SearchDocument(ctx, s.resourceIndexName, body, opts...)
	if err != nil {
		return nil, err
	}

	return convertSearchResult(resp)
}

// SearchByTerms performs a search operation with a map of keys and values and pagination information.
func (s *Storage) SearchByTerms(ctx context.Context, keysAndValues map[string]any, pagination *storage.Pagination) (*storage.SearchResult, error) {
	var opts []elasticsearch.Option
	if pagination != nil {
		opts = append(opts, elasticsearch.Pagination(pagination.Page, pagination.PageSize))
	}
	resp, err := s.client.SearchDocumentByTerms(ctx, s.resourceIndexName, keysAndValues, opts...)
	if err != nil {
		return nil, err
	}
	return convertSearchResult(resp)
}

// convertSearchResult converts an elasticsearch.SearchResponse to a storage.SearchResult.
func convertSearchResult(in *elasticsearch.SearchResponse) (*storage.SearchResult, error) {
	out := &storage.SearchResult{
		Total:     in.Hits.Total.Value,
		Resources: make([]*storage.Resource, len(in.Hits.Hits)),
	}

	for i, hit := range in.Hits.Hits {
		var err error

		if out.Resources[i], err = storage.Map2Resource(hit.Source); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// convertAggregationResult converts an elasticsearch.AggResults to a storage.AggregateResults.
func convertAggregationResult(in *elasticsearch.AggResults) *storage.AggregateResults {
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
