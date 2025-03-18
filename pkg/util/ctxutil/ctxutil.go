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

package ctxutil

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"

	"github.com/KusionStack/karpor/pkg/core/middleware"
)

// GetLogger returns the logger from the given context.
//
// Example:
//
//	logger := ctxutil.GetLogger(ctx)
func GetLogger(ctx context.Context) logr.Logger {
	if logger, ok := ctx.Value(middleware.APILoggerKey).(logr.Logger); ok {
		return logger
	}

	return klog.NewKlogr()
}
