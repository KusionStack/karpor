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
	"fmt"
	"os"

	clusterv1beta1 "github.com/KusionStack/karbour/pkg/apis/cluster/v1beta1"
	"github.com/KusionStack/karbour/pkg/clusterinstall"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/relationship"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/dominikbraun/graph/draw"
	yaml "gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterManager struct {
	config *Config
}

// NewClusterManager returns a new ClusterManager object
func NewClusterManager(config *Config) *ClusterManager {
	return &ClusterManager{
		config: config,
	}
}

// GetCluster returns the unstructured Cluster object for a given cluster
func (c *ClusterManager) GetCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	obj, err := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	return c.SanitizeUnstructuredCluster(ctx, obj)
}

// CreateCluster creates a new Cluster resource in the hub cluster and returns the created unstructured Cluster object
func (c *ClusterManager) CreateCluster(ctx context.Context, client *multicluster.MultiClusterClient, name, displayName, description, kubeconfig string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	// Make sure the cluster does not exist first
	currentObj, err := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if currentObj != nil {
		return nil, fmt.Errorf("cluster %s already exists. Try updating it instead", name)
	}

	// Create rest.Config from the incoming KubeConfig
	restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	// Convert the rest.Config to Cluster object and create it using dynamic client
	clusterObj, err := clusterinstall.ConvertKubeconfigToCluster(name, description, displayName, restConfig)
	if err != nil {
		return nil, err
	}
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(clusterObj)
	if err != nil {
		return nil, err
	}
	unstructuredCluster := &unstructured.Unstructured{Object: unstructuredMap}
	return client.DynamicClient.Resource(clusterGVR).Create(ctx, unstructuredCluster, metav1.CreateOptions{})
}

// UpdateCluster updates cluster by name with a full payload
func (c *ClusterManager) UpdateMetadata(ctx context.Context, client *multicluster.MultiClusterClient, name, displayName, description string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	// Make sure the cluster exists first
	currentObj, err := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if currentObj == nil {
		return nil, fmt.Errorf("cluster %s not found. Try creating it instead", name)
	}

	// Update the display name and description. Updating name and kubeconfig is not allowed here.
	currentObj.Object["spec"].(map[string]interface{})["displayName"] = displayName
	currentObj.Object["spec"].(map[string]interface{})["description"] = description
	return client.DynamicClient.Resource(clusterGVR).Update(ctx, currentObj, metav1.UpdateOptions{})
}

// UpdateCredential updates cluster credential by name and a new kubeconfig
func (c *ClusterManager) UpdateCredential(ctx context.Context, client *multicluster.MultiClusterClient, name, displayName, description, kubeconfig string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	// Make sure the cluster exists first
	currentObj, err := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if currentObj == nil {
		return nil, fmt.Errorf("cluster %s not found. Try creating it instead", name)
	}

	// Create new restConfig from updated kubeconfig
	restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	// Convert the rest.Config to Cluster object and update it using dynamic client
	clusterObj, err := clusterinstall.ConvertKubeconfigToCluster(name, displayName, description, restConfig)
	if err != nil {
		return nil, err
	}
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(clusterObj)
	if err != nil {
		return nil, err
	}
	unstructuredMap["metadata"].(map[string]interface{})["resourceVersion"] = currentObj.Object["metadata"].(map[string]interface{})["resourceVersion"]
	unstructuredCluster := &unstructured.Unstructured{Object: unstructuredMap}
	return client.DynamicClient.Resource(clusterGVR).Update(ctx, unstructuredCluster, metav1.UpdateOptions{})
}

// DeleteCluster deletes the cluster by name
func (c *ClusterManager) DeleteCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) error {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	return client.DynamicClient.Resource(clusterGVR).Delete(ctx, name, metav1.DeleteOptions{})
}

// DeleteCluster deletes the cluster by name
func (c *ClusterManager) ListCluster(ctx context.Context, client *multicluster.MultiClusterClient) (*unstructured.UnstructuredList, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	clusterList, err := client.DynamicClient.Resource(clusterGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	sanitizedClusterList := &unstructured.UnstructuredList{}
	for _, cluster := range clusterList.Items {
		sanitized, _ := c.SanitizeUnstructuredCluster(ctx, &cluster)
		sanitizedClusterList.Items = append(sanitizedClusterList.Items, *sanitized)
	}
	return sanitizedClusterList, nil
}

// GetYAMLForCluster returns the yaml byte array for a given cluster
func (c *ClusterManager) GetYAMLForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) ([]byte, error) {
	obj, err := c.GetCluster(ctx, client, name)
	if err != nil {
		return nil, err
	}
	sanitized, _ := c.SanitizeUnstructuredCluster(ctx, obj)
	return yaml.Marshal(sanitized)
}

// GetYAMLForCluster returns the yaml byte array for a given cluster
func (c *ClusterManager) GetNamespaceForCluster(ctx context.Context, client *multicluster.MultiClusterClient, cluster, namespace string) (*v1.Namespace, error) {
	return client.ClientSet.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
}

// GetDetailsForCluster returns ClusterDetail object for a given cluster
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

// ConvertGraphToMap returns a map[string]ClusterTopology for a given relationship.RelationshipGraph
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

// SanitizeUnstructuredCluster masks sensitive information
// within a Unstructured cluster object, such as user
// credentials and certificate data.
func (c *ClusterManager) SanitizeUnstructuredCluster(ctx context.Context, cluster *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	log := ctxutil.GetLogger(ctx)

	// Inform that the unmarshaling process has started.
	log.Info("Sanitizing unstructured cluster...")
	sanitized := cluster
	if token, ok := sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["serviceAccountToken"]; ok {
		sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["serviceAccountToken"] = maskContent(token.(string))
	}
	if x509, ok := sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["x509"]; ok {
		sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["x509"].(map[string]interface{})["certificate"] = maskContent(x509.(map[string]interface{})["certificate"].(string))
		sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["x509"].(map[string]interface{})["privateKey"] = maskContent(x509.(map[string]interface{})["privateKey"].(string))
	}
	if caBundle, ok := sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["caBundle"]; ok {
		sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["caBundle"] = maskContent(caBundle.(string))
	}
	if _, ok := sanitized.Object["metadata"].(map[string]interface{})["annotations"]; ok {
		sanitized.Object["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["kubectl.kubernetes.io/last-applied-configuration"] = "[redacted]"
	}
	return sanitized, nil
}
