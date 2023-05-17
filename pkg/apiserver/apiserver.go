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
	"net/http"

	"github.com/KusionStack/karbour/pkg/registry"
	clusterstorage "github.com/KusionStack/karbour/pkg/registry/cluster"
	searchstorage "github.com/KusionStack/karbour/pkg/registry/search"
	"k8s.io/apimachinery/pkg/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

// DefaultStaticDirectory is the default static directory for
// dashboard.
const DefaultStaticDirectory = "./static"

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	SearchStorageType      string
	ElasticSearchAddresses []string
	ElasticSearchName      string
	ElasticSearchPassword  string
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   *ExtraConfig
}

// APIServer contains state for a Kubernetes cluster master/api server.
type APIServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		cfg.ExtraConfig,
	}

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

// New returns a new instance of APIServer from the given config.
func (c completedConfig) New() (*APIServer, error) {
	genericServer, err := c.GenericConfig.New("karbour-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}
	s := &APIServer{
		GenericAPIServer: genericServer,
	}

	if err := InstallLegacyAPI(s.GenericAPIServer, c.GenericConfig.RESTOptionsGetter); err != nil {
		return nil, err
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
	if err := InstallAPIs(s.GenericAPIServer, c.GenericConfig.RESTOptionsGetter, restStorageProviders...); err != nil {
		return nil, err
	}

	klog.Infof("Dashboard's static directory use: %s", DefaultStaticDirectory)
	s.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix("/", http.FileServer(http.Dir(DefaultStaticDirectory)))

	return s, nil
}
