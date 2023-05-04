package esstorage

import (
	"bytes"
	"context"
	"encoding/json"

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

func (s *ESClient) searchByQuery(ctx context.Context, query map[string]interface{}) (*SearchResponse, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.indexName),
		s.client.Search.WithBody(buf),
	)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, &ESError{
			StatusCode: res.StatusCode,
			Message:    res.String(),
		}
	}
	sr := &SearchResponse{}
	if err := json.NewDecoder(res.Body).Decode(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *ESClient) deleteByQuery(ctx context.Context, query map[string]interface{}) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
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
