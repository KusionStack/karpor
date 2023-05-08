package search

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	searchv1beta1 "github.com/KusionStack/karbour/pkg/apis/search/v1beta1"
	"github.com/KusionStack/karbour/pkg/search/cache"
	"github.com/KusionStack/karbour/pkg/search/transform"
	"github.com/KusionStack/karbour/pkg/search/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type SingleClusterSyncManager interface {
	Start() error
	Started() bool
	Stop()
	Stopped() bool
	SetSyncResources(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	IsSyncerExist(schema.GroupVersionResource) bool
	ClusterConfig() *rest.Config
}

type singleClusterSyncManager struct {
	clusterName   string
	clusterConfig *rest.Config
	dynamicClient dynamic.Interface

	ctx                 context.Context
	cancel              context.CancelFunc
	startLock           sync.RWMutex
	startOnce, stopOnce sync.Once
	started, stopped    bool
	wg                  wait.Group

	syncResources   atomic.Value
	ch              chan struct{}
	resourceSyncers sync.Map
	logger          klog.Logger
}

func NewSingleClusterSyncManager(baseContext context.Context, clusterName string, config *rest.Config) (SingleClusterSyncManager, error) {
	config = rest.CopyConfig(config)
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	innerCtx, innerCancel := context.WithCancel(baseContext)
	syncer := &singleClusterSyncManager{
		clusterName:   clusterName,
		clusterConfig: config,
		dynamicClient: dynamicClient,
		ctx:           innerCtx,
		cancel:        innerCancel,
		ch:            make(chan struct{}, 1),
		logger:        klog.FromContext(baseContext).WithValues("cluster", clusterName),
	}
	return syncer, nil
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

func (s *singleClusterSyncManager) Start() error {
	s.startOnce.Do(func() {
		s.logger.Info("start sync manager")

		go s.process()

		s.startLock.Lock()
		s.started = true
		s.startLock.Unlock()
	})
	return nil
}

func (s *singleClusterSyncManager) Stop() {
	s.stopOnce.Do(func() {
		s.logger.Info("stop sync manager")
		defer s.logger.Info("sync manager stopped")

		defer close(s.ch)

		s.startLock.Lock()
		s.stopped = true
		s.startLock.Unlock()

		s.cancel()

		s.logger.Info("waiting for syncers to stop")
		s.wg.Wait()
		s.logger.Info("syncers are all stopped")
	})
}

func (s *singleClusterSyncManager) SetSyncResources(newSyncResources map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule) {
	s.logger.Info("update sync resources")
	if s.Stopped() {
		return
	}

	s.syncResources.Store(newSyncResources)
	select {
	case s.ch <- struct{}{}:
	default:
	}
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

		s.handleSyncResources()
	}
}

func (s *singleClusterSyncManager) handleSyncResources() {
	desiredSyncResources := s.syncResources.Load().(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	if len(desiredSyncResources) == 0 {
		return
	}

	for gvr, rsr := range desiredSyncResources {
		if s.Stopped() {
			return
		}

		syncer, ok := s.getSyncer(gvr)
		if ok && !reflect.DeepEqual(syncer.ResourceSyncRule(), rsr) {
			s.logger.Info("ResourceSyncRule was changed", "gvr", gvr)
			s.stopSyncer(gvr, syncer)
			ok = false
		}
		if !ok {
			syncer, err := s.createSyncer(gvr, rsr)
			if err != nil {
				// TODO: expose error
				s.logger.Error(err, "failed to build resource syncer", "grv", gvr)
				continue
			}
			s.wg.StartWithContext(s.ctx, syncer.Run)
		}
	}

	syncersToStop := make(map[schema.GroupVersionResource]*resourceSyncer)
	s.resourceSyncers.Range(func(k, v any) bool {
		gvr := k.(schema.GroupVersionResource)
		if _, ok := desiredSyncResources[gvr]; !ok {
			syncersToStop[gvr] = v.(*resourceSyncer)
		}
		return true
	})
	for gvr, syncer := range syncersToStop {
		s.stopSyncer(gvr, syncer)
	}
}

func (s *singleClusterSyncManager) IsSyncerExist(gvr schema.GroupVersionResource) bool {
	_, found := s.getSyncer(gvr)
	return found
}

func (s *singleClusterSyncManager) getSyncer(gvr schema.GroupVersionResource) (*resourceSyncer, bool) {
	val, ok := s.resourceSyncers.Load(gvr)
	if !ok {
		return nil, false
	}
	return val.(*resourceSyncer), true
}

func (s *singleClusterSyncManager) createSyncer(gvr schema.GroupVersionResource, rsr *searchv1beta1.ResourceSyncRule) (*resourceSyncer, error) {
	s.logger.Info("create resource syncer", "gvr", gvr)
	selectors, err := parseSelectors(rsr)
	if err != nil {
		return nil, fmt.Errorf("parse selectors error: %v", err)
	}

	transform, err := s.parseTransformer(rsr)
	if err != nil {
		return nil, fmt.Errorf("parse transform rule for resource %v error: %v", gvr, err)
	}

	var resyncPeriod time.Duration
	if rsr.ResyncPeriod != nil {
		resyncPeriod = rsr.ResyncPeriod.Duration
	}

	informer := cache.NewResourceInformer(cache.ResourceSyncConfig{
		ClusterName:    s.clusterName,
		DynamicClient:  s.dynamicClient,
		GVR:            gvr,
		Namespace:      rsr.Namespace,
		Selectors:      selectors,
		JSONPathParser: utils.DefaultJSONPathParser,
		ResyncPeriod:   resyncPeriod,
		Handler:        s,
		TransformFunc:  transform,
	})

	syncer := &resourceSyncer{
		rsr:      rsr,
		informer: informer,
		stopped:  make(chan struct{}),
	}

	s.resourceSyncers.Store(gvr, syncer)
	return syncer, nil
}

func (s *singleClusterSyncManager) stopSyncer(gvr schema.GroupVersionResource, syncer *resourceSyncer) {
	s.logger.Info("start to stop the resource syncer", "gvr", gvr)
	<-syncer.Stop()
	s.resourceSyncers.Delete(gvr)
	s.logger.Info("resource syncer was stopped", "resource", gvr)
}

func (s *singleClusterSyncManager) ClusterConfig() *rest.Config {
	return rest.CopyConfig(s.clusterConfig)
}

func (s *singleClusterSyncManager) OnAdd(obj interface{}) error {
	// TODO: implement the method
	return nil
}

func (s *singleClusterSyncManager) OnUpdate(newObj interface{}) error {
	// TODO: implement the method
	return nil
}

func (s *singleClusterSyncManager) OnDelete(obj interface{}) error {
	// TODO: implement the method
	return nil
}

func parseSelectors(rsr *searchv1beta1.ResourceSyncRule) ([]utils.Selector, error) {
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
			selector.Field = &utils.FieldsSelector{
				Selector:        fields.SelectorFromSet(fields.Set(s.FieldSelector.MatchFields)),
				ServerSupported: s.FieldSelector.SeverSupported,
			}
		}
		selectors = append(selectors, selector)
	}
	return selectors, nil
}

func parseGVR(rsr *searchv1beta1.ResourceSyncRule) (schema.GroupVersionResource, error) {
	gv, err := schema.ParseGroupVersion(rsr.APIVersion)
	if err != nil {
		return schema.GroupVersionResource{}, fmt.Errorf("invalid group version %q", rsr.APIVersion)
	}
	return gv.WithResource(rsr.Resource), nil
}

func (s *singleClusterSyncManager) parseTransformer(rsr *searchv1beta1.ResourceSyncRule) (clientgocache.TransformFunc, error) {
	t := rsr.Transform
	if t == nil {
		return nil, nil
	}

	transformer, err := transform.NewTransformer(t.Type, t.ValueTemplate, s.clusterName)
	if err != nil {
		return nil, fmt.Errorf("failed to create transformer: %v", err)
	}
	return transformer.Transform, nil
}

type resourceSyncer struct {
	rsr      *searchv1beta1.ResourceSyncRule
	informer clientgocache.Controller
	ctx      context.Context
	cancel   context.CancelFunc
	stopped  chan struct{}
	runOnce  sync.Once
}

func (r *resourceSyncer) Run(ctx context.Context) {
	r.runOnce.Do(func() {
		r.ctx, r.cancel = context.WithCancel(ctx)
		r.informer.Run(r.ctx.Done())
		close(r.stopped)
	})
}

func (r *resourceSyncer) ResourceSyncRule() *searchv1beta1.ResourceSyncRule {
	return r.rsr
}

func (r *resourceSyncer) Stop() chan struct{} {
	r.cancel()
	return r.stopped
}
