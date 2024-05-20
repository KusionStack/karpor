/*
Copyright 2020 The Kubernetes Authors.
Copyright 2023 The KusionStack Authors.

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

package generator

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/keyutil"
)

const kubeConfigTemplate = `apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://karpor-server.karpor.svc:7443
  name: karpor
contexts:
- context:
    cluster: karpor
    user: karpor
  name: default
current-context: default
kind: Config
users:
- name: karpor
  user:
    client-certificate-data: %s
    client-key-data: %s`

type Generator struct {
	clientSet      *kubernetes.Clientset
	namespace      string
	certName       string
	kubeConfigName string
}

func NewGenerator(cfg *rest.Config, namespace string, certName string, kubeConfigName string) (*Generator, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cfg is required buit was nil")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required but was empty")
	}
	if certName == "" {
		return nil, fmt.Errorf("certName is required but was empty")
	}
	if kubeConfigName == "" {
		return nil, fmt.Errorf("kubeConfigName is required but was empty")
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		certName:       certName,
		kubeConfigName: kubeConfigName,
		clientSet:      clientSet,
		namespace:      namespace,
	}, nil
}

// Generate is a function that orchestrates the creation and application of certificates and kubeconfig necessary for a karpor sever.
func (g *Generator) Generate(ctx context.Context) error {
	caCert, caKey, kubeConfig, err := generateConfig()
	if err != nil {
		return err
	}

	err = g.applyConfig(ctx, caCert, caKey, kubeConfig)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) applyConfig(ctx context.Context, caCert *x509.Certificate, caKey crypto.Signer, kubeConfig string) error {
	err := g.applyCertToSecret(ctx, caCert, caKey)
	if err != nil {
		return err
	}
	err = g.applyKubeConfigToConfigMap(ctx, kubeConfig)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) applyCertToSecret(ctx context.Context, caCert *x509.Certificate, caKey crypto.Signer) error {
	caCertData := EncodeCertPEM(caCert)
	caKeyData, err := keyutil.MarshalPrivateKeyToPEM(caKey)
	if err != nil {
		return fmt.Errorf("unable to marshal private key to PEM %s", err)
	}

	secret := applycorev1.Secret(g.certName, g.namespace)
	secret.StringData = map[string]string{
		"ca.crt": string(caCertData),
		"ca.key": string(caKeyData),
	}
	_, err = g.clientSet.CoreV1().Secrets(g.namespace).Apply(ctx, secret, metav1.ApplyOptions{FieldManager: "cert-generator", Force: true})
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) applyKubeConfigToConfigMap(ctx context.Context, kubeConfig string) error {
	cm := applycorev1.ConfigMap(g.kubeConfigName, g.namespace)
	cm.Data = map[string]string{
		"config": kubeConfig,
	}
	_, err := g.clientSet.CoreV1().ConfigMaps(g.namespace).Apply(ctx, cm, metav1.ApplyOptions{FieldManager: "cert-generator", Force: true})
	if err != nil {
		return err
	}
	return nil
}

func generateConfig() (*x509.Certificate, crypto.Signer, string, error) {
	caCert, caKey, err := generateCA()
	if err != nil {
		return nil, nil, "", err
	}
	cert, key, err := generateCert(caCert, caKey)
	if err != nil {
		return nil, nil, "", err
	}
	kubeConfig, err := generateAdminKubeconfig(cert, key)
	if err != nil {
		return nil, nil, "", err
	}
	return caCert, caKey, kubeConfig, nil
}

func generateCA() (*x509.Certificate, crypto.Signer, error) {
	caConfig := Config{
		CommonName:   "kubernetes",
		Organization: nil,
		Year:         100,
	}
	caCert, caKey, err := NewCaCertAndKey(caConfig)
	if err != nil {
		return nil, nil, err
	}
	return caCert, caKey, nil
}

func generateCert(caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, crypto.Signer, error) {
	certConfig := Config{
		CAName:       "kubernetes",
		CommonName:   "kubernetes-admin",
		Organization: []string{"system:masters"},
		Year:         100,
		AltNames:     AltNames{},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	cert, key, err := NewCaCertAndKeyFromRoot(certConfig, caCert, caKey)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

func generateAdminKubeconfig(cert *x509.Certificate, key crypto.Signer) (string, error) {
	certData := EncodeCertPEM(cert)
	keyData, err := keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		return "", fmt.Errorf("unable to marshal private key to PEM %s", err)
	}
	return fmt.Sprintf(kubeConfigTemplate, base64.StdEncoding.EncodeToString(certData), base64.StdEncoding.EncodeToString(keyData)), nil
}
