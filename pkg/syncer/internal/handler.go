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

package internal

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	ctrlhandler "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// EventHandler adapts a controller-runtime's EventHandler interface to a cache.ResourceHandler interface.
type EventHandler struct {
	EventHandler ctrlhandler.EventHandler
	Queue        workqueue.RateLimitingInterface
	Predicates   []predicate.Predicate
}

func (e EventHandler) OnAdd(obj interface{}) {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		klog.Errorf("invalid object type. Expected *unstructured.Unstructured, but got %T", obj)
	}

	ce := event.CreateEvent{Object: u}
	for _, p := range e.Predicates {
		if !p.Create(ce) {
			return
		}
	}

	e.EventHandler.Create(ce, e.Queue)
}

func (e EventHandler) OnUpdate(oldObj, newObj interface{}) {
	u, ok := newObj.(*unstructured.Unstructured)
	if !ok {
		klog.Errorf("invalid object type. Expected *unstructured.Unstructured, but got %T", newObj)
		return
	}

	ue := event.UpdateEvent{ObjectNew: u}
	for _, p := range e.Predicates {
		if !p.Update(ue) {
			return
		}
	}

	e.EventHandler.Update(ue, e.Queue)
}

func (e EventHandler) OnDelete(obj interface{}) {
	var ok bool
	if _, ok = obj.(client.Object); !ok {
		tombstone, ok := obj.(clientgocache.DeletedFinalStateUnknown)
		if !ok {
			klog.Errorf("invalid object type. Expected cache.DeletedFinalStateUnknown, but got %T", obj)
			return
		}
		obj = tombstone.Obj
	}

	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		klog.Errorf("invalid object type. Expected *unstructured.Unstructured, but got %T", obj)
		return
	}

	d := event.DeleteEvent{Object: u}

	for _, p := range e.Predicates {
		if !p.Delete(d) {
			return
		}
	}

	e.EventHandler.Delete(d, e.Queue)
}
