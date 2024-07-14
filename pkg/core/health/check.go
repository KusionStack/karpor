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
	"k8s.io/kubernetes/pkg/kubelet/server"
	"net/http"
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
			NewServerCheck(serv, sg),
			NewSyncerCheck(sync),
			NewStorageCheck(sg),
		},
		FailureNotification: FailureNotification{Threshold: 1},
	}

	return NewHandler(conf)
}

func NewServerCheck(serv server.Server, storage storage.Storage) Check {
	return nil
}

func NewSyncerCheck(sync syncer.ResourceSyncer) Check {
	return nil
}

func NewStorageCheck(sg storage.Storage) Check {
	return NewStorageCheckHandler(sg)
}
