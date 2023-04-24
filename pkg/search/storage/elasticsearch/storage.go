package elasticsearch

import (
	"context"
	"fmt"

	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v8"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ storage.Storage = &Storage{}

type Storage struct {
	client    *elasticsearch.Client
	indexName string
}

func NewStorage(cfg elasticsearch.Config) (*Storage, error) {
	cl, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	indexName := defaultIndexName
	err = createIndex(cl, defaultMapping, indexName)
	if err != nil {
		return nil, err
	}

	return &Storage{
		client:    cl,
		indexName: indexName,
	}, nil
}

func (s *Storage) Get(ctx context.Context, cluster, namespace, apiVerson, kind, name string) (runtime.Object, error) {
	query := esquery.Bool().Must(
		esquery.Term(apiVersionKey, apiVerson),
		esquery.Term(kindKey, kind),
		esquery.Term(nameKey, name),
		esquery.Term(namespaceKey, namespace),
		esquery.Term(clusterKey, cluster),
	).Map()

	sr, err := s.search(ctx, query)
	if err != nil {
		return nil, err
	}

	resources := sr.GetResources()
	if cnt := len(resources); cnt != 1 {
		return nil, fmt.Errorf("query result expect %d, got %d", 1, cnt)
	}

	unStructContent, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resources[0])
	if err != nil {
		return nil, err
	}
	unObj := &unstructured.Unstructured{}
	unObj.SetUnstructuredContent(unStructContent)
	return unObj, nil
}

func (s *Storage) Delete(ctx context.Context, cluster, namespace, apiVerson, kind, name string) error {
	query := esquery.Bool().Must(
		esquery.Term(apiVersionKey, apiVerson),
		esquery.Term(kindKey, kind),
		esquery.Term(nameKey, name),
		esquery.Term(namespaceKey, namespace),
		esquery.Term(clusterKey, cluster),
	).Map()

	if err := s.deleteByQuery(ctx, query); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Create(ctx context.Context, cluster string, obj runtime.Object) error {
	return s.insert(ctx, cluster, obj)
}

func (s *Storage) Update(ctx context.Context, cluster string, obj runtime.Object) error {
	return s.insert(ctx, cluster, obj)
}
