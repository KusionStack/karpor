/*
Copyright The Alipay Authors.

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

package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type CredentialType string

const (
	CredentialTypeServiceAccountToken CredentialType = "ServiceAccountToken"
	CredentialTypeX509Certificate     CredentialType = "X509Certificate"
	CredentialTypeUnifiedIdentity     CredentialType = "UnifiedIdentity"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterExtension is an extension type to access a cluster
type ClusterExtension struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterExtensionSpec `json:"spec"`
	// +optional
	Status ClusterExtensionStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterExtensionList is a list of ClusterExtension objects.
type ClusterExtensionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ClusterExtension `json:"items"`
}

type ClusterExtensionSpec struct {
	Provider string        `json:"provider"`
	Access   ClusterAccess `json:"access"`
	// +optional
	Finalized *bool `json:"finalized,omitempty"`
}

type ClusterExtensionStatus struct {
	// +optional
	Healthy bool `json:"healthy,omitempty"`
}

type ClusterAccess struct {
	Endpoint string `json:"endpoint"`
	// +optional
	CABundle []byte `json:"caBundle,omitempty"`
	// +optional
	Insecure *bool `json:"insecure,omitempty"`
	// +optional
	Credential *ClusterAccessCredential `json:"credential,omitempty"`
}

type ClusterAccessCredential struct {
	Type CredentialType `json:"type"`
	// +optional
	ServiceAccountToken string `json:"serviceAccountToken,omitempty"`
	// +optional
	X509 *X509 `json:"x509,omitempty"`
}

type X509 struct {
	Certificate []byte `json:"certificate"`
	PrivateKey  []byte `json:"privateKey"`
}
