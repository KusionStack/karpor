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
	"crypto"
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	clusterv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	searchv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"github.com/KusionStack/karpor/pkg/syncer/utils"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func Test_buildClusterConfigInSyncer(t *testing.T) {
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
			got, err := buildClusterConfigInSyncer(tt.cluster)
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
	tests := []struct {
		name    string
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
		name             string
		storage          storage.ResourceStorage
		highAvailability bool
		storageAddresses []string
		externalEndpoint string
		agentImageTag    string
		caCert           *x509.Certificate
		caKey            crypto.Signer
	}{
		{
			"test nil",
			nil,
			false,
			[]string{"127.0.0.1"},
			"127.0.0.1",
			"v1.0.0",
			nil,
			nil,
		},
		{
			"test not nil",
			&elasticsearch.Storage{},
			false,
			[]string{"127.0.0.1"},
			"127.0.0.1",
			"v1.0.0",
			nil,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSyncReconciler(tt.storage, tt.highAvailability, tt.storageAddresses,
				tt.externalEndpoint, tt.agentImageTag)
			require.Equal(t, got.storage, tt.storage)
			require.Equal(t, got.highAvailability, tt.highAvailability)
			require.Equal(t, got.storageAddresses, tt.storageAddresses)
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
				{Group: "", Version: "v1", Resource: "pods"}: {APIVersion: "v1", Resource: "pods"},
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
			got, _, err := r.getNormalizedResources(context.TODO(), tt.registry)
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
			allResources, _, err := r.getResources(context.TODO(), tt.cluster)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, allResources)
			}
		})
	}
}

func TestSyncReconciler_processWildcardResources(t *testing.T) {
	tests := []struct {
		name         string
		wildcards    []*searchv1beta1.ResourceSyncRule
		apiResources []metav1.APIResource
		clusterName  string
		want         []*searchv1beta1.ResourceSyncRule
		wantErr      bool
	}{
		{
			name: "successfully process wildcard resources",
			wildcards: []*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "samplecontroller.k8s.io/v1alpha1",
					Resource:   "*",
				},
			},
			apiResources: []metav1.APIResource{
				{
					Name:       "foos",
					Namespaced: true,
				},
				{
					Name:       "bars",
					Namespaced: false,
				},
			},
			clusterName: "cluster1",
			want: []*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "samplecontroller.k8s.io/v1alpha1",
					Resource:   "foos",
				},
				{
					APIVersion: "samplecontroller.k8s.io/v1alpha1",
					Resource:   "bars",
				},
			},
			wantErr: false,
		},
		{
			name: "skip subresources",
			wildcards: []*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "samplecontroller.k8s.io/v1alpha1",
					Resource:   "*",
				},
			},
			apiResources: []metav1.APIResource{
				{
					Name:       "foos/status",
					Namespaced: true,
				},
				{
					Name:       "foos",
					Namespaced: true,
				},
			},
			clusterName: "cluster1",
			want: []*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "samplecontroller.k8s.io/v1alpha1",
					Resource:   "foos",
				},
			},
			wantErr: false,
		},
		{
			name: "skip cluster-scoped resources when namespace specified",
			wildcards: []*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "samplecontroller.k8s.io/v1alpha1",
					Resource:   "*",
					Namespace:  "default",
				},
			},
			apiResources: []metav1.APIResource{
				{
					Name:       "bars",
					Namespaced: false,
				},
				{
					Name:       "foos",
					Namespaced: true,
				},
			},
			clusterName: "cluster1",
			want: []*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "samplecontroller.k8s.io/v1alpha1",
					Resource:   "foos",
					Namespace:  "default",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid group version",
			wildcards: []*searchv1beta1.ResourceSyncRule{
				{
					APIVersion: "invalid",
					Resource:   "*",
				},
			},
			apiResources: []metav1.APIResource{},
			clusterName:  "cluster1",
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockManager := &mock.Mock{}
			if tt.wantErr {
				mockManager.On("GetAPIResources", mock.Anything).Return(nil, fmt.Errorf("failed to get API resources"))
			} else {
				mockManager.On("GetAPIResources", mock.Anything).Return(
					&metav1.APIResourceList{
						GroupVersion: tt.wildcards[0].APIVersion,
						APIResources: tt.apiResources,
					},
					nil,
				)
			}

			r := &SyncReconciler{}
			got, err := r.processWildcardResources(context.TODO(), tt.wildcards, &fakeSingleClusterSyncManager{mockManager}, tt.clusterName)

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
			err := r.handleClusterAddOrUpdate(context.TODO(), tt.cluster, buildClusterConfigInSyncer)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSyncReconciler_dispatchResources(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		cluster       *clusterv1beta1.Cluster
		dynamicClient dynamic.Interface
		objects       []runtime.Object
		wantErr       bool
	}{
		{
			name: "test no error",
			ctx:  context.Background(),
			cluster: &clusterv1beta1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster1",
				},
				Spec: clusterv1beta1.ClusterSpec{
					Mode:  clusterv1beta1.PullClusterMode,
					Level: 2,
				},
			},
			dynamicClient: &dynamicfake.FakeDynamicClient{},
			objects: []runtime.Object{
				&searchv1beta1.SyncRegistry{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters:      []string{"cluster1"},
						SyncResources: []searchv1beta1.ResourceSyncRule{{APIVersion: "v1", Resource: "pods", TransformRefName: "tfr1", TrimRefName: "tr1"}},
					},
				},
				&searchv1beta1.TransformRule{ObjectMeta: metav1.ObjectMeta{Name: "tfr1"}, Spec: searchv1beta1.TransformRuleSpec{}},
				&searchv1beta1.TrimRule{ObjectMeta: metav1.ObjectMeta{Name: "tr1"}, Spec: searchv1beta1.TrimRuleSpec{}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.Mock(utils.CreateOrUpdateUnstructured).Return(nil).Build()
			mockey.Mock((*SyncReconciler).getUnstructuredRegistries).Return(nil).Build()
			defer mockey.UnPatchAll()

			r := &SyncReconciler{client: fake.NewClientBuilder().WithRuntimeObjects(tt.objects...).WithScheme(scheme.Scheme).Build()}
			err := r.dispatchResources(tt.ctx, tt.dynamicClient, tt.cluster)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSyncReconciler_getUnstructuredRegistries(t *testing.T) {
	tests := []struct {
		name                  string
		ctx                   context.Context
		cluster               *clusterv1beta1.Cluster
		unstructuredObjectMap map[schema.GroupVersionResource][]unstructured.Unstructured
		objects               []runtime.Object
		wantErr               bool
	}{
		{
			name: "test transform and trim rule",
			cluster: &clusterv1beta1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster1",
				},
				Spec: clusterv1beta1.ClusterSpec{
					Mode:  clusterv1beta1.PullClusterMode,
					Level: 2,
				},
			},
			unstructuredObjectMap: map[schema.GroupVersionResource][]unstructured.Unstructured{},
			objects: []runtime.Object{
				&searchv1beta1.SyncRegistry{
					ObjectMeta: metav1.ObjectMeta{
						ResourceVersion: "1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters:      []string{"cluster1"},
						SyncResources: []searchv1beta1.ResourceSyncRule{{APIVersion: "v1", Resource: "pods", TransformRefName: "tfr1", TrimRefName: "tr1"}},
					},
				},
				&searchv1beta1.TransformRule{ObjectMeta: metav1.ObjectMeta{Name: "tfr1"}, Spec: searchv1beta1.TransformRuleSpec{}},
				&searchv1beta1.TrimRule{ObjectMeta: metav1.ObjectMeta{Name: "tr1"}, Spec: searchv1beta1.TrimRuleSpec{}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SyncReconciler{client: fake.NewClientBuilder().WithRuntimeObjects(tt.objects...).WithScheme(scheme.Scheme).Build()}
			err := r.getUnstructuredRegistries(tt.ctx, tt.cluster, tt.unstructuredObjectMap)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSyncReconciler_renderYamlFile(t *testing.T) {
	tests := []struct {
		name     string
		cluster  *clusterv1beta1.Cluster
		certData string
		keyData  string
		want     string
		wantErr  bool
	}{
		{
			name: "test pull mode",
			cluster: &clusterv1beta1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster1",
				},
				Spec: clusterv1beta1.ClusterSpec{
					Mode:  clusterv1beta1.PushClusterMode,
					Level: 2,
				},
			},
			certData: "cert",
			keyData:  "key",
			want:     renderResForPush,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SyncReconciler{
				storageAddresses: []string{"https://localhost:6443"},
				agentImageTag:    "latest",
				externalEndpoint: "https://localhost:6443",
			}
			got, err := r.renderYamlFile(tt.cluster)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

const (
	renderResForPush = "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: karpor\nspec:\n  finalizers:\n  - kubernetes\n---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: karpor-agent\n  namespace: karpor\nspec:\n  replicas: 1\n  revisionHistoryLimit: 10\n  selector:\n    matchLabels:\n      app.kubernetes.io/component: karpor-agent\n      app.kubernetes.io/instance: karpor\n      app.kubernetes.io/name: karpor\n  strategy:\n    rollingUpdate:\n      maxSurge: 25%\n      maxUnavailable: 25%\n    type: RollingUpdate\n  template:\n    metadata:\n      labels:\n        app.kubernetes.io/component: karpor-agent\n        app.kubernetes.io/instance: karpor\n        app.kubernetes.io/name: karpor\n    spec:\n      containers:\n      - args:\n        - agent\n        - --elastic-search-addresses=https://localhost:6443 \n        - --cluster-name=cluster1\n        - --cluster-mode=push\n        command:\n        - /karpor\n        image: kusionstack/karpor:latest\n        imagePullPolicy: IfNotPresent\n        name: karpor-agent\n        ports:\n        - containerPort: 7443\n          protocol: TCP\n        resources:\n          limits:\n            cpu: 500m\n            ephemeral-storage: 10Gi\n            memory: 1Gi\n          requests:\n            cpu: 250m\n            ephemeral-storage: 2Gi\n            memory: 256Mi\n      dnsPolicy: ClusterFirst\n      restartPolicy: Always\n      terminationGracePeriodSeconds: 30\n---\napiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleBinding\nmetadata:\n  name: karpor\nroleRef:\n  apiGroup: rbac.authorization.k8s.io\n  kind: ClusterRole\n  name: cluster-admin\nsubjects:\n- kind: ServiceAccount\n  name: default\n  namespace: karpor\n"
)
