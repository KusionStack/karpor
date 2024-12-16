package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

// DiagnoseLogsHandler handles the HTTP streaming response for log diagnosis
func (a *AIManager) DiagnoseLogsHandler(ctx context.Context, logs []string, language string, w http.ResponseWriter) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create channel for diagnosis events
	eventChan := make(chan *DiagnosisEvent, 10)
	go func() {
		if err := a.DiagnoseLogs(ctx, logs, language, eventChan); err != nil {
			// Error already sent through eventChan
			return
		}
	}()

	// Stream events to client
	for event := range eventChan {
		data, err := json.Marshal(event)
		if err != nil {
			continue
		}
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
}
