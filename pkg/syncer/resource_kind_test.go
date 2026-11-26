// Copyright The Karpor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syncer

import (
	"testing"

	searchv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	clienttesting "k8s.io/client-go/testing"
)

func Test_singleClusterSyncManager_resourceKind(t *testing.T) {
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &clienttesting.Fake{}}
	discoveryClient.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "networking.k8s.io/v1",
			APIResources: []metav1.APIResource{
				{Name: "ingresses", Kind: "Ingress"},
				{Name: "ingresses/status", Kind: "Ingress"},
				{Name: "networkpolicies", Kind: "NetworkPolicy"},
			},
		},
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{Name: "endpoints", Kind: "Endpoints"},
				{Name: "pods", Kind: "Pod"},
			},
		},
	}
	s := &singleClusterSyncManager{discoveryClient: discoveryClient, logger: logr.Discard()}

	tests := []struct {
		apiVersion string
		resource   string
		want       string
	}{
		{"networking.k8s.io/v1", "ingresses", "Ingress"},
		{"networking.k8s.io/v1", "networkpolicies", "NetworkPolicy"},
		{"v1", "endpoints", "Endpoints"},
		{"v1", "pods", "Pod"},
		{"v1", "unknowns", ""},
		{"example.com/v1", "widgets", ""},
	}
	for _, tt := range tests {
		t.Run(tt.apiVersion+"/"+tt.resource, func(t *testing.T) {
			rsr := &searchv1beta1.ResourceSyncRule{APIVersion: tt.apiVersion, Resource: tt.resource}
			require.Equal(t, tt.want, s.resourceKind(rsr))
		})
	}
}
