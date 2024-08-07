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
	"context"
	"fmt"
	"os"

	"github.com/KusionStack/karpor/pkg/util/certgenerator"
	"github.com/KusionStack/karpor/pkg/version"
	"github.com/spf13/cobra"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/cli"
)

func main() {
	ctx := server.SetupSignalContext()
	command := NewCertGeneratorCommand(ctx)
	code := cli.Run(command)
	os.Exit(code)
}

func NewCertGeneratorCommand(ctx context.Context) *cobra.Command {
	options := NewCertOptions()
	cmd := &cobra.Command{
		Use:   "gen-cert",
		Short: "Generate CA and kubeconfig",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCertGenerator(ctx, options)
		},
	}
	options.AddFlags(cmd.Flags())
	return cmd
}

func runCertGenerator(ctx context.Context, options *CertOptions) error {
	var cfg *rest.Config
	var ns string
	var err error

	if options.Version {
		fmt.Println(version.GetVersion())
		return nil
	}

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
