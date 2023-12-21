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
	"github.com/KusionStack/karbour/pkg/core/manager/cluster"
	"github.com/KusionStack/karbour/pkg/core/manager/resource"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
	"k8s.io/apiserver/pkg/server"
)

// GetTopology returns an HTTP handler function that returns a topology map for
// a Kubernetes resource. It utilizes a ResourceManager to execute the logic.
//
//	@Summary		GetTopology returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Description	This endpoint returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Tags			resource
//	@Produce		json
//	@Success		200	{object}	map[string]resource.ResourceTopology	"map from string to resource.ResourceTopology"
//	@Failure		400	{string}	string									"Bad Request"
//	@Failure		401	{string}	string									"Unauthorized"
//	@Failure		404	{string}	string									"Not Found"
//	@Failure		405	{string}	string									"Method Not Allowed"
//	@Failure		429	{string}	string									"Too Many Requests"
//	@Failure		500	{string}	string									"Internal Server Error"
//	@Router			/api/v1/insight/topology [get]
func GetTopology(resourceMgr *resource.ResourceManager, clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		loc, err := resource.BuildLocatorFromQuery(r)
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
		if ok && (locType == core.Resource || locType == core.NonNamespacedResource) {
			resourceTopologyMap, err := resourceMgr.GetTopologyForResource(r.Context(), client, loc)
			if err != nil {
				render.Render(w, r, handler.FailureResponse(ctx, err))
				return
			}
			render.JSON(w, r, handler.SuccessResponse(ctx, resourceTopologyMap))
		} else if ok && locType == core.Cluster {
			clusterTopologyMap, err := clusterMgr.GetTopologyForCluster(r.Context(), client, loc.Cluster)
			if err != nil {
				render.Render(w, r, handler.FailureResponse(ctx, err))
				return
			}
			render.JSON(w, r, handler.SuccessResponse(ctx, clusterTopologyMap))
		} else if ok && locType == core.Namespace {
			namespaceTopologyMap, err := clusterMgr.GetTopologyForClusterNamespace(r.Context(), client, loc.Cluster, loc.Namespace)
			if err != nil {
				render.Render(w, r, handler.FailureResponse(ctx, err))
				return
			}
			render.JSON(w, r, handler.SuccessResponse(ctx, namespaceTopologyMap))
		} else {
			render.Render(w, r, handler.FailureResponse(ctx, fmt.Errorf("no applicable locator type found")))
		}
	}
}
