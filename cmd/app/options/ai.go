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
	Backend     string
	AuthToken   string
	BaseURL     string
	Model       string
	Temperature float32
	TopP        float32
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
	config.Backend = o.Backend
	config.AuthToken = o.AuthToken
	config.BaseURL = o.BaseURL
	config.Model = o.Model
	config.Temperature = o.Temperature
	config.TopP = o.TopP
	return nil
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *AIOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.Backend, "ai-backend", defaultBackend, "The ai backend")
	fs.StringVar(&o.AuthToken, "ai-auth-token", "", "The ai auth token")
	fs.StringVar(&o.BaseURL, "ai-base-url", "", "The ai base url")
	fs.StringVar(&o.Model, "ai-model", defaultModel, "The ai model")
	fs.Float32Var(&o.Temperature, "ai-temperature", defaultTemperature, "The ai temperature")
	fs.Float32Var(&o.TopP, "ai-top-p", defaultTopP, "The ai top-p")
}
