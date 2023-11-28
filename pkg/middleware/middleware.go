package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"k8s.io/klog/v2"
)

type contextKey struct {
	name string
}

var APILoggerKey = &contextKey{"logger"}

func APILogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if reqIDs, ok := r.Header[middleware.RequestIDHeader]; ok && len(reqIDs) > 0 && len(reqIDs[0]) > 0 {
			requestID := reqIDs[0]
			logger := klog.FromContext(r.Context()).WithValues("requestID", requestID)
			ctx = context.WithValue(r.Context(), APILoggerKey, logger)
		}

		// continue serving request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DefaultLogger(next http.Handler) http.Handler {
	return middleware.Logger(next)
}
