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

package utils

import "testing"

func TestResourceKind(t *testing.T) {
	tests := []struct {
		resource string
		kind     string
		want     string
	}{
		{"ingresses", "Ingress", "Ingress"},
		{"storageclasses", "StorageClass", "StorageClass"},
		{"networkpolicies", "NetworkPolicy", "NetworkPolicy"},
		{"endpoints", "Endpoints", "Endpoints"},
		{"pods", "", "pod"},
		{"", "", ""},
	}
	for _, tt := range tests {
		if got := ResourceKind(tt.resource, tt.kind); got != tt.want {
			t.Errorf("ResourceKind(%q, %q) = %q, want %q", tt.resource, tt.kind, got, tt.want)
		}
	}
}
