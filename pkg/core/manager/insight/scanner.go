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

package insight

import (
	"context"
	"fmt"

	"kusionstack.io/karpor/pkg/core/entity"
	"kusionstack.io/karpor/pkg/infra/scanner"
	"kusionstack.io/karpor/pkg/infra/search/storage"
	"kusionstack.io/karpor/pkg/util/ctxutil"
	"github.com/pkg/errors"
)

// Audit performs the audit on Kubernetes manifests with the specified resourceGroup
// and returns the issues found during the audit.
func (i *InsightManager) Audit(ctx context.Context, resourceGroup entity.ResourceGroup, noCache bool) (scanner.ScanResult, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting audit with specified condition in AuditManager ...")

	if noCache {
		log.Info("Scan without cache for resourceGroup", "resourceGroup", resourceGroup)
		return i.scanFor(ctx, resourceGroup, true)
	} else {
		if auditData, exist := i.scanCache.Get(resourceGroup.Hash()); exist {
			log.Info("Cache hit for resourceGroup", "resourceGroup", resourceGroup)
			return auditData, nil
		} else {
			log.Info("Cache miss for resourceGroup", "resourceGroup", resourceGroup)
			return i.scanFor(ctx, resourceGroup, false)
		}
	}
}

// Score calculates a score based on the severity and total number of issues
// identified during the audit. It aggregates statistics on different severity
// levels and generates a cumulative score.
func (i *InsightManager) Score(ctx context.Context, resourceGroup entity.ResourceGroup, noCache bool) (*ScoreData, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting calculate score with specified issues list in AuditManager ...")

	scanResult, err := i.Audit(ctx, resourceGroup, noCache)
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

// scanFor is a helper function that performs the actual scanning for the given resourceGroup using the context and returns the scan result and error.
func (i *InsightManager) scanFor(ctx context.Context, resourceGroup entity.ResourceGroup, noCache bool) (scanner.ScanResult, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)

	searchQuery := resourceGroup.ToSQL()
	searchPattern := storage.SQLPatternType
	pageSizeIteration := 100
	pageIteration := 1

	var result scanner.ScanResult
	for {
		log.Info("Starting search in AuditManager ...",
			"searchQuery", searchQuery, "searchPattern", searchPattern,
			"searchPageSize", pageSizeIteration, "searchPage", pageIteration)

		res, err := i.search.Search(ctx, searchQuery, searchPattern, &storage.Pagination{Page: pageIteration, PageSize: pageSizeIteration})
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

	i.scanCache.Set(resourceGroup.Hash(), result)
	log.Info("Added data to cache for resourceGroup", "resourceGroup", resourceGroup)

	return result, nil
}
