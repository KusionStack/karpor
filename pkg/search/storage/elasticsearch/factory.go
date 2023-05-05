package elasticsearch

import (
	"github.com/KusionStack/karbour/pkg/search/storage"
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
