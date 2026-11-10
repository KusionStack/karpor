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

	"github.com/KusionStack/karpor/pkg/util/ctxutil"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdv1 "k8s.io/client-go/tools/clientcmd/api/v1"
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
		host = net.JoinHostPort(u.Hostname(), "443") // Kubernetes API default port is 443
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

// resolveCurrentContext returns the cluster and user entries referenced by the
// current context of the provided KubeConfig, in the same way as clientcmd
// resolves them when the cluster is created.
func resolveCurrentContext(config *KubeConfig) (*ClusterEntry, *UserEntry, error) {
	if config.CurrentContext == "" {
		return nil, nil, ErrMissingCurrentContext
	}

	var current *Context
	for i := range config.Contexts {
		if config.Contexts[i].Name == config.CurrentContext {
			current = &config.Contexts[i].Context
			break
		}
	}
	if current == nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrContextNotFound, config.CurrentContext)
	}

	var clusterEntry *ClusterEntry
	for i := range config.Clusters {
		if config.Clusters[i].Name == current.Cluster {
			clusterEntry = &config.Clusters[i]
			break
		}
	}
	if clusterEntry == nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrContextClusterNotFound, current.Cluster)
	}

	var userEntry *UserEntry
	for i := range config.Users {
		if config.Users[i].Name == current.User {
			userEntry = &config.Users[i]
			break
		}
	}
	if userEntry == nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrContextUserNotFound, current.User)
	}

	return clusterEntry, userEntry, nil
}

// buildClientConfigFromKubeConfig generates a clientConfig from the cluster and
// user of the current context in the provided KubeConfig.
func buildClientConfigFromKubeConfig(config *KubeConfig) (*rest.Config, error) {
	clusterEntry, userEntry, err := resolveCurrentContext(config)
	if err != nil {
		return nil, err
	}

	// Create an initial rest.Config object.
	clientConfig := &rest.Config{}
	// Set the API server and authentication details.
	cluster := clusterEntry.Cluster
	clientConfig.Host = cluster.Server
	clientConfig.Insecure = cluster.Insecure // allow insecure connection
	if len(cluster.CertificateAuthorityData) > 0 {
		if plain, err := base64.StdEncoding.DecodeString(cluster.CertificateAuthorityData); err != nil {
			return nil, errors.Wrapf(
				err,
				"invalid certificate-authority-data for cluster %s",
				clusterEntry.Name,
			)
		} else {
			clientConfig.CAData = plain
		}
	}

	user := userEntry.User
	clientConfig.Username = user.Username
	clientConfig.Password = user.Password
	clientConfig.BearerToken = user.Token
	if len(user.Token) > 0 || (len(user.Password) > 0 && len(user.Username) > 0) {
		// prefer token or username/password
		return clientConfig, nil
	}
	if user.Exec != nil {
		clientConfig.ExecProvider = buildExecProvider(user.Exec)
		return clientConfig, nil
	}
	if plain, err := base64.StdEncoding.DecodeString(user.ClientCertificateData); err != nil {
		return nil, fmt.Errorf(
			"invalid client-certificate-data for user %s: %v",
			userEntry.Name,
			err,
		)
	} else {
		clientConfig.CertData = plain
	}
	if plain, err := base64.StdEncoding.DecodeString(user.ClientKeyData); err != nil {
		return nil, fmt.Errorf(
			"invalid client-key-data for user %s: %v",
			userEntry.Name,
			err,
		)
	} else {
		clientConfig.KeyData = plain
	}

	return clientConfig, nil
}

// buildExecProvider converts the exec section of a KubeConfig user into the
// exec credential plugin configuration used by client-go.
func buildExecProvider(exec *ExecConfig) *clientcmdapi.ExecConfig {
	env := make([]clientcmdapi.ExecEnvVar, 0, len(exec.Env))
	for _, e := range exec.Env {
		env = append(env, clientcmdapi.ExecEnvVar{Name: e.Name, Value: e.Value})
	}

	// Default the interactive mode in the same way as clientcmd does when
	// loading a kubeconfig, so that the validated config matches the one used
	// after the cluster is created.
	defaulted := &clientcmdv1.ExecConfig{
		APIVersion:      exec.APIVersion,
		InteractiveMode: clientcmdv1.ExecInteractiveMode(exec.InteractiveMode),
	}
	clientcmdv1.SetDefaults_ExecConfig(defaulted)

	return &clientcmdapi.ExecConfig{
		Command:            exec.Command,
		Args:               exec.Args,
		Env:                env,
		APIVersion:         exec.APIVersion,
		InstallHint:        exec.InstallHint,
		ProvideClusterInfo: exec.ProvideClusterInfo,
		InteractiveMode:    clientcmdapi.ExecInteractiveMode(defaulted.InteractiveMode),
	}
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
	// Environment variables of the exec plugin usually carry credentials, such
	// as AWS_SECRET_ACCESS_KEY for AWS EKS.
	if env, ok, _ := unstructured.NestedFieldNoCopy(sanitized.Object, "spec", "access", "credential", "execConfig", "env"); ok {
		envVars, _ := env.([]interface{})
		for _, envVar := range envVars {
			if envVar, ok := envVar.(map[string]interface{}); ok {
				if value, ok := envVar["value"].(string); ok && value != "" {
					envVar["value"] = maskContent(value)
				}
			}
		}
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
