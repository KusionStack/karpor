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

package kubeaudit

import (
	"sync"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/scanner"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
)

var _ scanner.ScanResult = &scanResult{}

// scanResult implements the scanner.ScanResult interface and represents the
// result of scanning Kubernetes resources.
type scanResult struct {
	issueResourceMap  map[scanner.Issue]scanner.ResourceList         // Map of issues to resources
	resourceIssueMap  map[entity.ResourceGroupHash]scanner.IssueList // Map of resources to issues
	resourceGroupMap  map[entity.ResourceGroupHash]*storage.Resource // Map of resourceGroup to resources
	relationshipExist map[relationship]struct{}                      // Map to track relationships
	lock              sync.RWMutex                                   // Mutex for concurrent access
}

// relationship represents the relationship between an issue and a resourceGroups.
type relationship struct {
	scanner.Issue
	entity.ResourceGroupHash
}

// NewScanResult creates a new instance of scanResult.
func NewScanResult() scanner.ScanResult {
	return newScanResult()
}

// newScanResult creates and returns a new instance of scanResult.
func newScanResult() *scanResult {
	return &scanResult{
		issueResourceMap:  make(map[scanner.Issue]scanner.ResourceList),
		resourceIssueMap:  make(map[entity.ResourceGroupHash]scanner.IssueList),
		resourceGroupMap:  make(map[entity.ResourceGroupHash]*storage.Resource),
		relationshipExist: map[relationship]struct{}{},
		lock:              sync.RWMutex{},
	}
}

// ByIssue returns the map of issues to resources.
func (sr *scanResult) ByIssue() map[scanner.Issue]scanner.ResourceList {
	return sr.issueResourceMap
}

// ByResource returns the map of resources to issues.
func (sr *scanResult) ByResource() map[entity.ResourceGroupHash]scanner.IssueList {
	return sr.resourceIssueMap
}

// IssueTotal calculates the total number of issues.
func (sr *scanResult) IssueTotal() int {
	sr.lock.RLock()
	defer sr.lock.RUnlock()

	totalByIssueMap, totalByResourceMap := 0, 0

	// Calculate total issues by scanning issueResourceMap
	for _, resources := range sr.issueResourceMap {
		totalByIssueMap += len(resources)
	}

	// Calculate total issues by scanning resourceIssueMap
	for _, issues := range sr.resourceIssueMap {
		totalByResourceMap += len(issues)
	}

	// Compare the totals and return -1 if they differ, otherwise return the total
	if totalByIssueMap != totalByResourceMap {
		return -1
	}

	return totalByIssueMap
}

// MergeFrom merges the results from another scanResult into the current one.
func (sr *scanResult) MergeFrom(result scanner.ScanResult) {
	if result == nil {
		return
	}

	newResult, ok := result.(*scanResult)
	if !ok {
		panic("Invalid type for merging")
	}

	if sr == nil {
		sr = newScanResult()
	}

	for resourceGroup, issues := range newResult.resourceIssueMap {
		if resource, exist := newResult.resourceGroupMap[resourceGroup]; exist {
			sr.add(resource, issues)
		}
	}
}

// add adds a resource with its associated issues to the scanResult.
func (sr *scanResult) add(resource *storage.Resource, issues []*scanner.Issue) {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	if resource == nil {
		return
	}

	if len(issues) == 0 {
		issues = make([]*scanner.Issue, 0)
	}

	if _, exist := sr.resourceGroupMap[resource.ResourceGroup.Hash()]; !exist {
		sr.resourceGroupMap[resource.ResourceGroup.Hash()] = resource
	}

	if _, ok := sr.resourceIssueMap[resource.ResourceGroup.Hash()]; !ok {
		sr.resourceIssueMap[resource.ResourceGroup.Hash()] = make([]*scanner.Issue, 0)
	}

	for _, issue := range issues {
		if _, ok := sr.issueResourceMap[*issue]; !ok {
			sr.issueResourceMap[*issue] = make(scanner.ResourceList, 0)
		}
	}

	for i := 0; i < len(issues); i++ {
		issue := issues[i]

		rel := relationship{
			Issue:             *issue,
			ResourceGroupHash: resource.ResourceGroup.Hash(),
		}

		if _, exist := sr.relationshipExist[rel]; exist {
			continue
		}

		sr.resourceIssueMap[resource.ResourceGroup.Hash()] = append(sr.resourceIssueMap[resource.ResourceGroup.Hash()], issue)
		sr.issueResourceMap[*issue] = append(sr.issueResourceMap[*issue], resource)
		sr.relationshipExist[rel] = struct{}{}
	}
}
