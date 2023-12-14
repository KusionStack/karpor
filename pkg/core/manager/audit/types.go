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

import "github.com/KusionStack/karbour/pkg/scanner"

// ScoreData encapsulates the results of scoring an audited manifest. It provides
// a numerical score along with statistics about the total number of issues and
// their severities.
type ScoreData struct {
	// Score represents the calculated score of the audited manifest based on
	// the number and severity of issues. It provides a quantitative measure
	// of the security posture of the resources in the manifest.
	Score float64 `json:"score"`

	// IssuesTotal is the total count of all issues found during the audit.
	// This count can be used to understand the overall number of problems
	// that need to be addressed.
	IssuesTotal int `json:"issuesTotal"`

	// SeveritySum is the sum of severity scores of all issues, which can be
	// used to gauge the cumulative severity of all problems found.
	SeveritySum int `json:"severitySum"`

	// SeverityStatistic is a mapping of severity levels to their respective
	// number of occurrences. It allows for a quick overview of the distribution
	// of issues across different severity categories.
	SeverityStatistic map[string]int `json:"severityStatistic"`
}

// AuditData represents the aggregated data of scanner issues, including the
// original list of issues and their aggregated count based on title.
type AuditData struct {
	ResourcesTotal int              `json:"resourcesTotal"`
	Aggregated     map[string]int   `json:"aggregated"`
	Issues         []*scanner.Issue `json:"issues"`
}

// NewAuditData initializes an AuditData instance by aggregating the counts of
// each issue's title from the provided list of issues.
func NewAuditData(issues []*scanner.Issue, total int) *AuditData {
	data := &AuditData{
		ResourcesTotal: total,
		Issues:         issues,
		Aggregated:     make(map[string]int),
	}

	// Aggregate counts of each issue's title
	for _, issue := range data.Issues {
		data.Aggregated[issue.Title]++
	}

	return data
}
