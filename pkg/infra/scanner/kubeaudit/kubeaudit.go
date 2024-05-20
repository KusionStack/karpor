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

// Package kubeaudit wraps the kubeaudit library to provide a scanner.KubeScanner
// implementation for auditing Kubernetes resources against common security concerns.
package kubeaudit

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"kusionstack.io/karpor/pkg/core/entity"
	"kusionstack.io/karpor/pkg/infra/scanner"
	"kusionstack.io/karpor/pkg/infra/search/storage"
	"kusionstack.io/karpor/pkg/util/cache"
	"kusionstack.io/karpor/pkg/util/ctxutil"
	"kusionstack.io/karpor/pkg/util/safeutil"
	kubeauditpkg "github.com/elliotxx/kubeaudit"
	"github.com/elliotxx/kubeaudit/auditors/all"
	"github.com/elliotxx/kubeaudit/config"
	"github.com/elliotxx/safe"
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
	c              *cache.Cache[entity.ResourceGroupHash, scanner.ScanResult]
}

// New creates a new instance of a kubeaudit-based scanner with the specified
// attention level.
// The attentionLevel sets a threshold, and only issues that meet or exceed this
// threshold are included in the audit results.
// For example, if the attentionLevel is set to "Medium", then only issues
// classified at the "Medium" level or higher ("Medium", "High", "Critical")
// will be returned to the caller.
func New(attentionLevel scanner.IssueSeverityLevel) (scanner.KubeScanner, error) {
	const defaultExpiration = 10 * time.Minute

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
		c:              cache.NewCache[entity.ResourceGroupHash, scanner.ScanResult](defaultExpiration),
	}, nil
}

// Default creates a default instance of a kubeaudit-based scanner with the
// default attention level.
func Default() (scanner.KubeScanner, error) {
	return New(scanner.Low)
}

// Name returns the name of the kubeaudit scanner.
func (s *kubeauditScanner) Name() string {
	return ScannerName
}

// Scan audits the provided Kubernetes resources and returns security issues
// found during the scan.
func (s *kubeauditScanner) Scan(ctx context.Context, noCache bool, resources ...*storage.Resource) (scanner.ScanResult, error) {
	var wg sync.WaitGroup
	wg.Add(len(resources))

	resultsChan := make(chan scanner.ScanResult, len(resources))
	errChan := make(chan error, len(resources))

	for _, res := range resources {
		go func(res *storage.Resource) {
			defer safe.HandleCrash(safeutil.RecoverHandler(ctx, errChan))
			defer wg.Done()

			resYAML, err := yaml.Marshal(res.Object)
			if err != nil {
				errChan <- err
				return
			}

			result, err := s.scanManifest(ctx, noCache, res, resYAML)
			if err != nil {
				errChan <- err
				return
			}

			resultsChan <- result
		}(res)
	}

	go func() {
		defer safe.HandleCrash(safeutil.RecoverHandler(ctx, errChan))
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	allResults := newScanResult()

	for report := range resultsChan {
		allResults.MergeFrom(report)
	}

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return allResults, nil
}

// scanManifest performs the actual scanning on the Kubernetes manifest and
// returns the scan result.
func (s *kubeauditScanner) scanManifest(ctx context.Context, noCache bool, resource *storage.Resource, manifest []byte) (scanner.ScanResult, error) {
	if noCache {
		return s.scanManifestFor(ctx, resource, manifest)
	} else {
		if scanResult, exist := s.c.Get(resource.ResourceGroup.Hash()); exist {
			return scanResult, nil
		} else {
			return s.scanManifestFor(ctx, resource, manifest)
		}
	}
}

// scanManifestFor is a method that is likely intended for internal use within the kubeauditScanner. It scans the provided manifest for issues and returns the scan result.
func (s *kubeauditScanner) scanManifestFor(ctx context.Context, resource *storage.Resource, manifest []byte) (scanner.ScanResult, error) {
	// Extract the context and logger from the request.
	log := ctxutil.GetLogger(ctx)

	report, err := s.kubeAuditor.AuditSingleManifest(bytes.NewBuffer(manifest))
	if err != nil {
		return nil, err
	}

	results := report.RawResults()
	if len(results) != 1 {
		log.Info("invalid manifest", "manifest", string(manifest))
		return nil, fmt.Errorf("the scan result number should be equal to 1")
	}
	result := results[0]

	r := newScanResult()
	issues := []*scanner.Issue{}
	for _, auditResult := range result.GetAuditResults() {
		newIssue := AuditResult2Issue(auditResult)
		if int(newIssue.Severity) >= int(s.attentionLevel) {
			issues = append(issues, newIssue)
		}
	}
	r.add(resource, issues)
	s.c.Set(resource.ResourceGroup.Hash(), r)

	return r, nil
}
