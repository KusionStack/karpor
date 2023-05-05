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
	"context"

	"github.com/KusionStack/karbour/pkg/search/storage"
)

func (s *ESClient) Search(ctx context.Context, queries []storage.Query) (*storage.SearchResult, error) {
	esQuery, err := ParseQueries(queries)
	if err != nil {
		return nil, err
	}

	res, err := s.searchByQuery(ctx, esQuery)
	if err != nil {
		return nil, err
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

func (s *ESClient) SearchByString(ctx context.Context, queryString string) (*storage.SearchResult, error) {
	queries, err := Parse(queryString)
	if err != nil {
		return nil, err
	}
	return s.Search(ctx, queries)
}
