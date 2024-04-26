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

package insight

import (
	"context"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/multicluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GetDetailsForCluster returns ClusterDetail object for a given cluster
func (i *InsightManager) GetDetailsForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (*ClusterDetail, error) {
	serverVersion, _ := client.ClientSet.DiscoveryClient.ServerVersion()
	// Get the list of nodes
	nodes, err := client.ClientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var memoryCapacity, cpuCapacity, podCapacity int64
	for _, node := range nodes.Items {
		memoryCapacity += node.Status.Capacity.Memory().Value()
		cpuCapacity += node.Status.Capacity.Cpu().Value()
		podCapacity += node.Status.Capacity.Pods().Value()
	}
	return &ClusterDetail{
		NodeCount:      len(nodes.Items),
		ServerVersion:  serverVersion.String(),
		MemoryCapacity: memoryCapacity,
		CPUCapacity:    cpuCapacity,
		PodsCapacity:   podCapacity,
	}, nil
}

// GetResourceSummary returns the unstructured cluster object summary for a given cluster. Possibly will add more metrics to it in the future.
func (i *InsightManager) GetResourceSummary(ctx context.Context, client *multicluster.MultiClusterClient, resourceGroup *entity.ResourceGroup) (*ResourceSummary, error) {
	obj, err := i.GetResource(ctx, client, resourceGroup)
	if err != nil {
		return nil, err
	}

	return &ResourceSummary{
		Resource: entity.ResourceGroup{
			Name:       obj.GetName(),
			Namespace:  obj.GetNamespace(),
			APIVersion: obj.GetAPIVersion(),
			Cluster:    resourceGroup.Cluster,
			Kind:       obj.GetKind(),
		},
		CreationTimestamp: obj.GetCreationTimestamp(),
		ResourceVersion:   obj.GetResourceVersion(),
		UID:               obj.GetUID(),
	}, nil
}

// GetGVKSummary returns the unstructured cluster object summary for a given GVK. Possibly will add more metrics to it in the future.
func (i *InsightManager) GetGVKSummary(ctx context.Context, client *multicluster.MultiClusterClient, resourceGroup *entity.ResourceGroup) (*GVKSummary, error) {
	gvkCount, err := i.CountResourcesByGVK(ctx, client, resourceGroup)
	if err != nil {
		return nil, err
	}
	gv, err := schema.ParseGroupVersion(resourceGroup.APIVersion)
	if err != nil {
		return nil, err
	}
	return &GVKSummary{
		Cluster: resourceGroup.Cluster,
		Group:   gv.Group,
		Version: gv.Version,
		Kind:    resourceGroup.Kind,
		Count:   gvkCount,
	}, nil
}

// GetNamespaceSummary returns the unstructured cluster object summary for a given namespace. Possibly will add more metrics to it in the future.
func (i *InsightManager) GetNamespaceSummary(ctx context.Context, client *multicluster.MultiClusterClient, resourceGroup *entity.ResourceGroup) (*NamespaceSummary, error) {
	namespaceCount, err := i.CountByResourceGroup(ctx, client, resourceGroup)
	if err != nil {
		return nil, err
	}
	topFiveCount := GetTopResultsFromMap(namespaceCount)
	return &NamespaceSummary{
		Cluster:    resourceGroup.Cluster,
		Namespace:  resourceGroup.Namespace,
		CountByGVK: topFiveCount,
	}, nil
}

// GetResourceGroupSummary returns a summary of a resource group, including details about its resources and their distribution.
func (i *InsightManager) GetResourceGroupSummary(ctx context.Context, client *multicluster.MultiClusterClient, resourceGroup *entity.ResourceGroup) (*ResourceGroupSummary, error) {
	count, err := i.CountByResourceGroup(ctx, client, resourceGroup)
	if err != nil {
		return nil, err
	}
	topFiveCount := GetTopResultsFromMap(count)
	return &ResourceGroupSummary{
		ResourceGroup: resourceGroup,
		CountByGVK:    topFiveCount,
	}, nil
}
