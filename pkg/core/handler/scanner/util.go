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

package scanner

import (
	"sort"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/core/manager/ai"
	"github.com/KusionStack/karpor/pkg/infra/scanner"
)

// convertScanResultToAuditData converts the scanner.ScanResult to an AuditData
// structure containing aggregated issue and resource data.
func convertScanResultToAuditData(sr scanner.ScanResult) *ai.AuditData {
	issueGroups := make([]*ai.IssueGroup, 0, len(sr.ByIssue()))
	bySeverity := map[string]int{}

	// Iterate through each issue in the ScanResult and create corresponding
	// IssueGroup entries.
	for issue, resources := range sr.ByIssue() {
		issueGroup := &ai.IssueGroup{
			Issue:          issue,
			ResourceGroups: []entity.ResourceGroup{},
		}

		// For each resource tied to the issue, create a ResourceGroup and increment
		// severity count.
		for _, resource := range resources {
			issueGroup.ResourceGroups = append(issueGroup.ResourceGroups, resource.ResourceGroup)
			bySeverity[issue.Severity.String()]++
		}
		issueGroups = append(issueGroups, issueGroup)
	}

	// Custom sorting function for IssueGroups
	sort.Slice(issueGroups, func(i, j int) bool {
		// First, sort by Severity from high to low.
		if issueGroups[i].Issue.Severity > issueGroups[j].Issue.Severity {
			return true
		} else if issueGroups[i].Issue.Severity < issueGroups[j].Issue.Severity {
			return false
		}

		// If Severities are equal, sort by ResourceGroups array size from high to
		// low.
		return len(issueGroups[i].ResourceGroups) > len(issueGroups[j].ResourceGroups)
	})

	// Construct the AuditData structure.
	return &ai.AuditData{
		IssueTotal:    sr.IssueTotal(),
		ResourceTotal: len(sr.ByResource()),
		BySeverity:    bySeverity,
		IssueGroups:   issueGroups,
	}
}
