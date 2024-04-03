package syncer

import (
	"testing"

	"github.com/KusionStack/karbour/pkg/kubernetes/apis/search/v1beta1"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_genUnObj(t *testing.T) {
	type args struct {
		sr  v1beta1.ResourceSyncRule
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    *unstructured.Unstructured
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				sr:  v1beta1.ResourceSyncRule{APIVersion: "v1", Resource: "pods"},
				key: "ns1/name1",
			},
			want: &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "pod",
				"metadata":   map[string]interface{}{"name": "name1", "namespace": "ns1"},
			}},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				sr:  v1beta1.ResourceSyncRule{},
				key: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := genUnObj(tt.args.sr, tt.args.key)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
