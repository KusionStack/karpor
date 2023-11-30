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

	"github.com/KusionStack/karbour/pkg/search/storage"
)

type CtxTyp string

type Resource struct {
	Name       string
	Namespace  string
	APIVersion string
	Kind       string
	Cluster    string
}

const (
	searchQueryKey             = "query"
	patternTypeKey             = "patternType"
	resourceNameQueryKey       = "name"
	resourceNamespaceQueryKey  = "namespace"
	resourceClusterQueryKey    = "cluster"
	resourceAPIVersionQueryKey = "apiVersion"
	resourceKindQueryKey       = "kind"
)

func SearchQueryFrom(ctx context.Context) (string, bool) {
	query, ok := ctx.Value(CtxTyp(searchQueryKey)).(string)
	if !ok {
		return "", false
	}
	return query, true
}

func PatternTypeFrom(ctx context.Context) (string, bool) {
	patternType, ok := ctx.Value(CtxTyp(patternTypeKey)).(string)
	if !ok {
		return "", false
	}
	return patternType, true
}

func ResourceDetailFrom(ctx context.Context) (Resource, bool) {
	res := Resource{}
	resourceName, ok := ctx.Value(CtxTyp(resourceNameQueryKey)).(string)
	if !ok {
		return res, false
	}
	namespace, ok := ctx.Value(CtxTyp(resourceNamespaceQueryKey)).(string)
	if !ok {
		return res, false
	}
	cluster, ok := ctx.Value(CtxTyp(resourceClusterQueryKey)).(string)
	if !ok {
		return res, false
	}
	apiVersion, ok := ctx.Value(CtxTyp(resourceAPIVersionQueryKey)).(string)
	if !ok {
		return res, false
	}
	kind, ok := ctx.Value(CtxTyp(resourceKindQueryKey)).(string)
	if !ok {
		return res, false
	}
	return Resource{
		Name:       resourceName,
		Namespace:  namespace,
		Cluster:    cluster,
		APIVersion: apiVersion,
		Kind:       kind,
	}, true
}

func SearchFilter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = FromQueryToContext(req, searchQueryKey, "")
		req = FromQueryToContext(req, patternTypeKey, storage.DSLPatternType)
		req = FromQueryToContext(req, resourceNameQueryKey, "")
		req = FromQueryToContext(req, resourceNamespaceQueryKey, "")
		req = FromQueryToContext(req, resourceClusterQueryKey, "")
		req = FromQueryToContext(req, resourceAPIVersionQueryKey, "")
		req = FromQueryToContext(req, resourceKindQueryKey, "")
		handler.ServeHTTP(w, req)
	})
}

func FromQueryToContext(req *http.Request, key string, defaultVal string) *http.Request {
	query := req.URL.Query()
	queryVal, ok := query[key]
	var val string
	if !ok {
		if defaultVal == "" {
			return req
		}
		val = defaultVal
	} else {
		query.Del(key)
		val = queryVal[0]
	}
	return req.WithContext(context.WithValue(req.Context(), CtxTyp(key), val))
}
