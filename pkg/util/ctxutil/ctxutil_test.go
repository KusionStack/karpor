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
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	"k8s.io/klog/v2/klogr"

	"github.com/KusionStack/karpor/pkg/core/middleware"
)

func TestGetLogger(t *testing.T) {
	mockLogger := klogr.New().WithName("mock")
	tests := []struct {
		name           string
		ctx            context.Context
		expectedLogger logr.Logger // Expected logger type
	}{
		{
			name: "Logger in context",
			ctx:  context.WithValue(context.Background(), middleware.APILoggerKey, mockLogger),
			// Expect the logger type to be klogr.Logger.
			expectedLogger: mockLogger,
		},
		{
			name: "Logger not in context",
			ctx:  context.Background(),
			// Expect the logger type to be klogr.Logger as default.
			expectedLogger: klogr.New(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := GetLogger(tt.ctx)
			// Check if the returned logger type matches the expected logger
			// type.
			require.IsType(t, tt.expectedLogger, logger)
		})
	}
}
