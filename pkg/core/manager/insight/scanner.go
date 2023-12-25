package insight

import (
	"context"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/pkg/errors"
)

// Audit performs the audit on Kubernetes manifests with the specified locator
// and returns the issues found during the audit.
func (m *InsightManager) Audit(ctx context.Context, locator core.Locator) (scanner.ScanResult, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting audit with specified condition in AuditManager ...")

	searchQuery := locator.ToSQL()
	searchPattern := storage.SQLPatternType
	pageSizeIteration := 100
	pageIteration := 1

	if auditData, exist := m.scanCache.Get(locator); exist {
		log.Info("Cache hit for locator", "locator", locator)
		return auditData, nil
	} else {
		log.Info("Cache miss for locator", "locator", locator)

		var result scanner.ScanResult
		for {
			log.Info("Starting search in AuditManager ...",
				"searchQuery", searchQuery, "searchPattern", searchPattern,
				"searchPageSize", pageSizeIteration, "searchPage", pageIteration)

			res, err := m.search.Search(ctx, searchQuery, searchPattern, pageSizeIteration, pageIteration)
			if err != nil {
				return nil, err
			}

			log.Info("Finish current search", "overview", res.Overview())

			newResult, err := m.scanner.Scan(ctx, res.Resources...)
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

		m.scanCache.Set(locator, result)
		log.Info("Added data to cache for locator", "locator", locator)

		return result, nil
	}
}

// Score calculates a score based on the severity and total number of issues
// identified during the audit. It aggregates statistics on different severity
// levels and generates a cumulative score.
func (m *InsightManager) Score(ctx context.Context, locator core.Locator) (*ScoreData, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting calculate score with specified issues list in AuditManager ...")

	scanResult, err := m.Audit(ctx, locator)
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
	if len(scanResult.ByResource()) > 0 {
		scoreTotal /= float64(len(scanResult.ByResource()))
	}

	// Prepare the score data including the total, sum and statistics.
	data := &ScoreData{
		Score:             scoreTotal,
		ResourceTotal:     len(scanResult.ByResource()),
		IssuesTotal:       scanResult.IssueTotal(),
		SeverityStatistic: severityStats,
	}

	return data, nil
}
