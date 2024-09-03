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

package authn

import (
	"net/http"

	"github.com/KusionStack/karpor/pkg/core/handler"
)

// Get returns an HTTP handler that determine whether the user's token can pass authentication.
//
// @Summary      Get returns an authn result of user's token.
// @Description  This endpoint returns an authn result.
// @Tags         authn
// @Produce      json
// @Success      200          {string}  string                     "OK"
// @Failure      400          {string}  string                     "Bad Request"
// @Failure      401          {string}  string                     "Unauthorized"
// @Failure      404          {string}  string                     "Not Found"
// @Failure      405          {string}  string                     "Method Not Allowed"
// @Failure      429          {string}  string                     "Too Many Requests"
// @Failure      500          {string}  string                     "Internal Server Error"
// @Router       /authn [get]
func Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// No action is taken as the API server will handle authentication.
		handler.HandleResult(w, r, r.Context(), nil, "")
	}
}
