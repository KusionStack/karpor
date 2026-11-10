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
	"net"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// TestMaskContent tests the correctness of the maskContent function.
func TestMaskContent(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "NormalCase",
			content:  "test content",
			expected: "9473************************2157",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := maskContent(tc.content)
			require.Equal(t, tc.expected, actual)
		})
	}
}

// mockConn is a mock implementation of the net.Conn interface for testing
// purposes.
type mockConn struct {
	net.Conn
}

func (m *mockConn) Close() error {
	return nil
}

// TestCheckEndpointConnectivity tests the functionality of the
// checkEndpointConnectivity function.
func TestCheckEndpointConnectivity(t *testing.T) {
	// Setup mock to record the dialed address and return a mock connection and
	// nil error for the DialTimeout call.
	var dialedAddress string
	mockey.Mock(net.DialTimeout).To(func(network, address string, timeout time.Duration) (net.Conn, error) {
		dialedAddress = address
		return &mockConn{}, nil
	}).Build()
	defer mockey.UnPatchAll()

	testCases := []struct {
		name            string
		endpoint        string
		expectedAddress string
		expectation     require.ErrorAssertionFunc
	}{
		{
			name:            "ValidEndpoint",
			endpoint:        "https://example.com",
			expectedAddress: "example.com:443",
			expectation:     require.NoError,
		},
		{
			name:            "EndpointWithPort",
			endpoint:        "https://example.com:6443",
			expectedAddress: "example.com:6443",
			expectation:     require.NoError,
		},
		{
			name:            "IPv6Endpoint",
			endpoint:        "https://[fd00::1]/k8s/clusters/c-xxxxx",
			expectedAddress: "[fd00::1]:443",
			expectation:     require.NoError,
		},
		{
			name:            "IPv6EndpointWithPort",
			endpoint:        "https://[fd00::1]:6443",
			expectedAddress: "[fd00::1]:6443",
			expectation:     require.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dialedAddress = ""
			actual := checkEndpointConnectivity(tc.endpoint)
			tc.expectation(t, actual)
			require.Equal(t, tc.expectedAddress, dialedAddress)
		})
	}
}

// TestBuildClientConfigFromKubeConfig tests the functionality of the
// buildClientConfigFromKubeConfig function.
func TestBuildClientConfigFromKubeConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      *KubeConfig
		expectation require.ErrorAssertionFunc
	}{
		{
			name: "ValidConfig",
			config: &KubeConfig{
				Clusters: []ClusterEntry{
					{
						Name: "cluster1",
						Cluster: Cluster{
							Server: "https://example.com",
							CertificateAuthorityData: base64.StdEncoding.EncodeToString(
								[]byte("certificate_authority_data"),
							),
						},
					},
				},
				Users: []UserEntry{
					{
						Name: "user1",
						User: User{
							Username: "username",
							Password: "password",
							ClientCertificateData: base64.StdEncoding.EncodeToString(
								[]byte("client_certificate_data"),
							),
							ClientKeyData: base64.StdEncoding.EncodeToString(
								[]byte("client_key_data"),
							),
						},
					},
				},
				Contexts: []ContextEntry{
					{Name: "context1", Context: Context{Cluster: "cluster1", User: "user1"}},
				},
				CurrentContext: "context1",
			},
			expectation: require.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := buildClientConfigFromKubeConfig(tc.config)
			tc.expectation(t, err)
			require.NotNil(t, actual)
		})
	}
}

// TestBuildClientConfigFromKubeConfigWithExec tests that
// buildClientConfigFromKubeConfig converts the exec section of a KubeConfig
// user into the exec provider of the client config.
func TestBuildClientConfigFromKubeConfigWithExec(t *testing.T) {
	testCases := []struct {
		name                    string
		exec                    *ExecConfig
		expectedInteractiveMode clientcmdapi.ExecInteractiveMode
	}{
		{
			name: "Default interactive mode for v1beta1",
			exec: &ExecConfig{
				APIVersion: "client.authentication.k8s.io/v1beta1",
				Command:    "aws",
				Args:       []string{"eks", "get-token", "--cluster-name", "cluster1"},
				Env:        []ExecEnvVar{{Name: "AWS_ACCESS_KEY_ID", Value: "access-key-id"}},
			},
			expectedInteractiveMode: clientcmdapi.IfAvailableExecInteractiveMode,
		},
		{
			name: "Keep explicit interactive mode",
			exec: &ExecConfig{
				APIVersion:         "client.authentication.k8s.io/v1",
				Command:            "aws",
				Args:               []string{"eks", "get-token", "--cluster-name", "cluster1"},
				Env:                []ExecEnvVar{{Name: "AWS_ACCESS_KEY_ID", Value: "access-key-id"}},
				InstallHint:        "install aws-cli",
				ProvideClusterInfo: true,
				InteractiveMode:    "Never",
			},
			expectedInteractiveMode: clientcmdapi.NeverExecInteractiveMode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := buildClientConfigFromKubeConfig(&KubeConfig{
				Clusters: []ClusterEntry{{Name: "cluster1", Cluster: Cluster{Server: "https://example.com"}}},
				Users:    []UserEntry{{Name: "user1", User: User{Exec: tc.exec}}},
				Contexts: []ContextEntry{
					{Name: "context1", Context: Context{Cluster: "cluster1", User: "user1"}},
				},
				CurrentContext: "context1",
			})
			require.NoError(t, err)
			require.Equal(t, &clientcmdapi.ExecConfig{
				Command:            tc.exec.Command,
				Args:               tc.exec.Args,
				Env:                []clientcmdapi.ExecEnvVar{{Name: "AWS_ACCESS_KEY_ID", Value: "access-key-id"}},
				APIVersion:         tc.exec.APIVersion,
				InstallHint:        tc.exec.InstallHint,
				ProvideClusterInfo: tc.exec.ProvideClusterInfo,
				InteractiveMode:    tc.expectedInteractiveMode,
			}, actual.ExecProvider)
			require.Empty(t, actual.CertData)
			require.Empty(t, actual.KeyData)
		})
	}
}

// TestSanitizeUnstructuredCluster tests the functionality of the
// SanitizeUnstructuredCluster function.
func TestSanitizeUnstructuredCluster(t *testing.T) {
	testCases := []struct {
		name        string
		cluster     *unstructured.Unstructured
		expectation require.ErrorAssertionFunc
	}{
		{
			name: "ValidCluster",
			cluster: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"access": map[string]interface{}{
							"credential": map[string]interface{}{
								"serviceAccountToken": "token",
								"x509": map[string]interface{}{
									"certificate": "certificate_data",
									"privateKey":  "private_key_data",
								},
							},
						},
					},
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{},
					},
				},
			},
			expectation: require.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := SanitizeUnstructuredCluster(context.Background(), tc.cluster)
			tc.expectation(t, err)
			require.NotNil(t, actual)
		})
	}
}

// TestSortUnstructuredList tests the functionality of the SortUnstructuredList
// function.
func TestSortUnstructuredList(t *testing.T) {
	testCases := []struct {
		name        string
		unList      *unstructured.UnstructuredList
		criteria    SortCriteria
		descending  bool
		expectation require.ErrorAssertionFunc
	}{
		{
			name: "ValidList-ByName",
			unList: &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"name": "b",
							},
						},
					},
					{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"name": "a",
							},
						},
					},
				},
			},
			criteria:    ByName,
			descending:  false,
			expectation: require.NoError,
		},
		{
			name: "ValidList-ByTimestamp",
			unList: &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"creationTimestamp": time.Now(),
							},
						},
					},
					{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"creationTimestamp": time.Now().
									Add(time.Duration(time.Minute.Minutes())),
							},
						},
					},
				},
			},
			criteria:    ByTimestamp,
			descending:  false,
			expectation: require.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := SortUnstructuredList(tc.unList, tc.criteria, tc.descending)
			tc.expectation(t, err)
			require.NotNil(t, actual)
		})
	}
}
