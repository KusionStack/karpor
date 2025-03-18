/*
Copyright The Karpor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"

	"github.com/KusionStack/karpor/cmd/karpor/app"
)

// @title           Karpor
// @version         1.0
// @description     Karpor is a brand new Kubernetes visualization tool that focuses on search, insights, and AI at its core

func main() {
	ctx := genericapiserver.SetupSignalContext()

	cmd := app.NewServerCommand(ctx)
	syncCmd := app.NewSyncerCommand(ctx)
	agentCmd := app.NewAgentCommand(ctx)

	cmd.AddCommand(syncCmd)
	cmd.AddCommand(agentCmd)

	code := cli.Run(cmd)
	os.Exit(code)
}
