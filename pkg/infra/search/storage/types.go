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

	"github.com/KusionStack/karbour/pkg/core"
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

type Storage interface {
	Get(ctx context.Context, cluster string, obj runtime.Object) error
	Save(ctx context.Context, cluster string, obj runtime.Object) error
	Delete(ctx context.Context, cluster string, obj runtime.Object) error
}

type Query struct {
	Key      string
	Values   []string
	Operator string
}

type SearchStorage interface {
	Search(ctx context.Context, queryString, patternType string, pageSize, page int) (*SearchResult, error)
}

type SearchStorageGetter interface {
	GetSearchStorage() (SearchStorage, error)
}

type SearchResult struct {
	Total     int
	Resources []*Resource
}

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
