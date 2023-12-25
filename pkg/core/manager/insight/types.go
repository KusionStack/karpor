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

package insight

import (
	"time"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/util/cache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type InsightConfig struct {
	Verbose bool `json:"verbose"`
}

type InsightManager struct {
	config                *InsightConfig
	clusterTopologyCache  *cache.Cache[core.Locator, map[string]ClusterTopology]
	resourceTopologyCache *cache.Cache[core.Locator, map[string]ResourceTopology]
}

// NewInsightManager returns a new InsightManager object
func NewInsightManager(config *InsightConfig) *InsightManager {
	return &InsightManager{
		config:                config,
		clusterTopologyCache:  cache.NewCache[core.Locator, map[string]ClusterTopology](10 * time.Minute),
		resourceTopologyCache: cache.NewCache[core.Locator, map[string]ResourceTopology](10 * time.Minute),
	}
}

// Resource-related

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

// Cluster-related

type ClusterTopology struct {
	GroupVersionKind string            `json:"groupVersionKind"`
	Count            int               `json:"count"`
	Relationship     map[string]string `json:"relationship"`
}

type ClusterDetail struct {
	NodeCount      int    `json:"nodeCount"`
	ServerVersion  string `json:"serverVersion"`
	MemoryCapacity int64  `json:"memoryCapacity"`
	CPUCapacity    int64  `json:"cpuCapacity"`
	PodsCapacity   int64  `json:"podsCapacity"`
}
