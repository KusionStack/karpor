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

package app

import (
	"context"
	"expvar"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/KusionStack/karpor/cmd/app/options"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"github.com/KusionStack/karpor/pkg/server"
	proxyutil "github.com/KusionStack/karpor/pkg/util/proxy"
	"github.com/KusionStack/karpor/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/discovery"
	"k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/util/keyutil"
	"k8s.io/klog/v2"
	authzmodes "k8s.io/kubernetes/pkg/kubeapiserver/authorizer/modes"
	"k8s.io/kubernetes/pkg/serviceaccount"
	netutils "k8s.io/utils/net"
)

const defaultEtcdPathPrefix = "/registry/karpor"

// Options contains state for master/api server
type Options struct {
	RecommendedOptions   *options.RecommendedOptions
	SearchStorageOptions *options.SearchStorageOptions
	CoreOptions          *options.CoreOptions

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
		CoreOptions:          options.NewCoreOptions(),
		StdOut:               out,
		StdErr:               errOut,
	}
	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = schema.GroupVersions(scheme.Versions)
	o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"localhost", o.AlternateDNS, []net.IP{netutils.ParseIPSloppy("127.0.0.1")},
	); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}
	o.RecommendedOptions.Admission.DisablePlugins = []string{"MutatingAdmissionWebhook", "NamespaceLifecycle", "ValidatingAdmissionWebhook", "ValidatingAdmissionPolicy"}
	o.RecommendedOptions.Authorization.Modes = []string{"RBAC"}
	o.RecommendedOptions.ServerRun.CorsAllowedOriginList = []string{".*"}
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

	expvar.Publish("CoreOptions", expvar.Func(func() interface{} {
		return o.CoreOptions
	}))
	expvar.Publish("StorageOptions", expvar.Func(func() interface{} {
		return o.SearchStorageOptions
	}))
	expvar.Publish("Version", expvar.Func(func() interface{} {
		return version.GetVersion()
	}))

	cmd := &cobra.Command{
		Use:   "karpor",
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
	o.CoreOptions.AddFlags(fs)
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
	// generate token issuer
	if len(o.RecommendedOptions.Authentication.ServiceAccounts.Issuers) == 0 || o.RecommendedOptions.Authentication.ServiceAccounts.Issuers[0] == "" {
		return fmt.Errorf("no valid serviceaccounts issuer")
	}
	if o.RecommendedOptions.ServiceAccountSigningKeyFile == "" {
		return fmt.Errorf("no valid serviceaccounts signing key file")
	}
	sk, err := keyutil.PrivateKeyFromFile(o.RecommendedOptions.ServiceAccountSigningKeyFile)
	if err != nil {
		return fmt.Errorf("failed to parse service-account-issuer-key-file: %w", err)
	}
	o.RecommendedOptions.ServiceAccountIssuer, err = serviceaccount.JWTTokenGenerator(o.RecommendedOptions.Authentication.ServiceAccounts.Issuers[0], sk)
	if err != nil {
		return fmt.Errorf("no valid serviceaccounts signing key file: %w", err)
	}
	o.RecommendedOptions.ServiceAccountTokenMaxExpiration = o.RecommendedOptions.Authentication.ServiceAccounts.MaxExpiration

	return nil
}

// Config returns config for the api server given Options
func (o *Options) Config() (*server.Config, error) {
	config := &server.Config{
		GenericConfig: genericapiserver.NewRecommendedConfig(scheme.Codecs),
		ExtraConfig:   &registry.ExtraConfig{},
	}
	// always allow access if readOnlyMode is open
	if o.CoreOptions.ReadOnlyMode {
		o.RecommendedOptions.Authorization.Modes = []string{authzmodes.ModeAlwaysAllow}
	}
	if err := o.RecommendedOptions.ApplyTo(config.GenericConfig); err != nil {
		return nil, err
	}
	if err := o.RecommendedOptions.ApplyToExtraConfig(config.ExtraConfig); err != nil {
		return nil, err
	}
	if err := o.SearchStorageOptions.ApplyTo(config.ExtraConfig); err != nil {
		return nil, err
	}
	if err := o.CoreOptions.ApplyTo(config.ExtraConfig); err != nil {
		return nil, err
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

	serv, err := config.Complete().New()
	if err != nil {
		return err
	}

	serv.GenericAPIServer.AddPostStartHookOrDie("start-server-informers", func(context genericapiserver.PostStartHookContext) error {
		config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
		return nil
	})

	serv.GenericAPIServer.AddPostStartHookOrDie("register-default-config", server.ConfigRegister)

	return serv.GenericAPIServer.PrepareRun().Run(stopCh)
}
