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
