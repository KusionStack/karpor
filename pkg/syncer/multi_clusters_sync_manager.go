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

package syncer

import (
	"context"
	"fmt"
	"sync"

	"kusionstack.io/karpor/pkg/infra/search/storage"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

// MultiClusterSyncManager defines the interface for synchronizing across multiple clusters.
type MultiClusterSyncManager interface {
	Create(ctx context.Context, clusterName string, config *rest.Config) (SingleClusterSyncManager, error)
	Start(ctx context.Context, clusterName string) error
	Stop(ctx context.Context, clusterName string)
	GetForCluster(clusterName string) (SingleClusterSyncManager, bool)
}

// multiClusterSyncManager is the concrete implementation of the MultiClusterSyncManager interface.
type multiClusterSyncManager struct {
	storage    storage.ResourceStorage
	controller controller.Controller

	managers map[string]SingleClusterSyncManager
	sync.RWMutex
}

// NewMultiClusterSyncManager creates a new MultiClusterSyncManager instance with the given context, controller and storage.
func NewMultiClusterSyncManager(baseContext context.Context, controller controller.Controller, storage storage.ResourceStorage) MultiClusterSyncManager {
	return &multiClusterSyncManager{
		managers:   make(map[string]SingleClusterSyncManager),
		controller: controller,
		storage:    storage,
	}
}

// Create creates a SingleClusterSyncManager for the specified cluster using the provided context and configuration.
func (s *multiClusterSyncManager) Create(ctx context.Context, clusterName string, config *rest.Config) (SingleClusterSyncManager, error) {
	mgr, ok := s.GetForCluster(clusterName)
	if ok {
		// already exist, just return
		return mgr, nil
	}

	mgr, err := NewSingleClusterSyncManager(ctx, clusterName, config, s.controller, s.storage)
	if err != nil {
		return nil, err
	}

	s.Lock()
	defer s.Unlock()
	s.managers[clusterName] = mgr
	return mgr, nil
}

// Start starts the synchronization process for the specified cluster within the given context.
func (s *multiClusterSyncManager) Start(ctx context.Context, clusterName string) error {
	mgr, ok := s.GetForCluster(clusterName)
	if !ok {
		return fmt.Errorf("SingleClusterSyncManager for cluster %q does not exist", clusterName)
	}
	return mgr.Start(ctx)
}

// Stop stops the synchronization process for the specified cluster.
func (s *multiClusterSyncManager) Stop(ctx context.Context, clusterName string) {
	mgr, found := s.GetForCluster(clusterName)
	if !found {
		return
	}

	mgr.Stop(ctx)

	s.Lock()
	defer s.Unlock()
	delete(s.managers, clusterName)
}

// GetForCluster retrieves the SingleClusterSyncManager instance for the specified cluster if it exists.
func (s *multiClusterSyncManager) GetForCluster(clusterName string) (SingleClusterSyncManager, bool) {
	s.RLock()
	defer s.RUnlock()
	syncer, ok := s.managers[clusterName]
	return syncer, ok
}
