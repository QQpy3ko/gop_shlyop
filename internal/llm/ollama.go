package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gop_shlyop/internal/config"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultModel          = "gemma:2b-instruct"
	defaultRequestTimeout = 30 * time.Second
)

type Client struct {
	host       string
	model      string
	httpClient *http.Client
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
}

func NewClient(cfg config.OllamaConfig) *Client {
	return &Client{
		host:  cfg.Host,
		model: defaultModel,
		httpClient: &http.Client{
			Timeout: defaultRequestTimeout,
		},
	}
}

func (c *Client) AnalyzeSentiment(ctx context.Context, text string) (string, error) {
	prompt := fmt.Sprintf("Определи тональность следующего отзыва. Ответь одним словом: 'positive', 'negative' или 'neutral'.\n\nОтзыв: \"%s\"", text)

	reqPayload := ollamaGenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ollama request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned non-200 status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read ollama response body: %w", err)
	}

	var ollamaResp ollamaGenerateResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal ollama response: %w", err)
	}

	sentiment := strings.TrimSpace(strings.ToLower(ollamaResp.Response))

	// Basic validation of the response
	switch sentiment {
	case "positive", "negative", "neutral":
		return sentiment, nil
	default:
		return "neutral", fmt.Errorf("unexpected sentiment from ollama: %s", sentiment)
	}
}
