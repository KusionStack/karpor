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

package utils

import (
	"context"

	"github.com/KusionStack/karbour/pkg/infra/search/storage/elasticsearch"
	"github.com/elliotxx/esquery"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

// Importer defines the interface for importing data to a specified storage.
type Importer interface {
	ImportTo(ctx context.Context, store cache.Store) error
}

var _ Importer = &ESImporter{}

// ESImporter is the struct responsible for importing data to an Elasticsearch storage.
type ESImporter struct {
	cluster  string
	esClient *elasticsearch.Storage
	gvr      schema.GroupVersionResource
}

// NewESImporter creates a new instance of the ESImporter struct with the provided Elasticsearch client, cluster name, and GroupVersionResource.
func NewESImporter(esClient *elasticsearch.Storage, cluster string, gvr schema.GroupVersionResource) *ESImporter {
	return &ESImporter{
		cluster:  cluster,
		esClient: esClient,
		gvr:      gvr,
	}
}

// ImportTo implements the Importer interface by importing data to the specified cache.Store within the given context. The method should be implemented to perform the actual data import.
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

	for _, r := range sr.Resources {
		obj := &unstructured.Unstructured{}
		obj.SetUnstructuredContent(r.Object)
		err = store.Add(obj)
		if err != nil {
			return err
		}
	}
	return nil
}
