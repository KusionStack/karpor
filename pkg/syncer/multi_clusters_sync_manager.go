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

package syncer

import (
	"context"
	"fmt"
	"sync"

	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

type MultiClusterSyncManager interface {
	Create(ctx context.Context, clusterName string, config *rest.Config) (SingleClusterSyncManager, error)
	Start(ctx context.Context, clusterName string) error
	Stop(ctx context.Context, clusterName string)
	GetForCluster(clusterName string) (SingleClusterSyncManager, bool)
}

type multiClusterSyncManager struct {
	storage    storage.Storage
	controller controller.Controller

	managers map[string]SingleClusterSyncManager
	sync.RWMutex
}

func NewMultiClusterSyncManager(baseContext context.Context, controller controller.Controller, storage storage.Storage) MultiClusterSyncManager {
	return &multiClusterSyncManager{
		managers:   make(map[string]SingleClusterSyncManager),
		controller: controller,
		storage:    storage,
	}
}

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

func (s *multiClusterSyncManager) Start(ctx context.Context, clusterName string) error {
	mgr, ok := s.GetForCluster(clusterName)
	if !ok {
		return fmt.Errorf("SingleClusterSyncManager for cluster %q does not exist", clusterName)
	}
	return mgr.Start(ctx)
}

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

func (s *multiClusterSyncManager) GetForCluster(clusterName string) (SingleClusterSyncManager, bool) {
	s.RLock()
	defer s.RUnlock()
	syncer, ok := s.managers[clusterName]
	return syncer, ok
}
