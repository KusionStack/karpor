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
	"encoding/base64"
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
	"k8s.io/client-go/kubernetes"
)

type ClusterManager struct {
	config *Config
}

func NewClusterManager(config *Config) *ClusterManager {
	return &ClusterManager{
		config: config,
	}
}

// GetCluster returns the unstructured Cluster object for a given cluster
func (c *ClusterManager) GetCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	obj, _ := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	return obj, nil
}

// GetYAMLForCluster returns the yaml byte array for a given cluster
func (c *ClusterManager) GetYAMLForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) ([]byte, error) {
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
func (c *ClusterManager) GetNamespaceForCluster(ctx context.Context, client *multicluster.MultiClusterClient, cluster, namespace string) (*v1.Namespace, error) {
	namespaceObj, err := client.ClientSet.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return namespaceObj, nil
}

func (c *ClusterManager) GetDetailsForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (*ClusterDetail, error) {
	serverVersion, _ := client.ClientSet.DiscoveryClient.ServerVersion()
	// Get the list of nodes
	nodes, err := client.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
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

// GetTopologyForCluster returns a map that describes topology for a given cluster
func (c *ClusterManager) GetTopologyForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (map[string]ClusterTopology, error) {
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
func (c *ClusterManager) GetTopologyForClusterNamespace(ctx context.Context, client *multicluster.MultiClusterClient, cluster, namespace string) (map[string]ClusterTopology, error) {
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

func (c *ClusterManager) ConvertGraphToMap(rg *relationship.RelationshipGraph) map[string]ClusterTopology {
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

// SanitizeKubeConfigWithYAML takes a plain KubeConfig YAML string and returns
// a sanitized version with sensitive information masked.
func (c *ClusterManager) SanitizeKubeConfigWithYAML(ctx context.Context, plain string) (sanitize string, err error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)

	// Inform that the unmarshaling process has started.
	log.Info("Unmarshal the yaml file into the KubeConfig struct in SanitizeKubeConfigWithYAML")

	// Prepare KubeConfig structure to hold unmarshaled data.
	var config KubeConfig

	// Convert YAML to KubeConfig struct.
	if err = yaml.Unmarshal([]byte(plain), &config); err != nil {
		return "", err
	}

	// Perform sanitization of the KubeConfig data.
	c.SanitizeKubeConfigFor(&config)

	// Convert the sanitized KubeConfig back to YAML.
	result, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// SanitizeKubeConfigFor masks sensitive information within a KubeConfig object,
// such as user credentials and certificate data.
func (c *ClusterManager) SanitizeKubeConfigFor(config *KubeConfig) {
	// Iterate over each user and sanitize sensitive fields.
	for i := range config.Users {
		user := &config.Users[i].User
		if user.ClientCertificateData != "" {
			user.ClientCertificateData = maskContent(user.ClientCertificateData)
		}
		if user.ClientKeyData != "" {
			user.ClientKeyData = maskContent(user.ClientKeyData)
		}
		if user.Token != "" {
			user.Token = maskContent(user.Token)
		}
		if user.Username != "" {
			user.Username = maskContent(user.Username)
		}
		if user.Password != "" {
			user.Password = maskContent(user.Password)
		}
	}

	// Iterate over each cluster and sanitize certificate authority data.
	for i := range config.Clusters {
		cluster := &config.Clusters[i].Cluster
		if cluster.CertificateAuthorityData != "" {
			cluster.CertificateAuthorityData = maskContent(cluster.CertificateAuthorityData)
		}
	}
}

// ValidateKubeConfigWithYAML unmarshals YAML content into KubeConfig and validates it.
func (c *ClusterManager) ValidateKubeConfigWithYAML(ctx context.Context, plain string) (string, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)

	// Inform that the unmarshaling process has started.
	log.Info("Unmarshaling the YAML file into the KubeConfig struct in ValidateKubeConfigWithYAML")

	// Prepare KubeConfig structure to hold unmarshaled data.
	var config KubeConfig

	// Convert YAML to KubeConfig struct.
	if err := yaml.Unmarshal([]byte(plain), &config); err != nil {
		return "", err
	}

	// Validate the KubeConfig.
	return c.ValidateKubeConfigFor(ctx, &config)
}

// ValidateKubeConfigFor validates the provided KubeConfig.
func (c *ClusterManager) ValidateKubeConfigFor(ctx context.Context, config *KubeConfig) (string, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)

	// Validate if KubeConfig API version and kind are empty.
	if config.APIVersion == "" {
		return "", ErrMissingAPIVersion
	}
	if config.Kind == "" {
		return "", ErrMissingKind
	}

	// Check for at least one cluster and user defined.
	if len(config.Clusters) == 0 {
		return "", ErrMissingClusterEntry
	}
	if len(config.Users) == 0 {
		return "", ErrMissingUserEntry
	}

	// Validate cluster information including server address and certificate.
	for _, clusterEntry := range config.Clusters {
		if clusterEntry.Name == "" {
			return "", ErrMissingClusterName
		}
		cluster := clusterEntry.Cluster
		if cluster.Server == "" {
			return "", ErrMissingClusterServer
		}
		if cluster.CertificateAuthorityData == "" {
			return "", ErrMissingCertificateAuthority
		}
		if _, err := base64.StdEncoding.DecodeString(cluster.CertificateAuthorityData); err != nil {
			return "", ErrInvalidCertificateAuthority
		}

		// Check cluster server address connectivity
		if err := checkEndpointConnectivity(cluster.Server); err != nil {
			return "", ErrClusterServerConnectivity
		}
	}

	// Use the provided KubeConfig to build the clientConfig.
	clientConfig, err := buildClientConfigFromKubeConfig(config)
	if err != nil {
		return "", ErrBuildClientConfig
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return "", ErrCreateClientSet
	}

	// Try fetching server version information to test connectivity.
	if info, err := clientset.Discovery().ServerVersion(); err != nil {
		return "", ErrGetServerVersion
	} else {
		log.Info("KubeConfig is valid and the cluster is reachable.", "serverVersion", info.String())
		return info.String(), nil
	}
}
