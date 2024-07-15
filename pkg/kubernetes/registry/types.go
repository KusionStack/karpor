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

package registry

import (
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

// RESTStorageProvider is a factory type for REST storage.
type RESTStorageProvider interface {
	GroupName() string
	NewRESTStorage(
		apiResourceConfigSource serverstorage.APIResourceConfigSource,
		restOptionsGetter generic.RESTOptionsGetter,
	) (genericapiserver.APIGroupInfo, error)
}

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	SearchStorageType      string
	ElasticSearchAddresses []string
	ElasticSearchUsername  string
	ElasticSearchPassword  string
	ReadOnlyMode           bool
	GithubBadge            bool
	APIServerSecurePort    int
}
