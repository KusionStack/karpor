// Copyright The Karpor Authors.
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
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/spf13/pflag"
)

type AIOptions struct {
	AIBackend     string
	AIAuthToken   string
	AIBaseURL     string
	AIModel       string
	AITemperature float32
	AITopP        float32
	// proxy options
	AIProxyEnabled bool
	AIHTTPProxy    string
	AIHTTPSProxy   string
	AINoProxy      string
}

const (
	defaultBackend     = "openai"
	defaultModel       = "gpt-3.5-turbo"
	defaultTemperature = 1
	defaultTopP        = 1
)

func NewAIOptions() *AIOptions {
	return &AIOptions{}
}

func (o *AIOptions) Validate() []error {
	return nil
}

func (o *AIOptions) ApplyTo(config *registry.ExtraConfig) error {
	// Apply the AIOptions to the provided config
	config.AIBackend = o.AIBackend
	config.AIAuthToken = o.AIAuthToken
	config.AIBaseURL = o.AIBaseURL
	if o.AIBackend != defaultBackend && o.AIModel == defaultModel {
		config.AIModel = ""
	} else {
		config.AIModel = o.AIModel
	}
	config.AITemperature = o.AITemperature
	config.AITopP = o.AITopP
	config.AIProxyEnabled = o.AIProxyEnabled
	config.AIHTTPProxy = o.AIHTTPProxy
	config.AIHTTPSProxy = o.AIHTTPSProxy
	config.AINoProxy = o.AINoProxy
	return nil
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *AIOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.AIBackend, "ai-backend", defaultBackend, "The ai backend")
	fs.StringVar(&o.AIAuthToken, "ai-auth-token", "", "The ai auth token")
	fs.StringVar(&o.AIBaseURL, "ai-base-url", "", "The ai base url")
	fs.StringVar(&o.AIModel, "ai-model", defaultModel, "The ai model")
	fs.Float32Var(&o.AITemperature, "ai-temperature", defaultTemperature, "The ai temperature")
	fs.Float32Var(&o.AITopP, "ai-top-p", defaultTopP, "The ai top-p")
	fs.BoolVar(&o.AIProxyEnabled, "ai-proxy-enabled", false, "The ai proxy enable")
	fs.StringVar(&o.AIHTTPProxy, "ai-http-proxy", "", "The ai http proxy")
	fs.StringVar(&o.AIHTTPSProxy, "ai-https-proxy", "", "The ai https proxy")
	fs.StringVar(&o.AINoProxy, "ai-no-proxy", "", "The ai no-proxy")
}
