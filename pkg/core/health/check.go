/*
Copyright The Karpor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package health

import (
	"net/http"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/go-chi/chi/v5"
)

// Register registers the livez and readyz handlers to the specified
// router.
func Register(r *chi.Mux, sg storage.Storage) {
	r.Get("/livez", NewLivezHandler())
	r.Get("/readyz", NewReadyzHandler(sg))
}

// NewLivezHandler creates a new liveness check handler that can be
// used to check if the application is running.
func NewLivezHandler() http.HandlerFunc {
	conf := HandlerConfig{
		Verbose: false,
		// checkList is a list of healthcheck to run.
		Checks: []Check{
			NewPingCheck(),
		},
		FailureNotification: FailureNotification{Threshold: 1},
	}

	return NewHandler(conf)
}

// NewReadyzHandler creates a new readiness check handler that can be
// used to check if the application is ready to serve traffic.
// Currently, this feature only contains a storage check and a ping check.
func NewReadyzHandler(sg storage.Storage) http.HandlerFunc {
	conf := HandlerConfig{
		Verbose: true,
		// checkList is a list of healthcheck to run.
		Checks: []Check{
			NewPingCheck(),
			NewStorageCheck(sg),
		},
		FailureNotification: FailureNotification{Threshold: 1},
	}

	return NewHandler(conf)
}

func NewStorageCheck(sg storage.Storage) Check {
	return NewStorageCheckHandler(sg)
}
