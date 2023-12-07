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
	"github.com/KusionStack/karbour/pkg/relationship"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/dominikbraun/graph/draw"
	yaml "gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

type ClusterController struct {
	config *Config
}

func NewClusterController(config *Config) *ClusterController {
	return &ClusterController{
		config: config,
	}
}

func (c *ClusterController) GetCluster(ctx context.Context, hubClient *dynamic.DynamicClient, name string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	obj, _ := hubClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	return obj, nil
}

// GetYAMLForCluster returns the yaml byte array for a given cluster
func (c *ClusterController) GetYAMLForCluster(ctx context.Context, hubClient *dynamic.DynamicClient, name string) ([]byte, error) {
	obj, err := c.GetCluster(ctx, hubClient, name)
	if err != nil {
		return nil, err
	}
	objYAML, err := yaml.Marshal(obj.Object)
	if err != nil {
		return nil, err
	}
	return objYAML, nil
}

// GetTopologyForCluster returns a map that describes topology for a given cluster
func (c *ClusterController) GetTopologyForCluster(ctx context.Context, spokeClient *dynamic.DynamicClient, discoveryClient *discovery.DiscoveryClient, name string) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	// Build relationship graph based on GVK
	graph, rg, _ := relationship.BuildRelationshipGraph(ctx, spokeClient)
	// Count resources in all namespaces
	log.Info("Retrieving topology for cluster", "clusterName", name)
	rg, err := rg.CountRelationshipGraph(ctx, spokeClient, discoveryClient, "")
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
func (c *ClusterController) GetTopologyForClusterNamespace(ctx context.Context, spokeClient *dynamic.DynamicClient, discoveryClient *discovery.DiscoveryClient, cluster, namespace string) (map[string]ClusterTopology, error) {
	log := ctxutil.GetLogger(ctx)

	// Build relationship graph based on GVK
	graph, rg, _ := relationship.BuildRelationshipGraph(ctx, spokeClient)
	// Only count resources that belong to a specific namespace
	log.Info("Retrieving topology", "namespace", namespace, "cluster", cluster)
	rg, err := rg.CountRelationshipGraph(ctx, spokeClient, discoveryClient, namespace)
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
