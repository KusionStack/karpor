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
)

// MCPServer is the main object that holds the necessary fields and components for the mcp server component to expose the storage backend via resources, tools and prompts
type MCPStorageServer struct {
	Storage              storage.Storage  //the storage backend to be exposed by the MCPStorageServer
	MCPServer            *server.MCPServer //manages configurations of tools, resources, prompts and their handlers
	sseServerBaseURL     string           //baseURL on which the sse server is exposed
	sseServerBaseURLAddr string           //port on which the sse server is exposed
	sseServer            *server.SSEServer //manages configuration of the HTTP based SSE server on which the MCPServer is exposed
}

// MCPToolName is a string tag for an mcp server tool
type MCPToolName string

// MCPResourceName is a string tag for an mcp server resource
type MCPResourceName string

// MCPPromptName is a string tag for an mcp server prompt
type MCPPromptName string
