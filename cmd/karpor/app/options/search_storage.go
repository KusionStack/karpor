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
	"encoding/json"

	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
	"github.com/spf13/pflag"
)

type ElasticSearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

type SearchStorageOptions struct {
	SearchStorageType string
	SearchAddresses   []string
	SearchUsername    string
	SearchPassword    string
}

func NewSearchStorageOptions() *SearchStorageOptions {
	return &SearchStorageOptions{}
}

func (o *SearchStorageOptions) Validate() []error {
	return nil
}

func (o *SearchStorageOptions) ApplyTo(config *registry.ExtraConfig) error {
	config.SearchStorageType = o.SearchStorageType
	config.SearchAddresses = o.SearchAddresses
	config.SearchUsername = o.SearchUsername
	config.SearchPassword = o.SearchPassword
	return nil
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *SearchStorageOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.SearchStorageType, "search-storage-type", "", "The search storage type")
	fs.StringSliceVar(&o.SearchAddresses, "search-addresses", nil, "The search address")
	fs.StringVar(&o.SearchUsername, "search-username", "", "The search username")
	fs.StringVar(&o.SearchPassword, "search-password", "", "The search password")
}

// MarshalJSON is custom marshalling function for masking sensitive field values
func (o SearchStorageOptions) MarshalJSON() ([]byte, error) {
	type tempOptions SearchStorageOptions
	o2 := tempOptions(o)
	if o2.SearchPassword != "" {
		o2.SearchPassword = MaskString
	}
	return json.Marshal(&o2)
}
