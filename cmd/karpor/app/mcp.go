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
	"fmt"

	_ "github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	_  "github.com/elastic/go-elasticsearch/v8" // esclient
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type mcpOptions struct {
	SSEPort string
	ElasticSearchAddresses []string
}

func NewMCPOptions() *mcpOptions {
	return &mcpOptions{}
}

func (o *mcpOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.SSEPort,"MCP SSE server exposure port", ":7999", "The address expossing the mcp server")
	fs.StringSliceVar(&o.ElasticSearchAddresses, "elastic-search-addresses", nil, "The elastic search address")
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

func mcpRun(ctx context.Context, options *mcpOptions) error {

	//TODO integrate MCP-GOLANG SSE Server functionality here
	// https://mcpgolang.com/introduction
	// rough outline:
	// - Logger setup
	// - registering the exposure points of the elastic search storage
	// - elastic search client setup :
	// 		- map from KusionStack/karpor/pkg/infra/search/storage/elasticsearch
	// - running the server
	//

	fmt.Println("TODO: yet to implement mcp functionality")
	fmt.Println("see /cmd/karpor/app/mcp.go for further directives")
	return nil
}
