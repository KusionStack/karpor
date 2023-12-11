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

package cache

import (
	"fmt"
	"strconv"

	hashstructure "github.com/mitchellh/hashstructure/v2"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/cache"
)

type ResourceCache struct {
	cacheStorage cache.ThreadSafeStore
	keyFunc      cache.KeyFunc
}

type ResourcecHash struct {
	ResourceVersion string
	Hash            string
}

func NewResourceCache() *ResourceCache {
	return &ResourceCache{
		keyFunc:      cache.DeletionHandlingMetaNamespaceKeyFunc,
		cacheStorage: cache.NewThreadSafeStore(cache.Indexers{}, cache.Indices{}),
	}
}

func (c *ResourceCache) Add(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return cache.KeyError{Obj: obj, Err: err}
	}

	if item, err := c.newCacheItem(obj); err != nil {
		return err
	} else {
		c.cacheStorage.Add(key, item)
		return nil
	}
}

func (c *ResourceCache) Update(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return cache.KeyError{Obj: obj, Err: err}
	}

	if item, err := c.newCacheItem(obj); err != nil {
		return err
	} else {
		c.cacheStorage.Update(key, item)
		return nil
	}
}

func (c *ResourceCache) Delete(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return cache.KeyError{Obj: obj, Err: err}
	}

	c.cacheStorage.Delete(key)
	return nil
}

func (c *ResourceCache) List() []interface{} {
	return c.cacheStorage.List()
}

func (c *ResourceCache) ListKeys() []string {
	return c.cacheStorage.ListKeys()
}

func (c *ResourceCache) Get(obj interface{}) (item interface{}, exists bool, err error) {
	key, err := c.keyFunc(obj)
	if err != nil {
		return nil, false, cache.KeyError{Obj: obj, Err: err}
	}
	item, exists = c.cacheStorage.Get(key)
	return
}

func (c *ResourceCache) GetByKey(key string) (item interface{}, exists bool, err error) {
	item, exists = c.cacheStorage.Get(key)
	return
}

func (c *ResourceCache) Replace(objs []interface{}, resourceVersion string) error {
	items := make(map[string]interface{})
	for _, obj := range objs {
		key, err := c.keyFunc(obj)
		if err != nil {
			return cache.KeyError{Obj: obj, Err: err}
		}

		if item, err := c.newCacheItem(obj); err != nil {
			return err
		} else {
			items[key] = item
		}
	}
	c.cacheStorage.Replace(items, resourceVersion)
	return nil
}

func (c *ResourceCache) Resync() error {
	// Nothing to do
	return nil
}

func (c *ResourceCache) newCacheItem(obj interface{}) (*ResourcecHash, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	// exclude resource version from the calculation
	resourceVersion := accessor.GetResourceVersion()
	defer accessor.SetResourceVersion(resourceVersion)
	accessor.SetResourceVersion("")

	hash, err := c.hash(obj)
	if err != nil {
		return nil, err
	}

	return &ResourcecHash{
		ResourceVersion: resourceVersion,
		Hash:            hash,
	}, nil
}

func (c *ResourceCache) hash(obj interface{}) (string, error) {
	hash, err := hashstructure.Hash(obj, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(hash, 10), nil
}

func (c *ResourceCache) IsNewer(obj interface{}) (bool, error) {
	cachedItem, exist, err := c.Get(obj)
	if err != nil {
		return false, err
	}
	if !exist {
		return true, nil
	}

	// first, compare the resource version
	rh := cachedItem.(*ResourcecHash)
	compare, err := CompareResourverVersion(obj, rh.ResourceVersion)
	if err != nil {
		return false, err
	}
	if compare <= 0 {
		return false, nil
	}

	// if resource version is newer, compare the hash
	newItem, err := c.newCacheItem(obj)
	if err != nil {
		return false, err
	}
	return newItem.Hash != rh.Hash, nil
}

func CompareResourverVersion(obj interface{}, resourceVersion string) (int, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return -1, err
	}

	rv1, err := parseResourceVersion(accessor.GetResourceVersion())
	if err != nil {
		return -1, err
	}
	rv2, _ := parseResourceVersion(resourceVersion)
	if rv1 == rv2 {
		return 0, nil
	}
	if rv1 < rv2 {
		return -1, nil
	}
	return 1, nil
}

func parseResourceVersion(resourceVersion string) (uint64, error) {
	if resourceVersion == "" || resourceVersion == "0" {
		return 0, nil
	}
	version, err := strconv.ParseUint(resourceVersion, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid resource version %s", resourceVersion)
	}
	return version, nil
}
