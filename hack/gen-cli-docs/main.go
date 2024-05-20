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

package main

import (
	"os"

	"kusionstack.io/karpor/cmd/app"
	"github.com/spf13/cobra/doc"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

const (
	cliDocsDir = "docs/cli"
)

func main() {
	genCliDocs(cliDocsDir)
}

func genCliDocs(cliDir string) {
	if err := os.MkdirAll(cliDir, os.ModePerm); err != nil {
		klog.Fatalf("failed to create directory: %v", err)
	}

	ctx := genericapiserver.SetupSignalContext()
	cmd := app.NewServerCommand(ctx)
	syncCmd := app.NewSyncerCommand(ctx)
	cmd.AddCommand(syncCmd)

	if err := doc.GenMarkdownTree(cmd, cliDir); err != nil {
		klog.Fatal("failed to generate markdown document: %v", err)
	}
}
