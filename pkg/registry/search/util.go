package search

import (
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/KusionStack/karbour/pkg/search/storage"
)

func NewSearchStorage(c registry.ExtraConfig) (storage.SearchStorage, error) {
	storage := RESTStorageProvider{
		SearchStorageType:      c.SearchStorageType,
		ElasticSearchAddresses: c.ElasticSearchAddresses,
		ElasticSearchName:      c.ElasticSearchName,
		ElasticSearchPassword:  c.ElasticSearchPassword,
	}

	searchStorageGetter, err := storage.SearchStorageGetter()
	if err != nil {
		return nil, err
	}

	return searchStorageGetter.GetSearchStorage()
}
