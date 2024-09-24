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
	"github.com/KusionStack/karpor/pkg/infra/ai"
	"github.com/KusionStack/karpor/pkg/kubernetes/registry"
)

type AIManager struct {
	client ai.AIProvider
}

// NewAIManager returns a new AIManager object
func NewAIManager(c registry.ExtraConfig) (*AIManager, error) {
	if c.AIAuthToken == "" {
		return nil, ErrMissingAuthToken
	}
	aiClient := ai.NewClient(c.AIBackend)
	if err := aiClient.Configure(ai.ConvertToAIConfig(c)); err != nil {
		return nil, err
	}

	return &AIManager{
		client: aiClient,
	}, nil
}

// CheckAIManager check if the AI manager is created
func CheckAIManager(aiMgr *AIManager) error {
	if aiMgr == nil {
		return ErrMissingAuthToken
	}
	return nil
}
