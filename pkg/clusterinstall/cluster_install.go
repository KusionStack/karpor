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

package clusterinstall

import (
	"context"
	"fmt"

	clusterv1beta1 "github.com/KusionStack/karbour/pkg/apis/cluster/v1beta1"
	"github.com/KusionStack/karbour/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ConvertKubeconfigToCluster(name string, cfg *rest.Config) (*clusterv1beta1.Cluster, error) {
	cluster := clusterv1beta1.Cluster{}
	cluster.Name = name
	access := clusterv1beta1.ClusterAccess{}
	if cfg.Insecure {
		access.CABundle = cfg.CAData
	} else {
		access.Insecure = &cfg.Insecure
	}
	credential := &clusterv1beta1.ClusterAccessCredential{}
	if cfg.KeyData != nil && cfg.CertData != nil {
		credential.Type = clusterv1beta1.CredentialTypeX509Certificate
		credential.X509 = &clusterv1beta1.X509{
			Certificate: cfg.CertData,
			PrivateKey:  cfg.KeyData,
		}
	} else if cfg.BearerToken != "" {
		credential.Type = clusterv1beta1.CredentialTypeServiceAccountToken
		credential.ServiceAccountToken = cfg.BearerToken
	} else {
		return nil, fmt.Errorf("failed to parse credential from kubeconfig")
	}
	access.Credential = credential
	cluster.Spec.Access = access
	return &cluster, nil
}

func ProbeWithHealthz(cfg *rest.Config) error {
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}
	var statusCode int
	if err := client.RESTClient().
		Get().
		AbsPath("/healthz").
		Do(context.TODO()).
		StatusCode(&statusCode).
		Error(); err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("status code is %d, not 200", statusCode)
	}
	return nil
}

func CountClusters(cfg *rest.Config) (int, error) {
	client, err := versioned.NewForConfig(cfg)
	if err != nil {
		return 0, err
	}
	clusters, err := client.ClusterV1beta1().Clusters().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(clusters.Items), nil
}
