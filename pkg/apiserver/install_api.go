package apiserver

import (
	"fmt"

	"github.com/KusionStack/karbour/pkg/registry"
	podstore "github.com/KusionStack/karbour/pkg/registry/core/pod"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	namespacestore "k8s.io/kubernetes/pkg/registry/core/namespace/storage"
	secretstore "k8s.io/kubernetes/pkg/registry/core/secret/storage"
	serviceaccountstore "k8s.io/kubernetes/pkg/registry/core/serviceaccount/storage"
)

type LegecyRESTStorageProvider interface {
	ResourceName() string
	NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, bool, error)
}

func InstallAPIs(genericAPIServer *genericapiserver.GenericAPIServer, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...registry.RESTStorageProvider) error {
	for _, restStorageProvider := range restStorageProviders {
		groupName := restStorageProvider.GroupName()
		apiGroupInfo, err := restStorageProvider.NewRESTStorage(restOptionsGetter)
		if err != nil {
			return fmt.Errorf("problem initializing API group %q : %v", groupName, err)
		}

		if len(apiGroupInfo.VersionedResourcesStorageMap) == 0 {
			// If we have no storage for any resource configured, this API group is effectively disabled.
			// This can happen when an entire API group, version, or development-stage (alpha, beta, GA) is disabled.
			klog.Infof("API group %q is not enabled, skipping.", groupName)
			continue
		}

		if err = genericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
			return fmt.Errorf("problem install API group %q: %v", groupName, err)
		}

		klog.Infof("Enabling API group %q.", groupName)
	}
	return nil
}

func InstallLegacyAPI(genericAPIServer *genericapiserver.GenericAPIServer, restOptionsGetter generic.RESTOptionsGetter) error {
	apiGroupInfo := genericapiserver.APIGroupInfo{
		PrioritizedVersions:          legacyscheme.Scheme.PrioritizedVersionsForGroup(""),
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		Scheme:                       legacyscheme.Scheme,
		ParameterCodec:               legacyscheme.ParameterCodec,
		NegotiatedSerializer:         legacyscheme.Codecs,
	}
	storage := map[string]rest.Storage{}

	namespaceStorage, namespaceStatusStorage, namespaceFinalizeStorage, err := namespacestore.NewREST(restOptionsGetter)
	if err != nil {
		return err
	}
	storage["namespaces"] = namespaceStorage
	storage["namespaces/status"] = namespaceStatusStorage
	storage["namespaces/finalize"] = namespaceFinalizeStorage

	secretStorage, err := secretstore.NewREST(restOptionsGetter)
	if err != nil {
		return err
	}
	storage["secrets"] = secretStorage

	serviceAccountStorage, err := serviceaccountstore.NewREST(restOptionsGetter, nil, nil, 0, nil, nil, false)
	if err != nil {
		return err
	}
	storage["serviceaccounts"] = serviceAccountStorage

	podStorage, err := podstore.NewStorage(restOptionsGetter)
	if err != nil {
		return nil
	}
	storage["pods"] = podStorage.Pod
	storage["pods/status"] = podStorage.Status

	apiGroupInfo.VersionedResourcesStorageMap["v1"] = storage

	if err := genericAPIServer.InstallLegacyAPIGroup(genericapiserver.DefaultLegacyAPIPrefix, &apiGroupInfo); err != nil {
		return fmt.Errorf("error in registering group versions: %v", err)
	}
	return nil
}
