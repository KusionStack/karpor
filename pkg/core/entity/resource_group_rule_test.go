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

package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceGroupRule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    *ResourceGroupRule
		wantErr error
	}{
		{
			name: "ValidRule",
			rule: &ResourceGroupRule{
				Name: "valid-rule",
			},
			wantErr: nil,
		},
		{
			name:    "NilRule",
			rule:    nil,
			wantErr: fmt.Errorf("resource group rule is nil"),
		},
		{
			name:    "EmptyName",
			rule:    &ResourceGroupRule{},
			wantErr: fmt.Errorf("resource group rule must have a name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.rule.Validate()
			if tt.wantErr != nil {
				require.EqualError(t, gotErr, tt.wantErr.Error())
			} else {
				require.NoError(t, gotErr)
			}
		})
	}
}
