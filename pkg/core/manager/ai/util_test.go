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

package ai

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestExtractSelectSQL tests the correctness of the ExtractSelectSQL function.
func TestExtractSelectSQL(t *testing.T) {
	testCases := []struct {
		name     string
		sql      string
		expected string
	}{
		{
			name: "NormalCase",
			sql: "Q: 所有kind=namespace " +
				"Schema_links: [kind, namespace] " +
				"SQL: select * from resources where kind='namespace';",
			expected: "select * from resources where kind='namespace'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ExtractSelectSQL(tc.sql)
			require.Equal(t, tc.expected, actual)
		})
	}
}

// TestIsInvalidQuery tests the IsInvalidQuery function.
func TestIsInvalidQuery(t *testing.T) {
	testCases := []struct {
		name     string
		sql      string
		expected bool
	}{
		{
			name:     "ValidQueryWithoutError",
			sql:      "select * from resources where kind='namespace';",
			expected: false,
		},
		{
			name:     "InvalidQuery",
			sql:      "Error",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := IsInvalidQuery(tc.sql)
			require.Equal(t, tc.expected, actual)
		})
	}
}
