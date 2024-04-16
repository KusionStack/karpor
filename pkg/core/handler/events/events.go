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

package events

import (
	"fmt"
	"net/http"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/insight"
	"github.com/KusionStack/karbour/pkg/infra/multicluster"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
	"k8s.io/apiserver/pkg/server"
)

// GetEvents returns an HTTP handler function that returns events for a
// Kubernetes resource. It utilizes an InsightManager to execute the logic.
//
// @Summary      GetEvents returns events for a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
// @Description  This endpoint returns events for a Kubernetes resource YAML by name, namespace, cluster, apiVersion and kind.
// @Tags         insight
// @Produce      json
// @Param        cluster     query     string                     false  "The specified cluster name, such as 'example-cluster'"
// @Param        apiVersion  query     string                     false  "The specified apiVersion, such as 'apps/v1'. Should be percent-encoded"
// @Param        kind        query     string                     false  "The specified kind, such as 'Deployment'"
// @Param        namespace   query     string                     false  "The specified namespace, such as 'default'"
// @Param        name        query     string                     false  "The specified resource name, such as 'foo'"
// @Success      200         {array}   unstructured.Unstructured  "List of events"
// @Failure      400         {string}  string                     "Bad Request"
// @Failure      401         {string}  string                     "Unauthorized"
// @Failure      404         {string}  string                     "Not Found"
// @Failure      405         {string}  string                     "Method Not Allowed"
// @Failure      429         {string}  string                     "Too Many Requests"
// @Failure      500         {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/insight/events [get]
func GetEvents(insightMgr *insight.InsightManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		loc, err := core.NewResourceGroupFromQuery(r)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		logger.Info("Getting events for resourceGroup...", "resourceGroup", loc)

		client, err := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, loc.Cluster)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		locType, ok := loc.GetType()
		if !ok {
			render.Render(w, r, handler.FailureResponse(ctx, fmt.Errorf("unable to determine resource group type")))
			return
		}

		switch locType {
		case core.Resource, core.NonNamespacedResource:
			resourceEvents, err := insightMgr.GetResourceEvents(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, resourceEvents)
		case core.Cluster:
			clusterEvents, err := insightMgr.GetClusterEvents(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, clusterEvents)
		case core.ClusterGVKNamespace:
			clusterGVKEvents, err := insightMgr.GetNamespaceGVKEvents(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, clusterGVKEvents)
		case core.Namespace:
			namespaceEvents, err := insightMgr.GetNamespaceEvents(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, namespaceEvents)
		case core.GVK:
			gvkEvents, err := insightMgr.GetGVKEvents(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, gvkEvents)
		default:
			render.Render(w, r, handler.FailureResponse(ctx, fmt.Errorf("no applicable resource group type found")))
		}
	}
}
