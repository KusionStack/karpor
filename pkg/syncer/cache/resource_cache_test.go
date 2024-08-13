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

package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNewer(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
			"f1":                       "foo",
		})

	cache.Add(u1)

	u2 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "2",
			"f1":                       "foo",
			"f2":                       "bar",
		})
	newer, err := cache.IsNewer(u2)
	require.NoError(t, err)
	require.True(t, newer)
}

func TestIsNewer_MapOrderShouldNotMatter(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"f1": "foo",
			"f2": "bar",
		})

	cache.Add(u1)

	u2 := makeUnstructured("default", "test",
		map[string]interface{}{
			"f2": "bar",
			"f1": "foo",
		})
	newer, err := cache.IsNewer(u2)
	require.NoError(t, err)
	require.False(t, newer, "map with different order should be taken as same")
}

func TestIsNewer_ResourceVersionIsExcluded(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
		})

	cache.Add(u1)

	u2 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "2",
		})

	newer, err := cache.IsNewer(u2)
	require.NoError(t, err)
	require.False(t, newer, "resource version should be excluded from hash caculation")
}

func TestResourceCache_Get(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
		})

	cache.Add(u1)

	u2 := makeUnstructured("default", "test", map[string]interface{}{
		"metadata.resourceVersion": "1",
	})

	item, exists, err := cache.Get(u2)
	require.NoError(t, err)
	require.True(t, exists)
	require.NotNil(t, item)
}

func TestResourceCache_GetByKey(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
		})

	cache.Add(u1)

	item, exists, err := cache.GetByKey("default/test")
	require.NoError(t, err)
	require.True(t, exists)
	require.NotNil(t, item)
}

func TestResourceCache_List(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
		})

	cache.Add(u1)

	list := cache.List()
	require.Len(t, list, 1)
}

func TestResourceCache_ListKeys(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
		})

	cache.Add(u1)

	list := cache.ListKeys()
	require.Len(t, list, 1)
	require.Equal(t, "default/test", list[0])
}

func TestResourceCache_Update(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
		})

	cache.Add(u1)

	u2 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "2",
		})

	err := cache.Update(u2)
	require.NoError(t, err)
}

func TestResourceCache_Delete(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
		})

	cache.Add(u1)

	u2 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "2",
		})
	err := cache.Delete(u2)
	require.NoError(t, err)
}

func TestResourceCache_Replace(t *testing.T) {
	cache := NewResourceCache()
	u1 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "1",
		})

	cache.Add(u1)

	u2 := makeUnstructured("default", "test",
		map[string]interface{}{
			"metadata.resourceVersion": "2",
		})
	err := cache.Replace([]any{u2}, "v1")
	require.NoError(t, err)
}
