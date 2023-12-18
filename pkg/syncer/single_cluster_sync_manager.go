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
	"reflect"
	"sync"
	"sync/atomic"

	searchv1beta1 "github.com/KusionStack/karbour/pkg/apis/search/v1beta1"
	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/go-logr/logr"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

type SingleClusterSyncManager interface {
	Start(context.Context) error
	Started() bool
	Stop(context.Context)
	Stopped() bool
	UpdateSyncResources(context.Context, []*searchv1beta1.ResourceSyncRule) error
	HasSyncResource(schema.GroupVersionResource) bool
	ClusterConfig() *rest.Config
}

type singleClusterSyncManager struct {
	clusterName   string
	clusterConfig *rest.Config
	dynamicClient dynamic.Interface
	controller    controller.Controller

	ctx                 context.Context
	cancel              context.CancelFunc
	startLock           sync.RWMutex
	startOnce, stopOnce sync.Once
	started, stopped    bool
	wg                  wait.Group

	syncResources atomic.Value // map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule
	ch            chan struct{}
	syncers       sync.Map // map[schema.GroupVersionResource]*ResourceSyncer
	storage       storage.Storage

	logger logr.Logger
}

func NewSingleClusterSyncManager(baseContext context.Context,
	clusterName string,
	config *rest.Config,
	controller controller.Controller,
	storage storage.Storage,
) (SingleClusterSyncManager, error) {
	config = rest.CopyConfig(config)
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	innerCtx, innerCancel := context.WithCancel(baseContext)
	mgr := &singleClusterSyncManager{
		clusterName:   clusterName,
		clusterConfig: config,
		dynamicClient: dynamicClient,
		ctx:           innerCtx,
		cancel:        innerCancel,
		ch:            make(chan struct{}, 1),
		controller:    controller,
		storage:       storage,
		logger:        ctrl.LoggerFrom(baseContext).WithName("single-cluster-manager").WithValues("cluster", clusterName),
	}
	return mgr, nil
}

func (s *singleClusterSyncManager) Started() bool {
	s.startLock.RLock()
	defer s.startLock.RUnlock()
	return s.started
}

func (s *singleClusterSyncManager) Stopped() bool {
	s.startLock.RLock()
	defer s.startLock.RUnlock()
	return s.stopped
}

func (s *singleClusterSyncManager) Start(ctx context.Context) error {
	s.startOnce.Do(func() {
		s.logger.Info("start sync manager")

		go s.process()

		s.startLock.Lock()
		s.started = true
		s.startLock.Unlock()
	})
	return nil
}

func (s *singleClusterSyncManager) Stop(ctx context.Context) {
	s.stopOnce.Do(func() {
		s.logger.Info("start to stop the single cluster sync manager")
		defer s.logger.Info("single cluster sync manager was stopped")

		defer close(s.ch)

		s.startLock.Lock()
		s.stopped = true
		s.startLock.Unlock()

		s.cancel()

		s.logger.Info("waiting for resource syncers to stop")
		s.wg.Wait()
		s.logger.Info("all the resource syncers was stopped")
	})
}

func (s *singleClusterSyncManager) UpdateSyncResources(ctx context.Context, syncResources []*searchv1beta1.ResourceSyncRule) error {
	if s.Stopped() {
		return nil
	}

	byGVR := make(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	for _, r := range syncResources {
		gvr, err := parseGVR(r)
		if err != nil {
			return err
		}
		if _, exist := byGVR[gvr]; exist {
			return fmt.Errorf("found duplicate ResourceSyncRule defination for resource %q of cluster %q", gvr, s.clusterName)
		}
		byGVR[gvr] = r
	}

	s.syncResources.Store(byGVR)
	select {
	case s.ch <- struct{}{}:
	default:
	}
	return nil
}

func (s *singleClusterSyncManager) process() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.ch:
		}

		select {
		case <-s.ctx.Done():
			return
		default:
		}

		// TODO: retry on error
		if err := s.handleSyncResourcesUpdate(s.ctx); err != nil {
			s.logger.Error(err, "failed to update resources")
		}
	}
}

func (s *singleClusterSyncManager) handleSyncResourcesUpdate(ctx context.Context) error {
	var merr error

	desiredSyncResources := s.syncResources.Load().(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	for gvr, rsr := range desiredSyncResources {
		if s.Stopped() {
			return nil
		}

		syncer, exist := s.getSyncer(gvr)
		if exist && !reflect.DeepEqual(syncer.SyncRule(), rsr) {
			s.logger.Info("ResourceSyncRule has been updated", "grv", gvr)
			if err := s.stopResource(ctx, syncer); err != nil {
				merr = multierr.Append(merr, fmt.Errorf("error stopping syncing resource, gvr: %v, err: %v", gvr, err))
				continue
			}
			exist = false
		}
		if !exist {
			s.startResource(ctx, gvr, rsr)
		}
	}

	syncersToStop := make(map[schema.GroupVersionResource]*ResourceSyncer)
	s.syncers.Range(func(k, v any) bool {
		gvr := k.(schema.GroupVersionResource)
		if _, ok := desiredSyncResources[gvr]; !ok {
			syncersToStop[gvr] = v.(*ResourceSyncer)
		}
		return true
	})
	for _, syncer := range syncersToStop {
		s.stopResource(s.ctx, syncer)
	}
	return merr
}

func (s *singleClusterSyncManager) startResource(ctx context.Context, gvr schema.GroupVersionResource, rsr *searchv1beta1.ResourceSyncRule) {
	s.logger.Info("create resource syncer", "rsr", rsr)
	syncer := NewResourceSyncer(s.clusterName, s.dynamicClient, *rsr, s.storage)

	s.syncers.Store(gvr, syncer)
	s.controller.Watch(syncer.Source(), handler.Funcs{
		CreateFunc: func(ce event.CreateEvent, rli workqueue.RateLimitingInterface) {
			syncer.OnAdd(ce.Object)
		},
		UpdateFunc: func(ue event.UpdateEvent, rli workqueue.RateLimitingInterface) {
			syncer.OnUpdate(ue.ObjectNew)
		},
		DeleteFunc: func(de event.DeleteEvent, rli workqueue.RateLimitingInterface) {
			syncer.OnUpdate(de.Object)
		},
		GenericFunc: func(ge event.GenericEvent, rli workqueue.RateLimitingInterface) {
			syncer.OnGeneric(ge.Object)
		},
	})

	s.wg.StartWithContext(s.ctx, func(ctx context.Context) {
		if err := syncer.Run(ctx); err != nil {
			s.logger.Error(err, "failed to start syncer", "gvr", gvr)
		}
	})
}

func (s *singleClusterSyncManager) stopResource(ctx context.Context, syncer *ResourceSyncer) error {
	s.logger.Info("start to stop resource", "rsr", syncer.SyncRule())
	return syncer.Stop(ctx)
}

func (s *singleClusterSyncManager) ClusterConfig() *rest.Config {
	return rest.CopyConfig(s.clusterConfig)
}

func (s *singleClusterSyncManager) HasSyncResource(gvr schema.GroupVersionResource) bool {
	_, found := s.getSyncer(gvr)
	return found
}

func (s *singleClusterSyncManager) getSyncer(gvr schema.GroupVersionResource) (*ResourceSyncer, bool) {
	val, ok := s.syncers.Load(gvr)
	if !ok {
		return nil, false
	}
	return val.(*ResourceSyncer), true
}
