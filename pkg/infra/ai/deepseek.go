package ai

import (
	"context"
	"errors"
	"io"

	deepseek "github.com/sashabaranov/go-openai"
)

const (
	defaultDeepseekBaseURL = "https://api.deepseek.com/v1"
)

type DeepseekClient struct {
	client      *deepseek.Client
	model       string
	temperature float32
	topP        float32
}

func (c *DeepseekClient) Configure(cfg AIConfig) error {
	defaultConfig := deepseek.DefaultConfig(cfg.AuthToken)
	defaultConfig.BaseURL = defaultDeepseekBaseURL
	if cfg.BaseURL != "" {
		defaultConfig.BaseURL = cfg.BaseURL
	}

	if cfg.ProxyEnabled {
		defaultConfig.HTTPClient.Transport = GetProxyHTTPClient(cfg)
	}

	client := deepseek.NewClientWithConfig(defaultConfig)
	if client == nil {
		return errors.New("error creating Deepseek client")
	}

	c.client = client
	c.model = cfg.Model
	c.temperature = cfg.Temperature
	c.topP = cfg.TopP
	return nil
}

func (c *DeepseekClient) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx, deepseek.ChatCompletionRequest{
		Model: c.model,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: c.temperature,
		TopP:        c.topP,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no completion choices returned from response")
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *DeepseekClient) GenerateStream(ctx context.Context, prompt string) (<-chan string, error) {
	stream, err := c.client.CreateChatCompletionStream(ctx, deepseek.ChatCompletionRequest{
		Model: c.model,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: c.temperature,
		TopP:        c.topP,
		Stream:      true,
	})
	if err != nil {
		return nil, err
	}

	// Create buffered channel for response chunks
	resultChan := make(chan string, 100)

	// Start goroutine to handle streaming response
	go func() {
		defer close(resultChan)
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				// Send error as a special message
				resultChan <- "ERROR: " + err.Error()
				return
			}

			// Send non-empty content chunks
			if len(response.Choices) > 0 {
				chunk := response.Choices[0].Delta.Content
				if chunk != "" {
					resultChan <- chunk
				}
			}
		}
	}()

	return resultChan, nil
}
