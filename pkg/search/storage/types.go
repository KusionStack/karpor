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

	"k8s.io/apimachinery/pkg/runtime"
)

const (
	Equals         = "="
	DSLPatternType = "dsl"
	SQLPatternType = "sql"
)

type Storage interface {
	Get(ctx context.Context, cluster string, obj runtime.Object) error
	Create(ctx context.Context, cluster string, obj runtime.Object) error
	Update(ctx context.Context, cluster string, obj runtime.Object) error
	Delete(ctx context.Context, cluster string, obj runtime.Object) error
}

type Query struct {
	Key      string
	Values   []string
	Operator string
}

type SearchStorage interface {
	Search(ctx context.Context, queryString, patternType string) (*SearchResult, error)
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
