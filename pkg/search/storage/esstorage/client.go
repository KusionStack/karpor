package esstorage

import (
	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/elastic/go-elasticsearch/v8"
)

const (
	apiVersionKey = "apiVersion"
	kindKey       = "kind"
	nameKey       = "name"
	namespaceKey  = "namespace"
	clusterKey    = "cluster"
	objectKey     = "object"
)

var (
	_ storage.Storage  = &ESClient{}
	_ storage.Searcher = &ESClient{}
)

type ESClient struct {
	client    *elasticsearch.Client
	indexName string
}

func NewESClient(cfg elasticsearch.Config) (*ESClient, error) {
	cl, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	indexName := defaultIndexName
	err = createIndex(cl, defaultMapping, indexName)
	if err != nil {
		return nil, err
	}

	return &ESClient{
		client:    cl,
		indexName: indexName,
	}, nil
}
