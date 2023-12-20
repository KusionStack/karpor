package audit

import (
	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/scanner"
)

func convertScanResultToAuditData(sr scanner.ScanResult) *AuditData {
	issueGroups := make([]*IssueGroup, 0, len(sr.ByIssue()))
	bySeverity := map[string]int{}

	for issue, resources := range sr.ByIssue() {
		issueGroup := &IssueGroup{
			Issue:    issue,
			Locators: []*core.Locator{},
		}
		for _, resource := range resources {
			locator := core.LocatorFor(*resource)
			issueGroup.Locators = append(issueGroup.Locators, locator)
			bySeverity[issue.Severity.String()] += 1
		}
		issueGroups = append(issueGroups, issueGroup)
	}

	auditData := &AuditData{
		IssueTotal:    sr.IssueTotal(),
		ResourceTotal: len(sr.ByResource()),
		BySeverity:    bySeverity,
		IssueGroups:   issueGroups,
	}

	return auditData
}
