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

package entity

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ResourceGroupType represents the type of a ResourceGroup.
type ResourceGroupType int

// Enumerated constants representing different types of ResourceGroups.
const (
	Cluster ResourceGroupType = iota
	GVK
	Namespace
	ClusterGVKNamespace
	Resource
	NonNamespacedResource
)

// ResourceGroup represents information required to locate a resource or multi resources.
type ResourceGroup struct {
	Cluster    string `json:"cluster" yaml:"cluster"`
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	Name       string `json:"name" yaml:"name"`
}

// NewResourceGroupFromQuery creates a ResourceGroup from an HTTP request query parameters.
func NewResourceGroupFromQuery(r *http.Request) (ResourceGroup, error) {
	cluster := r.URL.Query().Get("cluster")
	if cluster == "" {
		return ResourceGroup{}, fmt.Errorf("cluster cannot be empty")
	}

	apiVersion := r.URL.Query().Get("apiVersion")
	if r.URL.RawPath != "" {
		apiVersion, _ = url.PathUnescape(apiVersion)
	}

	return ResourceGroup{
		Cluster:    cluster,
		APIVersion: apiVersion,
		Kind:       r.URL.Query().Get("kind"),
		Namespace:  r.URL.Query().Get("namespace"),
		Name:       r.URL.Query().Get("name"),
	}, nil
}

// ToSQL generates a SQL query string based on the ResourceGroup.
func (c *ResourceGroup) ToSQL() string {
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

// GetType returns the type of ResourceGroup and a boolean indicating success.
func (c *ResourceGroup) GetType() (ResourceGroupType, bool) {
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
