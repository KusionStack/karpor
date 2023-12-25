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

package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"k8s.io/klog/v2"
)

// APILoggerKey is a context key used for associating a logger with a request.
var APILoggerKey = &contextKey{"logger"}

// APILogger is a middleware that injects a logger, configured with a request ID,
// into the request context for use throughout the request's lifecycle.
func APILogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Retrieve the request ID from the context and create a logger with it.
		if requestID := middleware.GetReqID(ctx); len(requestID) > 0 {
			logger := klog.FromContext(ctx).WithValues("requestID", requestID)
			ctx = context.WithValue(ctx, APILoggerKey, logger)
		}

		// Continue serving the request with the new context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DefaultLogger is a middleware that provides basic request logging using chi's
// built-in Logger middleware.
func DefaultLogger(next http.Handler) http.Handler {
	return middleware.Logger(next)
}
