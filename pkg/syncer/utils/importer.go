package utils

import (
	"context"

	"github.com/KusionStack/karbour/pkg/infra/search/storage/elasticsearch"
	"github.com/aquasecurity/esquery"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

type Importer interface {
	ImportTo(ctx context.Context, store cache.Store) error
}

var _ Importer = &ESImporter{}

type ESImporter struct {
	cluster  string
	esClient *elasticsearch.ESClient
	gvr      schema.GroupVersionResource
}

func NewESImporter(esClient *elasticsearch.ESClient, cluster string, gvr schema.GroupVersionResource) *ESImporter {
	return &ESImporter{
		cluster:  cluster,
		esClient: esClient,
		gvr:      gvr,
	}
}

func (e *ESImporter) ImportTo(ctx context.Context, store cache.Store) error {
	resource := e.gvr.Resource
	kind := resource[0 : len(resource)-1]
	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term("cluster", e.cluster),
		esquery.Term("apiVersion", e.gvr.GroupVersion().String()),
		esquery.Term("kind", kind),
	).Map()
	sr, err := e.esClient.SearchByQuery(ctx, query, nil)
	if err != nil {
		return err
	}

	for _, r := range sr.GetResources() {
		obj := &unstructured.Unstructured{}
		obj.SetUnstructuredContent(r.Object)
		err = store.Add(obj)
		if err != nil {
			return err
		}
	}
	return nil
}
