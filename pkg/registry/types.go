package registry

import (
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// RESTStorageProvider is a factory type for REST storage.
type RESTStorageProvider interface {
	GroupName() string
	NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error)
}
