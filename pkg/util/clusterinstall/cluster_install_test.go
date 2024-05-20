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

package clusterinstall

import (
	"testing"

	clusterv1beta1 "kusionstack.io/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

func TestConvertKubeconfigToCluster(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		description string
		cfg         *rest.Config
		wantCluster *clusterv1beta1.Cluster
		wantErr     bool
	}{
		// Test case with secure setup (non-insecure, with certificate data)
		{
			name:        "SecureCluster",
			displayName: "My Cluster",
			description: "This is a description",
			cfg: &rest.Config{
				Host: "https://securehost:6443",
				TLSClientConfig: rest.TLSClientConfig{
					Insecure: false,
					CertData: []byte("mock"),
					KeyData:  []byte("mock"),
				},
			},
			wantErr: false,
		},
		// Test case with insecure setup (insecure flag on)
		{
			name:        "InsecureCluster",
			displayName: "My Insecure Cluster",
			description: "This is an insecure cluster",
			cfg: &rest.Config{
				Host: "http://insecurehost:8080",
				TLSClientConfig: rest.TLSClientConfig{
					Insecure: true,
					CertData: []byte("mock"),
					KeyData:  []byte("mock"),
				},
			},
			wantErr: false,
		},
		// Test case with service account token
		{
			name:        "TokenCluster",
			displayName: "My Token Cluster",
			description: "This cluster uses a token credential",
			cfg: &rest.Config{
				Host:        "https://tokenhost:6443",
				BearerToken: "tokenDataHere",
				TLSClientConfig: rest.TLSClientConfig{
					Insecure: false,
					CertData: []byte("mock"),
					KeyData:  []byte("mock"),
				},
			},
			wantErr: false,
		},
		// Test case for failure due to lack of credential data
		{
			name:        "NoCredsCluster",
			displayName: "My NoCreds Cluster",
			description: "This cluster lacks credentials",
			cfg: &rest.Config{
				Host: "https://nocreds:6443",
				TLSClientConfig: rest.TLSClientConfig{
					Insecure: false,
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function under test
			cluster, err := ConvertKubeconfigToCluster(tc.name, tc.displayName, tc.description, tc.cfg)
			// Assert that an error occurred when expected
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cluster)
				// Assert the cluster fields are what we expect
				require.Equal(t, tc.name, cluster.Name)
				require.Equal(t, tc.description, cluster.Spec.Description)
				wantDisplayName := tc.displayName
				if wantDisplayName == "" {
					wantDisplayName = tc.name
				}
				require.Equal(t, wantDisplayName, cluster.Spec.DisplayName)
				require.Equal(t, tc.cfg.Host, cluster.Spec.Access.Endpoint)
				if tc.cfg.Insecure {
					require.True(t, *cluster.Spec.Access.Insecure)
				}
				if tc.cfg.CAData != nil {
					require.Equal(t, tc.cfg.CAData, cluster.Spec.Access.CABundle)
				}
				if tc.cfg.CertData != nil && tc.cfg.KeyData != nil {
					require.Equal(t, clusterv1beta1.CredentialTypeX509Certificate, cluster.Spec.Access.Credential.Type)
					require.Equal(t, tc.cfg.CertData, cluster.Spec.Access.Credential.X509.Certificate)
					require.Equal(t, tc.cfg.KeyData, cluster.Spec.Access.Credential.X509.PrivateKey)
				} else if tc.cfg.BearerToken != "" {
					require.Equal(t, clusterv1beta1.CredentialTypeServiceAccountToken, cluster.Spec.Access.Credential.Type)
					require.Equal(t, tc.cfg.BearerToken, cluster.Spec.Access.Credential.ServiceAccountToken)
				}
			}
		})
	}
}
