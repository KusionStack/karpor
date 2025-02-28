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
	"sync/atomic"
	"time"
)

// RateLimiter limits the max retry times for each event
type RateLimiter struct {
	maxRetries int
	items      sync.Map
}

// NewRateLimiter create a new RateLimiter instance
func NewRateLimiter(maxRetries int) *RateLimiter {
	return &RateLimiter{
		maxRetries: maxRetries,
	}
}

// When return time during of next reconcile
func (r *RateLimiter) When(item interface{}) time.Duration {
	counter, _ := r.items.LoadOrStore(item, new(int32))
	retries := atomic.LoadInt32(counter.(*int32))

	if retries >= int32(r.maxRetries) {
		return time.Duration(1<<63 - 1) // math.MaxInt64
	}

	atomic.AddInt32(counter.(*int32), 1)

	return 10 * time.Second
}

// Forget forgets event
func (r *RateLimiter) Forget(item interface{}) {
	r.items.Delete(item)
}

// NumRequeues return retry times of special event
func (r *RateLimiter) NumRequeues(item interface{}) int {
	counter, ok := r.items.Load(item)
	if !ok {
		return 0
	}
	return int(atomic.LoadInt32(counter.(*int32)))
}
