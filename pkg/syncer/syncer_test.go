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
	"strings"
	"testing"
	"time"

	"kusionstack.io/karpor/pkg/infra/search/storage/elasticsearch"
	"kusionstack.io/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/workqueue"
)

func Test_genUnObj(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name string                     // Name of the test case
		sr   v1beta1.ResourceSyncRule   // Input parameter to genUnObj
		key  string                     // Input parameter to genUnObj
		want *unstructured.Unstructured // Expected result
	}{
		{
			name: "Single key without namespace",
			sr: v1beta1.ResourceSyncRule{
				APIVersion: "v1",
				Resource:   "pods",
			},
			key: "mypod",
			want: func() *unstructured.Unstructured {
				obj := &unstructured.Unstructured{}
				obj.SetAPIVersion("v1")
				obj.SetKind("pod")
				obj.SetName("mypod")
				return obj
			}(),
		},
		{
			name: "Key with namespace",
			sr: v1beta1.ResourceSyncRule{
				APIVersion: "v1",
				Resource:   "services",
			},
			key: "mynamespace/myservice",
			want: func() *unstructured.Unstructured {
				obj := &unstructured.Unstructured{}
				obj.SetAPIVersion("v1")
				obj.SetKind("service")
				obj.SetNamespace("mynamespace")
				obj.SetName("myservice")
				return obj
			}(),
		},
	}
	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			un := genUnObj(tt.sr, tt.key)
			require.Equal(t, tt.want, un, "The generated unstructured object does not match the expected value.")
		})
	}
}

func TestResourceSyncer_sync(t *testing.T) {
	tests := []struct {
		name    string
		exist   bool
		item    any
		wantErr bool
	}{
		{
			"test exist",
			true,
			&unstructured.Unstructured{},
			false,
		},
		{
			"test not exist",
			false,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewResourceSyncer("cluster1", nil, v1beta1.ResourceSyncRule{APIVersion: "v1", Resource: "services"}, &elasticsearch.Storage{})
			m1 := mockey.Mock((*informerSource).GetByKey).Return(tt.item, tt.exist, nil).Build()
			m2 := mockey.Mock((*elasticsearch.Storage).SaveResource).Return(nil).Build()
			m3 := mockey.Mock((*elasticsearch.Storage).DeleteResource).Return(nil).Build()
			defer m3.UnPatch()
			defer m2.UnPatch()
			defer m1.UnPatch()
			err := s.sync(context.TODO(), "test")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestResourceSyncer_Run(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"test no error",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewResourceSyncer("cluster1", nil, v1beta1.ResourceSyncRule{APIVersion: "v1", Resource: "services"}, &elasticsearch.Storage{})
			m := mockey.Mock((*informerSource).HasSynced).Return(true).Build()
			defer m.UnPatch()
			ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
			defer cancel()
			err := s.Run(ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestResourceSyncer_OnAdd(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test normally",
			args: args{
				name:      "ns1",
				namespace: "name1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &fakeQueue{}
			s := &ResourceSyncer{
				queue: q,
			}
			pod := &corev1.Pod{}
			pod.Namespace = tt.args.namespace
			pod.Name = tt.args.name
			s.OnAdd(pod)
			require.Equal(t, strings.Join([]string{tt.args.namespace, tt.args.name}, "/"), q.queue[0])
		})
	}
}

func TestResourceSyncer_OnUpdate(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test normally",
			args: args{
				name:      "ns1",
				namespace: "name1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &fakeQueue{}
			s := &ResourceSyncer{
				queue: q,
			}
			pod := &corev1.Pod{}
			pod.Namespace = tt.args.namespace
			pod.Name = tt.args.name
			s.OnUpdate(pod)
			require.Equal(t, strings.Join([]string{tt.args.namespace, tt.args.name}, "/"), q.queue[0])
		})
	}
}

func TestResourceSyncer_OnGeneric(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test normally",
			args: args{
				name:      "ns1",
				namespace: "name1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &fakeQueue{}
			s := &ResourceSyncer{
				queue: q,
			}
			pod := &corev1.Pod{}
			pod.Namespace = tt.args.namespace
			pod.Name = tt.args.name
			s.OnGeneric(pod)
			require.Equal(t, strings.Join([]string{tt.args.namespace, tt.args.name}, "/"), q.queue[0])
		})
	}
}

func TestResourceSyncer_OnDelete(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test normally",
			args: args{
				name:      "ns1",
				namespace: "name1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &fakeQueue{}
			s := &ResourceSyncer{
				queue: q,
			}
			pod := &corev1.Pod{}
			pod.Namespace = tt.args.namespace
			pod.Name = tt.args.name
			s.OnDelete(pod)
			require.Equal(t, strings.Join([]string{tt.args.namespace, tt.args.name}, "/"), q.queue[0])
		})
	}
}

type fakeQueue struct {
	workqueue.RateLimitingInterface
	queue []interface{}
}

func (q *fakeQueue) Add(item interface{}) {
	if q.queue == nil {
		q.queue = []interface{}{}
	}
	q.queue = append(q.queue, item)
}
