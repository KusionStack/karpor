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

	clustermanager "github.com/KusionStack/karbour/pkg/core/manager/cluster"
	"github.com/KusionStack/karbour/pkg/infra/multicluster"
)

// Statistics is a method of the InsightManager struct which provides statistical information.
//
// Parameters:
// - ctx (context.Context): The context object for managing the lifecycle of the request.
//
// Returns:
// - *Statistics: A pointer to a Statistics struct containing the aggregated statistics.
// - error: An error if one occurred during the retrieval of statistics.
func (i *InsightManager) Statistics(ctx context.Context) (*Statistics, error) {
	// Get count of resources.
	resourceCount, err := i.resource.CountResources(ctx)
	if err != nil {
		return nil, err
	}
	// Get count of resource group rules.
	rgrCount, err := i.resourceGroupRule.CountResourceGroupRules(ctx)
	if err != nil {
		return nil, err
	}

	// Get count of clusters.
	clusterMgr := clustermanager.NewClusterManager()
	client, err := multicluster.BuildMultiClusterClient(ctx, i.genericConfig.LoopbackClientConfig, "")
	if err != nil {
		return nil, err
	}

	summary, err := clusterMgr.CountCluster(ctx, client, i.genericConfig.LoopbackClientConfig)
	if err != nil {
		return nil, err
	}

	return &Statistics{
		ClusterCount:           summary.TotalCount,
		ResourceCount:          resourceCount,
		ResourceGroupRuleCount: rgrCount,
	}, nil
}
