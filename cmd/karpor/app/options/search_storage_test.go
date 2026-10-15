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

package options

import (
	"encoding/json"
	"testing"

	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestSearchStorageOptionsApplyToIncludesSQLiteDSN(t *testing.T) {
	opts := &SearchStorageOptions{
		SearchStorageType:      "sqlite",
		ElasticSearchAddresses: []string{"http://localhost:9200"},
		ElasticSearchUsername:  "elastic",
		ElasticSearchPassword:  "secret",
		SQLiteDSN:              "file:/tmp/karpor-search.db?mode=rwc",
	}
	cfg := &registry.ExtraConfig{}

	require.NoError(t, opts.ApplyTo(cfg))
	require.Equal(t, "sqlite", cfg.SearchStorageType)
	require.Equal(t, []string{"http://localhost:9200"}, cfg.ElasticSearchAddresses)
	require.Equal(t, "elastic", cfg.ElasticSearchUsername)
	require.Equal(t, "secret", cfg.ElasticSearchPassword)
	require.Equal(t, "file:/tmp/karpor-search.db?mode=rwc", cfg.SQLiteDSN)
}

func TestSearchStorageOptionsAddFlagsIncludesSQLiteDSN(t *testing.T) {
	opts := &SearchStorageOptions{}
	fs := pflag.NewFlagSet("search-storage", pflag.ContinueOnError)
	opts.AddFlags(fs)

	require.NoError(t, fs.Parse([]string{
		"--search-storage-type=sqlite",
		"--sqlite-dsn=file:karpor-test.db?mode=rwc",
	}))

	require.Equal(t, "sqlite", opts.SearchStorageType)
	require.Equal(t, "file:karpor-test.db?mode=rwc", opts.SQLiteDSN)
	require.NotNil(t, fs.Lookup("sqlite-dsn"))
}

func TestSearchStorageOptionsMarshalJSONMasksElasticPassword(t *testing.T) {
	opts := SearchStorageOptions{
		SearchStorageType:     "sqlite",
		ElasticSearchPassword: "secret",
		SQLiteDSN:             "file:karpor-search.db?mode=rwc",
	}

	data, err := json.Marshal(opts)
	require.NoError(t, err)

	require.Contains(t, string(data), `"ElasticSearchPassword":"`+MaskString+`"`)
	require.Contains(t, string(data), `"SQLiteDSN":"file:karpor-search.db?mode=rwc"`)
	require.NotContains(t, string(data), "secret")
}
