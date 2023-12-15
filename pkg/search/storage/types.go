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
	"context"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
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

type Resource struct {
	Cluster    string                 `json:"cluster"`
	Namespace  string                 `json:"namespace"`
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Name       string                 `json:"name"`
	Object     map[string]interface{} `json:"object"`
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
		resYAML, err := yaml.Marshal(res)
		if err != nil {
			return "", err
		}

		yamlString += string(resYAML) + "\n---\n"
	}

	return yamlString, nil
}
