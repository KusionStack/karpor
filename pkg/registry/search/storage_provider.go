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
	"fmt"

	"github.com/KusionStack/karbour/pkg/apis/search"
	"github.com/KusionStack/karbour/pkg/infra/search/storage"
	"github.com/KusionStack/karbour/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/KusionStack/karbour/pkg/registry/search/syncclusterresources"
	"github.com/KusionStack/karbour/pkg/registry/search/transformrule"
	"github.com/KusionStack/karbour/pkg/scheme"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

var _ registry.RESTStorageProvider = &RESTStorageProvider{}

type RESTStorageProvider struct {
	SearchStorageType      string
	ElasticSearchAddresses []string
	ElasticSearchName      string
	ElasticSearchPassword  string
}

func (p RESTStorageProvider) GroupName() string {
	return search.GroupName
}

func (p RESTStorageProvider) NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(search.GroupName, scheme.Scheme, scheme.ParameterCodec, scheme.Codecs)
	searchStorageGetter, err := p.SearchStorageGetter()
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}

	storageMap, err := p.v1beta1Storage(restOptionsGetter, searchStorageGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}

	apiGroupInfo.VersionedResourcesStorageMap["v1beta1"] = storageMap
	return apiGroupInfo, nil
}

func (p RESTStorageProvider) v1beta1Storage(restOptionsGetter generic.RESTOptionsGetter, searchStorageGetter storage.SearchStorageGetter) (map[string]rest.Storage, error) {
	v1beta1Storage := map[string]rest.Storage{}

	syncClusterResources, syncClusterResourcesStatus, err := syncclusterresources.NewREST(restOptionsGetter)
	if err != nil {
		return map[string]rest.Storage{}, err
	}
	v1beta1Storage["syncclusterresources"] = syncClusterResources
	v1beta1Storage["syncclusterresources/status"] = syncClusterResourcesStatus

	transformRule, err := transformrule.NewREST(restOptionsGetter)
	if err != nil {
		return map[string]rest.Storage{}, err
	}
	v1beta1Storage["transformrules"] = transformRule

	return v1beta1Storage, nil
}

func (p RESTStorageProvider) SearchStorageGetter() (storage.SearchStorageGetter, error) {
	switch p.SearchStorageType {
	case "elasticsearch":
		return elasticsearch.NewSearchStorageGetter(p.ElasticSearchAddresses, p.ElasticSearchName, p.ElasticSearchPassword), nil
	default:
		return nil, fmt.Errorf("invalid search storage type %s", p.SearchStorageType)
	}
}
