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

	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/cch123/elasticsql"
)

func (s *ESClient) Search(ctx context.Context, queryStr string, patternType string, pageSize, page int) (*storage.SearchResult, error) {
	var res *SearchResponse
	var err error
	switch patternType {
	case storage.DSLPatternType:
		res, err = s.searchByDSL(ctx, queryStr, pageSize, page)
		if err != nil {
			return nil, err
		}
	case storage.SQLPatternType:
		res, err = s.searchBySQL(ctx, queryStr, pageSize, page)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid type %s", patternType)
	}

	rt := &storage.SearchResult{}
	rt.Total = res.Hits.Total.Value
	hits := res.Hits.Hits
	resources := make([]*storage.Resource, len(hits))
	for i := range hits {
		resources[i] = hits[i].Source
	}
	rt.Resources = resources
	return rt, nil
}

func (s *ESClient) searchByDSL(ctx context.Context, dslStr string, pageSize, page int) (*SearchResponse, error) {
	queries, err := Parse(dslStr)
	if err != nil {
		return nil, err
	}
	esQuery, err := ParseQueries(queries)
	if err != nil {
		return nil, err
	}
	res, err := s.searchByQuery(ctx, esQuery, pageSize, page)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ESClient) searchBySQL(ctx context.Context, sqlStr string, pageSize, page int) (*SearchResponse, error) {
	dsl, _, err := elasticsql.Convert(sqlStr)
	if err != nil {
		return nil, err
	}
	return s.search(ctx, strings.NewReader(dsl), pageSize, page)
}

func (s *ESClient) searchByQuery(ctx context.Context, query map[string]interface{}, pageSize, page int) (*SearchResponse, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}
	return s.search(ctx, buf, pageSize, page)
}

func (s *ESClient) search(ctx context.Context, body io.Reader, pageSize, page int) (*SearchResponse, error) {
	from := (page - 1) * pageSize
	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.indexName),
		s.client.Search.WithBody(body),
		s.client.Search.WithSize(pageSize),
		s.client.Search.WithFrom(from),
	)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, &ESError{
			StatusCode: res.StatusCode,
			Message:    res.String(),
		}
	}
	sr := &SearchResponse{}
	if err := json.NewDecoder(res.Body).Decode(sr); err != nil {
		return nil, err
	}
	return sr, nil
}
