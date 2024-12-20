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
	"strings"
)

// InterpretEvent represents a single event in the YAML interpretation stream
type InterpretEvent struct {
	Type    string `json:"type"`              // Event type: start, chunk, error, complete
	Content string `json:"content,omitempty"` // Event content or error message
}

// InterpretYAML performs AI interpretation of YAML content and sends events through the channel
func (a *AIManager) InterpretYAML(ctx context.Context, yaml, language string, eventChan chan<- *InterpretEvent) error {
	defer close(eventChan)

	// Send start event
	eventChan <- &InterpretEvent{Type: "start"}

	// Validate input
	if yaml == "" {
		eventChan <- &InterpretEvent{
			Type:    "error",
			Content: "YAML content cannot be empty",
		}
		return fmt.Errorf("YAML content cannot be empty")
	}

	// Build prompt from template
	prompt := fmt.Sprintf(ServicePromptMap[YAMLInterpretType], language, yaml)

	// Get AI service client
	if a.client == nil {
		eventChan <- &InterpretEvent{
			Type:    "error",
			Content: "AI service not configured",
		}
		return fmt.Errorf("AI service not configured")
	}

	// Stream completion from AI service
	stream, err := a.client.GenerateStream(ctx, prompt)
	if err != nil {
		eventChan <- &InterpretEvent{
			Type:    "error",
			Content: fmt.Sprintf("Failed to start AI service: %v", err),
		}
		return fmt.Errorf("failed to start AI service: %w", err)
	}

	// Process stream
	var fullContent strings.Builder
	for chunk := range stream {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if strings.HasPrefix(chunk, "ERROR:") {
				eventChan <- &InterpretEvent{
					Type:    "error",
					Content: fmt.Sprintf("AI service error: %v", strings.TrimPrefix(chunk, "ERROR: ")),
				}
				return fmt.Errorf("AI service error: %v", chunk)
			}

			fullContent.WriteString(chunk)
			eventChan <- &InterpretEvent{
				Type:    "chunk",
				Content: chunk,
			}
		}
	}

	// Send complete event
	eventChan <- &InterpretEvent{
		Type:    "complete",
		Content: fullContent.String(),
	}
	return nil
}
