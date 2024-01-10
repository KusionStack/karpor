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

//nolint:tagliatelle
package elasticsearch

import "github.com/KusionStack/karbour/pkg/infra/search/storage"

type SearchResponse struct {
	ScrollID string `json:"_scroll_id"`
	Took     int    `json:"took"`
	TimeOut  bool   `json:"time_out"`
	Hits     *Hits  `json:"hits"`
}

type Hits struct {
	Total    *Total  `json:"total"`
	MaxScore float32 `json:"max_score"`
	Hits     []*Hit  `json:"hits"`
}

type Total struct {
	Value    int    `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type Hit struct {
	Index  string            `json:"_index"`
	ID     string            `json:"_id"`
	Score  float32           `json:"_score"`
	Source *storage.Resource `json:"_source"`
}

func (s *SearchResponse) GetResources() []*storage.Resource {
	if s == nil || s.Hits == nil {
		return nil
	}
	rt := make([]*storage.Resource, len(s.Hits.Hits))
	for i, hit := range s.Hits.Hits {
		rt[i] = hit.Source
	}
	return rt
}
