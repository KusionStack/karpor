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
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
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
	ks := createTestScanner(t)

	// Run the ScanManifest method with the loaded YAML file.
	issues, err := ks.ScanManifest(context.Background(), file)
	if err != nil {
		t.Errorf("ScanManifest returned an error: %s", err)
	}

	commonCheck(t, issues)
}

// TestKubeauditScannerScan tests scanning a Deployment with kubeauditScanner.
func TestKubeauditScannerScan(t *testing.T) {
	// Create a test Deployment object.
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}

	// Create a new kubeauditScanner instance.
	ks := createTestScanner(t)

	// Run the Scan method with the Deployment object.
	issues, err := ks.Scan(context.Background(), deployment)
	if err != nil {
		t.Errorf("Scan returned an error: %s", err)
	}

	commonCheck(t, issues)
}

func commonCheck(t *testing.T, issues []*scanner.Issue) {
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

func createTestScanner(t *testing.T) scanner.KubeScanner {
	ks, err := New()
	if err != nil {
		t.Fatalf("Failed to create kubeauditScanner: %s", err)
	}
	return ks
}
