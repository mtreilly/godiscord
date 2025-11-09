AGENTS.md — Discord Go SDK (vibe CLI Integration)

Purpose
- Build a production-ready Go SDK for Discord interactions (webhooks, slash commands, bot APIs) to be integrated into the vibe CLI.
- Agent-friendly: consistent APIs, JSON outputs, excellent error handling, and comprehensive testing.
- Inspired by patterns from `../agent-mobile`, `~/vibe-engineering`, and Discord best practices.

Scope
- Discord webhooks (send messages, embeds, files)
- Bot API client (messages, channels, guilds, interactions)
- Slash commands and component interactions
- Rate limiting, retries, and error handling
- Configuration management (YAML, env vars)
- Integration-ready for vibe CLI

Quick Start (Local Dev)
- Prerequisites: Go 1.21+, Discord bot token, webhook URLs
- Verify tools:
  - `go version`
  - `which ag fd ast-grep` (code search tools)
- Set up environment:
  ```bash
  cp .env.example .env
  # Edit .env with DISCORD_BOT_TOKEN and DISCORD_WEBHOOK
  ```
- Run examples:
  ```bash
  cd gosdk
  go run examples/webhook/main.go
  go run examples/bot/main.go
  ```

Project Structure
```
agent-discord/
├── AGENTS.md              # This file
├── README.md              # Project overview
├── discord-bot/           # Old Python implementation (reference only)
├── gosdk/                 # Go SDK (main development)
│   ├── discord/           # Discord API packages
│   │   ├── client/        # Core API client
│   │   ├── webhook/       # Webhook functionality
│   │   ├── interactions/  # Slash commands & components
│   │   ├── gateway/       # WebSocket gateway (future)
│   │   └── types/         # Shared types and models
│   ├── config/            # Configuration management
│   ├── logger/            # Structured logging
│   ├── examples/          # Usage examples
│   ├── go.mod             # Go module definition
│   └── go.sum             # Dependency checksums
└── docs/
    ├── design/            # Design principles and patterns
    ├── plans/             # Project plans and roadmaps
    ├── guides/            # How-to guides
    ├── manual/            # API reference and manuals
    └── progress/          # Status updates and tracking
```

Code Search & Discovery

FORBIDDEN TOOLS
- NEVER use: `find`, `grep`, `ls -R`, `cat` (for searching)
- Exception: Only if user explicitly requests it

Tool Selection Matrix (MANDATORY)

| Task | ONLY Use | Example |
|------|----------|---------|
| Function/struct defs | `ast-grep` | `ast-grep --lang go -p 'func $NAME($$$) $$$ { $$$ }'` |
| Interface definitions | `ast-grep` | `ast-grep --lang go -p 'type $NAME interface { $$$ }'` |
| Import statements | `ast-grep` | `ast-grep --lang go -p 'import "$PKG"'` |
| File discovery | `fd` | `fd -e go` |
| Directory structure | `fd` + `tree` | `fd -t d \| tree --fromfile -L 2` |
| Content search | `ag` | `fd -e go \| ag --file-list - 'pattern'` |

Recipes
```bash
fd -e go | ag --file-list - 'http.Client|context.Context'  # HTTP usage
fd -e go | ag --file-list - 'TODO|FIXME'                   # Find TODOs
fd -e go | ag --file-list - 'json:"'                        # JSON tags
fd -HI -t f -E .git | ag --file-list - 'token|secret'      # Secrets scan
fd -t d -E .git -E vendor | tree --fromfile -L 2            # Repo structure
ast-grep --lang go -p 'type $NAME struct { $$$ }'          # Find structs
```

Configuration Management

Config File (config.yaml)
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
  rate_limit_strategy: adaptive

logging:
  level: info
  format: json
  output: artifacts/logs/discord.log
```

Environment Variables
- `DISCORD_BOT_TOKEN` — Bot authentication token
- `DISCORD_WEBHOOK` — Default webhook URL
- `DISCORD_WEBHOOK_*` — Additional webhook URLs (alerts, builds, etc.)
- `DISCORD_LOG_LEVEL` — Logging level (debug, info, warn, error)

Go Development Workflow

Code Organization
- Package-based organization (not feature-based for SDK)
- Packages: `discord/client`, `discord/webhook`, `discord/types`, etc.
- Interfaces for all external dependencies (HTTP, logging, config)
- Clear separation: API client, business logic, configuration

Testing
- Unit tests: `*_test.go` files alongside implementation
- Table-driven tests for comprehensive coverage
- Golden tests for JSON serialization
- Integration tests: `examples/` directory with runnable code
- Coverage target: >80% for core packages

Building
```bash
cd gosdk
go build ./...                    # Build all packages
go test ./...                     # Run all tests
go test -v -cover ./...           # Tests with coverage
go test -race ./...               # Race detection
```

Linting & Formatting
```bash
gofmt -w .                        # Format code
go vet ./...                      # Static analysis
golangci-lint run                 # Comprehensive linting (if installed)
```

Debugging
1. Read full error + stack trace
2. Add debug logging with structured fields
3. Use `go test -v` for verbose test output
4. Test incrementally after fixes
5. Document fix in commit message

Docs Organization
- `docs/plans/` — project and architecture plans
- `docs/guides/` — how-to guides (usage, integration, deployment)
  - `docs/guides/RATE_LIMITS.md` — strategy selection + troubleshooting
  - `docs/guides/WEBHOOKS.md` — end-to-end webhook workflows
  - `docs/guides/INTERACTIONS.md` — slash commands, components, modals, and server guidance
  - `docs/guides/GATEWAY.md` — gateway connection, sharding, and observability
  - `docs/guides/PHASE6.md` — Phase 6 advanced features/tests summary
  - `docs/guides/CLI_EXAMPLES.md` — quick usage patterns for the new CLI
  - `docs/guides/MIGRATION.md` — migrating from the Python bot to this SDK
  - `docs/guides/CLI_RELEASE.md` — release playbook for the CLI bundle
- `docs/manual/` — API reference and detailed manuals
- `docs/design/` — design principles and patterns (adapted from vibe-engineering)
- `docs/progress/` — status updates and phase tracking
- Root files: `README.md` (overview), `AGENTS.md` (this guide)

Design Standards (Critical)
- Adhere to `docs/design/CLI_DESIGN_PRINCIPLES.md` and `docs/design/CLI_PATTERNS_COOKBOOK.md`
- All packages should follow these principles: context support, proper error handling, structured logging, JSON-serializable types
- When deviating, document rationale in `docs/OPEN_QUESTIONS.md` and propose updates to design docs if needed

Git Workflow
- Commit early and often. Use: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`
- Scope commits by package, e.g., `feat(webhook):`, `fix(client):`, `docs(guides):`
- Before risky changes: commit, then proceed. Use `git status`, `git diff`, `git log` continuously
- Keep commits atomic and focused on single changes

Open Questions (Living Log)
- Maintain `docs/OPEN_QUESTIONS.md` actively
- When you encounter uncertainty, a blocked decision, or design tradeoffs, add an entry
- Keep entries concise with context, options, and next experiments
- Close items by linking to resolving commits/PRs/docs
- Treat this as part of handoff hygiene: leave open threads visible for the next agent

Integration with vibe CLI
- SDK designed to be imported as a Go module
- Packages expose clean interfaces for CLI commands
- Configuration integrates with vibe's config system
- Logging integrates with vibe's logging framework
- Examples demonstrate CLI integration patterns

Initial Roadmap (Phases)

Phase 1 (Current): Foundation
- [x] Scaffold project structure
- [x] Create AGENTS.md and design docs
- [ ] Implement core types package
- [ ] Basic webhook client with retries
- [ ] Configuration management
- [ ] Structured logging

Phase 2: Core Features
- [ ] Full webhook API (messages, embeds, files)
- [ ] Bot API client (messages, channels)
- [ ] Rate limiting and backoff
- [ ] Comprehensive error types
- [ ] Unit tests for core packages

Phase 3: Advanced Features
- [ ] Slash commands registration and handling
- [ ] Component interactions (buttons, select menus)
- [ ] Embed builder with validation
- [ ] File uploads and attachments

Phase 4: Integration & Polish
- [ ] Integration examples
- [ ] vibe CLI integration guide
- [ ] Performance benchmarks
- [ ] Documentation completion
- [ ] API stability review

Phase 5: Gateway (Future)
- [ ] WebSocket gateway connection
- [ ] Event handling framework
- [ ] Presence and status management
- [ ] Voice support (if needed)

Next Actions for Agents
- Review `docs/OPEN_QUESTIONS.md` for active discussions
- Check `docs/progress/STATUS.md` for current work
- Start with webhook implementation (most common use case)
- Write comprehensive tests alongside implementation
- Document all public APIs with godoc comments
- Add examples for each major feature

Reference: Old Python Implementation
- The `discord-bot/` directory contains the old Python implementation
- Use as reference for features and patterns, but NOT for direct translation
- Go implementation should follow Go idioms and best practices
- Key learnings: slash commands, task queues, registry patterns

Testing & Validation
- All exported functions must have tests
- Use table-driven tests for comprehensive coverage
- Mock external dependencies (HTTP, Discord API)
- Integration tests using real Discord test servers (when available)
- Golden tests for JSON marshaling/unmarshaling
- `go test -race ./discord/webhook` should stay green; run before handing off
- Optional real-webhook smoke tests live behind `-tags integration` and require `DISCORD_WEBHOOK`

Common Tasks

Send Webhook Message
```bash
cd gosdk
go run examples/webhook/main.go
```

Run Tests
```bash
cd gosdk
go test ./...
go test -v -cover ./discord/webhook
go test -race ./discord/webhook
go test ./discord/webhook -run Golden
DISCORD_WEBHOOK=... go test -tags integration ./discord/webhook
```

Format & Lint
```bash
cd gosdk
gofmt -w .
go vet ./...
```

Generate Documentation
```bash
cd gosdk
go doc -all ./discord/webhook
```

Borrowed Patterns (Do This)
- Small, atomic commits with conventional prefixes
- Interfaces for all external dependencies
- Context propagation throughout
- Structured logging with levels
- Table-driven tests
- Error wrapping with context
- Configuration precedence: params > env > config > defaults

Anti-Patterns (Don't Do This)
- Global state or singletons
- Blocking operations without context
- Silent failures or swallowed errors
- Inconsistent error handling
- Missing godoc comments
- Untested code paths

Discord API Resources
- Official docs: https://discord.com/developers/docs
- Rate limits: https://discord.com/developers/docs/topics/rate-limits
- Webhook guide: https://discord.com/developers/docs/resources/webhook
- Interactions: https://discord.com/developers/docs/interactions/application-commands

Quick Reference

Environment Setup
```bash
export DISCORD_BOT_TOKEN="your-bot-token"
export DISCORD_WEBHOOK="https://discord.com/api/webhooks/..."
export DISCORD_LOG_LEVEL="debug"
```

Run Examples
```bash
cd gosdk/examples/webhook && go run main.go
cd gosdk/examples/bot && go run main.go
```

Test Individual Package
```bash
cd gosdk
go test -v ./discord/webhook
go test -v -race ./discord/client
```

Build for vibe CLI Integration
```bash
cd gosdk
go build -o ../bin/discord-cli ./cmd/cli
```
