package utils

import (
	"context"
	"testing"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
	"github.com/KusionStack/karpor/pkg/syncer/cache"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestImportTo(t *testing.T) {
	// Set up mocks for dynamic client
	mockey.Mock((*elasticsearch.Storage).SearchByQuery).Return(&storage.SearchResult{
		Resources: []*storage.Resource{
			{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "default",
						"namespace": "existing-namespace",
					},
				},
			},
		},
	}, nil).Build()
	defer mockey.UnPatchAll()

	// Test cases
	testCases := []struct {
		name           string
		gvr            schema.GroupVersionResource
		expectedLength int
		expectError    bool
	}{
		{
			name: "Success - ImportTo",
			gvr: schema.GroupVersionResource{
				Version:  "v1",
				Resource: "Pod",
			},
			expectedLength: 1,
			expectError:    false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call ListKeys method
			esImporter := NewESImporter(&elasticsearch.Storage{}, "defalut", tc.gvr)
			err := esImporter.ImportTo(context.TODO(), cache.NewResourceCache())

			// Check error expectation
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}
