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

package handler

import (
	"context"
	"net/http"
	"time"

	appmiddleware "github.com/KusionStack/karpor/pkg/core/middleware"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// SuccessMessage is the default success message for successful responses.
const SuccessMessage = "OK"

// Response creates a standard API response renderer.
func Response(ctx context.Context, data any, err error, statusCode int) render.Renderer {
	resp := &response{}

	// Set the Message and Data fields based on the error parameter.
	if err == nil {
		resp.Message = SuccessMessage
		resp.Data = data
	} else {
		resp.Message = err.Error()
	}

	// Set the Success fields based on the error and statusCode parameters.
	if err == nil || statusCode == http.StatusNotFound {
		resp.Success = true
	} else {
		resp.Success = false
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
	return Response(ctx, nil, err, http.StatusInternalServerError)
}

// SuccessResponse creates a response renderer for a successful request.
func SuccessResponse(ctx context.Context, data any) render.Renderer {
	return Response(ctx, data, nil, http.StatusOK)
}

// NotFoundResponse creates a response renderer for a not found request.
func NotFoundResponse(ctx context.Context, err error) render.Renderer {
	return Response(ctx, nil, err, http.StatusNotFound)
}

// FailureRender renders a failed response and status code and respond to the
// client request.
func FailureRender(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error {
	render.Status(r, http.StatusInternalServerError)
	respRender := FailureResponse(ctx, err)
	return render.Render(w, r, respRender)
}

// SuccessRender renders a success response and status code and respond to the
// client request.
func SuccessRender(ctx context.Context, w http.ResponseWriter, r *http.Request, data any) error {
	render.Status(r, http.StatusOK)
	respRender := SuccessResponse(ctx, data)
	return render.Render(w, r, respRender)
}

// NotFoundRender renders a not found response and status code and respond to the
// client request.
func NotFoundRender(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error {
	render.Status(r, http.StatusNotFound)
	respRender := NotFoundResponse(ctx, err)
	return render.Render(w, r, respRender)
}
