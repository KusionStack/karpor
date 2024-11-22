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

package entity

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceGroupHash(t *testing.T) {
	tests := []struct {
		name         string
		rg           ResourceGroup
		wantHash     string
		wantEqual    []ResourceGroup
		wantNotEqual []ResourceGroup
	}{
		{
			name: "SingleResourceGroup",
			rg: ResourceGroup{
				Cluster:     "test-cluster",
				APIVersion:  "v1",
				Kind:        "Pod",
				Namespace:   "default",
				Name:        "test-pod",
				Labels:      map[string]string{"app": "myapp", "env": "dev"},
				Annotations: map[string]string{"note": "test"},
			},
			wantHash: "test-cluster-v1-Pod-default-test-pod-app:myapp-env:dev-note:test-",
		},
		{
			name: "ResourceGroupWithDifferentLabels",
			rg: ResourceGroup{
				Cluster:     "test-cluster",
				APIVersion:  "v1",
				Kind:        "Pod",
				Namespace:   "default",
				Name:        "test-pod",
				Labels:      map[string]string{"env": "prod", "app": "myapp"},
				Annotations: map[string]string{"note": "test"},
			},
			wantHash: "test-cluster-v1-Pod-default-test-pod-app:myapp-env:prod-note:test-",
			wantEqual: []ResourceGroup{
				{
					Cluster:     "test-cluster",
					APIVersion:  "v1",
					Kind:        "Pod",
					Namespace:   "default",
					Name:        "test-pod",
					Labels:      map[string]string{"app": "myapp", "env": "prod"},
					Annotations: map[string]string{"note": "test"},
				},
			},
			wantNotEqual: []ResourceGroup{
				{
					Cluster:     "test-cluster",
					APIVersion:  "v1",
					Kind:        "Pod",
					Namespace:   "default",
					Name:        "test-pod",
					Labels:      map[string]string{"env": "staging"},
					Annotations: map[string]string{"note": "test"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash := tt.rg.Hash()

			require.Equal(t, tt.wantHash, string(gotHash))

			for _, equalRg := range tt.wantEqual {
				require.Equal(t, equalRg.Hash(), gotHash)
			}

			for _, notEqualRg := range tt.wantNotEqual {
				require.NotEqual(t, notEqualRg.Hash(), gotHash)
			}
		})
	}
}

func TestResourceGroupToSQL(t *testing.T) {
	tests := []struct {
		name    string
		rg      ResourceGroup
		wantSQL string
	}{
		{
			name: "FullResourceGroup",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
				Name:       "test-pod",
				Labels: map[string]string{
					"app": "myapp",
					"env": "dev",
				},
				Annotations: map[string]string{
					"note": "test",
				},
			},
			wantSQL: "SELECT * from resources WHERE cluster='test-cluster' AND apiVersion='v1' AND kind='Pod' AND namespace='default' AND name='test-pod' AND `annotations.note`='test' AND `labels.app`='myapp' AND `labels.env`='dev'",
		},
		{
			name: "ResourceGroupWithoutLabelsOrAnnotations",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
				Name:       "test-pod",
			},
			wantSQL: "SELECT * from resources WHERE cluster='test-cluster' AND apiVersion='v1' AND kind='Pod' AND namespace='default' AND name='test-pod'",
		},
		{
			name:    "EmptyResourceGroup",
			rg:      ResourceGroup{},
			wantSQL: "SELECT * from resources",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL := tt.rg.ToSQL()
			require.Equal(t, tt.wantSQL, gotSQL)
		})
	}
}

func TestResourceGroupGetType(t *testing.T) {
	tests := []struct {
		name    string
		rg      ResourceGroup
		want    ResourceGroupType
		success bool
	}{
		{
			name: "Cluster",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "",
				Kind:       "",
				Namespace:  "",
				Name:       "",
			},
			want:    Cluster,
			success: true,
		},
		{
			name: "GVK",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "",
				Name:       "",
			},
			want:    GVK,
			success: true,
		},
		{
			name: "Namespace",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "",
				Kind:       "",
				Namespace:  "namespace-1",
				Name:       "",
			},
			want:    Namespace,
			success: true,
		},
		{
			name: "ClusterGVKNamespace",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "v1",
				Kind:       "Namespace",
				Namespace:  "namespace-1",
			},
			want:    ClusterGVKNamespace,
			success: true,
		},
		{
			name: "NamespaceScopedResource",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
				Name:       "pod-1",
			},
			want:    Resource,
			success: true,
		},
		{
			name: "ClusterScopedResource",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "v1",
				Kind:       "Node",
				Namespace:  "",
				Name:       "node-1",
			},
			want:    NonNamespacedResource,
			success: true,
		},
		{
			name: "CustomResourceGroup",
			rg: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Labels: map[string]string{
					"app": "myapp",
				},
				Annotations: map[string]string{
					"note": "test",
				},
			},
			want:    Custom,
			success: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, success := tt.rg.GetType()
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.success, success)
		})
	}
}

func TestNewResourceGroupFromQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    ResourceGroup
		wantErr bool
	}{
		{
			name:  "ValidQuery",
			query: "/api?cluster=test-cluster&apiVersion=v1&kind=Pod&namespace=default&name=test-pod&labels=app.kubernetes.io%2Fname=myapp,env=dev&annotations=note=test",
			want: ResourceGroup{
				Cluster:    "test-cluster",
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
				Name:       "test-pod",
				Labels: map[string]string{
					"app.kubernetes.io/name": "myapp",
					"env":                    "dev",
				},
				Annotations: map[string]string{
					"note": "test",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, tt.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			got, err := NewResourceGroupFromQuery(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResourceGroupFromQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceGroupFromQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}
