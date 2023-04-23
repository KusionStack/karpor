// Copyright 2017 The Karbour Authors.
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

package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KusionStack/karbour/pkg/search/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
)

type ResourceHandler interface {
	OnAdd(obj interface{}) error
	OnUpdate(newObj interface{}) error
	OnDelete(obj interface{}) error
}

type ResourceHandlerFuncs struct {
	AddFunc    func(obj interface{}) error
	UpdateFunc func(newObj interface{}) error
	DeleteFunc func(obj interface{}) error
}

func (r ResourceHandlerFuncs) OnAdd(obj interface{}) error {
	return r.AddFunc(obj)
}

func (r ResourceHandlerFuncs) OnUpdate(newObj interface{}) error {
	return r.UpdateFunc(newObj)
}

func (r ResourceHandlerFuncs) OnDelete(obj interface{}) error {
	return r.DeleteFunc(obj)
}

type ResourceSyncConfig struct {
	ClusterName    string
	DynamicClient  dynamic.Interface
	GVR            schema.GroupVersionResource
	Namespace      string
	Selectors      []utils.Selector
	JSONPathParser *utils.JSONPathParser
	TransformFunc  cache.TransformFunc
	ResyncPeriod   time.Duration
	Handler        ResourceHandler
}

func NewResourceInformer(config ResourceSyncConfig) cache.Controller {
	informerCache := NewResourceCache()
	fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KnownObjects:          informerCache,
		EmitDeltaTypeReplaced: true,
	})

	selectors := utils.MultiSelectors(config.Selectors)
	var listWatchSelector *utils.Selector
	if len(selectors) == 1 {
		selector := selectors[0]
		if selector.ServerSupported() {
			listWatchSelector = &selector
			selectors = nil
		}
	}

	applyToList := func(options *metav1.ListOptions) {
		if listWatchSelector != nil {
			if listWatchSelector.Label != nil {
				options.LabelSelector = listWatchSelector.Label.String()
			}
			if listWatchSelector.Field != nil {
				options.FieldSelector = listWatchSelector.Field.String()
			}
		}
	}

	cfg := &cache.Config{
		Queue: fifo,
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				applyToList(&options)
				return config.DynamicClient.Resource(config.GVR).Namespace(config.Namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				applyToList(&options)
				return config.DynamicClient.Resource(config.GVR).Namespace(config.Namespace).Watch(context.TODO(), options)
			},
		},
		ObjectType:       &unstructured.Unstructured{},
		FullResyncPeriod: config.ResyncPeriod,
		RetryOnError:     true,
		Process: func(d interface{}) error {
			deltas, ok := d.(cache.Deltas)
			if !ok {
				return errors.New("object given as Process argument is not Deltas")
			}

			// only process the latest delta
			newest := deltas.Newest()
			obj := newest.Object

			// filter
			if selectors != nil {
				u, ok := obj.(*unstructured.Unstructured)
				if !ok {
					return fmt.Errorf("unexpected type '%T', should be *unstructured.Unstructured", obj)
				}
				matched, err := selectors.Matches(utils.SelectableUnstructured(u, config.JSONPathParser))
				if err != nil {
					return err
				}
				if !matched {
					return nil
				}
			}

			// transform
			if config.TransformFunc != nil {
				transformed, err := config.TransformFunc(obj)
				if err != nil {
					return fmt.Errorf("error transforming object: %v", err)
				}
				obj = transformed
			}

			h := config.Handler
			switch newest.Type {
			case cache.Sync, cache.Replaced, cache.Added, cache.Updated:
				if _, exists, err := informerCache.Get(obj); err == nil && exists {
					if newer, err := informerCache.IsNewer(obj); err != nil {
						return err
					} else if !newer {
						return nil
					}

					if err := informerCache.Update(obj); err != nil {
						return err
					}
					return h.OnUpdate(obj)
				} else {
					if err := informerCache.Add(obj); err != nil {
						return err
					}
					return h.OnAdd(obj)
				}
			case cache.Deleted:
				if err := informerCache.Delete(obj); err != nil {
					return err
				}
				return h.OnDelete(obj)
			}
			return nil
		},
	}

	return cache.New(cfg)
}
