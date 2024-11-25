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

package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/elliotxx/esquery"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Purger defines the interface for pruning data in storage.
type Purger interface {
	Purge(ctx context.Context, syncBefore time.Time) error
}

var _ Purger = (*ESPurger)(nil)

// NewESPurger creates an ESPurger which implements the Purger interface.
func NewESPurger(esClient *elasticsearch.Storage, cluster string, gvr schema.GroupVersionResource,
	store cache.Store, onPurge func(obj client.Object),
) *ESPurger {
	return &ESPurger{
		cluster:  cluster,
		esClient: esClient,
		gvr:      gvr,
		onPurge:  onPurge,
		store:    store,
		logger:   ctrl.Log.WithName(fmt.Sprintf("%s-es-purger", gvr.Resource)),
	}
}

type ESPurger struct {
	cluster  string
	esClient *elasticsearch.Storage
	gvr      schema.GroupVersionResource
	onPurge  func(obj client.Object)
	store    cache.Store
	logger   logr.Logger
}

// Purge calls onPurge for objects that do not exist in the cache but have not been deleted in ES.
func (e *ESPurger) Purge(ctx context.Context, syncBefore time.Time) error {
	resource := e.gvr.Resource
	kind := resource[0 : len(resource)-1]

	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term("cluster", e.cluster),
		esquery.Term("apiVersion", e.gvr.GroupVersion().String()),
		esquery.Term("kind", kind),
		esquery.Term("deleted", false),
		esquery.Range("syncAt").Lte(syncBefore),
	).Map()

	sr, err := e.esClient.SearchByQuery(ctx, query, nil)
	if err != nil {
		return err
	}

	for _, r := range sr.Resources {
		obj := &unstructured.Unstructured{}
		obj.SetUnstructuredContent(r.Object)
		key, err := cache.MetaNamespaceKeyFunc(obj)
		if err != nil {
			e.logger.Error(err, "error in getting object key")
			continue
		}

		_, exist, err := e.store.GetByKey(key)
		if err != nil {
			e.logger.Error(err, "error in getting object by key")
			continue
		}

		if !exist {
			e.logger.V(1).Info("found an object that should be purged", "key", key)
			e.onPurge(obj)
		}
	}

	return nil
}
