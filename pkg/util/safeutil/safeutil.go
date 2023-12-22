package safeutil

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/elliotxx/safe"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

// LoggerRecoverHandler returns a recover handler by the given logger.
//
// Example:
//
//	func() {
//	  defer safe.HandleCrash(RecoverHandler(ctx, errChan))
//	  ...
//	}
func RecoverHandler(ctx context.Context, errChan chan error) safe.RecoverHandler {
	return func(r any) {
		logger := ctxutil.GetLogger(ctx)
		err := fmt.Errorf("%v", r)
		if errChan != nil {
			errChan <- err
		}
		logger.Error(err, "Recovered from panic", "stackTrace", string(debug.Stack()))
	}
}

// LoggerRecoverHandler returns a recover handler by the given logger.
//
// Example:
//
//	func() {
//	  defer safe.HandleCrash(LoggerRecoverHandler(ctxutil.GetLogger(ctx)))
//	  ...
//	}
func LoggerRecoverHandler(logger logr.Logger) safe.RecoverHandler {
	return func(r any) {
		err := fmt.Errorf("%v", r)
		logger.Error(err, "Recovered from panic", "stackTrace", string(debug.Stack()))
	}
}

// Go starts a recoverable goroutine with a new logger (klog.NewKlogr()).
//
// Example:
//
//	safeutil.Go(func(){...})
func Go(do safe.DoFunc) {
	safe.GoR(do, LoggerRecoverHandler(klog.NewKlogr()))
}

// GoL starts a recoverable goroutine with a given logger.
//
// Example:
//
//	safeutil.GoL(func(){...}, logger)
func GoL(do safe.DoFunc, logger logr.Logger) {
	safe.GoR(do, LoggerRecoverHandler(logger))
}
