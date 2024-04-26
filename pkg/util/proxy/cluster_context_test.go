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

package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestWithCluster tests the WithCluster function.
func TestWithCluster(t *testing.T) {
	tests := []struct {
		name     string
		parent   context.Context
		cluster  string
		expected string
	}{
		{
			name:     "With valid cluster context",
			parent:   context.Background(),
			cluster:  "test-cluster",
			expected: "test-cluster",
		},
		{
			name:     "With empty cluster context",
			parent:   context.Background(),
			cluster:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			ctx := WithCluster(tt.parent, tt.cluster)
			cluster, ok := ClusterFrom(ctx)

			// Assert
			require.Equal(t, tt.expected, cluster)
			require.True(t, ok)
		})
	}
}

// TestClusterFrom tests the ClusterFrom function.
func TestClusterFrom(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
		ok       bool
	}{
		{
			name:     "With valid cluster context",
			ctx:      WithCluster(context.Background(), "test-cluster"),
			expected: "test-cluster",
			ok:       true,
		},
		{
			name:     "Without cluster context",
			ctx:      context.Background(),
			expected: "",
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			cluster, ok := ClusterFrom(tt.ctx)

			// Assert
			require.Equal(t, tt.expected, cluster)
			require.Equal(t, tt.ok, ok)
		})
	}
}

// TestWithProxyByCluster tests the WithProxyByCluster middleware.
func TestWithProxyByCluster(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		expectedPath string
	}{
		{
			name:         "With valid cluster context",
			ctx:          WithCluster(context.Background(), "test-cluster"),
			expectedPath: "/apis/cluster.karpor.com/v1beta1/clusters/test-cluster/proxy/test-path",
		},
		{
			name:         "Without cluster context",
			ctx:          context.Background(),
			expectedPath: "/test-path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				// Assert
				require.Equal(t, tt.expectedPath, req.URL.Path)
			})

			// Execute
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
			req = req.WithContext(tt.ctx)
			WithProxyByCluster(handler).ServeHTTP(rr, req)
		})
	}
}
