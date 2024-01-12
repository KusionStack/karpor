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

	"github.com/KusionStack/karbour/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karbour/pkg/kubernetes/scheme"
	"github.com/KusionStack/karbour/pkg/syncer"
	esclient "github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

type syncerOptions struct {
	MetricsAddr            string
	ProbeAddr              string
	ElasticSearchAddresses []string
}

func NewSyncerOptions() *syncerOptions {
	return &syncerOptions{}
}

func (o *syncerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	fs.StringVar(&o.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	fs.StringSliceVar(&o.ElasticSearchAddresses, "elastic-search-addresses", nil, "The elastic search address.")
}

func NewSyncerCommand(ctx context.Context) *cobra.Command {
	options := NewSyncerOptions()
	cmd := &cobra.Command{
		Use:   "syncer",
		Short: "start a resource syncer to sync resource from clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(ctx, options)
		},
	}
	options.AddFlags(cmd.Flags())
	return cmd
}

func run(ctx context.Context, options *syncerOptions) error {
	ctrl.SetLogger(klog.NewKlogr())
	log := ctrl.Log.WithName("setup")

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
	es, err := elasticsearch.NewESClient(esclient.Config{
		Addresses: options.ElasticSearchAddresses,
	})
	if err != nil {
		log.Error(err, "unable to init elasticsearch client")
		return err
	}

	if err = syncer.NewSyncReconciler(es).SetupWithManager(mgr); err != nil {
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
