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
	"io"

	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func generateIndexRequest(cluster string, obj runtime.Object) (id string, body []byte, err error) {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return
	}

	body, err = json.Marshal(map[string]interface{}{
		apiVersionKey: obj.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		kindKey:       obj.GetObjectKind().GroupVersionKind().Kind,
		nameKey:       metaObj.GetName(),
		namespaceKey:  metaObj.GetNamespace(),
		clusterKey:    cluster,
		objectKey:     metaObj,
	})
	if err != nil {
		return
	}
	id = string(metaObj.GetUID())
	return
}

func generateQuery(cluster, namespace, name string, obj runtime.Object) map[string]interface{} {
	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term(apiVersionKey, obj.GetObjectKind().GroupVersionKind().GroupVersion().String()),
		esquery.Term(kindKey, obj.GetObjectKind().GroupVersionKind().Kind),
		esquery.Term(nameKey, name),
		esquery.Term(namespaceKey, namespace),
		esquery.Term(clusterKey, cluster),
	).Map()
	return query
}

func (s *ESClient) insertObj(ctx context.Context, cluster string, obj runtime.Object) error {
	id, body, err := generateIndexRequest(cluster, obj)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Index:      s.indexName,
	}
	res, err := req.Do(ctx, s.client)
	if err != nil {
		return err
	}

	if res.IsError() {
		return &ESError{
			StatusCode: res.StatusCode,
			Message:    res.String(),
		}
	}
	return nil
}

func (s *ESClient) deleteByQuery(ctx context.Context, query map[string]interface{}) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return err
	}
	return s.delete(ctx, buf)
}

func (s *ESClient) delete(ctx context.Context, body io.Reader) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return err
	}
	req := esapi.DeleteByQueryRequest{
		Index: []string{s.indexName},
		Body:  buf,
	}
	res, err := req.Do(ctx, s.client)
	if err != nil {
		return err
	}

	if res.IsError() {
		return &ESError{
			StatusCode: res.StatusCode,
			Message:    res.String(),
		}
	}
	return nil
}
