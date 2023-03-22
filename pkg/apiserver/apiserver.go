/*
Copyright The Alipay Authors.

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

	"code.alipay.com/multi-cluster/karbour/pkg/registry"
	clusterstorage "code.alipay.com/multi-cluster/karbour/pkg/registry/cluster/clusterextension"
	searchstorage "code.alipay.com/multi-cluster/karbour/pkg/registry/search/clusterextension"

	"k8s.io/apimachinery/pkg/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	// Place you custom config here.
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
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
		&cfg.ExtraConfig,
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

	restStorageProviders := []registry.RESTStorageProvider{
		clusterstorage.RESTStorageProvider{},
		searchstorage.RESTStorageProvider{},
	}

	for _, restStorageProvider := range restStorageProviders {
		groupName := restStorageProvider.GroupName()
		apiGroupInfo, err := restStorageProvider.NewRESTStorage(c.GenericConfig.RESTOptionsGetter)
		if err != nil {
			return nil, fmt.Errorf("problem initializing API group %q : %v", groupName, err)
		}

		if len(apiGroupInfo.VersionedResourcesStorageMap) == 0 {
			// If we have no storage for any resource configured, this API group is effectively disabled.
			// This can happen when an entire API group, version, or development-stage (alpha, beta, GA) is disabled.
			klog.Infof("API group %q is not enabled, skipping.", groupName)
			continue
		}

		if err = s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
			return nil, fmt.Errorf("problem install API group %q: %v", groupName, err)
		}

		klog.Infof("Enabling API group %q.", groupName)
	}

	return s, nil
}
