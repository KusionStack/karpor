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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KusionStack/karbour/pkg/core/entity"
	"github.com/KusionStack/karbour/pkg/infra/persistence/elasticsearch"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/kubernetes/scheme"
	esv8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elliotxx/esquery"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	runtimejson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

var (
	_ storage.Storage                  = &Storage{}
	_ storage.ResourceStorage          = &Storage{}
	_ storage.ResourceGroupRuleStorage = &Storage{}
	_ storage.SearchStorage            = &Storage{}
)

// Storage is the struct that holds the necessary fields for interacting with the Elasticsearch cluster.
type Storage struct {
	client                     *elasticsearch.Client
	resourceIndexName          string
	resourceGroupRuleIndexName string
	objectEncoder              runtime.Encoder
}

// NewStorage creates and returns a new instance of the Storage struct with the provided Elasticsearch configuration.
func NewStorage(cfg esv8.Config) (*Storage, error) {
	cl, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if err = cl.CreateIndex(context.TODO(), defaultResourceIndexName, strings.NewReader(defaultResourceMapping)); err != nil {
		return nil, err
	}

	if err = cl.CreateIndex(context.TODO(), defaultResourceGroupRuleIndexName, strings.NewReader(defaultResourceGroupRuleMapping)); err != nil {
		return nil, err
	}

	// Check if the default resource group rule exists, if not, create it.
	if err = createResourceGroupRuleIfNotExists(cl, "namespace"); err != nil {
		return nil, err
	}

	return &Storage{
		client:                     cl,
		resourceIndexName:          defaultResourceIndexName,
		resourceGroupRuleIndexName: defaultResourceGroupRuleIndexName,
		objectEncoder: runtimejson.NewSerializerWithOptions(
			runtimejson.DefaultMetaFactory,
			scheme.Scheme,
			scheme.Scheme,
			runtimejson.SerializerOptions{Yaml: false, Pretty: true, Strict: true}),
	}, nil
}

// createResourceGroupRuleIfNotExists checks if a resource group rule exists and creates it if it does not.
func createResourceGroupRuleIfNotExists(cl *elasticsearch.Client, ruleName string) error {
	// Refresh the index before searching to ensure real-time data.
	if err := cl.Refresh(context.Background(), defaultResourceGroupRuleIndexName); err != nil {
		return err
	}

	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term(resourceKeyName, ruleName),
	).Map()
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return err
	}
	resp, err := cl.SearchDocument(context.TODO(), defaultResourceGroupRuleIndexName, buf)
	if err != nil {
		return err
	}

	if resp.Hits.Total.Value == 0 {
		// If specified resource group rule not found, create it
		id := entity.UUID()
		nowTime := metav1.Now()
		body, err := json.Marshal(map[string]interface{}{
			resourceGroupRuleKeyID:          id,
			resourceGroupRuleKeyName:        ruleName,
			resourceGroupRuleKeyDescription: fmt.Sprintf("Default resource group rule for %s", ruleName),
			resourceGroupRuleKeyFields:      []string{ruleName},
			resourceGroupRuleKeyCreatedAt:   &nowTime,
			resourceGroupRuleKeyUpdatedAt:   &nowTime,
		})
		err = cl.SaveDocument(context.TODO(), defaultResourceGroupRuleIndexName, id, bytes.NewReader(body))
		if err != nil {
			return err
		}
	}
	return nil
}
