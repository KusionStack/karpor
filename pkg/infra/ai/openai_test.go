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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const openAIChatReply = `{"id":"1","object":"chat.completion","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`

// requestBodyFor configures an OpenAIClient against a test server and returns the
// JSON body sent by Generate.
func requestBodyFor(t *testing.T, temperature, topP float32) map[string]any {
	t.Helper()

	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(raw, &body))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(openAIChatReply))
	}))
	t.Cleanup(server.Close)

	c := &OpenAIClient{}
	require.NoError(t, c.Configure(AIConfig{
		AuthToken:   "test-token",
		Model:       "gpt-5.6-terra",
		BaseURL:     server.URL,
		Temperature: temperature,
		TopP:        topP,
	}))

	_, err := c.Generate(context.Background(), "hello")
	require.NoError(t, err)
	return body
}

func TestOpenAIClient_Generate_OmitsDefaultTopP(t *testing.T) {
	body := requestBodyFor(t, 1, 1)

	assert.NotContains(t, body, "top_p")
	assert.EqualValues(t, 1, body["temperature"])
}

func TestOpenAIClient_Generate_SendsCustomTopP(t *testing.T) {
	body := requestBodyFor(t, 1, 0.5)

	assert.EqualValues(t, 0.5, body["top_p"])
}
