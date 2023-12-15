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

package audit

import (
	"context"
	"strings"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/KusionStack/karbour/pkg/registry/search"
	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/KusionStack/karbour/pkg/scanner/kubeaudit"
	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
)

// AuditManager manages the auditing process of Kubernetes manifests using
// a KubeScanner.
type AuditManager struct {
	ks scanner.KubeScanner
	ss storage.SearchStorage
}

// NewAuditManager initializes a new instance of AuditManager with a KubeScanner.
func NewAuditManager(c *registry.ExtraConfig) (*AuditManager, error) {
	// Create a new Kubernetes scanner instance.
	kubeauditScanner, err := kubeaudit.Default()
	if err != nil {
		return nil, err
	}

	storage := search.RESTStorageProvider{
		SearchStorageType:      c.SearchStorageType,
		ElasticSearchAddresses: c.ElasticSearchAddresses,
		ElasticSearchName:      c.ElasticSearchName,
		ElasticSearchPassword:  c.ElasticSearchPassword,
	}
	searchStorageGetter, err := storage.SearchStorageGetter()
	if err != nil {
		return nil, err
	}
	searchStorage, err := searchStorageGetter.GetSearchStorage()
	if err != nil {
		return nil, err
	}

	return &AuditManager{
		ks: kubeauditScanner,
		ss: searchStorage,
	}, nil
}

// Audit performs the audit on Kubernetes manifests with the specified locator
// and returns the issues found during the audit.
func (m *AuditManager) Audit(ctx context.Context, locator *core.Locator) (*AuditData, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting audit with specified condition in AuditManager ...")

	var total int
	searchQuery := locator.ToSQL()
	searchPattern := storage.SQLPatternType
	searchPageSize := 100
	searchPage := 1

	var allIssues []*scanner.Issue
	for {
		log.Info("Starting search in AuditManager ...",
			"searchQuery", searchQuery, "searchPattern", searchPattern, "searchPageSize", searchPageSize, "searchPage", searchPage)

		res, err := m.ss.Search(ctx, searchQuery, searchPattern, searchPageSize, searchPage)
		if err != nil {
			return nil, err
		}
		total = res.Total

		log.Info("Finish current search", "overview", res.Overview())

		manifests, err := res.ToYAML()
		if err != nil {
			return nil, err
		}

		issues, err := m.ks.ScanManifest(ctx, strings.NewReader(manifests))
		if err != nil {
			return nil, err
		}
		allIssues = append(allIssues, issues...)

		if len(res.Resources) < searchPageSize {
			break
		}
		searchPage++
	}

	return NewAuditData(allIssues, total), nil
}

// Audit performs a security audit on the provided manifest, returning a list
// of issues discovered during scanning.
func (m *AuditManager) AuditManifest(ctx context.Context, manifest string) ([]*scanner.Issue, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting audit of the specified manifest in AuditManager ...")

	// Execute the scan using the scanner's ScanManifest method.
	return m.ks.ScanManifest(ctx, strings.NewReader(manifest))
}

// Score calculates a score based on the severity and total number of issues
// identified during the audit. It aggregates statistics on different severity
// levels and generates a cumulative score.
func (m *AuditManager) Score(ctx context.Context, issues []*scanner.Issue) (*ScoreData, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting calculate score with specified issues list in AuditManager ...")

	// Initialize variables to calculate the score.
	issueTotal, severitySum := len(issues), 0
	severityStats := map[string]int{}

	// Summarize severity statistics for all issues.
	for _, issue := range issues {
		severitySum += int(issue.Severity)
		severityStats[issue.Severity.String()] += 1
	}

	// Use the aggregated data to calculate the score.
	score := CalculateScore(issueTotal, severitySum)

	// Prepare the score data including the total, sum and statistics.
	data := &ScoreData{
		Score:             score,
		IssuesTotal:       issueTotal,
		SeveritySum:       severitySum,
		SeverityStatistic: severityStats,
	}

	return data, nil
}
