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

package detail

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/core/handler"
	"github.com/KusionStack/karpor/pkg/core/manager/cluster"
	"github.com/KusionStack/karpor/pkg/core/manager/insight"
	"github.com/KusionStack/karpor/pkg/infra/multicluster"
	"github.com/KusionStack/karpor/pkg/util/ctxutil"
	"k8s.io/apiserver/pkg/server"
)

// GetDetail returns an HTTP handler function that returns a Kubernetes
// resource. It utilizes an InsightManager to execute the logic.
//
// @Summary      GetDetail returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
// @Description  This endpoint returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
// @Tags         insight
// @Produce      json
// @Param        format      query     string                     false  "The format of the response. Either in json or yaml. Default to json"
// @Param        cluster     query     string                     false  "The specified cluster name, such as 'example-cluster'"
// @Param        apiVersion  query     string                     false  "The specified apiVersion, such as 'apps/v1'. Should be percent-encoded"
// @Param        kind        query     string                     false  "The specified kind, such as 'Deployment'"
// @Param        namespace   query     string                     false  "The specified namespace, such as 'default'"
// @Param        name        query     string                     false  "The specified resource name, such as 'foo'"
// @Success      200         {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400         {string}  string                     "Bad Request"
// @Failure      401         {string}  string                     "Unauthorized"
// @Failure      404         {string}  string                     "Not Found"
// @Failure      405         {string}  string                     "Method Not Allowed"
// @Failure      429         {string}  string                     "Too Many Requests"
// @Failure      500         {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/insight/detail [get]
func GetDetail(clusterMgr *cluster.ClusterManager, insightMgr *insight.InsightManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		logger.Info("Getting resources...")
		outputFormat := r.URL.Query().Get("format")

		resourceGroup, err := entity.NewResourceGroupFromQuery(r)
		if err != nil {
			handler.FailureRender(ctx, w, r, err)
			return
		}
		logger.Info("Getting resource detail for resourceGroup...", "resourceGroup", resourceGroup)

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
			if strings.ToLower(outputFormat) == "yaml" {
				resourceYAML, err := insightMgr.GetYAMLForResource(r.Context(), client, &resourceGroup)
				handler.HandleResult(w, r, ctx, err, string(resourceYAML))
			} else {
				resourceUnstructured, err := insightMgr.GetResource(r.Context(), client, &resourceGroup)
				handler.HandleResult(w, r, ctx, err, resourceUnstructured)
			}
		case entity.Namespace:
			if strings.ToLower(outputFormat) == "yaml" {
				namespaceYAML, err := clusterMgr.GetNamespaceYAML(r.Context(), client, resourceGroup.Namespace)
				handler.HandleResult(w, r, ctx, err, string(namespaceYAML))
			} else {
				namespace, err := clusterMgr.GetNamespace(r.Context(), client, resourceGroup.Namespace)
				handler.HandleResult(w, r, ctx, err, namespace)
			}
		default:
			handler.FailureRender(ctx, w, r, fmt.Errorf("unsupported resource group type: %v", resourceGroupType))
		}
	}
}
