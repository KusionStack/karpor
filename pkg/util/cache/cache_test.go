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
	"sync"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	expiration := 100 * time.Millisecond
	cache := NewCache[int, string](expiration)

	key := 42
	value := "test value"

	cache.Set(key, value)

	// Check if the value is retrieved correctly
	retrievedValue, exists := cache.Get(key)
	if !exists {
		t.Errorf("Expected value '%s' to exist in cache, but it doesn't.", value)
	}
	if retrievedValue != value {
		t.Errorf("Expected value '%s', got '%s'", value, retrievedValue)
	}

	// Wait for the value to expire
	time.Sleep(expiration + 50*time.Millisecond)

	// Check if the value is expired
	_, exists = cache.Get(key)
	if exists {
		t.Error("Expected value to be expired, but it still exists in cache.")
	}
}

func TestCache_SetAndGet_Concurrent(t *testing.T) {
	expiration := 100 * time.Millisecond
	cache := NewCache[int, string](expiration)

	key := 42
	value := "test value"

	// Concurrently set and get the value
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		cache.Set(key, value)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		retrievedValue, exists := cache.Get(key)
		if !exists || retrievedValue != value {
			t.Errorf("Concurrent Set/Get: Expected value '%s', got '%s'", value, retrievedValue)
		}
	}()

	wg.Wait()
}

func TestCache_ExpiredKeyIsDeleted(t *testing.T) {
	expiration := 100 * time.Millisecond
	cache := NewCache[int, string](expiration)

	key := 42
	value := "test value"

	cache.Set(key, value)

	// Wait for the value to expire
	time.Sleep(expiration + 50*time.Millisecond)

	// Access the expired key
	_, exists := cache.Get(key)
	if exists {
		t.Error("Expected expired key to be automatically deleted from the cache.")
	}
}
