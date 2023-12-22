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

package resource

import (
	"github.com/KusionStack/karbour/pkg/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type ResourceConfig struct {
	Verbose bool `json:"verbose"`
}

type ResourceSummary struct {
	Resource          core.Locator `json:"resource"`
	CreationTimestamp metav1.Time  `json:"creationTimestamp"`
	ResourceVersion   string       `json:"resourceVersion"`
	UID               types.UID    `json:"uid"`
}

type ResourceEvents struct {
	Resource       core.Locator `json:"resource"`
	Count          int          `json:"count"`
	Reason         string       `json:"reason"`
	Source         string       `json:"source"`
	Type           string       `json:"type"`
	LastTimestamp  metav1.Time  `json:"firstTimestamp"`
	FirstTimestamp metav1.Time  `json:"lastTimestamp"`
}

type ResourceTopology struct {
	Identifier string   `json:"identifier"`
	Parents    []string `json:"parents"`
	Children   []string `json:"children"`
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
