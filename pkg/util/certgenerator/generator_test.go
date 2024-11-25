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
	"crypto/x509"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	testcore "k8s.io/client-go/kubernetes/fake"
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
				fakeClientset := testcore.NewSimpleClientset(secret)
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
				fakeClientset := testcore.NewSimpleClientset(configMap)
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
