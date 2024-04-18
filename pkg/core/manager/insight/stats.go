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
