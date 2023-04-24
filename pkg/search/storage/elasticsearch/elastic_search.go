package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *Storage) search(ctx context.Context, query map[string]interface{}) (*SearchResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, &ESError{
			StatusCode: res.StatusCode,
			Message:    res.String(),
		}
	}
	var sr *SearchResponse
	if err := json.NewDecoder(res.Body).Decode(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *Storage) insert(ctx context.Context, cluster string, obj runtime.Object) error {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	body, err := json.Marshal(map[string]interface{}{
		apiVersionKey: obj.GetObjectKind().GroupVersionKind().GroupVersion().String(), // s.storageGroupResource.Group,
		kindKey:       obj.GetObjectKind().GroupVersionKind().Kind,
		nameKey:       metaObj.GetName(),
		namespaceKey:  metaObj.GetNamespace(),
		clusterKey:    cluster,
		objectKey:     metaObj,
	})
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		DocumentID: string(metaObj.GetUID()),
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

func (s *Storage) deleteByQuery(ctx context.Context, query map[string]interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err
	}

	req := esapi.DeleteByQueryRequest{
		Index: []string{s.indexName},
		Body:  &buf,
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

func (s *Storage) deleteByObj(ctx context.Context, obj runtime.Object) error {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	uid := string(metaObj.GetUID())
	if len(uid) == 0 {
		return nil
	}

	req := esapi.DeleteRequest{
		Index:      s.indexName,
		DocumentID: uid,
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
