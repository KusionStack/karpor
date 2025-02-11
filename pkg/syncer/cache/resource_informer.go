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

package cache

import (
	"errors"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

// TransformFunc allows for transforming an object before it will be processed
// and put into the controller cache and before the corresponding handlers will
// be called on it.
// TransformFunc (similarly to ResourceEventHandler functions) should be able
// to correctly handle the tombstone of type cache.DeletedFinalStateUnknown
//
// The most common usage pattern is to clean-up some parts of the object to
// reduce component memory usage if a given component doesn't care about them.
// given controller doesn't care for them
type TransformFunc func(interface{}) (interface{}, error)

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
func NewResourceInformer(lw clientgocache.ListerWatcher,
	selector ResourceSelector,
	transform TransformFunc,
	resyncPeriod time.Duration,
	handler ResourceHandler,
	knownObjects clientgocache.KeyListerGetter,
) clientgocache.Controller {
	informerCache := NewResourceCache()
	fifo := clientgocache.NewDeltaFIFOWithOptions(clientgocache.DeltaFIFOOptions{
		KnownObjects:          knownObjects,
		EmitDeltaTypeReplaced: true,
	})

	doProcess := func(obj interface{}, dType clientgocache.DeltaType) error {
		// transform
		if transform != nil {
			if _, ok := obj.(clientgocache.DeletedFinalStateUnknown); !ok {
				transformed, err := transform(obj)
				if err != nil {
					return fmt.Errorf("error transforming object: %v, delta type: %s", err, dType)
				}
				obj = transformed
			}
		}

		switch dType {
		case clientgocache.Sync, clientgocache.Replaced, clientgocache.Added, clientgocache.Updated:
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
		case clientgocache.Deleted:
			if err := informerCache.Delete(obj); err != nil {
				return err
			}
			return handler.OnDelete(obj)
		}
		return nil
	}

	cfg := &clientgocache.Config{
		Queue:            fifo,
		ObjectType:       &unstructured.Unstructured{},
		FullResyncPeriod: resyncPeriod,
		RetryOnError:     true,
		ListerWatcher: &clientgocache.ListWatch{
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

			deltas, ok := d.(clientgocache.Deltas)
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

	return clientgocache.New(cfg)
}

// NewInformerWithTransformer returns a Store and a controller for populating
// the store while also providing event notifications. You should only used
// the returned Store for Get/List operations; Add/Modify/Deletes will cause
// the event notifications to be faulty.
// The given transform function will be called on all objects before they will
// put into the Store and corresponding Add/Modify/Delete handlers will
// be invoked for them.
func NewInformerWithTransformer(
	lw clientgocache.ListerWatcher,
	objType runtime.Object,
	resyncPeriod time.Duration,
	h clientgocache.ResourceEventHandler,
	transformer TransformFunc,
) (clientgocache.Store, clientgocache.Controller) {
	// This will hold the client state, as we know it.
	clientState := clientgocache.NewStore(clientgocache.DeletionHandlingMetaNamespaceKeyFunc)

	return clientState, newInformer(lw, objType, resyncPeriod, h, clientState, transformer)
}

// newInformer returns a controller for populating the store while also
// providing event notifications.
//
// Parameters
//   - lw is list and watch functions for the source of the resource you want to
//     be informed of.
//   - objType is an object of the type that you expect to receive.
//   - resyncPeriod: if non-zero, will re-list this often (you will get OnUpdate
//     calls, even if nothing changed). Otherwise, re-list will be delayed as
//     long as possible (until the upstream source closes the watch or times out,
//     or you stop the controller).
//   - h is the object you want notifications sent to.
//   - clientState is the store you want to populate
func newInformer(
	lw clientgocache.ListerWatcher,
	objType runtime.Object,
	resyncPeriod time.Duration,
	h clientgocache.ResourceEventHandler,
	clientState clientgocache.Store,
	transformer TransformFunc,
) clientgocache.Controller {
	// This will hold incoming changes. Note how we pass clientState in as a
	// KeyLister, that way resync operations will result in the correct set
	// of update/delete deltas.
	fifo := clientgocache.NewDeltaFIFOWithOptions(clientgocache.DeltaFIFOOptions{
		KnownObjects:          clientState,
		EmitDeltaTypeReplaced: true,
	})

	cfg := &clientgocache.Config{
		Queue:            fifo,
		ListerWatcher:    lw,
		ObjectType:       objType,
		FullResyncPeriod: resyncPeriod,
		RetryOnError:     false,

		Process: func(obj interface{}) error {
			if deltas, ok := obj.(clientgocache.Deltas); ok {
				return processDeltas(h, clientState, transformer, deltas)
			}
			return errors.New("object given as Process argument is not Deltas")
		},
	}
	return clientgocache.New(cfg)
}

// Multiplexes updates in the form of a list of Deltas into a Store, and informs
// a given handler of events OnUpdate, OnAdd, OnDelete
func processDeltas(
	// Object which receives event notifications from the given deltas
	handler clientgocache.ResourceEventHandler,
	clientState clientgocache.Store,
	transformer TransformFunc,
	deltas clientgocache.Deltas,
) error {
	// from oldest to newest
	for _, d := range deltas {
		obj := d.Object
		if transformer != nil {
			var err error
			obj, err = transformer(obj)
			if err != nil {
				return err
			}
		}

		switch d.Type {
		case clientgocache.Sync, clientgocache.Replaced, clientgocache.Added, clientgocache.Updated:
			if old, exists, err := clientState.Get(obj); err == nil && exists {
				if err := clientState.Update(obj); err != nil {
					return err
				}
				handler.OnUpdate(old, obj)
			} else {
				if err := clientState.Add(obj); err != nil {
					return err
				}
				handler.OnAdd(obj)
			}
		case clientgocache.Deleted:
			if err := clientState.Delete(obj); err != nil {
				return err
			}
			handler.OnDelete(obj)
		}
	}
	return nil
}
