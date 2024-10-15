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
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karpor/pkg/syncer/transform"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
	sprig "github.com/Masterminds/sprig/v3"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	defaultWorkers = 10

	// purgeMarker indicates it's time to prune the storage.
	// As k8s object name cannot use underscore, we can use this name in workqueue without collision.
	purgeMarker = "__purge_marker__"
)

// deleted is a type that represents a deleted Kubernetes object.
type deleted struct {
	client.Object
}

// ResourceSyncer is the main struct that holds the necessary fields and methods for the resource syncer component.
type ResourceSyncer struct {
	source  SyncSource
	storage storage.ResourceStorage

	queue  workqueue.RateLimitingInterface
	ctx    context.Context
	cancel context.CancelFunc

	logger logr.Logger

	transformFunc clientgocache.TransformFunc
	startTime     time.Time
}

// NewResourceSyncer creates a new instance of the ResourceSyncer with the given parameters.
func NewResourceSyncer(cluster string, dynamicClient dynamic.Interface, rsr v1beta1.ResourceSyncRule, storage storage.ResourceStorage) *ResourceSyncer {
	source := NewSource(cluster, dynamicClient, rsr, storage)
	return &ResourceSyncer{
		source:  source,
		storage: storage,
		queue:   workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), fmt.Sprintf("%s/%s-sync-queue", rsr.APIVersion, rsr.Resource)),
		logger:  ctrl.Log.WithName(fmt.Sprintf("%s-syncer", source.SyncRule().Resource)),
	}
}

// Source returns the SyncSource associated with the ResourceSyncer.
func (s *ResourceSyncer) Source() SyncSource {
	return s.source
}

// SyncRule returns the ResourceSyncRule associated with the ResourceSyncer.
func (s *ResourceSyncer) SyncRule() v1beta1.ResourceSyncRule {
	return s.source.SyncRule()
}

// Stop stops the ResourceSyncer and cleans up resources.
func (s *ResourceSyncer) Stop(ctx context.Context) error {
	if err := s.source.Stop(ctx); err != nil {
		return errors.Wrap(err, "failed to stop the source")
	}
	s.cancel()
	return nil
}

// OnAdd handles the addition of a Kubernetes object.
func (s *ResourceSyncer) OnAdd(obj client.Object) {
	s.enqueue(obj)
}

// OnUpdate handles updates to a Kubernetes object.
func (s *ResourceSyncer) OnUpdate(obj client.Object) {
	s.enqueue(obj)
}

// OnDelete handles the deletion of a Kubernetes object.
func (s *ResourceSyncer) OnDelete(obj client.Object) {
	s.enqueue(deleted{Object: obj})
}

// OnGeneric handles generic events for a Kubernetes object.
func (s *ResourceSyncer) OnGeneric(obj client.Object) {
	s.enqueue(obj)
}

// enqueue adds a Kubernetes object to the work queue for processing.
func (s *ResourceSyncer) enqueue(obj client.Object) {
	key, _ := clientgocache.MetaNamespaceKeyFunc(obj)
	s.queue.Add(key)
}

// Run starts the ResourceSyncer and its workers to process Kubernetes object events.
func (s *ResourceSyncer) Run(ctx context.Context) error {
	s.startTime = time.Now()

	s.ctx, s.cancel = context.WithCancel(ctx)

	defer utilruntime.HandleCrash()
	defer s.queue.ShutDown()

	s.logger.Info("Starting resource syncer")

	// Wait for the caches to be synced before starting workers
	s.logger.Info("Waiting for informer caches to sync")

	if transformFunc, err := s.parseTransformer(); err != nil {
		s.logger.Error(err, "error in parsing transform rule")
	} else {
		s.transformFunc = transformFunc
	}

	if ok := clientgocache.WaitForCacheSync(s.ctx.Done(), s.source.HasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	// We push the purgeMarker after cacheSync, meaning that when the purgeMarker
	// is read from the queue, almost all resources have been synced.
	s.queue.Add(purgeMarker)

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

// runWorker is the main worker loop for the ResourceSyncer, processing items from the work queue.
func (s *ResourceSyncer) runWorker(ctx context.Context) {
	for s.processNextWorkItem(ctx) {
	}
}

func (s *ResourceSyncer) purgeStorage(ctx context.Context) {
	err := s.storage.Refresh(ctx)
	if err != nil {
		s.logger.Error(err, "error in refreshing storage")
		return
	}

	r := s.SyncRule()
	gvr, err := parseGVR(&r)
	if err != nil {
		s.logger.Error(err, "error in parsing GVR")
		return
	}

	// TODO: Use interface instead of struct
	esPurger := utils.NewESPurger(s.storage.(*elasticsearch.Storage), s.source.Cluster(), gvr, s.source, s.OnDelete)
	if err := esPurger.Purge(ctx, s.startTime); err != nil {
		s.logger.Error(err, "error in purging ES")
		return
	}
}

// processNextWorkItem processes the next work item from the queue, returning true if work continues.
func (s *ResourceSyncer) processNextWorkItem(ctx context.Context) bool {
	item, shutdown := s.queue.Get()
	if shutdown {
		return false
	}
	key := item.(string)

	if key == purgeMarker {
		s.purgeStorage(ctx)
		return true
	}

	func() {
		defer s.queue.Done(item)

		if err := s.sync(ctx, key); err != nil {
			// Retry 12 times, about 20 seconds.
			if s.queue.NumRequeues(item) < 12 {
				s.queue.AddRateLimited(item)
				return
			} else {
				s.logger.Error(err, "retry reached max times", "key", key)
			}
		}
		s.queue.Forget(item)
	}()
	return true
}

func (s *ResourceSyncer) saveResource(ctx context.Context, obj runtime.Object) error {
	return s.storage.SaveResource(ctx, s.source.Cluster(), obj)
}

func (s *ResourceSyncer) deleteResource(ctx context.Context, obj runtime.Object) error {
	remainAfterDeleted := s.source.SyncRule().RemainAfterDeleted
	if remainAfterDeleted {
		return s.storage.SoftDeleteResource(ctx, s.source.Cluster(), obj)
	}

	return s.storage.DeleteResource(ctx, s.source.Cluster(), obj)
}

// sync synchronizes the specified resource based on the key provided.
func (s *ResourceSyncer) sync(ctx context.Context, key string) error {
	val, exists, err := s.source.GetByKey(key)
	if err != nil {
		return err
	}

	if exists && s.transformFunc != nil {
		val, err = s.transformFunc(val)
		if err != nil {
			return err
		}
	}

	var op string
	if exists {
		op = "save"
		obj := val.(*unstructured.Unstructured)
		err = s.saveResource(ctx, obj)
	} else {
		op = "delete"
		obj := genUnObj(s.SyncRule(), key)
		err = s.deleteResource(ctx, obj)
		if errors.Is(err, elasticsearch.ErrNotFound) {
			s.logger.Error(err, "failed to sync", "key", key, "op", op)
			err = nil
		}
	}

	if err != nil {
		s.logger.Error(err, "failed to sync", "key", key, "op", op)
		return err
	}

	s.logger.V(1).Info("successfully sync", "key", key, "op", op)
	return nil
}

// parseTransformer creates and returns a transformation function for the informerSource based on the ResourceSyncRule's transformers.
func (s *ResourceSyncer) parseTransformer() (clientgocache.TransformFunc, error) {
	t := s.source.SyncRule().Transform
	if t == nil {
		return nil, nil
	}

	fn, found := transform.GetTransformFunc(t.Type)
	if !found {
		return nil, fmt.Errorf("unsupported transform type %q", t.Type)
	}

	tmpl, err := newTemplate(t.ValueTemplate, s.source.Cluster())
	if err != nil {
		return nil, errors.Wrap(err, "invalid transform template")
	}

	return func(obj interface{}) (ret interface{}, err error) {
		defer func() {
			if err != nil {
				s.logger.Error(err, "error in transforming object")
			}
		}()

		u, ok := obj.(*unstructured.Unstructured)
		if !ok {
			return nil, fmt.Errorf("transform: object's type should be *unstructured.Unstructured, but received %T", obj)
		}

		templateData := struct {
			*unstructured.Unstructured
			Cluster string
		}{
			Unstructured: u,
			Cluster:      s.source.Cluster(),
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, templateData); err != nil {
			return nil, errors.Wrap(err, "transform: error rendering template")
		}
		return fn(obj, buf.String())
	}, nil
}

// genUnObj creates a new unstructured.Unstructured object based on the ResourceSyncRule and key.
func genUnObj(sr v1beta1.ResourceSyncRule, key string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion(sr.APIVersion)
	obj.SetKind(sr.Resource[0 : len(sr.Resource)-1])
	keys := strings.Split(key, "/")
	if len(keys) == 1 {
		obj.SetName(keys[0])
	} else if len(keys) == 2 {
		obj.SetNamespace(keys[0])
		obj.SetName(keys[1])
	}
	return obj
}

// newTemplate creates and returns a new text template from the provided string, which can be used for processing templates in the syncer.
func newTemplate(tmpl, cluster string) (*template.Template, error) {
	clusterFuncs, _ := transform.GetClusterTmplFuncs(cluster)
	return template.New("transformTemplate").Funcs(sprig.FuncMap()).Funcs(clusterFuncs).Parse(tmpl)
}
