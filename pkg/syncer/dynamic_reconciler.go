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
	"reflect"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgocache "k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	syncercache "github.com/KusionStack/karpor/pkg/syncer/cache"
	"github.com/KusionStack/karpor/pkg/syncer/jsonextracter"
	"github.com/KusionStack/karpor/pkg/syncer/transform"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
	"github.com/KusionStack/karpor/pkg/util/jsonpath"
)

type DynamicReconciler struct {
	ctx         context.Context
	gvk         schema.GroupVersionKind
	clusterName string
	logger      logr.Logger
	storage     storage.ResourceStorage

	lock          sync.Mutex
	syncRule      v1beta1.ResourceSyncRule
	transformFunc syncercache.TransformFunc
	trimFunc      syncercache.TransformFunc
	client        client.Client
}

func NewDynamicReconciler(ctx context.Context, clusterName string, gvk schema.GroupVersionKind, storage storage.ResourceStorage) *DynamicReconciler {
	return &DynamicReconciler{
		ctx:         ctx,
		clusterName: clusterName,
		gvk:         gvk,
		logger:      ctrl.Log.WithName(gvk.String()),
		storage:     storage,
		syncRule:    utils.ZeroVal,
	}
}

func (r *DynamicReconciler) SetupWithManager(mgr manager.Manager) error {
	// step 1: set client for DynamicReconciler
	r.client = mgr.GetClient()

	// step 2: set manager
	c, err := controller.New(r.gvk.String(), mgr, controller.Options{
		Reconciler:              r,
		MaxConcurrentReconciles: 10,
	})
	if err != nil {
		return err
	}
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(r.gvk)
	err = c.Watch(&source.Kind{Type: obj}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// step 3: register transform func for template pkg
	if err := transform.RegisterClusterTmplFunc(r.clusterName, "objectRef", r.getObject); err != nil {
		r.logger.Error(err, "error in registering tmpl func")
	}

	return nil
}

func (r *DynamicReconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, retErr error) {
	defer func() {
		if err := recover(); err != nil {
			r.logger.Error(errors.New("expected no panic; recovered"), "")
		}
	}()

	r.logger.WithValues("kind", r.gvk)

	if reflect.DeepEqual(r.syncRule, utils.ZeroVal) {
		rule, exist := utils.GetSyncGVK(r.gvk)
		if exist && !reflect.DeepEqual(rule, utils.ZeroVal) {
			r.syncRule = rule

			if rule.Transform != nil {
				transformFunc, err := parseTransformer(ctx, rule.Transform, r.clusterName)
				if err != nil {
					return reconcile.Result{}, err
				}
				r.transformFunc = transformFunc
			}

			if rule.Trim != nil {
				trimFunc, err := parseTrimer(ctx, rule.Trim)
				if err != nil {
					return reconcile.Result{}, err
				}
				r.trimFunc = trimFunc
			}

			r.logger.Info("get rule now, start to sync.")
		} else {
			r.logger.Info("wait rule ready, restart in 8 seconds...")
			return reconcile.Result{RequeueAfter: 8 * time.Second}, nil
		}
	}

	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(r.gvk)
	err := r.client.Get(ctx, req.NamespacedName, obj)
	if err != nil {
		if apierrors.IsNotFound(err) {
			obj := genUnObj(r.syncRule, req.String())
			err = r.storage.DeleteResource(ctx, r.clusterName, obj)
			if errors.Is(err, elasticsearch.ErrNotFound) {
				r.logger.Error(err, "failed to delete", "key", req.String())
				return reconcile.Result{}, nil
			}
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, err
	}

	// trimer
	if r.trimFunc != nil {
		val, err := r.transformFunc(obj)
		if err != nil {
			return reconcile.Result{}, err
		}
		obj = val.(*unstructured.Unstructured)
	}
	// transformer
	if r.transformFunc != nil {
		val, err := r.transformFunc(obj)
		if err != nil {
			return reconcile.Result{}, err
		}
		obj = val.(*unstructured.Unstructured)
	}
	err = r.storage.SaveResource(ctx, r.clusterName, obj)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to save resource")
	}

	return reconcile.Result{}, nil
}

// parseTransformer creates and returns a transformation function for the informerSource based on the ResourceSyncRule's transformers.
func parseTransformer(ctx context.Context, t *v1beta1.TransformRuleSpec, clusterName string) (syncercache.TransformFunc, error) {
	fn, found := transform.GetTransformFunc(t.Type)
	if !found {
		return nil, fmt.Errorf("unsupported transform type %q", t.Type)
	}

	tmpl, err := newTemplate(t.ValueTemplate, clusterName)
	if err != nil {
		return nil, errors.Wrap(err, "invalid transform template")
	}

	return func(obj interface{}) (ret interface{}, err error) {
		defer func() {
			if err != nil {
				logr.FromContext(ctx).Error(err, "error in transforming object")
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
			Cluster:      clusterName,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, templateData); err != nil {
			return nil, errors.Wrap(err, "transform: error rendering template")
		}
		return fn(obj, buf.String())
	}, nil
}

func parseTrimer(ctx context.Context, t *v1beta1.TrimRuleSpec) (syncercache.TransformFunc, error) {
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
				logr.FromContext(ctx).Error(err, "error in triming object")
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

func (r *DynamicReconciler) getObject(apiVersion, kind, namespace, name string) (interface{}, error) {
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return nil, err
	}

	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gv.Group,
		Version: gv.Version,
		Kind:    kind,
	})
	err = r.client.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, obj)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, elasticsearch.ErrNotFound
		}
	}
	return obj, nil
}
