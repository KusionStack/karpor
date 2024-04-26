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
	"sync"
	"time"
)

// Cache manages the caching of items based on keys with
// expiration time for cached items.
type Cache[K comparable, V any] struct {
	cache      map[K]*CacheItem[V]
	mu         sync.RWMutex
	expiration time.Duration
}

// CacheItem represents an item stored in the cache along with its expiration
// time.
type CacheItem[V any] struct {
	Data       V
	ExpiryTime time.Time
}

// NewCache creates a new Cache instance with a specified expiration time.
func NewCache[K comparable, V any](expiration time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		cache:      make(map[K]*CacheItem[V]),
		expiration: expiration,
	}
}

// Get retrieves an item from the cache based on the provided key. It returns
// the data and a boolean indicating if the data exists and hasn't expired.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exist := c.cache[key]
	if !exist {
		return zeroValue[V](), false
	}

	if time.Now().After(item.ExpiryTime) {
		delete(c.cache, key)
		return zeroValue[V](), false
	}

	return item.Data, true
}

// Set adds or updates an item in the cache with the provided key and data.
func (c *Cache[K, V]) Set(key K, data V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &CacheItem[V]{
		Data:       data,
		ExpiryTime: time.Now().Add(c.expiration),
	}
}

// zeroValue returns the zero value of type V.
func zeroValue[V any]() V {
	var zero V
	return zero
}
