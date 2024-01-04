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
	"fmt"

	clusterv1beta1 "github.com/KusionStack/karbour/pkg/kubernetes/apis/cluster/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func ConvertKubeconfigToCluster(name, displayName, description string, cfg *rest.Config) (*clusterv1beta1.Cluster, error) {
	cluster := clusterv1beta1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterv1beta1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
	}
	cluster.Name = name
	cluster.Spec.Description = description
	if displayName != "" {
		cluster.Spec.DisplayName = displayName
	} else {
		cluster.Spec.DisplayName = name
	}
	access := clusterv1beta1.ClusterAccess{}
	if !cfg.Insecure {
		access.CABundle = cfg.CAData
	} else {
		access.Insecure = &cfg.Insecure
	}
	access.Endpoint = cfg.Host
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
