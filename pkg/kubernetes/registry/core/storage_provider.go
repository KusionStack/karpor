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

package core

import (
	"fmt"
	"time"

	podstore "github.com/KusionStack/karpor/pkg/kubernetes/registry/core/pod"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	namespacestore "k8s.io/kubernetes/pkg/registry/core/namespace/storage"
	secretstore "k8s.io/kubernetes/pkg/registry/core/secret/storage"
	serviceaccountstore "k8s.io/kubernetes/pkg/registry/core/serviceaccount/storage"
	"k8s.io/kubernetes/pkg/serviceaccount"
)

const GroupName = "core"

type RESTStorageProvider struct {
	ServiceAccountIssuer        serviceaccount.TokenGenerator
	ServiceAccountMaxExpiration time.Duration
}

func NewRESTStorageProvider(serviceAccouuntIssuer serviceaccount.TokenGenerator,
	serviceAccountMaxExpiration time.Duration) *RESTStorageProvider {
	return &RESTStorageProvider{
		ServiceAccountIssuer:        serviceAccouuntIssuer,
		ServiceAccountMaxExpiration: serviceAccountMaxExpiration,
	}
}

func (p RESTStorageProvider) GroupName() string {
	return GroupName
}

func (p RESTStorageProvider) NewRESTStorage(
	restOptionsGetter generic.RESTOptionsGetter,
) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.APIGroupInfo{
		PrioritizedVersions:          scheme.Scheme.PrioritizedVersionsForGroup(""),
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		Scheme:                       scheme.Scheme,
		ParameterCodec:               scheme.ParameterCodec,
		NegotiatedSerializer:         scheme.Codecs,
	}
	storage := map[string]rest.Storage{}

	namespaceStorage, namespaceStatusStorage, namespaceFinalizeStorage, err := namespacestore.NewREST(
		restOptionsGetter,
	)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}
	storage["namespaces"] = namespaceStorage
	storage["namespaces/status"] = namespaceStatusStorage
	storage["namespaces/finalize"] = namespaceFinalizeStorage

	secretStorage, err := secretstore.NewREST(restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}
	storage["secrets"] = secretStorage

	podStorage, err := podstore.NewStorage(restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}
	storage["pods"] = podStorage.Pod
	storage["pods/status"] = podStorage.Status

	var serviceAccountStorage *serviceaccountstore.REST
	if p.ServiceAccountIssuer != nil {
		serviceAccountStorage, err = serviceaccountstore.NewREST(
			restOptionsGetter,
			p.ServiceAccountIssuer,
			nil,
			p.ServiceAccountMaxExpiration,
			podStorage.Pod.Store,
			secretStorage.Store,
			false,
		)
	} else {
		serviceAccountStorage, err = serviceaccountstore.NewREST(restOptionsGetter, nil, nil, 0, nil, nil, false)
	}
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}
	storage["serviceaccounts"] = serviceAccountStorage
	if serviceAccountStorage.Token != nil {
		storage["serviceaccounts/token"] = serviceAccountStorage.Token
	} else {
		return apiGroupInfo, fmt.Errorf("get token rest storage failed")
	}

	apiGroupInfo.VersionedResourcesStorageMap["v1"] = storage
	return apiGroupInfo, nil
}
