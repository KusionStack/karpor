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
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
)

// NewSearchStorage creates a new instance of a search storage component using the provided extra configuration.
func NewSearchStorage(c registry.ExtraConfig) (storage.SearchStorage, error) {
	storage := RESTStorageProvider{
		SearchStorageType: c.SearchStorageType,
		SearchAddresses:   c.SearchAddresses,
		SearchName:        c.SearchUsername,
		SearchPassword:    c.SearchPassword,
	}

	searchStorageGetter, err := storage.SearchStorageGetter()
	if err != nil {
		return nil, err
	}

	return searchStorageGetter.GetSearchStorage()
}

// NewResourceStorage creates a new instance of a resource storage component using the provided extra configuration.
func NewResourceStorage(c registry.ExtraConfig) (storage.ResourceStorage, error) {
	storage := RESTStorageProvider{
		SearchStorageType: c.SearchStorageType,
		SearchAddresses:   c.SearchAddresses,
		SearchName:        c.SearchUsername,
		SearchPassword:    c.SearchPassword,
	}

	resourceStorageGetter, err := storage.ResourceStorageGetter()
	if err != nil {
		return nil, err
	}

	return resourceStorageGetter.GetResourceStorage()
}

// NewResourceGroupRuleStorage creates a new instance of a resource group rule storage component using the provided extra configuration.
func NewResourceGroupRuleStorage(c registry.ExtraConfig) (storage.ResourceGroupRuleStorage, error) {
	storage := RESTStorageProvider{
		SearchStorageType: c.SearchStorageType,
		SearchAddresses:   c.SearchAddresses,
		SearchName:        c.SearchUsername,
		SearchPassword:    c.SearchPassword,
	}

	resourceGroupRuleStorageGetter, err := storage.ResourceGroupRuleStorageGetter()
	if err != nil {
		return nil, err
	}

	return resourceGroupRuleStorageGetter.GetResourceGroupRuleStorage()
}

// NewGeneralStorage creates a new instance of a general storage component using the provided extra configuration.
func NewGeneralStorage(c registry.ExtraConfig) (storage.Storage, error) {
	storage := RESTStorageProvider{
		SearchStorageType: c.SearchStorageType,
		SearchAddresses:   c.SearchAddresses,
		SearchName:        c.SearchUsername,
		SearchPassword:    c.SearchPassword,
	}

	generalStorageGetter, err := storage.GeneralStorageGetter()
	if err != nil {
		return nil, err
	}

	return generalStorageGetter.GetGeneralStorage()
}
