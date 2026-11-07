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

package search

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
)

type fakeSearchStorage struct {
	storage.SearchStorage
	got *storage.Pagination
}

func (f *fakeSearchStorage) Search(_ context.Context, _, _ string, p *storage.Pagination) (*storage.SearchResult, error) {
	f.got = p
	return &storage.SearchResult{}, nil
}

func TestSearchForResourcePagination(t *testing.T) {
	cases := []struct {
		name         string
		query        string
		wantPage     int
		wantPageSize int
	}{
		{"defaults", "", 1, 10},
		{"page size 1 is honored", "&pageSize=1&page=2", 2, 1},
		{"page size 5", "&pageSize=5", 1, 5},
		{"zero page size falls back to default", "&pageSize=0", 1, 10},
		{"negative values fall back to defaults", "&pageSize=-3&page=-1", 1, 10},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fs := &fakeSearchStorage{}
			h := SearchForResource(nil, nil, fs)
			req := httptest.NewRequest(http.MethodGet, "/rest-api/v1/search?query=select%20*%20from%20resources&pattern=sql"+tc.query, nil)
			h(httptest.NewRecorder(), req)
			if fs.got == nil {
				t.Fatal("storage.Search was not called")
			}
			if fs.got.Page != tc.wantPage || fs.got.PageSize != tc.wantPageSize {
				t.Errorf("got page=%d pageSize=%d, want page=%d pageSize=%d",
					fs.got.Page, fs.got.PageSize, tc.wantPage, tc.wantPageSize)
			}
		})
	}
}
