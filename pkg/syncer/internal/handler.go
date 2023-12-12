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
	"fmt"

	"github.com/KusionStack/karbour/pkg/syncer/cache"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	ctrlhandler "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ cache.ResourceHandler = EventHandler{}

// EventHandler adapts a controller-runtime's EventHandler interface to a cache.ResourceHandler interface.
type EventHandler struct {
	EventHandler ctrlhandler.EventHandler
	Queue        workqueue.RateLimitingInterface
	Predicates   []predicate.Predicate
}

func (e EventHandler) OnAdd(obj interface{}) error {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("invalid object type. Expected *unstructured.Unstructured, but got %T", obj)
	}

	ce := event.CreateEvent{Object: u}
	for _, p := range e.Predicates {
		if !p.Create(ce) {
			return nil
		}
	}

	e.EventHandler.Create(ce, e.Queue)
	return nil
}

func (e EventHandler) OnUpdate(newObj interface{}) error {
	u, ok := newObj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("invalid object type. Expected *unstructured.Unstructured, but got %T", newObj)
	}

	ue := event.UpdateEvent{ObjectNew: u}
	for _, p := range e.Predicates {
		if !p.Update(ue) {
			return nil
		}
	}

	e.EventHandler.Update(ue, e.Queue)
	return nil
}

func (e EventHandler) OnDelete(obj interface{}) error {
	var ok bool
	if _, ok = obj.(client.Object); !ok {
		tombstone, ok := obj.(clientgocache.DeletedFinalStateUnknown)
		if !ok {
			return fmt.Errorf("invalid object type. Expected cache.DeletedFinalStateUnknown, but got %T", obj)
		}
		obj = tombstone.Obj
	}

	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("invalid object type. Expected *unstructured.Unstructured, but got %T", obj)
	}

	d := event.DeleteEvent{Object: u}

	for _, p := range e.Predicates {
		if !p.Delete(d) {
			return nil
		}
	}

	e.EventHandler.Delete(d, e.Queue)
	return nil
}
