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

// Package scanner provides an interface and primitives for implementing scanners
// that check Kubernetes resources for various issues, such as security vulnerabilities,
// misconfigurations, and best practices.
package scanner

import (
	"io"

	"k8s.io/apimachinery/pkg/runtime"
)

// IssueSeverityLevel defines the severity level of an issue.
// It is an enumeration starting from 0 (Low) and increases with severity.
const (
	Low      IssueSeverityLevel = iota // Low indicates a minor issue that should be addressed.
	Medium                             // Medium indicates a potential issue that may have a moderate impact.
	High                               // High indicates a serious issue that has a significant impact.
	Critical                           // Critical indicates an extremely serious issue that must be addressed immediately.
)

// KubeScanner is an interface for scanners that analyze Kubernetes resources.
// Each scanner should implement this interface to provide scanning functionality.
type KubeScanner interface {
	Name() string                                       // Name returns the name of the scanner.
	Scan(resources ...runtime.Object) ([]*Issue, error) // Scan accepts one or more Kubernetes resources and returns a slice of issues found.
	ScanManifest(manifest io.Reader) ([]*Issue, error)  // Scan accepts a Kubernetes manifest and returns a slice of issues found.
}

// Issue represents a particular finding or problem discovered by a scanner.
// It encapsulates the details of the issue such as the scanner's name, its severity,
// and a human-readable title and message.
type Issue struct {
	Scanner  string             // Scanner is the name of the scanner that discovered the issue.
	Severity IssueSeverityLevel // Severity indicates how critical the issue is, using the IssueSeverityLevel constants.
	Title    string             // Title is a brief summary of the issue.
	Message  string             // Message provides a detailed human-readable description of the issue.
}

// IssueSeverityLevel represents the severity level of an issue.
type IssueSeverityLevel int

// String returns the string representation of the IssueSeverityLevel.
func (s IssueSeverityLevel) String() string {
	switch s {
	case Low:
		return "Low" // Represents low severity level.
	case Medium:
		return "Medium" // Represents medium severity level.
	case High:
		return "High" // Represents high severity level.
	case Critical:
		return "Critical" // Represents critical severity level.
	default:
		return "Unknown" // Indicates an unknown severity level.
	}
}
