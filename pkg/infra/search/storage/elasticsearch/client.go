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
	esv8 "github.com/elastic/go-elasticsearch/v8"
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
	_ storage.Storage       = &Storage{}
	_ storage.SearchStorage = &Storage{}
)

type Storage struct {
	client    *elasticsearch.Client
	indexName string
}

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
	}, nil
}
