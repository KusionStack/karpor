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
	"sort"
	"strings"

	docs "github.com/KusionStack/karbour/api/openapispec"
	audithandler "github.com/KusionStack/karbour/pkg/core/handler/audit"
	clusterhandler "github.com/KusionStack/karbour/pkg/core/handler/cluster"
	detailhandler "github.com/KusionStack/karbour/pkg/core/handler/detail"
	eventshandler "github.com/KusionStack/karbour/pkg/core/handler/events"
	searchhandler "github.com/KusionStack/karbour/pkg/core/handler/search"
	summaryhandler "github.com/KusionStack/karbour/pkg/core/handler/summary"
	topologyhandler "github.com/KusionStack/karbour/pkg/core/handler/topology"
	auditmanager "github.com/KusionStack/karbour/pkg/core/manager/audit"
	clustermanager "github.com/KusionStack/karbour/pkg/core/manager/cluster"
	insightmanager "github.com/KusionStack/karbour/pkg/core/manager/insight"
	searchmanager "github.com/KusionStack/karbour/pkg/core/manager/search"
	appmiddleware "github.com/KusionStack/karbour/pkg/middleware"
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/KusionStack/karbour/pkg/registry/search"
	"github.com/KusionStack/karbour/pkg/search/storage"
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

	// Initialize managers, storage for the different core components of the API.
	searchStorage, err := search.NewSearchStorage(*extraConfig)
	if err != nil {
		return nil, err
	}
	clusterMgr := clustermanager.NewClusterManager(&clustermanager.ClusterConfig{
		Verbose: false,
	})
	insightMgr := insightmanager.NewInsightManager(&insightmanager.InsightConfig{
		Verbose: false,
	})
	searchMgr := searchmanager.NewSearchManager(&searchmanager.SearchConfig{
		Verbose: false,
	})
	auditMgr, err := auditmanager.NewAuditManager(searchStorage)
	if err != nil {
		return nil, err
	}

	// Set up the root routes.
	docs.SwaggerInfo.BasePath = "/"
	router.Get("/docs/*", httpswagger.Handler())

	// Set up the API routes for version 1 of the API.
	router.Route("/api/v1", func(r chi.Router) {
		setupAPIV1(r, clusterMgr, insightMgr, searchMgr, auditMgr, searchStorage, genericConfig)
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
	clusterMgr *clustermanager.ClusterManager,
	insightMgr *insightmanager.InsightManager,
	searchMgr *searchmanager.SearchManager,
	auditMgr *auditmanager.AuditManager,
	searchStorage storage.SearchStorage,
	genericConfig *genericapiserver.CompletedConfig,
) {
	// Define API routes for 'cluster', 'search', and 'insight', etc.
	r.Route("/clusters", func(r chi.Router) {
		r.Get("/", clusterhandler.List(clusterMgr, genericConfig))
	})

	r.Route("/cluster", func(r chi.Router) {
		r.Route("/{clusterName}", func(r chi.Router) {
			r.Get("/", clusterhandler.Get(clusterMgr, genericConfig))
			r.Post("/", clusterhandler.Create(clusterMgr, genericConfig))
			r.Put("/", clusterhandler.UpdateMetadata(clusterMgr, genericConfig))
			r.Delete("/", clusterhandler.Delete(clusterMgr, genericConfig))
		})
		r.Post("/config/file", clusterhandler.UpdateKubeConfig)
		r.Post("/config/validate", clusterhandler.ValidateKubeConfig(clusterMgr))
	})

	r.Route("/search", func(r chi.Router) {
		r.Get("/", searchhandler.SearchForResource(searchMgr, searchStorage))
	})

	r.Route("/insight", func(r chi.Router) {
		r.Get("/audit", audithandler.Audit(auditMgr))
		r.Get("/score", audithandler.Score(auditMgr))
		r.Get("/topology", topologyhandler.GetTopology(insightMgr, genericConfig))
		r.Get("/summary", summaryhandler.GetSummary(insightMgr, genericConfig))
		r.Get("/events", eventshandler.GetEvents(insightMgr, genericConfig))
		r.Get("/detail", detailhandler.GetDetail(insightMgr, genericConfig))
	})
}

// listEndpoints generates a list of all routes registered in the router.
func listEndpoints(r chi.Router) []string {
	var endpoints []string
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		endpoint := fmt.Sprintf("%s\t%s", method, route)
		endpoints = append(endpoints, endpoint)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Walking routes error: %s\n", err.Error())
	}
	sort.Strings(endpoints)
	return endpoints
}
