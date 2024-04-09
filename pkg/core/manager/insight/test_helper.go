package insight

import (
	"context"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// mockNamespaceableResource is a mock implementation of
// dynamic.NamespaceableResourceInterface.
type mockNamespaceableResource struct {
	dynamic.NamespaceableResourceInterface
}

func (m *mockNamespaceableResource) Namespace(namespace string) dynamic.ResourceInterface {
	return &mockResource{}
}

func (m *mockNamespaceableResource) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return &unstructured.UnstructuredList{
		Object: map[string]interface{}{"kind": "List", "apiVersion": "v1"},
		Items: []unstructured.Unstructured{
			*newMockConfigmap("default", "existing-configmap"),
		},
	}, nil
}

// mockResource is a mock implementation of dynamic.ResourceInterface.
type mockResource struct {
	dynamic.ResourceInterface
}

func (m *mockResource) Get(ctx context.Context, name string, options metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error) {
	if name == "existing-configmap" {
		return newMockConfigmap("default", name), nil
	}
	if name == "existing-secret" {
		return newMockSecret("default", name), nil
	}
	return nil, errors.NewNotFound(schema.GroupResource{Group: "", Resource: ""}, name)
}

// mockSearchStorage is an in-memory implementation of the SearchStorage
// interface for testing purposes.
type mockSearchStorage struct{}

// Search implements the search operation returning a single mock resource.
func (m *mockSearchStorage) Search(ctx context.Context, queryString, patternType string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	return &storage.SearchResult{
		Total: 1,
		Resources: []*storage.Resource{{
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Namespace:  "default",
				Name:       "existing-configmap",
			},
			Object: newMockConfigmap("default", "existing-configmap").Object,
		}},
	}, nil
}

func newMockConfigmap(namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"data": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
}

func newMockSecret(namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"data": map[string]interface{}{
				"key1": "sensitive-value1",
				"key2": "sensitive-value2",
			},
		},
	}
}
