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

package server

import (
	"fmt"
	"net/http"
	"strings"

	docs "github.com/KusionStack/karbour/api/openapispec"
	audithandler "github.com/KusionStack/karbour/pkg/core/handler/audit"
	clusterhandler "github.com/KusionStack/karbour/pkg/core/handler/cluster"
	confighandler "github.com/KusionStack/karbour/pkg/core/handler/config"
	resourcehandler "github.com/KusionStack/karbour/pkg/core/handler/resource"
	auditmanager "github.com/KusionStack/karbour/pkg/core/manager/audit"
	clustermanager "github.com/KusionStack/karbour/pkg/core/manager/cluster"
	"github.com/KusionStack/karbour/pkg/core/manager/config"
	resourcemanager "github.com/KusionStack/karbour/pkg/core/manager/resource"
	appmiddleware "github.com/KusionStack/karbour/pkg/middleware"
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpswagger "github.com/swaggo/http-swagger"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// NewCoreServer creates and configures an instance of chi.Mux with the given
// configuration and extra configuration parameters.
func NewCoreServer(
	genericConfig *genericapiserver.CompletedConfig,
	extraConfig *registry.ExtraConfig,
) (*chi.Mux, error) {
	router := chi.NewRouter()

	// Set up middlewares for logging, recovery, and timing, etc.
	router.Use(middleware.RequestID)
	router.Use(appmiddleware.DefaultLogger)
	router.Use(appmiddleware.APILogger)
	router.Use(appmiddleware.Timing)
	router.Use(middleware.Recoverer)

	// Initialize managers for the different core components of the API.
	configMgr := config.NewManager(&config.Config{
		Verbose: false,
	})
	clusterMgr := clustermanager.NewClusterManager(&clustermanager.Config{
		Verbose: false,
	})
	resourceMgr := resourcemanager.NewResourceManager(&resourcemanager.ResourceConfig{
		Verbose: false,
	})
	auditMgr, err := auditmanager.NewAuditManager()
	if err != nil {
		return nil, err
	}

	// Set up the root routes.
	docs.SwaggerInfo.BasePath = "/"
	router.Route("/", func(r chi.Router) {
		r.Get("/docs/*", httpswagger.Handler())
	})

	// Set up the API routes for version 1 of the API.
	router.Route("/api/v1", func(r chi.Router) {
		setupAPIV1(r, configMgr, clusterMgr, resourceMgr, auditMgr, genericConfig, extraConfig)
	})

	// Endpoint to list all available endpoints in the router.
	router.Get("/endpoints", func(w http.ResponseWriter, req *http.Request) {
		endpoints := listEndpoints(router)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(endpoints, "\n")))
	})

	return router, nil
}

// setupAPIV1 configures routing for the API version 1, grouping routes by
// resource type and setting up proper handlers.
func setupAPIV1(
	r chi.Router,
	configMgr *config.Manager,
	clusterMgr *clustermanager.ClusterManager,
	resourceMgr *resourcemanager.ResourceManager,
	auditMgr *auditmanager.AuditManager,
	genericConfig *genericapiserver.CompletedConfig,
	extraConfig *registry.ExtraConfig,
) {
	// Define API routes for 'config', 'cluster', 'resource', and 'audit', etc.
	r.Route("/config", func(r chi.Router) {
		r.Get("/", confighandler.Get(configMgr))
		// r.Delete("/", confighandler.Delete(configMgr))
		// r.Post("/", confighandler.Post(configMgr))
		// r.Put("/", confighandler.Put(configMgr))
	})

	r.Route("/cluster", func(r chi.Router) {
		// Define cluster specific routes.
		r.Route("/", func(r chi.Router) {
			r.Post("/", clusterhandler.Create(clusterMgr, genericConfig))
			r.Get("/", clusterhandler.List(clusterMgr, genericConfig))
		})
		r.Route("/{clusterName}", func(r chi.Router) {
			r.Get("/", clusterhandler.Get(clusterMgr, genericConfig))
			r.Put("/", clusterhandler.UpdateMetadata(clusterMgr, genericConfig))
			r.Delete("/", clusterhandler.Delete(clusterMgr, genericConfig))
			r.Get("/yaml", clusterhandler.GetYAML(clusterMgr, genericConfig))
			r.Get("/detail", clusterhandler.GetDetail(clusterMgr, genericConfig))
			r.Get("/topology", clusterhandler.GetTopology(clusterMgr, genericConfig))
			r.Get("/namespace/{namespaceName}", clusterhandler.GetNamespace(clusterMgr, genericConfig))
			r.Get("/namespace/{namespaceName}/topology", clusterhandler.GetNamespaceTopology(clusterMgr, genericConfig))
		})
		r.Post("/config/file", clusterhandler.UpdateKubeConfig)
		r.Post("/config/validate", clusterhandler.ValidateKubeConfig(clusterMgr))
	})

	r.Route("/resource", func(r chi.Router) {
		r.Route("/search", func(r chi.Router) {
			r.Get("/", resourcehandler.SearchForResource(resourceMgr, extraConfig))
		})
		r.Route("/cluster/{clusterName}/{apiVersion}/namespace/{namespaceName}/{kind}/name/{resourceName}", func(r chi.Router) {
			r.Get("/", resourcehandler.Get(resourceMgr, genericConfig))
			r.Get("/yaml", resourcehandler.GetYAML(resourceMgr, genericConfig))
			r.Get("/topology", resourcehandler.GetTopology(resourceMgr, genericConfig))
		})
	})

	r.Route("/audit", func(r chi.Router) {
		r.Post("/", audithandler.Audit(auditMgr))
		r.Post("/score", audithandler.Score(auditMgr))
	})
}

// listEndpoints generates a list of all routes registered in the router.
func listEndpoints(r chi.Router) []string {
	var endpoints []string
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		endpoint := fmt.Sprintf("%s %s", method, route)
		endpoints = append(endpoints, endpoint)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Walking routes error: %s\n", err.Error())
	}
	return endpoints
}
