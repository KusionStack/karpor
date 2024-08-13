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

	"github.com/KusionStack/karpor/pkg/util/ctxutil"
	"github.com/go-chi/render"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// HandleResult is a handler function that writes the response to the HTTP response writer based on the provided error and data.
func HandleResult(w http.ResponseWriter, r *http.Request, ctx context.Context, err error, data any) {
	if err != nil {
		render.Render(w, r, failureResponse(ctx, err))
		return
	}
	render.JSON(w, r, successResponse(ctx, data))
}

// RemoveUnstructuredManagedFields remove managedFields information within a Unstructured
func RemoveUnstructuredManagedFields(
	ctx context.Context,
	yaml *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
	log := ctxutil.GetLogger(ctx)

	// Inform that the unmarshaling process has started.
	log.Info("Sanitizing unstructured cluster...")
	sanitized := yaml
	if _, ok := sanitized.Object["metadata"].(map[string]interface{})["managedFields"]; ok {
		sanitized.Object["metadata"].(map[string]interface{})["managedFields"] = "[redacted]"
	}
	return sanitized, nil
}
