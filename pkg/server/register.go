// Copyright The Karbour Authors.
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

package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/KusionStack/karbour/config"
	"github.com/KusionStack/karbour/pkg/kubernetes/scheme"
	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// ConfigRegister is a function that registers the server configuration.
func ConfigRegister(hookContext genericapiserver.PostStartHookContext) error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-hookContext.StopCh
		cancel()
	}()
	client := dynamic.NewForConfigOrDie(hookContext.LoopbackClientConfig)
	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(hookContext.LoopbackClientConfig)
	cache := memory.NewMemCacheClient(discoveryClient)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cache)
	for _, data := range config.DefaultConfig {
		if err := createResources(ctx, client, mapper, data); err != nil {
			return err
		}
	}
	return nil
}

// createResources is a function that creates multiple Kubernetes resources based on the provided data.
func createResources(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	reader := yaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(data)))
	for {
		b, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		if len(bytes.TrimSpace(b)) == 0 {
			continue
		}
		if err = createResource(ctx, client, mapper, b); err != nil {
			return fmt.Errorf("failed to create resource: %v", err)
		}
	}
	return nil
}

// createResource is a function that creates a single Kubernetes resource based on the provided data.
func createResource(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, data []byte) error {
	obj, gvk, err := scheme.Codecs.UniversalDeserializer().Decode(data, nil, &unstructured.Unstructured{})
	if err != nil {
		return fmt.Errorf("could not decode data: %v", err)
	}
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("decoded into incorrect type, got %T, wanted %T", obj, &unstructured.Unstructured{})
	}
	m, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("could not get REST mapping for %s: %v", gvk, err)
	}
	_, err = client.Resource(m.Resource).Namespace(u.GetNamespace()).Create(ctx, u, metav1.CreateOptions{})
	if err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}
