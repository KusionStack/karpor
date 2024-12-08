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

package route

import (
	"errors"
	"expvar"

	docs "github.com/KusionStack/karpor/api/openapispec"
	aggregatorhandler "github.com/KusionStack/karpor/pkg/core/handler/aggregator"
	authnhandler "github.com/KusionStack/karpor/pkg/core/handler/authn"
	clusterhandler "github.com/KusionStack/karpor/pkg/core/handler/cluster"
	detailhandler "github.com/KusionStack/karpor/pkg/core/handler/detail"
	endpointhandler "github.com/KusionStack/karpor/pkg/core/handler/endpoint"
	eventshandler "github.com/KusionStack/karpor/pkg/core/handler/events"
	resourcegrouphandler "github.com/KusionStack/karpor/pkg/core/handler/resourcegroup"
	resourcegrouprulehandler "github.com/KusionStack/karpor/pkg/core/handler/resourcegrouprule"
	scannerhandler "github.com/KusionStack/karpor/pkg/core/handler/scanner"
	searchhandler "github.com/KusionStack/karpor/pkg/core/handler/search"
	statshandler "github.com/KusionStack/karpor/pkg/core/handler/stats"
	summaryhandler "github.com/KusionStack/karpor/pkg/core/handler/summary"
	topologyhandler "github.com/KusionStack/karpor/pkg/core/handler/topology"
	healthhandler "github.com/KusionStack/karpor/pkg/core/health"
	aimanager "github.com/KusionStack/karpor/pkg/core/manager/ai"
	clustermanager "github.com/KusionStack/karpor/pkg/core/manager/cluster"
	insightmanager "github.com/KusionStack/karpor/pkg/core/manager/insight"
	resourcegroupmanager "github.com/KusionStack/karpor/pkg/core/manager/resourcegroup"
	searchmanager "github.com/KusionStack/karpor/pkg/core/manager/search"
	appmiddleware "github.com/KusionStack/karpor/pkg/core/middleware"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry/search"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpswagger "github.com/swaggo/http-swagger/v2"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
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
	resourceStorage, err := search.NewResourceStorage(*extraConfig)
	if err != nil {
		return nil, err
	}
	resourceGroupRuleStorage, err := search.NewResourceGroupRuleStorage(*extraConfig)
	if err != nil {
		return nil, err
	}
	generalStorage, err := search.NewGeneralStorage(*extraConfig)
	if err != nil {
		return nil, err
	}

	insightMgr, err := insightmanager.NewInsightManager(searchStorage, resourceStorage, resourceGroupRuleStorage, genericConfig)
	if err != nil {
		return nil, err
	}
	resourceGroupMgr, err := resourcegroupmanager.NewResourceGroupManager(resourceGroupRuleStorage)
	if err != nil {
		return nil, err
	}
	aiMgr, err := aimanager.NewAIManager(*extraConfig)
	if err != nil {
		if errors.Is(err, aimanager.ErrMissingAuthToken) {
			klog.Warning("Auth token is empty.")
		} else {
			return nil, err
		}
	}

	clusterMgr := clustermanager.NewClusterManager()
	searchMgr := searchmanager.NewSearchManager()

	// Set up the API routes for version 1 of the API.
	router.Route("/rest-api/v1", func(r chi.Router) {
		setupRestAPIV1(r,
			aiMgr,
			clusterMgr,
			insightMgr,
			resourceGroupMgr,
			searchMgr,
			searchStorage,
			genericConfig)
	})

	// Set up the root routes.
	docs.SwaggerInfo.BasePath = "/"
	router.Get("/docs/*", httpswagger.Handler())

	// Endpoint to list all available endpoints in the router.
	router.Get("/endpoints", endpointhandler.Endpoints(router))

	// Expose server configuration and runtime statistics.
	router.Get("/server-configs", expvar.Handler().ServeHTTP)

	healthhandler.Register(router, generalStorage)
	return router, nil
}

// setupRestAPIV1 configures routing for the API version 1, grouping routes by
// resource type and setting up proper handlers.
func setupRestAPIV1(
	r chi.Router,
	aiMgr *aimanager.AIManager,
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
		r.Get("/", searchhandler.SearchForResource(searchMgr, aiMgr, searchStorage))
	})

	r.Route("/insight", func(r chi.Router) {
		r.Get("/stats", statshandler.GetStatistics(insightMgr))
		r.Get("/audit", scannerhandler.Audit(insightMgr))
		r.Get("/score", scannerhandler.Score(insightMgr))
		r.Get("/topology", topologyhandler.GetTopology(clusterMgr, insightMgr, genericConfig))
		r.Get("/summary", summaryhandler.GetSummary(insightMgr, genericConfig))
		r.Get("/events", eventshandler.GetEvents(insightMgr, genericConfig))
		r.Get("/detail", detailhandler.GetDetail(clusterMgr, insightMgr, genericConfig))
		r.Get("/aggregator/pod/{cluster}/{namespace}/{name}/log", aggregatorhandler.GetPodLogs(clusterMgr, genericConfig))
		r.Get("/aggregator/event/{cluster}/{namespace}/{name}", aggregatorhandler.GetEvents(clusterMgr, genericConfig))
	})

	r.Route("/resource-group-rule", func(r chi.Router) {
		r.Get("/{resourceGroupRuleName}", resourcegrouprulehandler.Get(resourceGroupMgr))
		r.Post("/", resourcegrouprulehandler.Create(resourceGroupMgr))
		r.Put("/", resourcegrouprulehandler.Update(resourceGroupMgr))
		r.Delete("/{resourceGroupRuleName}", resourcegrouprulehandler.Delete(resourceGroupMgr))
	})
	r.Get("/resource-group-rules", resourcegrouprulehandler.List(resourceGroupMgr))
	r.Get("/resource-groups/{resourceGroupRuleName}", resourcegrouphandler.List(resourceGroupMgr))
	r.Get("/authn", authnhandler.Get())
}
