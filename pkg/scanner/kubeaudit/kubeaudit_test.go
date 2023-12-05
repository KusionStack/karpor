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

package kubeaudit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKubeauditScannerScanManifest(t *testing.T) {
	// Load the testdata file containing the real YAML manifest for a Kubernetes
	// Deployment.
	testdataDir := "testdata"
	deploymentFile := "deployment.yaml"
	deploymentFilePath := filepath.Join(testdataDir, deploymentFile)

	file, err := os.Open(deploymentFilePath)
	if err != nil {
		t.Fatalf("Failed to open test data file: %s", err)
	}
	defer file.Close()

	// Create a new kubeauditScanner instance.
	ks, err := New()
	if err != nil {
		t.Fatalf("Failed to create kubeauditScanner: %s", err)
	}

	// Run the ScanManifest method with the loaded YAML file.
	issues, err := ks.ScanManifest(file)
	if err != nil {
		t.Errorf("ScanManifest returned an error: %s", err)
	}

	// Check if any issues were detected.
	if len(issues) == 0 {
		t.Error("Expected to find issues, but none were detected")
	}

	// Print the issues for debugging.
	for _, issue := range issues {
		t.Logf("Issue found: %v", issue)
	}

	// Perform more detailed checks on the issues.
	if assert.Equal(t, 1, len(issues), "Unexpected number of issues found") {
		assert.Equal(t, scanner.High.String(), issues[0].Severity.String(), "Unexpected issue severity")
		assert.Equal(t, "AutomountServiceAccountTokenTrueAndDefaultSA", issues[0].Title, "Unexpected issue title")
	}
}

// TestKubeauditScannerScan tests scanning a Deployment with kubeauditScanner.
func TestKubeauditScannerScan(t *testing.T) {
	// Create a test Deployment object.
	deployment := createTestDeployment("test-deployment", 3)

	// Create our kubeauditScanner (for the full implementation, you must have the New() function available in your package).
	ks, err := New()
	if err != nil {
		t.Fatalf("Failed to create kubeauditScanner: %s", err)
	}

	// Run the Scan method with the Deployment object.
	issues, err := ks.Scan(deployment)
	if err != nil {
		t.Errorf("Scan returned an error: %s", err)
	}

	// Validate the issues (the specifics will depend on the expected output of kubeaudit).
	// For now, we will just check that the `Scan` function does not return an error and we get some issues back.
	if len(issues) == 0 {
		t.Errorf("Expected to find issues with the test deployment, but none were found")
	}

	// Print the issues for debugging.
	for _, issue := range issues {
		t.Logf("Issue found: %v", issue)
	}

	// Perform more detailed checks on the issues.
	assert.Equal(t, 0, len(issues), "Unexpected number of issues found")
	// if assert.Equal(t, 1, len(issues), "Unexpected number of issues found") {
	// 	assert.Equal(t, scanner.High.String(), issues[0].Severity.String(), "Unexpected issue severity")
	// 	assert.Equal(t, "AutomountServiceAccountTokenTrueAndDefaultSA", issues[0].Title, "Unexpected issue title")
	// }
}

func createTestDeployment(name string, replicas int32) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "web",
						},
					},
				},
			},
		},
	}
}
