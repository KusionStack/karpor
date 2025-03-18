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

package app

import (
	"context"
	"crypto"
	"crypto/x509"

	"k8s.io/klog/v2"

	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"github.com/KusionStack/karpor/pkg/syncer"
	"github.com/KusionStack/karpor/pkg/util/certgenerator"

	esclient "github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

const (
	defaultCertFile = "/etc/karpor/ca/ca.crt"
	defaultKeyFile  = "/etc/karpor/ca/ca.key"
)

type syncerOptions struct {
	HighAvailability bool
	OnlyPushMode     bool

	MetricsAddr            string
	ProbeAddr              string
	ElasticSearchAddresses []string

	ExternalEndpoint string
	AgentImageTag    string

	CaCertFile string
	CaKeyFile  string
}

func NewSyncerOptions() *syncerOptions {
	return &syncerOptions{}
}

func (o *syncerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.HighAvailability, "high-availability", false, "Whether to use high-availability feature.")
	fs.BoolVar(&o.OnlyPushMode, "only-push-mode", false, "Only push mode in high availability feature.")

	fs.StringVar(&o.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	fs.StringVar(&o.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	fs.StringSliceVar(&o.ElasticSearchAddresses, "elastic-search-addresses", nil, "The elastic search address.")
	fs.StringVar(&o.ExternalEndpoint, "external-addresses", "", "The external address that expose to user cluster in pull mode.")
	fs.StringVar(&o.AgentImageTag, "agent-image-tag", "v0.0.0", "The agent image tag.")

	fs.StringVar(&o.CaCertFile, "ca-cert-file", defaultCertFile, "Root CA certificate file for karpor server.")
	fs.StringVar(&o.CaKeyFile, "ca-key-file", defaultKeyFile, "Root KEY file for karpor server..")
}

func NewSyncerCommand(ctx context.Context) *cobra.Command {
	options := NewSyncerOptions()
	cmd := &cobra.Command{
		Use:   "syncer",
		Short: "start a resource syncer to sync resource from clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSyncer(ctx, options)
		},
	}
	options.AddFlags(cmd.Flags())
	return cmd
}

func runSyncer(ctx context.Context, options *syncerOptions) error {
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
	//nolint:contextcheck
	es, err := elasticsearch.NewStorage(esclient.Config{
		Addresses: options.ElasticSearchAddresses,
	})
	if err != nil {
		log.Error(err, "unable to init elasticsearch client")
		return err
	}

	var caCert *x509.Certificate
	var caKey crypto.Signer
	if options.HighAvailability && !options.OnlyPushMode {
		caCert, caKey, err = certgenerator.LoadCertificate(options.CaCertFile, options.CaKeyFile)
		if err != nil {
			log.Error(err, "unable to load certificate")
			return err
		}
	}

	//nolint:contextcheck
	if err = syncer.NewSyncReconciler(es, options.HighAvailability, options.ElasticSearchAddresses, options.ExternalEndpoint, options.AgentImageTag, caCert, caKey).SetupWithManager(mgr); err != nil {
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
