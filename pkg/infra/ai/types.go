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

package ai

import (
	"context"

	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
)

const (
	AzureProvider       = "azureopenai"
	HuggingFaceProvider = "huggingface"
	OpenAIProvider      = "openai"
)

const (
	Text2sqlType = "Text2sql"
	SQLFixType   = "SqlFix"
)

var clients = map[string]AIProvider{
	AzureProvider:       &AzureAIClient{},
	HuggingFaceProvider: &HuggingfaceClient{},
	OpenAIProvider:      &OpenAIClient{},
}

// AIProvider is an interface all AI clients.
type AIProvider interface {
	// Configure sets up the AI service with the provided configuration.
	Configure(config AIConfig) error
	// Generate generates a response from the AI service based on
	// the provided prompt and service type.
	Generate(ctx context.Context, prompt string) (string, error)
	// GenerateStream generates a streaming response from the AI service
	// based on the provided prompt. It returns a channel that will receive
	// chunks of the response as they are generated.
	GenerateStream(ctx context.Context, prompt string) (<-chan string, error)
}

// AIConfig represents the configuration settings for an AI client.
type AIConfig struct {
	Name         string
	AuthToken    string
	BaseURL      string
	Model        string
	Temperature  float32
	TopP         float32
	ProxyEnabled bool
	HTTPProxy    string
	HTTPSProxy   string
	NoProxy      string
}

func ConvertToAIConfig(c registry.ExtraConfig) AIConfig {
	return AIConfig{
		Name:         c.AIBackend,
		AuthToken:    c.AIAuthToken,
		BaseURL:      c.AIBaseURL,
		Model:        c.AIModel,
		Temperature:  c.AITemperature,
		TopP:         c.AITopP,
		ProxyEnabled: c.AIProxyEnabled,
		HTTPProxy:    c.AIHTTPProxy,
		HTTPSProxy:   c.AIHTTPSProxy,
		NoProxy:      c.AINoProxy,
	}
}

// NewClient returns a new AIProvider object
func NewClient(name string) AIProvider {
	if client, exists := clients[name]; exists {
		return client
	}
	// default client
	return &OpenAIClient{}
}
