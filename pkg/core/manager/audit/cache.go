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

package audit

import (
	"sync"
	"time"

	"github.com/KusionStack/karbour/pkg/core"
)

// Cache manages the caching of AuditData based on core.Locator keys with
// expiration time for cached items.
type Cache struct {
	cache      map[core.Locator]*CacheItem
	mu         sync.RWMutex
	expiration time.Duration
}

// CacheItem represents an item stored in the cache along with its expiration time.
type CacheItem struct {
	Data       *AuditData
	ExpiryTime time.Time
}

// NewCache creates a new Cache instance with a specified expiration time.
func NewCache(expiration time.Duration) *Cache {
	return &Cache{
		cache:      make(map[core.Locator]*CacheItem),
		expiration: expiration,
	}
}

// Get retrieves an item from the cache based on the provided locator. It returns
// the AuditData and a boolean indicating if the data exists and hasn't expired.
func (c *Cache) Get(locator core.Locator) (*AuditData, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exist := c.cache[locator]
	if !exist {
		return nil, false
	}

	if time.Now().After(item.ExpiryTime) {
		delete(c.cache, locator)
		return nil, false
	}

	return item.Data, true
}

// Set adds or updates an item in the cache with the provided locator and AuditData.
func (c *Cache) Set(locator core.Locator, data *AuditData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[locator] = &CacheItem{
		Data:       data,
		ExpiryTime: time.Now().Add(c.expiration),
	}
}
