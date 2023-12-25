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

package summary

import (
	"fmt"
	"net/http"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/insight"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
	"k8s.io/apiserver/pkg/server"
)

// GetSummary returns an HTTP handler function that returns a Kubernetes
// resource summary. It utilizes an InsightManager to execute the logic.
//
//	@Summary		Get returns a Kubernetes resource summary by name, namespace, cluster, apiVersion and kind.
//	@Description	This endpoint returns a Kubernetes resource summary by name, namespace, cluster, apiVersion and kind.
//	@Tags			insight
//	@Produce		json
//	@Param			cluster		query		string					false	"The specified cluster name, such as 'example-cluster'"
//	@Param			apiVersion	query		string					false	"The specified apiVersion, such as 'apps/v1'. Should be percent-encoded"
//	@Param			kind		query		string					false	"The specified kind, such as 'Deployment'"
//	@Param			namespace	query		string					false	"The specified namespace, such as 'default'"
//	@Param			name		query		string					false	"The specified resource name, such as 'foo'"
//	@Success		200			{object}	insight.ResourceSummary	"Resource Summary"
//	@Failure		400			{string}	string					"Bad Request"
//	@Failure		401			{string}	string					"Unauthorized"
//	@Failure		404			{string}	string					"Not Found"
//	@Failure		405			{string}	string					"Method Not Allowed"
//	@Failure		429			{string}	string					"Too Many Requests"
//	@Failure		500			{string}	string					"Internal Server Error"
//	@Router			/api/v1/insight/summary [get]
func GetSummary(insightMgr *insight.InsightManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		loc, err := core.NewLocatorFromQuery(r)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		logger.Info("Getting summary for locator...", "locator", loc)

		client, err := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, loc.Cluster)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		locType, ok := loc.GetType()
		if !ok {
			render.Render(w, r, handler.FailureResponse(ctx, fmt.Errorf("unable to determine locator type")))
			return
		}

		switch locType {
		case core.Resource, core.NonNamespacedResource:
			resourceSummary, err := insightMgr.GetResourceSummary(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, resourceSummary)
		case core.Cluster:
			clusterDetail, err := insightMgr.GetDetailsForCluster(r.Context(), client, loc.Cluster)
			handler.HandleResult(w, r, ctx, err, clusterDetail)
		case core.Namespace:
			namespaceSummary, err := insightMgr.GetNamespaceSummary(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, namespaceSummary)
		case core.GVK:
			gvkSummary, err := insightMgr.GetGVKSummary(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, gvkSummary)
		default:
			render.Render(w, r, handler.FailureResponse(ctx, fmt.Errorf("no applicable locator type found")))
		}
	}
}
