package ai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepseekClient_Configure(t *testing.T) {
	tests := []struct {
		name    string
		config  AIConfig
		wantErr bool
	}{
		{
			name: "Base Config",
			config: AIConfig{
				AuthToken: "test-token",
				Model:     "deepseek-r1",
				BaseURL:   "",
			},
			wantErr: false,
		},
		{
			name: "Custom BaseURL",
			config: AIConfig{
				AuthToken: "test-token",
				Model:     "deepseek-r1",
				BaseURL:   "https://custom.deepseek.api",
			},
			wantErr: false,
		},
		{
			name: "Enable Proxy",
			config: AIConfig{
				AuthToken:    "test-token",
				Model:        "deepseek-r1",
				ProxyEnabled: true,
				HTTPProxy:    "http://proxy.example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DeepseekClient{}
			err := c.Configure(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.config.Model, c.model)
				assert.NotNil(t, c.client)
			}
		})
	}
}

func TestDeepseekClient_Generate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "POST", r.Method)

		// Mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"choices": [{
				"message": {
					"content": "test response"
				}
			}]
		}`))
	}))
	defer server.Close()

	client := &DeepseekClient{}
	err := client.Configure(AIConfig{
		AuthToken: "test-token",
		BaseURL:   server.URL,
		Model:     "deepseek-r1",
	})
	assert.NoError(t, err)

	// Test generate
	resp, err := client.Generate(context.Background(), "test prompt")
	assert.NoError(t, err)
	assert.Equal(t, "test response", resp)
}

func TestDeepseekClient_GenerateStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "POST", r.Method)

		// Set stream response
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Mock stream data
		w.Write([]byte(`data: {"choices":[{"delta":{"content":"first"}}]}
data: {"choices":[{"delta":{"content":"second"}}]}
data: {"choices":[{"delta":{"content":"third"}}]}
data: [DONE]
`))
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	client := &DeepseekClient{}
	err := client.Configure(AIConfig{
		AuthToken: "test-token",
		BaseURL:   server.URL,
		Model:     "deepseek-v3",
	})
	assert.NoError(t, err)

	// Test stream generate
	stream, err := client.GenerateStream(context.Background(), "test prompt")
	assert.NoError(t, err)

	// Collect stream response
	var result string
	for chunk := range stream {
		if len(chunk) >= 6 && chunk[:6] == "ERROR:" {
			t.Fatalf("received error: %s", chunk[6:])
		}
		result += chunk
	}

	assert.Equal(t, "firstsecondthird", result)
}
