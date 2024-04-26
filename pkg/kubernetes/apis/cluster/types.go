/*
Copyright The Karpor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CredentialType string

const (
	CredentialTypeServiceAccountToken CredentialType = "ServiceAccountToken"
	CredentialTypeX509Certificate     CredentialType = "X509Certificate"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is an extension type to access a cluster
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"` //nolint:tagliatelle

	Spec   ClusterSpec   `json:"spec"`
	Status ClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList is a list of Cluster objects.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"` //nolint:tagliatelle

	Items []Cluster `json:"items"`
}

type ClusterSpec struct {
	Provider string        `json:"provider"`
	Access   ClusterAccess `json:"access"`
	// +optional
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Finalized   *bool  `json:"finalized,omitempty"`
}

type ClusterStatus struct {
	Healthy bool `json:"healthy,omitempty"`
}

type ClusterAccess struct {
	Endpoint string `json:"endpoint"`
	// +optional
	CABundle   []byte                   `json:"caBundle,omitempty"`
	Insecure   *bool                    `json:"insecure,omitempty"`
	Credential *ClusterAccessCredential `json:"credential,omitempty"`
}

type ClusterAccessCredential struct {
	Type CredentialType `json:"type"`
	// +optional
	ServiceAccountToken string `json:"serviceAccountToken,omitempty"`
	X509                *X509  `json:"x509,omitempty"`
}

type X509 struct {
	Certificate []byte `json:"certificate"`
	PrivateKey  []byte `json:"privateKey"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Path is the target api path of the proxy request.
	// e.g. "/healthz", "/api/v1"
	Path string `json:"path,omitempty"`
}
