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

// Copied from github.com/elliotxx/healthcheck

package health

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	stringBuilderPool = sync.Pool{
		New: func() any {
			return &strings.Builder{}
		},
	}
	ErrHealthCheckNamesConflict = errors.New("health check names conflict")
	ErrHealthCheckFailed        = errors.New("health check failed")
)

// Check represents a health check that can be run to check the status
// of a service.
type Check interface {
	Pass(context.Context) bool
	Name() string
}

// HandlerConfig represents the configuration for the health check
// handler.
type HandlerConfig struct {
	Verbose  bool
	Excludes []string
	Checks   []Check
	FailureNotification
}

// NewCheckStatuses creates a new CheckStatuses instance with the
// specified capacity.
func NewCheckStatuses(n int) *CheckStatuses {
	return &CheckStatuses{
		m: make(map[string]bool, n),
	}
}

// FailureNotification represents the configuration for failure
// notifications.
type FailureNotification struct {
	Chan      chan error
	Threshold uint32
}

type CheckStatuses struct {
	m map[string]bool
	sync.RWMutex
}

// Get returns the value and existence status for the specified key.
func (cs *CheckStatuses) Get(k string) (bool, bool) {
	cs.RLock()
	defer cs.RUnlock()
	v, existed := cs.m[k]
	return v, existed
}

// Set sets the value for the specified key.
func (cs *CheckStatuses) Set(k string, v bool) {
	cs.Lock()
	defer cs.Unlock()
	cs.m[k] = v
}

// Delete deletes the specified key from the map.
func (cs *CheckStatuses) Delete(k string) {
	cs.Lock()
	defer cs.Unlock()
	delete(cs.m, k)
}

// Len returns the number of items in the map.
func (cs *CheckStatuses) Len() int {
	cs.RLock()
	defer cs.RUnlock()

	return len(cs.m)
}

// Each calls the specified function for each key/value pair in the
// map.
func (cs *CheckStatuses) Each(f func(k string, v bool)) {
	cs.RLock()
	defer cs.RUnlock()

	for k, v := range cs.m {
		f(k, v)
	}
}

// String returns a string representation of the check statuses.
// If verbose is true, the output includes pass/fail status for each
// check.
func (cs *CheckStatuses) String(verbose bool) string {
	passNames := make([]string, 0, cs.Len())
	failedNames := make([]string, 0, cs.Len())
	allPass := true
	cs.Each(func(name string, pass bool) {
		if pass {
			passNames = append(passNames, name)
		} else {
			failedNames = append(failedNames, name)
			allPass = false
		}
	})

	sort.Strings(passNames)
	sort.Strings(failedNames)

	if verbose {
		b := stringBuilderPool.Get().(*strings.Builder)
		defer stringBuilderPool.Put(b)
		defer b.Reset()

		for _, name := range passNames {
			b.WriteString("[+] " + name + " ok\n")
		}

		for _, name := range failedNames {
			b.WriteString("[-] " + name + " fail\n")
		}

		if allPass {
			b.WriteString("health check passed")
		} else {
			b.WriteString("health check failed")
		}

		return b.String()
	}

	if allPass {
		return "OK"
	}
	return "Fail"
}

// pingCheck is a simple health check that always returns true.
type pingCheck struct{}

// NewPingCheck creates a new ping health check.
func NewPingCheck() Check {
	return &pingCheck{}
}

// Pass always returns true for the ping health check.
func (c *pingCheck) Pass(context.Context) bool {
	return true
}

// Name returns the name of the ping health check.
func (c *pingCheck) Name() string {
	return "Ping"
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

		ctx := r.Context()

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
				pass := captureCheck.Pass(ctx)
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
