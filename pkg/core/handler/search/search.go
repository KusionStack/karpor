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

package search

import (
	"net/http"
	"strconv"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/search"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// SearchForResource returns an HTTP handler function that returns an
// array of Kubernetes runtime Object matched using the query from
// context. It utilizes a SearchManager to execute the logic.
//
// @Summary      SearchForResource returns an array of Kubernetes runtime Object matched using the query from context.
// @Description  This endpoint returns an array of Kubernetes runtime Object matched using the query from context.
// @Tags         search
// @Produce      json
// @Param        query     query     string          true   "The query to use for search. Required"
// @Param        pattern   query     string          true   "The search pattern. Can be either sql or dsl. Required"
// @Param        pageSize  query     string          false  "The size of the page. Default to 10"
// @Param        page      query     string          false  "The current page to fetch. Default to 1"
// @Success      200       {array}   runtime.Object  "Array of runtime.Object"
// @Failure      400       {string}  string          "Bad Request"
// @Failure      401       {string}  string          "Unauthorized"
// @Failure      404       {string}  string          "Not Found"
// @Failure      405       {string}  string          "Method Not Allowed"
// @Failure      429       {string}  string          "Too Many Requests"
// @Failure      500       {string}  string          "Internal Server Error"
// @Router       /rest-api/v1/search [get]
func SearchForResource(searchMgr *search.SearchManager, searchStorage storage.SearchStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		// Extract URL query parameters with default value
		searchQuery := r.URL.Query().Get("query")
		searchPattern := r.URL.Query().Get("pattern")
		searchPageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		searchPage, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if searchPageSize <= 1 {
			searchPageSize = 10
		}
		if searchPage <= 1 {
			searchPage = 1
		}

		logger.Info("Searching for resources...", "page", searchPage, "pageSize", searchPageSize)

		res, err := searchStorage.Search(ctx, searchQuery, searchPattern, searchPageSize, searchPage)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		rt := &search.UniResourceList{}
		for _, r := range res.Resources {
			unObj := &unstructured.Unstructured{}
			unObj.SetUnstructuredContent(r.Object)
			rt.Items = append(rt.Items, search.UniResource{
				Cluster: r.Cluster,
				Object:  unObj,
			})
		}
		rt.Total = res.Total
		rt.CurrentPage = searchPage
		rt.PageSize = searchPageSize
		render.JSON(w, r, handler.SuccessResponse(ctx, rt))
	}
}
