package kubeaudit

import (
	"sync"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/KusionStack/karbour/pkg/search/storage"
)

var _ scanner.ScanResult = &scanResult{}

// scanResult implements the scanner.ScanResult interface and represents the
// result of scanning Kubernetes resources.
type scanResult struct {
	issueResourceMap  map[scanner.Issue]scanner.ResourceList // Map of issues to resources
	resourceIssueMap  map[core.Locator]scanner.IssueList     // Map of resources to issues
	locatorMap        map[core.Locator]*storage.Resource     // Map of locators to resources
	relationshipExist map[relationship]struct{}              // Map to track relationships
	lock              sync.RWMutex                           // Mutex for concurrent access
}

// relationship represents the relationship between an issue and a locator.
type relationship struct {
	scanner.Issue
	core.Locator
}

// NewScanResult creates a new instance of scanResult.
func NewScanResult() scanner.ScanResult {
	return newScanResult()
}

// newScanResult creates and returns a new instance of scanResult.
func newScanResult() *scanResult {
	return &scanResult{
		issueResourceMap:  make(map[scanner.Issue]scanner.ResourceList),
		resourceIssueMap:  make(map[core.Locator]scanner.IssueList),
		locatorMap:        make(map[core.Locator]*storage.Resource),
		relationshipExist: map[relationship]struct{}{},
		lock:              sync.RWMutex{},
	}
}

// ByIssue returns the map of issues to resources.
func (sr *scanResult) ByIssue() map[scanner.Issue]scanner.ResourceList {
	return sr.issueResourceMap
}

// ByResource returns the map of resources to issues.
func (sr *scanResult) ByResource() map[core.Locator]scanner.IssueList {
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

	for locator, issues := range newResult.resourceIssueMap {
		if resource, exist := newResult.locatorMap[locator]; exist {
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

	locator := core.LocatorFor(resource)

	if _, exist := sr.locatorMap[locator]; !exist {
		sr.locatorMap[locator] = resource
	}

	if _, ok := sr.resourceIssueMap[locator]; !ok {
		sr.resourceIssueMap[locator] = make([]*scanner.Issue, 0)
	}

	for _, issue := range issues {
		if _, ok := sr.issueResourceMap[*issue]; !ok {
			sr.issueResourceMap[*issue] = make(scanner.ResourceList, 0)
		}
	}

	for i := 0; i < len(issues); i++ {
		issue := issues[i]

		rel := relationship{
			Issue:   *issue,
			Locator: locator,
		}

		if _, exist := sr.relationshipExist[rel]; exist {
			continue
		}

		sr.resourceIssueMap[locator] = append(sr.resourceIssueMap[locator], issue)
		sr.issueResourceMap[*issue] = append(sr.issueResourceMap[*issue], resource)
		sr.relationshipExist[rel] = struct{}{}
	}
}
