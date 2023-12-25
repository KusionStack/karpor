package insight

import (
	"context"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/multicluster"
	topologyutil "github.com/KusionStack/karbour/pkg/util/topology"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "sigs.k8s.io/yaml"
)

// GetResource returns the unstructured cluster object for a given cluster
func (m *InsightManager) GetResource(
	ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator,
) (*unstructured.Unstructured, error) {
	resourceGVR, err := topologyutil.GetGVRFromGVK(loc.APIVersion, loc.Kind)
	if err != nil {
		return nil, err
	}
	return client.DynamicClient.
		Resource(resourceGVR).
		Namespace(loc.Namespace).
		Get(ctx, loc.Name, metav1.GetOptions{})
}

// GetYAMLForResource returns the yaml byte array for a given cluster
func (m *InsightManager) GetYAMLForResource(
	ctx context.Context, client *multicluster.MultiClusterClient, loc *core.Locator,
) ([]byte, error) {
	obj, err := m.GetResource(ctx, client, loc)
	if err != nil {
		return nil, err
	}
	return k8syaml.Marshal(obj.Object)
}
