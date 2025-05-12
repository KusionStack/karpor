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

	_ "github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	_ "github.com/KusionStack/karpor/pkg/mcp"
	_ "github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

type mcpOptions struct {
	SSEPort                string
	// TODO update mcpOptions to use the generic storage interface
	// should be able handle accept multiple storage backends
	// ElasticSearchAddresses []string
}

func NewMCPOptions() *mcpOptions {
	return &mcpOptions{}
}

func (o *mcpOptions) AddFlags(fs *pflag.FlagSet) {
	// TODO chart out how to handle multiple generic storage backends
	fs.StringVar(&o.SSEPort, "MCP SSE server exposure port", ":7999", "The address expossing the mcp server")
	// fs.StringSliceVar(&o.ElasticSearchAddresses, "elastic-search-addresses", nil, "The elastic search address")
}

func NewMCPCommand(ctx context.Context) *cobra.Command {
	options := NewMCPOptions()
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "start a storage mcp server to enable natural language interaction capabilities with the storage backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mcpRun(ctx, options)
		},
	}
	options.AddFlags(cmd.Flags())
	return cmd
}

// TODO update mcpOptions to use the generic storage interface
//nolint:unparam
func mcpRun(_ context.Context, options *mcpOptions) error {
	ctrl.SetLogger(klog.NewKlogr())
	log := ctrl.Log.WithName("mcp")

	//TODO update so that this receives generic storage backends (more than 1)
	log.Info("Starting MCP SSE server",
		"port", options.SSEPort, )


	//TODO update to use the generic storage interface for initialization
	//nolint:contextcheck
	// es, err := elasticsearch.NewStorage(esclient.Config{
	// 	Addresses: options.ElasticSearchAddresses,
	// })
	// if err != nil {
	// 	log.Error(err, "unable to init elasticsearch client")
	// 	return err
	// }
	// log.Info("Acquired elasticsearch storage backend", "esStorage", es)


	//TODO pickup syncer operations patterns for running the mcp server from app/syncer.go

	log.Info("TODO: yet to implement mcp functionality")
	log.Info("see /cmd/karpor/app/mcp.go for further directives")


	return nil
}
