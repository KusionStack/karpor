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
	"github.com/pkg/errors"
)

// Get returns an HTTP handler function that reads a resourcegrouprule
// detail. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      Get returns a ResourceGroupRule by name.
// @Description  This endpoint returns a ResourceGroupRule by name.
// @Tags         resourcegrouprule
// @Produce      json
// @Param        resourceGroupRuleName  path      string                     true  "The name of the resource group rule"
// @Success      200                    {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400                    {string}  string                     "Bad Request"
// @Failure      401                    {string}  string                     "Unauthorized"
// @Failure      404                    {string}  string                     "Not Found"
// @Failure      405                    {string}  string                     "Method Not Allowed"
// @Failure      429                    {string}  string                     "Too Many Requests"
// @Failure      500                    {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rule/{resourceGroupRuleName} [get]
func Get(resourceGroupMgr *resourcegroup.ResourceGroupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		name := chi.URLParam(r, "resourceGroupRuleName")
		if len(name) == 0 {
			render.Render(w, r, handler.FailureResponse(ctx, errors.New("resource group rule name cannot be empty")))
			return
		}

		logger.Info("Getting resourceGroupRule...", "resourceGroupRule", name)

		// Use the ResourceGroupManager to get the resource group rule.
		data, err := resourceGroupMgr.GetResourceGroupRule(ctx, name)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Render the response in the requested format.
		render.JSON(w, r, handler.SuccessResponse(ctx, data))
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
// @Param        request  body      ResourceGroupRulePayload   true  "resourceGroupRule to create (either plain text or JSON format)"
// @Success      200      {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400      {string}  string                     "Bad Request"
// @Failure      401      {string}  string                     "Unauthorized"
// @Failure      404      {string}  string                     "Not Found"
// @Failure      405      {string}  string                     "Method Not Allowed"
// @Failure      429      {string}  string                     "Too Many Requests"
// @Failure      500      {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rule [post]
func Create(resourceGroupMgr *resourcegroup.ResourceGroupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		// Decode the request body into the payload.
		var payload ResourceGroupRulePayload
		if err := payload.Decode(r); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		if payload.Name == "" {
			render.Render(w, r, handler.FailureResponse(ctx, errors.New("resource group rule name cannot be empty")))
			return
		}

		logger.Info("Creating resourceGroupRule...", "resourceGroupRule", payload.Name)

		// Use the ResourceGroupManager to create the resource group rule.
		rgr := payload.ToEntity()
		if err := resourceGroupMgr.CreateResourceGroupRule(ctx, rgr); err != nil {
			if !errors.Is(err, resourcegroup.ErrResourceGroupRuleAlreadyExists) {
				render.Render(w, r, handler.FailureResponse(ctx, err))
				return
			}
		}

		// Render the created resource group rule.
		render.JSON(w, r, handler.SuccessResponse(ctx, payload))
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
// @Param        request  body      ResourceGroupRulePayload   true  "resourceGroupRule to update (either plain text or JSON format)"
// @Success      200      {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400      {string}  string                     "Bad Request"
// @Failure      401      {string}  string                     "Unauthorized"
// @Failure      404      {string}  string                     "Not Found"
// @Failure      405      {string}  string                     "Method Not Allowed"
// @Failure      429      {string}  string                     "Too Many Requests"
// @Failure      500      {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rule [put]
func Update(resourceGroupMgr *resourcegroup.ResourceGroupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		// Decode the request body into the payload.
		var payload ResourceGroupRulePayload
		if err := payload.Decode(r); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}
		if payload.Name == "" {
			render.Render(w, r, handler.FailureResponse(ctx, errors.New("resource group rule name cannot be empty")))
			return
		}

		logger.Info("Updating resourceGroupRule metadata...", "resourceGroupRule", payload.Name)

		// Use the ResourceGroupManager to update the resource group rule.
		rgr := payload.ToEntity()
		if err := resourceGroupMgr.UpdateResourceGroupRule(ctx, rgr.Name, rgr); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Render the updated resource group rule.
		render.JSON(w, r, handler.SuccessResponse(ctx, payload))
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
func List(resourceGroupMgr *resourcegroup.ResourceGroupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		logger.Info("Listing resourceGroupRules...")

		// Use the ResourceGroupManager to list resource group rules.
		rules, err := resourceGroupMgr.ListResourceGroupRules(ctx)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Render the list of resource group rules.
		render.JSON(w, r, handler.SuccessResponse(ctx, rules))
	}
}

// Delete returns an HTTP handler function that deletes a ResourceGroupRule
// resource. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      Delete removes a ResourceGroupRule by name.
// @Description  This endpoint deletes the ResourceGroupRule by name.
// @Tags         resourcegrouprule
// @Produce      json
// @Param        resourceGroupRuleName  path      string  true  "The name of the resource group rule"
// @Success      200                    {string}  string  "Operation status"
// @Failure      400                    {string}  string  "Bad Request"
// @Failure      401                    {string}  string  "Unauthorized"
// @Failure      404                    {string}  string  "Not Found"
// @Failure      405                    {string}  string  "Method Not Allowed"
// @Failure      429                    {string}  string  "Too Many Requests"
// @Failure      500                    {string}  string  "Internal Server Error"
// @Router       /rest-api/v1/resource-group-rule/{resourceGroupRuleName} [delete]
func Delete(resourceGroupMgr *resourcegroup.ResourceGroupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		name := chi.URLParam(r, "resourceGroupRuleName")
		if len(name) == 0 {
			render.Render(w, r, handler.FailureResponse(ctx, errors.New("resource group rule name cannot be empty")))
			return
		}

		logger.Info("Deleting resourceGroupRule...", "resourceGroupRule", name)

		// Use the ResourceGroupManager to delete the resource group rule.
		if err := resourceGroupMgr.DeleteResourceGroupRule(ctx, name); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Render a success response.
		render.JSON(w, r, handler.SuccessResponse(ctx, "ResourceGroupRule deleted successfully"))
	}
}
