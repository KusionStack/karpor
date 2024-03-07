package app

import (
	"context"
	"os"

	"github.com/KusionStack/karbour/pkg/certgenerator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	fs.StringVar(&o.KubeConfigName, "kubeconfig-name", o.CertName, "The name of the configmap used to store the kubeconfig.")
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
	var cs *kubernetes.Clientset
	var ns string
	var err error
	if options.KubeConfig == "" && options.Namespace == "" {
		cs, ns, err = getInClusterConfig()
		if err != nil {
			return err
		}
	}

	generator, err := certgenerator.NewGenerator(cs, ns, options.CertName, options.KubeConfigName)
	if err != nil {
		return err
	}

	err = generator.Generate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func getInClusterConfig() (*kubernetes.Clientset, string, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, "", err
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, "", err
	}

	namespace, err := os.ReadFile(inClusterNamespace)
	if err != nil {
		return nil, "", err
	}
	return cs, string(namespace), nil
}
