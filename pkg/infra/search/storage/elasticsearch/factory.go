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

package elasticsearch

import (
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/elastic/go-elasticsearch/v8"
)

var _ storage.SearchStorageGetter = &SearchStorageGetter{}

type SearchStorageGetter struct {
	cfg *Config
}

func (s *SearchStorageGetter) GetSearchStorage() (storage.SearchStorage, error) {
	esClient, err := NewESClient(elasticsearch.Config{
		Addresses: s.cfg.Addresses,
		Username:  s.cfg.UserName,
		Password:  s.cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

type Config struct {
	Addresses []string `env:"ES_ADDRESSES"`
	UserName  string   `env:"ES_USER"`
	Password  string   `env:"ES_PASSWORD"`
}

func NewSearchStorageGetter(addresses []string, userName, password string) *SearchStorageGetter {
	cfg := &Config{
		Addresses: addresses,
		UserName:  userName,
		Password:  password,
	}

	return &SearchStorageGetter{
		cfg,
	}
}
