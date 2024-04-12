// Copyright The Karbour Authors.
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
	"strings"

	"github.com/KusionStack/karbour/pkg/infra/persistence/elasticsearch"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/util/sql2es"
	"github.com/pkg/errors"
)

// Pagination defines the struct for pagination which contains page number and page size.
type Pagination struct {
	Page     int
	PageSize int
}

// Search performs a search operation with the given query string, pattern type, and pagination settings.
func (s *Storage) Search(ctx context.Context, queryStr string, patternType string, pagination *storage.Pagination) (*storage.SearchResult, error) {
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
	dsl, _, err := sql2es.Convert(sqlStr)
	if err != nil {
		return nil, err
	}
	return s.search(ctx, strings.NewReader(dsl), pagination)
}

// SearchByQuery performs a search operation using a query map and pagination settings.
func (s *Storage) SearchByQuery(ctx context.Context, query map[string]interface{}, pagination *storage.Pagination) (*storage.SearchResult, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}
	return s.search(ctx, buf, pagination)
}

// search performs a search operation using an io.Reader as the query body and pagination settings.
func (s *Storage) search(ctx context.Context, body io.Reader, pagination *storage.Pagination) (*storage.SearchResult, error) {
	var opts []elasticsearch.Option
	if pagination != nil {
		opts = append(opts, elasticsearch.Pagination(pagination.Page, pagination.PageSize))
	}
	resp, err := s.client.SearchDocument(ctx, s.indexName, body, opts...)
	if err != nil {
		return nil, err
	}

	sr := &storage.SearchResult{
		Total:     resp.Hits.Total.Value,
		Resources: make([]*storage.Resource, len(resp.Hits.Hits)),
	}

	for i, hit := range resp.Hits.Hits {
		sr.Resources[i], err = storage.Map2Resource(hit.Source)
		if err != nil {
			return nil, err
		}
	}

	return sr, nil
}
