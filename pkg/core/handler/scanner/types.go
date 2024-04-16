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

package scanner

import (
	"github.com/KusionStack/karbour/pkg/core/entity"
	"github.com/KusionStack/karbour/pkg/infra/scanner"
)

// AuditData represents the aggregated data of scanner issues, including the
// original list of issues and their aggregated count based on title.
type AuditData struct {
	IssueTotal    int            `json:"issueTotal"`
	ResourceTotal int            `json:"resourceTotal"`
	BySeverity    map[string]int `json:"bySeverity"`
	IssueGroups   []*IssueGroup  `json:"issueGroups"`
}

// IssueGroup represents a group of resourceGroups tied to a specific issue.
type IssueGroup struct {
	Issue    scanner.Issue        `json:"issue"`
	ResourceGroups []entity.ResourceGroup `json:"resourceGroups"`
}
