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

package resourcegroup

import (
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/manager/resourcegroup"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
)

// List returns an HTTP handler function that lists all ResourceGroup
// resources. It utilizes a ResourceGroupManager to execute the logic.
//
// @Summary      List lists all ResourceGroups.
// @Description  This endpoint lists all ResourceGroups.
// @Tags         resourcegroup
// @Produce      json
// @Param        summary     query     bool                       false  "Whether to display summary or not. Default to false"
// @Param        orderBy     query     string                     false  "The order to list the resourceGroup. Default to order by name"
// @Param        descending  query     bool                       false  "Whether to sort the list in descending order. Default to false"
// @Success      200         {array}   unstructured.Unstructured  "List of resourceGroup objects"
// @Failure      400         {string}  string                     "Bad Request"
// @Failure      401         {string}  string                     "Unauthorized"
// @Failure      404         {string}  string                     "Not Found"
// @Failure      405         {string}  string                     "Method Not Allowed"
// @Failure      429         {string}  string                     "Too Many Requests"
// @Failure      500         {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/resource-groups/{resourceGroupRuleName} [get]
func List(resourceGroupMgr *resourcegroup.ResourceGroupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		logger.Info("Listing resourceGroups by resourceGroupRule ...")

		// TODO: need to implement
	}
}
