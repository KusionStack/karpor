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

// Package scanner provides an interface and primitives for implementing scanners
// that check Kubernetes resources for various issues, such as security vulnerabilities,
// misconfigurations, and best practices.
package scanner

import (
	"context"
	"encoding/json"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
)

// IssueSeverityLevel defines the severity levels for issues identified by
// scanners.
const (
	Safe     IssueSeverityLevel = 0 // Safe indicates the absence of any security risk or an informational finding that does not require action.
	Low      IssueSeverityLevel = 1 // Low indicates a minor issue that should be addressed.
	Medium   IssueSeverityLevel = 2 // Medium indicates a potential issue that may have a moderate impact.
	High     IssueSeverityLevel = 3 // High indicates a serious issue that has a significant impact.
	Critical IssueSeverityLevel = 5 // Critical indicates an extremely serious issue that must be addressed immediately.
)

// KubeScanner is an interface for scanners that analyze Kubernetes resources.
// Each scanner should implement this interface to provide scanning functionality.
type KubeScanner interface {
	Name() string                                                                               // Name returns the name of the scanner.
	Scan(ctx context.Context, noCache bool, resources ...*storage.Resource) (ScanResult, error) // Scan accepts one or more Kubernetes resources and returns a slice of issues found.
}

// ScanResult defines the interface for the result of a scan.
type ScanResult interface {
	ByIssue() map[Issue]ResourceList
	ByResource() map[entity.ResourceGroupHash]IssueList
	IssueTotal() int
	MergeFrom(result ScanResult)
}

type (
	// ResourceList is a slice of storage resources.
	ResourceList []*storage.Resource
	// IssueList is a slice of issues.
	IssueList []*Issue
)

// Issue represents a particular finding or problem discovered by a scanner.
// It encapsulates the details of the issue such as the scanner's name, its severity,
// and a human-readable title and message.
type Issue struct {
	// Scanner is the name of the scanner that discovered the issue.
	Scanner string `json:"scanner" yaml:"scanner"`
	// Severity indicates how critical the issue is, using the IssueSeverityLevel constants.
	Severity IssueSeverityLevel `json:"severity" yaml:"severity"`
	// Title is a brief summary of the issue.
	Title string `json:"title" yaml:"title"`
	// Message provides a detailed human-readable description of the issue.
	Message string `json:"message" yaml:"message"`
}

// IssueSeverityLevel represents the severity level of an issue.
type IssueSeverityLevel int

// String returns the string representation of the IssueSeverityLevel.
func (s IssueSeverityLevel) String() string {
	switch s {
	case Safe:
		return "Safe"
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// MarshalJSON implements the json.Marshaler interface for IssueSeverityLevel.
func (s IssueSeverityLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
