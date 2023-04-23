/*
Copyright The Karbour Authors.

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

package search

import (
	"github.com/KusionStack/karbour/pkg/apis/search"
	"github.com/KusionStack/karbour/pkg/registry"
	uniresourceregistry "github.com/KusionStack/karbour/pkg/registry/search/uniresource"
	"github.com/KusionStack/karbour/pkg/scheme"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

var _ registry.RESTStorageProvider = &RESTStorageProvider{}

type RESTStorageProvider struct{}

func (p RESTStorageProvider) GroupName() string {
	return search.GroupName
}

func (p RESTStorageProvider) NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(search.GroupName, scheme.Scheme, scheme.ParameterCodec, scheme.Codecs)
	v1beta1 := map[string]rest.Storage{}
	apiGroupInfo.VersionedResourcesStorageMap["v1beta1"] = v1beta1
	v1beta1["uniresources"] = uniresourceregistry.NewREST()
	return apiGroupInfo, nil
}
