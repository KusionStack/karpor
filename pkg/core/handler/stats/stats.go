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

package stats

import (
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/insight"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
)

// GetStatistics returns an HTTP handler function that returns a statistics
// info. It utilizes an InsightManager to execute the logic.
//
// @Summary      Get returns a global statistics info.
// @Description  This endpoint returns a global statistics info.
// @Tags         insight
// @Produce      json
// @Success      200         {object}  insight.Statistics "Global statistics info"
// @Failure      400         {string}  string                   "Bad Request"
// @Failure      401         {string}  string                   "Unauthorized"
// @Failure      404         {string}  string                   "Not Found"
// @Failure      405         {string}  string                   "Method Not Allowed"
// @Failure      429         {string}  string                   "Too Many Requests"
// @Failure      500         {string}  string                   "Internal Server Error"
// @Router       /rest-api/v1/insight/stats [get]
func GetStatistics(insightMgr *insight.InsightManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		logger.Info("Getting statistics ...")
		statistics, err := insightMgr.Statistics(ctx)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		render.Render(w, r, handler.SuccessResponse(ctx, statistics))
	}
}
