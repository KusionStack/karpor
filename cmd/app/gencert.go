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

package app

import (
	"context"
	"os"

	"github.com/KusionStack/karbour/pkg/certgenerator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultCertName       = "karbour-certification"
	defaultKubeConfigName = "karbour-kubeconfig"
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

func NewCertGeneratorCommand(ctx context.Context) *cobra.Command {
	options := NewCertOptions()
	cmd := &cobra.Command{
		Use:   "gen-cert",
		Short: "Generate CA and kubeconfig",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGertGenerator(ctx, options)
		},
	}
	options.AddFlags(cmd.Flags())
	return cmd
}

func runGertGenerator(ctx context.Context, options *CertOptions) error {
	var cfg *rest.Config
	var ns string
	var err error

	if options.KubeConfig == "" {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags("", options.KubeConfig)
		if err != nil {
			return err
		}
	}

	if options.Namespace == "" {
		var b []byte
		b, err = os.ReadFile(inClusterNamespace)
		if err != nil {
			return err
		}
		ns = string(b)
	} else {
		ns = options.Namespace
	}

	generator, err := certgenerator.NewGenerator(cfg, ns, options.CertName, options.KubeConfigName)
	if err != nil {
		return err
	}

	err = generator.Generate(ctx)
	if err != nil {
		return err
	}
	return nil
}
