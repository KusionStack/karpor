package esstorage

import (
	"context"
	"fmt"

	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *ESClient) Create(ctx context.Context, cluster string, obj runtime.Object) error {
	return s.Insert(ctx, cluster, obj)
}

func (s *ESClient) Update(ctx context.Context, cluster string, obj runtime.Object) error {
	return s.Insert(ctx, cluster, obj)
}

func (s *ESClient) Delete(ctx context.Context, cluster string, obj runtime.Object) error {
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

func (s *ESClient) Get(ctx context.Context, cluster, namespace, apiVerson, kind, name string) (runtime.Object, error) {
	query := esquery.Bool().Must(
		esquery.Term(apiVersionKey, apiVerson),
		esquery.Term(kindKey, kind),
		esquery.Term(nameKey, name),
		esquery.Term(namespaceKey, namespace),
		esquery.Term(clusterKey, cluster),
	).Map()

	sr, err := s.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	resources := sr.GetResources()
	if cnt := len(resources); cnt != 1 {
		return nil, fmt.Errorf("query result expect %d, got %d", 1, cnt)
	}

	unStructContent, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resources[0])
	if err != nil {
		return nil, err
	}
	unObj := &unstructured.Unstructured{}
	unObj.SetUnstructuredContent(unStructContent)
	return unObj, nil
}

func (s *ESClient) List(ctx context.Context, options *storage.ListOptions) (runtime.Object, error) {
	query := esquery.Bool().Must(
		esquery.Term(apiVersionKey, options.APIVersions),
		esquery.Term(kindKey, options.Kinds),
		esquery.Term(nameKey, options.Names),
		esquery.Term(namespaceKey, options.Namespaces),
		esquery.Term(clusterKey, options.Clusters),
	).Map()

	sr, err := s.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	unObjs := &unstructured.UnstructuredList{}
	for _, resource := range sr.GetResources() {
		unStructContent, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resource)
		if err != nil {
			return nil, err
		}
		unObj := &unstructured.Unstructured{}
		unObj.SetUnstructuredContent(unStructContent)
		unObjs.Items = append(unObjs.Items, *unObj)
	}
	return unObjs, nil
}
