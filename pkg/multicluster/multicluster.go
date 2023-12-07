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

package multicluster

import (
	"context"
	"net"
	"net/url"

	clusterv1beta1 "github.com/KusionStack/karbour/pkg/apis/cluster/v1beta1"
	"github.com/KusionStack/karbour/pkg/scheme"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

// BuildDynamicClient returns a dynamic client based on the cluster name in the request
func BuildDynamicClient(ctx context.Context, hubConfig *restclient.Config, name string) (*dynamic.DynamicClient, error) {
	// Create the dynamic client using loopback hubConfig for Karbour apiserver
	hubClient, err := dynamic.NewForConfig(hubConfig)
	if err != nil {
		return nil, err
	}
	// If name is empty, return the dynamic client for the hub cluster
	if name == "" {
		return hubClient, nil
	}
	// otherwise, return the dynamic client for the spoke cluster
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	spokeUnstructured, err := hubClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	spokeObj, err := unstructuredToRuntimeObject(spokeUnstructured)
	if err != nil {
		return nil, err
	}
	spokeConfig, err := NewConfigFromCluster(spokeObj.(*clusterv1beta1.Cluster))
	if err != nil {
		return nil, err
	}
	spokeClient, err := dynamic.NewForConfig(spokeConfig)
	if err != nil {
		return nil, err
	}
	return spokeClient, nil
}

// BuildDiscoveryClient returns a discovery client based on the cluster name in the request
func BuildDiscoveryClient(ctx context.Context, hubConfig *restclient.Config, name string) (*discovery.DiscoveryClient, error) {
	// Create the dynamic client using loopback hubConfig for Karbour apiserver
	hubClient, err := dynamic.NewForConfig(hubConfig)
	if err != nil {
		return nil, err
	}
	clusterGVR := clusterv1beta1.SchemeGroupVersion.WithResource("clusters")
	spokeUnstructured, err := hubClient.Resource(clusterGVR).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	spokeObj, err := unstructuredToRuntimeObject(spokeUnstructured)
	if err != nil {
		return nil, err
	}
	spokeConfig, err := NewConfigFromCluster(spokeObj.(*clusterv1beta1.Cluster))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	// Create the discovery client
	clientset := kubernetes.NewForConfigOrDie(spokeConfig)
	discoveryClient := clientset.DiscoveryClient
	return discoveryClient, nil
}

// unstructuredToRuntimeObject converts an unstructured.Unstructured into a typed runtime.Object
func unstructuredToRuntimeObject(u *unstructured.Unstructured) (runtime.Object, error) {
	obj, err := scheme.Scheme.New(u.GroupVersionKind())
	if err != nil {
		return nil, err
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// NewConfigFromCluster takes in a v1beta1.Cluster object and return the corresponding restclient.Config object for client-go
func NewConfigFromCluster(c *clusterv1beta1.Cluster) (*restclient.Config, error) {
	cfg := &restclient.Config{}
	cfg.Host = c.Spec.Access.Endpoint
	cfg.CAData = c.Spec.Access.CABundle
	if c.Spec.Access.Insecure != nil && *c.Spec.Access.Insecure {
		cfg.TLSClientConfig = restclient.TLSClientConfig{Insecure: true}
	}
	switch c.Spec.Access.Credential.Type {
	case clusterv1beta1.CredentialTypeServiceAccountToken:
		cfg.BearerToken = c.Spec.Access.Credential.ServiceAccountToken
	case clusterv1beta1.CredentialTypeX509Certificate:
		cfg.CertData = c.Spec.Access.Credential.X509.Certificate
		cfg.KeyData = c.Spec.Access.Credential.X509.PrivateKey
	}
	u, err := url.Parse(c.Spec.Access.Endpoint)
	if err != nil {
		return nil, err
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}
	cfg.ServerName = host // apiserver may listen on SNI cert
	return cfg, nil
}
