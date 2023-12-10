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

package handler

import (
	"context"
	"time"

	appmiddleware "github.com/KusionStack/karbour/pkg/middleware"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// SuccessMessage is the default success message for successful responses.
const SuccessMessage = "OK"

// Response creates a standard API response renderer.
func Response(ctx context.Context, data any, err error) render.Renderer {
	resp := &response{}

	// Set the Success and Message fields based on the error parameter.
	if err == nil {
		resp.Success = true
		resp.Message = SuccessMessage
		resp.Data = data
	} else {
		resp.Success = false
		resp.Message = err.Error()
	}

	// Include the request trace ID if available.
	if requestID := middleware.GetReqID(ctx); len(requestID) > 0 {
		resp.TraceID = requestID
	}

	// Calculate and include timing details if a start time is set.
	if startTime := appmiddleware.GetStartTime(ctx); !startTime.IsZero() {
		endTime := time.Now()
		resp.StartTime = &startTime
		resp.EndTime = &endTime
		resp.CostTime = Duration(endTime.Sub(startTime))
	}

	return resp
}

// FailureResponse creates a response renderer for a failed request.
func FailureResponse(ctx context.Context, err error) render.Renderer {
	return Response(ctx, nil, err)
}

// SuccessResponse creates a response renderer for a successful request.
func SuccessResponse(ctx context.Context, data any) render.Renderer {
	return Response(ctx, data, nil)
}
