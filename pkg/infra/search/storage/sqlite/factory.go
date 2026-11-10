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

package sqlite

import "github.com/KusionStack/karpor/pkg/infra/search/storage"

var (
	_ storage.SearchStorageGetter            = &SearchStorageGetter{}
	_ storage.ResourceStorageGetter          = &ResourceStorageGetter{}
	_ storage.ResourceGroupRuleStorageGetter = &ResourceGroupRuleStorageGetter{}
	_ storage.GeneralStorageGetter           = &GeneralStorageGetter{}
)

type SearchStorageGetter struct {
	cfg *Config
}

func (s *SearchStorageGetter) GetSearchStorage() (storage.SearchStorage, error) {
	return NewStorage(s.cfg)
}

type ResourceStorageGetter struct {
	cfg *Config
}

func (s *ResourceStorageGetter) GetResourceStorage() (storage.ResourceStorage, error) {
	return NewStorage(s.cfg)
}

type ResourceGroupRuleStorageGetter struct {
	cfg *Config
}

func (s *ResourceGroupRuleStorageGetter) GetResourceGroupRuleStorage() (storage.ResourceGroupRuleStorage, error) {
	return NewStorage(s.cfg)
}

type GeneralStorageGetter struct {
	cfg *Config
}

func (s *GeneralStorageGetter) GetGeneralStorage() (storage.Storage, error) {
	return NewStorage(s.cfg)
}

type Config struct {
	DataSourceName string `env:"SQLITE_DSN"`
}

func NewSearchStorageGetter(dataSourceName string) *SearchStorageGetter {
	return &SearchStorageGetter{cfg: &Config{DataSourceName: dataSourceName}}
}

func NewResourceStorageGetter(dataSourceName string) *ResourceStorageGetter {
	return &ResourceStorageGetter{cfg: &Config{DataSourceName: dataSourceName}}
}

func NewResourceGroupRuleStorageGetter(dataSourceName string) *ResourceGroupRuleStorageGetter {
	return &ResourceGroupRuleStorageGetter{cfg: &Config{DataSourceName: dataSourceName}}
}

func NewGeneralStorageGetter(dataSourceName string) *GeneralStorageGetter {
	return &GeneralStorageGetter{cfg: &Config{DataSourceName: dataSourceName}}
}
