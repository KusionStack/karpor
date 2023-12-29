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
	"time"

	"github.com/KusionStack/karbour/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/syncer/cache"
	"github.com/go-logr/logr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const defaultWorkers = 10

type syncCache struct {
	objs map[string]client.Object
	sync.RWMutex
}

func (c *syncCache) Save(obj client.Object) bool {
	c.Lock()
	defer c.Unlock()

	key, _ := clientgocache.MetaNamespaceKeyFunc(obj)
	cached, ok := c.objs[key]
	if ok {
		// only override if resource version is newer
		compare, _ := cache.CompareResourverVersion(obj.GetResourceVersion(), cached.GetResourceVersion())
		if compare <= 0 {
			return false
		}
	}
	c.objs[key] = obj
	return true
}

func (c *syncCache) Remove(obj client.Object) {
	c.Lock()
	defer c.Unlock()

	key, _ := clientgocache.MetaNamespaceKeyFunc(obj)
	cached, ok := c.objs[key]
	if !ok {
		return
	}
	if obj.GetResourceVersion() == cached.GetResourceVersion() {
		delete(c.objs, key)
	}
}

func (c *syncCache) Get(key string) (client.Object, bool) {
	c.RLock()
	defer c.RUnlock()
	obj, ok := c.objs[key]
	return obj, ok
}

type deleted struct {
	client.Object
}

type ResourceSyncer struct {
	source  SyncSource
	storage storage.Storage

	queue workqueue.RateLimitingInterface
	// resources that being processed or to be processed, but only keep the latest resource version
	pendingResources *syncCache

	ctx    context.Context
	cancel context.CancelFunc

	logger logr.Logger
}

func NewResourceSyncer(cluster string, dynamicClient dynamic.Interface, rsr v1beta1.ResourceSyncRule, storage storage.Storage) *ResourceSyncer {
	source := NewSource(cluster, dynamicClient, rsr)
	return &ResourceSyncer{
		source:           source,
		storage:          storage,
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), fmt.Sprintf("%s/%s-sync-queue", rsr.APIVersion, rsr.Resource)),
		pendingResources: &syncCache{objs: make(map[string]client.Object)},
		logger:           ctrl.Log.WithName(fmt.Sprintf("%s-syncer", source.SyncRule().Resource)),
	}
}

func (s *ResourceSyncer) Source() SyncSource {
	return s.source
}

func (s *ResourceSyncer) SyncRule() v1beta1.ResourceSyncRule {
	return s.source.SyncRule()
}

func (s *ResourceSyncer) Stop(ctx context.Context) error {
	if err := s.source.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop the source: %v", err)
	}
	s.cancel()
	return nil
}

func (s *ResourceSyncer) OnAdd(obj client.Object) {
	s.enqueue(obj)
}

func (s *ResourceSyncer) OnUpdate(obj client.Object) {
	s.enqueue(obj)
}

func (s *ResourceSyncer) OnDelete(obj client.Object) {
	s.enqueue(deleted{Object: obj})
}

func (s *ResourceSyncer) OnGeneric(obj client.Object) {
	s.enqueue(obj)
}

func (s *ResourceSyncer) enqueue(obj client.Object) {
	if s.pendingResources.Save(obj) {
		key, _ := clientgocache.MetaNamespaceKeyFunc(obj)
		s.queue.Add(key)
	}
}

func (s *ResourceSyncer) Run(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	defer utilruntime.HandleCrash()
	defer s.queue.ShutDown()

	s.logger.Info("Starting resource syncer")

	// Wait for the caches to be synced before starting workers
	s.logger.Info("Waiting for informer caches to sync")

	if ok := clientgocache.WaitForCacheSync(s.ctx.Done(), s.source.HasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	workers := s.source.SyncRule().MaxConcurrent
	if workers <= 0 {
		workers = defaultWorkers
	}
	s.logger.Info("Starting workers", "count", workers)
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(s.ctx, s.runWorker, time.Second)
	}

	s.logger.Info("Started workers")
	<-s.ctx.Done()
	s.logger.Info("Shutting down workers")

	return nil
}

func (s *ResourceSyncer) runWorker(ctx context.Context) {
	for s.processNextWorkItem(ctx) {
	}
}

func (s *ResourceSyncer) processNextWorkItem(ctx context.Context) bool {
	obj, shutdown := s.queue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer s.queue.Done(obj)

		key, ok := obj.(string)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := s.sync(ctx, key); err != nil {
			s.queue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s/%s': %s, requeuing", s.source.SyncRule().Resource, key, err.Error())
		}

		s.queue.Forget(obj)
		s.logger.Info("Successfully synced", "event", key)
		return nil
	}(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (s *ResourceSyncer) sync(ctx context.Context, key string) error {
	obj, ok := s.pendingResources.Get(key)
	if !ok {
		s.logger.WithValues("resourceKey", key).Info("ignore sync, resource object already processed")
		return nil
	}

	_, isDeleted := obj.(deleted)
	obj = obj.DeepCopyObject().(client.Object)
	cluster := s.source.Cluster()

	var err error
	if isDeleted {
		err = s.storage.Delete(ctx, cluster, obj)
	} else {
		err = s.storage.Save(ctx, cluster, obj)
	}
	if err != nil {
		return err
	}

	// obj was successfully processed, remove it from cache
	s.pendingResources.Remove(obj)
	return nil
}
