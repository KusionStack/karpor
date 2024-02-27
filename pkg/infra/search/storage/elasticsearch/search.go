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

type Pagination struct {
	Page     int
	PageSize int
}

func (e *ESClient) Search(ctx context.Context, queryStr string, patternType string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	var sr *storage.SearchResult
	var err error

	switch patternType {
	case storage.DSLPatternType:
		sr, err = e.searchByDSL(ctx, queryStr, pagination)
		if err != nil {
			return nil, errors.Wrap(err, "search by DSL failed")
		}
	case storage.SQLPatternType:
		sr, err = e.searchBySQL(ctx, queryStr, pagination)
		if err != nil {
			return nil, errors.Wrap(err, "search by SQL failed")
		}
	default:
		return nil, fmt.Errorf("invalid type %s", patternType)
	}

	return sr, nil
}

func (e *ESClient) searchByDSL(ctx context.Context, dslStr string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	queries, err := Parse(dslStr)
	if err != nil {
		return nil, err
	}
	esQuery, err := ParseQueries(queries)
	if err != nil {
		return nil, err
	}
	res, err := e.SearchByQuery(ctx, esQuery, pagination)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *ESClient) searchBySQL(ctx context.Context, sqlStr string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	dsl, _, err := sql2es.Convert(sqlStr)
	if err != nil {
		return nil, err
	}
	return e.search(ctx, strings.NewReader(dsl), pagination)
}

func (e *ESClient) SearchByQuery(ctx context.Context, query map[string]interface{}, pagination *storage.Pagination) (*storage.SearchResult, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}
	return e.search(ctx, buf, pagination)
}

func (e *ESClient) search(ctx context.Context, body io.Reader, pagination *storage.Pagination) (*storage.SearchResult, error) {
	var opts []elasticsearch.Option
	if pagination != nil {
		opts = append(opts, elasticsearch.Pagination(pagination.Page, pagination.PageSize))
	}
	resp, err := e.client.SearchDocument(ctx, e.indexName, body, opts...)
	if err != nil {
		return nil, err
	}

	sr := &storage.SearchResult{
		Total:     resp.Hits.Total.Value,
		Resources: make([]*storage.Resource, len(resp.Hits.Hits)),
	}

	for i, hit := range resp.Hits.Hits {
		sr.Resources[i] = storage.Map2Resource(hit.Source)
	}

	return sr, nil
}
