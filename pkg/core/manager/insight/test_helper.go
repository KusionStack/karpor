package insight

import (
	"context"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/infra/multicluster"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

func mockMultiClusterClient() *multicluster.MultiClusterClient {
	return &multicluster.MultiClusterClient{
		ClientSet: &kubernetes.Clientset{
			DiscoveryClient: &discovery.DiscoveryClient{},
		},
		DynamicClient: &dynamic.DynamicClient{},
		MetricsClient: &v1beta1.MetricsV1beta1Client{},
	}
}

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
	if name == "existing-pod" {
		return newMockPod("default", name), nil
	}
	return nil, errors.NewNotFound(schema.GroupResource{Group: "", Resource: ""}, name)
}

func (m *mockResource) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return &unstructured.UnstructuredList{
		Object: map[string]interface{}{"kind": "List", "apiVersion": "v1"},
		Items: []unstructured.Unstructured{
			*newMockConfigmap("default", "existing-configmap"),
		},
	}, nil
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
				Kind:       "Pod",
				Namespace:  "default",
				Name:       "existing-pod",
			},
			Object: newMockPod("default", "existing-pod").Object,
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

func newMockPod(namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":  "test-container",
						"image": "nginx:latest",
					},
				},
			},
		},
	}
}

func mockClusterTopologyMapForCluster() map[string]ClusterTopology {
	return map[string]ClusterTopology{
		".v1.Node": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Node",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
			},
		},
		".v1.PersistentVolume": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "PersistentVolume",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.PersistentVolumeClaim": "child",
			},
		},
		".v1.PersistentVolumeClaim": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "PersistentVolumeClaim",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.PersistentVolume": "parent",
				".v1.Pod":              "parent",
			},
		},
		".v1.Pod": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Node":                  "parent",
				".v1.PersistentVolumeClaim": "child",
				".v1.Secret":                "child",
				".v1.Service":               "parent",
				"apps.cafe.cloud.alipay.com.v1alpha1.InPlaceSet": "parent",
				"apps.v1.ReplicaSet":                             "parent",
				"policy.v1beta1.PodDisruptionBudget":             "parent",
			},
		},
		".v1.Secret": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Secret",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "parent",
			},
		},
		".v1.Service": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Service",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
			},
		},
		"apps.cafe.cloud.alipay.com.v1alpha1.CafeDeployment": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "apps.cafe.cloud.alipay.com/v1alpha1",
				Kind:       "CafeDeployment",
			},
			Count: 1,
			Relationship: map[string]string{
				"apps.cafe.cloud.alipay.com.v1alpha1.InPlaceSet": "child",
			},
		},
		"apps.cafe.cloud.alipay.com.v1alpha1.InPlaceSet": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "apps.cafe.cloud.alipay.com/v1alpha1",
				Kind:       "InPlaceSet",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
				"apps.cafe.cloud.alipay.com.v1alpha1.CafeDeployment": "parent",
			},
		},
		"apps.v1.Deployment": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
			Count: 1,
			Relationship: map[string]string{
				"apps.v1.ReplicaSet": "child",
			},
		},
		"apps.v1.ReplicaSet": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "apps/v1",
				Kind:       "ReplicaSet",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod":            "child",
				"apps.v1.Deployment": "parent",
			},
		},
		"policy.v1beta1.PodDisruptionBudget": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "policy/v1beta1",
				Kind:       "PodDisruptionBudget",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
			},
		},
	}
}

func mockClusterTopologyMapForClusterNamespace() map[string]ClusterTopology {
	return map[string]ClusterTopology{
		".v1.Node": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Node",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
			},
		},
		".v1.PersistentVolume": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "PersistentVolume",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.PersistentVolumeClaim": "child",
			},
		},
		".v1.PersistentVolumeClaim": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "PersistentVolumeClaim",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.PersistentVolume": "parent",
				".v1.Pod":              "parent",
			},
		},
		".v1.Pod": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Node":                  "parent",
				".v1.PersistentVolumeClaim": "child",
				".v1.Secret":                "child",
				".v1.Service":               "parent",
				"apps.cafe.cloud.alipay.com.v1alpha1.InPlaceSet": "parent",
				"apps.v1.ReplicaSet":                             "parent",
				"policy.v1beta1.PodDisruptionBudget":             "parent",
			},
		},
		".v1.Secret": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Secret",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "parent",
			},
		},
		".v1.Service": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Service",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
			},
		},
		"apps.cafe.cloud.alipay.com.v1alpha1.CafeDeployment": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "apps.cafe.cloud.alipay.com/v1alpha1",
				Kind:       "CafeDeployment",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				"apps.cafe.cloud.alipay.com.v1alpha1.InPlaceSet": "child",
			},
		},
		"apps.cafe.cloud.alipay.com.v1alpha1.InPlaceSet": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "apps.cafe.cloud.alipay.com/v1alpha1",
				Kind:       "InPlaceSet",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
				"apps.cafe.cloud.alipay.com.v1alpha1.CafeDeployment": "parent",
			},
		},
		"apps.v1.Deployment": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				"apps.v1.ReplicaSet": "child",
			},
		},
		"apps.v1.ReplicaSet": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "apps/v1",
				Kind:       "ReplicaSet",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod":            "child",
				"apps.v1.Deployment": "parent",
			},
		},
		"policy.v1beta1.PodDisruptionBudget": {
			Locator: core.Locator{
				Cluster:    "existing-cluster",
				APIVersion: "policy/v1beta1",
				Kind:       "PodDisruptionBudget",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
			},
		},
	}
}
