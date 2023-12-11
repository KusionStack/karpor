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

// Package kubeaudit wraps the kubeaudit library to provide a scanner.KubeScanner
// implementation for auditing Kubernetes resources against common security concerns.
package kubeaudit

import (
	"bytes"
	"context"
	"io"

	"github.com/KusionStack/karbour/pkg/scanner"
	kubeauditpkg "github.com/elliotxx/kubeaudit"
	"github.com/elliotxx/kubeaudit/auditors/all"
	"github.com/elliotxx/kubeaudit/config"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/yaml"
)

// ScannerName is the name of the scanner.
const ScannerName = "KubeAudit"

// Ensure that kubeauditScanner implements the scanner.KubeScanner interface.
var _ scanner.KubeScanner = &kubeauditScanner{}

// kubeauditScanner is an implementation of scanner.KubeScanner that utilizes
// the functionality from the kubeaudit package to perform security audits.
type kubeauditScanner struct {
	kubeAuditor    *kubeauditpkg.Kubeaudit
	attentionLevel scanner.IssueSeverityLevel
	serializer     *json.Serializer
}

// New creates a new instance of a kubeaudit-based scanner with the specified
// attention level.
// The attentionLevel sets a threshold, and only issues that meet or exceed this
// threshold are included in the audit results.
// For example, if the attentionLevel is set to "Medium", then only issues
// classified at the "Medium" level or higher ("Medium", "High", "Critical")
// will be returned to the caller.
func New(attentionLevel scanner.IssueSeverityLevel) (scanner.KubeScanner, error) {
	// Initialize auditors with the kubeaudit configuration.
	auditors, err := all.Auditors(config.KubeauditConfig{})
	if err != nil {
		return nil, err
	}

	// Create a new kubeauditor instance with the configured auditors.
	kubeAuditor, err := kubeauditpkg.New(auditors)
	if err != nil {
		return nil, err
	}

	// Prepare a JSON serializer for serializing the Kubernetes resources.
	serializer := json.NewSerializerWithOptions(
		json.DefaultMetaFactory, nil, nil,
		json.SerializerOptions{Yaml: true, Pretty: false, Strict: false},
	)

	// Default attentionLevel to Low if it's invalid (less than zero).
	if int(attentionLevel) < 0 {
		attentionLevel = scanner.Low
	}

	return &kubeauditScanner{
		kubeAuditor:    kubeAuditor,
		attentionLevel: attentionLevel,
		serializer:     serializer,
	}, nil
}

// New creates a default instance of a kubeaudit-based scanner with the default
// attention level.
func Default() (scanner.KubeScanner, error) {
	return New(scanner.Low)
}

// Name returns the name of the kubeaudit scanner.
func (s *kubeauditScanner) Name() string {
	return ScannerName
}

// Scan audits the provided Kubernetes resources and returns a list of
// security issues found, if any. It serializes the runtime.Object to JSON
// and then uses kubeaudit to perform the auditing.
func (s *kubeauditScanner) Scan(ctx context.Context, resources ...runtime.Object) ([]*scanner.Issue, error) {
	manifest, err := s.serializeObjectsToYAML(resources...)
	if err != nil {
		return nil, err
	}

	return s.ScanManifest(ctx, manifest)
}

// Scan audits the provided Kubernetes resources manifest and returns a list of
// security issues found.
func (s *kubeauditScanner) ScanManifest(ctx context.Context, manifest io.Reader) ([]*scanner.Issue, error) {
	// Audit the specific manifest.
	report, err := s.kubeAuditor.AuditManifest("", manifest)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to collect issues.
	issues := []*scanner.Issue{}

	for _, result := range report.Results() {
		// Process the audit results and convert them to scanner.Issue.
		for _, auditResult := range result.GetAuditResults() {
			newIssue := AuditResult2Issue(auditResult)
			if int(newIssue.Severity) >= int(s.attentionLevel) {
				issues = append(issues, newIssue)
			}
		}
	}

	// Return the list of discovered issues.
	return issues, nil
}

// serializeObjectsToYAML concatenates multiple runtime.Object instances into a
// single YAML string.
func (s *kubeauditScanner) serializeObjectsToYAML(objects ...runtime.Object) (io.Reader, error) {
	var yamlBuffer bytes.Buffer
	for i, obj := range objects {
		// Serialize the object into YAML bytes.
		data, err := runtime.Encode(s.serializer, obj)
		if err != nil {
			return nil, err
		}

		// Convert JSON bytes to YAML.
		yamlData, err := yaml.JSONToYAML(data)
		if err != nil {
			return nil, err
		}

		// Write YAML to buffer, adding the separator if necessary.
		if i > 0 {
			if _, err := yamlBuffer.WriteString("---\n"); err != nil {
				return nil, err
			}
		}
		if _, err := yamlBuffer.Write(yamlData); err != nil {
			return nil, err
		}

		// Append a newline after each object for readability.
		if _, err := yamlBuffer.WriteRune('\n'); err != nil {
			return nil, err
		}
	}

	return &yamlBuffer, nil
}
