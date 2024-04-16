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

package scanner

import (
	"net/http"
	"strconv"

	"github.com/KusionStack/karbour/pkg/core/entity"
	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/insight"
	_ "github.com/KusionStack/karbour/pkg/infra/scanner"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
)

// Audit handles the auditing process based on the specified resource group.
//
// @Summary      Audit based on resource group.
// @Description  This endpoint audits based on the specified resource group.
// @Tags         insight
// @Produce      json
// @Param        cluster     query     string     false  "The specified cluster name, such as 'example-cluster'"
// @Param        apiVersion  query     string     false  "The specified apiVersion, such as 'apps/v1'"
// @Param        kind        query     string     false  "The specified kind, such as 'Deployment'"
// @Param        namespace   query     string     false  "The specified namespace, such as 'default'"
// @Param        name        query     string     false  "The specified resource name, such as 'foo'"
// @Param        forceNew    query     bool       false  "Switch for forced scanning, default is 'false'"
// @Success      200         {object}  AuditData  "Audit results"
// @Failure      400         {string}  string     "Bad Request"
// @Failure      401         {string}  string     "Unauthorized"
// @Failure      429         {string}  string     "Too Many Requests"
// @Failure      404         {string}  string     "Not Found"
// @Failure      500         {string}  string     "Internal Server Error"
// @Router       /rest-api/v1/insight/audit [get]
func Audit(insight *insight.InsightManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)

		// Begin the auditing process, logging the start.
		log.Info("Starting audit with specified resourceGroup in handler ...")

		// Decode the query parameters into the resourceGroup.
		resourceGroup, err := entity.NewResourceGroupFromQuery(r)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		forceNew, _ := strconv.ParseBool(r.URL.Query().Get("forceNew"))

		// Log successful decoding of the request body.
		log.Info("Successfully decoded the query parameters to resourceGroup", "resourceGroup", resourceGroup)

		// Perform the audit using the manager and the provided manifest.
		scanResult, err := insight.Audit(ctx, resourceGroup, forceNew)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		data := convertScanResultToAuditData(scanResult)

		render.JSON(w, r, handler.SuccessResponse(ctx, data))
	}
}

// Score returns an HTTP handler function that calculates a score for the
// audited manifest. It utilizes an AuditManager to compute the score based
// on detected issues.
//
// @Summary      ScoreHandler calculates a score for the audited manifest.
// @Description  This endpoint calculates a score for the provided manifest based on the number and severity of issues detected during the audit.
// @Tags         insight
// @Produce      json
// @Param        cluster     query     string             false  "The specified cluster name, such as 'example-cluster'"
// @Param        apiVersion  query     string             false  "The specified apiVersion, such as 'apps/v1'"
// @Param        kind        query     string             false  "The specified kind, such as 'Deployment'"
// @Param        namespace   query     string             false  "The specified namespace, such as 'default'"
// @Param        name        query     string             false  "The specified resource name, such as 'foo'"
// @Param        forceNew    query     bool               false  "Switch for forced compute score, default is 'false'"
// @Success      200         {object}  insight.ScoreData  "Score calculation result"
// @Failure      400         {string}  string             "Bad Request"
// @Failure      401         {string}  string             "Unauthorized"
// @Failure      429         {string}  string             "Too Many Requests"
// @Failure      404         {string}  string             "Not Found"
// @Failure      500         {string}  string             "Internal Server Error"
// @Router       /rest-api/v1/insight/score [get]
func Score(insightMgr *insight.InsightManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)

		// Begin the auditing process, logging the start.
		log.Info("Starting calculate score with specified resourceGroup in handler...")

		// Decode the query parameters into the resourceGroup.
		resourceGroup, err := entity.NewResourceGroupFromQuery(r)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		forceNew, _ := strconv.ParseBool(r.URL.Query().Get("forceNew"))

		// Log successful decoding of the request body.
		log.Info("Successfully decoded the query parameters to resourceGroup", "resourceGroup", resourceGroup)

		// Calculate score using the audit issues.
		data, err := insightMgr.Score(ctx, resourceGroup, forceNew)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		render.JSON(w, r, handler.SuccessResponse(ctx, data))
	}
}
