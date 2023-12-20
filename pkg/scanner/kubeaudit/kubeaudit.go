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
	"fmt"
	"io"
	"sync"

	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/KusionStack/karbour/pkg/search/storage"
	kubeauditpkg "github.com/elliotxx/kubeaudit"
	"github.com/elliotxx/kubeaudit/auditors/all"
	"github.com/elliotxx/kubeaudit/config"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
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

	// Default attentionLevel to Low if it's invalid (less than zero).
	if int(attentionLevel) < 0 {
		attentionLevel = scanner.Low
	}

	return &kubeauditScanner{
		kubeAuditor:    kubeAuditor,
		attentionLevel: attentionLevel,
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

func (s *kubeauditScanner) Scan(ctx context.Context, resources ...*storage.Resource) (scanner.ScanResult, error) {
	var wg sync.WaitGroup
	wg.Add(len(resources))

	resultsChan := make(chan scanner.ScanResult, len(resources))
	errChan := make(chan error, len(resources))

	for _, res := range resources {
		go func(res storage.Resource) {
			defer wg.Done()

			resYAML, err := yaml.Marshal(res.Object)
			if err != nil {
				errChan <- err
				return
			}

			report, err := s.scanManifest(ctx, res.Cluster, bytes.NewBuffer(resYAML))
			if err != nil {
				errChan <- err
				return
			}

			resultsChan <- report
		}(*res)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	allReports := newScanResult()

	for report := range resultsChan {
		allReports.MergeBy(report)
	}

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return allReports, nil
}

func (s *kubeauditScanner) scanManifest(ctx context.Context, cluster string, manifest io.Reader) (scanner.ScanResult, error) {
	report, err := s.kubeAuditor.AuditManifest("", manifest)
	if err != nil {
		return nil, err
	}

	results := report.Results()
	if len(results) == 0 {
		return nil, nil
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("the scan result number should be greater than or equal to 1")
	}
	result := results[0]

	resource, err := runtimeObjectToResource(cluster, result.GetResource().Object())
	if err != nil {
		return nil, err
	}

	r := newScanResult()
	issues := []*scanner.Issue{}
	for _, auditResult := range result.GetAuditResults() {
		newIssue := AuditResult2Issue(auditResult)
		if int(newIssue.Severity) >= int(s.attentionLevel) {
			issues = append(issues, newIssue)
		}
	}
	r.add(resource, issues)

	return r, nil
}

func runtimeObjectToResource(cluster string, obj runtime.Object) (*storage.Resource, error) {
	m, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	objMap, err := objectToMap(obj)
	if err != nil {
		return nil, err
	}

	return &storage.Resource{
		Cluster:    cluster,
		APIVersion: obj.GetObjectKind().GroupVersionKind().Version,
		Kind:       obj.GetObjectKind().GroupVersionKind().Kind,
		Namespace:  m.GetNamespace(),
		Name:       m.GetName(),
		Object:     objMap,
	}, nil
}
