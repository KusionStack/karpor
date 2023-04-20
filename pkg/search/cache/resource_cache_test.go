package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNewer(t *testing.T) {
	assert := assert.New(t)

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
	assert.NoError(err)
	assert.True(newer)
}

func TestIsNewer_MapOrderShouldNotMatter(t *testing.T) {
	assert := assert.New(t)

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
	assert.NoError(err)
	assert.False(newer, "map with different order should be taken as same")
}

func TestIsNewer_ResourceVersionIsExcluded(t *testing.T) {
	assert := assert.New(t)

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
	assert.NoError(err)
	assert.False(newer, "resource version should be excluded from hash caculation")
}
