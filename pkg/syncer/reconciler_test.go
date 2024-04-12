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
	"testing"

	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/infra/search/storage/elasticsearch"
	clusterv1beta1 "github.com/KusionStack/karbour/pkg/kubernetes/apis/cluster/v1beta1"
	searchv1beta1 "github.com/KusionStack/karbour/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karbour/pkg/kubernetes/scheme"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Test_buildClusterConfig(t *testing.T) {
	tests := []struct {
		name    string
		cluster *clusterv1beta1.Cluster
		want    *rest.Config
		wantErr bool
	}{
		{
			"test error",
			&clusterv1beta1.Cluster{},
			nil,
			true,
		},
		{
			name: "test no error",
			cluster: &clusterv1beta1.Cluster{
				Spec: clusterv1beta1.ClusterSpec{
					Access: clusterv1beta1.ClusterAccess{
						Endpoint: "https://localhost:6443",
						CABundle: []byte("ca"),
						Credential: &clusterv1beta1.ClusterAccessCredential{
							Type: clusterv1beta1.CredentialTypeX509Certificate,
							X509: &clusterv1beta1.X509{
								Certificate: []byte("cert"),
								PrivateKey:  []byte("key"),
							},
						},
					},
				},
			},
			want: &rest.Config{
				Host: "https://localhost:6443",
				TLSClientConfig: rest.TLSClientConfig{
					CAData:   []byte("ca"),
					CertData: []byte("cert"),
					KeyData:  []byte("key"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildClusterConfig(tt.cluster)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSyncReconciler_SetupWithManager(t *testing.T) {
	type fields struct {
		storage    storage.Storage
		client     client.Client
		controller controller.Controller
		mgr        MultiClusterSyncManager
	}
	type args struct {
		mgr controllerruntime.Manager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SyncReconciler{}
			err := r.SetupWithManager(nil)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSyncReconciler_CreateEvent(t *testing.T) {
	tests := []struct {
		name string
		ce   event.CreateEvent
		want reconcile.Request
	}{
		{
			name: "test normal",
			ce: event.CreateEvent{
				Object: &searchv1beta1.SyncRegistry{
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"cluster1"},
					},
				},
			},
			want: reconcile.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &fakeQueue{}
			cluster1 := &clusterv1beta1.Cluster{}
			cluster1.Name = "cluster1"
			r := &SyncReconciler{
				client: fake.NewClientBuilder().WithRuntimeObjects(cluster1).WithScheme(scheme.Scheme).Build(),
			}
			r.CreateEvent(tt.ce, q)
			require.Equal(t, tt.want, q.queue[0])
		})
	}
}

func TestSyncReconciler_UpdateEvent(t *testing.T) {
	tests := []struct {
		name string
		ue   event.UpdateEvent
		want []interface{}
	}{
		{
			name: "test changed",
			ue: event.UpdateEvent{
				ObjectOld: &searchv1beta1.SyncRegistry{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters:             []string{"cluster1"},
						SyncResourcesRefName: "sr-old",
					},
				},
				ObjectNew: &searchv1beta1.SyncRegistry{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "2",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters:             []string{"cluster1"},
						SyncResourcesRefName: "sr-new",
					},
				},
			},
			want: []interface{}{reconcile.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}}},
		},
		{
			name: "test unchanged",
			ue: event.UpdateEvent{
				ObjectOld: &searchv1beta1.SyncRegistry{
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"cluster1"},
					},
				},
				ObjectNew: &searchv1beta1.SyncRegistry{
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"cluster1"},
					},
				},
			},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &fakeQueue{queue: []interface{}{}}
			cluster1 := &clusterv1beta1.Cluster{}
			cluster1.Name = "cluster1"
			r := &SyncReconciler{
				client: fake.NewClientBuilder().WithRuntimeObjects(cluster1).WithScheme(scheme.Scheme).Build(),
			}
			r.UpdateEvent(tt.ue, q)
			require.Equal(t, tt.want, q.queue)
		})
	}
}

func TestSyncReconciler_DeleteEvent(t *testing.T) {
	tests := []struct {
		name string
		de   event.DeleteEvent
		want reconcile.Request
	}{
		{
			name: "test normally",
			de: event.DeleteEvent{
				Object: &searchv1beta1.SyncRegistry{
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"cluster1"},
					},
				},
			},
			want: reconcile.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &fakeQueue{}
			cluster1 := &clusterv1beta1.Cluster{}
			cluster1.Name = "cluster1"
			r := &SyncReconciler{
				client: fake.NewClientBuilder().WithRuntimeObjects(cluster1).WithScheme(scheme.Scheme).Build(),
			}
			r.DeleteEvent(tt.de, q)
			require.Equal(t, tt.want, q.queue[0])
		})
	}
}

func TestSyncReconciler_Reconcile(t *testing.T) {
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
			r := &SyncReconciler{
				client: fake.NewClientBuilder().WithRuntimeObjects(tt.cluster).WithScheme(scheme.Scheme).Build(),
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

func TestSyncReconciler_getNormalizedResource(t *testing.T) {
	tests := []struct {
		name    string
		rsr     *searchv1beta1.ResourceSyncRule
		tfr     *searchv1beta1.TransformRule
		want    *searchv1beta1.ResourceSyncRule
		wantErr bool
	}{
		{
			name: "test ref exist",
			rsr:  &searchv1beta1.ResourceSyncRule{TransformRefName: "tfr1"},
			tfr:  &searchv1beta1.TransformRule{ObjectMeta: metav1.ObjectMeta{Name: "tfr1"}, Spec: searchv1beta1.TransformRuleSpec{Type: "patch"}},
			want: &searchv1beta1.ResourceSyncRule{TransformRefName: "tfr1", Transform: &searchv1beta1.TransformRuleSpec{Type: "patch"}},
		},
		{
			name: "test ref not exist",
			rsr:  &searchv1beta1.ResourceSyncRule{},
			want: &searchv1beta1.ResourceSyncRule{},
		},
		{
			name:    "test ref miss match",
			rsr:     &searchv1beta1.ResourceSyncRule{TransformRefName: "tfr1"},
			wantErr: true,
		},
		{
			name:    "test both exist",
			rsr:     &searchv1beta1.ResourceSyncRule{TransformRefName: "tfr1", Transform: &searchv1beta1.TransformRuleSpec{Type: "patch"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tfrs []runtime.Object
			if tt.tfr != nil {
				tfrs = append(tfrs, tt.tfr)
			}
			r := &SyncReconciler{
				client: fake.NewClientBuilder().WithRuntimeObjects(tfrs...).WithScheme(scheme.Scheme).Build(),
			}
			got, err := r.getNormalizedResource(context.TODO(), tt.rsr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_isMatched(t *testing.T) {
	type args struct {
		registry *searchv1beta1.SyncRegistry
		cluster  *clusterv1beta1.Cluster
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"test label selector",
			args{
				registry: &searchv1beta1.SyncRegistry{
					Spec: searchv1beta1.SyncRegistrySpec{
						ClusterLabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"k1": "v1"},
						},
					},
				},
				cluster: &clusterv1beta1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "cluster1",
						Labels: map[string]string{"k1": "v1"},
					},
				},
			},
			true,
			false,
		},
		{
			"test name match",
			args{
				registry: &searchv1beta1.SyncRegistry{
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"cluster1"},
						ClusterLabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"k1": "v1"},
						},
					},
				},
				cluster: &clusterv1beta1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "cluster1",
						Labels: map[string]string{"k1": "v1"},
					},
				},
			},
			true,
			false,
		},
		{
			"test wildcard",
			args{
				registry: &searchv1beta1.SyncRegistry{
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"*"},
						ClusterLabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"k1": "v1"},
						},
					},
				},
				cluster: &clusterv1beta1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "cluster1",
						Labels: map[string]string{"k1": "v1"},
					},
				},
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isMatched(tt.args.registry, tt.args.cluster)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSyncReconciler_getRegistries(t *testing.T) {
	tests := []struct {
		name    string
		cluster *clusterv1beta1.Cluster
		srs     []runtime.Object
		want    []searchv1beta1.SyncRegistry
		wantErr bool
	}{
		{
			"test no error",
			&clusterv1beta1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster1",
				},
			},
			[]runtime.Object{
				&searchv1beta1.SyncRegistry{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"cluster1"},
					},
				},
			},
			[]searchv1beta1.SyncRegistry{
				{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{Clusters: []string{"cluster1"}},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithRuntimeObjects(tt.srs...).WithScheme(scheme.Scheme).Build()
			r := &SyncReconciler{
				client: cl,
			}
			got, err := r.getRegistries(context.TODO(), tt.cluster)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewSyncReconciler(t *testing.T) {
	tests := []struct {
		name    string
		storage storage.Storage
	}{
		{
			"test nil",
			nil,
		},
		{
			"test not nil",
			&elasticsearch.Storage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSyncReconciler(tt.storage)
			require.Equal(t, got.storage, tt.storage)
		})
	}
}

func TestSyncReconciler_getNormalizedResources(t *testing.T) {
	tests := []struct {
		name     string
		registry *searchv1beta1.SyncRegistry
		srs      []runtime.Object
		want     map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule
		wantErr  bool
	}{
		{
			name: "test",
			registry: &searchv1beta1.SyncRegistry{
				Spec: searchv1beta1.SyncRegistrySpec{
					SyncResources: []searchv1beta1.ResourceSyncRule{{APIVersion: "v1", Resource: "pods"}},
				},
			},
			want: map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule{
				{"", "v1", "pods"}: {APIVersion: "v1", Resource: "pods"},
			},
			wantErr: false,
		},
		{
			name: "test ref",
			registry: &searchv1beta1.SyncRegistry{
				Spec: searchv1beta1.SyncRegistrySpec{
					SyncResourcesRefName: "sr1",
				},
			},
			srs: []runtime.Object{
				&searchv1beta1.SyncResources{
					ObjectMeta: metav1.ObjectMeta{Name: "sr1"},
					Spec:       searchv1beta1.SyncResourcesSpec{SyncResources: []searchv1beta1.ResourceSyncRule{{APIVersion: "v1", Resource: "pods"}}},
				},
			},
			want: map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule{
				{Group: "", Version: "v1", Resource: "pods"}: {APIVersion: "v1", Resource: "pods"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SyncReconciler{client: fake.NewClientBuilder().WithRuntimeObjects(tt.srs...).WithScheme(scheme.Scheme).Build()}
			got, err := r.getNormalizedResources(context.TODO(), tt.registry)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSyncReconciler_getResources(t *testing.T) {
	tests := []struct {
		name    string
		cluster *clusterv1beta1.Cluster
		srs     []runtime.Object
		want    []*searchv1beta1.ResourceSyncRule
		wantErr bool
	}{
		{
			name: "test no error",
			cluster: &clusterv1beta1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster1",
				},
			},
			srs: []runtime.Object{
				&searchv1beta1.SyncRegistry{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters:      []string{"cluster1"},
						SyncResources: []searchv1beta1.ResourceSyncRule{{APIVersion: "v1", Resource: "pods"}},
					},
				},
			},
			want: []*searchv1beta1.ResourceSyncRule{{APIVersion: "v1", Resource: "pods"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SyncReconciler{client: fake.NewClientBuilder().WithRuntimeObjects(tt.srs...).WithScheme(scheme.Scheme).Build()}
			got, err := r.getResources(context.TODO(), tt.cluster)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSyncReconciler_handleClusterAddOrUpdate(t *testing.T) {
	tests := []struct {
		name    string
		cluster *clusterv1beta1.Cluster
		srs     []runtime.Object
		config  *rest.Config
		exist   bool
		wantErr bool
	}{
		{
			name: "test exist",
			cluster: &clusterv1beta1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster1",
				},
				Spec: clusterv1beta1.ClusterSpec{
					Access: clusterv1beta1.ClusterAccess{
						Endpoint: "https://localhost:6443",
						CABundle: []byte("ca"),
						Credential: &clusterv1beta1.ClusterAccessCredential{
							Type: clusterv1beta1.CredentialTypeX509Certificate,
							X509: &clusterv1beta1.X509{
								Certificate: []byte("cert"),
								PrivateKey:  []byte("key"),
							},
						},
					},
				},
			},
			srs: []runtime.Object{
				&searchv1beta1.SyncRegistry{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters:      []string{"cluster1"},
						SyncResources: []searchv1beta1.ResourceSyncRule{{APIVersion: "v1", Resource: "pods"}},
					},
				},
			},
			config: &rest.Config{
				Host: "https://localhost:6443",
				TLSClientConfig: rest.TLSClientConfig{
					CAData:   []byte("ca"),
					CertData: []byte("cert"),
					KeyData:  []byte("key"),
				},
			},
			exist: true,
		},
		{
			name: "test not exist",
			cluster: &clusterv1beta1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster1",
				},
				Spec: clusterv1beta1.ClusterSpec{
					Access: clusterv1beta1.ClusterAccess{
						Endpoint: "https://localhost:6443",
						CABundle: []byte("ca"),
						Credential: &clusterv1beta1.ClusterAccessCredential{
							Type: clusterv1beta1.CredentialTypeX509Certificate,
							X509: &clusterv1beta1.X509{
								Certificate: []byte("cert"),
								PrivateKey:  []byte("key"),
							},
						},
					},
				},
			},
			srs: []runtime.Object{
				&searchv1beta1.SyncRegistry{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters:      []string{"cluster1"},
						SyncResources: []searchv1beta1.ResourceSyncRule{{APIVersion: "v1", Resource: "pods"}},
					},
				},
			},
			config: &rest.Config{
				Host: "https://localhost:6443",
				TLSClientConfig: rest.TLSClientConfig{
					CAData:   []byte("ca"),
					CertData: []byte("cert"),
					KeyData:  []byte("key"),
				},
			},
			exist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1 := &mock.Mock{}
			m1.On("UpdateSyncResources", mock.Anything, mock.Anything).Return(nil)
			m1.On("ClusterConfig").Return(tt.config)
			m2 := &mock.Mock{}
			m2.On("GetForCluster", mock.Anything).Return(&fakeSingleClusterSyncManager{m1}, tt.exist)
			m2.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(&fakeSingleClusterSyncManager{m1}, nil)
			m2.On("Start", mock.Anything, mock.Anything).Return(nil)
			r := &SyncReconciler{
				mgr:    &fakeMultiClusterSyncManager{m2},
				client: fake.NewClientBuilder().WithRuntimeObjects(tt.srs...).WithScheme(scheme.Scheme).Build(),
			}
			err := r.handleClusterAddOrUpdate(context.TODO(), tt.cluster)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
