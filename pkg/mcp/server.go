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

// NewMCPStorageServer yields a storage initialized MCPServer struct
func NewMCPStorageServer(storageEndpoint storage.Storage, sseBaseURL string, sseAddr string) MCPStorageServer {
	mcpServer := server.NewMCPServer("MCP Storage Server", "0.0.1")
	sseServer := server.NewSSEServer(mcpServer, sseBaseURL)
	return MCPStorageServer{
		Storage:              storageEndpoint,
		MCPServer:            mcpServer,
		sseServerBaseURL:     sseBaseURL,
		sseServerBaseURLAddr: sseAddr,
		sseServer:            sseServer,
	}
}

// Serve starts an SSE server on the baseaddr provided
func (m *MCPStorageServer) Serve() error {
	return m.sseServer.Start(m.sseServerBaseURLAddr)
}
