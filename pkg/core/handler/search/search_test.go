// Copyright The Karpor Authors.
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

package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name         string
		pageSize     int
		page         int
		wantPageSize int
		wantPage     int
	}{
		{"missing values use defaults", 0, 0, 10, 1},
		{"negative values use defaults", -5, -1, 10, 1},
		{"page size 1 is kept", 1, 1, 1, 1},
		{"explicit values are kept", 25, 3, 25, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPageSize, gotPage := normalizePagination(tt.pageSize, tt.page)
			assert.Equal(t, tt.wantPageSize, gotPageSize)
			assert.Equal(t, tt.wantPage, gotPage)
		})
	}
}
