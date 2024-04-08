package syncer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

var _ MultiClusterSyncManager = &fakeMultiClusterSyncManager{}

type fakeMultiClusterSyncManager struct {
	mock *mock.Mock
}

func (f *fakeMultiClusterSyncManager) Create(ctx context.Context, clusterName string, config *rest.Config) (SingleClusterSyncManager, error) {
	args := f.mock.Called(ctx, clusterName, config)
	var scym SingleClusterSyncManager
	if args.Get(0) != nil {
		scym = args.Get(0).(SingleClusterSyncManager)
	}
	return scym, args.Error(1)
}

func (f *fakeMultiClusterSyncManager) Start(ctx context.Context, clusterName string) error {
	args := f.mock.Called(ctx, clusterName)
	return args.Error(0)
}

func (f *fakeMultiClusterSyncManager) Stop(ctx context.Context, clusterName string) {
	f.mock.Called(ctx, clusterName)
}

func (f *fakeMultiClusterSyncManager) GetForCluster(clusterName string) (SingleClusterSyncManager, bool) {
	args := f.mock.Called(clusterName)
	var scym SingleClusterSyncManager
	if args.Get(0) != nil {
		scym = args.Get(0).(SingleClusterSyncManager)
	}
	return scym, args.Bool(1)
}

func Test_multiClusterSyncManager_Create(t *testing.T) {
	tests := []struct {
		name    string
		config  *rest.Config
		wantErr bool
	}{
		{
			"test1",
			&rest.Config{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMultiClusterSyncManager(nil, nil, nil)
			_, err := s.Create(context.TODO(), "cluster1", tt.config)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_multiClusterSyncManager_Start(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		ctx := context.TODO()
		cluster := "cluster1"
		m := &mock.Mock{}
		m.On("Start", mock.Anything, mock.Anything).Return(nil)
		s := &multiClusterSyncManager{
			managers: map[string]SingleClusterSyncManager{
				cluster: &fakeSingleClusterSyncManager{m},
			},
		}
		err := s.Start(ctx, cluster)
		require.NoError(t, err)
		m.AssertCalled(t, "Start", mock.Anything, mock.Anything)
	})
}

func Test_multiClusterSyncManager_Stop(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		ctx := context.TODO()
		cluster := "cluster1"
		m := &mock.Mock{}
		m.On("Stop", mock.Anything, mock.Anything).Return()
		s := &multiClusterSyncManager{
			managers: map[string]SingleClusterSyncManager{
				cluster: &fakeSingleClusterSyncManager{m},
			},
		}
		s.Stop(ctx, cluster)
		m.AssertCalled(t, "Stop", mock.Anything, mock.Anything)
	})
}
