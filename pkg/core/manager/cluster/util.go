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
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

// maskContent to apply MD5 hash and mask the content.
func maskContent(content string) string {
	// Apply MD5 hash
	hash := fmt.Sprintf("%x", md5.Sum([]byte(content)))

	// Calculate the range for masking
	maskLength := len(hash) * 3 / 4 // Three quarters of the hash length
	start := len(hash) / 8          // Start masking a quarter in
	end := start + maskLength       // End masking
	masked := hash[:start] + strings.Repeat("*", maskLength) + hash[end:]

	return masked
}

// checkEndpointConnectivity checks the network connectivity of the Kubernetes
// API endpoint.
func checkEndpointConnectivity(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	host := u.Host
	if u.Port() == "" {
		host = fmt.Sprintf("%s:443", u.Hostname()) // Kubernetes API default port is 443
	}

	// Set timeout duration
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

// buildClientConfigFromKubeConfig generates a clientConfig from the provided
// KubeConfig.
func buildClientConfigFromKubeConfig(config *KubeConfig) (*rest.Config, error) {
	// Create an initial rest.Config object.
	clientConfig := &rest.Config{}

	// Set the API server and authentication details.
	if len(config.Clusters) > 0 {
		cluster := config.Clusters[0].Cluster
		clientConfig.Host = cluster.Server
		if plain, err := base64.StdEncoding.DecodeString(cluster.CertificateAuthorityData); err != nil {
			return nil, errors.Wrapf(
				err,
				"invalid certificate-authority-data for cluster %s",
				config.Clusters[0].Name,
			)
		} else {
			clientConfig.CAData = plain
		}
	}

	if len(config.Users) > 0 {
		user := config.Users[0].User
		clientConfig.Username = user.Username
		clientConfig.Password = user.Password
		if plain, err := base64.StdEncoding.DecodeString(user.ClientCertificateData); err != nil {
			return nil, fmt.Errorf(
				"invalid client-certificate-data for user %s: %v",
				config.Users[0].Name,
				err,
			)
		} else {
			clientConfig.CertData = plain
		}
		if plain, err := base64.StdEncoding.DecodeString(user.ClientKeyData); err != nil {
			return nil, fmt.Errorf(
				"invalid client-key-data for user %s: %v",
				config.Users[0].Name,
				err,
			)
		} else {
			clientConfig.KeyData = plain
		}
	}

	return clientConfig, nil
}

// SanitizeUnstructuredCluster masks sensitive information within a Unstructured
// cluster object, such as user credentials and certificate data.
func SanitizeUnstructuredCluster(
	ctx context.Context,
	cluster *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
	log := ctxutil.GetLogger(ctx)

	// Inform that the unmarshaling process has started.
	log.Info("Sanitizing unstructured cluster...")
	sanitized := cluster
	if token, ok := sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["serviceAccountToken"]; ok {
		sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["serviceAccountToken"] = maskContent(
			token.(string),
		)
	}
	if x509, ok := sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["x509"]; ok &&
		x509 != nil {
		sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["x509"].(map[string]interface{})["certificate"] = []byte(
			maskContent(x509.(map[string]interface{})["certificate"].(string)),
		)
		sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["credential"].(map[string]interface{})["x509"].(map[string]interface{})["privateKey"] = []byte(
			maskContent(x509.(map[string]interface{})["privateKey"].(string)),
		)
	}
	if caBundle, ok := sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["caBundle"]; ok {
		sanitized.Object["spec"].(map[string]interface{})["access"].(map[string]interface{})["caBundle"] = []byte(
			maskContent(caBundle.(string)),
		)
	}
	if _, ok := sanitized.Object["metadata"].(map[string]interface{})["annotations"]; ok {
		sanitized.Object["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["kubectl.kubernetes.io/last-applied-configuration"] = "[redacted]"
	}
	return sanitized, nil
}

// SortUnstructuredList returns a sorted unstructured.UnstructuredList based on criteria
func SortUnstructuredList(
	unList *unstructured.UnstructuredList,
	criteria SortCriteria,
	descending bool,
) (*unstructured.UnstructuredList, error) {
	switch criteria {
	case ByTimestamp:
		sort.Slice(unList.Items, func(i, j int) bool {
			if descending {
				return unList.Items[i].GetCreationTimestamp().
					UTC().
					After(unList.Items[j].GetCreationTimestamp().UTC())
			} else {
				return unList.Items[i].GetCreationTimestamp().UTC().Before(unList.Items[j].GetCreationTimestamp().UTC())
			}
		})
	case ByName:
		sort.Slice(unList.Items, func(i, j int) bool {
			if descending {
				return unList.Items[i].GetName() > unList.Items[j].GetName()
			} else {
				return unList.Items[i].GetName() < unList.Items[j].GetName()
			}
		})
	}
	return unList, nil
}
