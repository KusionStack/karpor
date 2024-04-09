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

package route

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/kubernetes/registry"
	"github.com/KusionStack/karbour/pkg/kubernetes/registry/search"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	"k8s.io/apiserver/pkg/server"
)

// mockSearchStorage is an in-memory implementation of the SearchStorage
// interface for testing purposes.
type mockSearchStorage struct{}

// Search implements the search operation returning a single mock resource.
func (m *mockSearchStorage) Search(ctx context.Context, queryString, patternType string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	return &storage.SearchResult{
		Total: 1,
		Resources: []*storage.Resource{{
			Locator: core.Locator{
				Cluster:   "mock-cluster",
				Namespace: "mock-namespace",
				Name:      "mock-name",
			},
			Object: map[string]interface{}{},
		}},
	}, nil
}

// TestNewCoreRoute will test the NewCoreRoute function with different
// configurations.
func TestNewCoreRoute(t *testing.T) {
	tests := []struct {
		name         string
		extraConfig  registry.ExtraConfig
		expectError  bool
		expectRoutes []string
	}{
		{
			name: "successful with normal mode",
			extraConfig: registry.ExtraConfig{
				SearchStorageType: "elasticsearch",
			}, // Normal mode.
			expectError: false,
			expectRoutes: []string{
				"/endpoints",
				"/server-configs",
				"/rest-api/v1/search/",
			},
		},
	}

	// Initialize dummy server config.
	genericConfig := &server.CompletedConfig{}

	// Mock the NewSearchStorage function to return a mock storage instead of
	// actual implementation.
	mockey.Mock(search.NewSearchStorage).Return(&mockSearchStorage{}, nil).Build()
	defer mockey.UnPatchAll()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, err := NewCoreRoute(genericConfig, &tt.extraConfig)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				for _, route := range tt.expectRoutes {
					req := httptest.NewRequest(http.MethodGet, route, nil)
					rr := httptest.NewRecorder()
					router.ServeHTTP(rr, req)

					// Assert status code is not 404 to ensure the route exists.
					require.NotEqual(t, http.StatusNotFound, rr.Code, "Route should exist: %s", route)
				}
			}
		})
	}
}
