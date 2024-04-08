package syncer

import (
	"context"
	"testing"

	"github.com/KusionStack/karbour/pkg/kubernetes/apis/search/v1beta1"
	searchv1beta1 "github.com/KusionStack/karbour/pkg/kubernetes/apis/search/v1beta1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
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

func Test_singleClusterSyncManager_UpdateSyncResources(t *testing.T) {
	tests := []struct {
		name          string
		syncResources []*searchv1beta1.ResourceSyncRule
		wantErr       bool
	}{
		{
			name: "test1",
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
