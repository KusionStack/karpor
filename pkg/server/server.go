/*
Copyright The Karpor Authors.

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

package server

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/KusionStack/karpor/pkg/core/route"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	clusterstorage "github.com/KusionStack/karpor/pkg/kubernetes/registry/cluster"
	corestorage "github.com/KusionStack/karpor/pkg/kubernetes/registry/core"
	searchstorage "github.com/KusionStack/karpor/pkg/kubernetes/registry/search"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"github.com/KusionStack/karpor/ui"
	"github.com/go-chi/chi/v5"
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/klog/v2"
	rbacrest "k8s.io/kubernetes/pkg/registry/rbac/rest"
)

// KarporServer is the carrier of the main process of Karpor.
type KarporServer struct {
	*genericapiserver.GenericAPIServer
	mux *chi.Mux
	err error
}

// InstallKubernetesServer installs various resource-specific REST storage
// implementations on the KarporServer. This method is part of the
// bootstrapping process of setting up the API server.
func (s *KarporServer) InstallKubernetesServer(c *CompletedConfig) *KarporServer {
	if s.err != nil {
		return s
	}
	err := s.InstallLegacyAPI(c)
	if err != nil {
		s.err = err
		return s
	}

	// Initialize REST storage providers for the server.
	restStorageProviders := []registry.RESTStorageProvider{
		clusterstorage.RESTStorageProvider{},
		searchstorage.RESTStorageProvider{
			SearchStorageType: c.ExtraConfig.SearchStorageType,
			SearchAddresses:   c.ExtraConfig.SearchAddresses,
			SearchName:        c.ExtraConfig.SearchUsername,
			SearchPassword:    c.ExtraConfig.SearchPassword,
		},
		rbacrest.RESTStorageProvider{Authorizer: c.GenericConfig.Authorization.Authorizer},
	}
	apiResourceConfigSource := serverstorage.NewResourceConfig()
	apiResourceConfigSource.EnableVersions(scheme.Versions...)
	err = s.InstallAPIs(
		apiResourceConfigSource,
		c.GenericConfig.RESTOptionsGetter,
		restStorageProviders...)
	if err != nil {
		s.err = err
	}
	return s
}

// InstallCoreServer installs the core server (handling non kubernetes-like API,
// regular HTTP requests) onto the KarporServer. This is typically the server
// that serves the user interface assets.
func (s *KarporServer) InstallCoreServer(c *CompletedConfig) *KarporServer {
	if s.err != nil {
		return s
	}

	// Instantiate and set up the core route.
	if mux, err := route.NewCoreRoute(&c.GenericConfig, c.ExtraConfig); err == nil {
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

// InstallPublicFileServer sets up the server to serve public files.
// It is used to serve files like stylesheets, scripts, and images for the
// karpor dashboard.
func (s *KarporServer) InstallPublicFileServer() *KarporServer {
	if s.err != nil {
		return s
	}

	// Get the web root of dashboard from embedded filesystem.
	webRootFS, err := fs.Sub(ui.Embedded, "build")
	if err != nil {
		klog.Warningf(
			"Failed to get web root directory from embedded filesystem as %s",
			err.Error(),
		)
	}

	// Set up the router to serve static files when not found by other routes.
	s.mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := webRootFS.Open("index.html")
		if err != nil {
			klog.Warningf(
				"Failed to open dashboard index.html from embedded filesystem as %s",
				err.Error(),
			)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		b, err := io.ReadAll(indexFile)
		if err != nil {
			klog.Warningf("Failed to read dashboard index.html as %s", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	publicHandler := http.StripPrefix("/public", http.FileServer(http.FS(webRootFS)))
	s.mux.Mount("/public", publicHandler)

	return s
}

// Error returns any errors that have occurred during the setup of the
// KarporServer. It is designed to be called after the configuration steps
// to ensure any issues are captured and reported.
func (s *KarporServer) Error() error {
	return s.err
}

// InstallLegacyAPI installs legacy API groups and resources into the server.
func (s *KarporServer) InstallLegacyAPI(c *CompletedConfig) error {
	// Installing core API group
	coreProvider := corestorage.NewRESTStorageProvider(c.ExtraConfig.ServiceAccountIssuer, c.ExtraConfig.ServiceAccountMaxExpiration)
	coreGroupName := coreProvider.GroupName()
	coreGroupInfo, err := coreProvider.NewRESTStorage(c.GenericConfig.RESTOptionsGetter)
	if err != nil {
		return fmt.Errorf("problem initializing legacy API group %q: %v", coreGroupName, err)
	}
	if err = s.GenericAPIServer.InstallLegacyAPIGroup(genericapiserver.DefaultLegacyAPIPrefix, &coreGroupInfo); err != nil {
		return fmt.Errorf("problem installing legacy API group %q: %v", coreGroupName, err)
	}
	klog.Infof("Enabling legacy API group %q.", coreGroupName)
	return nil
}

func (s *KarporServer) InstallAPIs(
	apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter,
	restStorageProviders ...registry.RESTStorageProvider,
) error {
	apiGroupsInfo := []*genericapiserver.APIGroupInfo{}
	for _, restStorageProvider := range restStorageProviders {
		groupName := restStorageProvider.GroupName()
		apiGroupInfo, err := restStorageProvider.NewRESTStorage(
			apiResourceConfigSource,
			restOptionsGetter,
		)
		if err != nil {
			return fmt.Errorf("problem initializing API group %q: %v", groupName, err)
		}
		if len(apiGroupInfo.VersionedResourcesStorageMap) == 0 {
			// Skip the setup of this API group if it is effectively disabled
			// (no resources configured).
			klog.Infof("API group %q is not enabled, skipping.", groupName)
			continue
		}

		klog.Infof("Enabling API group %q.", groupName)

		if postHookProvider, ok := restStorageProvider.(genericapiserver.PostStartHookProvider); ok {
			name, hook, err := postHookProvider.PostStartHook()
			if err != nil {
				klog.Fatalf("Error building PostStartHook: %v", err)
			}
			s.GenericAPIServer.AddPostStartHookOrDie(name, hook)
		}

		apiGroupsInfo = append(apiGroupsInfo, &apiGroupInfo)
	}

	if err := s.GenericAPIServer.InstallAPIGroups(apiGroupsInfo...); err != nil {
		return fmt.Errorf("error in registering group versions: %v", err)
	}
	return nil
}
