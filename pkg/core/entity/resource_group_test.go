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

package entity

import "testing"

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

			if string(gotHash) != tt.wantHash {
				t.Errorf("Hash() = %v, want %v", gotHash, tt.wantHash)
			}

			for _, equalRg := range tt.wantEqual {
				if gotHash != equalRg.Hash() {
					t.Errorf("Hash() of %v should be equal to %v", equalRg, tt.rg)
				}
			}

			for _, notEqualRg := range tt.wantNotEqual {
				if gotHash == notEqualRg.Hash() {
					t.Errorf("Hash() of %v should not be equal to %v", notEqualRg, tt.rg)
				}
			}
		})
	}
}
