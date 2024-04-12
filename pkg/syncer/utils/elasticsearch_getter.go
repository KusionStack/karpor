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
	"fmt"
	"strings"

	"github.com/KusionStack/karbour/pkg/infra/search/storage/elasticsearch"
	"github.com/aquasecurity/esquery"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8scache "k8s.io/client-go/tools/cache"
)

var _ k8scache.KeyListerGetter = &ESListerGetter{}

type ESListerGetter struct {
	cluster  string
	esClient *elasticsearch.Storage
	gvr      schema.GroupVersionResource
}

func NewESListerGetter(cluster string, esClient *elasticsearch.Storage, gvr schema.GroupVersionResource) *ESListerGetter {
	return &ESListerGetter{
		cluster:  cluster,
		esClient: esClient,
		gvr:      gvr,
	}
}

// ListKeys returns a list of keys for the resources managed by the ESListerGetter.
func (e *ESListerGetter) ListKeys() []string {
	resource := e.gvr.Resource
	kind := resource[0 : len(resource)-1]
	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term("cluster", e.cluster),
		esquery.Term("apiVersion", e.gvr.GroupVersion().String()),
		esquery.Term("kind", kind),
	).Map()
	sr, err := e.esClient.SearchByQuery(context.Background(), query, nil)
	if err != nil {
		return nil
	}
	rt := []string{}
	for _, r := range sr.Resources {
		name, _, _ := unstructured.NestedString(r.Object, "metadata", "name")
		ns, _, _ := unstructured.NestedString(r.Object, "metadata", "namespace")
		var key string
		if ns != "" && name != "" {
			key = ns + "/" + name
		} else if name != "" {
			key = name
		}
		if key != "" {
			rt = append(rt, key)
		}
	}
	return rt
}

// GetByKey retrieves the value associated with the provided key from the managed resources.
func (e *ESListerGetter) GetByKey(key string) (value interface{}, exists bool, err error) {
	s := strings.Split(key, "/")
	var name, ns string
	switch len(s) {
	case 1:
		name = s[0]
	case 2:
		ns = s[0]
		name = s[1]
	default:
		return nil, false, fmt.Errorf("invalid key:%s", key)
	}

	resource := e.gvr.Resource
	kind := resource[0 : len(resource)-1]
	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term("cluster", e.cluster),
		esquery.Term("apiVersion", e.gvr.GroupVersion().String()),
		esquery.Term("kind", kind),
		esquery.Term("namespace", ns),
		esquery.Term("name", name)).Map()
	sr, err := e.esClient.SearchByQuery(context.Background(), query, nil)
	if err != nil {
		return nil, false, err
	}
	resources := sr.Resources
	if len(resources) != 1 {
		return nil, false, fmt.Errorf("query result expected 1, got %d", len(resources))
	}

	unObj := &unstructured.Unstructured{}
	unObj.SetUnstructuredContent(resources[0].Object)
	return unObj, true, nil
}
