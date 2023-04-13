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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	searchv1beta1 "github.com/KusionStack/karbour/pkg/apis/search/v1beta1"
	versioned "github.com/KusionStack/karbour/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/KusionStack/karbour/pkg/generated/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/KusionStack/karbour/pkg/generated/listers/search/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SyncResourcesInformer provides access to a shared informer and lister for
// SyncResourceses.
type SyncResourcesInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.SyncResourcesLister
}

type syncResourcesInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewSyncResourcesInformer constructs a new informer for SyncResources type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSyncResourcesInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSyncResourcesInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredSyncResourcesInformer constructs a new informer for SyncResources type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSyncResourcesInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SearchV1beta1().SyncResourceses().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SearchV1beta1().SyncResourceses().Watch(context.TODO(), options)
			},
		},
		&searchv1beta1.SyncResources{},
		resyncPeriod,
		indexers,
	)
}

func (f *syncResourcesInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSyncResourcesInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *syncResourcesInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&searchv1beta1.SyncResources{}, f.defaultInformer)
}

func (f *syncResourcesInformer) Lister() v1beta1.SyncResourcesLister {
	return v1beta1.NewSyncResourcesLister(f.Informer().GetIndexer())
}
