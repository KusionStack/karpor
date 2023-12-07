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
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/spf13/pflag"
)

type ElasticSearchConfig struct {
	Addresses []string
	UserName  string
	Password  string
}

type SearchStorageOptions struct {
	SearchStorageType      string
	ElasticSearchAddresses []string
	ElasticSearchName      string
	ElasticSearchPassword  string
}

func NewSearchStorageOptions() *SearchStorageOptions {
	return &SearchStorageOptions{}
}

func (o *SearchStorageOptions) Validate() []error {
	return nil
}

func (o *SearchStorageOptions) ApplyTo(config *registry.ExtraConfig) error {
	config.SearchStorageType = o.SearchStorageType
	config.ElasticSearchAddresses = o.ElasticSearchAddresses
	config.ElasticSearchName = o.ElasticSearchName
	config.ElasticSearchPassword = o.ElasticSearchPassword
	return nil
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *SearchStorageOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.SearchStorageType, "search-storage-type", "", "The search storage type")
	fs.StringSliceVar(&o.ElasticSearchAddresses, "elastic-search-addresses", nil, "The elastic search address")
	fs.StringVar(&o.ElasticSearchName, "elastic-search-username", "", "The elastic search username")
	fs.StringVar(&o.ElasticSearchPassword, "elastic-search-password", "", "The elastic search password")
}
