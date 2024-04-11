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

package elasticsearch

import (
	"context"
	"strings"

	"github.com/KusionStack/karbour/pkg/infra/persistence/elasticsearch"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/kubernetes/scheme"
	esv8 "github.com/elastic/go-elasticsearch/v8"
	"k8s.io/apimachinery/pkg/runtime"
	runtimejson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

const (
	clusterKey           = "cluster"
	apiVersionKey        = "apiVersion"
	kindKey              = "kind"
	namespaceKey         = "namespace"
	nameKey              = "name"
	labelsKey            = "labels"
	annotationsKey       = "annotations"
	creationTimestampKey = "creationTimestamp"
	deletionTimestampKey = "deletionTimestamp"
	ownerReferencesKey   = "ownerReferences"
	resourceVersionKey   = "resourceVersion"
	contentKey           = "content"
)

var (
	_ storage.Storage       = &Storage{}
	_ storage.SearchStorage = &Storage{}
)

// Storage is the struct that holds the necessary fields for interacting with the Elasticsearch cluster.
type Storage struct {
	client        *elasticsearch.Client
	indexName     string
	objectEncoder runtime.Encoder
}

// NewStorage creates and returns a new instance of the Storage struct with the provided Elasticsearch configuration.
func NewStorage(cfg esv8.Config) (*Storage, error) {
	cl, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if err = cl.CreateIndex(context.TODO(), defaultIndexName, strings.NewReader(defaultMapping)); err != nil {
		return nil, err
	}

	return &Storage{
		client:    cl,
		indexName: defaultIndexName,
		objectEncoder: runtimejson.NewSerializerWithOptions(
			runtimejson.DefaultMetaFactory,
			scheme.Scheme,
			scheme.Scheme,
			runtimejson.SerializerOptions{Yaml: false, Pretty: true, Strict: true}),
	}, nil
}
