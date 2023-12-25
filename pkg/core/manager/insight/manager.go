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
	"github.com/KusionStack/karbour/pkg/infra/scanner"
	"github.com/KusionStack/karbour/pkg/infra/scanner/kubeaudit"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/util/cache"
)

type InsightManager struct {
	search                storage.SearchStorage
	scanner               scanner.KubeScanner
	scanCache             *cache.Cache[core.Locator, scanner.ScanResult]
	clusterTopologyCache  *cache.Cache[core.Locator, map[string]ClusterTopology]
	resourceTopologyCache *cache.Cache[core.Locator, map[string]ResourceTopology]
}

// NewInsightManager returns a new InsightManager object
func NewInsightManager(searchStorage storage.SearchStorage) (*InsightManager, error) {
	const defaultExpiration = 10 * time.Minute

	// Create a new Kubernetes scanner instance.
	kubeauditScanner, err := kubeaudit.Default()
	if err != nil {
		return nil, err
	}

	return &InsightManager{
		scanner:               kubeauditScanner,
		search:                searchStorage,
		scanCache:             cache.NewCache[core.Locator, scanner.ScanResult](defaultExpiration),
		clusterTopologyCache:  cache.NewCache[core.Locator, map[string]ClusterTopology](defaultExpiration),
		resourceTopologyCache: cache.NewCache[core.Locator, map[string]ResourceTopology](defaultExpiration),
	}, nil
}
