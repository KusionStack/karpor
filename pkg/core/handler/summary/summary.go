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

package summary

import (
	"fmt"
	"net/http"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/core/handler"
	"github.com/KusionStack/karpor/pkg/core/manager/insight"
	"github.com/KusionStack/karpor/pkg/infra/multicluster"
	"github.com/KusionStack/karpor/pkg/util/ctxutil"
	"k8s.io/apiserver/pkg/server"
)

// GetSummary returns an HTTP handler function that returns a Kubernetes
// resource summary. It utilizes an InsightManager to execute the logic.
//
// @Summary      Get returns a Kubernetes resource summary by name, namespace, cluster, apiVersion and kind.
// @Description  This endpoint returns a Kubernetes resource summary by name, namespace, cluster, apiVersion and kind.
// @Tags         insight
// @Produce      json
// @Param        cluster     query     string                   false  "The specified cluster name, such as 'example-cluster'"
// @Param        apiVersion  query     string                   false  "The specified apiVersion, such as 'apps/v1'. Should be percent-encoded"
// @Param        kind        query     string                   false  "The specified kind, such as 'Deployment'"
// @Param        namespace   query     string                   false  "The specified namespace, such as 'default'"
// @Param        name        query     string                   false  "The specified resource name, such as 'foo'"
// @Success      200         {object}  insight.ResourceSummary  "Resource Summary"
// @Failure      400         {string}  string                   "Bad Request"
// @Failure      401         {string}  string                   "Unauthorized"
// @Failure      404         {string}  string                   "Not Found"
// @Failure      405         {string}  string                   "Method Not Allowed"
// @Failure      429         {string}  string                   "Too Many Requests"
// @Failure      500         {string}  string                   "Internal Server Error"
// @Router       /rest-api/v1/insight/summary [get]
func GetSummary(insightMgr *insight.InsightManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		resourceGroup, err := entity.NewResourceGroupFromQuery(r)
		if err != nil {
			handler.FailureRender(ctx, w, r, err)
			return
		}
		logger.Info("Getting summary for resourceGroup...", "resourceGroup", resourceGroup)

		client, err := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, resourceGroup.Cluster)
		if err != nil {
			handler.FailureRender(ctx, w, r, err)
			return
		}

		resourceGroupType, ok := resourceGroup.GetType()
		if !ok {
			handler.FailureRender(ctx, w, r, fmt.Errorf("unable to determine resource group type"))
			return
		}

		switch resourceGroupType {
		case entity.Resource, entity.NonNamespacedResource:
			resourceSummary, err := insightMgr.GetResourceSummary(r.Context(), client, &resourceGroup)
			handler.HandleResult(w, r, ctx, err, resourceSummary)
		case entity.Cluster:
			clusterDetail, err := insightMgr.GetDetailsForCluster(r.Context(), client, resourceGroup.Cluster)
			handler.HandleResult(w, r, ctx, err, clusterDetail)
		case entity.Namespace:
			namespaceSummary, err := insightMgr.GetNamespaceSummary(r.Context(), client, &resourceGroup)
			handler.HandleResult(w, r, ctx, err, namespaceSummary)
		case entity.GVK:
			gvkSummary, err := insightMgr.GetGVKSummary(r.Context(), client, &resourceGroup)
			handler.HandleResult(w, r, ctx, err, gvkSummary)
		case entity.Custom:
			rgSummary, err := insightMgr.GetResourceGroupSummary(r.Context(), client, &resourceGroup)
			handler.HandleResult(w, r, ctx, err, rgSummary)
		default:
			handler.FailureRender(ctx, w, r, fmt.Errorf("unsupported resource group type"))
		}
	}
}
