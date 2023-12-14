// Copyright The Karbour Authors.
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
	"errors"
	"testing"
)

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

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("Test case '%s' failed. Expected error: %v, Got error: %v", test.name, test.expectedErr, err)
			}
		})
	}
}
