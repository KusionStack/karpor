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
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/spf13/pflag"
)

type CoreOptions struct {
	ReadOnlyMode bool
	GithubBadge  bool
}

func NewCoreOptions() *CoreOptions {
	return &CoreOptions{}
}

func (o *CoreOptions) Validate() []error {
	return nil
}

func (o *CoreOptions) ApplyTo(config *registry.ExtraConfig) error {
	config.ReadOnlyMode = o.ReadOnlyMode
	config.GithubBadge = o.GithubBadge
	return nil
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *CoreOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.BoolVar(&o.ReadOnlyMode, "read-only-mode", false, "turn on the read only mode")
	fs.BoolVar(&o.GithubBadge, "github-badge", false, "whether to display the github badge")
}
