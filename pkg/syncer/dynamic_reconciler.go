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
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	"github.com/KusionStack/karpor/pkg/syncer/transform"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
)

type DynamicReconciler struct {
	ctx         context.Context
	gvk         schema.GroupVersionKind
	clusterName string
	logger      logr.Logger
	storage     storage.ResourceStorage

	syncRule      v1beta1.ResourceSyncRule
	transformFunc syncercache.TransformFunc
	trimFunc      syncercache.TransformFunc
	client        client.Client

	scheme *runtime.Scheme
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
	r.scheme = mgr.GetScheme()

	// step 2: set manager
	c, err := controller.New(r.gvk.String(), mgr, controller.Options{
		Reconciler:              r,
		MaxConcurrentReconciles: 10,
		RateLimiter:             utils.NewRateLimiter(12),
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
	logger := r.logger
	logger.WithValues("kind", r.gvk).WithValues("name", req.NamespacedName)

	defer func() {
		if retErr != nil {
			logger.Error(retErr, "reconciliation error")
		}
		if err := recover(); err != nil {
			logger.Error(errors.New("expected no panic; recovered"), "")
		}
	}()

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

			logger.Info("get rule now, start to sync.")
		} else {
			logger.Info("wait rule ready, restart in 20 seconds...")
			return reconcile.Result{RequeueAfter: 20 * time.Second}, nil
		}
	}

	/*
	    If using protobuf type to receive unstructured Object in the 1.12 Kubernetes, it will decode failed.
	   	To be compatible with old version kubernetes, it uses object.Object rather than unstructured Object.
	*/
	runtimeObj, err := r.scheme.New(r.gvk)
	if err != nil {
		runtimeObj = &unstructured.Unstructured{}
	}
	obj, err := runtimeObjectToObject(runtimeObj)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.client.Get(ctx, req.NamespacedName, obj)
	if err != nil {
		if apierrors.IsNotFound(err) {
			obj := genUnObj(r.syncRule, req.String())
			err = r.storage.DeleteResource(ctx, r.clusterName, obj)
			if errors.Is(err, elasticsearch.ErrNotFound) {
				logger.Error(err, "failed to delete resource")
				return reconcile.Result{}, nil
			}
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, err
	}

	unstructuredObj, err := objectToUnstructured(obj)
	if err != nil {
		return reconcile.Result{}, err
	}
	// trimer
	if r.trimFunc != nil {
		val, err := r.trimFunc(unstructuredObj)
		if err != nil {
			return reconcile.Result{}, err
		}
		unstructuredObj = val.(*unstructured.Unstructured)
	}
	// transformer
	if r.transformFunc != nil {
		val, err := r.transformFunc(unstructuredObj)
		if err != nil {
			return reconcile.Result{}, err
		}
		unstructuredObj = val.(*unstructured.Unstructured)
	}
	err = r.storage.SaveResource(ctx, r.clusterName, unstructuredObj)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to save resource")
	}

	return reconcile.Result{}, nil
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
		return nil, err
	}
	return obj, nil
}

func objectToUnstructured(obj interface{}) (*unstructured.Unstructured, error) {
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
	return unstructuredObj, nil
}

func runtimeObjectToObject(runtimeObj interface{}) (client.Object, error) {
	obj, ok := runtimeObj.(client.Object)
	if !ok {
		return nil, errors.New("error in convert runtimeObj")
	}
	return obj, nil
}
