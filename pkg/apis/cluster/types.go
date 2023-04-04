/*
Copyright The Karbour Authors.

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
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is an extension type to access a cluster
type Cluster struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ClusterSpec
	Status ClusterStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList is a list of Cluster objects.
type ClusterList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Cluster
}

type ClusterSpec struct {
	Provider  string
	Access    ClusterAccess
	Finalized *bool
}

type ClusterStatus struct {
	Healthy bool
}

type ClusterAccess struct {
	Endpoint   string
	CABundle   []byte
	Insecure   *bool
	Credential *ClusterAccessCredential
}

type ClusterAccessCredential struct {
	Type                CredentialType
	ServiceAccountToken string
	X509                *X509
}

type X509 struct {
	Certificate []byte
	PrivateKey  []byte
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterProxyOptions struct {
	metav1.TypeMeta

	// Path is the target api path of the proxy request.
	// e.g. "/healthz", "/api/v1"
	Path string
}
