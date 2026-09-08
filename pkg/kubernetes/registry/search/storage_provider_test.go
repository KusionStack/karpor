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

package search

import (
	"path/filepath"
	"testing"

	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	sqlitestorage "github.com/KusionStack/karpor/pkg/infra/search/storage/sqlite"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/stretchr/testify/require"
)

func TestRESTStorageProviderStorageGetters(t *testing.T) {
	tests := []struct {
		name        string
		provider    RESTStorageProvider
		assertTypes func(t *testing.T, p RESTStorageProvider)
		wantErr     string
	}{
		{
			name: "elasticsearch",
			provider: RESTStorageProvider{
				SearchStorageType:      "elasticsearch",
				ElasticSearchAddresses: []string{"http://localhost:9200"},
			},
			assertTypes: func(t *testing.T, p RESTStorageProvider) {
				t.Helper()

				searchGetter, err := p.SearchStorageGetter()
				require.NoError(t, err)
				require.IsType(t, &elasticsearch.SearchStorageGetter{}, searchGetter)

				resourceGetter, err := p.ResourceStorageGetter()
				require.NoError(t, err)
				require.IsType(t, &elasticsearch.ResourceStorageGetter{}, resourceGetter)

				ruleGetter, err := p.ResourceGroupRuleStorageGetter()
				require.NoError(t, err)
				require.IsType(t, &elasticsearch.ResourceGroupRuleStorageGetter{}, ruleGetter)

				generalGetter, err := p.GeneralStorageGetter()
				require.NoError(t, err)
				require.IsType(t, &elasticsearch.GeneralStorageGetter{}, generalGetter)
			},
		},
		{
			name: "sqlite",
			provider: RESTStorageProvider{
				SearchStorageType: "sqlite",
				SQLiteDSN:         filepath.Join(t.TempDir(), "karpor-search.db"),
			},
			assertTypes: func(t *testing.T, p RESTStorageProvider) {
				t.Helper()

				searchGetter, err := p.SearchStorageGetter()
				require.NoError(t, err)
				require.IsType(t, &sqlitestorage.SearchStorageGetter{}, searchGetter)

				resourceGetter, err := p.ResourceStorageGetter()
				require.NoError(t, err)
				require.IsType(t, &sqlitestorage.ResourceStorageGetter{}, resourceGetter)

				ruleGetter, err := p.ResourceGroupRuleStorageGetter()
				require.NoError(t, err)
				require.IsType(t, &sqlitestorage.ResourceGroupRuleStorageGetter{}, ruleGetter)

				generalGetter, err := p.GeneralStorageGetter()
				require.NoError(t, err)
				require.IsType(t, &sqlitestorage.GeneralStorageGetter{}, generalGetter)
			},
		},
		{
			name:     "invalid",
			provider: RESTStorageProvider{SearchStorageType: "zincsearch"},
			wantErr:  "invalid search storage type zincsearch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != "" {
				_, err := tt.provider.SearchStorageGetter()
				require.EqualError(t, err, tt.wantErr)
				_, err = tt.provider.ResourceStorageGetter()
				require.EqualError(t, err, "invalid resource storage type zincsearch")
				_, err = tt.provider.ResourceGroupRuleStorageGetter()
				require.EqualError(t, err, "invalid resource group rule storage type zincsearch")
				_, err = tt.provider.GeneralStorageGetter()
				require.EqualError(t, err, "invalid general storage type zincsearch")
				return
			}
			tt.assertTypes(t, tt.provider)
		})
	}
}

func TestNewStorageFactoriesWithSQLite(t *testing.T) {
	cfg := registry.ExtraConfig{
		SearchStorageType: "sqlite",
		SQLiteDSN:         filepath.Join(t.TempDir(), "karpor-search.db"),
	}

	searchStorage, err := NewSearchStorage(cfg)
	require.NoError(t, err)
	require.IsType(t, &sqlitestorage.Storage{}, searchStorage)
	require.NoError(t, searchStorage.(*sqlitestorage.Storage).Close())

	resourceStorage, err := NewResourceStorage(cfg)
	require.NoError(t, err)
	require.IsType(t, &sqlitestorage.Storage{}, resourceStorage)
	require.NoError(t, resourceStorage.(*sqlitestorage.Storage).Close())

	ruleStorage, err := NewResourceGroupRuleStorage(cfg)
	require.NoError(t, err)
	require.IsType(t, &sqlitestorage.Storage{}, ruleStorage)
	require.NoError(t, ruleStorage.(*sqlitestorage.Storage).Close())

	generalStorage, err := NewGeneralStorage(cfg)
	require.NoError(t, err)
	require.IsType(t, &sqlitestorage.Storage{}, generalStorage)
	require.NoError(t, generalStorage.(*sqlitestorage.Storage).Close())
}
