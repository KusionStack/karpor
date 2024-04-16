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

package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/core/entity"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	Equals         = "="
	DSLPatternType = "dsl"
	SQLPatternType = "sql"
)

// ResourceStorage interface defines the basic operations for storage.
type Storage interface {
	ResourceStorage
	ResourceGroupRuleStorage
	SearchStorage
}

// ResourceStorage interface defines the basic operations for resource storage.
type ResourceStorage interface {
	GetResource(ctx context.Context, cluster string, obj runtime.Object) error
	SaveResource(ctx context.Context, cluster string, obj runtime.Object) error
	DeleteResource(ctx context.Context, cluster string, obj runtime.Object) error
	DeleteAllResources(ctx context.Context, cluster string) error
}

// ResourceGroupRuleStorage interface defines the basic operations for resource
// group rule storage.
type ResourceGroupRuleStorage interface {
	GetResourceGroupRule(ctx context.Context, name string) (*entity.ResourceGroupRule, error)
	SaveResourceGroupRule(ctx context.Context, data *entity.ResourceGroupRule) error
	DeleteResourceGroupRule(ctx context.Context, name string) error
	ListResourceGroupRules(ctx context.Context) ([]*entity.ResourceGroupRule, error)
}

// Storage interface defines the basic operations for resource storage.
type SearchStorage interface {
	Search(ctx context.Context, queryString, patternType string, pagination *Pagination) (*SearchResult, error)
}

type SearchStorageGetter interface {
	GetSearchStorage() (SearchStorage, error)
}

type ResourceGroupRuleStorageGetter interface {
	GetResourceGroupRuleStorage() (ResourceGroupRuleStorage, error)
}

// Query represents the query parameters for searching resources.
type Query struct {
	Key      string
	Values   []string
	Operator string
}

// Pagination defines the parameters for pagination in search results.
type Pagination struct {
	Page     int
	PageSize int
}

// SearchResult contains the search results and total count.
type SearchResult struct {
	Total     int
	Resources []*Resource
}

// AggregateResults is assumed to be a struct that holds aggregation results.
type AggregateResults struct {
	Buckets []Bucket
	Total   int
}

// Bucket is assumed to be a struct that holds individual bucket data.
type Bucket struct {
	Keys  []string
	Count int
}

// Overview returns a brief summary of the search result.
func (r *SearchResult) Overview() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total: %d\n", r.Total))
	sb.WriteString("Resources:\n")
	for _, res := range r.Resources {
		sb.WriteString(fmt.Sprintf("- Cluster: %s, Namespace: %s, Kind: %s, Name: %s\n",
			res.Cluster, res.Namespace, res.Kind, res.Name))
	}

	return sb.String()
}

// ToYAML returns the search result in YAML format.
func (r *SearchResult) ToYAML() (string, error) {
	if len(r.Resources) == 0 {
		return "", nil
	}

	var yamlString string
	for _, res := range r.Resources {
		resYAML, err := yaml.Marshal(res.Object)
		if err != nil {
			return "", err
		}

		yamlString += string(resYAML) + "\n---\n"
	}

	return yamlString, nil
}

// Resource represents a Kubernetes resource with additional metadata.
type Resource struct {
	core.Locator `json:",inline" yaml:",inline"`
	Object       map[string]interface{} `json:"object"`
}

// NewResource creates a new Resource instance based on the provided bytes
// and cluster. It decodes the YAML bytes into an unstructured object and
// constructs a Resource with the relevant fields.
func NewResource(cluster string, b []byte) (*Resource, error) {
	// Ensure the cluster name is not empty.
	if len(cluster) == 0 {
		return nil, fmt.Errorf("cluster cannot be empty")
	}

	// Initialize an unstructured object for decoding data.
	obj := &unstructured.Unstructured{}

	// Create a YAML or JSON decoder with the provided bytes.
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), len(b))
	if err := decoder.Decode(obj); err != nil {
		return nil, err
	}

	// Build and return the Resource object with decoded data and cluster info.
	return &Resource{
		Locator: core.Locator{
			Cluster:    cluster,
			Namespace:  obj.GetNamespace(),
			APIVersion: obj.GetAPIVersion(),
			Kind:       obj.GetKind(),
			Name:       obj.GetName(),
		},
		Object: obj.Object,
	}, nil
}

// Map2Resource converts a map to a Resource object.
func Map2Resource(in map[string]interface{}) (*Resource, error) {
	out := &Resource{}
	out.Cluster = in["cluster"].(string)
	out.APIVersion = in["apiVersion"].(string)
	out.Kind = in["kind"].(string)
	out.Namespace = in["namespace"].(string)
	out.Name = in["name"].(string)

	content := in["content"].(string)
	obj := &unstructured.Unstructured{}
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewBufferString(content), len(content))
	if err := decoder.Decode(obj); err != nil {
		return nil, err
	}
	out.Object = obj.Object
	return out, nil
}

// Map2ResourceGroupRule converts a map to a ResourceGroupRule object.
func Map2ResourceGroupRule(in map[string]interface{}) (*entity.ResourceGroupRule, error) {
	out := &entity.ResourceGroupRule{}
	out.ID = in["id"].(string)
	out.Name = in["name"].(string)
	out.Description = in["description"].(string)
	out.Fields = in["fields"].([]string)
	out.CreatedAt = in["createdAt"].(time.Time)
	out.DeletedAt = in["deletedAt"].(time.Time)
	out.UpdatedAt = in["updatedAt"].(time.Time)
	return out, nil
}
