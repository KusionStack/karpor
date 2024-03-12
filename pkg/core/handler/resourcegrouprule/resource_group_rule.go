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

package resourcegrouprule

import (
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/resourcegroup"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"k8s.io/apiserver/pkg/server"
)

// Get returns an HTTP handler function that reads a resourcegrouprule
// detail. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      Get returns a ResourceGroupRule by name.
// @Description  This endpoint returns a ResourceGroupRule by name.
// @Tags         resourcegrouprule
// @Produce      json
// @Param        format  query     string                     false  "The format of the response. Either in json or yaml"
// @Success      200     {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400     {string}  string                     "Bad Request"
// @Failure      401     {string}  string                     "Unauthorized"
// @Failure      404     {string}  string                     "Not Found"
// @Failure      405     {string}  string                     "Method Not Allowed"
// @Failure      429     {string}  string                     "Too Many Requests"
// @Failure      500     {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rule/{resourceGroupRule} [get]
func Get(resourceGroupMgr *resourcegroup.ResourceGroupManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		resourceGroupRule := chi.URLParam(r, "resourceGroupRule")
		logger.Info("Getting resourceGroupRule...", "resourceGroupRule", resourceGroupRule)

		// TODO: need to implement
	}
}

// Create returns an HTTP handler function that creates a ResourceGroupRule
// resource. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      Create creates a ResourceGroupRule.
// @Description  This endpoint creates a new ResourceGroupRule using the payload.
// @Tags         resourcegrouprule
// @Accept       plain
// @Accept       json
// @Produce      json
// @Param        request  body      ResourceGroupRulePayload             true  "resourceGroupRule to create (either plain text or JSON format)"
// @Success      200      {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400      {string}  string                     "Bad Request"
// @Failure      401      {string}  string                     "Unauthorized"
// @Failure      404      {string}  string                     "Not Found"
// @Failure      405      {string}  string                     "Method Not Allowed"
// @Failure      429      {string}  string                     "Too Many Requests"
// @Failure      500      {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rule/{resourceGroupRule} [post]
func Create(resourceGroupMgr *resourcegroup.ResourceGroupManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		resourceGroupRule := chi.URLParam(r, "resourceGroupRule")
		logger.Info("Creating resourceGroupRule...", "resourceGroupRule", resourceGroupRule)

		// Decode the request body into the payload.
		payload := &ResourceGroupRulePayload{}
		if err := payload.Decode(r); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// TODO: need to implement
	}
}

// Update returns an HTTP handler function that updates a ResourceGroupRule
// resource. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      Update updates the ResourceGroupRule metadata by name.
// @Description  This endpoint updates the display name and description of an existing ResourceGroupRule.
// @Tags         resourcegrouprule
// @Accept       plain
// @Accept       json
// @Produce      json
// @Param        request  body      ResourceGroupRulePayload             true  "resourceGroupRule to update (either plain text or JSON format)"
// @Success      200      {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400      {string}  string                     "Bad Request"
// @Failure      401      {string}  string                     "Unauthorized"
// @Failure      404      {string}  string                     "Not Found"
// @Failure      405      {string}  string                     "Method Not Allowed"
// @Failure      429      {string}  string                     "Too Many Requests"
// @Failure      500      {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rule/{resourceGroupRule}  [put]
func Update(resourceGroupMgr *resourcegroup.ResourceGroupManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		resourceGroupRule := chi.URLParam(r, "resourceGroupRule")
		logger.Info("Updating resourceGroupRule metadata...", "resourceGroupRule", resourceGroupRule)

		// Decode the request body into the payload.
		payload := &ResourceGroupRulePayload{}
		if err := payload.Decode(r); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// TODO: need to implement
	}
}

// List returns an HTTP handler function that lists all ResourceGroupRule
// resources. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      List lists all ResourceGroupRules.
// @Description  This endpoint lists all ResourceGroupRules.
// @Tags         resourcegrouprule
// @Produce      json
// @Param        summary     query     bool                       false  "Whether to display summary or not. Default to false"
// @Param        orderBy     query     string                     false  "The order to list the resourceGroupRule. Default to order by name"
// @Param        descending  query     bool                       false  "Whether to sort the list in descending order. Default to false"
// @Success      200         {array}   unstructured.Unstructured  "List of resourceGroupRule objects"
// @Failure      400         {string}  string                     "Bad Request"
// @Failure      401         {string}  string                     "Unauthorized"
// @Failure      404         {string}  string                     "Not Found"
// @Failure      405         {string}  string                     "Method Not Allowed"
// @Failure      429         {string}  string                     "Too Many Requests"
// @Failure      500         {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rules [get]
func List(resourceGroupMgr *resourcegroup.ResourceGroupManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		logger.Info("Listing resourceGroupRules...")

		// TODO: need to implement
	}
}

// Delete returns an HTTP handler function that deletes a ResourceGroupRule
// resource. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      Delete removes a ResourceGroupRule by name.
// @Description  This endpoint deletes the ResourceGroupRule by name.
// @Tags         resourcegrouprule
// @Produce      json
// @Success      200  {string}  string  "Operation status"
// @Failure      400  {string}  string  "Bad Request"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      404  {string}  string  "Not Found"
// @Failure      405  {string}  string  "Method Not Allowed"
// @Failure      429  {string}  string  "Too Many Requests"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rule/{resourceGroupRule} [delete]
func Delete(resourceGroupMgr *resourcegroup.ResourceGroupManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		resourceGroupRule := chi.URLParam(r, "resourceGroupRule")
		logger.Info("Deleting resourceGroupRule...", "resourceGroupRule", resourceGroupRule)

		// TODO: need to implement
	}
}
