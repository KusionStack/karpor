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

package audit

import (
	"io"
	"net/http"

	"github.com/go-chi/render"
)

// decode detects the correct decoder for use on an HTTP request and
// marshals into a given interface.
func decode(r *http.Request, payload *Payload) error {
	// Check if the content type is plain text, read it as such.
	if render.GetRequestContentType(r) == render.ContentTypePlainText {
		// Read the request body.
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close() // Ensure the body is closed after reading.
		if err != nil {
			return err
		}
		// Set the read content as the manifest payload.
		payload.Manifest = string(body)
	} else {
		// For non-plain text, decode the JSON body into the payload.
		if err := render.DecodeJSON(r.Body, payload); err != nil {
			return err
		}
	}

	return nil
}
