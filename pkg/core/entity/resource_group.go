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

package entity

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type (
	// ResourceGroupType represents the type of a ResourceGroup.
	ResourceGroupType int
	// ResourceGroupType represents the hash value of a ResourceGroup.
	ResourceGroupHash string
)

// Enumerated constants representing different types of ResourceGroups.
const (
	Cluster ResourceGroupType = iota
	GVK
	Namespace
	ClusterGVKNamespace
	Resource
	NonNamespacedResource
	Custom
)

// ResourceGroup represents information required to locate a resource or multi resources.
type ResourceGroup struct {
	Cluster     string            `json:"cluster,omitempty" yaml:"cluster,omitempty"`
	APIVersion  string            `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind        string            `json:"kind,omitempty" yaml:"kind,omitempty"`
	Namespace   string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// Hash returns a unique string representation of the ResourceGroup that can be
// used as a cache key.
func (rg *ResourceGroup) Hash() ResourceGroupHash {
	// Create a slice of keys from the Labels map and sort them to ensure
	// consistent ordering.
	labelKeys := make([]string, 0, len(rg.Labels))
	for k := range rg.Labels {
		labelKeys = append(labelKeys, k)
	}
	sort.Strings(labelKeys)

	// Create a slice of keys from the Annotations map and sort them to ensure
	// consistent ordering.
	annotationKeys := make([]string, 0, len(rg.Annotations))
	for k := range rg.Annotations {
		annotationKeys = append(annotationKeys, k)
	}
	sort.Strings(annotationKeys)

	// Create a hash from the sorted keys and values of Labels and Annotations.
	var hash strings.Builder
	hash.WriteString(rg.Cluster + "-")
	hash.WriteString(rg.APIVersion + "-")
	hash.WriteString(rg.Kind + "-")
	hash.WriteString(rg.Namespace + "-")
	hash.WriteString(rg.Name + "-")
	for _, k := range labelKeys {
		hash.WriteString(k + ":")
		hash.WriteString(rg.Labels[k] + "-")
	}
	for _, k := range annotationKeys {
		hash.WriteString(k + ":")
		hash.WriteString(rg.Annotations[k] + "-")
	}

	return ResourceGroupHash(hash.String())
}

// ToSQL generates a SQL query string based on the ResourceGroup.
func (rg *ResourceGroup) ToSQL() string {
	conditions := []string{}

	if rg.Cluster != "" {
		conditions = append(conditions, fmt.Sprintf("cluster='%s'", rg.Cluster))
	}
	if rg.APIVersion != "" {
		conditions = append(conditions, fmt.Sprintf("apiVersion='%s'", rg.APIVersion))
	}
	if rg.Kind != "" {
		conditions = append(conditions, fmt.Sprintf("kind='%s'", rg.Kind))
	}
	if rg.Namespace != "" {
		conditions = append(conditions, fmt.Sprintf("namespace='%s'", rg.Namespace))
	}
	if rg.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name='%s'", rg.Name))
	}
	for k, v := range rg.Annotations {
		conditions = append(conditions, fmt.Sprintf("`annotations.%s`='%s'", k, v))
	}
	for k, v := range rg.Labels {
		conditions = append(conditions, fmt.Sprintf("`labels.%s`='%s'", k, v))
	}

	if len(conditions) > 0 {
		return "SELECT * from resources WHERE " + strings.Join(conditions, " AND ")
	} else {
		return "SELECT * from resources"
	}
}

// GetType returns the type of ResourceGroup and a boolean indicating success.
func (rg *ResourceGroup) GetType() (ResourceGroupType, bool) {
	if rg.Cluster == "" || len(rg.Labels) != 0 || len(rg.Annotations) != 0 {
		return Custom, true
	}

	// Cluster is not empty
	if rg.APIVersion != "" && rg.Kind != "" && rg.Namespace != "" && rg.Name != "" {
		return Resource, true
	}
	if rg.APIVersion != "" && rg.Kind != "" && rg.Namespace == "" && rg.Name != "" {
		return NonNamespacedResource, true
	}
	if rg.APIVersion != "" && rg.Kind != "" && rg.Namespace != "" && rg.Name == "" {
		return ClusterGVKNamespace, true
	}
	if rg.APIVersion != "" && rg.Kind != "" && rg.Namespace == "" && rg.Name == "" {
		return GVK, true
	}
	if rg.APIVersion == "" && rg.Kind == "" && rg.Namespace != "" && rg.Name == "" {
		return Namespace, true
	}
	if rg.APIVersion == "" && rg.Kind == "" && rg.Namespace == "" && rg.Name == "" {
		return Cluster, true
	}
	return Custom, true
}

// NewResourceGroupFromQuery creates a ResourceGroup from an HTTP request query parameters.
//
// Examples:
// - url?apiVersion=v1&kind=Pod&labels=app.kubernetes.io/name=mockapp,env=prod
func NewResourceGroupFromQuery(r *http.Request) (ResourceGroup, error) {
	// Parse the query parameters.
	labelsRaw := r.URL.Query().Get("labels")
	annotationsRaw := r.URL.Query().Get("annotations")
	cluster := r.URL.Query().Get("cluster")
	apiVersion := r.URL.Query().Get("apiVersion")
	if r.URL.RawPath != "" {
		apiVersion, _ = url.PathUnescape(apiVersion)
	}

	// Convert the raw query parameters to maps.
	var labels map[string]string
	var annotations map[string]string
	if len(labelsRaw) > 0 {
		labels = make(map[string]string)
		// Each label is expected to be in the format "key=value".
		parts := strings.SplitN(labelsRaw, "=", 2)
		if len(parts) == 2 {
			labels[parts[0]] = parts[1]
		}
	}
	if len(annotationsRaw) > 0 {
		annotations = make(map[string]string)
		// Each annotation is expected to be in the format "key=value".
		parts := strings.SplitN(annotationsRaw, "=", 2)
		if len(parts) == 2 {
			annotations[parts[0]] = parts[1]
		}
	}

	// Construct a resource group instance.
	rg := ResourceGroup{
		Cluster:     cluster,
		APIVersion:  apiVersion,
		Kind:        r.URL.Query().Get("kind"),
		Namespace:   r.URL.Query().Get("namespace"),
		Name:        r.URL.Query().Get("name"),
		Labels:      labels,
		Annotations: annotations,
	}

	return rg, nil
}
