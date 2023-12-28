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

package search

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type SearchManager struct{}

// NewSearchManager returns a new SearchManager object
func NewSearchManager() *SearchManager {
	return &SearchManager{}
}

type UniResource struct {
	Cluster string         `json:"cluster"`
	Object  runtime.Object `json:"object"`
}

type UniResourceList struct {
	metav1.TypeMeta
	Items       []UniResource `json:"items"`
	Total       int           `json:"total"`
	CurrentPage int           `json:"currentPage"`
	PageSize    int           `json:"pageSize"`
}
