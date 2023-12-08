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

	clusterhandler "github.com/KusionStack/karbour/pkg/core/handler/cluster"
	confighandler "github.com/KusionStack/karbour/pkg/core/handler/config"
	resourcehandler "github.com/KusionStack/karbour/pkg/core/handler/resource"
	clustermanager "github.com/KusionStack/karbour/pkg/core/manager/cluster"
	"github.com/KusionStack/karbour/pkg/core/manager/config"
	resourcemanager "github.com/KusionStack/karbour/pkg/core/manager/resource"
	appmiddleware "github.com/KusionStack/karbour/pkg/middleware"
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

func NewCoreServer(
	genericConfig *genericapiserver.CompletedConfig,
	extraConfig *registry.ExtraConfig,
) *chi.Mux {
	router := chi.NewRouter()

	// Set up middlewares
	router.Use(middleware.RequestID)
	router.Use(appmiddleware.AuditLogger)
	router.Use(appmiddleware.APILogger)
	router.Use(middleware.Recoverer)

	// Set up the core api router
	configMgr := config.NewManager(&config.Config{
		Verbose: false,
	})
	clusterMgr := clustermanager.NewClusterManager(&clustermanager.Config{
		Verbose: false,
	})
	resourceMgr := resourcemanager.NewResourceManager(&resourcemanager.ResourceConfig{
		Verbose: false,
	})

	router.Route("/api/v1", func(r chi.Router) {
		setupAPIV1(r, configMgr, clusterMgr, resourceMgr, genericConfig, extraConfig)
	})

	router.Get("/endpoints", func(w http.ResponseWriter, req *http.Request) {
		endpoints := listEndpoints(router)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(endpoints, "\n")))
	})

	return router
}

func setupAPIV1(
	r chi.Router,
	configMgr *config.Manager,
	clusterMgr *clustermanager.ClusterManager,
	resourceMgr *resourcemanager.ResourceManager,
	genericConfig *genericapiserver.CompletedConfig,
	extraConfig *registry.ExtraConfig,
) {
	r.Route("/config", func(r chi.Router) {
		r.Get("/", confighandler.Get(configMgr))
		// r.Delete("/", confighandler.Delete(configMgr))
		// r.Post("/", confighandler.Post(configMgr))
		// r.Put("/", confighandler.Put(configMgr))
	})

	r.Route("/cluster", func(r chi.Router) {
		r.Route("/{clusterName}", func(r chi.Router) {
			r.Get("/", clusterhandler.Get(clusterMgr, genericConfig))
			r.Get("/yaml", clusterhandler.GetYAML(clusterMgr, genericConfig))
			r.Get("/detail", clusterhandler.GetDetail(clusterMgr, genericConfig))
			r.Get("/topology", clusterhandler.GetTopology(clusterMgr, genericConfig))
			r.Get("/namespace/{namespaceName}", clusterhandler.GetNamespace(clusterMgr, genericConfig))
			r.Get("/namespace/{namespaceName}/topology", clusterhandler.GetNamespaceTopology(clusterMgr, genericConfig))
		})
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
}

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
