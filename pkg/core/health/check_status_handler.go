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
	"github.com/pkg/errors"
	"sort"
	"strings"
	"sync"
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
	Pass() bool
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
	Threshold uint32
	Chan      chan error
}

type CheckStatuses struct {
	sync.RWMutex
	m map[string]bool
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
func (c *pingCheck) Pass() bool {
	return true
}

// Name returns the name of the ping health check.
func (c *pingCheck) Name() string {
	return "Ping"
}
