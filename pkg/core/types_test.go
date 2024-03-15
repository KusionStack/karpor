package core

import "testing"

func TestSortCustomResourceGroup(t *testing.T) {
	tests := []struct {
		name    string
		in      map[string]any
		exp     string
		wantErr bool
	}{
		{
			name:    "test1",
			in:      map[string]any{"banana": 3, "apple": 5, "pear": 6, "orange": 2},
			exp:     `{"apple":5,"banana":3,"orange":2,"pear":6}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortCustomResourceGroup(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("SortCustomResourceGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.exp {
				t.Errorf("SortCustomResourceGroup() got = %v, want %v", got, tt.exp)
			}
		})
	}
}
