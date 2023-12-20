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
	"time"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/KusionStack/karbour/pkg/scanner/kubeaudit"
	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/KusionStack/karbour/pkg/util/cache"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
)

// AuditManager manages the auditing process of Kubernetes manifests using
// a KubeScanner.
type AuditManager struct {
	ks scanner.KubeScanner
	ss storage.SearchStorage
	c  *cache.Cache[core.Locator, scanner.ScanResult]
}

// NewAuditManager initializes a new instance of AuditManager with a KubeScanner.
func NewAuditManager(searchStorage storage.SearchStorage) (*AuditManager, error) {
	const defaultExpiration = 30 * time.Minute

	// Create a new Kubernetes scanner instance.
	kubeauditScanner, err := kubeaudit.Default()
	if err != nil {
		return nil, err
	}

	return &AuditManager{
		ks: kubeauditScanner,
		ss: searchStorage,
		c:  cache.NewCache[core.Locator, scanner.ScanResult](defaultExpiration),
	}, nil
}

// Audit performs the audit on Kubernetes manifests with the specified locator
// and returns the issues found during the audit.
func (m *AuditManager) Audit(ctx context.Context, locator core.Locator) (scanner.ScanResult, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting audit with specified condition in AuditManager ...")

	searchQuery := locator.ToSQL()
	searchPattern := storage.SQLPatternType
	pageSizeIteration := 100
	pageIteration := 1

	if auditData, exist := m.c.Get(locator); exist {
		log.Info("Cache hit for locator", "locator", locator)
		return auditData, nil
	} else {
		log.Info("Cache miss for locator", "locator", locator)

		var result scanner.ScanResult
		for {
			log.Info("Starting search in AuditManager ...",
				"searchQuery", searchQuery, "searchPattern", searchPattern, "searchPageSize", pageSizeIteration, "searchPage", pageIteration)

			res, err := m.ss.Search(ctx, searchQuery, searchPattern, pageSizeIteration, pageIteration)
			if err != nil {
				return nil, err
			}

			log.Info("Finish current search", "overview", res.Overview())

			newResult, err := m.ks.Scan(ctx, res.Resources...)
			if err != nil {
				return nil, err
			}
			if result == nil {
				result = newResult
			} else {
				result.MergeBy(newResult)
			}

			if len(res.Resources) < pageSizeIteration {
				break
			}
			pageIteration++
		}

		m.c.Set(locator, result)
		log.Info("Added data to cache for locator", "locator", locator)

		return result, nil
	}
}

// Score calculates a score based on the severity and total number of issues
// identified during the audit. It aggregates statistics on different severity
// levels and generates a cumulative score.
func (m *AuditManager) Score(ctx context.Context, locator core.Locator) (*ScoreData, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting calculate score with specified issues list in AuditManager ...")

	scanResult, err := m.Audit(ctx, locator)
	if err != nil {
		return nil, err
	}

	// Initialize variables to calculate the score.
	issueTotal, severitySum := scanResult.IssueTotal(), 0
	severityStats := map[string]int{}

	// Summarize severity statistics for all issues.
	for issue, resources := range scanResult.ByIssue() {
		for range resources {
			severitySum += int(issue.Severity)
			severityStats[issue.Severity.String()] += 1
		}
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
