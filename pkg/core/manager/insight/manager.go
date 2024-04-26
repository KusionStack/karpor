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

package insight

import (
	"time"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/scanner"
	"github.com/KusionStack/karpor/pkg/infra/scanner/kubeaudit"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/util/cache"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

type InsightManager struct {
	search                storage.SearchStorage
	resource              storage.ResourceStorage
	resourceGroupRule     storage.ResourceGroupRuleStorage
	scanner               scanner.KubeScanner
	scanCache             *cache.Cache[entity.ResourceGroupHash, scanner.ScanResult]
	clusterTopologyCache  *cache.Cache[entity.ResourceGroupHash, map[string]ClusterTopology]
	resourceTopologyCache *cache.Cache[entity.ResourceGroupHash, map[string]ResourceTopology]
	genericConfig         *genericapiserver.CompletedConfig
}

// NewInsightManager returns a new InsightManager object
func NewInsightManager(
	searchStorage storage.SearchStorage,
	resourceStorage storage.ResourceStorage,
	resourceGroupRuleStorage storage.ResourceGroupRuleStorage,
	genericConfig *genericapiserver.CompletedConfig,
) (*InsightManager, error) {
	const defaultExpiration = 10 * time.Minute

	// Create a new Kubernetes scanner instance.
	kubeauditScanner, err := kubeaudit.Default()
	if err != nil {
		return nil, err
	}

	return &InsightManager{
		search:                searchStorage,
		resource:              resourceStorage,
		resourceGroupRule:     resourceGroupRuleStorage,
		scanner:               kubeauditScanner,
		scanCache:             cache.NewCache[entity.ResourceGroupHash, scanner.ScanResult](defaultExpiration),
		clusterTopologyCache:  cache.NewCache[entity.ResourceGroupHash, map[string]ClusterTopology](defaultExpiration),
		resourceTopologyCache: cache.NewCache[entity.ResourceGroupHash, map[string]ResourceTopology](defaultExpiration),
		genericConfig:         genericConfig,
	}, nil
}
