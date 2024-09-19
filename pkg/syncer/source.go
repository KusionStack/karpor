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
	"time"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karpor/pkg/syncer/internal"
	"github.com/KusionStack/karpor/pkg/syncer/jsonextracter"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
	"github.com/KusionStack/karpor/pkg/util/jsonpath"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlhandler "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	defaultResyncPeriod = 1 * time.Hour
)

// SyncSource defines the interface for sources that can be synced, including methods for interacting with the Kubernetes API and cache.
type SyncSource interface {
	source.Source
	clientgocache.Store
	Cluster() string
	SyncRule() v1beta1.ResourceSyncRule
	Stop(context.Context) error
	HasSynced() bool
}

// informerSource is a struct that implements the SyncSource interface, providing functionality for syncing resources using informers.
type informerSource struct {
	cluster string
	v1beta1.ResourceSyncRule
	storage storage.ResourceStorage

	client   dynamic.Interface
	cache    clientgocache.Store
	informer clientgocache.Controller

	ctx     context.Context
	cancel  context.CancelFunc
	stopped chan struct{}

	logger logr.Logger
}

func (s *informerSource) Add(obj interface{}) error {
	return s.cache.Add(obj)
}

func (s *informerSource) Update(obj interface{}) error {
	return s.cache.Update(obj)
}

func (s *informerSource) Delete(obj interface{}) error {
	return s.cache.Delete(obj)
}

func (s *informerSource) List() []interface{} {
	return s.cache.List()
}

func (s *informerSource) ListKeys() []string {
	return s.cache.ListKeys()
}

func (s *informerSource) Get(obj interface{}) (item interface{}, exists bool, err error) {
	return s.cache.Get(obj)
}

func (s *informerSource) GetByKey(key string) (item interface{}, exists bool, err error) {
	return s.cache.GetByKey(key)
}

func (s *informerSource) Replace(i []interface{}, s2 string) error {
	return s.cache.Replace(i, s2)
}

func (s *informerSource) Resync() error {
	return s.cache.Resync()
}

// NewSource creates a new instance of informerSource with the provided parameters, including cluster name, Kubernetes client, sync rule, and storage.
func NewSource(cluster string, client dynamic.Interface, rsr v1beta1.ResourceSyncRule, storage storage.ResourceStorage) SyncSource {
	return &informerSource{
		cluster:          cluster,
		storage:          storage,
		ResourceSyncRule: rsr,
		client:           client,
		stopped:          make(chan struct{}),
		logger:           ctrl.Log.WithName(fmt.Sprintf("%s-syncer-source", rsr.Resource)),
	}
}

func (s *informerSource) Cluster() string {
	return s.cluster
}

func (s *informerSource) SyncRule() v1beta1.ResourceSyncRule {
	return s.ResourceSyncRule
}

// Start initializes and starts the informerSource, setting up informers and handlers for resource syncing based on the provided context, event handler, workqueue, and predicates.
func (s *informerSource) Start(ctx context.Context, handler ctrlhandler.EventHandler, queue workqueue.RateLimitingInterface, predicates ...predicate.Predicate) error {
	cache, informer, err := s.createInformer(ctx, handler, queue, predicates...)
	if err != nil {
		return err
	}
	s.cache = cache
	s.informer = informer

	s.ctx, s.cancel = context.WithCancel(ctx)
	go func() {
		s.informer.Run(s.ctx.Done())
		close(s.stopped)
	}()

	return nil
}

// Stop gracefully shuts down the informerSource, stopping informers and canceling the context.
func (s *informerSource) Stop(ctx context.Context) error {
	s.cancel()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil
		}
		return errors.New("timed out waiting for source to stop")
	case <-s.stopped:
		return nil
	}
}

// createInformer sets up and returns the informer and controller for the informerSource, using the provided context, event handler, workqueue, and predicates.
func (s *informerSource) createInformer(ctx context.Context, handler ctrlhandler.EventHandler, queue workqueue.RateLimitingInterface, predicates ...predicate.Predicate) (clientgocache.Store, clientgocache.Controller, error) {
	gvr, err := parseGVR(&s.ResourceSyncRule)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error parsing GroupVersionResource")
	}

	selectors, err := parseSelectors(s.ResourceSyncRule)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing selectors: %v", selectors)
	}

	trim, err := s.parseTrimer()
	if err != nil {
		return nil, nil, errors.Wrap(err, "error parsing trim rule")
	}

	resyncPeriod := defaultResyncPeriod
	if s.ResyncPeriod != nil {
		resyncPeriod = s.ResyncPeriod.Duration
	}

	lw := &clientgocache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return s.client.Resource(gvr).Namespace(s.Namespace).List(context.TODO(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return s.client.Resource(gvr).Namespace(s.Namespace).Watch(context.TODO(), options)
		},
	}

	h := &internal.EventHandler{EventHandler: handler, Queue: queue, Predicates: predicates}
	cache, informer := clientgocache.NewTransformingInformer(lw, &unstructured.Unstructured{}, resyncPeriod, h, trim)
	// TODO: Use interface instead of struct
	importer := utils.NewESImporter(s.storage.(*elasticsearch.Storage), s.cluster, gvr)
	if err = importer.ImportTo(ctx, cache); err != nil {
		return nil, nil, err
	}
	return cache, informer, nil
}

func (s *informerSource) HasSynced() bool {
	return s.informer.HasSynced()
}

// parseGVR extracts and returns the GroupVersionResource information from the provided ResourceSyncRule.
func parseGVR(rsr *v1beta1.ResourceSyncRule) (schema.GroupVersionResource, error) {
	gv, err := schema.ParseGroupVersion(rsr.APIVersion)
	if err != nil {
		return schema.GroupVersionResource{}, fmt.Errorf("invalid group version %q", rsr.APIVersion)
	}
	return gv.WithResource(rsr.Resource), nil
}

// parseSelectors extracts and returns the list of Selectors from the provided ResourceSyncRule.
func parseSelectors(rsr v1beta1.ResourceSyncRule) ([]utils.Selector, error) {
	if len(rsr.Selectors) == 0 {
		return nil, nil
	}

	selectors := make([]utils.Selector, 0, len(rsr.Selectors))
	for _, s := range rsr.Selectors {
		var selector utils.Selector
		if s.LabelSelector != nil {
			labelSelector, err := metav1.LabelSelectorAsSelector(s.LabelSelector)
			if err != nil {
				return nil, err
			}
			selector.Label = labelSelector
		}
		if s.FieldSelector != nil {
			selector.Field = utils.FieldsSelector{
				Selector:        fields.SelectorFromSet(fields.Set(s.FieldSelector.MatchFields)),
				ServerSupported: s.FieldSelector.ServerSupported,
			}
		}
		selectors = append(selectors, selector)
	}
	return selectors, nil
}

func (s *informerSource) parseTrimer() (clientgocache.TransformFunc, error) {
	t := s.ResourceSyncRule.Trim
	if t == nil || len(t.Retain.JSONPaths) == 0 {
		return nil, nil
	}

	extracters := make([]jsonextracter.Extracter, 0, len(t.Retain.JSONPaths))
	for _, p := range t.Retain.JSONPaths {
		p, err := jsonpath.RelaxedJSONPathExpression(p)
		if err != nil {
			return nil, err
		}

		ex, err := jsonextracter.BuildExtracter(p, true)
		if err != nil {
			return nil, err
		}
		extracters = append(extracters, ex)
	}

	trimFunc := func(obj interface{}) (ret interface{}, err error) {
		defer func() {
			if err != nil {
				s.logger.Error(err, "error in triming object")
				ret, err = obj, nil
			}
		}()

		if d, ok := obj.(clientgocache.DeletedFinalStateUnknown); ok {
			// Since we import ES data into informer cache at startup, the
			// resource that was deleted during the restart will generate
			// DeletedFinalStateUnknown.
			// We unwarp the object here, so there is no need for following
			// steps including event handler to care about DeletedFinalStateUnknown.
			obj = d.Obj
		}

		u, ok := obj.(*unstructured.Unstructured)
		if !ok {
			return nil, fmt.Errorf("trim: object's type should be *unstructured.Unstructured, but received %T", obj)
		}

		merged, err := jsonextracter.Merge(extracters, u.Object)
		if err != nil {
			return nil, err
		}

		unObj := &unstructured.Unstructured{Object: merged}
		return unObj, nil
	}

	return trimFunc, nil
}
