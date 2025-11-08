package webhook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/yourusername/agent-discord/gosdk/discord/types"
)

const (
	// MaxFileSize is the maximum size for a single file (25MB)
	MaxFileSize = 25 * 1024 * 1024

	// MaxTotalSize is the maximum total size for all files (8MB for free, 100MB for nitro)
	// Using conservative 8MB limit
	MaxTotalSize = 8 * 1024 * 1024

	// MaxFiles is the maximum number of files per message
	MaxFiles = 10
)

// FileAttachment represents a file to be uploaded via webhook
type FileAttachment struct {
	// Name is the filename (e.g., "image.png")
	Name string

	// ContentType is the MIME type (e.g., "image/png")
	// If empty, defaults to "application/octet-stream"
	ContentType string

	// Reader provides the file content
	Reader io.Reader

	// Size is the file size in bytes (optional, for validation)
	Size int64
}

// Validate checks if the file attachment is valid
func (f *FileAttachment) Validate() error {
	if f.Name == "" {
		return &types.ValidationError{
			Field:   "name",
			Message: "filename is required",
		}
	}

	if f.Reader == nil {
		return &types.ValidationError{
			Field:   "reader",
			Message: "file reader is required",
		}
	}

	if f.Size > MaxFileSize {
		return &types.ValidationError{
			Field:   "size",
			Message: fmt.Sprintf("file size %d exceeds maximum %d bytes (25MB)", f.Size, MaxFileSize),
		}
	}

	return nil
}

// SendWithFiles sends a webhook message with file attachments
func (c *Client) SendWithFiles(ctx context.Context, msg *types.WebhookMessage, files []FileAttachment) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("invalid webhook message: %w", err)
	}

	if len(files) == 0 {
		return &types.ValidationError{
			Field:   "files",
			Message: "at least one file is required (use Send for messages without files)",
		}
	}

	if len(files) > MaxFiles {
		return &types.ValidationError{
			Field:   "files",
			Message: fmt.Sprintf("too many files: %d (maximum %d)", len(files), MaxFiles),
		}
	}

	// Validate all files
	var totalSize int64
	for i, file := range files {
		if err := file.Validate(); err != nil {
			return fmt.Errorf("file %d validation failed: %w", i, err)
		}
		totalSize += file.Size
	}

	if totalSize > MaxTotalSize {
		return &types.ValidationError{
			Field:   "files",
			Message: fmt.Sprintf("total file size %d exceeds maximum %d bytes (8MB)", totalSize, MaxTotalSize),
		}
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add JSON payload
	if err := c.writeJSONPayload(writer, msg); err != nil {
		return fmt.Errorf("failed to write JSON payload: %w", err)
	}

	// Add files
	for i, file := range files {
		if err := c.writeFile(writer, i, file); err != nil {
			return fmt.Errorf("failed to write file %d: %w", i, err)
		}
	}

	// Close multipart writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Send with retry
	return c.sendMultipartWithRetry(ctx, body.Bytes(), writer.FormDataContentType())
}

// writeJSONPayload writes the webhook message as JSON to the multipart form
func (c *Client) writeJSONPayload(writer *multipart.Writer, msg *types.WebhookMessage) error {
	// Create form field for JSON payload
	part, err := writer.CreateFormField("payload_json")
	if err != nil {
		return err
	}

	// Marshal message to JSON
	data, err := marshalJSON(msg)
	if err != nil {
		return err
	}

	// Write JSON to form field
	_, err = part.Write(data)
	return err
}

// writeFile writes a file attachment to the multipart form
func (c *Client) writeFile(writer *multipart.Writer, index int, file FileAttachment) error {
	// Create form file with unique field name
	fieldName := fmt.Sprintf("file%d", index)

	// Set content type if provided
	contentType := file.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Create part with custom headers
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, file.Name),
	}
	h["Content-Type"] = []string{contentType}

	part, err := writer.CreatePart(h)
	if err != nil {
		return err
	}

	// Copy file content to part
	_, err = io.Copy(part, file.Reader)
	return err
}

// sendMultipartWithRetry sends a multipart request with retry logic
func (c *Client) sendMultipartWithRetry(ctx context.Context, body []byte, contentType string) error {
	var lastErr error
	backoff := c.timeout / 30 // Start with ~1 second

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-waitWithBackoff(backoff):
				backoff *= 2
			}
		}

		req, err := http.NewRequestWithContext(ctx, "POST", c.webhookURL, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", contentType)
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

		// Handle error response (reuse existing logic)
		apiErr := c.parseErrorResponse(resp)
		resp.Body.Close()

		// Handle rate limiting
		if resp.StatusCode == 429 {
			if apiErr.RetryAfter > 0 {
				backoff = backoffFromSeconds(apiErr.RetryAfter)
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
		return fmt.Errorf("multipart request failed after %d attempts: %w", c.maxRetries+1, lastErr)
	}

	return fmt.Errorf("multipart request failed after %d attempts", c.maxRetries+1)
}
