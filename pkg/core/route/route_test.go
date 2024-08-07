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

package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry/search"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	"k8s.io/apiserver/pkg/server"
)

// mockSearchStorage is an in-memory implementation of the SearchStorage
// interface for testing purposes.
type mockSearchStorage struct {
	storage.SearchStorage
}

// mockResourceStorage is an in-memory implementation of the ResourceStorage
// interface for testing purposes.
type mockResourceStorage struct {
	storage.ResourceStorage
}

// mockResourceGroupRuleStorage is an in-memory implementation of the
// ResourceGroupRuleStorage interface for testing purposes.
type mockResourceGroupRuleStorage struct {
	storage.ResourceGroupRuleStorage
}

// mockGeneralStorage is an in-memory implementation of the Storage interface
// for testing purposes.
type mockGeneralStorage struct {
	storage.Storage
}

// TestNewCoreRoute will test the NewCoreRoute function with different
// configurations.
func TestNewCoreRoute(t *testing.T) {
	// Mock the NewSearchStorage and NewResourceGroupRuleStorage function to
	// return a mock storage instead of actual implementation.
	mockey.Mock(search.NewSearchStorage).Return(&mockSearchStorage{}, nil).Build()
	mockey.Mock(search.NewResourceStorage).Return(&mockResourceStorage{}, nil).Build()
	mockey.Mock(search.NewResourceGroupRuleStorage).Return(&mockResourceGroupRuleStorage{}, nil).Build()
	mockey.Mock(search.NewGeneralStorage).Return(&mockGeneralStorage{}, nil).Build()
	defer mockey.UnPatchAll()

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
				"/rest-api/v1/search/", // fixme: this may result in a nil pointer
				"/livez",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, err := NewCoreRoute(&server.CompletedConfig{}, &tt.extraConfig)
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
