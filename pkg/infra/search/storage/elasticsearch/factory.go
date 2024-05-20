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

package elasticsearch

import (
	"kusionstack.io/karpor/pkg/infra/search/storage"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	_ storage.SearchStorageGetter            = &SearchStorageGetter{}
	_ storage.ResourceStorageGetter          = &ResourceStorageGetter{}
	_ storage.ResourceGroupRuleStorageGetter = &ResourceGroupRuleStorageGetter{}
)

// SearchStorageGetter represents a structure for getting search storage instances.
type SearchStorageGetter struct {
	cfg *Config
}

// GetSearchStorage retrieves and returns a search storage instance based on the provided configuration.
func (s *SearchStorageGetter) GetSearchStorage() (storage.SearchStorage, error) {
	esClient, err := NewStorage(elasticsearch.Config{
		Addresses: s.cfg.Addresses,
		Username:  s.cfg.UserName,
		Password:  s.cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

// ResourceStorageGetter represents a structure for getting resource storage
// instances.
type ResourceStorageGetter struct {
	cfg *Config
}

// GetResourceStorage retrieves and returns a resource storage instance based on
// the provided configuration.
func (s *ResourceStorageGetter) GetResourceStorage() (storage.ResourceStorage, error) {
	esClient, err := NewStorage(elasticsearch.Config{
		Addresses: s.cfg.Addresses,
		Username:  s.cfg.UserName,
		Password:  s.cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

// ResourceGroupRuleStorageGetter represents a structure for getting resource
// group rule storage instances.
type ResourceGroupRuleStorageGetter struct {
	cfg *Config
}

// GetResourceGroupRuleStorage retrieves and returns a resource group rule
// storage instance based on the provided configuration.
func (s *ResourceGroupRuleStorageGetter) GetResourceGroupRuleStorage() (storage.ResourceGroupRuleStorage, error) {
	esClient, err := NewStorage(elasticsearch.Config{
		Addresses: s.cfg.Addresses,
		Username:  s.cfg.UserName,
		Password:  s.cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

// Config defines the configuration structure for Elasticsearch storage.
type Config struct {
	Addresses []string `env:"ES_ADDRESSES"`
	UserName  string   `env:"ES_USER"`
	Password  string   `env:"ES_PASSWORD"`
}

// NewSearchStorageGetter creates a new instance of the SearchStorageGetter with
// the given Elasticsearch addresses, user name, and password.
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

// NewResourceStorageGetter creates a new instance of the ResourceStorageGetter
// with the given Elasticsearch addresses, user name, and password.
func NewResourceStorageGetter(addresses []string, userName, password string) *ResourceStorageGetter {
	cfg := &Config{
		Addresses: addresses,
		UserName:  userName,
		Password:  password,
	}

	return &ResourceStorageGetter{
		cfg,
	}
}

// NewResourceGroupRuleStorageGetter creates a new instance of the
// ResourceGroupRuleStorageGetter with the given Elasticsearch addresses, user
// name, and password.
func NewResourceGroupRuleStorageGetter(addresses []string, userName, password string) *ResourceGroupRuleStorageGetter {
	cfg := &Config{
		Addresses: addresses,
		UserName:  userName,
		Password:  password,
	}

	return &ResourceGroupRuleStorageGetter{
		cfg,
	}
}
