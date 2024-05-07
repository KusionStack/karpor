package main

import (
	"os"

	"github.com/KusionStack/karpor/cmd/app"
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
