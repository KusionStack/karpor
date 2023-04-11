package cluster

import (
	"github.com/KusionStack/karbour/pkg/apis/cluster"
	"github.com/KusionStack/karbour/pkg/registry"
	"github.com/KusionStack/karbour/pkg/scheme"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

var _ registry.RESTStorageProvider = &RESTStorageProvider{}

type RESTStorageProvider struct{}

func (p RESTStorageProvider) GroupName() string {
	return cluster.GroupName
}

func (p RESTStorageProvider) NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(cluster.GroupName, scheme.Scheme, scheme.ParameterCodec, scheme.Codecs)

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
