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

package multicluster

import (
	"testing"

	clusterv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	"github.com/stretchr/testify/require"
)

// TestNewConfigFromCluster tests the NewConfigFromCluster function with
// various cluster endpoints and credential types.
func TestNewConfigFromCluster(t *testing.T) {
	testCases := []struct {
		name               string
		endpoint           string
		credentialType     clusterv1beta1.CredentialType
		expectedServerName string
	}{
		{
			name:               "Endpoint with port",
			endpoint:           "https://example.com:6443",
			credentialType:     clusterv1beta1.CredentialTypeServiceAccountToken,
			expectedServerName: "example.com",
		},
		{
			name:               "Endpoint without port",
			endpoint:           "https://example.com/k8s/clusters/c-xxxxx",
			credentialType:     clusterv1beta1.CredentialTypeServiceAccountToken,
			expectedServerName: "example.com",
		},
		{
			name:               "IPv6 endpoint with port",
			endpoint:           "https://[fd00::1]:6443",
			credentialType:     clusterv1beta1.CredentialTypeX509Certificate,
			expectedServerName: "fd00::1",
		},
		{
			name:               "IPv6 endpoint without port",
			endpoint:           "https://[fd00::1]",
			credentialType:     clusterv1beta1.CredentialTypeX509Certificate,
			expectedServerName: "fd00::1",
		},
		{
			name:               "Exec credential leaves server name empty",
			endpoint:           "https://example.com",
			credentialType:     clusterv1beta1.CredentialTypeOIDC,
			expectedServerName: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cluster := &clusterv1beta1.Cluster{
				Spec: clusterv1beta1.ClusterSpec{
					Access: clusterv1beta1.ClusterAccess{
						Endpoint: tc.endpoint,
						Credential: &clusterv1beta1.ClusterAccessCredential{
							Type:                tc.credentialType,
							ServiceAccountToken: "token",
							X509:                &clusterv1beta1.X509{},
							ExecConfig:          &clusterv1beta1.ExecConfig{Command: "aws"},
						},
					},
				},
			}

			cfg, err := NewConfigFromCluster(cluster)
			require.NoError(t, err)
			require.Equal(t, tc.endpoint, cfg.Host)
			require.Equal(t, tc.expectedServerName, cfg.ServerName)
		})
	}
}
