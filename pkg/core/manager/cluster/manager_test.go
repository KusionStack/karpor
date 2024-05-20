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

package cluster

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"kusionstack.io/karpor/pkg/infra/multicluster"
	clusterv1beta1 "kusionstack.io/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// TestGetCluster tests the GetCluster method of the ClusterManager for various
// scenarios.
func TestGetCluster(t *testing.T) {
	manager := NewClusterManager()
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	testCases := []struct {
		name        string
		clusterName string
		expectError bool
	}{
		{
			name:        "Sanitize existing cluster",
			clusterName: "existing-cluster",
			expectError: false,
		},
		{
			name:        "Attempt to get non-existing cluster",
			clusterName: "non-existing-cluster",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cluster, err := manager.GetCluster(
				context.TODO(),
				&multicluster.MultiClusterClient{},
				tc.clusterName,
			)
			if tc.expectError {
				require.Error(t, err, "Expected an error when getting non-existing cluster.")
			} else {
				require.NoError(t, err, "Did not expect an error when getting existing cluster.")
				require.NotNil(t, cluster, "Expected a non-nil sanitized cluster.")

				// Assert the sensitive value is sanitized
				realSensitiveValue, _, _ := unstructured.NestedString(cluster.Object, "spec", "access", "credential", "serviceAccountToken")
				require.Contains(t, realSensitiveValue, "***", "Expected the serviceAccountToken to be sanitized.")

				realSensitiveValue, _, _ = unstructured.NestedString(cluster.Object, "metadata", "annotations", "kubectl.kubernetes.io/last-applied-configuration")
				require.Equal(t, "[redacted]", realSensitiveValue, "Expected the spec.metadata.annotations.kubectl.kubernetes.io/last-applied-configuration to be sanitized.")

				realSensitiveValue = getByteSliceFieldValue(cluster, "spec", "access", "credential", "x509", "certificate")
				require.Contains(t, realSensitiveValue, "***", "Expected the spec.access.credential.x509.certificate to be sanitized.")

				realSensitiveValue = getByteSliceFieldValue(cluster, "spec", "access", "credential", "x509", "privateKey")
				require.Contains(t, realSensitiveValue, "***", "Expected the spec.access.credential.x509.privateKey to be sanitized.")
			}
		})
	}
}

// TestValidateKubeConfigFor tests the ValidateKubeConfigFor method.
func TestValidateKubeConfigFor(t *testing.T) {
	tests := []struct {
		name        string
		inputConfig *KubeConfig
		expectedErr error
	}{
		{
			name: "Invalid Configuration",
			inputConfig: &KubeConfig{
				APIVersion: "v1",
				Kind:       "Config",
				Clusters:   []ClusterEntry{{Name: "cluster1", Cluster: Cluster{Server: "1.2.3.4"}}},
				Users: []UserEntry{{
					Name: "user1",
					User: User{
						ClientCertificateData: "",
						ClientKeyData:         "",
					},
				}},
			},
			expectedErr: ErrMissingCertificateAuthority,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clusterManager := &ClusterManager{}
			_, err := clusterManager.ValidateKubeConfigFor(context.Background(), test.inputConfig)

			require.ErrorIs(t, err, test.expectedErr)
		})
	}
}

// TestCreateCluster tests the CreateCluster method of the ClusterManager for
// various scenarios.
func TestCreateCluster(t *testing.T) {
	manager := NewClusterManager()
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	testCases := []struct {
		name                 string
		clusterName          string
		displayName          string
		description          string
		kubeconfig           string
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name:        "Create new cluster successfully",
			clusterName: "new-cluster",
			displayName: "New Cluster",
			description: "This is a new cluster.",
			kubeconfig:  newMockKubeConfig(),
			expectError: false,
		},
		{
			name:                 "Attempt to create an existing cluster",
			clusterName:          "existing-cluster",
			displayName:          "Existing Cluster",
			description:          "This cluster already exists.",
			kubeconfig:           newMockKubeConfig(),
			expectError:          true,
			expectedErrorMessage: "cluster existing-cluster already exists. Try updating it instead",
		},
		{
			name:                 "Invalid kubeconfig",
			clusterName:          "invalid-kubeconfig-cluster",
			displayName:          "Invalid KubeConfig Cluster",
			description:          "This cluster has invalid kubeconfig.",
			kubeconfig:           "invalid",
			expectError:          true,
			expectedErrorMessage: "failed to parse kubeconfig",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cluster, err := manager.CreateCluster(
				context.TODO(),
				&multicluster.MultiClusterClient{},
				tc.clusterName,
				tc.displayName,
				tc.description,
				tc.kubeconfig,
			)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					tc.expectedErrorMessage,
					"Unexpected error message received.",
				)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cluster, "Expected a non-nil cluster object.")
				require.Equal(t, tc.clusterName, cluster.GetName(), "Cluster name mismatch.")
			}
		})
	}
}

// TestUpdateMetadata tests the UpdateMetadata method of the ClusterManager for
// various scenarios.
func TestUpdateMetadata(t *testing.T) {
	manager := NewClusterManager()
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	testCases := []struct {
		name          string
		clusterName   string
		displayName   string
		description   string
		expectError   bool
		expectedError string
	}{
		{
			name:        "Update metadata successfully",
			clusterName: "existing-cluster",
			displayName: "Updated Cluster",
			description: "This cluster has been updated.",
			expectError: false,
		},
		{
			name:          "Attempt to update non-existing cluster",
			clusterName:   "non-existing-cluster",
			displayName:   "Updated Cluster",
			description:   "This cluster has been updated.",
			expectError:   true,
			expectedError: "\"non-existing-cluster\" not found",
		},
		{
			name:        "Update metadata with empty display name",
			clusterName: "existing-cluster",
			displayName: "",
			description: "Updated Cluster",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updatedCluster, err := manager.UpdateMetadata(
				context.TODO(),
				&multicluster.MultiClusterClient{},
				tc.clusterName,
				tc.displayName,
				tc.description,
			)
			if tc.expectError {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					tc.expectedError,
					"Unexpected error message received.",
				)
			} else {
				require.NoError(t, err)
				require.NotNil(t, updatedCluster, "Expected a non-nil updated cluster object.")
				if len(tc.displayName) == 0 {
					require.Equal(t, tc.clusterName, updatedCluster.Object["spec"].(map[string]interface{})["displayName"].(string), "Display name mismatch.")
				} else {
					require.Equal(t, tc.displayName, updatedCluster.Object["spec"].(map[string]interface{})["displayName"].(string), "Display name mismatch.")
				}
				require.Equal(t, tc.description, updatedCluster.Object["spec"].(map[string]interface{})["description"], "Description mismatch.")
			}
		})
	}
}

// TestUpdateCredential tests the UpdateCredential method of the ClusterManager
// for various scenarios.
func TestUpdateCredential(t *testing.T) {
	manager := NewClusterManager()
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	testCases := []struct {
		name          string
		clusterName   string
		kubeconfig    string
		expectError   bool
		expectedError string
	}{
		{
			name:        "Update credential successfully",
			clusterName: "existing-cluster",
			kubeconfig:  newMockKubeConfig(),
			expectError: false,
		},
		{
			name:          "Attempt to update credential for non-existing cluster",
			clusterName:   "non-existing-cluster",
			kubeconfig:    newMockKubeConfig(),
			expectError:   true,
			expectedError: "\"non-existing-cluster\" not found",
		},
		{
			name:          "Update credential with invalid kubeconfig",
			clusterName:   "existing-cluster",
			kubeconfig:    "invalid-kubeconfig",
			expectError:   true,
			expectedError: "failed to parse kubeconfig",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updatedCluster, err := manager.UpdateCredential(
				context.TODO(),
				&multicluster.MultiClusterClient{},
				tc.clusterName,
				tc.kubeconfig,
			)
			if tc.expectError {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					tc.expectedError,
					"Unexpected error message received.",
				)
			} else {
				require.NoError(t, err)
				require.NotNil(t, updatedCluster, "Expected a non-nil updated cluster object.")
			}
		})
	}
}

// TestDeleteCluster tests the DeleteCluster method of the ClusterManager for
// various scenarios.
func TestDeleteCluster(t *testing.T) {
	manager := NewClusterManager()
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	testCases := []struct {
		name        string
		clusterName string
		expectError bool
	}{
		{
			name:        "Delete existing cluster successfully",
			clusterName: "existing-cluster",
			expectError: false,
		},
		{
			name:        "Attempt to delete non-existing cluster",
			clusterName: "non-existing-cluster",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.DeleteCluster(
				context.TODO(),
				&multicluster.MultiClusterClient{},
				tc.clusterName,
			)
			if tc.expectError {
				require.Error(t, err, "Expected an error when deleting non-existing cluster.")
			} else {
				require.NoError(t, err, "Did not expect an error when deleting existing cluster.")
			}
		})
	}
}

// TestListCluster tests the ListCluster method for retrieving a list of
// clusters.
func TestListCluster(t *testing.T) {
	manager := NewClusterManager()
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	testCases := []struct {
		name        string
		orderBy     SortCriteria
		descending  bool
		expectError bool
	}{
		{
			name:        "List clusters ordered by Name ascending",
			orderBy:     ByName,
			descending:  false,
			expectError: false,
		},
		{
			name:        "List clusters ordered by CreationDate descending",
			orderBy:     ByTimestamp,
			descending:  true,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := manager.ListCluster(
				context.TODO(),
				&multicluster.MultiClusterClient{},
				tc.orderBy,
				tc.descending,
			)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

// TestGetYAMLForCluster tests the GetYAMLForCluster method.
func TestGetYAMLForCluster(t *testing.T) {
	manager := NewClusterManager()
	mockey.Mock((*dynamic.DynamicClient).Resource).Return(&mockNamespaceableResource{}).Build()
	defer mockey.UnPatchAll()

	testCases := []struct {
		name         string
		clusterName  string
		expectedErr  bool
		expectedYAML string
	}{
		{
			name:        "Get YAML for existing cluster",
			clusterName: "existing-cluster",
			expectedErr: false,
			expectedYAML: `apiVersion: cluster.karpor.io/v1beta1
kind: Cluster
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: '[redacted]'
  name: existing-cluster
spec:
  access:
    credential:
      caBundle: sensitive-ca-bundle
      serviceAccountToken: 0873************************bff6
      x509:
        certificate: ZmNjNCoqKioqKioqKioqKioqKioqKioqKioqKmNmMGM=
        privateKey: M2I5NioqKioqKioqKioqKioqKioqKioqKioqKjY1MzY=
  description: mock-description
  displayName: Existing Cluster
`,
		},
		{
			name:        "Attempt to get YAML for non-existing cluster",
			clusterName: "non-existing-cluster",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			yamlData, err := manager.GetYAMLForCluster(
				context.TODO(),
				&multicluster.MultiClusterClient{},
				tc.clusterName,
			)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, yamlData, "Expected a non-nil YAML data.")
				require.Equal(t, tc.expectedYAML, string(yamlData))
			}
		})
	}
}

// TestSanitizeKubeConfigWithYAML tests the SanitizeKubeConfigWithYAML method.
func TestSanitizeKubeConfigWithYAML(t *testing.T) {
	manager := NewClusterManager()

	testCases := []struct {
		name                        string
		plainKubeConfig             string
		expectedSanitizedKubeConfig string
		expectedErr                 bool
		expectedErrorMessage        string
	}{
		{
			name:            "Sanitize valid kubeconfig",
			plainKubeConfig: newMockKubeConfig(),
			expectedSanitizedKubeConfig: `apiVersion: v1
kind: Config
clusters:
    - name: test-cluster
      cluster:
        server: https://127.0.0.1:6443
        certificate-authority-data: 080a************************929b
contexts:
    - name: test-context
      context:
        cluster: test-cluster
        user: test-user
current-context: test-context
users:
    - name: test-user
      user:
        token: 6b26************************99af
`,
			expectedErr: false,
		},
		{
			name:                 "Attempt to sanitize invalid kubeconfig",
			plainKubeConfig:      "invalid",
			expectedErr:          true,
			expectedErrorMessage: "failed to parse kubeconfig",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sanitizedConfig, err := manager.SanitizeKubeConfigWithYAML(
				context.TODO(),
				tc.plainKubeConfig,
			)

			if tc.expectedErr {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					tc.expectedErrorMessage,
					"Unexpected error message received.",
				)
			} else {
				require.NoError(t, err)
				require.NotNil(t, sanitizedConfig, "Expected a non-nil sanitized config.")
				require.Equal(t, tc.expectedSanitizedKubeConfig, sanitizedConfig)
			}
		})
	}
}

// TestValidateKubeConfigWithYAML tests the ValidateKubeConfigWithYAML method.
func TestValidateKubeConfigWithYAML(t *testing.T) {
	manager := NewClusterManager()

	testCases := []struct {
		name                 string
		plainKubeConfig      string
		expectedErr          bool
		expectedErrorMessage string
	}{
		{
			name:                 "Validate unreachable kubeconfig",
			plainKubeConfig:      newMockKubeConfig(),
			expectedErr:          true,
			expectedErrorMessage: ErrClusterServerConnectivity.Error(),
		},
		{
			name:                 "Validate invalid kubeconfig",
			plainKubeConfig:      "invalid",
			expectedErr:          true,
			expectedErrorMessage: "failed to parse kubeconfig",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validatedVersion, err := manager.ValidateKubeConfigWithYAML(
				context.TODO(),
				tc.plainKubeConfig,
			)

			if tc.expectedErr {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					tc.expectedErrorMessage,
					"Unexpected error message received.",
				)
			} else {
				require.NoError(t, err)
				require.NotNil(t, validatedVersion, "Expected a non-nil validated version.")
			}
		})
	}
}

// newUnstructured creates an unstructured object with the provided details.
func newUnstructured(apiVersion, kind, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

// newMockKubeConfig generates a mock kubeconfig with sensitive information
// included.
func newMockKubeConfig() string {
	return fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- name: test-cluster
  cluster:
    server: https://127.0.0.1:6443
    certificate-authority-data: %s
users:
- name: test-user
  user:
    token: fake-token
contexts:
- name: test-context
  context:
    cluster: test-cluster
    user: test-user
current-context: test-context
`, base64.StdEncoding.EncodeToString([]byte("sensitive-certificate-data")))
}

// newMockCluster is a helper function that creates a mock cluster unstructured
// object.
func newMockCluster(name string) *unstructured.Unstructured {
	// Create the base unstructured object
	unsanitizedCluster := newUnstructured(
		clusterv1beta1.SchemeGroupVersion.String(),
		"Cluster",
		name,
	)

	// Populate the object with the mock data
	unsanitizedCluster.Object["spec"] = map[string]interface{}{
		"displayName": "Existing Cluster",
		"description": "mock-description",
		"access": map[string]interface{}{
			"credential": map[string]interface{}{
				"serviceAccountToken": "sensitive-token",
				"x509": map[string]interface{}{
					"certificate": "sensitive-certificate",
					"privateKey":  "sensitive-private-key",
				},
				"caBundle": "sensitive-ca-bundle",
			},
		},
	}

	// Set annotations on the object
	unsanitizedCluster.SetAnnotations(map[string]string{
		"kubectl.kubernetes.io/last-applied-configuration": "sensitive-configuration",
	})

	return unsanitizedCluster
}

// getByteSliceFieldValue retrieves a byte slice field value from the
// unstructured object.
func getByteSliceFieldValue(obj *unstructured.Unstructured, fields ...string) string {
	if nestedField, found, _ := unstructured.NestedFieldNoCopy(obj.Object, fields...); found {
		if bytes, ok := nestedField.([]byte); ok {
			return string(bytes)
		}
	}
	return ""
}

// mockNamespaceableResource is a mock implementation of
// dynamic.NamespaceableResourceInterface.
type mockNamespaceableResource struct {
	dynamic.NamespaceableResourceInterface
}

// Get retrieves the cluster with the provided name.
func (m *mockNamespaceableResource) Get(
	ctx context.Context,
	name string,
	options metav1.GetOptions,
	subresources ...string,
) (*unstructured.Unstructured, error) {
	if name == "existing-cluster" {
		return newMockCluster("existing-cluster"), nil
	}
	return nil, errors.NewNotFound(clusterv1beta1.Resource("cluster"), name)
}

// List retrieves a list of clusters.
func (m *mockNamespaceableResource) List(
	ctx context.Context,
	opts metav1.ListOptions,
) (*unstructured.UnstructuredList, error) {
	unsanitizedCluster := newMockCluster("existing-cluster")

	return &unstructured.UnstructuredList{
		Object: map[string]interface{}{"kind": "List", "apiVersion": "v1"},
		Items: []unstructured.Unstructured{
			*unsanitizedCluster,
		},
	}, nil
}

// Create creates a new cluster.
func (m *mockNamespaceableResource) Create(
	ctx context.Context,
	obj *unstructured.Unstructured,
	options metav1.CreateOptions,
	subresources ...string,
) (*unstructured.Unstructured, error) {
	return obj, nil
}

// Update updates an existing cluster.
func (m *mockNamespaceableResource) Update(
	ctx context.Context,
	obj *unstructured.Unstructured,
	options metav1.UpdateOptions,
	subresources ...string,
) (*unstructured.Unstructured, error) {
	return obj, nil
}

// Delete deletes the cluster with the provided name.
func (m *mockNamespaceableResource) Delete(
	ctx context.Context,
	name string,
	options metav1.DeleteOptions,
	subresources ...string,
) error {
	if name == "existing-cluster" {
		return nil
	}
	return errors.NewNotFound(clusterv1beta1.Resource("cluster"), name)
}
