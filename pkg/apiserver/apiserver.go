/*
Copyright The Karbour Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apiserver

import (
	"fmt"
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/server"
	"github.com/KusionStack/karbour/pkg/registry"
	clusterstorage "github.com/KusionStack/karbour/pkg/registry/cluster"
	searchstorage "github.com/KusionStack/karbour/pkg/registry/search"
	"github.com/go-chi/chi/v5"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

// KarbourServer is the carrier of the main process of Karbour.
type KarbourServer struct {
	*genericapiserver.GenericAPIServer
	mux *chi.Mux
	err error
}

// InstallGatewayServer installs various resource-specific REST storage
// implementations on the KarbourServer. This method is part of the
// bootstrapping process of setting up the API server.
func (s *KarbourServer) InstallGatewayServer(c *CompletedConfig) *KarbourServer {
	if s.err != nil {
		return s
	}

	// Initialize REST storage providers for the server.
	restStorageProviders := []registry.RESTStorageProvider{
		clusterstorage.RESTStorageProvider{},
		searchstorage.RESTStorageProvider{
			SearchStorageType:      c.ExtraConfig.SearchStorageType,
			ElasticSearchAddresses: c.ExtraConfig.ElasticSearchAddresses,
			ElasticSearchName:      c.ExtraConfig.ElasticSearchName,
			ElasticSearchPassword:  c.ExtraConfig.ElasticSearchPassword,
		},
	}

	// Attempt to set up each storage provider on the server.
	for _, restStorageProvider := range restStorageProviders {
		groupName := restStorageProvider.GroupName()
		apiGroupInfo, err := restStorageProvider.NewRESTStorage(
			c.GenericConfig.RESTOptionsGetter,
		)
		if err != nil {
			// Capture initialization errors and prevent further setup attempts.
			s.err = fmt.Errorf("problem initializing API group %q: %v", groupName, err)
			return s
		}

		if len(apiGroupInfo.VersionedResourcesStorageMap) == 0 {
			// Skip the setup of this API group if it is effectively disabled
			// (no resources configured).
			klog.Infof("API group %q is not enabled, skipping.", groupName)
			continue
		}

		// Install the API group on the GenericAPIServer.
		if err = s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
			// Capture any errors encountered during installation.
			s.err = fmt.Errorf("problem installing API group %q: %v", groupName, err)
			return s
		}

		klog.Infof("Enabling API group %q.", groupName)
	}

	return s
}

// InstallCoreServer installs the core server (handling non kubernetes-like API,
// regular HTTP requests) onto the KarbourServer. This is typically the server
// that serves the user interface assets.
func (s *KarbourServer) InstallCoreServer(c *CompletedConfig) *KarbourServer {
	if s.err != nil {
		return s
	}

	// Instantiate and set up the core server.
	if mux, err := server.NewCoreServer(&c.GenericConfig, c.ExtraConfig); err == nil {
		s.mux = mux
	} else {
		// Capture any errors encountered during core server setup.
		s.err = err
		return s
	}

	// Mount the core server's Mux to the GenericAPIServer's non-API request
	// handler.
	s.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix("/", s.mux)

	return s
}

// InstallStaticFileServer sets up the server to serve static files.
// It is used to serve files like stylesheets, scripts, and images for the
// karbour dashboard.
func (s *KarbourServer) InstallStaticFileServer() *KarbourServer {
	if s.err != nil {
		return s
	}

	// Define the directory where the static files are located.
	const DefaultStaticDirectory = "./static"

	// Log the use of the static directory for the dashboard.
	klog.Infof("Dashboard's static directory use: %s", DefaultStaticDirectory)

	// Set up the router to serve static files when not found by other routes.
	s.mux.NotFound(http.FileServer(http.Dir(DefaultStaticDirectory)).ServeHTTP)

	return s
}

// Error returns any errors that have occurred during the setup of the
// KarbourServer. It is designed to be called after the configuration steps
// to ensure any issues are captured and reported.
func (s *KarbourServer) Error() error {
	return s.err
}
