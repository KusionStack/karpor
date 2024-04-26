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

package middleware

import (
	"context"
	"net/http"
	"time"
)

// StartTimeKey is a context key used for storing the start time of a request.
var StartTimeKey = &contextKey{"startTime"}

// Timing is a middleware that captures the current time at the start of a request
// and stores it in the request context. This start time can be used to measure
// request processing duration.
func Timing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Set the start time in the context if it hasn't already been set.
		if GetStartTime(ctx).IsZero() {
			ctx = context.WithValue(ctx, StartTimeKey, time.Now())
		}

		// Continue serving the request with the updated context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetStartTime returns the start time from the given context if one is present.
// If the start time is not present or the context is nil, returns the zero time.
func GetStartTime(ctx context.Context) time.Time {
	if ctx == nil {
		// Return zero time if the context is nil.
		return time.Time{}
	}
	if startTime, ok := ctx.Value(StartTimeKey).(time.Time); ok {
		// Return the start time if it's present in the context.
		return startTime
	}
	// Return zero time if the start time is not found in the context.
	return time.Time{}
}
