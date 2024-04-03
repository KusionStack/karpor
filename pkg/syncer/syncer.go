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
	"strings"
	"time"

	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karbour/pkg/kubernetes/apis/search/v1beta1"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const defaultWorkers = 10

type deleted struct {
	client.Object
}

type ResourceSyncer struct {
	source  SyncSource
	storage storage.Storage

	queue  workqueue.RateLimitingInterface
	ctx    context.Context
	cancel context.CancelFunc

	logger logr.Logger
}

func NewResourceSyncer(cluster string, dynamicClient dynamic.Interface, rsr v1beta1.ResourceSyncRule, storage storage.Storage) *ResourceSyncer {
	source := NewSource(cluster, dynamicClient, rsr, storage)
	return &ResourceSyncer{
		source:  source,
		storage: storage,
		queue:   workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), fmt.Sprintf("%s/%s-sync-queue", rsr.APIVersion, rsr.Resource)),
		logger:  ctrl.Log.WithName(fmt.Sprintf("%s-syncer", source.SyncRule().Resource)),
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
		return errors.Wrap(err, "failed to stop the source")
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
	key, _ := clientgocache.MetaNamespaceKeyFunc(obj)
	s.queue.Add(key)
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
		//nolint:contextcheck
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
	item, shutdown := s.queue.Get()
	if shutdown {
		return false
	}
	key := item.(string)
	func() {
		defer s.queue.Done(item)

		if err := s.sync(ctx, key); err != nil {
			s.queue.AddRateLimited(item)
			return
		}
		s.queue.Forget(item)
	}()
	return true
}

func (s *ResourceSyncer) sync(ctx context.Context, key string) error {
	op, err := s.syncStorage(ctx, key)
	if err != nil {
		s.logger.Error(err, "failed to sync", "key", key, "operation", op)
		return err
	}

	s.logger.V(1).Info("successfully sync", "key", key, "operation", op)
	return nil
}

func (s *ResourceSyncer) syncStorage(ctx context.Context, key string) (string, error) {
	val, exists, err := s.source.GetByKey(key)
	if err != nil {
		return "none", err
	}

	if exists {
		err = s.storage.Save(ctx, s.source.Cluster(), val.(*unstructured.Unstructured))
		return "save", err
	}

	obj, err := genUnObj(s.SyncRule(), key)
	if err != nil {
		return "none", err
	}

	op := "delete"
	err = s.storage.Delete(ctx, s.source.Cluster(), obj)
	if err != nil {
		if errors.Is(err, elasticsearch.ErrNotFound) {
			return op, nil
		}
		return op, err
	}
	return op, nil
}

func genUnObj(sr v1beta1.ResourceSyncRule, key string) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion(sr.APIVersion)
	if len(sr.Resource) == 0 {
		return nil, fmt.Errorf("resource name is not provided in ResourceSyncRule")
	}
	obj.SetKind(sr.Resource[0 : len(sr.Resource)-1])
	if len(key) == 0 {
		return nil, fmt.Errorf("key is not provided")
	}
	keys := strings.Split(key, "/")
	if len(keys) == 1 {
		obj.SetName(keys[0])
	} else if len(keys) == 2 {
		obj.SetNamespace(keys[0])
		obj.SetName(keys[1])
	}
	return obj, nil
}
