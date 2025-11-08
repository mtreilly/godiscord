# Project Structure

Generated: 2025-11-08

## Overview

This repository contains a Go SDK for Discord interactions, designed to be integrated into the vibe CLI. The old Python Discord bot implementation is preserved in `discord-bot/` for reference only.

## Directory Layout

```
agent-discord/
├── AGENTS.md                      # Agent collaboration guide ⭐
├── README.md                      # Project overview and quick start ⭐
├── PROJECT_STRUCTURE.md           # This file
├── .env.example                   # Environment variable template
├── .gitignore                     # Git ignore rules
│
├── docs/                          # Documentation
│   ├── OPEN_QUESTIONS.md          # Active design discussions ⭐
│   ├── design/                    # Design principles and patterns
│   │   ├── _INDEX.md
│   │   ├── CLI_DESIGN_PRINCIPLES.md ⭐
│   │   └── CLI_PATTERNS_COOKBOOK.md ⭐
│   ├── plans/                     # Project plans and roadmaps
│   │   └── ROADMAP.md             # Development roadmap ⭐
│   ├── progress/                  # Status tracking
│   │   └── STATUS.md              # Current status ⭐
│   ├── guides/                    # How-to guides (future)
│   └── manual/                    # API reference (future)
│
├── gosdk/                         # Go SDK (main development) ⭐
│   ├── README.md                  # SDK-specific README
│   ├── go.mod                     # Go module definition
│   ├── go.sum                     # Dependency checksums
│   │
│   ├── discord/                   # Discord API packages
│   │   ├── types/                 # Core types and models
│   │   │   ├── errors.go          # Error definitions
│   │   │   ├── message.go         # Message types
│   │   │   └── webhook.go         # Webhook types
│   │   ├── webhook/               # Webhook client
│   │   │   ├── webhook.go         # Implementation
│   │   │   └── webhook_test.go    # Tests
│   │   ├── client/                # Bot API client (future)
│   │   └── interactions/          # Slash commands (future)
│   │
│   ├── config/                    # Configuration management
│   │   └── config.go              # YAML config + env vars
│   │
│   ├── logger/                    # Structured logging
│   │   └── logger.go              # Logger implementation
│   │
│   ├── examples/                  # Usage examples
│   │   ├── webhook/               # Webhook example
│   │   │   ├── main.go
│   │   │   └── README.md
│   │   └── bot/                   # Bot example (future)
│   │
│   └── cmd/                       # CLI tools (future)
│       └── cli/                   # vibe CLI integration
│
└── discord-bot/                   # Old Python bot (reference only)
    └── ...                        # Preserved for reference
```

## Key Files to Read First

1. **AGENTS.md** - Start here for development workflow and collaboration
2. **README.md** - Project overview, quick start, and usage examples
3. **docs/design/CLI_DESIGN_PRINCIPLES.md** - Core design principles
4. **docs/design/CLI_PATTERNS_COOKBOOK.md** - Practical patterns and examples
5. **docs/OPEN_QUESTIONS.md** - Active design discussions
6. **docs/progress/STATUS.md** - Current development status
7. **docs/plans/ROADMAP.md** - Development roadmap

## Package Organization

### gosdk/discord/types
Core types, models, and error definitions:
- `errors.go`: Typed errors (RateLimitError, APIError, etc.)
- `message.go`: Message, User, Embed types
- `webhook.go`: WebhookMessage with validation

### gosdk/discord/webhook
Webhook client implementation:
- Send messages via Discord webhooks
- Automatic retries with exponential backoff
- Rate limit handling
- Context support

### gosdk/config
Configuration management:
- YAML file parsing with env var substitution
- Default configuration
- Precedence: params > env > config > defaults

### gosdk/logger
Structured logging:
- Multiple levels (debug, info, warn, error)
- JSON and text formats
- Field-based logging

## Current Status

### ✅ Implemented (Phase 1 Complete)
- Project structure and documentation
- Core types package
- Webhook client with retry logic
- Configuration management
- Structured logging
- Basic tests and examples

### 🚧 In Progress (Phase 2)
- Full webhook API (files, threads, edit/delete)
- Bot API client (channels, messages, guilds)
- Enhanced rate limiting
- Expanded test coverage

### 📋 Planned
- Phase 3: Slash commands and component interactions
- Phase 4: vibe CLI integration and API stability
- Phase 5: Gateway (WebSocket) support

## Development Commands

### Build
```bash
cd gosdk
go build ./...
```

### Test
```bash
go test ./...
go test -v -cover ./...
go test -race ./...
```

### Run Examples
```bash
export DISCORD_WEBHOOK="https://discord.com/api/webhooks/..."
cd gosdk/examples/webhook
go run main.go
```

### Code Search
```bash
# Find functions
ast-grep --lang go -p 'func $NAME($$$) $$$ { $$$ }'

# Find structs
ast-grep --lang go -p 'type $NAME struct { $$$ }'

# Search content
fd -e go | ag --file-list - 'pattern'
```

## Integration Points

### vibe CLI
- SDK designed as importable Go module
- Clean interfaces for CLI commands
- Configuration compatible with vibe's config system
- Examples in `gosdk/cmd/cli/` (future)

### Old Python Bot
- Reference implementation in `discord-bot/`
- Use for understanding features, NOT for direct translation
- Go implementation follows Go idioms

## Documentation Standards

- **All exported symbols**: Godoc comments required
- **Design decisions**: Document in `docs/OPEN_QUESTIONS.md`
- **Patterns**: Reference `docs/design/CLI_PATTERNS_COOKBOOK.md`
- **Status updates**: Update `docs/progress/STATUS.md`

## Git Workflow

Conventional commits with scopes:
```
feat(webhook): add file upload support
fix(client): handle rate limit edge case
docs(guides): add integration guide
test(webhook): add retry logic tests
refactor(types): simplify error handling
```

## References

- Discord API: https://discord.com/developers/docs
- Go best practices: https://go.dev/doc/effective_go
- Inspired by: `~/vibe-engineering`, `../agent-mobile`
