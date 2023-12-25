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

package kubeaudit

import (
	"github.com/KusionStack/karbour/pkg/infra/scanner"
	kubeauditpkg "github.com/elliotxx/kubeaudit"
)

// AuditResult2Issue converts a kubeaudit.AuditResult to a scanner.Issue,
// which can be used to report security findings in a standardized format.
func AuditResult2Issue(auditResult *kubeauditpkg.AuditResult) *scanner.Issue {
	return &scanner.Issue{
		// Scanner is the name of the scanner that identified the issue.
		Scanner: ScannerName,

		// Severity represents the severity level of the issue as determined by
		// the audit result.
		Severity: ConvertSeverity(auditResult.Severity),

		// Title is the rule that was violated, which serves as a concise
		// description of the issue.
		Title: auditResult.Rule,

		// Message provides a detailed, human-readable description of the issue.
		Message: auditResult.Message,
	}
}

// ConvertSeverity translates a kubeaudit.SeverityLevel into a
// scanner.IssueSeverityLevel, which standardizes severity levels across
// different scanners.
func ConvertSeverity(level kubeauditpkg.SeverityLevel) scanner.IssueSeverityLevel {
	switch level {
	case kubeauditpkg.Warn:
		// Low severity corresponds to warnings in kubeaudit findings.
		return scanner.Low
	case kubeauditpkg.Error:
		// High severity corresponds to errors in kubeaudit findings.
		return scanner.High
	default:
		// Safe represents no security risk or an informational finding.
		return scanner.Safe
	}
}
