package search

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/client-go/rest"
)

type MultiClusterSyncManager interface {
	ForCluster(clusterName string, config *rest.Config) (SingleClusterSyncManager, error)
	Start(clusterName string) error
	Stop(clusterName string)
	GetSingleClusterSyncManager(clusterName string) (SingleClusterSyncManager, bool)
}

type multiClusterSyncManager struct {
	syncers map[string]SingleClusterSyncManager
	sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMultiClusterSyncManager(baseContext context.Context) MultiClusterSyncManager {
	innerCtx, innerCancel := context.WithCancel(baseContext)
	return &multiClusterSyncManager{
		syncers: make(map[string]SingleClusterSyncManager),
		ctx:     innerCtx,
		cancel:  innerCancel,
	}
}

func (s *multiClusterSyncManager) ForCluster(clusterName string, config *rest.Config) (SingleClusterSyncManager, error) {
	mgr, ok := s.GetSingleClusterSyncManager(clusterName)
	if ok {
		// already exist, just return
		return mgr, nil
	}

	mgr, err := NewSingleClusterSyncManager(s.ctx, clusterName, config)
	if err != nil {
		return nil, fmt.Errorf("failed to build sync manager for cluster %q: %v", clusterName, err)
	}
	s.Lock()
	defer s.Unlock()
	s.syncers[clusterName] = mgr
	return mgr, nil
}

func (s *multiClusterSyncManager) Start(clusterName string) error {
	mgr, ok := s.GetSingleClusterSyncManager(clusterName)
	if !ok {
		return fmt.Errorf("SingleClusterSyncManager for cluster %q does not exist", clusterName)
	}
	return mgr.Start()
}

func (s *multiClusterSyncManager) Stop(clusterName string) {
	mgr, found := s.GetSingleClusterSyncManager(clusterName)
	if !found {
		return
	}
	mgr.Stop()

	s.Lock()
	defer s.Unlock()
	delete(s.syncers, clusterName)
}

func (s *multiClusterSyncManager) GetSingleClusterSyncManager(clusterName string) (SingleClusterSyncManager, bool) {
	s.RLock()
	defer s.RUnlock()
	syncer, ok := s.syncers[clusterName]
	return syncer, ok
}
