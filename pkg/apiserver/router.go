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

package apiserver

import (
	"fmt"
	"net/http"
	"strings"

	clustercontroller "github.com/KusionStack/karbour/pkg/controller/cluster"
	"github.com/KusionStack/karbour/pkg/controller/config"
	resourcecontroller "github.com/KusionStack/karbour/pkg/controller/resource"
	clusterhandler "github.com/KusionStack/karbour/pkg/handler/cluster"
	confighandler "github.com/KusionStack/karbour/pkg/handler/config"
	resourcehandler "github.com/KusionStack/karbour/pkg/handler/resource"
	appmiddleware "github.com/KusionStack/karbour/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewCoreServer(c *CompletedConfig) *chi.Mux {
	router := chi.NewRouter()

	// Set up middlewares
	router.Use(middleware.RequestID)
	router.Use(appmiddleware.AuditLogger)
	router.Use(appmiddleware.APILogger)
	router.Use(middleware.Recoverer)

	// Set up the core api router
	configCtrl := config.NewController(&config.Config{
		Verbose: false,
	})
	clusterCtrl := clustercontroller.NewClusterController(&clustercontroller.Config{
		Verbose: false,
	})
	resourceCtrl := resourcecontroller.NewResourceController(&resourcecontroller.ResourceConfig{
		Verbose: false,
	})

	router.Route("/api/v1", func(r chi.Router) {
		setupAPIV1(r, configCtrl, clusterCtrl, resourceCtrl, c)
	})

	router.Get("/endpoints", func(w http.ResponseWriter, req *http.Request) {
		endpoints := listEndpoints(router)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(endpoints, "\n")))
	})

	return router
}

func setupAPIV1(r chi.Router, configCtrl *config.Controller, clusterCtrl *clustercontroller.ClusterController, resourceCtrl *resourcecontroller.ResourceController, c *CompletedConfig) {
	r.Route("/config", func(r chi.Router) {
		r.Get("/", confighandler.Get(configCtrl))
		// r.Delete("/", confighandler.Delete(configCtrl))
		// r.Post("/", confighandler.Post(configCtrl))
		// r.Put("/", confighandler.Put(configCtrl))
	})

	r.Route("/cluster", func(r chi.Router) {
		r.Route("/{clusterName}", func(r chi.Router) {
			r.Get("/", clusterhandler.Get(clusterCtrl, &c.GenericConfig))
			r.Get("/yaml", clusterhandler.GetYAML(clusterCtrl, &c.GenericConfig))
			r.Get("/topology", clusterhandler.GetTopology(clusterCtrl, &c.GenericConfig))
			r.Get("/namespace/{namespaceName}/topology", clusterhandler.GetNamespaceTopology(clusterCtrl, &c.GenericConfig))
		})
	})

	r.Route("/resource", func(r chi.Router) {
		r.Route("/search", func(r chi.Router) {
			r.Get("/", resourcehandler.SearchForResource(resourceCtrl, c.ExtraConfig))
		})
		r.Route("/cluster/{clusterName}/{apiVersion}/namespace/{namespaceName}/{kind}/name/{resourceName}", func(r chi.Router) {
			r.Get("/", resourcehandler.Get(resourceCtrl, &c.GenericConfig))
			r.Get("/yaml", resourcehandler.GetYAML(resourceCtrl, &c.GenericConfig))
			r.Get("/topology", resourcehandler.GetTopology(resourceCtrl, &c.GenericConfig))
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
