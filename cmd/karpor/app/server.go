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
	"time"

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

	"github.com/KusionStack/karpor/cmd/karpor/app/options"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"github.com/KusionStack/karpor/pkg/server"
	proxyutil "github.com/KusionStack/karpor/pkg/util/proxy"
	"github.com/KusionStack/karpor/pkg/version"
)

const (
	defaultEtcdPathPrefix     = "/registry/karpor"
	defaultTokenIssuer        = "karpor"
	defaultTokenMaxExpiration = 8760 * time.Hour
)

// Options contains state for master/api server
type Options struct {
	RecommendedOptions   *options.RecommendedOptions
	SearchStorageOptions *options.SearchStorageOptions
	CoreOptions          *options.CoreOptions
	AIOptions            *options.AIOptions

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
		AIOptions:            options.NewAIOptions(),
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
	o.RecommendedOptions.Admission.DisablePlugins = []string{"MutatingAdmissionWebhook", "NamespaceLifecycle", "ValidatingAdmissionWebhook"}
	o.RecommendedOptions.Authorization.Modes = []string{"RBAC"}
	o.RecommendedOptions.ServerRun.CorsAllowedOriginList = []string{".*"}
	return o, nil
}

// NewServerCommand provides a CLI handler for 'start master' command
// with a default Options.
func NewServerCommand(ctx context.Context) *cobra.Command {
	o, err := NewOptions(os.Stdout, os.Stderr)
	if err != nil {
		klog.Error(err, "Unable to initialize command options")
		klog.Flush()
	}

	expvar.Publish("CoreOptions", expvar.Func(func() interface{} {
		return o.CoreOptions
	}))
	expvar.Publish("StorageOptions", expvar.Func(func() interface{} {
		return o.SearchStorageOptions
	}))
	expvar.Publish("AIOptions", expvar.Func(func() interface{} {
		displayOpts := *o.AIOptions
		displayOpts.AIAuthToken = "[hidden]"
		return &displayOpts
	}))
	expvar.Publish("Version", expvar.Func(func() interface{} {
		return version.GetVersion()
	}))

	cmd := &cobra.Command{
		Use:   "karpor",
		Short: "Launch an API server",
		Long:  "Launch an API server",
		RunE: func(c *cobra.Command, args []string) error {
			if o.CoreOptions.Version {
				fmt.Println(version.GetVersion())
				return nil
			}
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
	o.AIOptions.AddFlags(fs)
}

// Validate validates Options
func (o *Options) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	errors = append(errors, o.SearchStorageOptions.Validate()...)
	errors = append(errors, o.AIOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in fields required to have valid data
func (o *Options) Complete() error {
	// generate token issuer
	if len(o.RecommendedOptions.Authentication.ServiceAccounts.Issuers) == 0 || o.RecommendedOptions.Authentication.ServiceAccounts.Issuers[0] == "" {
		o.RecommendedOptions.Authentication.ServiceAccounts.Issuers = []string{defaultTokenIssuer}
	}

	// set default token max expiration
	o.RecommendedOptions.ServiceAccountTokenMaxExpiration = defaultTokenMaxExpiration
	if o.RecommendedOptions.Authentication.ServiceAccounts.MaxExpiration != 0 {
		o.RecommendedOptions.ServiceAccountTokenMaxExpiration = o.RecommendedOptions.Authentication.ServiceAccounts.MaxExpiration
	}

	// complete two content-related keys with each other
	if o.RecommendedOptions.ServiceAccountSigningKeyFile == "" && (len(o.RecommendedOptions.Authentication.ServiceAccounts.KeyFiles) == 0 ||
		o.RecommendedOptions.Authentication.ServiceAccounts.KeyFiles[0] == "") {
		return fmt.Errorf("no valid serviceaccounts signing key file")
	}
	if o.RecommendedOptions.ServiceAccountSigningKeyFile == "" {
		o.RecommendedOptions.ServiceAccountSigningKeyFile = o.RecommendedOptions.Authentication.ServiceAccounts.KeyFiles[0]
	}
	if len(o.RecommendedOptions.Authentication.ServiceAccounts.KeyFiles) == 0 {
		o.RecommendedOptions.Authentication.ServiceAccounts.KeyFiles = []string{o.RecommendedOptions.ServiceAccountSigningKeyFile}
	}

	// create token generator
	sk, err := keyutil.PrivateKeyFromFile(o.RecommendedOptions.ServiceAccountSigningKeyFile)
	if err != nil {
		return fmt.Errorf("failed to parse key-file for token generator: %w", err)
	}
	o.RecommendedOptions.ServiceAccountIssuer, err = serviceaccount.JWTTokenGenerator(o.RecommendedOptions.Authentication.ServiceAccounts.Issuers[0], sk)
	if err != nil {
		return fmt.Errorf("create token generator failed: %w", err)
	}
	return nil
}

// Config returns config for the api server given Options
func (o *Options) Config() (*server.Config, error) {
	config := &server.Config{
		GenericConfig: genericapiserver.NewRecommendedConfig(scheme.Codecs),
		ExtraConfig:   &registry.ExtraConfig{},
	}
	if !o.CoreOptions.EnableRBAC {
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
	if err := o.AIOptions.ApplyTo(config.ExtraConfig); err != nil {
		return nil, err
	}

	config.GenericConfig.BuildHandlerChainFunc = func(handler http.Handler, c *genericapiserver.Config) http.Handler {
		handler = genericapiserver.DefaultBuildHandlerChain(handler, c)
		handler = proxyutil.WithProxyByCluster(handler)
		if !o.CoreOptions.EnableRBAC {
			// remove http header "Authorization" if RBAC is not enabled
			handler = removeAuthorizationHeader(handler)
		}
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

func removeAuthorizationHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newReq := r.WithContext(context.Background())
		newReq.Header.Del("Authorization")
		next.ServeHTTP(w, newReq)
	})
}
