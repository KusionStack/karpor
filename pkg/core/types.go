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

	"github.com/KusionStack/karbour/pkg/search/storage"
)

type LocatorType int

const (
	Cluster LocatorType = iota
	GVK
	Namespace
	ClusterGVKNamespace
	Resource
	NonNamespacedResource
)

type Locator struct {
	Cluster    string `json:"cluster" yaml:"cluster"`
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	Name       string `json:"name" yaml:"name"`
}

func LocatorFor(r *storage.Resource) Locator {
	if r.Cluster == "" {
		panic("cluster is empty")
	}

	return Locator{
		Cluster:    r.Cluster,
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Namespace:  r.Namespace,
		Name:       r.Name,
	}
}

func (c *Locator) ToSQL() string {
	var conditions []string

	if c.Cluster != "" {
		conditions = append(conditions, fmt.Sprintf("cluster='%s'", c.Cluster))
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

func (c *Locator) GetType() (LocatorType, bool) {
	if c.Cluster == "" {
		return -1, false
	}
	if c.APIVersion != "" && c.Kind != "" && c.Namespace != "" && c.Name != "" {
		return Resource, true
	} else if c.APIVersion != "" && c.Kind != "" && c.Name != "" {
		return NonNamespacedResource, true
	} else if c.APIVersion != "" && c.Kind != "" && c.Namespace != "" {
		return ClusterGVKNamespace, true
	} else if c.APIVersion != "" && c.Kind != "" {
		return GVK, true
	} else if c.Namespace != "" {
		// TODO: what if only apiversion is present but kind is not?
		return Namespace, true
	}
	return Cluster, true
}
