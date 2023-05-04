package esstorage

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *ESClient) Create(ctx context.Context, cluster string, obj runtime.Object) error {
	// TODO: check if the resource exists
	return s.insertObj(ctx, cluster, obj)
}

func (s *ESClient) Update(ctx context.Context, cluster string, obj runtime.Object) error {
	// TODO: check if the resource exists
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
	sr, err := s.searchByQuery(ctx, query)
	if err != nil {
		return err
	}

	resources := sr.GetResources()
	if len(resources) != 1 {
		return fmt.Errorf("query result expected 1, got %d", len(resources))
	}

	unStructContent, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resources[0])
	if err != nil {
		return err
	}
	unObj.SetUnstructuredContent(unStructContent)
	return nil
}
