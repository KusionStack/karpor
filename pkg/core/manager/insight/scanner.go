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

package insight

import (
	"context"
	"fmt"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/infra/scanner"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/pkg/errors"
)

// Audit performs the audit on Kubernetes manifests with the specified locator
// and returns the issues found during the audit.
func (i *InsightManager) Audit(ctx context.Context, locator core.Locator, noCache bool) (scanner.ScanResult, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting audit with specified condition in AuditManager ...")

	if noCache {
		log.Info("Scan without cache for locator", "locator", locator)
		return i.scanFor(ctx, locator, true)
	} else {
		if auditData, exist := i.scanCache.Get(locator); exist {
			log.Info("Cache hit for locator", "locator", locator)
			return auditData, nil
		} else {
			log.Info("Cache miss for locator", "locator", locator)
			return i.scanFor(ctx, locator, false)
		}
	}
}

// Score calculates a score based on the severity and total number of issues
// identified during the audit. It aggregates statistics on different severity
// levels and generates a cumulative score.
func (i *InsightManager) Score(ctx context.Context, locator core.Locator, noCache bool) (*ScoreData, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting calculate score with specified issues list in AuditManager ...")

	scanResult, err := i.Audit(ctx, locator, noCache)
	if err != nil {
		return nil, err
	}

	// Calculate the total score and severity statistics for each resource.
	var scoreTotal float64 = 0
	severityStats := map[string]int{}
	for _, issues := range scanResult.ByResource() {
		score, stats := CalculateResourceScore(issues)
		scoreTotal += score
		for k, v := range stats {
			severityStats[k] += v
		}
	}

	resourceTotal := len(scanResult.ByResource())
	if resourceTotal == 0 {
		scoreTotal = 100
	} else if resourceTotal > 0 {
		scoreTotal /= float64(resourceTotal)
	} else {
		return nil, fmt.Errorf("invalid resource total")
	}

	return &ScoreData{
		Score:             scoreTotal,
		ResourceTotal:     resourceTotal,
		IssuesTotal:       scanResult.IssueTotal(),
		SeverityStatistic: severityStats,
	}, nil
}

func (i *InsightManager) scanFor(ctx context.Context, locator core.Locator, noCache bool) (scanner.ScanResult, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)

	searchQuery := locator.ToSQL()
	searchPattern := storage.SQLPatternType
	pageSizeIteration := 100
	pageIteration := 1

	var result scanner.ScanResult
	for {
		log.Info("Starting search in AuditManager ...",
			"searchQuery", searchQuery, "searchPattern", searchPattern,
			"searchPageSize", pageSizeIteration, "searchPage", pageIteration)

		res, err := i.search.Search(ctx, searchQuery, searchPattern, &storage.Pagination{Page: pageSizeIteration, PageSize: pageIteration})
		if err != nil {
			return nil, err
		}

		log.Info("Finish current search", "overview", res.Overview())

		newResult, err := i.scanner.Scan(ctx, noCache, res.Resources...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan resources")
		}
		if result == nil {
			result = newResult
		} else {
			result.MergeFrom(newResult)
		}

		if len(res.Resources) < pageSizeIteration {
			break
		}
		pageIteration++
	}

	i.scanCache.Set(locator, result)
	log.Info("Added data to cache for locator", "locator", locator)

	return result, nil
}
