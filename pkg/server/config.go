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
	"github.com/KusionStack/karbour/pkg/kubernetes/registry"
	"k8s.io/apimachinery/pkg/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   *registry.ExtraConfig
}

// Complete fills in any fields not set that are required to have valid data.
// It's mutating the receiver.
func (cfg *Config) Complete() *CompletedConfig {
	c := &CompletedConfig{
		cfg.GenericConfig.Complete(),
		cfg.ExtraConfig,
	}

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return c
}

type CompletedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *registry.ExtraConfig
}

// New returns a new instance of APIServer from the given config.
func (c *CompletedConfig) New() (*KarbourServer, error) {
	genericServer, err := c.GenericConfig.New(
		"karbour-apiserver",
		genericapiserver.NewEmptyDelegate(),
	)
	if err != nil {
		return nil, err
	}

	s := &KarbourServer{
		GenericAPIServer: genericServer,
	}
	if err := s.InstallKubernetesServer(c).
		InstallCoreServer(c).
		InstallStaticFileServer().
		Error(); err != nil {
		return nil, err
	}

	return s, nil
}
