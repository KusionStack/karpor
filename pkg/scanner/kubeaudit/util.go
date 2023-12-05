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
	"github.com/KusionStack/karbour/pkg/scanner"
	kubeauditpkg "github.com/elliotxx/kubeaudit"
)

// AuditResult2Issue converts a kubeaudit.AuditResult to a scanner.Issue.
func AuditResult2Issue(auditResult *kubeauditpkg.AuditResult) *scanner.Issue {
	return &scanner.Issue{
		Scanner:  ScannerName,                                      // The name of the scanner that identified the issue.
		Severity: scanner.IssueSeverityLevel(auditResult.Severity), // The severity of the issue based on the audit result.
		Title:    auditResult.Rule,                                 // The rule that was violated, serving as the issue's title.
		Message:  auditResult.Message,                              // A human-readable message describing the issue.
	}
}
