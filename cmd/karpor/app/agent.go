package app

import (
	"context"

	esclient "github.com/elastic/go-elasticsearch/v8"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	clusterv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"github.com/KusionStack/karpor/pkg/syncer"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
)

type agentOptions struct {
	syncerOptions
	ClusterName string
	ClusterMode string
}

func NewAgentOptions() *agentOptions {
	return &agentOptions{
		syncerOptions: *NewSyncerOptions(),
	}
}

func (o *agentOptions) AddFlags(fs *pflag.FlagSet) {
	o.syncerOptions.AddFlags(fs)
	fs.StringVar(&o.ClusterName, "cluster-name", "", "The cluster name in hub cluster.")
	fs.StringVar(&o.ClusterMode, "cluster-mode", "pull", "The cluster mode.")
}

func NewAgentCommand(ctx context.Context) *cobra.Command {
	options := NewAgentOptions()
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "start a resource syncer agent which deployed in user cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			// use the same logical as the Non-HA syncer in controller cluster
			return runAgent(ctx, options)
		},
	}
	options.AddFlags(cmd.Flags())
	return cmd
}

func runAgent(ctx context.Context, options *agentOptions) error {
	ctrl.SetLogger(klogr.New())
	log := ctrl.Log.WithName("setup")

	if options.ClusterMode == clusterv1beta1.PushClusterMode {
		// apply crds
		dynamicClient, err := dynamic.NewForConfig(ctrl.GetConfigOrDie())
		if err != nil {
			return errors.Wrapf(err, "failed to build dynamic client for ageng")
		}
		err = utils.ApplyCrds(ctx, dynamicClient)
		if err != nil {
			return err
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme.Scheme,
		MetricsBindAddress:     options.MetricsAddr,
		HealthProbeBindAddress: options.ProbeAddr,
	})
	if err != nil {
		log.Error(err, "unable to start manager")
		return err
	}

	// TODO: add startup parameters to change the type of storage
	//nolint:contextcheck
	es, err := elasticsearch.NewStorage(esclient.Config{
		Addresses: options.ElasticSearchAddresses,
	})
	if err != nil {
		log.Error(err, "unable to init elasticsearch client")
		return err
	}

	//nolint:contextcheck
	if err = syncer.NewAgentReconciler(es, options.HighAvailability, options.ClusterName).SetupWithManager(mgr); err != nil {
		log.Error(err, "unable to create resource syncer")
		return err
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Error(err, "unable to set up health check")
		return err
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.Error(err, "unable to set up ready check")
		return err
	}

	log.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		log.Error(err, "problem running manager")
		return err
	}

	return nil
}
