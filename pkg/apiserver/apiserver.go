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

// KarbourServer contains state for a Kubernetes cluster master/api server.
type KarbourServer struct {
	*genericapiserver.GenericAPIServer
	mux *chi.Mux
	err error
}

func (s *KarbourServer) InstallGatewayServer(c *CompletedConfig) *KarbourServer {
	if s.err != nil {
		return s
	}

	restStorageProviders := []registry.RESTStorageProvider{
		clusterstorage.RESTStorageProvider{},
		searchstorage.RESTStorageProvider{
			SearchStorageType:      c.ExtraConfig.SearchStorageType,
			ElasticSearchAddresses: c.ExtraConfig.ElasticSearchAddresses,
			ElasticSearchName:      c.ExtraConfig.ElasticSearchName,
			ElasticSearchPassword:  c.ExtraConfig.ElasticSearchPassword,
		},
	}

	for _, restStorageProvider := range restStorageProviders {
		groupName := restStorageProvider.GroupName()
		apiGroupInfo, err := restStorageProvider.NewRESTStorage(c.GenericConfig.RESTOptionsGetter)
		if err != nil {
			s.err = fmt.Errorf("problem initializing API group %q : %v", groupName, err)
			return s
		}

		if len(apiGroupInfo.VersionedResourcesStorageMap) == 0 {
			// If we have no storage for any resource configured, this API group is effectively disabled.
			// This can happen when an entire API group, version, or development-stage (alpha, beta, GA) is disabled.
			klog.Infof("API group %q is not enabled, skipping.", groupName)
			continue
		}

		if err = s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
			s.err = fmt.Errorf("problem install API group %q: %v", groupName, err)
			return s
		}

		klog.Infof("Enabling API group %q.", groupName)
	}

	return s
}

func (s *KarbourServer) InstallCoreServer(c *CompletedConfig) *KarbourServer {
	if s.err != nil {
		return s
	}

	// Create the core server.
	if mux, err := server.NewCoreServer(&c.GenericConfig, c.ExtraConfig); err == nil {
		s.mux = mux
	} else {
		s.err = err
		return s
	}

	// Mount the core mux to NonGoRestfulMux of GenericAPIServer.
	s.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix("/", s.mux)

	return s
}

func (s *KarbourServer) InstallStaticFileServer() *KarbourServer {
	if s.err != nil {
		return s
	}

	// DefaultStaticDirectory is the default static directory for
	// dashboard.
	const DefaultStaticDirectory = "./static"

	// Set up the frontend router
	klog.Infof("Dashboard's static directory use: %s", DefaultStaticDirectory)
	s.mux.NotFound(http.FileServer(http.Dir(DefaultStaticDirectory)).ServeHTTP)

	return s
}

func (s *KarbourServer) Error() error {
	return s.err
}
