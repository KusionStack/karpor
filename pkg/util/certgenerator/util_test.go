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
	"crypto/x509"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPrivateKey(t *testing.T) {
	testCases := []struct {
		name        string
		keyType     x509.PublicKeyAlgorithm
		expectError bool
	}{
		{
			name:        "Success - RSA Key",
			keyType:     x509.UnknownPublicKeyAlgorithm,
			expectError: false,
		},
		{
			name:        "Success - ECDSA Key",
			keyType:     x509.ECDSA,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, err := NewPrivateKey(tc.keyType)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, key)
			} else {
				require.NoError(t, err)
				require.NotNil(t, key)
			}
		})
	}
}

func TestNewSelfSignedCACert(t *testing.T) {
	// Create a test private key first
	key, err := NewPrivateKey(x509.UnknownPublicKeyAlgorithm)
	require.NoError(t, err)

	testCases := []struct {
		name         string
		commonName   string
		organization []string
		year         time.Duration
		expectError  bool
	}{
		{
			name:         "Success",
			commonName:   "test-ca",
			organization: []string{"test-org"},
			year:         1,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cert, err := NewSelfSignedCACert(key, tc.commonName, tc.organization, tc.year)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, cert)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cert)
				require.Equal(t, tc.commonName, cert.Subject.CommonName)
				require.Equal(t, tc.organization, cert.Subject.Organization)
				require.True(t, cert.IsCA)
			}
		})
	}
}

func TestNewCaCertAndKey(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         Config
		expectError bool
	}{
		{
			name: "Success",
			cfg: Config{
				CommonName:   "test-ca",
				Organization: []string{"test-org"},
				Year:         1,
				AltNames: AltNames{
					DNSNames: map[string]string{"test": "test.example.com"},
					IPs:      map[string]net.IP{"test": net.ParseIP("127.0.0.1")},
				},
				Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cert, key, err := NewCaCertAndKey(tc.cfg)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, cert)
				require.Nil(t, key)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cert)
				require.NotNil(t, key)
				require.Equal(t, tc.cfg.CommonName, cert.Subject.CommonName)
				require.Equal(t, tc.cfg.Organization, cert.Subject.Organization)
				require.True(t, cert.IsCA)
			}
		})
	}
}

func TestNewCaCertAndKeyFromRoot(t *testing.T) {
	// Create root CA first
	rootKey, err := NewPrivateKey(x509.UnknownPublicKeyAlgorithm)
	require.NoError(t, err)
	rootCert, err := NewSelfSignedCACert(rootKey, "root-ca", []string{"root-org"}, 1)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		cfg         Config
		expectError bool
	}{
		{
			name: "Success",
			cfg: Config{
				CommonName:   "test-ca",
				Organization: []string{"test-org"},
				Year:         1,
				AltNames: AltNames{
					DNSNames: map[string]string{"test": "test.example.com"},
					IPs:      map[string]net.IP{"test": net.ParseIP("127.0.0.1")},
				},
				Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cert, key, err := NewCaCertAndKeyFromRoot(tc.cfg, rootCert, rootKey)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, cert)
				require.Nil(t, key)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cert)
				require.NotNil(t, key)
				require.Equal(t, tc.cfg.CommonName, cert.Subject.CommonName)
				require.Equal(t, tc.cfg.Organization, cert.Subject.Organization)
			}
		})
	}
}

func TestNewSignedCert(t *testing.T) {
	// Create CA first
	caKey, err := NewPrivateKey(x509.UnknownPublicKeyAlgorithm)
	require.NoError(t, err)
	caCert, err := NewSelfSignedCACert(caKey, "ca", []string{"org"}, 1)
	require.NoError(t, err)

	// Create key for the new cert
	key, err := NewPrivateKey(x509.UnknownPublicKeyAlgorithm)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		cfg         Config
		expectError bool
	}{
		{
			name: "Success",
			cfg: Config{
				CommonName:   "test-cert",
				Organization: []string{"test-org"},
				Year:         1,
				AltNames: AltNames{
					DNSNames: map[string]string{"test": "test.example.com"},
					IPs:      map[string]net.IP{"test": net.ParseIP("127.0.0.1")},
				},
				Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cert, err := NewSignedCert(tc.cfg, key, caCert, caKey)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, cert)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cert)
				require.Equal(t, tc.cfg.CommonName, cert.Subject.CommonName)
				require.Equal(t, tc.cfg.Organization, cert.Subject.Organization)
				require.False(t, cert.IsCA)
			}
		})
	}
}

func TestEncodeCertPEM(t *testing.T) {
	// Create a test certificate first
	key, err := NewPrivateKey(x509.UnknownPublicKeyAlgorithm)
	require.NoError(t, err)
	cert, err := NewSelfSignedCACert(key, "test-ca", []string{"test-org"}, 1)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		cert      *x509.Certificate
		expectNil bool
	}{
		{
			name:      "Success",
			cert:      cert,
			expectNil: false,
		},
		{
			name:      "Nil Certificate",
			cert:      nil,
			expectNil: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pemBytes := EncodeCertPEM(tc.cert)
			if tc.expectNil {
				require.Nil(t, pemBytes)
			} else {
				require.NotNil(t, pemBytes)
				require.Contains(t, string(pemBytes), "BEGIN CERTIFICATE")
				require.Contains(t, string(pemBytes), "END CERTIFICATE")
			}
		})
	}
}
