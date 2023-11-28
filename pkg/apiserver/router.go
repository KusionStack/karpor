package apiserver

import (
	"net/http"

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
	r := chi.NewRouter()

	// Set up middlewares
	r.Use(middleware.RequestID)
	r.Use(appmiddleware.DefaultLogger)
	r.Use(appmiddleware.APILogger)
	r.Use(middleware.Recoverer)

	// Set up the frontend router
	klog.Infof("Dashboard's static directory use: %s", DefaultStaticDirectory)
	r.NotFound(http.FileServer(http.Dir(DefaultStaticDirectory)).ServeHTTP)

	// Set up the core api router
	configCtrl := config.NewController(&config.Config{
		Verbose: false,
	})

	r.Route("/api/v1", func(r chi.Router) {
		setupApiV1(r, configCtrl)
	})

	return r
}

func setupApiV1(r chi.Router, configCtrl *config.Controller) {
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
