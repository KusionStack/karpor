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
	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/scanner"
)

// convertScanResultToAuditData converts the scanner.ScanResult to an AuditData
// structure containing aggregated issue and resource data.
func convertScanResultToAuditData(sr scanner.ScanResult) *AuditData {
	issueGroups := make([]*IssueGroup, 0, len(sr.ByIssue()))
	bySeverity := map[string]int{}

	// Iterate through each issue in the ScanResult and create corresponding
	// IssueGroup entries.
	for issue, resources := range sr.ByIssue() {
		issueGroup := &IssueGroup{
			Issue:    issue,
			Locators: []core.Locator{},
		}

		// For each resource tied to the issue, create a Locator and increment
		// severity count.
		for _, resource := range resources {
			locator := core.LocatorFor(resource)
			issueGroup.Locators = append(issueGroup.Locators, locator)
			bySeverity[issue.Severity.String()]++
		}
		issueGroups = append(issueGroups, issueGroup)
	}

	// Construct the AuditData structure.
	auditData := &AuditData{
		IssueTotal:    sr.IssueTotal(),
		ResourceTotal: len(sr.ByResource()),
		BySeverity:    bySeverity,
		IssueGroups:   issueGroups,
	}

	return auditData
}
