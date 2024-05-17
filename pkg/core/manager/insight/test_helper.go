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

package insight

import (
	"context"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/multicluster"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

// mockMultiClusterClient returns a mock MultiClusterClient for testing purposes.
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

// Namespace sets the namespace on the mock NamespaceableResource.
func (m *mockNamespaceableResource) Namespace(namespace string) dynamic.ResourceInterface {
	return &mockResource{}
}

// List retrieves a list of unstructured resources from the mock NamespaceableResource.
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

// Get retrieves a single unstructured resource from the mock ResourceInterface.
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

// List retrieves a list of unstructured resources from the mock ResourceInterface.
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
type mockSearchStorage struct {
	storage.SearchStorage
}

// mockResourceStorage is an in-memory implementation of the ResourceStorage
// interface for testing purposes.
type mockResourceStorage struct {
	storage.ResourceStorage
}

// mockResourceGroupRuleStorage is an in-memory implementation of the
// ResourceGroupRuleStorage interface for testing purposes.
type mockResourceGroupRuleStorage struct {
	storage.ResourceGroupRuleStorage
}

// Search implements the search operation returning a single mock resource.
func (m *mockSearchStorage) Search(ctx context.Context, queryString, patternType string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	return &storage.SearchResult{
		Total: 1,
		Resources: []*storage.Resource{{
			ResourceGroup: entity.ResourceGroup{
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

// newMockConfigmap creates a mock Unstructured object representing a ConfigMap resource.
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

// newMockSecret creates a mock Unstructured object representing a Secret resource.
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

// newMockPod creates a mock Unstructured object representing a Pod resource.
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

// mockClusterTopologyMapForCluster returns a mock map of ClusterTopology for testing purposes.
func mockClusterTopologyMapForCluster() map[string]ClusterTopology {
	return map[string]ClusterTopology{
		".v1.Node": {
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Node":                           "parent",
				".v1.PersistentVolumeClaim":          "child",
				".v1.Secret":                         "child",
				".v1.Service":                        "parent",
				"apps.v1.ReplicaSet":                 "parent",
				"policy.v1beta1.PodDisruptionBudget": "parent",
			},
		},
		".v1.Secret": {
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Service",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Pod": "child",
			},
		},
		"apps.v1.Deployment": {
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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

// mockClusterTopologyMapForClusterNamespace returns a mock map of ClusterTopology for testing purposes, focused on cluster namespaces.
func mockClusterTopologyMapForClusterNamespace() map[string]ClusterTopology {
	return map[string]ClusterTopology{
		".v1.Node": {
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
				Cluster:    "existing-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
			},
			Count: 1,
			Relationship: map[string]string{
				".v1.Node":                           "parent",
				".v1.PersistentVolumeClaim":          "child",
				".v1.Secret":                         "child",
				".v1.Service":                        "parent",
				"apps.v1.ReplicaSet":                 "parent",
				"policy.v1beta1.PodDisruptionBudget": "parent",
			},
		},
		".v1.Secret": {
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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
		"apps.v1.Deployment": {
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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
			ResourceGroup: entity.ResourceGroup{
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
