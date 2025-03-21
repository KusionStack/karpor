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

package meilisearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/KusionStack/karpor/pkg/infra/persistence/meilisearch"
	"time"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/elliotxx/esquery"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	resourceKeyID                = "id"
	resourceKeyCluster           = "cluster"
	resourceKeyAPIVersion        = "apiVersion"
	resourceKeyKind              = "kind"
	resourceKeyNamespace         = "namespace"
	resourceKeyName              = "name"
	resourceKeyLabels            = "labels"
	resourceKeyAnnotations       = "annotations"
	resourceKeyCreationTimestamp = "creationTimestamp"
	resourceKeyDeletionTimestamp = "deletionTimestamp"
	resourceKeyOwnerReferences   = "ownerReferences"
	resourceKeyResourceVersion   = "resourceVersion"
	resourceKeyContent           = "content"
	resourceKeySyncAt            = "syncAt"  // resource save/update/delete time
	resourceKeyDeleted           = "deleted" // indicates whether the resource is deleted in cluster
)

var ErrNotFound = fmt.Errorf("object not found")

// SaveResource stores an object in the Elasticsearch storage for the specified cluster.
func (s *Storage) SaveResource(ctx context.Context, cluster string, obj runtime.Object) error {
	out, err := s.generateResourceDocument(cluster, obj)
	if err != nil {
		return err
	}
	return s.client.SaveDocument(ctx, s.resourceIndexName, out)
}

// Refresh meilisearch storage does not need to refresh, it is automatically refreshed.
func (s *Storage) Refresh(_ context.Context) error {
	return nil
}

// SoftDeleteResource only sets the deleted field to true, not really deletes the data in storage.
func (s *Storage) SoftDeleteResource(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		// TODO: support other implement of runtime.Object
		return fmt.Errorf("only support *unstructured.Unstructured type")
	}

	if err := s.GetResource(ctx, cluster, unObj); err != nil {
		return err
	}

	newObject, err := s.generateResourceDocument(cluster, obj)
	if err != nil {
		return err
	}
	newObject[resourceKeyDeleted] = true
	return s.client.SaveDocument(ctx, s.resourceIndexName, newObject)
}

// DeleteResource removes an object from the Elasticsearch storage for the specified cluster.
func (s *Storage) DeleteResource(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		// TODO: support other implement of runtime.Object
		return fmt.Errorf("only support *unstructured.Unstructured type")
	}

	if err := s.GetResource(ctx, cluster, unObj); err != nil {
		return err
	}

	return s.client.DeleteDocument(ctx, s.resourceIndexName, string(unObj.GetUID()))
}

// GetResource retrieves an object from the Elasticsearch storage for the specified cluster.
func (s *Storage) GetResource(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		// TODO: support other implement of runtime.Object
		return fmt.Errorf("only support *unstructured.Unstructured type")
	}

	filter := generateResourceFilter(cluster, unObj.GetNamespace(), unObj.GetName(), unObj)
	resp, err := s.client.SearchDocument(ctx, s.resourceIndexName, &meilisearch.SearchRequest{Filter: filter})
	if err != nil {
		return err
	}

	if resp.TotalHits == 0 {
		return fmt.Errorf("no resource found for cluster: %s, namespace: %s, name: %s", cluster, unObj.GetNamespace(), unObj.GetName())
	}

	res, err := storage.Map2Resource(resp.Hits[0].(map[string]interface{}))
	if err != nil {
		return err
	}

	unObj.Object = res.Object
	return nil
}

// CountResources return a count of resources in the Elasticsearch storage.
func (s *Storage) CountResources(ctx context.Context) (int, error) {
	if resp, err := s.client.Count(ctx, s.resourceIndexName); err != nil {
		return 0, err
	} else {
		return int(resp.Count), nil
	}
}

// DeleteAllResources removes all resources from the Elasticsearch storage for the specified cluster.
func (s *Storage) DeleteAllResources(ctx context.Context, cluster string) error {
	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term(resourceKeyCluster, cluster),
	).Map()
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return err
	}
	return s.client.DeleteDocumentByQuery(ctx, s.resourceIndexName, buf)
}

// generateResourceDocument creates an resource document for Elasticsearch with
// the specified cluster and object.
func (s *Storage) generateResourceDocument(cluster string, obj runtime.Object) (out map[string]any, err error) {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer([]byte{})
	if err = s.objectEncoder.Encode(obj, buf); err != nil {
		return
	}

	out = map[string]any{
		resourceKeyID:                string(metaObj.GetUID()),
		resourceKeyCluster:           cluster,
		resourceKeyAPIVersion:        obj.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		resourceKeyKind:              obj.GetObjectKind().GroupVersionKind().Kind,
		resourceKeyNamespace:         metaObj.GetNamespace(),
		resourceKeyName:              metaObj.GetName(),
		resourceKeyLabels:            metaObj.GetLabels(),
		resourceKeyAnnotations:       metaObj.GetAnnotations(),
		resourceKeyCreationTimestamp: metaObj.GetCreationTimestamp(),
		resourceKeyDeletionTimestamp: metaObj.GetDeletionTimestamp(),
		resourceKeyOwnerReferences:   metaObj.GetOwnerReferences(),
		resourceKeyResourceVersion:   metaObj.GetResourceVersion(),
		resourceKeyContent:           buf.String(),
		resourceKeySyncAt:            time.Now(),
		resourceKeyDeleted:           false,
	}
	return
}

// generateResourceFilter creates a filter to search for an object in
// meilisearch based on resource's cluster, namespace, and name.
func generateResourceFilter(cluster, namespace, name string, obj runtime.Object) interface{} {
	return []string{
		generateFilter(resourceKeyAPIVersion, obj.GetObjectKind().GroupVersionKind().GroupVersion().String()),
		generateFilter(resourceKeyKind, obj.GetObjectKind().GroupVersionKind().Kind),
		generateFilter(resourceKeyName, name),
		generateFilter(resourceKeyNamespace, namespace),
		generateFilter(resourceKeyCluster, cluster),
	}
}

// CheckStorageHealth checks the health of the Elasticsearch storage by pinging the client.
func (s *Storage) CheckStorageHealth(ctx context.Context) error {
	return s.client.SearchLiveness(ctx)
}
