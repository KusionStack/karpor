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

package cache

import (
	"errors"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

// ResourceHandler defines the interface for handling resource events.
type ResourceHandler interface {
	OnAdd(obj interface{}) error
	OnUpdate(newObj interface{}) error
	OnDelete(obj interface{}) error
}

// ResourceHandlerFuncs is a struct that implements the ResourceHandler interface by wrapping functions.
type ResourceHandlerFuncs struct {
	AddFunc    func(obj interface{}) error
	UpdateFunc func(newObj interface{}) error
	DeleteFunc func(obj interface{}) error
}

// OnAdd calls the AddFunc field of ResourceHandlerFuncs, handling the add event.
func (r ResourceHandlerFuncs) OnAdd(obj interface{}) error {
	return r.AddFunc(obj)
}

// OnUpdate calls the UpdateFunc field of ResourceHandlerFuncs, handling the update event.
func (r ResourceHandlerFuncs) OnUpdate(newObj interface{}) error {
	return r.UpdateFunc(newObj)
}

// OnDelete calls the DeleteFunc field of ResourceHandlerFuncs, handling the delete event.
func (r ResourceHandlerFuncs) OnDelete(obj interface{}) error {
	return r.DeleteFunc(obj)
}

// ResourceSelector defines the interface for selecting resources based on certain criteria.
type ResourceSelector interface {
	ApplyToList(*metav1.ListOptions)
	Predicate(interface{}) bool
}

// NewResourceInformer creates a new informer that watches for resource events and handles them using the provided ResourceHandler.
func NewResourceInformer(lw cache.ListerWatcher,
	selector ResourceSelector,
	transform cache.TransformFunc,
	resyncPeriod time.Duration,
	handler ResourceHandler,
	knownObjects cache.KeyListerGetter,
) cache.Controller {
	informerCache := NewResourceCache()
	fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KnownObjects:          knownObjects,
		EmitDeltaTypeReplaced: true,
	})

	doProcess := func(obj interface{}, dType cache.DeltaType) error {
		// transform
		if transform != nil {
			if _, ok := obj.(cache.DeletedFinalStateUnknown); !ok {
				transformed, err := transform(obj)
				if err != nil {
					return fmt.Errorf("error transforming object: %v, delta type: %s", err, dType)
				}
				obj = transformed
			}
		}

		switch dType {
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
				return handler.OnUpdate(obj)
			} else {
				if err := informerCache.Add(obj); err != nil {
					return err
				}
				return handler.OnAdd(obj)
			}
		case cache.Deleted:
			if err := informerCache.Delete(obj); err != nil {
				return err
			}
			return handler.OnDelete(obj)
		}
		return nil
	}

	cfg := &cache.Config{
		Queue:            fifo,
		ObjectType:       &unstructured.Unstructured{},
		FullResyncPeriod: resyncPeriod,
		RetryOnError:     true,
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if selector != nil {
					selector.ApplyToList(&options)
				}
				return lw.List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if selector != nil {
					selector.ApplyToList(&options)
				}
				return lw.Watch(options)
			},
		},

		Process: func(d interface{}) (err error) {
			defer func() {
				if err != nil {
					klog.Errorf("resource informer: error processing item: %v", err)
				}
			}()

			deltas, ok := d.(cache.Deltas)
			if !ok {
				return errors.New("object given as Process argument is not Deltas")
			}

			// only process the latest delta
			newest := deltas.Newest()
			obj := newest.Object

			// filter
			if selector != nil && !selector.Predicate(obj) {
				return
			}
			err = doProcess(obj, newest.Type)
			return
		},
	}

	return cache.New(cfg)
}
