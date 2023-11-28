package search

import (
	"context"
	"fmt"
	"sync"

	"github.com/KusionStack/karbour/pkg/search/storage"
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
