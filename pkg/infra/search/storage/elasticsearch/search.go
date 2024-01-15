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

	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/cch123/elasticsql"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/pkg/errors"
)

type Pagination struct {
	Page     int
	PageSize int
}

func (s *ESClient) Search(ctx context.Context, queryStr string, patternType string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	var res *SearchResponse
	var err error

	switch patternType {
	case storage.DSLPatternType:
		res, err = s.searchByDSL(ctx, queryStr, pagination)
		if err != nil {
			return nil, errors.Wrap(err, "search by DSL failed")
		}
	case storage.SQLPatternType:
		res, err = s.searchBySQL(ctx, queryStr, pagination)
		if err != nil {
			return nil, errors.Wrap(err, "search by SQL failed")
		}
	default:
		return nil, fmt.Errorf("invalid type %s", patternType)
	}

	rt := &storage.SearchResult{
		Total:     res.Hits.Total.Value,
		Resources: make([]*storage.Resource, len(res.Hits.Hits)),
	}

	for i, hit := range res.Hits.Hits {
		rt.Resources[i] = hit.Source
	}

	return rt, nil
}

func (s *ESClient) searchByDSL(ctx context.Context, dslStr string, pagination *storage.Pagination) (*SearchResponse, error) {
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

func (s *ESClient) searchBySQL(ctx context.Context, sqlStr string, pagination *storage.Pagination) (*SearchResponse, error) {
	dsl, _, err := elasticsql.Convert(sqlStr)
	if err != nil {
		return nil, err
	}
	return s.search(ctx, strings.NewReader(dsl), pagination)
}

func (s *ESClient) SearchByQuery(ctx context.Context, query map[string]interface{}, pagination *storage.Pagination) (*SearchResponse, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}
	return s.search(ctx, buf, pagination)
}

func (s *ESClient) search(ctx context.Context, body io.Reader, pagination *storage.Pagination) (*SearchResponse, error) {
	opts := []func(*esapi.SearchRequest){
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.indexName),
		s.client.Search.WithBody(body),
	}
	if pagination != nil {
		from := (pagination.Page - 1) * pagination.PageSize
		opts = append(
			opts,
			s.client.Search.WithSize(pagination.PageSize),
			s.client.Search.WithFrom(from),
		)
	}
	resp, err := s.client.Search(opts...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, &ESError{
			StatusCode: resp.StatusCode,
			Message:    resp.String(),
		}
	}
	sr := &SearchResponse{}
	if err := json.NewDecoder(resp.Body).Decode(sr); err != nil {
		return nil, err
	}
	return sr, nil
}
