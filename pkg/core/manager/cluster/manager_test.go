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
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestClusterManager_SanitizeKubeConfigWithYAML(t *testing.T) {
	// Create a mock context with a logger (use your actual logger initialization)
	ctx := context.Background()

	// Create a fake KubeConfig in YAML format with sensitive data
	fakeYAML := `
apiVersion: v1
kind: Config
users:
- name: test-user
  user:
    client-certificate-data: sensitive-cert-data
    client-key-data: sensitive-key-data
    token: sensitive-token
    username: sensitive-username
    password: sensitive-password
clusters:
- name: test-cluster
  cluster:
    certificate-authority-data: sensitive-ca-data
`
	t.Log("Plain YAML output:")
	t.Log(fakeYAML)

	// Initialize ClusterManager
	cm := NewClusterManager(&Config{})

	// Call SanitizeKubeConfigWithYAML method with the fake data
	sanitizedYAML, err := cm.SanitizeKubeConfigWithYAML(ctx, fakeYAML)
	if err != nil {
		t.Fatalf("SanitizeKubeConfigWithYAML returned an unexpected error: %v", err)
	}

	// Verify that the sanitized YAML string does not contain any of the sensitive data
	var sanitizedConfig KubeConfig
	err = yaml.Unmarshal([]byte(sanitizedYAML), &sanitizedConfig)
	if err != nil {
		t.Fatalf("Error unmarshaling sanitized YAML: %v", err)
	}
	t.Log("Sanitized YAML output:")
	t.Log(sanitizedYAML)

	// Check if the sensitive data has been masked properly
	for _, user := range sanitizedConfig.Users {
		if containsSensitiveData(user.User.ClientCertificateData) ||
			containsSensitiveData(user.User.ClientKeyData) ||
			containsSensitiveData(user.User.Token) ||
			containsSensitiveData(user.User.Username) ||
			containsSensitiveData(user.User.Password) {
			t.Errorf("Sensitive data was not properly masked")
		}
	}
	for _, cluster := range sanitizedConfig.Clusters {
		if containsSensitiveData(cluster.Cluster.CertificateAuthorityData) {
			t.Errorf("Sensitive data was not properly masked")
		}
	}
}

// containsSensitiveData is a helper function to check if the given data is masked properly.
// In a real test, you would check if the data is masked according to the maskContent logic.
func containsSensitiveData(data string) bool {
	// This is where you check if the data is in masked format (e.g., does not contain sensitive values)
	// This is a placeholder function and should be replaced with actual logic to verify masked content.
	return strings.Contains(data, "sensitive")
}
