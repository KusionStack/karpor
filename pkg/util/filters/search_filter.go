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

package filters

import (
	"context"
	"net/http"
)

type searchType int

const (
	searchContextKey searchType = iota
	searchQueryKey              = "query"
)

func WithSearchQuery(parent context.Context, query string) context.Context {
	return context.WithValue(parent, searchContextKey, query)
}

func SearchQueryFrom(ctx context.Context) (string, bool) {
	query, ok := ctx.Value(searchContextKey).(string)
	if !ok {
		return "", false
	}
	return query, true
}

func SearchFilter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		searchQuery, ok := query[searchQueryKey]
		ctx := req.Context()
		if ok {
			ctx = WithSearchQuery(ctx, searchQuery[0])
			req = req.WithContext(ctx)
			query.Del(searchQueryKey)
		}
		handler.ServeHTTP(w, req)
	})
}
