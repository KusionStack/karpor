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
	"fmt"
	"github.com/KusionStack/karpor/pkg/infra/ai"
)

// ConvertTextToSQL converts natural language text to an SQL query
func (a *AIManager) ConvertTextToSQL(query string) (string, error) {
	servicePrompt := ai.ServicePromptMap[ai.Text2sqlType]
	prompt := fmt.Sprintf(servicePrompt, query)
	res, err := a.client.Generate(context.Background(), prompt)
	if err != nil {
		return "", err
	}
	return ExtractSelectSQL(res), nil
}

// FixSQL fix the error SQL
func (a *AIManager) FixSQL(sql string, query string, error string) (string, error) {
	servicePrompt := ai.ServicePromptMap[ai.SqlFixType]
	prompt := fmt.Sprintf(servicePrompt, query, sql, error)
	res, err := a.client.Generate(context.Background(), prompt)
	if err != nil {
		return "", err
	}
	return ExtractSelectSQL(res), nil
}
