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
