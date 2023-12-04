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

	confighandler "github.com/KusionStack/karbour/pkg/apis/config"
	"github.com/KusionStack/karbour/pkg/controller/config"
	appmiddleware "github.com/KusionStack/karbour/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"k8s.io/klog/v2"
)

// DefaultStaticDirectory is the default static directory for
// dashboard.
const DefaultStaticDirectory = "./static"

func NewCoreAPIs() http.Handler {
	router := chi.NewRouter()

	// Set up middlewares
	router.Use(middleware.RequestID)
	router.Use(appmiddleware.AuditLogger)
	router.Use(appmiddleware.APILogger)
	router.Use(middleware.Recoverer)

	// Set up the frontend router
	klog.Infof("Dashboard's static directory use: %s", DefaultStaticDirectory)
	router.NotFound(http.FileServer(http.Dir(DefaultStaticDirectory)).ServeHTTP)

	// Set up the core api router
	configCtrl := config.NewController(&config.Config{
		Verbose: false,
	})

	router.Route("/api/v1", func(r chi.Router) {
		setupAPIV1(r, configCtrl)
	})

	router.Get("/endpoints", func(w http.ResponseWriter, req *http.Request) {
		endpoints := listEndpoints(router)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(endpoints, "\n")))
	})

	return router
}

func setupAPIV1(r chi.Router, configCtrl *config.Controller) {
	r.Route("/config", func(r chi.Router) {
		r.Get("/", confighandler.Get(configCtrl))
		// r.Delete("/", confighandler.Delete(configCtrl))
		// r.Post("/", confighandler.Post(configCtrl))
		// r.Put("/", confighandler.Put(configCtrl))
	})

	// r.Route("/topology", func(r chi.Router) {
	// 	r.Get("/", topologyhandler.Get(topologyCtrl))
	// 	r.Delete("/", topologyhandler.Delete(topologyCtrl))
	// 	r.Post("/", topologyhandler.Post(topologyCtrl))
	// 	r.Put("/", topologyhandler.Put(topologyCtrl))
	// })
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
