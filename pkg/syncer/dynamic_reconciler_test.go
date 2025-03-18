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

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
)

func TestDynamicReconciler_Reconcile(t *testing.T) {
	tests := []struct {
		name    string
		gvk     schema.GroupVersionKind
		pod     *v1.Pod
		req     reconcile.Request
		wantErr bool
	}{
		{
			"test no error with pods",
			schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
			&v1.Pod{},
			reconcile.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDynamicReconciler(context.Background(), "cluster-name", tt.gvk, &elasticsearch.Storage{})
			r.client = fake.NewClientBuilder().WithRuntimeObjects(tt.pod).WithScheme(scheme.Scheme).Build()
			r.scheme = scheme.Scheme
			utils.SetSyncGVK(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}, v1beta1.ResourceSyncRule{
				Namespace: "karpor",
				Transform: &v1beta1.TransformRuleSpec{
					Type: "patch",
				},
				Trim: &v1beta1.TrimRuleSpec{},
			})

			m := mockey.Mock((*elasticsearch.Storage).SaveResource).Return(nil).Build()
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
