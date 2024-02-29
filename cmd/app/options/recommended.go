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

package options

import (
	"github.com/spf13/pflag"

	karbouropenapi "github.com/KusionStack/karbour/pkg/kubernetes/generated/openapi"
	k8sopenapi "github.com/KusionStack/karbour/pkg/kubernetes/openapi"
	"github.com/KusionStack/karbour/pkg/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/server"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/egressselector"
	"k8s.io/apiserver/pkg/server/filters"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/apiserver/pkg/util/feature"
	clientgoinformers "k8s.io/client-go/informers"
	clientgoclientset "k8s.io/client-go/kubernetes"
	"k8s.io/component-base/featuregate"
	"k8s.io/kube-openapi/pkg/common"
	kubeoptions "k8s.io/kubernetes/pkg/kubeapiserver/options"
)

// RecommendedOptions contains the recommended options for running an API server.
// If you add something to this list, it should be in a logical grouping.
// Each of them can be nil to leave the feature unconfigured on ApplyTo.
type RecommendedOptions struct {
	ServerRun      *options.ServerRunOptions
	Etcd           *options.EtcdOptions
	SecureServing  *options.SecureServingOptionsWithLoopback
	Audit          *options.AuditOptions
	Features       *options.FeatureOptions
	Authentication *kubeoptions.BuiltInAuthenticationOptions
	Authorization  *kubeoptions.BuiltInAuthorizationOptions

	// FeatureGate is a way to plumb feature gate through if you have them.
	FeatureGate featuregate.FeatureGate
	// ExtraAdmissionInitializers is called once after all ApplyTo from the options above, to pass the returned
	// admission plugin initializers to Admission.ApplyTo.
	ExtraAdmissionInitializers func(c *server.RecommendedConfig) ([]admission.PluginInitializer, error)
	Admission                  *options.AdmissionOptions
	// API Server Egress Selector is used to control outbound traffic from the API Server
	EgressSelector *options.EgressSelectorOptions
	// Traces contains options to control distributed request tracing.
	Traces *options.TracingOptions
}

func NewRecommendedOptions(prefix string, codec runtime.Codec) *RecommendedOptions {
	sso := options.NewSecureServingOptions()
	sso.HTTP2MaxStreamsPerConnection = 1000
	feature.DefaultMutableFeatureGate.Set("feature-gates=APIPriorityAndFairness=false")

	return &RecommendedOptions{
		ServerRun:     options.NewServerRunOptions(),
		Etcd:          options.NewEtcdOptions(storagebackend.NewDefaultConfig(prefix, codec)),
		SecureServing: sso.WithLoopback(),
		Authentication: kubeoptions.NewBuiltInAuthenticationOptions().
			WithAnonymous().
			WithClientCert().
			WithRequestHeader(),
		Authorization:              kubeoptions.NewBuiltInAuthorizationOptions(),
		Audit:                      options.NewAuditOptions(),
		Features:                   options.NewFeatureOptions(),
		FeatureGate:                feature.DefaultMutableFeatureGate,
		ExtraAdmissionInitializers: func(c *server.RecommendedConfig) ([]admission.PluginInitializer, error) { return nil, nil },
		Admission:                  options.NewAdmissionOptions(),
		EgressSelector:             options.NewEgressSelectorOptions(),
		Traces:                     options.NewTracingOptions(),
	}
}

func (o *RecommendedOptions) AddFlags(fs *pflag.FlagSet) {
	o.ServerRun.AddUniversalFlags(fs)
	o.Etcd.AddFlags(fs)
	o.SecureServing.AddFlags(fs)
	o.Authentication.AddFlags(fs)
	o.Authorization.AddFlags(fs)
	o.Audit.AddFlags(fs)
	o.Features.AddFlags(fs)
	o.Admission.AddFlags(fs)
	o.EgressSelector.AddFlags(fs)
	o.Traces.AddFlags(fs)
}

// ApplyTo adds RecommendedOptions to the server configuration.
// pluginInitializers can be empty, it is only need for additional initializers.
func (o *RecommendedOptions) ApplyTo(config *server.RecommendedConfig) error {
	genericConfig := &config.Config
	if err := o.ServerRun.ApplyTo(genericConfig); err != nil {
		return err
	}
	if err := o.Etcd.Complete(genericConfig.StorageObjectCountTracker, genericConfig.DrainedNotify(), genericConfig.AddPostStartHook); err != nil {
		return err
	}
	if err := o.Etcd.ApplyTo(genericConfig); err != nil {
		return err
	}
	if err := o.EgressSelector.ApplyTo(genericConfig); err != nil {
		return err
	}
	if err := o.Traces.ApplyTo(genericConfig.EgressSelector, genericConfig); err != nil {
		return err
	}
	if err := o.SecureServing.ApplyTo(&genericConfig.SecureServing, &genericConfig.LoopbackClientConfig); err != nil {
		return err
	}

	genericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(GetOpenAPIDefinitions, openapi.NewDefinitionNamer(scheme.Scheme))
	genericConfig.OpenAPIConfig.Info.Title = "Karbour"
	genericConfig.OpenAPIConfig.Info.Version = "0.1"

	genericConfig.LongRunningFunc = filters.BasicLongRunningRequestCheck(
		sets.NewString("watch", "proxy"),
		sets.NewString("attach", "exec", "proxy", "log", "portforward"),
	)

	kubeClientConfig := config.LoopbackClientConfig
	client, err := clientgoclientset.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}
	informer := clientgoinformers.NewSharedInformerFactory(client, kubeClientConfig.Timeout)
	config.ClientConfig = kubeClientConfig
	config.SharedInformerFactory = informer

	if err := o.Authentication.ApplyTo(&genericConfig.Authentication, genericConfig.SecureServing, config.EgressSelector, config.OpenAPIConfig, config.OpenAPIV3Config, client, informer); err != nil {
		return err
	}

	genericConfig.Authorization.Authorizer, genericConfig.RuleResolver, err = BuildAuthorizer(o.Authorization, config.EgressSelector, informer)
	if err != nil {
		return err
	}

	if err := o.Audit.ApplyTo(genericConfig); err != nil {
		return err
	}
	if err := o.Features.ApplyTo(genericConfig); err != nil {
		return err
	}

	if initializers, err := o.ExtraAdmissionInitializers(config); err != nil {
		return err
	} else if err := o.Admission.ApplyTo(genericConfig, config.SharedInformerFactory, config.ClientConfig, o.FeatureGate, initializers...); err != nil {
		return err
	}
	return nil
}

func (o *RecommendedOptions) Validate() []error {
	errors := []error{}
	errors = append(errors, o.ServerRun.Validate()...)
	errors = append(errors, o.Etcd.Validate()...)
	errors = append(errors, o.SecureServing.Validate()...)
	errors = append(errors, o.Authentication.Validate()...)
	errors = append(errors, o.Audit.Validate()...)
	errors = append(errors, o.Features.Validate()...)
	errors = append(errors, o.Admission.Validate()...)
	errors = append(errors, o.EgressSelector.Validate()...)
	errors = append(errors, o.Traces.Validate()...)

	return errors
}

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	ret := k8sopenapi.GetOpenAPIDefinitions(ref)
	for k, v := range karbouropenapi.GetOpenAPIDefinitions(ref) {
		ret[k] = v
	}
	return ret
}

// BuildAuthorizer constructs the authorizer
func BuildAuthorizer(authz *kubeoptions.BuiltInAuthorizationOptions, EgressSelector *egressselector.EgressSelector, versionedInformers clientgoinformers.SharedInformerFactory) (authorizer.Authorizer, authorizer.RuleResolver, error) {
	authorizationConfig := authz.ToAuthorizationConfig(versionedInformers)

	if EgressSelector != nil {
		egressDialer, err := EgressSelector.Lookup(egressselector.ControlPlane.AsNetworkContext())
		if err != nil {
			return nil, nil, err
		}
		authorizationConfig.CustomDial = egressDialer
	}

	return authorizationConfig.New()
}
