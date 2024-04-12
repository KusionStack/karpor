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

package syncer

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/KusionStack/karbour/pkg/kubernetes/apis/search/v1beta1"
	searchv1beta1 "github.com/KusionStack/karbour/pkg/kubernetes/apis/search/v1beta1"
	"github.com/bytedance/mockey"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var _ SingleClusterSyncManager = &fakeSingleClusterSyncManager{}

type fakeSingleClusterSyncManager struct {
	mock *mock.Mock
}

func (f *fakeSingleClusterSyncManager) Start(ctx context.Context) error {
	args := f.mock.Called(ctx)
	return args.Error(0)
}

func (f *fakeSingleClusterSyncManager) Started() bool {
	args := f.mock.Called()
	return args.Bool(0)
}

func (f *fakeSingleClusterSyncManager) Stop(ctx context.Context) {
	f.mock.Called(ctx)
	return
}

func (f *fakeSingleClusterSyncManager) Stopped() bool {
	args := f.mock.Called()
	return args.Bool(0)
}

func (f *fakeSingleClusterSyncManager) UpdateSyncResources(ctx context.Context, rules []*v1beta1.ResourceSyncRule) error {
	args := f.mock.Called(ctx, rules)
	return args.Error(0)
}

func (f *fakeSingleClusterSyncManager) HasSyncResource(resource schema.GroupVersionResource) bool {
	args := f.mock.Called(resource)
	return args.Bool(0)
}

func (f *fakeSingleClusterSyncManager) ClusterConfig() *rest.Config {
	args := f.mock.Called()
	if arg := args.Get(0); arg == nil {
		return nil
	} else {
		return arg.(*rest.Config)
	}
}

var _ controller.Controller = &fakeController{}

type fakeController struct {
	mock *mock.Mock
}

func (f *fakeController) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	args := f.mock.Called(ctx, request)
	return args.Get(0).(reconcile.Result), args.Error(1)
}

func (f *fakeController) Watch(src source.Source, eventhandler handler.EventHandler, predicates ...predicate.Predicate) error {
	args := f.mock.Called(src, eventhandler, predicates)
	return args.Error(0)
}

func (f *fakeController) Start(ctx context.Context) error {
	args := f.mock.Called(ctx)
	return args.Error(0)
}

func (f *fakeController) GetLogger() logr.Logger {
	args := f.mock.Called()
	return args.Get(0).(logr.Logger)
}

func Test_singleClusterSyncManager_UpdateSyncResources(t *testing.T) {
	tests := []struct {
		name          string
		syncResources []*searchv1beta1.ResourceSyncRule
		wantErr       bool
	}{
		{
			name: "test no error",
			syncResources: []*searchv1beta1.ResourceSyncRule{
				{APIVersion: "v1", Resource: "pods"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &singleClusterSyncManager{
				ch: make(chan struct{}),
			}
			go func() {
				<-s.ch
			}()
			err := s.UpdateSyncResources(context.TODO(), tt.syncResources)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_singleClusterSyncManager_Start(t *testing.T) {
	m := mockey.Mock((*singleClusterSyncManager).process).Return().Build()
	defer m.UnPatch()
	tests := []struct {
		name    string
		wantErr bool
		started bool
	}{
		{
			"test no error",
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &singleClusterSyncManager{logger: klog.NewKlogr()}
			err := s.Start(context.TODO())
			time.Sleep(1 * time.Second)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.started, s.Started())
			}
		})
	}
}

func Test_singleClusterSyncManager_Stop(t *testing.T) {
	tests := []struct {
		name    string
		stopped bool
	}{
		{
			"test stop",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			s := &singleClusterSyncManager{
				logger: klog.NewKlogr(),
				ch:     make(chan struct{}),
				ctx:    ctx,
				cancel: cancel,
			}
			s.Stop(context.TODO())
			require.Equal(t, tt.stopped, s.Stopped())
		})
	}
}

func Test_singleClusterSyncManager_process(t *testing.T) {
	tests := []struct {
		name     string
		prepFunc func(*singleClusterSyncManager)
		called   bool
	}{
		{
			"test cancel",
			func(s *singleClusterSyncManager) {
				s.cancel()
			},
			false,
		},
		{
			"test close",
			func(s *singleClusterSyncManager) {
				close(s.ch)
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			m := mockey.Mock((*singleClusterSyncManager).handleSyncResourcesUpdate).To(func(ctx context.Context) error {
				called = true
				return nil
			}).Build()
			defer m.UnPatch()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			s := &singleClusterSyncManager{
				logger: klog.NewKlogr(),
				ch:     make(chan struct{}),
				ctx:    ctx,
				cancel: cancel,
			}
			tt.prepFunc(s)
			s.process()
			require.Equal(t, tt.called, called)
		})
	}
}

func Test_singleClusterSyncManager_handleSyncResourcesUpdate(t *testing.T) {
	tests := []struct {
		name    string
		gvr     schema.GroupVersionResource
		rsr     *searchv1beta1.ResourceSyncRule
		wantErr bool
	}{
		{
			name: "test normally",
			gvr:  schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			rsr:  &searchv1beta1.ResourceSyncRule{APIVersion: "v1", Resource: "pods"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncResources := func() atomic.Value {
				val := atomic.Value{}
				val.Store(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule{
					tt.gvr: tt.rsr,
				})
				return val
			}()

			src := func() SyncSource {
				m1 := &mock.Mock{}
				cache := &fakeCache{m1}
				ctx1, cancel1 := context.WithCancel(context.Background())
				ch := make(chan struct{})
				close(ch)
				return &informerSource{cache: cache, ctx: ctx1, cancel: cancel1, stopped: ch}
			}()

			syncers := func() sync.Map {
				q := &fakeQueue{}
				ctx2, cancel2 := context.WithCancel(context.Background())
				rs := &ResourceSyncer{queue: q, source: src, ctx: ctx2, cancel: cancel2}
				sm := sync.Map{}
				sm.Store(tt.gvr, rs)
				return sm
			}()

			controller := func() controller.Controller {
				m2 := &mock.Mock{}
				m2.On("Watch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &fakeController{m2}
			}()

			s := &singleClusterSyncManager{
				ctx:           context.TODO(),
				controller:    controller,
				syncResources: syncResources,
				syncers:       syncers,
				stopped:       false,
				logger:        klog.NewKlogr(),
			}

			m := mockey.Mock((*wait.Group).StartWithContext).Return().Build()
			defer m.UnPatch()

			err := s.handleSyncResourcesUpdate(context.TODO())
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
