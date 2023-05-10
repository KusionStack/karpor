// Copyright The Karbour Authors.
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
	"github.com/KusionStack/karbour/pkg/apiserver"
	"github.com/spf13/pflag"
)

type StaticOptions struct {
	StaticDirectory string
}

func NewStaticOptions() *StaticOptions {
	return &StaticOptions{}
}

func (o *StaticOptions) Validate() []error {
	return nil
}

func (o *StaticOptions) ApplyTo(config *apiserver.ExtraConfig) error {
	config.StaticDirectory = o.StaticDirectory
	return nil
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *StaticOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.StaticDirectory, "static-dir", "./static", "static directory for web")
}
