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
	"context"
	"testing"

	"github.com/KusionStack/karpor/pkg/infra/ai"
	"github.com/stretchr/testify/require"
)

// fakeAIProvider is a minimal ai.AIProvider stub that always returns a
// fixed response, used to drive AIManager without a real AI backend.
type fakeAIProvider struct {
	response string
}

func (f *fakeAIProvider) Configure(ai.AIConfig) error { return nil }

func (f *fakeAIProvider) Generate(context.Context, string) (string, error) {
	return f.response, nil
}

func (f *fakeAIProvider) GenerateStream(context.Context, string) (<-chan string, error) {
	return nil, nil
}

// TestConvertTextToSQL tests the correctness of the ConvertTextToSQL function.
func TestConvertTextToSQL(t *testing.T) {
	testCases := []struct {
		name        string
		response    string
		expectedSQL string
		expectedErr error
	}{
		{
			// Regression test: a query like "Which pods are in error state"
			// legitimately produces a WHERE clause containing the word
			// "error" as a data value. This must not be rejected.
			name:        "ValidQueryWithErrorAsValue",
			response:    "SQL: select * from resources where content like '%Error%';",
			expectedSQL: "select * from resources where content like '%Error%'",
		},
		{
			name:        "InvalidQuery",
			response:    "Error",
			expectedErr: ErrInvalidQuery,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mgr := &AIManager{client: &fakeAIProvider{response: tc.response}}
			sql, err := mgr.ConvertTextToSQL("some query")
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedSQL, sql)
		})
	}
}
