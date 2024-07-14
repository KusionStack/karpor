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
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/syncer"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"k8s.io/kubernetes/pkg/kubelet/server"
	"net/http"
	"strings"
	"sync"
)

// Register registers the livez and readyz handlers to the specified
// router.
func Register(r *chi.Mux, serv server.Server, sync syncer.ResourceSyncer, sg storage.Storage) {
	r.Get("/livez", NewLivezHandler())
	r.Get("/readyz", NewReadyzHandler(serv, sync, sg))
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
func NewReadyzHandler(serv server.Server, sync syncer.ResourceSyncer, sg storage.Storage) http.HandlerFunc {
	conf := HandlerConfig{
		Verbose: true,
		// checkList is a list of healthcheck to run.
		Checks: []Check{
			NewPingCheck(),
			//NewServerCheck(serv),
			//NewSyncerCheck(sync),
			//NewStorageCheck(sg),
		},
		FailureNotification: FailureNotification{Threshold: 1},
	}

	return NewHandler(conf)
}

func NewHandler(conf HandlerConfig) http.HandlerFunc {
	var (
		mu            sync.Mutex
		failureInARow uint32
	)

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			eg         errgroup.Group
			httpStatus = http.StatusOK
		)

		// Process the request parameters.
		verbose := conf.Verbose
		excludes := conf.Excludes
		if r.URL.Query().Has("verbose") {
			verbose = true
		}
		if excludesStr := r.URL.Query().Get("excludes"); excludesStr != "" {
			excludes = strings.Split(excludesStr, ",")
		}

		// Create a new check statuses instance.
		statuses := NewCheckStatuses(len(conf.Checks))

		// Iterate over the list of health checks and execute them
		// concurrently.
		for _, check := range conf.Checks {
			// Capture the check variable to avoid race conditions.
			captureCheck := check

			eg.Go(func() error {
				// Get the name of the check and check if it already
				// exists in the statuses list.
				name := captureCheck.Name()

				if len(excludes) > 0 {
					for _, excludedName := range excludes {
						if excludedName == name {
							return nil
						}
					}
				}

				if _, ok := statuses.Get(name); ok {
					return ErrHealthCheckNamesConflict
				}

				// Execute the check and update the status list.
				pass := captureCheck.Pass()
				statuses.Set(name, pass)

				// If the check fails, return a failure error.
				if !pass {
					return ErrHealthCheckFailed
				}
				return nil
			})
		}

		// Wait for all the checks to complete.
		mu.Lock()
		if err := eg.Wait(); err != nil {
			// If any of the checks fail, set the HTTP status code to service
			// unavailable.
			httpStatus = http.StatusServiceUnavailable
			failureInARow++

			// Send a notification if the failure threshold is reached.
			if failureInARow >= conf.FailureNotification.Threshold &&
				conf.FailureNotification.Chan != nil {
				conf.FailureNotification.Chan <- err
			}
		} else if failureInARow != 0 && conf.FailureNotification.Chan != nil {
			// Reset the failure counter if all checks pass.
			failureInARow = 0
			conf.FailureNotification.Chan <- nil
		}
		mu.Unlock()

		// Return the status response as a string.
		w.WriteHeader(httpStatus)
		w.Write([]byte(statuses.String(verbose)))
	}
}
