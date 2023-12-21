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
	_ storage.Storage       = &ESClient{}
	_ storage.SearchStorage = &ESClient{}
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
