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

package search

import (
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/KusionStack/karbour/pkg/search/storage"
)

func NewSearchStorage(c registry.ExtraConfig) (storage.SearchStorage, error) {
	storage := RESTStorageProvider{
		SearchStorageType:      c.SearchStorageType,
		ElasticSearchAddresses: c.ElasticSearchAddresses,
		ElasticSearchName:      c.ElasticSearchUsername,
		ElasticSearchPassword:  c.ElasticSearchPassword,
	}

	searchStorageGetter, err := storage.SearchStorageGetter()
	if err != nil {
		return nil, err
	}

	return searchStorageGetter.GetSearchStorage()
}
