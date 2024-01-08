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
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/KusionStack/karbour/pkg/syncer/utils"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic/fake"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
)

type OP uint

const (
	ADD OP = iota
	UPDATE
	DELETE
)

type Record struct {
	op   OP
	data interface{}
}

func event(op OP, data interface{}) *Record {
	return &Record{op: op, data: data}
}

type ResourceRecorder struct {
	events []*Record
	sync.RWMutex
}

func (r *ResourceRecorder) RecordAdd(data interface{}) error {
	r.record(&Record{op: ADD, data: data})
	return nil
}

func (r *ResourceRecorder) RecordUpdate(data interface{}) error {
	r.record(&Record{op: UPDATE, data: data})
	return nil
}

func (r *ResourceRecorder) RecordDelete(data interface{}) error {
	r.record(&Record{op: DELETE, data: data})
	return nil
}

func (r *ResourceRecorder) List() []*Record {
	var events []*Record
	r.read(func() {
		events = append(events, r.events...)
	})
	return events
}

func (r *ResourceRecorder) read(cb func()) {
	r.RLock()
	defer r.RUnlock()
	cb()
}

func (r *ResourceRecorder) record(e *Record) {
	r.Lock()
	defer r.Unlock()
	r.events = append(r.events, e)
}

func (r *ResourceRecorder) resourceHandler() ResourceHandler {
	return ResourceHandlerFuncs{
		AddFunc:    r.RecordAdd,
		UpdateFunc: r.RecordUpdate,
		DeleteFunc: r.RecordDelete,
	}
}

func TestInformerWithSelectors(t *testing.T) {
	assert := assert.New(t)

	objs := []runtime.Object{
		makeUnstructured("default", "pod1", map[string]interface{}{"metadata.labels.foo": "bar1", "status.phase": "RUNNING"}),
		makeUnstructured("default", "pod2", map[string]interface{}{"metadata.labels.foo": "bar2"}),
		makeUnstructured("default", "pod3", map[string]interface{}{"status.phase": "RUNNING"}),
	}

	matchLabels := func(ls map[string]string) labels.Selector {
		s, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{MatchLabels: ls})
		return s
	}

	matchFields := func(fs map[string]string, serverSupported bool) utils.FieldsSelector {
		return utils.FieldsSelector{
			Selector:        fields.SelectorFromSet(fields.Set(fs)),
			ServerSupported: serverSupported,
		}
	}

	tests := []struct {
		name      string
		objs      []runtime.Object
		selectors []utils.Selector
		expected  []*Record
	}{
		{
			"match labels",
			objs,
			[]utils.Selector{
				{Label: matchLabels(map[string]string{"foo": "bar1"})},
			},
			[]*Record{event(ADD, objs[0])},
		},
		{
			"match fields",
			objs,
			[]utils.Selector{
				{Field: matchFields(map[string]string{"metadata.namespace": "default", "status.phase": "RUNNING"}, false)},
			},
			[]*Record{event(ADD, objs[0]), event(ADD, objs[2])},
		},
		{
			"multiple selectors",
			objs,
			[]utils.Selector{
				{Label: matchLabels(map[string]string{"foo": "bar1"})},
				{Label: matchLabels(map[string]string{"foo": "bar2"})},
			},
			[]*Record{event(ADD, objs[0]), event(ADD, objs[1])},
		},
	}

	for _, tt := range tests {
		recorder := new(ResourceRecorder)
		client := fake.NewSimpleDynamicClient(scheme(), objs...)
		gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

		lw := &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.Resource(gvr).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.Resource(gvr).Watch(context.TODO(), options)
			},
		}

		informer := NewResourceInformer(lw, utils.MultiSelectors(tt.selectors), nil, 0, recorder.resourceHandler())
		stop := make(chan struct{})
		defer close(stop)
		go informer.Run(stop)

		wait.Poll(100*time.Millisecond, wait.ForeverTestTimeout, func() (bool, error) {
			return informer.HasSynced(), nil
		})
		assert.True(informer.HasSynced(), "Expected HasSynced() to return true after the initial sync")
		assert.EqualValues(tt.expected, recorder.List())
	}
}

func TestInformerWithTransformer(t *testing.T) {
	assert := assert.New(t)

	objs := []runtime.Object{
		makeUnstructured("default", "pod1", map[string]interface{}{
			"metadata.resourceVersion": "1",
			"metadata.labels.foo":      "bar", // to be excluded
			"status.phase":             "Running",
		}),
		makeUnstructured("default", "pod2", map[string]interface{}{
			"metadata.resourceVersion": "2",
			"spec.foo":                 "bar", // to be excluded
			"status.phase":             "Pending",
		}),
	}

	tests := []struct {
		transFunc cache.TransformFunc
		updates   []struct {
			action string
			data   interface{}
		}
		expected []*Record
	}{
		{
			transFunc: func(source interface{}) (interface{}, error) {
				u, _ := source.(*unstructured.Unstructured)
				dest := u.NewEmptyInstance()
				unstructured.SetNestedField(dest.UnstructuredContent(), u.GetResourceVersion(), "metadata", "resourceVersion")
				unstructured.SetNestedField(dest.UnstructuredContent(), u.GetNamespace(), "metadata", "namespace")
				unstructured.SetNestedField(dest.UnstructuredContent(), u.GetName(), "metadata", "name")

				phase, _, _ := unstructured.NestedString(u.UnstructuredContent(), "status", "phase")
				unstructured.SetNestedField(dest.UnstructuredContent(), phase, "status", "phase")
				return dest, nil
			},
			updates: []struct {
				action string
				data   interface{}
			}{
				{
					"update",
					makeUnstructured("default", "pod1", map[string]interface{}{
						"metadata.resourceVersion": "3",
						"metadata.labels.foo":      "bar2",
						"status.phase":             "Running",
					}),
				},
				{
					"update",
					makeUnstructured("default", "pod2", map[string]interface{}{
						"metadata.resourceVersion": "4",
						"status.phase":             "Failed",
					}),
				},
			},
			expected: []*Record{
				event(ADD,
					makeUnstructured("default", "pod1", map[string]interface{}{
						"metadata.resourceVersion": "1",
						"status.phase":             "Running",
					}),
				),
				event(ADD,
					makeUnstructured("default", "pod2", map[string]interface{}{
						"metadata.resourceVersion": "2",
						"status.phase":             "Pending",
					}),
				),
				event(UPDATE,
					makeUnstructured("default", "pod2", map[string]interface{}{
						"metadata.resourceVersion": "4",
						"status.phase":             "Failed",
					}),
				),
			},
		},
	}

	for _, tt := range tests {
		recorder := new(ResourceRecorder)
		client := fake.NewSimpleDynamicClient(scheme(), objs...)
		gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

		lw := &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.Resource(gvr).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.Resource(gvr).Watch(context.TODO(), options)
			},
		}

		informer := NewResourceInformer(lw, nil, tt.transFunc, 0, recorder.resourceHandler())
		stop := make(chan struct{})
		defer close(stop)
		go informer.Run(stop)

		wait.Poll(100*time.Millisecond, wait.ForeverTestTimeout, func() (bool, error) {
			return informer.HasSynced(), nil
		})
		assert.True(informer.HasSynced(), "Expected HasSynced() to return true after the initial sync")

		for _, e := range tt.updates {
			switch e.action {
			case "create":
				obj := e.data.(*unstructured.Unstructured)
				client.Resource(gvr).Namespace(obj.GetNamespace()).Create(context.TODO(), obj, metav1.CreateOptions{})
			case "update":
				obj := e.data.(*unstructured.Unstructured)
				client.Resource(gvr).Namespace(obj.GetNamespace()).Update(context.TODO(), obj, metav1.UpdateOptions{})
			case "delete":
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(e.data)
				assert.NoError(err)
				parts := strings.Split(key, ",")
				assert.Equal(len(parts), 2)
				client.Resource(gvr).Namespace(parts[0]).Delete(context.TODO(), parts[1], metav1.DeleteOptions{})
			}
		}
		assert.Eventually(
			func() bool {
				return assert.EqualValues(tt.expected, recorder.List())
			},
			100*time.Millisecond, 10*time.Millisecond)
	}
}

func scheme() *runtime.Scheme {
	s := runtime.NewScheme()
	clientgoscheme.AddToScheme(s)
	return s
}

//nolint:unparam
func makeUnstructured(namespace, name string, fields map[string]interface{}) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetAPIVersion("v1")
	u.SetKind("pod")
	u.SetName(name)
	u.SetNamespace(namespace)
	for path, val := range fields {
		unstructured.SetNestedField(u.Object, val, strings.Split(path, ".")...)
	}
	return u
}
