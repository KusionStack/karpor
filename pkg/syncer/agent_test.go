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

//nolint:dupl
package syncer

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	clusterv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	searchv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
)

func TestAgentReconciler_Reconcile(t *testing.T) {
	tests := []struct {
		name    string
		cluster *clusterv1beta1.Cluster
		req     reconcile.Request
		wantErr bool
	}{
		{
			"test no error",
			&clusterv1beta1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster1"}},
			reconcile.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AgentReconciler{
				SyncReconciler: SyncReconciler{
					client: fake.NewClientBuilder().WithRuntimeObjects(tt.cluster).WithScheme(scheme.Scheme).Build(),
				},
			}
			m := mockey.Mock((*SyncReconciler).handleClusterAddOrUpdate).Return(nil).Build()
			defer m.UnPatch()
			_, err := r.Reconcile(context.TODO(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewAgentReconciler(t *testing.T) {
	tests := []struct {
		name        string
		storage     storage.ResourceStorage
		clusterName string
	}{
		{
			"test nil",
			nil,
			"example-cluster",
		},
		{
			"test not nil",
			&elasticsearch.Storage{},
			"example-cluster2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAgentReconciler(tt.storage, tt.clusterName)
			require.Equal(t, got.storage, tt.storage)
			require.Equal(t, got.clusterName, tt.clusterName)
		})
	}
}

func TestAgentReconciler_getResources(t *testing.T) {
	tests := []struct {
		name             string
		cluster          *clusterv1beta1.Cluster
		registries       []searchv1beta1.SyncRegistry
		resources        map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule
		wildResources    map[string]*searchv1beta1.ResourceSyncRule
		allResources     []*searchv1beta1.ResourceSyncRule
		pendingWildcards []*searchv1beta1.ResourceSyncRule
		wantErr          bool
	}{
		{
			"test no error",
			&clusterv1beta1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster1"}},
			[]searchv1beta1.SyncRegistry{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "example-registry1"},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"cluster1"},
					},
				},
			},
			map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule{
				v1.SchemeGroupVersion.WithResource("pods"): {
					APIVersion: "v1",
					Resource:   "pods",
				},
			},
			map[string]*searchv1beta1.ResourceSyncRule{
				"v1": {
					APIVersion: "v1",
					Resource:   "*",
				},
			},
			[]*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "v1",
					Resource:   "pods",
				},
			},
			[]*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "v1",
					Resource:   "*",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AgentReconciler{
				SyncReconciler: SyncReconciler{
					client: fake.NewClientBuilder().WithRuntimeObjects(tt.cluster).WithScheme(scheme.Scheme).Build(),
				},
			}
			r.gvrToGVKCache.Store(v1.SchemeGroupVersion.WithResource("pods"), schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"})

			mockey.Mock((*SyncReconciler).getRegistries).Return(tt.registries, nil).Build()
			mockey.Mock((*SyncReconciler).getNormalizedResources).Return(tt.resources, tt.wildResources, nil).Build()
			defer mockey.UnPatchAll()
			allResources, pendingWildcards, err := r.getResources(context.TODO(), tt.cluster)
			require.Equal(t, allResources, tt.allResources)
			require.Equal(t, pendingWildcards, tt.pendingWildcards)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
