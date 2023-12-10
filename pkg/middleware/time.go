package middleware

import (
	"context"
	"net/http"
	"time"
)

var StartTimeKey = &contextKey{"startTime"}

func Time(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if GetStartTime(ctx).IsZero() {
			ctx = context.WithValue(ctx, StartTimeKey, time.Now())
		}

		// continue serving request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetStartTime returns the start time from the given context if one is present.
// Returns the zero time if the start time cannot be found.
func GetStartTime(ctx context.Context) time.Time {
	if ctx == nil {
		return time.Time{}
	}
	if startTime, ok := ctx.Value(StartTimeKey).(time.Time); ok {
		return startTime
	}
	return time.Time{}
}
