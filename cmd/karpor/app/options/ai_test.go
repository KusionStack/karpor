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

package options

import (
	"testing"

	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestAIModelFlagDefault(t *testing.T) {
	opts := NewAIOptions()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	opts.AddFlags(fs)
	require.NoError(t, fs.Parse(nil))

	// gpt-3.5-turbo is shut down by OpenAI on 2026-10-23.
	require.Equal(t, "gpt-5.6-terra", opts.AIModel)
}

func TestAIModelDefaultIsKeptForOpenAI(t *testing.T) {
	opts := NewAIOptions()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	opts.AddFlags(fs)
	require.NoError(t, fs.Parse(nil))

	cfg := &registry.ExtraConfig{}
	require.NoError(t, opts.ApplyTo(cfg))
	require.Equal(t, "gpt-5.6-terra", cfg.AIModel)
}

func TestAIModelDefaultIsDroppedForOtherBackends(t *testing.T) {
	opts := NewAIOptions()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	opts.AddFlags(fs)
	require.NoError(t, fs.Parse([]string{"--ai-backend=huggingface"}))

	cfg := &registry.ExtraConfig{}
	require.NoError(t, opts.ApplyTo(cfg))
	require.Empty(t, cfg.AIModel)
}
