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

package cluster

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
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec ClusterExtensionSpec
	// +optional
	Status ClusterExtensionStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterExtensionList is a list of ClusterExtension objects.
type ClusterExtensionList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []ClusterExtension
}

type ClusterExtensionSpec struct {
	Provider string
	Access   ClusterAccess
	// +optional
	Finalized *bool
}

type ClusterExtensionStatus struct {
	// +optional
	Healthy bool
}

type ClusterAccess struct {
	Endpoint string
	// +optional
	CABundle []byte
	// +optional
	Insecure *bool
	// +optional
	Credential *ClusterAccessCredential
}

type ClusterAccessCredential struct {
	Type CredentialType
	// +optional
	ServiceAccountToken string
	// +optional
	X509 *X509
}

type X509 struct {
	Certificate []byte
	PrivateKey  []byte
}
