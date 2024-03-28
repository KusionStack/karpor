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

	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/aquasecurity/esquery"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var ErrNotFound = fmt.Errorf("object not found")

func (s *Storage) Save(ctx context.Context, cluster string, obj runtime.Object) error {
	id, body, err := s.generateIndexRequest(cluster, obj)
	if err != nil {
		return err
	}
	return s.client.SaveDocument(ctx, s.indexName, id, bytes.NewReader(body))
}

func (s *Storage) Delete(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		// TODO: support other implement of runtime.Object
		return fmt.Errorf("only support *unstructured.Unstructured type")
	}

	if err := s.Get(ctx, cluster, unObj); err != nil {
		return err
	}

	return s.client.DeleteDocument(ctx, s.indexName, string(unObj.GetUID()))
}

func (s *Storage) Get(ctx context.Context, cluster string, obj runtime.Object) error {
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
	resp, err := s.client.SearchDocument(ctx, s.indexName, buf)
	if err != nil {
		return err
	}

	res, err := storage.Map2Resource(resp.Hits.Hits[0].Source)
	if err != nil {
		return err
	}

	unObj.Object = res.Object
	return nil
}

func (s *Storage) DeleteAllResourcesInCluster(ctx context.Context, cluster string) error {
	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term(clusterKey, cluster),
	).Map()
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return err
	}
	return s.client.DeleteDocumentByQuery(ctx, s.indexName, buf)
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

func (s *Storage) generateIndexRequest(cluster string, obj runtime.Object) (id string, body []byte, err error) {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer([]byte{})
	if err = s.objectEncoder.Encode(obj, buf); err != nil {
		return
	}

	body, err = json.Marshal(map[string]interface{}{
		clusterKey:           cluster,
		apiVersionKey:        obj.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		kindKey:              obj.GetObjectKind().GroupVersionKind().Kind,
		namespaceKey:         metaObj.GetNamespace(),
		nameKey:              metaObj.GetName(),
		labelsKey:            metaObj.GetLabels(),
		annotationsKey:       metaObj.GetAnnotations(),
		creationTimestampKey: metaObj.GetCreationTimestamp(),
		deletionTimestampKey: metaObj.GetDeletionTimestamp(),
		ownerReferencesKey:   metaObj.GetOwnerReferences(),
		resourceVersionKey:   metaObj.GetResourceVersion(),
		contentKey:           buf.String(),
	})
	if err != nil {
		return
	}
	id = string(metaObj.GetUID())
	return
}
