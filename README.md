# Discord Go SDK

A production-ready Go SDK for Discord interactions.

## Features

- **Webhooks**: Send messages, embeds, and files via Discord webhooks
- **Bot API**: Interact with channels, messages, and guilds
- **Slash Commands**: Register and handle slash commands
- **Rate Limiting**: Automatic rate limit handling with exponential backoff
- **Error Handling**: Comprehensive typed errors for programmatic handling
- **Context Support**: Full context support for cancellation and timeouts
- **Testing**: Comprehensive test coverage with table-driven tests
- **Developer-Friendly**: JSON outputs, structured logging, deterministic behavior

## Quick Links

- **[AGENTS.md](AGENTS.md)** - Development workflow and collaboration guide
- **[Rate Limit Guide](docs/guides/RATE_LIMITS.md)** - Strategy/configuration reference
- **[Webhook Guide](docs/guides/WEBHOOKS.md)** - End-to-end webhook workflows
- **[Open Questions](docs/OPEN_QUESTIONS.md)** - Active design discussions

## Project Structure

```
godiscord/
├── AGENTS.md              # Agent collaboration guide ⭐
├── README.md              # This file
├── LICENSE                # MIT License
├── gosdk/                 # Go SDK (main development) ⭐
│   ├── discord/           # Discord API packages
│   │   ├── client/        # Core API client
│   │   ├── webhook/       # Webhook functionality ✅
│   │   ├── interactions/  # Slash commands
│   │   └── types/         # Shared types and models ✅
│   ├── config/            # Configuration management ✅
│   ├── logger/            # Structured logging ✅
│   ├── examples/          # Usage examples
│   └── go.mod             # Go module definition
└── docs/
    ├── design/            # Design principles and patterns ⭐
    ├── plans/             # Project plans ⭐
    ├── progress/          # Status tracking
    ├── guides/            # How-to guides
    └── OPEN_QUESTIONS.md  # Active design discussions ⭐
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- Discord webhook URL or bot token
- (Optional) Code search tools: `fd`, `ag`, `ast-grep`

### Installation

```bash
# Clone the repository
git clone https://github.com/mtreilly/godiscord.git
cd godiscord

# Navigate to Go SDK
cd gosdk

# Download dependencies
go mod download
```

### Usage

#### Webhooks

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/mtreilly/godiscord/gosdk/discord/types"
    "github.com/mtreilly/godiscord/gosdk/discord/webhook"
)

func main() {
    // Create webhook client
    client, err := webhook.NewClient(
        "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN",
        webhook.WithMaxRetries(3),
        webhook.WithTimeout(30*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Send simple message
    ctx := context.Background()
    if err := client.SendSimple(ctx, "Hello, Discord!"); err != nil {
        log.Fatal(err)
    }

    // Send message with embed
    msg := &types.WebhookMessage{
        Content: "Build completed!",
        Embeds: []types.Embed{
            {
                Title:       "✅ Success",
                Description: "All tests passed",
                Color:       0x00FF00,
                Fields: []types.EmbedField{
                    {Name: "Duration", Value: "2m 34s", Inline: true},
                    {Name: "Coverage", Value: "87.5%", Inline: true},
                },
            },
        },
    }
    if err := client.Send(ctx, msg); err != nil {
        log.Fatal(err)
    }
}
```

### Running Examples

```bash
# Set webhook URL
export DISCORD_WEBHOOK="https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN"

# Run webhook example
cd gosdk/examples/webhook
go run main.go

# Threaded webhook example (requires DISCORD_WEBHOOK_THREAD_ID or forum channel)
cd ../webhook-thread
go run main.go
```

## Configuration

### Environment Variables

```bash
DISCORD_BOT_TOKEN=your-bot-token
DISCORD_WEBHOOK=https://discord.com/api/webhooks/...
DISCORD_LOG_LEVEL=info  # debug, info, warn, error
```

### Config File (config.yaml)

```yaml
discord:
  bot_token: ${DISCORD_BOT_TOKEN}
  application_id: "your-app-id"
  webhooks:
    default: ${DISCORD_WEBHOOK}
    alerts: ${DISCORD_WEBHOOK_ALERTS}

client:
  timeout: 30s
  retries: 3
  rate_limit:
    strategy: adaptive
    backoff_base: 1s
    backoff_max: 60s

logging:
  level: info
  format: json
  output: stderr
```

## Development

### Building

```bash
cd gosdk
go build ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run tests with race detection
go test -race ./...

# Test specific package
go test -v ./discord/webhook

# Golden JSON fixtures
go test ./discord/webhook -run Golden

# Live webhook smoke (requires DISCORD_WEBHOOK)
DISCORD_WEBHOOK=... go test -tags integration ./discord/webhook
```

### Code Formatting

```bash
# Format code
gofmt -w .

# Run static analysis
go vet ./...

# Run linter (if installed)
golangci-lint run
```

### Code Search

Use the recommended tools from `AGENTS.md`:

```bash
# Find function definitions
ast-grep --lang go -p 'func $NAME($$$) $$$ { $$$ }'

# Find struct definitions
ast-grep --lang go -p 'type $NAME struct { $$$ }'

# Find files
fd -e go

# Search content
fd -e go | ag --file-list - 'pattern'
```

## Documentation

- **[AGENTS.md](AGENTS.md)**: Guide for agent collaboration and development workflow
- **[docs/design/](docs/design/)**: Design principles and patterns
- **[docs/OPEN_QUESTIONS.md](docs/OPEN_QUESTIONS.md)**: Active design discussions
- **Go docs**: Run `go doc -all ./discord/webhook` or similar

## Contributing

See [AGENTS.md](AGENTS.md) for development workflow and collaboration guidelines.

### Git Workflow

```bash
# Commit format
feat(webhook): add file upload support
fix(client): handle rate limit edge case
docs(guides): add integration guide
test(webhook): add retry logic tests
```

### Code Style

- Follow Go idioms and best practices
- Use interfaces for testability
- Always propagate context
- Wrap errors with context
- Write godoc comments for all exported symbols
- Table-driven tests for comprehensive coverage

## Testing Discord Integration

### Test Webhooks

Create a test Discord server and webhook:

1. Create a Discord server (or use an existing one)
2. Create a webhook: Server Settings → Integrations → Webhooks → New Webhook
3. Copy the webhook URL
4. Set `DISCORD_WEBHOOK` environment variable
5. Run examples

### Test Bot

1. Create a Discord application: https://discord.com/developers/applications
2. Create a bot and copy the token
3. Set `DISCORD_BOT_TOKEN` environment variable
4. Invite bot to your test server
5. Run bot examples

## License

MIT License - see [LICENSE](LICENSE) for details.

## References

- Discord API: https://discord.com/developers/docs
- Rate Limits: https://discord.com/developers/docs/topics/rate-limits
- Webhooks: https://discord.com/developers/docs/resources/webhook
- Gateway: https://discord.com/developers/docs/topics/gateway
