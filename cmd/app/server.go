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

package app

import (
	"fmt"
	"io"
	"net"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"code.alipay.com/ant-iac/karbour/pkg/apis/cluster/v1"
	"code.alipay.com/ant-iac/karbour/pkg/apiserver"
	informers "code.alipay.com/ant-iac/karbour/pkg/generated/informers/externalversions"
	sampleopenapi "code.alipay.com/ant-iac/karbour/pkg/generated/openapi"
	netutils "k8s.io/utils/net"
)

const defaultEtcdPathPrefix = "/registry/karbour"

// Options contains state for master/api server
type Options struct {
	RecommendedOptions *genericoptions.RecommendedOptions

	SharedInformerFactory informers.SharedInformerFactory
	StdOut                io.Writer
	StdErr                io.Writer

	AlternateDNS []string
}

// NewOptions returns a new Options
func NewOptions(out, errOut io.Writer) *Options {
	o := &Options{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			apiserver.Codecs.LegacyCodec(v1.SchemeGroupVersion),
		),

		StdOut: out,
		StdErr: errOut,
	}
	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(v1.SchemeGroupVersion, schema.GroupKind{Group: v1.GroupName})
	return o
}

// NewApiserverCommand provides a CLI handler for 'start master' command
// with a default Options.
func NewApiserverCommand(defaults *Options, stopCh <-chan struct{}) *cobra.Command {
	o := *defaults
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
			if err := o.RunServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)
	utilfeature.DefaultMutableFeatureGate.AddFlag(flags)

	return cmd
}

// Validate validates Options
func (o *Options) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in fields required to have valid data
func (o *Options) Complete() error {
	return nil
}

// Config returns config for the api server given Options
func (o *Options) Config() (*apiserver.Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", o.AlternateDNS, []net.IP{netutils.ParseIPSloppy("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)

	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(sampleopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(apiserver.Scheme))
	serverConfig.OpenAPIConfig.Info.Title = "Karbour"
	serverConfig.OpenAPIConfig.Info.Version = "0.1"

	if utilfeature.DefaultFeatureGate.Enabled(features.OpenAPIV3) {
		serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(sampleopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(apiserver.Scheme))
		serverConfig.OpenAPIV3Config.Info.Title = "Karbour"
		serverConfig.OpenAPIV3Config.Info.Version = "0.1"
	}

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   apiserver.ExtraConfig{},
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

	// server.GenericAPIServer.AddPostStartHookOrDie("start-server-informers", func(context genericapiserver.PostStartHookContext) error {
	// 	config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
	// 	o.SharedInformerFactory.Start(context.StopCh)
	// 	return nil
	// })

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
