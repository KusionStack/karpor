package apiserver

import (
	"github.com/KusionStack/karbour/pkg/registry"
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
	genericServer, err := c.GenericConfig.New("karbour-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &KarbourServer{
		GenericAPIServer: genericServer,
	}
	if err := s.InstallGatewayServer(c).
		InstallCoreServer(c).
		InstallStaticFileServer().
		Error(); err != nil {
		return nil, err
	}

	return s, nil
}
