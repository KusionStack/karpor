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
	"net/http"

	"github.com/hupe1980/go-huggingface"
)

type HuggingfaceClient struct {
	client      *huggingface.InferenceClient
	model       string
	temperature float32
}

func (c *HuggingfaceClient) Configure(cfg AIConfig) error {
	if cfg.ProxyEnabled {
		c.client = huggingface.NewInferenceClient(cfg.AuthToken, func(o *huggingface.InferenceClientOptions) {
			o.HTTPClient = &http.Client{
				Transport: GetProxyHTTPClient(cfg),
			}
		})
	} else {
		c.client = huggingface.NewInferenceClient(cfg.AuthToken)
	}

	c.model = cfg.Model
	c.temperature = cfg.Temperature
	return nil
}

func (c *HuggingfaceClient) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.TextGeneration(ctx, &huggingface.TextGenerationRequest{
		Inputs: prompt,
		Parameters: huggingface.TextGenerationParameters{
			Temperature: huggingface.PTR(float64(c.temperature)),
		},
		Options: huggingface.Options{
			WaitForModel: huggingface.PTR(true),
		},
		Model: c.model,
	})
	if err != nil {
		return "", err
	}

	return resp[0].GeneratedText[len(prompt):], nil
}

func (c *HuggingfaceClient) GenerateStream(ctx context.Context, prompt string) (<-chan string, error) {
	// Hugging Face API doesn't support streaming yet, so we simulate streaming output
	resultChan := make(chan string, 100)

	go func() {
		defer close(resultChan)

		resp, err := c.Generate(ctx, prompt)
		if err != nil {
			resultChan <- "ERROR: " + err.Error()
			return
		}

		// Simulate streaming by sending response character by character
		// Note: This is not ideal for production use
		// Consider waiting for native streaming support from Hugging Face API
		for _, char := range resp {
			select {
			case <-ctx.Done():
				return
			case resultChan <- string(char):
			}
		}
	}()

	return resultChan, nil
}
