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

package mcp

import (
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	_ "github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	_ "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	_ "sigs.k8s.io/controller-runtime"
)

// NewMCPStorageServer yields a new MCPStorageServer configuration
func NewMCPStorageServer(storageEndpoints []storage.Storage, sseBaseURL string, sseAddr string) MCPStorageServer {
	mcpServer := server.NewMCPServer("MCP Storage Server", "0.0.1")
	sseServer := server.NewSSEServer(mcpServer, sseBaseURL)
	return MCPStorageServer{
		Storages:             storageEndpoints,
		MCPServer:            mcpServer,
		sseServerBaseURL:     sseBaseURL,
		sseServerBaseURLAddr: sseAddr,
		sseServer:            sseServer,
	}
}


// TODO Modular design : storage backends should be pluggable
// design MCPStorage interface that should be satisfied by the tools
// all resource, tool, prompt handlers should use this MCPStorage interface
// Read-only access  to the storage backends needs to be configured
// the server need not

// for exposing a storage backend (DB queries): will be registering them as resources
// https://github.com/mark3labs/mcp-go/blob/main/README.md#resources

// Tools are supposed to have side effects and are not needed as of now in our case
// Prompts can be integrated later on, when patterns emerge

// Serve starts an SSE server on the baseaddr provided
// TODO integrate this Serve into the main cmd
func (m *MCPStorageServer) Serve() error {
	return m.sseServer.Start(m.sseServerBaseURLAddr)
}
