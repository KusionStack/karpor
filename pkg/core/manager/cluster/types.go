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

import "errors"

var (
	ErrMissingAPIVersion           = errors.New("apiVersion is required")
	ErrMissingKind                 = errors.New("kind is required")
	ErrMissingClusterEntry         = errors.New("at least one cluster entry is required")
	ErrMissingUserEntry            = errors.New("at least one user entry is required")
	ErrMissingClusterName          = errors.New("cluster name is required")
	ErrMissingClusterServer        = errors.New("cluster server is required")
	ErrBothInsecureAndCertificateAuthority = errors.New("certificate-authority-data and insecure-skip-tls-verify couldn't both be set")
	ErrMissingCertificateAuthority = errors.New("certificate-authority-data is required")
	ErrInvalidCertificateAuthority = errors.New("certificate-authority-data is invalid")
	ErrClusterServerConnectivity   = errors.New("cannot connect to the cluster server")
	ErrBuildClientConfig           = errors.New(
		"failed to create client config from provided KubeConfig",
	)
	ErrCreateClientSet  = errors.New("failed to create clientset")
	ErrGetServerVersion = errors.New("failed to connect to the cluster")
)

type SortCriteria int

const (
	ByTimestamp SortCriteria = iota
	ByName
)

type ClusterSummary struct {
	TotalCount        int      `json:"totalCount"`
	HealthyCount      int      `json:"healthyCount"`
	HealthyClusters   []string `json:"healthyClusters"`
	UnhealthyCount    int      `json:"unhealthyCount"`
	UnhealthyClusters []string `json:"unhealthyClusters"`
}

// KubeConfig represents the structure of a kubeconfig file
//
//nolint:tagliatelle
type KubeConfig struct {
	APIVersion     string         `yaml:"apiVersion"`
	Kind           string         `yaml:"kind"`
	Clusters       []ClusterEntry `yaml:"clusters"`
	Contexts       []ContextEntry `yaml:"contexts"`
	CurrentContext string         `yaml:"current-context"`
	Users          []UserEntry    `yaml:"users"`
}

// ClusterEntry represents each cluster entry in kubeconfig
type ClusterEntry struct {
	Name    string  `yaml:"name"`
	Cluster Cluster `yaml:"cluster"`
}

// Cluster contains the cluster information
//
//nolint:tagliatelle
type Cluster struct {
	Insecure bool 					`yaml:"insecure-skip-tls-verify,omitempty"`
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
}

// ContextEntry represents each context entry in kubeconfig
type ContextEntry struct {
	Name    string  `yaml:"name"`
	Context Context `yaml:"context"`
}

// Context contains the context information
type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

// UserEntry represents each user entry in kubeconfig
type UserEntry struct {
	Name string `yaml:"name"`
	User User   `yaml:"user"`
}

// User contains the user information
//
//nolint:tagliatelle
type User struct {
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
	Token                 string `yaml:"token,omitempty"`
	Username              string `yaml:"username,omitempty"`
	Password              string `yaml:"password,omitempty"`
}
