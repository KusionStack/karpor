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
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *ESClient) Save(ctx context.Context, cluster string, obj runtime.Object) error {
	return s.insertObj(ctx, cluster, obj)
}

func (s *ESClient) Delete(ctx context.Context, cluster string, obj runtime.Object) error {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	query := generateQuery(cluster, metaObj.GetNamespace(), metaObj.GetName(), obj)
	err = s.deleteByQuery(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (s *ESClient) Get(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		// TODO: support other implement of runtime.Object
		return fmt.Errorf("only support *unstructured.Unstructured type")
	}

	query := generateQuery(cluster, unObj.GetNamespace(), unObj.GetName(), unObj)
	sr, err := s.searchByQuery(ctx, query, 10, 0)
	if err != nil {
		return err
	}

	resources := sr.GetResources()
	if len(resources) != 1 {
		return fmt.Errorf("query result expected 1, got %d", len(resources))
	}

	unStructContent, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resources[0].Object)
	if err != nil {
		return err
	}
	unObj.SetUnstructuredContent(unStructContent)
	return nil
}
