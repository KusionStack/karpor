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
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hupe1980/go-huggingface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// newHuggingfaceClientWithBody returns a client whose Inference API answers every request with body.
func newHuggingfaceClientWithBody(body string) *HuggingfaceClient {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}
	return &HuggingfaceClient{
		client: huggingface.NewInferenceClient("test-token", func(o *huggingface.InferenceClientOptions) {
			o.HTTPClient = httpClient
		}),
		model: "test-model",
	}
}

func TestHuggingfaceClient_Generate(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{name: "response echoes the prompt", body: `[{"generated_text":"Say hi. Hi there"}]`, want: " Hi there"},
		{name: "response without the prompt is kept whole", body: `[{"generated_text":"Hi"}]`, want: "Hi"},
		{name: "empty result list is an error", body: `[]`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newHuggingfaceClientWithBody(tt.body)

			var got string
			var err error
			require.NotPanics(t, func() { got, err = c.Generate(context.Background(), "Say hi.") })

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
