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
	"testing"
	"time"

	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	clientgocache "k8s.io/client-go/tools/cache"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllertest"
)

func Test_informerSource_Add(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("Add", mock.Anything).Return(nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		pod := &corev1.Pod{}
		s.Add(pod)
		m.AssertCalled(t, "Add", pod)
	})
}

func Test_informerSource_Update(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("Update", mock.Anything).Return(nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		pod := &corev1.Pod{}
		s.Update(pod)
		m.AssertCalled(t, "Update", pod)
	})
}

func Test_informerSource_Delete(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("Delete", mock.Anything).Return(nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		pod := &corev1.Pod{}
		s.Delete(pod)
		m.AssertCalled(t, "Delete", pod)
	})
}

func Test_informerSource_List(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("List", mock.Anything).Return(nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		s.List()
		m.AssertCalled(t, "List")
	})
}

func Test_informerSource_ListKeys(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("ListKeys", mock.Anything).Return(nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		s.ListKeys()
		m.AssertCalled(t, "ListKeys")
	})
}

func Test_informerSource_Get(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("Get", mock.Anything).Return(nil, true, nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		pod := &corev1.Pod{}
		s.Get(pod)
		m.AssertCalled(t, "Get", pod)
	})
}

func Test_informerSource_GetByKey(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("GetByKey", mock.Anything).Return(nil, true, nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		s.GetByKey("")
		m.AssertCalled(t, "GetByKey", "")
	})
}

func Test_informerSource_Replace(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("Replace", mock.Anything, mock.Anything).Return(nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		s.Replace(nil, "")
		m.AssertCalled(t, "Replace", mock.Anything, mock.Anything)
	})
}

func Test_informerSource_Resync(t *testing.T) {
	t.Run("test function called", func(t *testing.T) {
		m := &mock.Mock{}
		m.On("Resync", mock.Anything).Return(nil)
		cache := &fakeCache{m}
		s := &informerSource{cache: cache}
		s.Resync()
		m.AssertCalled(t, "Resync")
	})
}

var _ clientgocache.Store = &fakeCache{}

type fakeCache struct {
	mock *mock.Mock
}

func (f *fakeCache) Add(obj interface{}) error {
	args := f.mock.Called(obj)
	return args.Error(0)
}

func (f *fakeCache) Update(obj interface{}) error {
	args := f.mock.Called(obj)
	return args.Error(0)
}

func (f *fakeCache) Delete(obj interface{}) error {
	args := f.mock.Called(obj)
	return args.Error(0)
}

func (f *fakeCache) List() []interface{} {
	args := f.mock.Called()
	if args[0] == nil {
		return nil
	}
	return args[0].([]interface{})
}

func (f *fakeCache) ListKeys() []string {
	args := f.mock.Called()
	if args[0] == nil {
		return nil
	}
	return args[0].([]string)
}

func (f *fakeCache) Get(obj interface{}) (item interface{}, exists bool, err error) {
	args := f.mock.Called(obj)
	return args[0], args.Bool(1), args.Error(2)
}

func (f *fakeCache) GetByKey(key string) (item interface{}, exists bool, err error) {
	args := f.mock.Called(key)
	return args[0], args.Bool(1), args.Error(2)
}

func (f *fakeCache) Replace(i []interface{}, s string) error {
	args := f.mock.Called(i, s)
	return args.Error(0)
}

func (f *fakeCache) Resync() error {
	args := f.mock.Called()
	return args.Error(0)
}

func Test_informerSource_Start(t *testing.T) {
	t.Run("test no error", func(t *testing.T) {
		mockey.Mock((*utils.ESImporter).ImportTo).Return(nil).Build()
		informer := &controllertest.FakeInformer{}
		mockey.Mock(clientgocache.NewTransformingInformer).Return(clientgocache.NewStore(clientgocache.DeletionHandlingMetaNamespaceKeyFunc), informer).Build()
		defer mockey.UnPatchAll()
		s := &informerSource{
			ResourceSyncRule: v1beta1.ResourceSyncRule{APIVersion: "v1", Resource: "pods"},
			storage:          &elasticsearch.Storage{},
			stopped:          make(chan struct{}),
		}
		err := s.Start(context.TODO(), nil, nil)
		require.NoError(t, err)
	})
}

func Test_informerSource_parseTransformer(t *testing.T) {
	t.Run("test no error", func(t *testing.T) {
		s := &informerSource{
			ResourceSyncRule: v1beta1.ResourceSyncRule{Transform: &v1beta1.TransformRuleSpec{Type: "patch"}},
		}
		_, err := s.parseTransformer()
		require.NoError(t, err)
	})
}

func Test_parseSelectors(t *testing.T) {
	tests := []struct {
		name    string
		rsr     v1beta1.ResourceSyncRule
		wantErr bool
	}{
		{
			name:    "test nil",
			rsr:     v1beta1.ResourceSyncRule{},
			wantErr: false,
		},
		{
			name: "test label selector",
			rsr: v1beta1.ResourceSyncRule{
				Selectors: []v1beta1.Selector{
					{LabelSelector: &metav1.LabelSelector{}},
				},
			},
			wantErr: false,
		},
		{
			name: "test field selector",
			rsr: v1beta1.ResourceSyncRule{
				Selectors: []v1beta1.Selector{
					{FieldSelector: nil},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseSelectors(tt.rsr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_informerSource_Stop(t *testing.T) {
	t.Run("test timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		s := &informerSource{
			stopped: make(chan struct{}),
			ctx:     ctx,
			cancel:  cancel,
		}
		ctx2, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
		defer cancel()
		err := s.Stop(ctx2)
		require.Error(t, err)
	})
	t.Run("test context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		s := &informerSource{
			stopped: make(chan struct{}),
			ctx:     ctx,
			cancel:  cancel,
		}
		ctx2, cancel2 := context.WithCancel(context.TODO())
		go func() {
			time.AfterFunc(1*time.Second, func() {
				cancel2()
			})
		}()
		err := s.Stop(ctx2)
		require.NoError(t, err)
	})
	t.Run("test stop", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		s := &informerSource{
			stopped: make(chan struct{}),
			ctx:     ctx,
			cancel:  cancel,
		}
		ctx2, cancel := context.WithCancel(context.TODO())
		defer cancel()
		go func() {
			time.AfterFunc(1*time.Second, func() {
				close(s.stopped)
			})
		}()
		err := s.Stop(ctx2)
		require.NoError(t, err)
	})
}

func Test_informerSource_HasSynced(t *testing.T) {
	s := &informerSource{
		informer: &controllertest.FakeInformer{},
	}
	res := s.HasSynced()
	require.False(t, res)
}
