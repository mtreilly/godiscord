# Discord Go SDK - Complete Implementation Plan

**Version**: 1.0
**Created**: 2025-11-08
**Status**: Active
**Target**: Production-ready Discord SDK with agentic workflow support

## Executive Summary

Build a comprehensive, production-ready Discord Go SDK that follows best practices, supports all major Discord API features, and is optimized for agentic workflows. The SDK will integrate seamlessly with the vibe CLI and provide excellent developer experience through clear interfaces, comprehensive testing, and thorough documentation.

## Core Principles for Agentic Workflows

### 1. Deterministic Behavior
- All operations produce consistent results given the same inputs
- No hidden state or side effects
- Explicit error handling with typed errors
- Predictable retry behavior

### 2. Observable Operations
- Structured logging for all operations
- JSON-serializable request/response types
- Operation metadata (timing, retries, rate limits)
- Debug mode with full request/response logging

### 3. Composable Architecture
- Small, focused packages with clear boundaries
- Interface-based design for mockability
- Middleware/plugin support for cross-cutting concerns
- Builder patterns for complex operations

### 4. Error Recovery
- Typed errors for programmatic handling
- Automatic retries with configurable strategies
- Rate limit handling with backoff
- Context support for cancellation

### 5. Testability
- All external dependencies injected via interfaces
- Mock implementations for testing
- Golden tests for JSON payloads
- Integration test helpers

---

# Phase 1: Foundation ✅ COMPLETE

**Duration**: 1 day
**Status**: ✅ Complete (2025-11-08)

## Deliverables ✅
- [x] Project structure and documentation
- [x] Core types (Message, Embed, User)
- [x] Error type hierarchy
- [x] Basic webhook client
- [x] Configuration system
- [x] Structured logger
- [x] Initial tests and examples

---

# Phase 2: Enhanced Webhook & Rate Limiting

**Duration**: 1 week
**Focus**: Complete webhook implementation, robust rate limiting
**Agent Tasks**: 6-8 atomic tasks

## 2.1: Webhook File Uploads (2 days)

### Task 2.1.1: Multipart Form Support
**Complexity**: Medium
**Dependencies**: Phase 1

**Implementation**:
```go
// gosdk/discord/webhook/multipart.go
type FileAttachment struct {
    Name        string
    ContentType string
    Reader      io.Reader
}

func (c *Client) SendWithFiles(ctx context.Context, msg *WebhookMessage, files []FileAttachment) error
```

**Steps**:
1. Create `multipart.go` with file attachment types
2. Implement multipart/form-data encoding
3. Support multiple files (up to 10 per message)
4. Handle file size limits (25MB per file, 8MB total)
5. Add Content-Disposition headers
6. Unit tests with mock files
7. Integration example with image upload

**Testing**:
- Table-driven tests for multipart encoding
- Mock file readers
- Size limit validation
- Content-Type handling

**2025-11-08 Update**:
- Added runtime byte counting to enforce per-file and aggregate limits even when attachment sizes are unknown.
- Detect attachment sizes via `Len()`/`io.Seeker` heuristics to validate totals before upload.
- Raised aggregate limit to match per-file limits to avoid rejecting valid single-file uploads; future work can make these configurable per-server tier.

**Agentic Considerations**:
- Clear error messages for file size violations
- JSON-serializable file metadata in logs
- Dry-run mode to validate without sending

### Task 2.1.2: Webhook Edit/Delete Operations
**Complexity**: Low
**Dependencies**: Task 2.1.1

**Implementation**:
```go
// Add to webhook.go
func (c *Client) Edit(ctx context.Context, messageID string, msg *WebhookMessage) error
func (c *Client) Delete(ctx context.Context, messageID string) error
func (c *Client) Get(ctx context.Context, messageID string) (*types.Message, error)
```

**Steps**:
1. Implement PATCH endpoint for editing
2. Implement DELETE endpoint
3. Implement GET endpoint for retrieving
4. Handle webhook token authentication
5. Update tests for all CRUD operations
6. Add examples for edit/delete workflows

**Testing**:
- Mock HTTP server responses
- Error handling (404, 403)
- Partial updates (only changed fields)

## 2.2: Advanced Rate Limiting (3 days)

### Task 2.2.1: Rate Limit Tracker
**Complexity**: High
**Dependencies**: Phase 1

**Architecture**:
```go
// gosdk/discord/ratelimit/tracker.go
type Bucket struct {
    Key       string
    Limit     int
    Remaining int
    Reset     time.Time
}

type Tracker interface {
    Wait(ctx context.Context, route string) error
    Update(route string, headers http.Header)
    GetBucket(route string) *Bucket
}

type MemoryTracker struct {
    buckets map[string]*Bucket
    mu      sync.RWMutex
}
```

**Steps**:
1. Create `ratelimit` package
2. Implement route-based bucketing
3. Parse Discord rate limit headers:
   - `X-RateLimit-Limit`
   - `X-RateLimit-Remaining`
   - `X-RateLimit-Reset`
   - `X-RateLimit-Bucket`
4. Implement waiting logic with context support
5. Add global rate limit handling
6. Thread-safe bucket updates
7. Comprehensive tests with concurrent requests

**Testing**:
- Concurrent request simulation
- Bucket expiry and cleanup
- Global rate limit scenarios
- Context cancellation during waits

**Agentic Considerations**:
- Export bucket state as JSON for monitoring
- Predictable wait times (no random jitter by default)
- Optional callback for rate limit events

#### Follow-up 2.2.1b: Route-aware bucket mapping _(added 2025-11-08 during Phase 2 review)_
- **Status**: ✅ Completed (2025-11-08)
- Problem: `MemoryTracker.Update` stores bucket data by `bucketKey`, but `Wait`/`GetBucket` read by route (`POST:https://...`). When Discord sends a bucket ID (most endpoints do), per-route waits never trigger and adaptive strategies never learn (see `gosdk/ratelimit/tracker.go:122-154`).
- Fix: Track both the canonical bucket (`bucketKey`) and route aliases. Store bucket structs by `bucketKey` and maintain a separate `routeToBucket` map so `Wait/GetBucket` can resolve quickly. Clean up both maps when a bucket expires.
- Steps:
  1. Extend `MemoryTracker` to include `routeToBucket map[string]string`.
  2. When updating, map the provided `route` to the resolved bucket key and reuse existing structs when possible.
  3. Update `Wait`/`GetBucket` to translate the incoming route through `routeToBucket`.
  4. Add tests covering: (a) Discord provides `bucketKey`, (b) Discord omits `bucketKey`, (c) bucket expiry removes stale aliases.
  5. Document the behavior in `docs/OPEN_QUESTIONS.md` once the approach is validated.

### Task 2.2.2: Rate Limit Strategies
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Medium  
**Dependencies**: Task 2.2.1

**Implementation**:
```go
// gosdk/discord/ratelimit/strategy.go
type Strategy interface {
    ShouldWait(bucket *Bucket) bool
    CalculateWait(bucket *Bucket) time.Duration
}

type ReactiveStrategy struct{}   // Wait on 429
type ProactiveStrategy struct{}  // Wait before hitting limit
type AdaptiveStrategy struct{}   // Learn from patterns
```

**Steps**:
1. Define Strategy interface
2. Implement Reactive (current behavior)
3. Implement Proactive (prevent 429s)
4. Implement Adaptive (learning-based)
5. Make strategy configurable
6. Add metrics collection
7. Benchmarks comparing strategies

**Testing**:
- ✅ `gosdk/ratelimit/strategy_test.go` covers Reactive, Proactive, and Adaptive behavior (table-driven tests for thresholds, safety margins, adaptive learning).
- TODO: Add longer-running simulations/benchmarks after tracker routing fix (see Follow-up 2.2.1b) so adaptive stats receive real bucket data.

**Verification Notes**:
- Implementation lives in `gosdk/ratelimit/strategy.go`. Ensure webhook client continues to default to `AdaptiveStrategy` (`gosdk/discord/webhook/webhook.go:61-84`).
- Adaptive `RecordRequest` is currently called only from webhook client; once bot client exists, reuse helpers to keep metrics consistent.

### Task 2.2.3: Integrate Rate Limiting
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Medium  
**Dependencies**: Tasks 2.2.1, 2.2.1b, 2.2.2

**What Shipped**:
1. **Webhook client parity**: `sendWithRetryToURL` now uses the shared `waitForRateLimit` helper + `buildRoute`, so every code path (JSON, multipart, CRUD) honors proactive waits and reactive tracker blocks even when no strategy is configured.
2. **Observability**: `waitForRateLimit` logs proactive and reactive waits, `c.recordStrategyOutcome` is reused everywhere, and rate-limit logs now include strategy names and durations.
3. **Configuration**: `ClientConfig` gained a `rate_limit` block with strategy/backoff knobs (`config/config.go`, new tests in `config_test.go`). ENV override `DISCORD_RATE_LIMIT_STRATEGY` and `.env.example` reflect the new option.
4. **Examples**: `gosdk/examples/webhook/main.go` reads the env var to demonstrate swapping strategies on the fly.
5. **Docs**: Added `docs/guides/RATE_LIMITS.md` covering strategies, configuration, and troubleshooting; STATUS/OPEN_QUESTIONS reference it.

**Next (Phase 3 tie-in)**:
- Bot API client constructors must accept injected trackers/strategies to share rate-limit state across packages.
- Add integration smoke tests against Discord’s staging server once CLI wiring is ready.

**Configuration**:
```yaml
client:
  rate_limit:
    strategy: adaptive  # reactive, proactive, adaptive
    global_limit: 50    # requests per second
    per_route_limit: true
    backoff_base: 1s
    backoff_max: 60s
```

## 2.3: Webhook Thread Support (1 day)

### Task 2.3.1: Thread Operations
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Low  
**Dependencies**: Task 2.1.1

**Highlights**:
- `types.WebhookMessage` enforces that only one of `thread_id` or `thread_name` is set, preventing invalid API calls.
- Webhook client already supports `SendToThread` and `CreateThread`; tests cover validation, query params, and multipart flows (`discord/webhook/thread_test.go`).
- Added runnable example `gosdk/examples/webhook-thread` plus `.env`/README updates so agents can try posting to a thread or creating forum threads locally.

**Remaining Considerations**:
- When the bot client lands (Phase 3), expose helper utilities for listing/validating thread IDs to help CLI users pick targets.

## 2.4: Testing & Documentation (1 day)

### Task 2.4.1: Comprehensive Tests
**Status**: ✅ Completed (2025-11-08)  

**Delivered**:
1. Coverage now sits at **82.6%** for `discord/webhook` (`go test ./discord/webhook -cover`).
2. Golden JSON tests guard serialized payloads (`discord/webhook/json_golden_test.go` + `testdata/golden/`).
3. Optional integration test behind `//go:build integration` allows real webhook smoke tests when `DISCORD_WEBHOOK` is set (`integration_test.go`).
4. Added `BenchmarkClientSend` for steady-state throughput measurements and `TestClientSendConcurrent` for race detection (validated via `go test -race ./discord/webhook`).
5. New test data + harness ensures multipart/rate-limit logging changes remain observable without HTTP regressions.

### Task 2.4.2: Documentation
**Status**: ✅ Completed (2025-11-08)  

**Delivered**:
1. Authored `docs/guides/WEBHOOKS.md` covering setup, sends, threads, rate limits, and testing workflow.
2. Updated AGENTS.md with guide references plus the canonical test commands (coverage, race, golden, integration).
3. README “Quick Links” + testing section now point to the webhook guide and advanced test targets.
4. Rate limit guide already live from Task 2.2.3; cross-links added so onboarding agents land on the right docs.
5. Verified exported symbols gained/retained godoc comments in recent packages (config/webhook/ratelimit); no undocumented exports remain from Phase 2 scope.

---

# Phase 3: Bot API Client

**Duration**: 2 weeks
**Focus**: Core Discord Bot API, channels, messages, guilds
**Agent Tasks**: 12-15 atomic tasks

## 3.1: HTTP Client Foundation (2 days) _(Current focus — kick-off 2025-11-08)_

### Task 3.1.1: Base HTTP Client
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: High  
**Dependencies**: Phase 2 (rate limiter, logger, config)

**Architecture**:
```go
// gosdk/discord/client/client.go
type Client struct {
    token       string
    httpClient  *http.Client
    rateLimiter ratelimit.Tracker
    logger      *logger.Logger
    baseURL     string
}

type Option func(*Client)

func NewClient(token string, opts ...Option) (*Client, error)

// Internal request methods
func (c *Client) do(ctx context.Context, method, path string, body interface{}) (*http.Response, error)
func (c *Client) get(ctx context.Context, path string, v interface{}) error
func (c *Client) post(ctx context.Context, path string, body, v interface{}) error
func (c *Client) patch(ctx context.Context, path string, body, v interface{}) error
func (c *Client) delete(ctx context.Context, path string) error
```

**Delivered**:
1. New `gosdk/discord/client` package with constructor + option set (HTTP client, base URL, logger, rate limiter, strategy, retries, timeout).
2. Shared rate-limit helper + adaptive strategy reuse from webhook; proactive/reactive waits logged via `logger.Debug`.
3. Authenticated `do` helper with JSON encode/decode, typed error parsing (`types.APIError`) and exponential backoff/429 handling.
4. Convenience wrappers `Get/Post/Patch/Delete` for downstream channel/guild helpers.
5. Test suite using `httptest.Server` to cover auth headers, retry behavior, context cancellation, API errors, and rate-limit waits; package now part of `go test ./...`.

**Testing**:
- Mock HTTP server for all status codes
- Authentication header verification
- Rate limit integration
- Retry behavior validation
- Middleware hook ordering + dry-run toggles

**Agentic Considerations**:
- Request/response logging to JSON
- Middleware for request tracing
- Dry-run mode (validate without executing)

### Task 3.1.2: Client Middleware System
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Medium  
**Dependencies**: Task 3.1.1

**Implementation**:
```go
// gosdk/discord/client/middleware.go
type Middleware func(next RequestHandler) RequestHandler

type RequestHandler func(ctx context.Context, req *http.Request) (*http.Response, error)

func (c *Client) Use(middleware Middleware)

// Built-in middleware
func LoggingMiddleware(logger *logger.Logger) Middleware
func RetryMiddleware(maxRetries int) Middleware
func MetricsMiddleware(collector MetricsCollector) Middleware
func TracingMiddleware(tracer Tracer) Middleware
```

**Delivered**:
1. Middleware primitives (`Request`, `Middleware`, `RequestHandler`) + `Client.Use`.
2. Built-in middleware: logging, retry (exponential backoff), metrics collector hook, dry-run short-circuit (foundation for tracing/custom ones).
3. Middleware-aware execution path inside `Client.do` so rate-limiter + HTTP transport remain shared.
4. Tests verify middleware ordering, retry behavior, metrics invocation, and dry-run bypass.

**Agentic Considerations**:
- Middleware for request/response capture
- Middleware for operation replay
- Middleware for cost tracking (API quota)
- Shared tracker guidance captured in OPEN_QUESTIONS Q6 (still open for CLI-scale coordination)

## 3.2: Channel Operations (3 days)

### Task 3.2.1: Channel Types and Models
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Medium  
**Dependencies**: Task 3.1.1

**Implementation**:
```go
// gosdk/discord/types/channel.go
type Channel struct {
    ID                   string            `json:"id"`
    Type                 ChannelType       `json:"type"`
    GuildID              string            `json:"guild_id,omitempty"`
    Name                 string            `json:"name,omitempty"`
    Topic                string            `json:"topic,omitempty"`
    Position             int               `json:"position,omitempty"`
    PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
    NSFW                 bool              `json:"nsfw,omitempty"`
    ParentID             string            `json:"parent_id,omitempty"`
    // ... more fields
}

type ChannelType int

const (
    ChannelTypeGuildText ChannelType = iota
    ChannelTypeDM
    ChannelTypeGuildVoice
    ChannelTypeGroupDM
    ChannelTypeGuildCategory
    // ... more types
)
```

**Delivered**:
1. `gosdk/discord/types/channel.go` with channel/permission/thread structs, enums, and flag constants that mirror Discord API payloads.
2. Validation helpers for channel objects + create params (name/topic limits, rate limit bounds, bitrate constraints).
3. Fluent `ChannelParamsBuilder` to assemble create payloads with compile-time guarantees.
4. JSON + validation tests covering the builder, marshaling, and failure scenarios (`channel_test.go`).

### Task 3.2.2: Channel CRUD Operations
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Medium  
**Dependencies**: Task 3.2.1

**Delivered**:
1. `gosdk/discord/client/channels.go` exposing `Channels()` service with `GetChannel`, `ModifyChannel`, `DeleteChannel`, `GetChannelMessages`.
2. Pagination + query builder with validation (limit bounds, around/before/after exclusivity).
3. Modify params reuse validation from types package, including audit-log reason header support.
4. Tests verifying payloads, headers, and query parameters via `httptest.Server` (`channels_test.go`).

**Examples/TODO**:
- Add CLI/godoc snippets once CRUD integrates with vibe CLI (tracked for Phase 4 docs).

### Task 3.2.3: Channel Message Operations
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Medium  
**Dependencies**: Task 3.2.2

**Delivered**:
1. `gosdk/discord/client/messages.go` with message service exposing Create/Edit/Delete/Get/BulkDelete helpers backed by the base client + middleware stack.
2. Reused existing `types.MessageCreateParams` + new `MessageEditParams`; added validation + tests for service-level inputs (empty IDs, >100 bulk) in `messages_test.go`.
3. Tests simulate Discord responses via `httptest.Server` to verify HTTP method, path, payload, and error handling per endpoint.
4. File uploads/components left for later (defer to webhook multipart path); ready for CLI integration for text/embeds flows.

**Next**: Reaction helpers (Task 3.3.1) + file upload parity once message CLI needs attachments.

## 3.3: Reaction Operations (1 day)

### Task 3.3.1: Reaction Methods
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Low  
**Dependencies**: Task 3.2.3

**Delivered**:
1. Reaction helpers added to `MessageService` (`CreateReaction`, `DeleteOwnReaction`, `DeleteUserReaction`, `GetReactions`, `DeleteAllReactions`).
2. Emoji encoding via URL escaping + validation to support unicode/custom emoji formats.
3. Pagination struct for `GetReactions`, reusing shared HTTP client + middleware stack.
4. Tests covering route construction, pagination query params, and validation paths (`messages_test.go`).

**Next**: Move on to guild operations (Task 3.4.x) once reactions integrate with CLI flows.

## 3.4: Guild Operations (3 days)

### Task 3.4.1: Guild Types and Models
**Status**: ✅ Completed (2025-11-08)  
**Complexity**: Medium  
**Dependencies**: Task 3.1.1

**Delivered**:
1. `gosdk/discord/types/guild.go` with Guild, Role, Member, Emoji, WelcomeScreen, and preview structs mirroring Discord’s API.
2. Validation helpers for guild/role objects plus unit tests covering happy-path JSON + validation failures (`guild_test.go`).
3. Types reference existing Channel/User/Message structures without circular deps, clearing the way for guild REST operations in Task 3.4.2.

### Task 3.4.2: Guild Operations
**Complexity**: Medium
**Dependencies**: Task 3.4.1

**Implementation**:
```go
// gosdk/discord/client/guilds.go
func (c *Client) GetGuild(ctx context.Context, guildID string, withCounts bool) (*types.Guild, error)
func (c *Client) GetGuildPreview(ctx context.Context, guildID string) (*types.GuildPreview, error)
func (c *Client) ModifyGuild(ctx context.Context, guildID string, params *ModifyGuildParams) (*types.Guild, error)
func (c *Client) GetGuildChannels(ctx context.Context, guildID string) ([]*types.Channel, error)
func (c *Client) CreateGuildChannel(ctx context.Context, guildID string, params *CreateChannelParams) (*types.Channel, error)
```

**Steps**:
1. Implement guild retrieval
2. Implement guild modification
3. Implement guild channel operations
4. Handle permissions
5. Tests with mock data
6. Examples

### Task 3.4.3: Role & Member Operations
**Complexity**: Medium
**Dependencies**: Task 3.4.2

**Implementation**:
```go
// Add to guilds.go
func (c *Client) GetGuildRoles(ctx context.Context, guildID string) ([]*types.Role, error)
func (c *Client) CreateGuildRole(ctx context.Context, guildID string, params *CreateRoleParams) (*types.Role, error)
func (c *Client) ModifyGuildRole(ctx context.Context, guildID, roleID string, params *ModifyRoleParams) (*types.Role, error)
func (c *Client) DeleteGuildRole(ctx context.Context, guildID, roleID string) error

func (c *Client) GetGuildMember(ctx context.Context, guildID, userID string) (*types.Member, error)
func (c *Client) ListGuildMembers(ctx context.Context, guildID string, params *ListMembersParams) ([]*types.Member, error)
func (c *Client) AddGuildMemberRole(ctx context.Context, guildID, userID, roleID string) error
func (c *Client) RemoveGuildMemberRole(ctx context.Context, guildID, userID, roleID string) error
```

**Steps**:
1. Implement role CRUD operations
2. Implement member operations
3. Handle role permissions
4. Pagination for member lists
5. Tests for all operations
6. Examples with role management

## 3.5: Testing & Documentation (2 days)

### Task 3.5.1: Client Integration Tests
**Steps**:
1. Create integration test suite (build tag: integration)
2. Mock Discord API server
3. End-to-end workflow tests
4. Performance benchmarks
5. Achieve >80% coverage on client package

### Task 3.5.2: Client Documentation
**Steps**:
1. Complete godoc for client package
2. Create bot client guide
3. Create authentication guide
4. Create error handling guide
5. Add examples for common workflows
6. Update AGENTS.md

---

# Phase 4: Slash Commands & Interactions

**Duration**: 2 weeks
**Focus**: Application commands, interactions, components
**Agent Tasks**: 10-12 atomic tasks

## 4.1: Interaction Models (2 days)

### Task 4.1.1: Interaction Types
**Complexity**: High
**Dependencies**: Phase 3

**Implementation**:
```go
// gosdk/discord/types/interaction.go
type Interaction struct {
    ID            string          `json:"id"`
    ApplicationID string          `json:"application_id"`
    Type          InteractionType `json:"type"`
    Data          *InteractionData `json:"data,omitempty"`
    GuildID       string          `json:"guild_id,omitempty"`
    ChannelID     string          `json:"channel_id,omitempty"`
    Member        *Member         `json:"member,omitempty"`
    User          *User           `json:"user,omitempty"`
    Token         string          `json:"token"`
    Version       int             `json:"version"`
    Message       *Message        `json:"message,omitempty"`
}

type InteractionType int

const (
    InteractionTypePing InteractionType = iota + 1
    InteractionTypeApplicationCommand
    InteractionTypeMessageComponent
    InteractionTypeApplicationCommandAutocomplete
    InteractionTypeModalSubmit
)

type InteractionData struct {
    ID       string                       `json:"id,omitempty"`
    Name     string                       `json:"name,omitempty"`
    Type     ApplicationCommandType       `json:"type,omitempty"`
    Resolved *ResolvedData                `json:"resolved,omitempty"`
    Options  []ApplicationCommandOption   `json:"options,omitempty"`
    CustomID string                       `json:"custom_id,omitempty"`
    Values   []string                     `json:"values,omitempty"`
}
```

**Steps**:
1. Define all interaction types
2. Define interaction data structures
3. Define resolved data types
4. Add JSON marshaling/unmarshaling
5. Validation methods
6. Comprehensive tests

### Task 4.1.2: Application Command Types
**Complexity**: High
**Dependencies**: Task 4.1.1

**Implementation**:
```go
// Add to interaction.go
type ApplicationCommand struct {
    ID                       string                     `json:"id,omitempty"`
    Type                     ApplicationCommandType     `json:"type,omitempty"`
    ApplicationID            string                     `json:"application_id,omitempty"`
    GuildID                  string                     `json:"guild_id,omitempty"`
    Name                     string                     `json:"name"`
    Description              string                     `json:"description"`
    Options                  []ApplicationCommandOption `json:"options,omitempty"`
    DefaultMemberPermissions *int64                     `json:"default_member_permissions,omitempty"`
    DMPermission             *bool                      `json:"dm_permission,omitempty"`
}

type ApplicationCommandType int

const (
    ApplicationCommandTypeChatInput ApplicationCommandType = iota + 1
    ApplicationCommandTypeUser
    ApplicationCommandTypeMessage
)

type ApplicationCommandOption struct {
    Type        ApplicationCommandOptionType `json:"type"`
    Name        string                       `json:"name"`
    Description string                       `json:"description"`
    Required    bool                         `json:"required,omitempty"`
    Choices     []ApplicationCommandChoice   `json:"choices,omitempty"`
    Options     []ApplicationCommandOption   `json:"options,omitempty"`
    // ... more fields
}
```

**Steps**:
1. Define command types (chat input, user, message)
2. Define option types (string, integer, boolean, user, channel, role, etc.)
3. Define choices for options
4. Add validation (name length, option limits)
5. Builder pattern for complex commands
6. Tests for all command types

## 4.2: Command Registration (2 days)

### Task 4.2.1: Command Management
**Complexity**: Medium
**Dependencies**: Task 4.1.2

**Implementation**:
```go
// gosdk/discord/client/commands.go
func (c *Client) GetGlobalApplicationCommands(ctx context.Context) ([]*types.ApplicationCommand, error)
func (c *Client) CreateGlobalApplicationCommand(ctx context.Context, cmd *types.ApplicationCommand) (*types.ApplicationCommand, error)
func (c *Client) EditGlobalApplicationCommand(ctx context.Context, cmdID string, cmd *types.ApplicationCommand) (*types.ApplicationCommand, error)
func (c *Client) DeleteGlobalApplicationCommand(ctx context.Context, cmdID string) error

func (c *Client) GetGuildApplicationCommands(ctx context.Context, guildID string) ([]*types.ApplicationCommand, error)
func (c *Client) CreateGuildApplicationCommand(ctx context.Context, guildID string, cmd *types.ApplicationCommand) (*types.ApplicationCommand, error)
func (c *Client) BulkOverwriteGlobalApplicationCommands(ctx context.Context, cmds []*types.ApplicationCommand) ([]*types.ApplicationCommand, error)
```

**Steps**:
1. Implement global command endpoints
2. Implement guild command endpoints
3. Implement bulk overwrite (atomic updates)
4. Handle command permissions
5. Tests with mock API
6. Examples with command registration

**Agentic Considerations**:
- Declarative command definitions (YAML/JSON)
- Command diff and sync utilities
- Version control for command schemas

### Task 4.2.2: Command Builder
**Complexity**: Medium
**Dependencies**: Task 4.2.1

**Implementation**:
```go
// gosdk/discord/interactions/builder.go
type CommandBuilder struct {
    cmd *types.ApplicationCommand
}

func NewSlashCommand(name, description string) *CommandBuilder
func (b *CommandBuilder) AddStringOption(name, description string, required bool) *CommandBuilder
func (b *CommandBuilder) AddIntegerOption(name, description string, required bool) *CommandBuilder
func (b *CommandBuilder) AddBooleanOption(name, description string, required bool) *CommandBuilder
func (b *CommandBuilder) AddUserOption(name, description string, required bool) *CommandBuilder
func (b *CommandBuilder) AddChannelOption(name, description string, required bool) *CommandBuilder
func (b *CommandBuilder) AddRoleOption(name, description string, required bool) *CommandBuilder
func (b *CommandBuilder) AddChoices(optionName string, choices ...ApplicationCommandChoice) *CommandBuilder
func (b *CommandBuilder) SetDefaultPermission(perm bool) *CommandBuilder
func (b *CommandBuilder) Build() (*types.ApplicationCommand, error)

// Example usage:
cmd, err := NewSlashCommand("greet", "Greet a user").
    AddUserOption("user", "User to greet", true).
    AddStringOption("message", "Custom message", false).
    Build()
```

**Steps**:
1. Create fluent builder API
2. Add all option types
3. Add validation in Build()
4. Add choice builders
5. Add subcommand/subcommand group support
6. Comprehensive tests
7. Examples showing builder usage

## 4.3: Interaction Responses (3 days)

### Task 4.3.1: Response Types
**Complexity**: High
**Dependencies**: Task 4.1.2

**Implementation**:
```go
// gosdk/discord/types/interaction_response.go
type InteractionResponse struct {
    Type InteractionResponseType     `json:"type"`
    Data *InteractionResponseData   `json:"data,omitempty"`
}

type InteractionResponseType int

const (
    InteractionResponseTypePong InteractionResponseType = iota + 1
    InteractionResponseTypeChannelMessageWithSource InteractionResponseType = 4
    InteractionResponseTypeDeferredChannelMessageWithSource
    InteractionResponseTypeDeferredUpdateMessage
    InteractionResponseTypeUpdateMessage
    InteractionResponseTypeApplicationCommandAutocompleteResult
    InteractionResponseTypeModal
)

type InteractionResponseData struct {
    TTS             bool              `json:"tts,omitempty"`
    Content         string            `json:"content,omitempty"`
    Embeds          []Embed           `json:"embeds,omitempty"`
    AllowedMentions *AllowedMentions  `json:"allowed_mentions,omitempty"`
    Flags           int               `json:"flags,omitempty"`
    Components      []Component       `json:"components,omitempty"`
    Attachments     []Attachment      `json:"attachments,omitempty"`
    Choices         []AutocompleteChoice `json:"choices,omitempty"`
    CustomID        string            `json:"custom_id,omitempty"`
    Title           string            `json:"title,omitempty"`
}
```

**Steps**:
1. Define all response types
2. Define response data structures
3. Handle ephemeral messages (flags)
4. Add validation
5. Tests for all response types

### Task 4.3.2: Interaction Client
**Complexity**: Medium
**Dependencies**: Task 4.3.1

**Implementation**:
```go
// gosdk/discord/interactions/client.go
type InteractionClient struct {
    baseClient *client.Client
}

func NewInteractionClient(c *client.Client) *InteractionClient

func (ic *InteractionClient) CreateInteractionResponse(ctx context.Context, interactionID, token string, resp *types.InteractionResponse) error
func (ic *InteractionClient) GetOriginalInteractionResponse(ctx context.Context, token string) (*types.Message, error)
func (ic *InteractionClient) EditOriginalInteractionResponse(ctx context.Context, token string, params *MessageEditParams) (*types.Message, error)
func (ic *InteractionClient) DeleteOriginalInteractionResponse(ctx context.Context, token string) error

func (ic *InteractionClient) CreateFollowupMessage(ctx context.Context, token string, params *MessageCreateParams) (*types.Message, error)
func (ic *InteractionClient) EditFollowupMessage(ctx context.Context, token, messageID string, params *MessageEditParams) (*types.Message, error)
func (ic *InteractionClient) DeleteFollowupMessage(ctx context.Context, token, messageID string) error
```

**Steps**:
1. Implement interaction response endpoint
2. Implement followup message endpoints
3. Handle interaction tokens (15-minute expiry)
4. Tests with mock responses
5. Examples with interaction flows

**Agentic Considerations**:
- Deferred response helpers (for long-running operations)
- Response templates
- Automatic followup for timeouts

### Task 4.3.3: Response Builders
**Complexity**: Low
**Dependencies**: Task 4.3.2

**Implementation**:
```go
// gosdk/discord/interactions/response_builder.go
type ResponseBuilder struct {
    resp *types.InteractionResponse
}

func NewMessageResponse(content string) *ResponseBuilder
func NewDeferredResponse() *ResponseBuilder
func NewModalResponse(customID, title string) *ResponseBuilder

func (b *ResponseBuilder) SetEphemeral() *ResponseBuilder
func (b *ResponseBuilder) AddEmbed(embed types.Embed) *ResponseBuilder
func (b *ResponseBuilder) AddComponents(components ...types.Component) *ResponseBuilder
func (b *ResponseBuilder) Build() *types.InteractionResponse
```

**Steps**:
1. Create builder patterns
2. Add convenience methods
3. Validation
4. Tests
5. Examples

## 4.4: Message Components (3 days)

### Task 4.4.1: Component Types
**Complexity**: High
**Dependencies**: Task 4.3.1

**Implementation**:
```go
// gosdk/discord/types/components.go
type Component interface {
    Type() ComponentType
}

type ComponentType int

const (
    ComponentTypeActionRow ComponentType = iota + 1
    ComponentTypeButton
    ComponentTypeSelectMenu
    ComponentTypeTextInput
)

type ActionRow struct {
    Type       ComponentType `json:"type"`
    Components []Component   `json:"components"`
}

type Button struct {
    Type     ComponentType `json:"type"`
    Style    ButtonStyle   `json:"style"`
    Label    string        `json:"label,omitempty"`
    Emoji    *Emoji        `json:"emoji,omitempty"`
    CustomID string        `json:"custom_id,omitempty"`
    URL      string        `json:"url,omitempty"`
    Disabled bool          `json:"disabled,omitempty"`
}

type ButtonStyle int

const (
    ButtonStylePrimary ButtonStyle = iota + 1
    ButtonStyleSecondary
    ButtonStyleSuccess
    ButtonStyleDanger
    ButtonStyleLink
)

type SelectMenu struct {
    Type        ComponentType    `json:"type"`
    CustomID    string           `json:"custom_id"`
    Options     []SelectOption   `json:"options,omitempty"`
    Placeholder string           `json:"placeholder,omitempty"`
    MinValues   int              `json:"min_values,omitempty"`
    MaxValues   int              `json:"max_values,omitempty"`
    Disabled    bool             `json:"disabled,omitempty"`
}

type SelectOption struct {
    Label       string  `json:"label"`
    Value       string  `json:"value"`
    Description string  `json:"description,omitempty"`
    Emoji       *Emoji  `json:"emoji,omitempty"`
    Default     bool    `json:"default,omitempty"`
}
```

**Steps**:
1. Define component interfaces
2. Define button types and styles
3. Define select menu types (string, user, role, channel, mentionable)
4. Define text input for modals
5. Custom JSON marshaling for polymorphic types
6. Validation (limits, required fields)
7. Comprehensive tests

### Task 4.4.2: Component Builders
**Complexity**: Medium
**Dependencies**: Task 4.4.1

**Implementation**:
```go
// gosdk/discord/interactions/components.go
type ButtonBuilder struct {
    button *types.Button
}

func NewButton(customID, label string, style types.ButtonStyle) *ButtonBuilder
func NewLinkButton(label, url string) *ButtonBuilder
func (b *ButtonBuilder) SetEmoji(emoji types.Emoji) *ButtonBuilder
func (b *ButtonBuilder) SetDisabled(disabled bool) *ButtonBuilder
func (b *ButtonBuilder) Build() *types.Button

type SelectMenuBuilder struct {
    menu *types.SelectMenu
}

func NewSelectMenu(customID string) *SelectMenuBuilder
func (b *SelectMenuBuilder) AddOption(label, value, description string) *SelectMenuBuilder
func (b *SelectMenuBuilder) SetPlaceholder(placeholder string) *SelectMenuBuilder
func (b *SelectMenuBuilder) SetMinMaxValues(min, max int) *SelectMenuBuilder
func (b *SelectMenuBuilder) Build() *types.SelectMenu

type ActionRowBuilder struct {
    row *types.ActionRow
}

func NewActionRow() *ActionRowBuilder
func (b *ActionRowBuilder) AddButton(button *types.Button) *ActionRowBuilder
func (b *ActionRowBuilder) AddSelectMenu(menu *types.SelectMenu) *ActionRowBuilder
func (b *ActionRowBuilder) Build() *types.ActionRow
```

**Steps**:
1. Create button builder
2. Create select menu builder
3. Create action row builder
4. Add validation in builders
5. Tests for all builders
6. Examples with component usage

### Task 4.4.3: Modal Support
**Complexity**: Medium
**Dependencies**: Task 4.4.2

**Implementation**:
```go
// Add to components.go
type Modal struct {
    CustomID   string          `json:"custom_id"`
    Title      string          `json:"title"`
    Components []ActionRow     `json:"components"`
}

type TextInput struct {
    Type        ComponentType   `json:"type"`
    CustomID    string          `json:"custom_id"`
    Style       TextInputStyle  `json:"style"`
    Label       string          `json:"label"`
    MinLength   int             `json:"min_length,omitempty"`
    MaxLength   int             `json:"max_length,omitempty"`
    Required    bool            `json:"required,omitempty"`
    Value       string          `json:"value,omitempty"`
    Placeholder string          `json:"placeholder,omitempty"`
}

type TextInputStyle int

const (
    TextInputStyleShort TextInputStyle = iota + 1
    TextInputStyleParagraph
)

type ModalBuilder struct {
    modal *Modal
}

func NewModal(customID, title string) *ModalBuilder
func (b *ModalBuilder) AddTextInput(customID, label string, style TextInputStyle, required bool) *ModalBuilder
func (b *ModalBuilder) Build() *Modal
```

**Steps**:
1. Define modal types
2. Define text input types
3. Create modal builder
4. Validation (5 text inputs max)
5. Tests
6. Examples with modal workflows

## 4.5: Interaction Server (2 days)

### Task 4.5.1: HTTP Interaction Handler
**Complexity**: High
**Dependencies**: Task 4.3.2

**Implementation**:
```go
// gosdk/discord/interactions/server.go
type Handler func(ctx context.Context, i *types.Interaction) (*types.InteractionResponse, error)

type Server struct {
    publicKey string
    handlers  map[string]Handler
    logger    *logger.Logger
}

func NewServer(publicKey string, opts ...ServerOption) *Server

func (s *Server) RegisterCommand(name string, handler Handler)
func (s *Server) RegisterComponent(customID string, handler Handler)
func (s *Server) RegisterModal(customID string, handler Handler)

func (s *Server) HandleInteraction(w http.ResponseWriter, r *http.Request)
func (s *Server) verifyRequest(r *http.Request) bool
```

**Steps**:
1. Create interaction server
2. Implement Ed25519 signature verification
3. Implement request routing by interaction type
4. Handle command interactions
5. Handle component interactions
6. Handle modal submissions
7. Auto-respond to ping interactions
8. Comprehensive tests with signed requests

**Agentic Considerations**:
- Request logging for debugging
- Dry-run mode for testing handlers
- Automatic error handling and logging

### Task 4.5.2: Interaction Router
**Complexity**: Medium
**Dependencies**: Task 4.5.1

**Implementation**:
```go
// Add to server.go
type Router struct {
    commands   map[string]Handler
    components map[string]Handler
    modals     map[string]Handler
    middleware []Middleware
}

func NewRouter() *Router

func (r *Router) Command(name string, handler Handler)
func (r *Router) Component(customID string, handler Handler)
func (r *Router) Modal(customID string, handler Handler)
func (r *Router) Use(middleware Middleware)

// Pattern matching for dynamic custom IDs
func (r *Router) ComponentPattern(pattern string, handler Handler)

type Middleware func(next Handler) Handler
```

**Steps**:
1. Create flexible router
2. Support exact and pattern matching
3. Add middleware support
4. Implement command groups
5. Tests for routing
6. Examples with complex routing

## 4.6: Testing & Documentation (2 days)

### Task 4.6.1: Interaction Tests
**Steps**:
1. Unit tests for all interaction types
2. Integration tests for command workflows
3. Mock interaction server
4. Component interaction tests
5. Achieve >80% coverage

### Task 4.6.2: Interaction Documentation
**Steps**:
1. Complete godoc for interactions package
2. Create slash commands guide
3. Create components guide
4. Create modals guide
5. Create interaction server guide
6. Add comprehensive examples
7. Update AGENTS.md

---

# Phase 5: Gateway (WebSocket)

**Duration**: 3 weeks
**Focus**: Real-time events via WebSocket
**Agent Tasks**: 15-18 atomic tasks

## 5.1: Gateway Foundation (1 week)

### Task 5.1.1: Gateway Types
**Complexity**: High
**Dependencies**: Phase 3

**Implementation**:
```go
// gosdk/discord/gateway/types.go
type OpCode int

const (
    OpCodeDispatch OpCode = iota
    OpCodeHeartbeat
    OpCodeIdentify
    OpCodePresenceUpdate
    OpCodeVoiceStateUpdate
    OpCodeResume OpCode = 6
    OpCodeReconnect
    OpCodeRequestGuildMembers
    OpCodeInvalidSession
    OpCodeHello
    OpCodeHeartbeatAck
)

type Payload struct {
    Op OpCode          `json:"op"`
    D  json.RawMessage `json:"d"`
    S  int             `json:"s,omitempty"`
    T  string          `json:"t,omitempty"`
}

type IdentifyPayload struct {
    Token      string               `json:"token"`
    Properties IdentifyProperties   `json:"properties"`
    Compress   bool                 `json:"compress,omitempty"`
    Intents    int                  `json:"intents"`
    Shard      []int                `json:"shard,omitempty"`
}
```

**Steps**:
1. Define all gateway opcodes
2. Define payload structures
3. Define identify/resume payloads
4. Define intents system
5. JSON marshaling
6. Tests

### Task 5.1.2: WebSocket Connection
**Complexity**: High
**Dependencies**: Task 5.1.1

**Implementation**:
```go
// gosdk/discord/gateway/connection.go
type Connection struct {
    token      string
    intents    int
    conn       *websocket.Conn
    mu         sync.Mutex
    heartbeat  *time.Ticker
    sessionID  string
    sequence   int
    logger     *logger.Logger
}

func NewConnection(token string, intents int, opts ...ConnectionOption) (*Connection, error)

func (c *Connection) Connect(ctx context.Context) error
func (c *Connection) Close() error
func (c *Connection) Send(ctx context.Context, payload *Payload) error
func (c *Connection) Receive(ctx context.Context) (*Payload, error)
```

**Steps**:
1. Implement WebSocket connection management
2. Handle TLS and compression
3. Implement send/receive with proper locking
4. Connection state management
5. Tests with mock WebSocket

### Task 5.1.3: Heartbeat & Reconnection
**Complexity**: High
**Dependencies**: Task 5.1.2

**Implementation**:
```go
// Add to connection.go
func (c *Connection) startHeartbeat(interval time.Duration)
func (c *Connection) stopHeartbeat()
func (c *Connection) sendHeartbeat(ctx context.Context) error

func (c *Connection) reconnect(ctx context.Context) error
func (c *Connection) resume(ctx context.Context) error
```

**Steps**:
1. Implement heartbeat ticker
2. Handle heartbeat ACKs
3. Detect zombie connections
4. Implement reconnection logic
5. Implement session resumption
6. Exponential backoff for reconnects
7. Tests for heartbeat and reconnection

**Agentic Considerations**:
- Connection state observable as JSON
- Reconnection events logged
- Configurable retry strategies

## 5.2: Event System (1 week)

### Task 5.2.1: Event Types
**Complexity**: High
**Dependencies**: Task 5.1.3

**Implementation**:
```go
// gosdk/discord/gateway/events.go
type Event interface {
    Type() string
}

type ReadyEvent struct {
    V          int      `json:"v"`
    User       *User    `json:"user"`
    Guilds     []*Guild `json:"guilds"`
    SessionID  string   `json:"session_id"`
    ResumeURL  string   `json:"resume_gateway_url"`
}

type MessageCreateEvent struct {
    *Message
}

type MessageUpdateEvent struct {
    *Message
}

type MessageDeleteEvent struct {
    ID        string `json:"id"`
    ChannelID string `json:"channel_id"`
    GuildID   string `json:"guild_id,omitempty"`
}

// ... define all 50+ event types
```

**Steps**:
1. Define all gateway events (50+ types)
2. Group by category (message, guild, channel, user, etc.)
3. JSON unmarshaling for each type
4. Event type registry
5. Comprehensive tests

### Task 5.2.2: Event Dispatcher
**Complexity**: High
**Dependencies**: Task 5.2.1

**Implementation**:
```go
// gosdk/discord/gateway/dispatcher.go
type EventHandler func(ctx context.Context, event Event) error

type Dispatcher struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
    logger   *logger.Logger
}

func NewDispatcher() *Dispatcher

func (d *Dispatcher) On(eventType string, handler EventHandler)
func (d *Dispatcher) OnMessageCreate(handler func(ctx context.Context, e *MessageCreateEvent) error)
func (d *Dispatcher) OnGuildCreate(handler func(ctx context.Context, e *GuildCreateEvent) error)
// ... type-safe handlers for common events

func (d *Dispatcher) Dispatch(ctx context.Context, eventType string, data json.RawMessage) error
```

**Steps**:
1. Create event dispatcher
2. Thread-safe handler registration
3. Type-safe handlers for common events
4. Generic handler for all events
5. Error handling and logging
6. Middleware support
7. Tests with mock events

**Agentic Considerations**:
- Event replay capability
- Event filtering/transformation
- Event persistence (optional)

### Task 5.2.3: Gateway Client
**Complexity**: High
**Dependencies**: Task 5.2.2

**Implementation**:
```go
// gosdk/discord/gateway/client.go
type Client struct {
    conn       *Connection
    dispatcher *Dispatcher
    intents    int
    status     string
    activity   *Activity
    logger     *logger.Logger
}

func NewClient(token string, intents int, opts ...ClientOption) (*Client, error)

func (c *Client) Connect(ctx context.Context) error
func (c *Client) Disconnect() error
func (c *Client) On(eventType string, handler EventHandler)

func (c *Client) UpdatePresence(ctx context.Context, status string, activity *Activity) error
func (c *Client) RequestGuildMembers(ctx context.Context, guildID string, query string, limit int) error
```

**Steps**:
1. Integrate connection and dispatcher
2. Handle identify/resume flows
3. Implement presence updates
4. Implement guild member requests
5. Graceful shutdown
6. Tests with full flow

## 5.3: Intents & Caching (3 days)

### Task 5.3.1: Intent System
**Complexity**: Medium
**Dependencies**: Task 5.1.1

**Implementation**:
```go
// gosdk/discord/gateway/intents.go
type Intent int

const (
    IntentGuilds Intent = 1 << iota
    IntentGuildMembers
    IntentGuildBans
    IntentGuildEmojis
    IntentGuildIntegrations
    IntentGuildWebhooks
    IntentGuildInvites
    IntentGuildVoiceStates
    IntentGuildPresences
    IntentGuildMessages
    IntentGuildMessageReactions
    IntentGuildMessageTyping
    IntentDirectMessages
    IntentDirectMessageReactions
    IntentDirectMessageTyping
    IntentMessageContent
    IntentGuildScheduledEvents
    IntentAutoModerationConfiguration
    IntentAutoModerationExecution
)

func AllIntents() int
func DefaultIntents() int
func (i Intent) Has(intent Intent) bool
```

**Steps**:
1. Define all intent flags
2. Helper functions for intent combinations
3. Documentation on privileged intents
4. Validation
5. Tests

### Task 5.3.2: State Cache
**Complexity**: High
**Dependencies**: Task 5.2.3

**Implementation**:
```go
// gosdk/discord/gateway/cache.go
type Cache interface {
    GetGuild(guildID string) (*Guild, bool)
    SetGuild(guild *Guild)
    RemoveGuild(guildID string)

    GetChannel(channelID string) (*Channel, bool)
    SetChannel(channel *Channel)

    GetMember(guildID, userID string) (*Member, bool)
    SetMember(guildID string, member *Member)

    // ... more cache methods
}

type MemoryCache struct {
    guilds   map[string]*Guild
    channels map[string]*Channel
    members  map[string]map[string]*Member
    mu       sync.RWMutex
}

func NewMemoryCache() *MemoryCache
```

**Steps**:
1. Define cache interface
2. Implement in-memory cache
3. Handle cache updates from events
4. TTL and eviction policies
5. Thread-safe operations
6. Cache statistics
7. Tests with concurrent access

**Agentic Considerations**:
- Cache dump/restore for debugging
- Cache statistics as JSON
- Pluggable cache backends (Redis, etc.)

### Task 5.3.3: Cache Integration
**Complexity**: Medium
**Dependencies**: Task 5.3.2

**Steps**:
1. Integrate cache with gateway client
2. Update cache on events (GUILD_CREATE, CHANNEL_UPDATE, etc.)
3. Provide helper methods to query cache
4. Make cache optional
5. Tests with event sequences
6. Examples with cache usage

## 5.4: Sharding (3 days)

### Task 5.4.1: Shard Manager
**Complexity**: High
**Dependencies**: Task 5.2.3

**Implementation**:
```go
// gosdk/discord/gateway/shard.go
type Shard struct {
    id     int
    total  int
    client *Client
}

type ShardManager struct {
    shards   []*Shard
    token    string
    intents  int
    logger   *logger.Logger
}

func NewShardManager(token string, shardCount int, intents int) *ShardManager

func (sm *ShardManager) Connect(ctx context.Context) error
func (sm *ShardManager) Disconnect() error
func (sm *ShardManager) On(eventType string, handler EventHandler)
func (sm *ShardManager) Broadcast(ctx context.Context, payload *Payload) error
```

**Steps**:
1. Implement shard identification
2. Implement shard manager
3. Staggered shard connections (5s delay)
4. Event aggregation from all shards
5. Shard-specific operations
6. Tests with multiple shards

### Task 5.4.2: Automatic Sharding
**Complexity**: Medium
**Dependencies**: Task 5.4.1

**Implementation**:
```go
// Add to shard.go
func (sm *ShardManager) AutoScale(ctx context.Context) error

type ShardingStrategy interface {
    Calculate(guildCount int) int
}

type RecommendedSharding struct{}
type FixedSharding struct{ Count int }
```

**Steps**:
1. Implement GET /gateway/bot endpoint
2. Parse recommended shard count
3. Implement sharding strategies
4. Auto-scale based on guild count
5. Tests
6. Examples

## 5.5: Testing & Documentation (2 days)

### Task 5.5.1: Gateway Tests
**Steps**:
1. Unit tests for all gateway components
2. Integration tests with mock WebSocket
3. Event flow tests
4. Reconnection scenario tests
5. Sharding tests
6. Achieve >80% coverage

### Task 5.5.2: Gateway Documentation
**Steps**:
1. Complete godoc for gateway package
2. Create gateway guide
3. Create intents guide
4. Create sharding guide
5. Create caching guide
6. Add comprehensive examples
7. Update AGENTS.md

---

# Phase 6: Advanced Features & Polish

**Duration**: 2 weeks
**Focus**: Permissions, embeds, utilities, performance
**Agent Tasks**: 8-10 atomic tasks

## 6.1: Permission System (3 days)

### Task 6.1.1: Permission Types
**Complexity**: High
**Dependencies**: Phase 3

**Implementation**:
```go
// gosdk/discord/permissions/permissions.go
type Permission int64

const (
    PermissionCreateInstantInvite Permission = 1 << iota
    PermissionKickMembers
    PermissionBanMembers
    PermissionAdministrator
    PermissionManageChannels
    PermissionManageGuild
    PermissionAddReactions
    PermissionViewAuditLog
    PermissionPrioritySpeaker
    PermissionStream
    PermissionViewChannel
    PermissionSendMessages
    PermissionSendTTSMessages
    PermissionManageMessages
    PermissionEmbedLinks
    PermissionAttachFiles
    PermissionReadMessageHistory
    PermissionMentionEveryone
    PermissionUseExternalEmojis
    PermissionViewGuildInsights
    PermissionConnect
    PermissionSpeak
    PermissionMuteMembers
    PermissionDeafenMembers
    PermissionMoveMembers
    PermissionUseVAD
    PermissionChangeNickname
    PermissionManageNicknames
    PermissionManageRoles
    PermissionManageWebhooks
    PermissionManageEmojis
    PermissionUseSlashCommands
    PermissionRequestToSpeak
    PermissionManageEvents
    PermissionManageThreads
    PermissionCreatePublicThreads
    PermissionCreatePrivateThreads
    PermissionUseExternalStickers
    PermissionSendMessagesInThreads
    PermissionUseEmbeddedActivities
    PermissionModerateMembers
)

func (p Permission) Has(perm Permission) bool
func (p Permission) Add(perm Permission) Permission
func (p Permission) Remove(perm Permission) Permission
func AllPermissions() Permission
```

**Steps**:
1. Define all permission bits
2. Implement permission helpers
3. Permission checking utilities
4. Tests for permission operations

### Task 6.1.2: Permission Calculator
**Complexity**: High
**Dependencies**: Task 6.1.1

**Implementation**:
```go
// Add to permissions.go
type PermissionCalculator struct {
    guild   *Guild
    channel *Channel
    member  *Member
}

func NewPermissionCalculator(guild *Guild, channel *Channel, member *Member) *PermissionCalculator

func (pc *PermissionCalculator) ComputeBasePermissions() Permission
func (pc *PermissionCalculator) ComputeOverwrites() Permission
func (pc *PermissionCalculator) Compute() Permission

func (pc *PermissionCalculator) Can(perm Permission) bool
func (pc *PermissionCalculator) CanManageChannel() bool
func (pc *PermissionCalculator) CanSendMessages() bool
// ... convenience methods for common checks
```

**Steps**:
1. Implement base permission calculation
2. Implement overwrite calculation
3. Handle administrator bypass
4. Handle role hierarchy
5. Comprehensive tests with complex scenarios
6. Examples

**Agentic Considerations**:
- Permission calculator as JSON (for debugging)
- Explain why permission was granted/denied

## 6.2: Embed Builder & Utilities (2 days)

### Task 6.2.1: Advanced Embed Builder
**Complexity**: Medium
**Dependencies**: Phase 1

**Implementation**:
```go
// gosdk/discord/embeds/builder.go
type Builder struct {
    embed *types.Embed
}

func New() *Builder
func (b *Builder) SetTitle(title string) *Builder
func (b *Builder) SetDescription(description string) *Builder
func (b *Builder) SetColor(color int) *Builder
func (b *Builder) SetURL(url string) *Builder
func (b *Builder) SetTimestamp(t time.Time) *Builder
func (b *Builder) SetFooter(text, iconURL string) *Builder
func (b *Builder) SetImage(url string) *Builder
func (b *Builder) SetThumbnail(url string) *Builder
func (b *Builder) SetAuthor(name, url, iconURL string) *Builder
func (b *Builder) AddField(name, value string, inline bool) *Builder
func (b *Builder) Build() (*types.Embed, error)

// Presets
func Success(title, description string) *Builder
func Error(title, description string) *Builder
func Warning(title, description string) *Builder
func Info(title, description string) *Builder
```

**Steps**:
1. Create fluent builder
2. Add validation (character limits)
3. Add color presets
4. Add template functions
5. Tests
6. Examples

### Task 6.2.2: Utility Functions
**Complexity**: Low
**Dependencies**: Phase 3

**Implementation**:
```go
// gosdk/discord/utils/utils.go
func ParseMention(mention string) (userID string, ok bool)
func FormatUserMention(userID string) string
func FormatChannelMention(channelID string) string
func FormatRoleMention(roleID string) string

func ParseEmoji(emoji string) (name, id string, animated bool)
func FormatEmoji(name, id string, animated bool) string

func Snowflake ToTime(snowflake string) time.Time
func TimeToSnowflake(t time.Time) string

func ChunkSlice[T any](slice []T, size int) [][]T
func RateLimitDelay(remaining, limit int, reset time.Time) time.Duration
```

**Steps**:
1. Implement mention parsing/formatting
2. Implement emoji parsing/formatting
3. Implement snowflake utilities
4. Implement helper functions
5. Tests for all utilities

## 6.3: Performance & Optimization (3 days)

### Task 6.3.1: Connection Pooling
**Complexity**: Medium
**Dependencies**: Phase 3

**Steps**:
1. Implement HTTP connection pooling
2. Configure MaxIdleConns and MaxIdleConnsPerHost
3. Connection reuse metrics
4. Benchmarks
5. Documentation

### Task 6.3.2: Request Batching
**Complexity**: High
**Dependencies**: Phase 3

**Implementation**:
```go
// gosdk/discord/client/batch.go
type Batcher struct {
    client *Client
    queue  chan *request
    wg     sync.WaitGroup
}

func (c *Client) NewBatcher() *Batcher
func (b *Batcher) AddMessage(channelID, content string)
func (b *Batcher) AddReaction(channelID, messageID, emoji string)
func (b *Batcher) Flush(ctx context.Context) error
```

**Steps**:
1. Implement request batching
2. Automatic flushing with timer
3. Respect rate limits
4. Tests
5. Examples

### Task 6.3.3: Caching Strategies
**Complexity**: Medium
**Dependencies**: Phase 5

**Steps**:
1. Implement cache warming
2. Implement cache invalidation strategies
3. LRU cache implementation
4. Cache metrics (hit rate, size)
5. Tests
6. Examples

### Task 6.3.4: Benchmarks & Profiling
**Complexity**: Medium
**Dependencies**: All previous phases

**Steps**:
1. Benchmark critical paths (message send, event dispatch)
2. Memory profiling
3. CPU profiling
4. Optimize hot paths
5. Document performance characteristics
6. Performance regression tests

## 6.4: Error Handling & Resilience (2 days)

### Task 6.4.1: Advanced Error Types
**Complexity**: Medium
**Dependencies**: Phase 1

**Implementation**:
```go
// Extend types/errors.go
type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    state        string
    failures     int
    lastFailure  time.Time
    mu           sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error

type RetryPolicy struct {
    MaxAttempts int
    BackoffBase time.Duration
    BackoffMax  time.Duration
    Jitter      bool
}

func (rp *RetryPolicy) Execute(ctx context.Context, fn func() error) error
```

**Steps**:
1. Implement circuit breaker
2. Implement retry policies
3. Implement timeout policies
4. Integrate with client
5. Tests
6. Examples

### Task 6.4.2: Health Checks
**Complexity**: Low
**Dependencies**: All clients

**Implementation**:
```go
// gosdk/discord/health/health.go
type Checker struct {
    client *client.Client
}

func (h *Checker) CheckAPI(ctx context.Context) error
func (h *Checker) CheckGateway(ctx context.Context) error
func (h *Checker) CheckWebhook(ctx context.Context, webhookURL string) error

type HealthReport struct {
    Timestamp time.Time         `json:"timestamp"`
    Status    string            `json:"status"`
    Checks    map[string]string `json:"checks"`
}

func (h *Checker) Report(ctx context.Context) (*HealthReport, error)
```

**Steps**:
1. Implement health check endpoints
2. Aggregate health status
3. JSON health reports
4. Tests
5. Examples

## 6.5: Testing & Documentation (2 days)

### Task 6.5.1: Final Testing
**Steps**:
1. End-to-end integration tests
2. Performance benchmarks
3. Load testing
4. Security review
5. Achieve >85% overall coverage

### Task 6.5.2: Complete Documentation
**Steps**:
1. Review all godoc
2. Create migration guide (from Python bot)
3. Create best practices guide
4. Create troubleshooting guide
5. Create performance tuning guide
6. Update all examples
7. Final AGENTS.md update

---

# Phase 7: vibe CLI Integration

**Duration**: 1 week
**Focus**: Integration with vibe CLI
**Agent Tasks**: 4-6 atomic tasks

## 7.1: CLI Commands (3 days)

### Task 7.1.1: CLI Structure
**Complexity**: Medium
**Dependencies**: All SDK packages

**Implementation**:
```go
// gosdk/cmd/discord/main.go
package main

import (
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "discord",
        Short: "Discord SDK CLI",
    }

    rootCmd.AddCommand(webhookCmd())
    rootCmd.AddCommand(messageCmd())
    rootCmd.AddCommand(channelCmd())
    rootCmd.AddCommand(guildCmd())
    rootCmd.AddCommand(interactionCmd())

    rootCmd.Execute()
}
```

**Steps**:
1. Set up cobra CLI framework
2. Implement root command
3. Implement webhook commands
4. Implement message commands
5. Implement channel commands
6. Implement guild commands
7. Tests for CLI

### Task 7.1.2: Configuration Integration
**Complexity**: Low
**Dependencies**: Task 7.1.1

**Steps**:
1. Integrate with vibe config system
2. Support config file discovery
3. Support environment variables
4. Support flags
5. Config validation
6. Examples

### Task 7.1.3: Output Formatting
**Complexity**: Low
**Dependencies**: Task 7.1.1

**Implementation**:
```go
// gosdk/cmd/discord/output/formatter.go
type Formatter interface {
    Format(v interface{}) ([]byte, error)
}

type JSONFormatter struct{}
type TableFormatter struct{}
type YAMLFormatter struct{}

func NewFormatter(format string) Formatter
```

**Steps**:
1. Implement JSON formatter
2. Implement table formatter
3. Implement YAML formatter
4. Support --output flag
5. Tests
6. Examples

## 7.2: Integration Examples (2 days)

### Task 7.2.1: Usage Examples
**Steps**:
1. Create webhook integration example
2. Create bot integration example
3. Create slash command integration example
4. Create event listener example
5. Document integration patterns

### Task 7.2.2: Migration Guide
**Steps**:
1. Create migration guide from Python bot
2. Side-by-side code comparisons
3. Common patterns translation
4. Performance comparison
5. Troubleshooting guide

## 7.3: Documentation & Release (2 days)

### Task 7.3.1: Integration Documentation
**Steps**:
1. Create vibe CLI integration guide
2. Document all CLI commands
3. Create configuration guide
4. Create examples repository
5. Update README with vibe integration

### Task 7.3.2: Release Preparation
**Steps**:
1. Version tagging strategy
2. Changelog generation
3. Release notes
4. GitHub releases
5. Package publishing

---

# Agentic Workflow Enhancements

## Cross-Cutting Concerns for All Phases

### 1. Observability
```go
// gosdk/discord/observability/trace.go
type Tracer interface {
    StartSpan(ctx context.Context, name string) (context.Context, Span)
}

type Span interface {
    SetTag(key string, value interface{})
    SetError(err error)
    Finish()
}

// Wrap all operations with tracing
func (c *Client) sendWithTrace(ctx context.Context, ...) {
    ctx, span := c.tracer.StartSpan(ctx, "discord.send_message")
    defer span.Finish()

    span.SetTag("channel_id", channelID)
    // ... operation
    if err != nil {
        span.SetError(err)
    }
}
```

**Apply to**:
- All HTTP requests
- All gateway events
- All rate limit waits
- All cache operations

### 2. Request/Response Recording
```go
// gosdk/discord/recorder/recorder.go
type Recorder interface {
    RecordRequest(ctx context.Context, req *http.Request)
    RecordResponse(ctx context.Context, resp *http.Response)
    RecordEvent(ctx context.Context, event Event)
}

type FileRecorder struct {
    path string
}

// Enable via config or flag
func (c *Client) WithRecorder(r Recorder) *Client
```

**Use cases**:
- Debugging failed operations
- Creating test fixtures
- Audit trails
- Replay for testing

### 3. Dry-Run Mode
```go
// All write operations support dry-run
type Config struct {
    DryRun bool `yaml:"dry_run"`
}

func (c *Client) SendMessage(ctx context.Context, ...) {
    if c.config.DryRun {
        c.logger.Info("DRY RUN: would send message",
            "channel_id", channelID,
            "content", content)
        return nil
    }
    // ... actual send
}
```

**Apply to**:
- All mutations (POST, PATCH, DELETE, PUT)
- Message sends
- Command registration
- Role/permission changes

### 4. Operation Templates
```go
// gosdk/discord/templates/templates.go
type Template struct {
    Name        string
    Description string
    Execute     func(ctx context.Context, params map[string]interface{}) error
}

// Example: Send build notification
var BuildNotification = Template{
    Name: "build_notification",
    Execute: func(ctx context.Context, params map[string]interface{}) error {
        // Template implementation
    },
}
```

**Common templates**:
- Build notifications
- Error alerts
- Status updates
- Report generation

### 5. Declarative Operations
```go
// gosdk/discord/declarative/spec.go
type Spec struct {
    Commands    []CommandSpec    `yaml:"commands"`
    Webhooks    []WebhookSpec    `yaml:"webhooks"`
    Permissions []PermissionSpec `yaml:"permissions"`
}

func (s *Spec) Apply(ctx context.Context, client *client.Client) error {
    // Reconcile desired state with actual state
}

// Usage:
// discord apply -f discord-config.yaml
```

**Benefits**:
- Version control for Discord configuration
- Infrastructure as code
- Reproducible setups
- Easy rollback

### 6. Metrics & Monitoring
```go
// gosdk/discord/metrics/metrics.go
type Collector interface {
    RecordRequest(method, endpoint string, statusCode int, duration time.Duration)
    RecordRateLimit(route string, remaining, limit int)
    RecordCacheHit(resource string)
    RecordCacheMiss(resource string)
    RecordEvent(eventType string)
}

type PrometheusCollector struct {
    // Prometheus metrics
}

func (c *Client) WithMetrics(collector Collector) *Client
```

**Metrics to track**:
- Request latency (p50, p95, p99)
- Rate limit usage
- Error rates
- Cache hit rates
- Event throughput
- Gateway connection status

### 7. Workflow Automation
```go
// gosdk/discord/workflow/workflow.go
type Step struct {
    Name     string
    Action   func(ctx context.Context) error
    OnError  func(ctx context.Context, err error) error
    Retries  int
    Timeout  time.Duration
}

type Workflow struct {
    Name  string
    Steps []Step
}

func (w *Workflow) Execute(ctx context.Context) error {
    // Execute steps with error handling, retries, timeouts
}

// Example workflow:
var DeploymentNotification = Workflow{
    Name: "deployment_notification",
    Steps: []Step{
        {Name: "send_start", Action: sendStartMessage},
        {Name: "wait_deployment", Action: waitForDeployment},
        {Name: "send_result", Action: sendResultMessage},
    },
}
```

---

# Success Criteria

## Per Phase

### Code Quality
- [ ] >80% test coverage (>85% for Phase 6+)
- [ ] All godoc comments complete
- [ ] No linter warnings
- [ ] All examples working

### Functionality
- [ ] All planned features implemented
- [ ] Integration tests passing
- [ ] Real Discord API tests passing (optional)

### Documentation
- [ ] Package documentation complete
- [ ] Usage guides written
- [ ] Examples provided
- [ ] AGENTS.md updated

### Agentic Readiness
- [ ] All operations JSON-loggable
- [ ] Dry-run mode supported
- [ ] Clear error messages
- [ ] Observable behavior

## Overall Project Success

### Technical
- [ ] Complete Discord API coverage
- [ ] <100ms p95 latency for REST operations
- [ ] <1% rate limit errors
- [ ] >90% uptime for gateway connections
- [ ] >85% test coverage overall

### Documentation
- [ ] Complete API documentation
- [ ] Comprehensive guides
- [ ] Working examples for all features
- [ ] Migration guide from Python bot

### Integration
- [ ] Successful vibe CLI integration
- [ ] CLI commands working
- [ ] Configuration compatible
- [ ] Examples running

### Agentic
- [ ] All operations deterministic
- [ ] Full observability
- [ ] Template library available
- [ ] Declarative config supported

---

# Appendix

## Estimated Effort

| Phase | Duration | Complexity | Agent Tasks |
|-------|----------|------------|-------------|
| Phase 1 | 1 day | Low | 8 | ✅ DONE
| Phase 2 | 1 week | Medium | 6-8 |
| Phase 3 | 2 weeks | High | 12-15 |
| Phase 4 | 2 weeks | High | 10-12 |
| Phase 5 | 3 weeks | Very High | 15-18 |
| Phase 6 | 2 weeks | Medium | 8-10 |
| Phase 7 | 1 week | Low | 4-6 |
| **Total** | **~10 weeks** | - | **~70 tasks** |

## Risk Factors

### Technical
- Discord API changes (mitigation: version pinning, monitoring)
- Rate limiting complexity (mitigation: comprehensive testing)
- WebSocket stability (mitigation: robust reconnection)
- Gateway event volume (mitigation: efficient caching)

### Process
- Scope creep (mitigation: strict phase boundaries)
- Test coverage gaps (mitigation: coverage requirements)
- Documentation lag (mitigation: docs-first approach)

## Dependencies

### External
- Go 1.21+
- Discord API access
- Test Discord server
- CI/CD environment

### Internal
- vibe CLI framework (for Phase 7)
- Configuration system
- Logging framework

## Testing Strategy

### Unit Tests
- All packages must have unit tests
- Table-driven tests preferred
- Mock external dependencies
- Target: >80% coverage

### Integration Tests
- Optional tests with real Discord API
- Build tag: `integration`
- Require test credentials
- Run in CI on schedule

### Benchmarks
- Critical paths benchmarked
- Performance regression tests
- Memory profiling

### End-to-End
- Complete workflow tests
- Multi-package integration
- Real-world scenarios

## Review Checkpoints

After each phase:
1. Code review
2. Test coverage review
3. Documentation review
4. Performance review
5. Agentic readiness review
6. Update OPEN_QUESTIONS.md
7. Update STATUS.md
