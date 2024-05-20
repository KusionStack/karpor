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

package resourcegroup

import (
	"net/http"

	"kusionstack.io/karpor/pkg/core/handler"
	"kusionstack.io/karpor/pkg/core/manager/resourcegroup"
	"kusionstack.io/karpor/pkg/util/ctxutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

// List returns an HTTP handler function that lists all ResourceGroup by
// specified rule. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      List lists all ResourceGroups by rule name.
// @Description  This endpoint lists all ResourceGroups.
// @Tags         resourcegroup
// @Produce      json
// @Param        resourceGroupRuleName  path      string                     true  "The name of the resource group rule"
// @Success      200                    {array}   unstructured.Unstructured  "List of resourceGroup objects"
// @Failure      400                    {string}  string                     "Bad Request"
// @Failure      401                    {string}  string                     "Unauthorized"
// @Failure      404                    {string}  string                     "Not Found"
// @Failure      405                    {string}  string                     "Method Not Allowed"
// @Failure      429                    {string}  string                     "Too Many Requests"
// @Failure      500                    {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-groups/{resourceGroupRuleName} [get]
func List(resourceGroupMgr *resourcegroup.ResourceGroupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		name := chi.URLParam(r, "resourceGroupRuleName")
		if len(name) == 0 {
			render.Render(w, r, handler.FailureResponse(ctx, errors.New("resource group rule name cannot be empty")))
			return
		}

		logger.Info("Listing resourceGroups by resourceGroupRule ...", "resourceGroupRule", name)

		// Use the ResourceGroupManager to list resource groups by specified rule.
		rgs, err := resourceGroupMgr.ListResourceGroupsBy(ctx, name)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Render the list of resource groups.
		render.JSON(w, r, handler.SuccessResponse(ctx, rgs))
	}
}
