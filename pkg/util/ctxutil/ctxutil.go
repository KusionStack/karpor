package ctxutil

import (
	"context"

	"github.com/KusionStack/karbour/pkg/middleware"
	"k8s.io/klog/v2"
)

// GetLogger returns the logger from the given context.
//
// Example:
//
//	logger := ctxutil.GetLogger(ctx)
func GetLogger(ctx context.Context) klog.Logger {
	if logger, ok := ctx.Value(middleware.APILoggerKey).(klog.Logger); ok {
		return logger
	}

	return klog.NewKlogr()
}
