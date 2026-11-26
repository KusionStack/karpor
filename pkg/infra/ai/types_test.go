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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProxyHTTPClient_NoProxyIgnoresPort(t *testing.T) {
	cfg := AIConfig{
		HTTPProxy: "http://proxy.example.com:3128",
		NoProxy:   "localhost,.internal.example.com",
	}
	transport := GetProxyHTTPClient(cfg)

	tests := []struct {
		name      string
		url       string
		wantProxy bool
	}{
		{"host with port is excluded", "http://localhost:11434/v1/chat", false},
		{"host without port is excluded", "http://localhost/v1/chat", false},
		{"suffix match with port is excluded", "http://llm.internal.example.com:8080/v1", false},
		{"other host with port is proxied", "http://api.example.org:11434/v1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			assert.NoError(t, err)
			proxyURL, err := transport.Proxy(req)
			assert.NoError(t, err)
			if tt.wantProxy {
				assert.NotNil(t, proxyURL)
			} else {
				assert.Nil(t, proxyURL)
			}
		})
	}
}

func TestGetProxyHTTPClient_NoProxyEntryWithPort(t *testing.T) {
	cfg := AIConfig{
		HTTPProxy: "http://proxy.example.com:3128",
		NoProxy:   "localhost:11434",
	}
	transport := GetProxyHTTPClient(cfg)

	tests := []struct {
		name      string
		url       string
		wantProxy bool
	}{
		{"matching port is excluded", "http://localhost:11434/v1/chat", false},
		{"other port is proxied", "http://localhost:8080/v1/chat", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			assert.NoError(t, err)
			proxyURL, err := transport.Proxy(req)
			assert.NoError(t, err)
			if tt.wantProxy {
				assert.NotNil(t, proxyURL)
			} else {
				assert.Nil(t, proxyURL)
			}
		})
	}
}
