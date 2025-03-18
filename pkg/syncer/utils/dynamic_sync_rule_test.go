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

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	searchv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
)

func TestGetSyncGVK(t *testing.T) {
	// Test cases
	testCases := []struct {
		name             string
		gvk              schema.GroupVersionKind
		resourceSyncRule searchv1beta1.ResourceSyncRule
		exist            bool
	}{
		{
			name:             "Success - ListKeys",
			gvk:              schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
			resourceSyncRule: ZeroVal,
			exist:            false,
		},
	}
	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceSyncRule, exist := GetSyncGVK(tc.gvk)
			require.Equal(t, tc.exist, exist)
			require.Equal(t, tc.resourceSyncRule, resourceSyncRule)
		})
	}
}
