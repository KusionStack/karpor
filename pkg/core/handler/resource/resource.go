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

package resource

import (
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/resource"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
	"k8s.io/apiserver/pkg/server"
)

// Get returns an HTTP handler function that returns a Kubernetes
// resource. It utilizes a ResourceManager to execute the logic.
//
//	@Summary		Get returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Description	This endpoint returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Tags			resource
//	@Produce		json
//	@Success		200	{object}	unstructured.Unstructured	"Unstructured object"
//	@Failure		400	{string}	string						"Bad Request"
//	@Failure		401	{string}	string						"Unauthorized"
//	@Failure		404	{string}	string						"Not Found"
//	@Failure		405	{string}	string						"Method Not Allowed"
//	@Failure		429	{string}	string						"Too Many Requests"
//	@Failure		500	{string}	string						"Internal Server Error"
//	@Router			/api/v1/resource/cluster/{clusterName}/{apiVersion}/namespace/{namespaceName}/{kind}/name/{resourceName}/ [get]
func Get(resourceMgr *resource.ResourceManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		logger.Info("Getting resources...")

		loc, err := resource.BuildLocatorFromQuery(r)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		client, err := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, loc.Cluster)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		resourceUnstructured, err := resourceMgr.GetResource(r.Context(), client, loc)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		render.JSON(w, r, handler.SuccessResponse(ctx, resourceUnstructured))
	}
}

// GetYAML returns an HTTP handler function that returns a Kubernetes
// resource YAML. It utilizes a ResourceManager to execute the logic.
//
//	@Summary		GetYAML returns a Kubernetes resource YAML by name, namespace, cluster, apiVersion and kind.
//	@Description	This endpoint returns a Kubernetes resource YAML by name, namespace, cluster, apiVersion and kind.
//	@Tags			resource
//	@Produce		json
//	@Success		200	{array}		byte	"Byte array"
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		404	{string}	string	"Not Found"
//	@Failure		405	{string}	string	"Method Not Allowed"
//	@Failure		429	{string}	string	"Too Many Requests"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/api/v1/resource/cluster/{clusterName}/{apiVersion}/namespace/{namespaceName}/{kind}/name/{resourceName}/yaml [get]
func GetYAML(resourceMgr *resource.ResourceManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		loc, err := resource.BuildLocatorFromQuery(r)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		logger.Info("Getting YAML for locator...", "locator", loc)

		client, err := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, loc.Cluster)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		result, err := resourceMgr.GetYAMLForResource(r.Context(), client, loc)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		render.JSON(w, r, handler.SuccessResponse(ctx, string(result)))
	}
}
