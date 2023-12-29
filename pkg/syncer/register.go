package syncer

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	embedConfig "github.com/KusionStack/karbour/config"
	"github.com/KusionStack/karbour/pkg/kubernetes/scheme"
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

func StrategyRegister(hookContext genericapiserver.PostStartHookContext) error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-hookContext.StopCh:
			cancel()
		}
	}()
	client := dynamic.NewForConfigOrDie(hookContext.LoopbackClientConfig)
	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(hookContext.LoopbackClientConfig)
	cache := memory.NewMemCacheClient(discoveryClient)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cache)
	return createResources(ctx, client, mapper, embedConfig.DefaultSyncStrategy)
}

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
		if err = createResource(ctx, client, mapper, data); err != nil {
			return fmt.Errorf("failed to create resource: %w", err)
		}
	}
	return nil
}

func createResource(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, data []byte) error {
	obj, gvk, err := scheme.Codecs.UniversalDeserializer().Decode(data, nil, &unstructured.Unstructured{})
	if err != nil {
		return fmt.Errorf("could not decode data: %w", err)
	}
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("decoded into incorrect type, got %T, wanted %T", obj, &unstructured.Unstructured{})
	}
	m, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("could not get REST mapping for %s: %w", gvk, err)
	}
	_, err = client.Resource(m.Resource).Namespace(u.GetNamespace()).Create(ctx, u, metav1.CreateOptions{})
	if err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}
