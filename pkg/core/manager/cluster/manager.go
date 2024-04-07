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
	"time"

	"github.com/KusionStack/karbour/pkg/infra/multicluster"
	clusterv1beta1 "github.com/KusionStack/karbour/pkg/kubernetes/apis/cluster/v1beta1"
	"github.com/KusionStack/karbour/pkg/util/clusterinstall"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	yaml "gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	k8syaml "sigs.k8s.io/yaml"
)

type ClusterManager struct{}

// NewClusterManager returns a new ClusterManager object
func NewClusterManager() *ClusterManager {
	return &ClusterManager{}
}

// GetCluster returns the unstructured Cluster object for a given cluster
func (c *ClusterManager) GetCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	obj, err := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return SanitizeUnstructuredCluster(ctx, obj)
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
	clusterObj, err := clusterinstall.ConvertKubeconfigToCluster(name, displayName, description, restConfig)
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

// UpdateMetadata updates cluster by name with a full payload
func (c *ClusterManager) UpdateMetadata(ctx context.Context, client *multicluster.MultiClusterClient, name, displayName, description string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	// Make sure the cluster exists first
	currentObj, err := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if displayName == "" {
		displayName = name
	}

	// Update the display name and description. Updating name and kubeconfig is not allowed in this method.
	currentObj.Object["spec"].(map[string]interface{})["displayName"] = displayName
	currentObj.Object["spec"].(map[string]interface{})["description"] = description
	return client.DynamicClient.Resource(clusterGVR).Update(ctx, currentObj, metav1.UpdateOptions{})
}

// UpdateCredential updates cluster credential by name and a new kubeconfig
func (c *ClusterManager) UpdateCredential(ctx context.Context, client *multicluster.MultiClusterClient, name, kubeconfig string) (*unstructured.Unstructured, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	// Make sure the cluster exists first
	currentObj, err := client.DynamicClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	displayName := currentObj.Object["spec"].(map[string]interface{})["displayName"].(string)
	description := currentObj.Object["spec"].(map[string]interface{})["description"].(string)

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

// ListCluster returns the full list of clusters in a specific order
func (c *ClusterManager) ListCluster(ctx context.Context, client *multicluster.MultiClusterClient, orderBy SortCriteria, descending bool) (*unstructured.UnstructuredList, error) {
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	clusterList, err := client.DynamicClient.Resource(clusterGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	sanitizedClusterList := &unstructured.UnstructuredList{}
	for _, cluster := range clusterList.Items {
		sanitized, _ := SanitizeUnstructuredCluster(ctx, &cluster)
		sanitizedClusterList.Items = append(sanitizedClusterList.Items, *sanitized)
	}
	return SortUnstructuredList(sanitizedClusterList, orderBy, descending)
}

// CountCluster returns the summary of clusters by their health status
func (c *ClusterManager) CountCluster(ctx context.Context, client *multicluster.MultiClusterClient, config *rest.Config) (*ClusterSummary, error) {
	clusterSummary := &ClusterSummary{
		HealthyClusters:   make([]string, 0),
		UnhealthyClusters: make([]string, 0),
	}
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	clusterList, err := client.DynamicClient.Resource(clusterGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	clusterSummary.TotalCount = len(clusterList.Items)
	for _, cluster := range clusterList.Items {
		spokeClient, err := multicluster.BuildMultiClusterClient(ctx, config, cluster.GetName())
		if err != nil {
			return nil, err
		}
		// Check the health of the API server by querying the /healthz endpoint
		var statusCode int
		if err = spokeClient.ClientSet.RESTClient().
			Get().
			AbsPath("/healthz").
			Timeout(1 * time.Second).
			Do(ctx).
			StatusCode(&statusCode).
			Error(); err != nil || statusCode != 200 {
			clusterSummary.UnhealthyClusters = append(clusterSummary.UnhealthyClusters, cluster.GetName())
			clusterSummary.UnhealthyCount++
		} else {
			clusterSummary.HealthyClusters = append(clusterSummary.HealthyClusters, cluster.GetName())
			clusterSummary.HealthyCount++
		}
	}
	return clusterSummary, nil
}

// GetYAMLForCluster returns the yaml byte array for a given cluster
func (c *ClusterManager) GetYAMLForCluster(ctx context.Context, client *multicluster.MultiClusterClient, name string) ([]byte, error) {
	obj, err := c.GetCluster(ctx, client, name)
	if err != nil {
		return nil, err
	}
	return k8syaml.Marshal(obj)
}

// GetNamespaceForCluster returns the yaml byte array for a given cluster
func (c *ClusterManager) GetNamespace(ctx context.Context, client *multicluster.MultiClusterClient, namespace string) (*unstructured.Unstructured, error) {
	// Typed clientset clears the TypeMeta (APIVersion and Kind) when decoding: https://github.com/kubernetes/kubernetes/issues/80609
	// e.g return client.ClientSet.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{}) returns an empty TypeMeta
	nsGVR := v1.SchemeGroupVersion.WithResource("namespaces")
	return client.DynamicClient.Resource(nsGVR).Get(ctx, namespace, metav1.GetOptions{})
}

// GetNamespaceForCluster returns the yaml byte array for a given cluster
func (c *ClusterManager) GetNamespaceYAML(ctx context.Context, client *multicluster.MultiClusterClient, namespace string) ([]byte, error) {
	obj, err := c.GetNamespace(ctx, client, namespace)
	if err != nil {
		return nil, err
	}
	return k8syaml.Marshal(obj)
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
