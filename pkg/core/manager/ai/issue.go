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

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/scanner"
)

// AuditData represents the aggregated data of scanner issues, including the
// original list of issues and their aggregated count based on title.
type AuditData struct {
	IssueTotal    int            `json:"issueTotal"`
	ResourceTotal int            `json:"resourceTotal"`
	BySeverity    map[string]int `json:"bySeverity"`
	IssueGroups   []*IssueGroup  `json:"issueGroups"`
}

// IssueGroup represents a group of resourceGroups tied to a specific issue.
type IssueGroup struct {
	Issue          scanner.Issue          `json:"issue"`
	ResourceGroups []entity.ResourceGroup `json:"resourceGroups"`
}

// InterpretIssues performs AI interpretation of scanner issues and sends events through the channel
func (a *AIManager) InterpretIssues(ctx context.Context, auditData *AuditData, language string, eventChan chan<- *InterpretEvent) error {
	defer close(eventChan)

	// Send start event
	eventChan <- &InterpretEvent{Type: "start"}

	// Validate input
	if auditData == nil || len(auditData.IssueGroups) == 0 {
		eventChan <- &InterpretEvent{
			Type:    "error",
			Content: "No issues to interpret",
		}
		return fmt.Errorf("no issues to interpret")
	}

	// Build issue summary
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Total Issues: %d\n", auditData.IssueTotal))
	summary.WriteString(fmt.Sprintf("Total Resources: %d\n", auditData.ResourceTotal))
	summary.WriteString("\nSeverity Distribution:\n")
	for severity, count := range auditData.BySeverity {
		summary.WriteString(fmt.Sprintf("- %s: %d\n", severity, count))
	}

	// Build issue details
	summary.WriteString("\nIssue Details:\n")
	for _, group := range auditData.IssueGroups {
		summary.WriteString(fmt.Sprintf("\n## Issue\n"))
		summary.WriteString(fmt.Sprintf("Title: %s\n", group.Issue.Title))
		summary.WriteString(fmt.Sprintf("Severity: %s\n", group.Issue.Severity))
		summary.WriteString(fmt.Sprintf("Scanner: %s\n", group.Issue.Scanner))
		if group.Issue.Message != "" {
			summary.WriteString(fmt.Sprintf("Message: %s\n", group.Issue.Message))
		}

		summary.WriteString("\nAffected Resources:\n")
		for _, rg := range group.ResourceGroups {
			summary.WriteString(fmt.Sprintf("- %s/%s (%s)\n", rg.Namespace, rg.Name, rg.Kind))
		}
	}

	// Build prompt from template
	prompt := fmt.Sprintf(ServicePromptMap[IssueInterpretType], language, summary.String())

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
