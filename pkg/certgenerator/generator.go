package certgenerator

import (
	"context"
	"crypto"
	"crypto/x509"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/keyutil"

	"sigs.k8s.io/yaml"
)

type Generator struct {
	clientSet      *kubernetes.Clientset
	namespace      string
	certName       string
	kubeConfigName string
}

func NewGenerator(clientSet *kubernetes.Clientset, namespace string, certName string, kubeConfigName string) (*Generator, error) {
	if clientSet == nil {
		return nil, fmt.Errorf("clientSet is required buit was nil")
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

	return &Generator{
		certName:       certName,
		kubeConfigName: kubeConfigName,
		clientSet:      clientSet,
		namespace:      namespace,
	}, nil
}

// Generate is a function that orchestrates the creation and application of certificates and kubeconfig necessary for a karbour sever.
func (g *Generator) Generate(ctx context.Context) error {
	exist, err := g.checkConfigExists()
	if err != nil {
		return err
	}
	if exist {
		// the config already exists, skipping generation
		return nil
	}

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

func (g *Generator) applyConfig(ctx context.Context, caCert *x509.Certificate, caKey crypto.Signer, kubeConfig *api.Config) error {
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

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: g.certName,
		},
		StringData: map[string]string{
			"ca.crt": string(caCertData),
			"ca.key": string(caKeyData),
		},
	}

	_, err = g.clientSet.CoreV1().Secrets(g.namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

func (g *Generator) applyKubeConfigToConfigMap(ctx context.Context, kubeConfig *api.Config) error {
	kubeConfigData, err := yaml.Marshal(kubeConfig)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: g.kubeConfigName,
		},
		Data: map[string]string{
			"config": string(kubeConfigData),
		},
	}

	_, err = g.clientSet.CoreV1().ConfigMaps(g.namespace).Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

// checkConfigExists determines if both the certificate and kubeconfig exist
func (g *Generator) checkConfigExists() (bool, error) {
	found1, err := g.checkCertExists()
	if err != nil {
		return false, err
	}
	found2, err := g.checkKubeConfigExists()
	if err != nil {
		return false, err
	}

	// return true if both the certificate and kubeconfig are found
	if found1 && found2 {
		return true, nil
	}
	return false, nil
}

func (g *Generator) checkKubeConfigExists() (bool, error) {
	_, err := g.clientSet.CoreV1().Secrets(g.namespace).Get(context.TODO(), g.certName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return false, err
		}
	}
	return true, nil
}

func (g *Generator) checkCertExists() (bool, error) {
	_, err := g.clientSet.CoreV1().ConfigMaps(g.namespace).Get(context.TODO(), g.kubeConfigName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return false, err
		}
	}
	return true, nil
}

func generateConfig() (*x509.Certificate, crypto.Signer, *api.Config, error) {
	caCert, caKey, err := generateCA()
	if err != nil {
		return nil, nil, nil, err
	}
	cert, key, err := generateCert(caCert, caKey)
	if err != nil {
		return nil, nil, nil, err
	}
	kubeConfig, err := genereateAdminKubeconfig(cert, key)
	if err != nil {
		return nil, nil, nil, err
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

func genereateAdminKubeconfig(cert *x509.Certificate, key crypto.Signer) (*api.Config, error) {
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
