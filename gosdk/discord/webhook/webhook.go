package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yourusername/agent-discord/gosdk/discord/types"
)

// Client represents a Discord webhook client
type Client struct {
	webhookURL string
	httpClient *http.Client
	maxRetries int
	timeout    time.Duration
}

// Option is a functional option for configuring the webhook client
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithMaxRetries sets the maximum number of retry attempts
func WithMaxRetries(retries int) Option {
	return func(c *Client) {
		c.maxRetries = retries
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// NewClient creates a new webhook client
func NewClient(webhookURL string, opts ...Option) (*Client, error) {
	if webhookURL == "" {
		return nil, &types.ValidationError{
			Field:   "webhookURL",
			Message: "webhook URL is required",
		}
	}

	c := &Client{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
		maxRetries: 3,
		timeout:    30 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.httpClient.Timeout = c.timeout

	return c, nil
}

// Send sends a message via the webhook
func (c *Client) Send(ctx context.Context, msg *types.WebhookMessage) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("invalid webhook message: %w", err)
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook message: %w", err)
	}

	return c.sendWithRetry(ctx, body)
}

// SendSimple sends a simple text message via the webhook
func (c *Client) SendSimple(ctx context.Context, content string) error {
	return c.Send(ctx, &types.WebhookMessage{
		Content: content,
	})
}

func (c *Client) sendWithRetry(ctx context.Context, body []byte) error {
	var lastErr error
	backoff := time.Second

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
			}
		}

		req, err := http.NewRequestWithContext(ctx, "POST", c.webhookURL, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "DiscordWebhook/1.0")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = &types.NetworkError{Op: "request", Err: err}
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			return nil
		}

		// Read error response
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		apiErr := &types.APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}

		// Parse JSON error if available
		var errData struct {
			Message    string                 `json:"message"`
			Code       int                    `json:"code"`
			Errors     map[string]interface{} `json:"errors"`
			RetryAfter float64                `json:"retry_after"`
		}
		if json.Unmarshal(respBody, &errData) == nil {
			apiErr.Message = errData.Message
			apiErr.Code = errData.Code
			apiErr.Errors = errData.Errors
			if errData.RetryAfter > 0 {
				apiErr.RetryAfter = int(errData.RetryAfter)
			}
		}

		// Handle rate limiting
		if resp.StatusCode == 429 {
			if apiErr.RetryAfter > 0 {
				backoff = time.Duration(apiErr.RetryAfter) * time.Second
			}
			lastErr = apiErr
			continue
		}

		// Don't retry client errors (except rate limits)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return apiErr
		}

		// Retry server errors
		lastErr = apiErr
	}

	if lastErr != nil {
		return fmt.Errorf("webhook request failed after %d attempts: %w", c.maxRetries+1, lastErr)
	}

	return fmt.Errorf("webhook request failed after %d attempts", c.maxRetries+1)
}

// parseErrorResponse parses an HTTP error response into an APIError
func (c *Client) parseErrorResponse(resp *http.Response) *types.APIError {
	respBody, _ := io.ReadAll(resp.Body)

	apiErr := &types.APIError{
		StatusCode: resp.StatusCode,
		Message:    string(respBody),
	}

	// Parse JSON error if available
	var errData struct {
		Message    string                 `json:"message"`
		Code       int                    `json:"code"`
		Errors     map[string]interface{} `json:"errors"`
		RetryAfter float64                `json:"retry_after"`
	}
	if json.Unmarshal(respBody, &errData) == nil {
		apiErr.Message = errData.Message
		apiErr.Code = errData.Code
		apiErr.Errors = errData.Errors
		if errData.RetryAfter > 0 {
			apiErr.RetryAfter = int(errData.RetryAfter)
		}
	}

	return apiErr
}

// marshalJSON marshals a value to JSON
func marshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// waitWithBackoff returns a channel that closes after the specified duration
func waitWithBackoff(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// backoffFromSeconds converts seconds to a duration for backoff
func backoffFromSeconds(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}
