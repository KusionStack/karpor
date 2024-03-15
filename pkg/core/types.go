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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// LocatorType represents the type of a Locator.
type LocatorType int

// Enumerated constants representing different types of Locators.
const (
	CustomResourceGroup LocatorType = iota
	Cluster
	GVK
	Namespace
	ClusterGVKNamespace
	Resource
	NonNamespacedResource
)

// Locator represents information required to locate a resource.
type Locator struct {
	Cluster             string `json:"cluster" yaml:"cluster"`
	APIVersion          string `json:"apiVersion" yaml:"apiVersion"`
	Kind                string `json:"kind" yaml:"kind"`
	Namespace           string `json:"namespace" yaml:"namespace"`
	Name                string `json:"name" yaml:"name"`
	CustomResourceGroup string `json:"customResourceGroup" yaml:"customResourceGroup"`
}

// NewLocatorFromQuery creates a Locator from an HTTP request query parameters.
func NewLocatorFromQuery(r *http.Request) (Locator, error) {
	apiVersion := r.URL.Query().Get("apiVersion")
	if r.URL.RawPath != "" {
		apiVersion, _ = url.PathUnescape(apiVersion)
	}

	customResourceGroup := r.URL.Query().Get("customResourceGroup")
	if customResourceGroup != "" {
		crg, err := ParseCustomResourceGroup(customResourceGroup)
		if err != nil {
			return Locator{}, errors.Wrap(err, "failed to parse custom resource group")
		}
		// The custom resource group will be used as the key for caching, so it needs to be sorted to ensure uniqueness.
		customResourceGroup, err = SortCustomResourceGroup(crg)
		if err != nil {
			return Locator{}, errors.Wrap(err, "failed to sort custom resource group")
		}
	}

	cluster := r.URL.Query().Get("cluster")
	if customResourceGroup == "" && cluster == "" {
		return Locator{}, fmt.Errorf("cluster cannot be empty")
	}

	return Locator{
		CustomResourceGroup: customResourceGroup,
		APIVersion:          apiVersion,
		Cluster:             cluster,
		Kind:                r.URL.Query().Get("kind"),
		Namespace:           r.URL.Query().Get("namespace"),
		Name:                r.URL.Query().Get("name"),
	}, nil
}

// ToSQL generates a SQL query string based on the Locator.
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

// GetType returns the type of Locator and a boolean indicating success.
func (c *Locator) GetType() (LocatorType, bool) {
	if c.CustomResourceGroup != "" {
		return CustomResourceGroup, true
	}
	if c.Cluster == "" {
		return -1, false
	}
	if c.APIVersion != "" && c.Kind != "" && c.Namespace != "" && c.Name != "" {
		return Resource, true
	}
	if c.APIVersion != "" && c.Kind != "" && c.Name != "" {
		return NonNamespacedResource, true
	}
	if c.APIVersion != "" && c.Kind != "" && c.Namespace != "" {
		return ClusterGVKNamespace, true
	}
	if c.APIVersion != "" && c.Kind != "" {
		return GVK, true
	}
	if c.Namespace != "" {
		// TODO: what if only apiversion is present but kind is not?
		return Namespace, true
	}
	return Cluster, true
}

// ParseCustomResourceGroup deserialize the input JSON string to a map
func ParseCustomResourceGroup(jsonStr string) (map[string]any, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return nil, err
	}
	return m, nil
}

// SortCustomResourceGroup takes a map[string]interface{} and returns a JSON string with sorted keys.
func SortCustomResourceGroup(m map[string]any) (string, error) {
	// Extract all keys from the map
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort the keys
	sort.Strings(keys)

	// Create a sorted map
	sortedMap := make(map[string]any, len(m))
	for _, k := range keys {
		sortedMap[k] = m[k]
	}

	// Serialize the sorted map to a JSON string
	sortedJSON, err := json.Marshal(sortedMap)
	if err != nil {
		return "", err
	}

	return string(sortedJSON), nil
}
