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

package insight

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/infra/multicluster"
	"github.com/KusionStack/karbour/pkg/infra/scanner"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	topologyutil "github.com/KusionStack/karbour/pkg/util/topology"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CalculateResourceScore calculates the resource score and severity statistics
// based on the provided issues.
func CalculateResourceScore(issues scanner.IssueList) (float64, map[string]int) {
	severityStats := map[string]int{}
	issueTotal, severitySum := len(issues), 0
	for _, issue := range issues {
		severitySum += int(issue.Severity)
		severityStats[issue.Severity.String()] += 1
	}
	return CalculateScore(issueTotal, severitySum), severityStats
}

// CalculateScore calculates the score based on the number of issues and their
// severity sum (in the range of 1-5).
// P is the number of issues, and S is the sum of the severity (range 1-5) of
// the issue S will not be less than P.
//
// Example:
// - When there is one high-level issue, P=1 and S=3.
// - When there are three high-level issues, P=3 and S=9.
// - When there are ten low-level issues, P=10 and S=10.
func CalculateScore(p, s int) float64 {
	a, b := -0.04, -0.06
	param := a*float64(p) + b*float64(s)
	return 100 * math.Exp(param)
}

// CountResourcesByGVK returns an int that corresponds to the count of a resource GVK defined using core.Locator
func (i *InsightManager) CountResourcesByGVK(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) (int, error) {
	if loc.Cluster == "" || loc.APIVersion == "" || loc.Kind == "" {
		return 0, fmt.Errorf("cluster, APIVersion and Kind in locator cannot be empty")
	}
	resourceGVR, err := topologyutil.GetGVRFromGVK(loc.APIVersion, loc.Kind)
	if err != nil {
		return 0, err
	}
	resList, err := client.DynamicClient.Resource(resourceGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(resList.Items), nil
}

// CountResourcesByGVK returns a map from string to int
func (i *InsightManager) CountResourcesByNamespace(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) (map[string]int, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	if loc.Cluster == "" || loc.Namespace == "" {
		return nil, fmt.Errorf("cluster and Namespace in locator cannot be empty")
	}
	counts := make(map[string]int)
	// Another option here is to retrieve the list of API resources and iterate over each using dynamic client
	// That will create more pressure on the spoke cluster and cause unexpected and unnecessary amount of time
	// We opted to use Elastic search as the source of the count
	searchQuery := loc.ToSQL()
	searchPattern := storage.SQLPatternType
	pageSizeIteration := 100
	pageIteration := 1

	log.Info("Starting search in InsightManager ...",
		"searchQuery", searchQuery, "searchPattern", searchPattern, "searchPageSize", pageSizeIteration, "searchPage", pageIteration)

	for {
		res, err := i.search.Search(ctx, searchQuery, searchPattern, pageSizeIteration, pageIteration)
		if err != nil {
			return nil, err
		}
		log.Info("Finish current search", "overview", res.Overview())

		for _, resource := range res.Resources {
			gvk := fmt.Sprintf("%s.%s", resource.Kind, resource.APIVersion)
			counts[gvk]++
		}
		if len(res.Resources) < pageSizeIteration {
			break
		}
		pageIteration++
	}

	return counts, nil
}

// GetTopResultsFromMap returns the top 5 results from the map sorted by value
// If map does not have 5 elements, return the full map
func GetTopResultsFromMap(m map[string]int) map[string]int {
	s := make([]KeyValuePair, 0)
	res := make(map[string]int, 0)

	for k, v := range m {
		s = append(s, KeyValuePair{k, v})
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].value > s[j].value
	})

	index := min(len(s), 5)
	for _, kv := range s[:index] {
		res[kv.key] = kv.value
	}

	return res
}

// min returns the smaller of two integers x and y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
