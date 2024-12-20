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

// DiagnosisEvent represents a diagnosis streaming event
type DiagnosisEvent struct {
	Type    string `json:"type"`    // Event type: start/chunk/error/complete
	Content string `json:"content"` // Event content
}

// DiagnoseLogs analyzes pod logs using LLM and returns diagnostic information through a streaming channel
func (a *AIManager) DiagnoseLogs(ctx context.Context, logs []string, language string, eventChan chan<- *DiagnosisEvent) error {
	defer close(eventChan)

	// Combine logs into a single string, limited to last 1000 lines
	if len(logs) > 1000 {
		logs = logs[len(logs)-1000:]
	}
	logsStr := strings.Join(logs, "\n")

	// Send start event
	eventChan <- &DiagnosisEvent{
		Type:    "start",
		Content: "Starting log analysis...",
	}

	// Get prompt template and add language instruction
	servicePrompt := ServicePromptMap[LogDiagnosisType]
	if language == "" {
		language = "English"
	}
	prompt := fmt.Sprintf(servicePrompt, language, logsStr)

	// Generate diagnosis using LLM with streaming
	stream, err := a.client.GenerateStream(ctx, prompt)
	if err != nil {
		errEvent := &DiagnosisEvent{
			Type:    "error",
			Content: fmt.Sprintf("Failed to analyze logs: %v", err),
		}
		eventChan <- errEvent
		return fmt.Errorf("failed to generate log diagnosis: %v", err)
	}

	var fullContent strings.Builder
	for chunk := range stream {
		if strings.HasPrefix(chunk, "ERROR:") {
			errEvent := &DiagnosisEvent{
				Type:    "error",
				Content: fmt.Sprintf("Failed to receive diagnosis: %v", strings.TrimPrefix(chunk, "ERROR: ")),
			}
			eventChan <- errEvent
			return fmt.Errorf("failed to receive diagnosis chunk: %v", chunk)
		}

		fullContent.WriteString(chunk)
		eventChan <- &DiagnosisEvent{
			Type:    "chunk",
			Content: chunk,
		}
	}

	// Send complete event
	eventChan <- &DiagnosisEvent{
		Type:    "complete",
		Content: fullContent.String(),
	}

	return nil
}
