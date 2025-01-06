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
	"fmt"
	"time"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/multicluster"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	PodStatusRunning    = "Running"
	PodStatusTerminated = "Terminated"
	PodStatusUnknown    = "Unknown"
	PodStatusWaiting    = "Waiting"
)

// GetDetailsForCluster returns ClusterDetail object for a given cluster
func (i *InsightManager) GetDetailsForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (*ClusterDetail, error) {
	// get server version and measure latency
	start := time.Now()
	serverVersionStr := "v0.0.0-unknown" // default version when error occurs
	serverVersion, err := client.ClientSet.DiscoveryClient.ServerVersion()
	if err == nil {
		serverVersionStr = serverVersion.String()
	}
	latency := time.Since(start).Milliseconds()

	// Get the list of nodes
	nodes, err := client.ClientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	var (
		memoryCapacity int64
		cpuCapacity    int64
		podCapacity    int64
		memoryUsage    float64
		cpuUsage       float64
		readyNodes     int
		notReadyNodes  int
		metricsEnabled bool
		cpuMetrics     ResourceMetrics
		memoryMetrics  ResourceMetrics
	)

	// Calculate capacities and node status
	for _, node := range nodes.Items {
		memoryCapacity += node.Status.Capacity.Memory().Value()
		cpuCapacity += node.Status.Capacity.Cpu().MilliValue()
		podCapacity += node.Status.Capacity.Pods().Value()

		// Count ready/not ready nodes
		isReady := false
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				if condition.Status == corev1.ConditionTrue {
					readyNodes++
					isReady = true
				}
				break
			}
		}
		if !isReady {
			notReadyNodes++
		}
	}

	// Get metrics if available
	if client.MetricsClient != nil {
		// Get current metrics
		nodeMetrics, err := client.MetricsClient.NodeMetricses().List(ctx, metav1.ListOptions{})
		if err == nil {
			metricsEnabled = true
			// Calculate current usage
			for _, metric := range nodeMetrics.Items {
				cpuUsage += float64(metric.Usage.Cpu().MilliValue())
				memoryUsage += float64(metric.Usage.Memory().Value())
			}

			// Get historical metrics
			now := metav1.Now()
			oneHourAgo := metav1.NewTime(now.Add(-1 * time.Hour))

			// Collect data points every 5 minutes, total 12 points
			interval := 5 * time.Minute
			for t := oneHourAgo.Time; t.Before(now.Time); t = t.Add(interval) {
				metrics, err := client.MetricsClient.NodeMetricses().List(ctx, metav1.ListOptions{
					TimeoutSeconds: &[]int64{5}[0], // Set timeout to 5 seconds
				})
				if err == nil {
					var cpuTotal, memoryTotal float64
					for _, metric := range metrics.Items {
						cpuTotal += float64(metric.Usage.Cpu().MilliValue())
						memoryTotal += float64(metric.Usage.Memory().Value())
					}
					timestamp := metav1.NewTime(t)
					cpuMetrics.Points = append(cpuMetrics.Points, MetricPoint{
						Timestamp: timestamp,
						Value:     cpuTotal,
					})
					memoryMetrics.Points = append(memoryMetrics.Points, MetricPoint{
						Timestamp: timestamp,
						Value:     memoryTotal,
					})
				}
			}
		}
	}

	// Get current pods count
	pods, err := client.ClientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	return &ClusterDetail{
		NodeCount:      len(nodes.Items),
		ServerVersion:  serverVersionStr,
		MemoryCapacity: memoryCapacity,
		MemoryUsage:    memoryUsage,
		CPUCapacity:    cpuCapacity,
		CPUUsage:       cpuUsage,
		PodsCapacity:   podCapacity,
		PodsUsage:      int64(len(pods.Items)),
		ReadyNodes:     readyNodes,
		NotReadyNodes:  notReadyNodes,
		MetricsEnabled: metricsEnabled,
		CPUMetrics:     cpuMetrics,
		MemoryMetrics:  memoryMetrics,
		Latency:        latency,
	}, nil
}

// GetResourceSummary returns the unstructured cluster object summary for a given cluster. Possibly will add more metrics to it in the future.
func (i *InsightManager) GetResourceSummary(ctx context.Context, client *multicluster.MultiClusterClient, resourceGroup *entity.ResourceGroup) (*ResourceSummary, error) {
	obj, err := i.GetResource(ctx, client, resourceGroup)
	if err != nil {
		return nil, err
	}

	Status := ""
	if obj.GetKind() == "Pod" {
		Status = GetPodStatus(obj.Object)
	}

	return &ResourceSummary{
		Resource: entity.ResourceGroup{
			Name:       obj.GetName(),
			Namespace:  obj.GetNamespace(),
			APIVersion: obj.GetAPIVersion(),
			Cluster:    resourceGroup.Cluster,
			Kind:       obj.GetKind(),
			Status:     Status,
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

// GetPodStatus returns the status of a pod
func GetPodStatus(obj map[string]any) string {
	containerStatuses, found, err := unstructured.NestedSlice(obj, "status", "containerStatuses")
	if err != nil || !found || len(containerStatuses) == 0 {
		return PodStatusUnknown
	}
	firstContainer, ok := containerStatuses[0].(map[string]any)
	if !ok {
		return PodStatusUnknown
	}
	state, found := firstContainer["state"]
	if !found {
		return PodStatusUnknown
	}
	stateMap, ok := state.(map[string]interface{})
	if !ok {
		return PodStatusUnknown
	}
	if stateMap["running"] != nil {
		return PodStatusRunning
	}
	if stateMap["waiting"] != nil {
		waitMap, ok := stateMap["waiting"].(map[string]any)
		if !ok {
			return PodStatusWaiting
		}
		if reason, ok := waitMap["reason"].(string); ok && reason != "" {
			return reason
		}
		return PodStatusWaiting
	}
	if stateMap["terminated"] != nil {
		return PodStatusTerminated
	}
	return PodStatusUnknown
}
