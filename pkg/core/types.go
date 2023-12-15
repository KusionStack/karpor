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

package core

import (
	"fmt"
	"strings"
)

type Locator struct {
	Cluster    string `json:"cluster" yaml:"cluster"`
	Group      string `json:"group" yaml:"group"`
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	Name       string `json:"name" yaml:"name"`
}

func (c *Locator) ToSQL() string {
	var conditions []string

	if c.Cluster != "" {
		conditions = append(conditions, fmt.Sprintf("cluster='%s'", c.Cluster))
	}
	if c.Group != "" {
		conditions = append(conditions, fmt.Sprintf("group='%s'", c.Group))
	}
	if c.APIVersion != "" {
		conditions = append(conditions, fmt.Sprintf("apiVersion='%s'", c.APIVersion))
	}
	if c.Kind != "" {
		conditions = append(conditions, fmt.Sprintf("kind='%s'", c.Kind))
	}
	if c.Namespace != "" {
		conditions = append(conditions, fmt.Sprintf("namespace='%s'", c.Namespace))
	}
	if c.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name='%s'", c.Name))
	}

	if len(conditions) > 0 {
		return "SELECT * from resources WHERE " + strings.Join(conditions, " AND ")
	} else {
		return "SELECT * from resources"
	}
}
