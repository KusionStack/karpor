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
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

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
	timeout := time.Duration(5 * time.Second)
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
			return nil, fmt.Errorf("invalid certificate-authority-data for cluster %s: %v", config.Clusters[0].Name, err)
		} else {
			clientConfig.CAData = plain
		}
	}

	if len(config.Users) > 0 {
		user := config.Users[0].User
		clientConfig.Username = user.Username
		clientConfig.Password = user.Password
		if plain, err := base64.StdEncoding.DecodeString(user.ClientCertificateData); err != nil {
			return nil, fmt.Errorf("invalid client-certificate-data for user %s: %v", config.Users[0].Name, err)
		} else {
			clientConfig.CertData = plain
		}
		if plain, err := base64.StdEncoding.DecodeString(user.ClientKeyData); err != nil {
			return nil, fmt.Errorf("invalid client-key-data for user %s: %v", config.Users[0].Name, err)
		} else {
			clientConfig.KeyData = plain
		}
	}

	return clientConfig, nil
}
