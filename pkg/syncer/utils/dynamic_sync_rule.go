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
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"

	searchv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
)

var (
	syncGVK sync.Map
	ZeroVal = searchv1beta1.ResourceSyncRule{}
)

func SetSyncGVK(gvk schema.GroupVersionKind, rule searchv1beta1.ResourceSyncRule) {
	syncGVK.Store(gvk.String(), rule)
}

func GetSyncGVK(gvk schema.GroupVersionKind) (searchv1beta1.ResourceSyncRule, bool) {
	rule, exist := syncGVK.Load(gvk.String())
	if !exist {
		return ZeroVal, false
	}
	return rule.(searchv1beta1.ResourceSyncRule), exist
}
