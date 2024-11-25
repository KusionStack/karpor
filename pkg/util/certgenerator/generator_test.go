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

package certgenerator

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/keyutil"
)

func TestNewGenerator(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            *rest.Config
		namespace      string
		certName       string
		kubeConfigName string
		expectError    bool
	}{
		{
			name:           "Success",
			cfg:            &rest.Config{},
			namespace:      "test-ns",
			certName:       "test-cert",
			kubeConfigName: "test-kubeconfig",
			expectError:    false,
		},
		{
			name:           "Error - Nil Config",
			cfg:            nil,
			namespace:      "test-ns",
			certName:       "test-cert",
			kubeConfigName: "test-kubeconfig",
			expectError:    true,
		},
		{
			name:           "Error - Empty Namespace",
			cfg:            &rest.Config{},
			namespace:      "",
			certName:       "test-cert",
			kubeConfigName: "test-kubeconfig",
			expectError:    true,
		},
		{
			name:           "Error - Empty CertName",
			cfg:            &rest.Config{},
			namespace:      "test-ns",
			certName:       "",
			kubeConfigName: "test-kubeconfig",
			expectError:    true,
		},
		{
			name:           "Error - Empty KubeConfigName",
			cfg:            &rest.Config{},
			namespace:      "test-ns",
			certName:       "test-cert",
			kubeConfigName: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock kubernetes.NewForConfig if needed
			if tc.cfg != nil {
				mockey.Mock(kubernetes.NewForConfig).Return(&kubernetes.Clientset{}, nil).Build()
				defer mockey.UnPatchAll()
			}

			generator, err := NewGenerator(tc.cfg, tc.namespace, tc.certName, tc.kubeConfigName)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, generator)
			} else {
				require.NoError(t, err)
				require.NotNil(t, generator)
				require.Equal(t, tc.namespace, generator.namespace)
				require.Equal(t, tc.certName, generator.certName)
				require.Equal(t, tc.kubeConfigName, generator.kubeConfigName)
			}
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
	// Create a test generator
	generator := &Generator{
		clientSet:      &kubernetes.Clientset{},
		namespace:      "test-ns",
		certName:       "test-cert",
		kubeConfigName: "test-kubeconfig",
	}

	testCases := []struct {
		name        string
		mockSetup   func()
		expectError bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				// Mock generateConfig
				mockey.Mock(generateConfig).Return(&x509.Certificate{}, nil, "test-kubeconfig", nil).Build()
				// Mock applyConfig
				mockey.Mock((*Generator).applyConfig).Return(nil).Build()
			},
			expectError: false,
		},
		{
			name: "Error - Generate Config Failed",
			mockSetup: func() {
				mockey.Mock(generateConfig).Return(nil, nil, "", errors.New("generate config failed")).Build()
			},
			expectError: true,
		},
		{
			name: "Error - Apply Config Failed",
			mockSetup: func() {
				mockey.Mock(generateConfig).Return(&x509.Certificate{}, nil, "test-kubeconfig", nil).Build()
				mockey.Mock((*Generator).applyConfig).Return(errors.New("apply config failed")).Build()
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()
			defer mockey.UnPatchAll()

			err := generator.Generate(context.Background())
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenerator_applyConfig(t *testing.T) {
	// Create a test generator
	generator := &Generator{
		clientSet:      &kubernetes.Clientset{},
		namespace:      "test-ns",
		certName:       "test-cert",
		kubeConfigName: "test-kubeconfig",
	}

	testCases := []struct {
		name        string
		mockSetup   func()
		expectError bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				mockey.Mock((*Generator).applyCertToSecret).Return(nil).Build()
				mockey.Mock((*Generator).applyKubeConfigToConfigMap).Return(nil).Build()
			},
			expectError: false,
		},
		{
			name: "Error - Apply Cert Failed",
			mockSetup: func() {
				mockey.Mock((*Generator).applyCertToSecret).Return(errors.New("apply cert failed")).Build()
			},
			expectError: true,
		},
		{
			name: "Error - Apply KubeConfig Failed",
			mockSetup: func() {
				mockey.Mock((*Generator).applyCertToSecret).Return(nil).Build()
				mockey.Mock((*Generator).applyKubeConfigToConfigMap).Return(errors.New("apply kubeconfig failed")).Build()
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()
			defer mockey.UnPatchAll()

			err := generator.applyConfig(context.Background(), &x509.Certificate{}, nil, "test-kubeconfig")
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenerator_applyCertToSecret(t *testing.T) {
	// Create a test generator
	generator := &Generator{
		clientSet:      &kubernetes.Clientset{},
		namespace:      "test-ns",
		certName:       "test-cert",
		kubeConfigName: "test-kubeconfig",
	}

	testCases := []struct {
		name        string
		mockSetup   func()
		expectError bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				// Mock EncodeCertPEM
				mockey.Mock(EncodeCertPEM).Return([]byte("test-cert")).Build()
				// Mock keyutil.MarshalPrivateKeyToPEM
				mockey.Mock(keyutil.MarshalPrivateKeyToPEM).Return([]byte("test-key"), nil).Build()
				// Create fake clientset with pre-created secret
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-cert",
						Namespace: "test-ns",
					},
				}
				fakeClientset := fake.NewSimpleClientset(secret)
				mockey.Mock((*kubernetes.Clientset).CoreV1).Return(fakeClientset.CoreV1()).Build()
				// Mock Secret Apply method
				secretsClient := fakeClientset.CoreV1().Secrets(generator.namespace)
				mockey.Mock(secretsClient.Apply).Return(&corev1.Secret{}, nil).Build()
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()
			defer mockey.UnPatchAll()

			err := generator.applyCertToSecret(context.Background(), &x509.Certificate{}, nil)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenerator_applyKubeConfigToConfigMap(t *testing.T) {
	// Create a test generator
	generator := &Generator{
		clientSet:      &kubernetes.Clientset{},
		namespace:      "test-ns",
		certName:       "test-cert",
		kubeConfigName: "test-kubeconfig",
	}

	testCases := []struct {
		name        string
		mockSetup   func()
		expectError bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				// Create fake clientset with pre-created configmap
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kubeconfig",
						Namespace: "test-ns",
					},
				}
				fakeClientset := fake.NewSimpleClientset(configMap)
				mockey.Mock((*kubernetes.Clientset).CoreV1).Return(fakeClientset.CoreV1()).Build()
				// Mock ConfigMap Apply method
				configMapsClient := fakeClientset.CoreV1().ConfigMaps(generator.namespace)
				mockey.Mock(configMapsClient.Apply).Return(&corev1.ConfigMap{}, nil).Build()
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()
			defer mockey.UnPatchAll()

			err := generator.applyKubeConfigToConfigMap(context.Background(), "test-kubeconfig")
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenerateConfig(t *testing.T) {
	testCases := []struct {
		name        string
		namespace   string
		mockSetup   func()
		expectError bool
	}{
		{
			name:      "Success",
			namespace: "test-ns",
			mockSetup: func() {
				// Mock generateCA
				mockey.Mock(generateCA).Return(&x509.Certificate{}, &rsa.PrivateKey{}, nil).Build()
				// Mock generateCert
				mockey.Mock(generateCert).Return(&x509.Certificate{}, &rsa.PrivateKey{}, nil).Build()
				// Mock generateAdminKubeconfig
				mockey.Mock(generateAdminKubeconfig).Return("test-kubeconfig", nil).Build()
			},
			expectError: false,
		},
		{
			name:      "Error - generateCA Failed",
			namespace: "test-ns",
			mockSetup: func() {
				// Mock generateCA with error
				mockey.Mock(generateCA).Return(nil, nil, errors.New("generateCA failed")).Build()
			},
			expectError: true,
		},
		{
			name:      "Error - generateCert Failed",
			namespace: "test-ns",
			mockSetup: func() {
				// Mock generateCA
				mockey.Mock(generateCA).Return(&x509.Certificate{}, &rsa.PrivateKey{}, nil).Build()
				// Mock generateCert with error
				mockey.Mock(generateCert).Return(nil, nil, errors.New("generateCert failed")).Build()
			},
			expectError: true,
		},
		{
			name:      "Error - generateAdminKubeconfig Failed",
			namespace: "test-ns",
			mockSetup: func() {
				// Mock generateCA
				mockey.Mock(generateCA).Return(&x509.Certificate{}, &rsa.PrivateKey{}, nil).Build()
				// Mock generateCert
				mockey.Mock(generateCert).Return(&x509.Certificate{}, &rsa.PrivateKey{}, nil).Build()
				// Mock generateAdminKubeconfig with error
				mockey.Mock(generateAdminKubeconfig).Return("", errors.New("generateAdminKubeconfig failed")).Build()
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()
			defer mockey.UnPatchAll()

			caCert, caKey, kubeConfig, err := generateConfig(tc.namespace)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, caCert)
				require.Nil(t, caKey)
				require.Empty(t, kubeConfig)
			} else {
				require.NoError(t, err)
				require.NotNil(t, caCert)
				require.NotNil(t, caKey)
				require.NotEmpty(t, kubeConfig)
			}
		})
	}
}

func TestGenerateCA(t *testing.T) {
	testCases := []struct {
		name        string
		mockSetup   func()
		expectError bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				// Mock NewCaCertAndKey
				mockey.Mock(NewCaCertAndKey).Return(&x509.Certificate{}, &rsa.PrivateKey{}, nil).Build()
			},
			expectError: false,
		},
		{
			name: "Error - NewCaCertAndKey Failed",
			mockSetup: func() {
				// Mock NewCaCertAndKey with error
				mockey.Mock(NewCaCertAndKey).Return(nil, nil, errors.New("NewCaCertAndKey failed")).Build()
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()
			defer mockey.UnPatchAll()

			cert, key, err := generateCA()
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, cert)
				require.Nil(t, key)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cert)
				require.NotNil(t, key)
			}
		})
	}
}

func TestGenerateCert(t *testing.T) {
	testCases := []struct {
		name        string
		mockSetup   func()
		expectError bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				// Mock NewCaCertAndKeyFromRoot
				mockey.Mock(NewCaCertAndKeyFromRoot).Return(&x509.Certificate{}, &rsa.PrivateKey{}, nil).Build()
			},
			expectError: false,
		},
		{
			name: "Error - NewCaCertAndKeyFromRoot Failed",
			mockSetup: func() {
				// Mock NewCaCertAndKeyFromRoot with error
				mockey.Mock(NewCaCertAndKeyFromRoot).Return(nil, nil, errors.New("NewCaCertAndKeyFromRoot failed")).Build()
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()
			defer mockey.UnPatchAll()

			cert, key, err := generateCert(&x509.Certificate{}, &rsa.PrivateKey{})
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, cert)
				require.Nil(t, key)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cert)
				require.NotNil(t, key)
			}
		})
	}
}

func TestGenerateAdminKubeconfig(t *testing.T) {
	testCases := []struct {
		name        string
		namespace   string
		caCert      *x509.Certificate
		caKey       *rsa.PrivateKey
		mockSetup   func()
		expectError bool
	}{
		{
			name:      "Success",
			namespace: "test-ns",
			caCert:    &x509.Certificate{},
			caKey:     &rsa.PrivateKey{},
			mockSetup: func() {
				// Mock generateCert
				mockey.Mock(generateCert).Return(&x509.Certificate{}, &rsa.PrivateKey{}, nil).Build()
				// Mock EncodeCertPEM
				mockey.Mock(EncodeCertPEM).Return([]byte("test-cert")).Build()
				// Mock keyutil.MarshalPrivateKeyToPEM
				mockey.Mock(keyutil.MarshalPrivateKeyToPEM).Return([]byte("test-key"), nil).Build()
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()
			defer mockey.UnPatchAll()

			kubeconfig, err := generateAdminKubeconfig(tc.namespace, tc.caCert, tc.caKey)
			if tc.expectError {
				require.Error(t, err)
				require.Empty(t, kubeconfig)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, kubeconfig)
				// Verify kubeconfig format
				require.Contains(t, kubeconfig, "apiVersion: v1")
				require.Contains(t, kubeconfig, "kind: Config")
				require.Contains(t, kubeconfig, "clusters:")
				require.Contains(t, kubeconfig, "contexts:")
				require.Contains(t, kubeconfig, "users:")
			}
		})
	}
}
