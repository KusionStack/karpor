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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var ErrNotFound = fmt.Errorf("object not found")

func (e *ESClient) Save(ctx context.Context, cluster string, obj runtime.Object) error {
	id, body, err := generateIndexRequest(cluster, obj)
	if err != nil {
		return err
	}
	return e.client.SaveDocument(ctx, e.indexName, id, bytes.NewReader(body))
}

func (e *ESClient) Delete(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		// TODO: support other implement of runtime.Object
		return fmt.Errorf("only support *unstructured.Unstructured type")
	}

	if err := e.Get(ctx, cluster, unObj); err != nil {
		return err
	}

	return e.client.DeleteDocument(ctx, e.indexName, string(unObj.GetUID()))
}

func (e *ESClient) Get(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		// TODO: support other implement of runtime.Object
		return fmt.Errorf("only support *unstructured.Unstructured type")
	}

	query := generateQuery(cluster, unObj.GetNamespace(), unObj.GetName(), unObj)
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return err
	}
	resp, err := e.client.SearchDocument(ctx, e.indexName, buf)
	if err != nil {
		return err
	}

	res, err := Convert(resp)
	if err != nil {
		return err
	}

	resources := res.GetResources()
	if len(resources) == 0 {
		return ErrNotFound
	}
	unObj.Object = resources[0].Object
	return nil
}
