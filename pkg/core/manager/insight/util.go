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
	"math"

	"context"
	"fmt"

	"github.com/KusionStack/karbour/pkg/scanner"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/multicluster"
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
func CountResourcesByGVK(ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator) (int, error) {
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
	// log := ctxutil.GetLogger(ctx)

	if loc.Cluster == "" || loc.Namespace == "" {
		return nil, fmt.Errorf("cluster and Namespace in locator cannot be empty")
	}
	counts := make(map[string]int)
	// Retrieve the list of API resources
	// apiResourceList, err := client.ClientSet.Discovery().ServerPreferredResources()
	// if err != nil {
	// 	return nil, err
	// }
	// for _, list := range apiResourceList {
	// 	groupVersion, err := schema.ParseGroupVersion(list.GroupVersion)
	// 	if err != nil {
	// 		fmt.Printf("Error parsing group version: %s\n", err.Error())
	// 		continue
	// 	}

	// 	for _, resource := range list.APIResources {
	// 		gvk := groupVersion.WithKind(resource.Kind)

	// 		// Skip resources that are not namespaced
	// 		if !resource.Namespaced {
	// 			continue
	// 		}

	// 		// Retrieve the list of resources for the given GVK
	// 		gvr := schema.GroupVersionResource{Group: gvk.Group, Version: gvk.Version, Resource: resource.Name}
	// 		timeoutSeconds := int64(1)
	// 		unstructuredList, err := client.DynamicClient.Resource(gvr).Namespace(loc.Namespace).List(context.TODO(), metav1.ListOptions{
	// 			TimeoutSeconds: &timeoutSeconds,
	// 		})
	// 		if err != nil {
	// 			fmt.Printf("Error listing %s: %s\n", gvk.String(), err.Error())
	// 			continue
	// 		}

	// 		// Count the resources
	// 		counts[resource.Kind] += len(unstructuredList.Items)
	// 	}
	// }

	// searchQuery := loc.ToSQL()
	// searchPattern := storage.SQLPatternType
	// pageSizeIteration := 100
	// pageIteration := 1

	// for {
	// 	log.Info("Starting search in AuditManager ...",
	// 		"searchQuery", searchQuery, "searchPattern", searchPattern, "searchPageSize", pageSizeIteration, "searchPage", pageIteration)

	// 	res, err := i.Search(ctx, searchQuery, searchPattern, pageSizeIteration, pageIteration)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	log.Info("Finish current search", "overview", res.Overview())

	// 	newResult, err := m.ks.Scan(ctx, res.Resources...)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "failed to scan resources")
	// 	}
	// 	if result == nil {
	// 		result = newResult
	// 	} else {
	// 		result.MergeFrom(newResult)
	// 	}

	// 	if len(res.Resources) < pageSizeIteration {
	// 		break
	// 	}
	// 	pageIteration++
	// }

	return counts, nil
}
