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

package topology

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

// GetTopology returns an HTTP handler function that returns a topology map for
// a Kubernetes resource. It utilizes an InsightManager to execute the logic.
//
//	@Summary		GetTopology returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Description	This endpoint returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Tags			insight
//	@Produce		json
//	@Param			cluster		query		string								false	"The specified cluster name, such as 'example-cluster'"
//	@Param			apiVersion	query		string								false	"The specified apiVersion, such as 'apps/v1'. Should be percent-encoded"
//	@Param			kind		query		string								false	"The specified kind, such as 'Deployment'"
//	@Param			namespace	query		string								false	"The specified namespace, such as 'default'"
//	@Param			name		query		string								false	"The specified resource name, such as 'foo'"
//	@Success		200			{object}	map[string]insight.ResourceTopology	"map from string to resource.ResourceTopology"
//	@Failure		400			{string}	string								"Bad Request"
//	@Failure		401			{string}	string								"Unauthorized"
//	@Failure		404			{string}	string								"Not Found"
//	@Failure		405			{string}	string								"Method Not Allowed"
//	@Failure		429			{string}	string								"Too Many Requests"
//	@Failure		500			{string}	string								"Internal Server Error"
//	@Router			/api/v1/insight/topology [get]
func GetTopology(insightMgr *insight.InsightManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		loc, err := core.NewLocatorFromQuery(r)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		logger.Info("Getting topology for locator...", "locator", loc)

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
			resourceTopologyMap, err := insightMgr.GetTopologyForResource(r.Context(), client, &loc)
			handler.HandleResult(w, r, ctx, err, resourceTopologyMap)
		case core.Cluster:
			clusterTopologyMap, err := insightMgr.GetTopologyForCluster(r.Context(), client, loc.Cluster)
			handler.HandleResult(w, r, ctx, err, clusterTopologyMap)
		case core.Namespace:
			namespaceTopologyMap, err := insightMgr.GetTopologyForClusterNamespace(r.Context(), client, loc.Cluster, loc.Namespace)
			handler.HandleResult(w, r, ctx, err, namespaceTopologyMap)
		default:
			render.Render(w, r, handler.FailureResponse(ctx, fmt.Errorf("no applicable locator type found")))
		}
	}
}
