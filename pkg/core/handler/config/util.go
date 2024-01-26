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

package config

import (
	"fmt"

	"github.com/KusionStack/karbour/pkg/kubernetes/registry"
)

func MaskSecretInConfig(extraConfig *registry.ExtraConfig) (*registry.ExtraConfig, error) {
	if extraConfig == nil {
		return nil, fmt.Errorf("extraConfig should not be empty")
	}
	maskedConfig := *extraConfig
	maskedConfig.ElasticSearchPassword = "redacted"
	return &maskedConfig, nil
}
