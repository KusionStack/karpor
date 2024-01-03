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

package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/KusionStack/karbour/cmd/app/options"
	karbouropenapi "github.com/KusionStack/karbour/pkg/kubernetes/generated/openapi"
	k8sopenapi "github.com/KusionStack/karbour/pkg/kubernetes/openapi"
	"github.com/KusionStack/karbour/pkg/kubernetes/registry"
	"github.com/KusionStack/karbour/pkg/kubernetes/scheme"
	"github.com/KusionStack/karbour/pkg/server"
	proxyutil "github.com/KusionStack/karbour/pkg/util/proxy"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/discovery"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/filters"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/pkg/common"
	netutils "k8s.io/utils/net"
)

const defaultEtcdPathPrefix = "/registry/karbour"

// Options contains state for master/api server
type Options struct {
	RecommendedOptions   *options.RecommendedOptions
	SearchStorageOptions *options.SearchStorageOptions

	StdOut io.Writer
	StdErr io.Writer

	AlternateDNS []string
}

// NewOptions returns a new Options
func NewOptions(out, errOut io.Writer) (*Options, error) {
	o := &Options{
		RecommendedOptions: options.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			scheme.Codecs.LegacyCodec(scheme.Versions...),
		),
		SearchStorageOptions: options.NewSearchStorageOptions(),
		StdOut:               out,
		StdErr:               errOut,
	}
	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = schema.GroupVersions(scheme.Versions)
	o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)
	// TODO: have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"localhost", o.AlternateDNS, []net.IP{netutils.ParseIPSloppy("127.0.0.1")},
	); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}
	return o, nil
}

// NewServerCommand provides a CLI handler for 'start master' command
// with a default Options.
func NewServerCommand(ctx context.Context) *cobra.Command {
	o, err := NewOptions(os.Stdout, os.Stderr)
	if err != nil {
		klog.Background().Error(err, "Unable to initialize command options")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	cmd := &cobra.Command{
		Short: "Launch an API server",
		Long:  "Launch an API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunServer(ctx.Done()); err != nil {
				return err
			}
			return nil
		},
	}

	o.AddFlags(cmd.Flags())
	return cmd
}

// AddFlags add flags to command
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.RecommendedOptions.AddFlags(fs)
	o.SearchStorageOptions.AddFlags(fs)
}

// Validate validates Options
func (o *Options) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	errors = append(errors, o.SearchStorageOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in fields required to have valid data
func (o *Options) Complete() error {
	return nil
}

// Config returns config for the api server given Options
func (o *Options) Config() (*server.Config, error) {
	config := &server.Config{
		GenericConfig: genericapiserver.NewRecommendedConfig(scheme.Codecs),
		ExtraConfig:   &registry.ExtraConfig{},
	}
	if err := o.RecommendedOptions.ApplyTo(config.GenericConfig); err != nil {
		return nil, err
	}
	if err := o.SearchStorageOptions.ApplyTo(config.ExtraConfig); err != nil {
		return nil, err
	}

	config.GenericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(
		GetOpenAPIDefinitions, openapi.NewDefinitionNamer(scheme.Scheme))
	config.GenericConfig.OpenAPIConfig.Info.Title = "Karbour"
	config.GenericConfig.OpenAPIConfig.Info.Version = "0.1"
	config.GenericConfig.LongRunningFunc = filters.BasicLongRunningRequestCheck(
		sets.NewString("watch", "proxy"),
		sets.NewString("attach", "exec", "proxy", "log", "portforward"),
	)
	if utilfeature.DefaultFeatureGate.Enabled(features.OpenAPIV3) {
		config.GenericConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(
			GetOpenAPIDefinitions, openapi.NewDefinitionNamer(scheme.Scheme))
		config.GenericConfig.OpenAPIV3Config.Info.Title = "Karbour"
		config.GenericConfig.OpenAPIV3Config.Info.Version = "0.1"
	}
	config.GenericConfig.BuildHandlerChainFunc = func(handler http.Handler, c *genericapiserver.Config) http.Handler {
		handler = genericapiserver.DefaultBuildHandlerChain(handler, c)
		handler = proxyutil.WithProxyByCluster(handler)
		return handler
	}

	config.GenericConfig.Config.EnableIndex = false
	// Define the discovery addresses for the API server
	config.GenericConfig.DiscoveryAddresses = discovery.DefaultAddresses{
		DefaultAddress: config.GenericConfig.LoopbackClientConfig.Host,
	}

	return config, nil
}

// RunServer starts a new APIServer given Options
func (o *Options) RunServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	server.GenericAPIServer.AddPostStartHookOrDie("start-server-informers", func(context genericapiserver.PostStartHookContext) error {
		config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
		return nil
	})

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	ret := k8sopenapi.GetOpenAPIDefinitions(ref)
	for k, v := range karbouropenapi.GetOpenAPIDefinitions(ref) {
		ret[k] = v
	}
	return ret
}
