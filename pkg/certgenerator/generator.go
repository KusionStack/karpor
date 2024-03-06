package cert

import (
	"context"
	"crypto"
	"crypto/x509"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/keyutil"

	"sigs.k8s.io/yaml"
)

// Generate is a function that orchestrates the creation and application of certificates and kubeconfig necessary for a Kubernetes cluster.
func Generate() error {
	caCert, caKey, err := GenerateCA()
	if err != nil {
		return err
	}
	cert, key, err := GenerateCert(caCert, caKey)
	if err != nil {
		return err
	}
	kubeConfig, err := GenereateAdminKubeconfig(cert, key)
	if err != nil {
		return err
	}
	err = ApplyCertToSecret(caCert, caKey)
	if err != nil {
		return err
	}
	err = ApplyKubeConfigToConfigMap(kubeConfig)
	if err != nil {
		return err
	}
	return nil
}

func GenerateCA() (*x509.Certificate, crypto.Signer, error) {
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

func GenerateCert(caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, crypto.Signer, error) {
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

func GenereateAdminKubeconfig(cert *x509.Certificate, key crypto.Signer) (*api.Config, error) {
	certData := EncodeCertPEM(cert)
	keyData, err := keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal private key to PEM %s", err)
	}

	return &api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: map[string]*api.Cluster{
			"karbour": {
				Server:                "https://karbour-server.karbour.svc:7443",
				InsecureSkipTLSVerify: true,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			"karbour": {
				ClientCertificateData: certData,
				ClientKeyData:         keyData,
			},
		},
		Contexts: map[string]*api.Context{
			"default": {
				Cluster:  "karbour",
				AuthInfo: "karbour",
			},
		},
		CurrentContext: "default",
	}, nil
}

func ApplyCertToSecret(caCert *x509.Certificate, caKey crypto.Signer) error {
	caCertData := EncodeCertPEM(caCert)
	caKeyData, err := keyutil.MarshalPrivateKeyToPEM(caKey)
	if err != nil {
		return fmt.Errorf("unable to marshal private key to PEM %s", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "karbour-certification",
		},
		StringData: map[string]string{
			"ca.crt": string(caCertData),
			"ca.key": string(caKeyData),
		},
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return err
	}

	_, err = cs.CoreV1().Secrets(string(namespace)).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

func ApplyKubeConfigToConfigMap(kubeConfig *api.Config) error {
	kubeConfigData, err := yaml.Marshal(kubeConfig)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "karbour-kubeconfig",
		},
		Data: map[string]string{
			"config": string(kubeConfigData),
		},
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return err
	}

	_, err = cs.CoreV1().ConfigMaps(string(namespace)).Create(context.TODO(), cm, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}
