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

package cluster

import (
	"kusionstack.io/karpor/pkg/kubernetes/apis/cluster"
	"kusionstack.io/karpor/pkg/kubernetes/registry"
	"kusionstack.io/karpor/pkg/kubernetes/scheme"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

var _ registry.RESTStorageProvider = &RESTStorageProvider{}

type RESTStorageProvider struct{}

func (p RESTStorageProvider) GroupName() string {
	return cluster.GroupName
}

func (p RESTStorageProvider) NewRESTStorage(
	apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter,
) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(
		cluster.GroupName,
		scheme.Scheme,
		scheme.ParameterCodec,
		scheme.Codecs,
	)

	v1beta1 := map[string]rest.Storage{}
	clusterStorage, err := NewREST(restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}

	v1beta1["clusters"] = clusterStorage.Cluster
	v1beta1["clusters/status"] = clusterStorage.Status
	v1beta1["clusters/proxy"] = clusterStorage.Proxy

	apiGroupInfo.VersionedResourcesStorageMap["v1beta1"] = v1beta1
	return apiGroupInfo, nil
}
