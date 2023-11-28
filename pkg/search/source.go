package search

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/KusionStack/karbour/pkg/apis/search/v1beta1"
	"github.com/KusionStack/karbour/pkg/search/cache"
	"github.com/KusionStack/karbour/pkg/search/internal"
	"github.com/KusionStack/karbour/pkg/search/transform"
	"github.com/KusionStack/karbour/pkg/search/utils"
	sprig "github.com/Masterminds/sprig/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	ctrlhandler "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type SyncSource interface {
	source.Source
	Cluster() string
	SyncRule() v1beta1.ResourceSyncRule
	Stop(context.Context) error
	HasSynced() bool
}

type informerSource struct {
	cluster string
	v1beta1.ResourceSyncRule

	client   dynamic.Interface
	informer clientgocache.Controller

	ctx     context.Context
	cancel  context.CancelFunc
	stopped chan struct{}
}

func NewSource(cluster string, client dynamic.Interface, rsr v1beta1.ResourceSyncRule) SyncSource {
	return &informerSource{
		cluster:          cluster,
		ResourceSyncRule: rsr,
		client:           client,
		stopped:          make(chan struct{}),
	}
}

func (s *informerSource) Cluster() string {
	return s.cluster
}

func (s *informerSource) SyncRule() v1beta1.ResourceSyncRule {
	return s.ResourceSyncRule
}

func (s *informerSource) Start(ctx context.Context, handler ctrlhandler.EventHandler, queue workqueue.RateLimitingInterface, predicates ...predicate.Predicate) error {
	informer, err := s.createInformer(handler, queue, predicates...)
	if err != nil {
		return err
	}
	s.informer = informer

	s.ctx, s.cancel = context.WithCancel(ctx)
	go func() {
		s.informer.Run(s.ctx.Done())
		close(s.stopped)
	}()

	return nil
}

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

func (s *informerSource) createInformer(handler ctrlhandler.EventHandler, queue workqueue.RateLimitingInterface, predicates ...predicate.Predicate) (clientgocache.Controller, error) {
	gvr, err := parseGVR(&s.ResourceSyncRule)
	if err != nil {
		return nil, fmt.Errorf("error parsing GroupVersionResource: %v", err)
	}

	selectors, err := parseSelectors(s.ResourceSyncRule)
	if err != nil {
		return nil, fmt.Errorf("error parsing selectors: %v", selectors)
	}

	transform, err := s.parseTransformer()
	if err != nil {
		return nil, fmt.Errorf("error parsing transform rule: %v", err)
	}

	var resyncPeriod time.Duration
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

	h := internal.EventHandler{EventHandler: handler, Queue: queue, Predicates: predicates}
	return cache.NewResourceInformer(lw, utils.MultiSelectors(selectors), transform, resyncPeriod, h), nil
}

func (s *informerSource) HasSynced() bool {
	return s.informer.HasSynced()
}

func parseGVR(rsr *v1beta1.ResourceSyncRule) (schema.GroupVersionResource, error) {
	gv, err := schema.ParseGroupVersion(rsr.APIVersion)
	if err != nil {
		return schema.GroupVersionResource{}, fmt.Errorf("invalid group version %q", rsr.APIVersion)
	}
	return gv.WithResource(rsr.Resource), nil
}

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
				ServerSupported: s.FieldSelector.SeverSupported,
			}
		}
		selectors = append(selectors, selector)
	}
	return selectors, nil
}

func (s *informerSource) parseTransformer() (clientgocache.TransformFunc, error) {
	t := s.ResourceSyncRule.Transform
	if t == nil {
		return nil, nil
	}

	fn, found := transform.GetTransformFunc(t.Type)
	if !found {
		return nil, fmt.Errorf("unsupported transform type %q", t.Type)
	}

	tmpl, err := newTemplate(t.ValueTemplate)
	if err != nil {
		return nil, fmt.Errorf("invalid transform template: %v", err)
	}

	return func(obj interface{}) (interface{}, error) {
		u, ok := obj.(*unstructured.Unstructured)
		if !ok {
			return nil, fmt.Errorf("transform: object's type should be *unstructured.Unstructured")
		}

		templateData := struct {
			*unstructured.Unstructured
			Cluster string
		}{
			Unstructured: u,
			Cluster:      s.cluster,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, templateData); err != nil {
			return nil, fmt.Errorf("transform: error rendering template: %v", err)
		}
		return fn(obj, buf.String())
	}, nil
}

func newTemplate(tmpl string) (*template.Template, error) {
	return template.New("transformTemplate").Funcs(sprig.FuncMap()).Parse(tmpl)
}
