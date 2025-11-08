# CLI Patterns Cookbook (Discord SDK)

Source inspiration: `~/vibe-engineering/docs/design/CLI_PATTERNS_COOKBOOK.md` and `../agent-mobile/docs/design/CLI_PATTERNS_COOKBOOK.md`

## Package Organization

### Basic Structure
```
gosdk/
├── discord/
│   ├── client/       # Core Discord API client
│   ├── webhook/      # Webhook management
│   ├── interactions/ # Slash commands and components
│   ├── gateway/      # WebSocket gateway (future)
│   └── types/        # Shared types and models
├── config/           # Configuration management
├── logger/           # Structured logging
└── examples/         # Usage examples
```

## Common Patterns

### Client Initialization
```go
// Prefer builder pattern for complex configuration
client, err := discord.NewClient(
    discord.WithToken(token),
    discord.WithTimeout(30*time.Second),
    discord.WithRetries(3),
    discord.WithLogger(logger),
)
```

### Error Handling
```go
// Use typed errors for programmatic handling
if errors.Is(err, discord.ErrRateLimited) {
    // Handle rate limit
}

// Wrap errors with context
return fmt.Errorf("failed to send message to channel %s: %w", channelID, err)
```

### Context Propagation
```go
// Always accept context as first parameter
func (c *Client) SendMessage(ctx context.Context, channelID, content string) error {
    // Use context for cancellation and timeouts
    req, err := http.NewRequestWithContext(ctx, "POST", url, body)
    // ...
}
```

## Configuration Management

### Config File (YAML)
```yaml
discord:
  bot_token: ${DISCORD_BOT_TOKEN}  # env var substitution
  application_id: "123456789"
  webhooks:
    default: ${DISCORD_WEBHOOK}
    alerts: ${DISCORD_WEBHOOK_ALERTS}

client:
  timeout: 30s
  retries: 3
  rate_limit_strategy: adaptive

logging:
  level: info
  format: json
  output: artifacts/logs/discord.log
```

### Environment Variables
```bash
DISCORD_BOT_TOKEN=...
DISCORD_WEBHOOK=...
DISCORD_LOG_LEVEL=debug
```

## JSON Serialization

### Message Format
```go
type Message struct {
    ID        string    `json:"id"`
    ChannelID string    `json:"channel_id"`
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}

// Marshal to JSON
data, err := json.Marshal(msg)
```

### Response Format
```json
{
  "success": true,
  "message_id": "1234567890",
  "channel_id": "9876543210",
  "timestamp": "2025-11-08T12:34:56Z"
}
```

## Retry Logic

### Exponential Backoff
```go
func (c *Client) sendWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error

    backoff := time.Second
    for attempt := 0; attempt < c.maxRetries; attempt++ {
        resp, err = c.httpClient.Do(req)
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }

        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(backoff):
            backoff *= 2
        }
    }
    return resp, err
}
```

## Rate Limiting

### Track Per-Route
```go
type RateLimiter struct {
    limits map[string]*routeLimit
    mu     sync.RWMutex
}

func (rl *RateLimiter) Wait(ctx context.Context, route string) error {
    // Wait for rate limit window if needed
}
```

## Logging Patterns

### Structured Logging
```go
logger.Info("sending message",
    "channel_id", channelID,
    "content_length", len(content),
    "attempt", attempt,
)

logger.Error("failed to send message",
    "error", err,
    "channel_id", channelID,
    "elapsed_ms", elapsed.Milliseconds(),
)
```

### Debug Mode
```go
if logger.IsDebug() {
    logger.Debug("request details",
        "method", req.Method,
        "url", req.URL.String(),
        "headers", req.Header,
    )
}
```

## Testing Patterns

### Table-Driven Tests
```go
func TestParseMessage(t *testing.T) {
    tests := []struct {
        name    string
        input   []byte
        want    *Message
        wantErr bool
    }{
        {"valid message", []byte(`{"id":"123"}`), &Message{ID: "123"}, false},
        {"empty input", []byte{}, nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseMessage(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Compare got with tt.want
        })
    }
}
```

### Mock Interfaces
```go
type DiscordClient interface {
    SendMessage(ctx context.Context, channelID, content string) error
}

type MockClient struct {
    SendMessageFunc func(ctx context.Context, channelID, content string) error
}

func (m *MockClient) SendMessage(ctx context.Context, channelID, content string) error {
    return m.SendMessageFunc(ctx, channelID, content)
}
```

## Examples

### Send Message
```go
ctx := context.Background()
client, _ := discord.NewClient(discord.WithToken(token))

err := client.SendMessage(ctx, channelID, "Hello, world!")
if err != nil {
    log.Fatalf("failed to send message: %v", err)
}
```

### Send Webhook
```go
webhook := discord.NewWebhook(webhookURL)
err := webhook.Send(ctx, &discord.WebhookMessage{
    Content: "Build completed successfully",
    Embeds: []discord.Embed{
        {
            Title:       "CI Build #123",
            Description: "All tests passed",
            Color:       0x00FF00,
        },
    },
})
```

### Create Slash Command
```go
cmd := &discord.SlashCommand{
    Name:        "status",
    Description: "Check bot status",
    Handler: func(ctx context.Context, i *discord.Interaction) error {
        return i.Respond(ctx, "Bot is running!")
    },
}

err := client.RegisterCommand(ctx, cmd)
```
