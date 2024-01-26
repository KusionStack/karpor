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

package config

import (
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/kubernetes/registry"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/chi/v5"

	genericapiserver "k8s.io/apiserver/pkg/server"
)

// GetServerConfig provides an endpoint to print server configs
//
// @Summary      Print server configs
// @Description  Print server configs
// @Tags         debug
// @Accept       plain
// @Produce      plain
// @Success      200  {string}  string  "Config printed successfully"
// @Router       /server-configs [get]
func GetServerConfig(router chi.Router, genericConfig *genericapiserver.CompletedConfig, extraConfig *registry.ExtraConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		logger.Info("Getting server config...")
		configData, err := MaskSecretInConfig(extraConfig)
		handler.HandleResult(w, r, ctx, err, configData)
	}
}
