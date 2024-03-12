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

package route

import (
	"expvar"

	docs "github.com/KusionStack/karbour/api/openapispec"
	clusterhandler "github.com/KusionStack/karbour/pkg/core/handler/cluster"
	detailhandler "github.com/KusionStack/karbour/pkg/core/handler/detail"
	endpointhandler "github.com/KusionStack/karbour/pkg/core/handler/endpoint"
	eventshandler "github.com/KusionStack/karbour/pkg/core/handler/events"
	resourcegrouphandler "github.com/KusionStack/karbour/pkg/core/handler/resourcegroup"
	resourcegrouprulehandler "github.com/KusionStack/karbour/pkg/core/handler/resourcegrouprule"
	scannerhandler "github.com/KusionStack/karbour/pkg/core/handler/scanner"
	searchhandler "github.com/KusionStack/karbour/pkg/core/handler/search"
	summaryhandler "github.com/KusionStack/karbour/pkg/core/handler/summary"
	topologyhandler "github.com/KusionStack/karbour/pkg/core/handler/topology"
	clustermanager "github.com/KusionStack/karbour/pkg/core/manager/cluster"
	insightmanager "github.com/KusionStack/karbour/pkg/core/manager/insight"
	resourcegroupmanager "github.com/KusionStack/karbour/pkg/core/manager/resourcegroup"
	searchmanager "github.com/KusionStack/karbour/pkg/core/manager/search"
	appmiddleware "github.com/KusionStack/karbour/pkg/core/middleware"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/kubernetes/registry"
	"github.com/KusionStack/karbour/pkg/kubernetes/registry/search"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpswagger "github.com/swaggo/http-swagger"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// NewCoreRoute creates and configures an instance of chi.Mux with the given
// configuration and extra configuration parameters.
func NewCoreRoute(
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
	if extraConfig.ReadOnlyMode {
		router.Use(appmiddleware.ReadOnlyMode)
	}

	// Initialize managers, storage for the different core components of the API.
	searchStorage, err := search.NewSearchStorage(*extraConfig)
	if err != nil {
		return nil, err
	}
	insightMgr, err := insightmanager.NewInsightManager(searchStorage)
	if err != nil {
		return nil, err
	}
	resourceGroupMgr := resourcegroupmanager.NewResourceGroupManager(searchStorage)
	clusterMgr := clustermanager.NewClusterManager()
	searchMgr := searchmanager.NewSearchManager()

	// Set up the API routes for version 1 of the API.
	router.Route("/rest-api/v1", func(r chi.Router) {
		setupRestAPIV1(r, clusterMgr, insightMgr, resourceGroupMgr, searchMgr, searchStorage, genericConfig)
	})

	// Set up the root routes.
	docs.SwaggerInfo.BasePath = "/"
	router.Get("/docs/*", httpswagger.Handler())

	// Endpoint to list all available endpoints in the router.
	router.Get("/endpoints", endpointhandler.Endpoints(router))

	// Endpoint to list all available endpoints in the router.
	router.Get("/server-configs", expvar.Handler().ServeHTTP)

	return router, nil
}

// setupRestAPIV1 configures routing for the API version 1, grouping routes by
// resource type and setting up proper handlers.
func setupRestAPIV1(
	r chi.Router,
	clusterMgr *clustermanager.ClusterManager,
	insightMgr *insightmanager.InsightManager,
	resourceGroupMgr *resourcegroupmanager.ResourceGroupManager,
	searchMgr *searchmanager.SearchManager,
	searchStorage storage.SearchStorage,
	genericConfig *genericapiserver.CompletedConfig,
) {
	// Define API routes for 'cluster', 'search', 'resourcegroup' and 'insight', etc.
	r.Route("/clusters", func(r chi.Router) {
		r.Get("/", clusterhandler.List(clusterMgr, genericConfig))
	})

	r.Route("/cluster", func(r chi.Router) {
		r.Route("/{clusterName}", func(r chi.Router) {
			r.Get("/", clusterhandler.Get(clusterMgr, genericConfig))
			r.Post("/", clusterhandler.Create(clusterMgr, genericConfig))
			r.Put("/", clusterhandler.Update(clusterMgr, genericConfig))
			r.Delete("/", clusterhandler.Delete(clusterMgr, genericConfig))
		})
		r.Post("/config/file", clusterhandler.UploadKubeConfig(clusterMgr))
		r.Post("/config/validate", clusterhandler.ValidateKubeConfig(clusterMgr))
	})

	r.Route("/search", func(r chi.Router) {
		r.Get("/", searchhandler.SearchForResource(searchMgr, searchStorage))
	})

	r.Route("/insight", func(r chi.Router) {
		r.Get("/audit", scannerhandler.Audit(insightMgr))
		r.Get("/score", scannerhandler.Score(insightMgr))
		r.Get("/topology", topologyhandler.GetTopology(insightMgr, genericConfig))
		r.Get("/summary", summaryhandler.GetSummary(insightMgr, genericConfig))
		r.Get("/events", eventshandler.GetEvents(insightMgr, genericConfig))
		r.Get("/detail", detailhandler.GetDetail(clusterMgr, insightMgr, genericConfig))
	})

	r.Route("/resource-group-rule", func(r chi.Router) {
		r.Route("/{resourceGroupRuleName}", func(r chi.Router) {
			r.Get("/", resourcegrouprulehandler.Get(resourceGroupMgr))
			r.Post("/", resourcegrouprulehandler.Create(resourceGroupMgr))
			r.Put("/", resourcegrouprulehandler.Update(resourceGroupMgr))
			r.Delete("/", resourcegrouprulehandler.Delete(resourceGroupMgr))
		})
	})
	r.Get("/resource-group/{resourceGroupName}", resourcegrouphandler.Get(resourceGroupMgr))
	r.Get("/resource-groups/", resourcegrouphandler.List(resourceGroupMgr))
}
