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

package cluster

import (
	"context"
	"os"

	clusterv1beta1 "github.com/KusionStack/karbour/pkg/apis/cluster/v1beta1"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/relationship"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/dominikbraun/graph/draw"
	yaml "gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ClusterController struct {
	config *Config
}

func NewClusterController(config *Config) *ClusterController {
	return &ClusterController{
		config: config,
	}
}

// GetCluster returns the unstructured Cluster object for a given cluster
func (c *ClusterController) GetCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	obj, _ := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	return obj, nil
}

// GetYAMLForCluster returns the yaml byte array for a given cluster
func (c *ClusterController) GetYAMLForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) ([]byte, error) {
	obj, err := c.GetCluster(ctx, client, name)
	if err != nil {
		return nil, err
	}
	objYAML, err := yaml.Marshal(obj.Object)
	if err != nil {
		return nil, err
	}
	return objYAML, nil
}

// GetYAMLForCluster returns the yaml byte array for a given cluster
func (c *ClusterController) GetNamespaceForCluster(ctx context.Context, client *multicluster.MultiClusterClient, cluster, namespace string) (*v1.Namespace, error) {
	namespaceObj, err := client.ClientSet.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return namespaceObj, nil
}

func (c *ClusterController) GetDetailsForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (*ClusterDetail, error) {
	serverVersion, _ := client.ClientSet.DiscoveryClient.ServerVersion()
	// Get the list of nodes
	nodes, err := client.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var memoryCapacity, cpuCapacity, podCapacity int64
	for _, node := range nodes.Items {
		memoryCapacity += node.Status.Capacity.Memory().Value()
		cpuCapacity += cpuCapacity + node.Status.Capacity.Cpu().Value()
		podCapacity += podCapacity + node.Status.Capacity.Pods().Value()
	}
	return &ClusterDetail{
		NodeCount:      len(nodes.Items),
		ServerVersion:  serverVersion.String(),
		MemoryCapacity: memoryCapacity,
		CPUCapacity:    cpuCapacity,
		PodsCapacity:   podCapacity,
	}, nil
}

// GetTopologyForCluster returns a map that describes topology for a given cluster
func (c *ClusterController) GetTopologyForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	// Build relationship graph based on GVK
	graph, rg, _ := relationship.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Count resources in all namespaces
	log.Info("Retrieving topology for cluster", "clusterName", name)
	rg, err := rg.CountRelationshipGraph(ctx, client.DynamicClient, client.ClientSet.DiscoveryClient, "")
	if err != nil {
		return nil, err
	}

	// Draw graph
	// TODO: This is drawn on the server side, not needed eventually
	file, _ := os.Create("./relationship.gv")
	_ = draw.DOT(graph, file)

	return c.ConvertGraphToMap(rg), nil
}

// GetTopologyForClusterNamespace returns a map that describes topology for a given namespace in a given cluster
func (c *ClusterController) GetTopologyForClusterNamespace(ctx context.Context, client *multicluster.MultiClusterClient, cluster, namespace string) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	// Build relationship graph based on GVK
	graph, rg, _ := relationship.BuildRelationshipGraph(ctx, client.DynamicClient)
	// Only count resources that belong to a specific namespace
	log.Info("Retrieving topology", "namespace", namespace, "cluster", cluster)
	rg, err := rg.CountRelationshipGraph(ctx, client.DynamicClient, client.ClientSet.DiscoveryClient, namespace)
	if err != nil {
		return nil, err
	}

	// Draw graph
	// TODO: This is drawn on the server side, not needed eventually
	file, _ := os.Create("./relationship.gv")
	_ = draw.DOT(graph, file)

	return c.ConvertGraphToMap(rg), nil
}

func (c *ClusterController) ConvertGraphToMap(rg *relationship.RelationshipGraph) map[string]ClusterTopology {
	m := make(map[string]ClusterTopology)
	for _, rgn := range rg.RelationshipNodes {
		rgnMap := rgn.ConvertToMap()
		m[rgn.GetHash()] = ClusterTopology{
			GroupVersionKind: rgn.GetHash(),
			Count:            rgn.ResourceCount,
			Relationship:     rgnMap,
		}
	}
	return m
}
