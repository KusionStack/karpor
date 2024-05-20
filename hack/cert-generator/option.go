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

package main

import (
	"github.com/spf13/pflag"
)

const (
	defaultCertName       = "karpor-certification"
	defaultKubeConfigName = "karpor-kubeconfig"
	inClusterNamespace    = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

type CertOptions struct {
	KubeConfig     string
	Namespace      string
	CertName       string
	KubeConfigName string
}

func NewCertOptions() *CertOptions {
	return &CertOptions{
		CertName:       defaultCertName,
		KubeConfigName: defaultKubeConfigName,
	}
}

func (o *CertOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.KubeConfig, "kubeconfig", o.KubeConfig, "The path of kubeconfig")
	fs.StringVar(&o.Namespace, "namespace", o.Namespace, "The namespace to store the CA and kubeconfig")
	fs.StringVar(&o.CertName, "ca-name", o.CertName, "The name of the secret used to store the CA certificate.")
	fs.StringVar(&o.KubeConfigName, "kubeconfig-name", o.KubeConfigName, "The name of the configmap used to store the kubeconfig.")
}
