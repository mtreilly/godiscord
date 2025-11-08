# Implementation Plan - Quick Reference

**Full Plan**: [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)

## Phase Overview

```
Phase 1: Foundation [1 day] ✅ COMPLETE
├── Project structure
├── Core types & errors
├── Basic webhook client
├── Config & logging
└── Initial tests

Phase 2: Enhanced Webhook & Rate Limiting [1 week]
├── 2.1: Webhook file uploads (2 days)
├── 2.2: Advanced rate limiting (3 days)
├── 2.3: Webhook thread support (1 day)
└── 2.4: Testing & docs (1 day)

Phase 3: Bot API Client [2 weeks]
├── 3.1: HTTP client foundation (2 days)
├── 3.2: Channel operations (3 days)
├── 3.3: Reaction operations (1 day)
├── 3.4: Guild operations (3 days)
└── 3.5: Testing & docs (2 days)

Phase 4: Slash Commands & Interactions [2 weeks]
├── 4.1: Interaction models (2 days)
├── 4.2: Command registration (2 days)
├── 4.3: Interaction responses (3 days)
├── 4.4: Message components (3 days)
├── 4.5: Interaction server (2 days)
└── 4.6: Testing & docs (2 days)

Phase 5: Gateway (WebSocket) [3 weeks]
├── 5.1: Gateway foundation (1 week)
├── 5.2: Event system (1 week)
├── 5.3: Intents & caching (3 days)
├── 5.4: Sharding (3 days)
└── 5.5: Testing & docs (2 days)

Phase 6: Advanced Features & Polish [2 weeks]
├── 6.1: Permission system (3 days)
├── 6.2: Embed builder & utilities (2 days)
├── 6.3: Performance & optimization (3 days)
├── 6.4: Error handling & resilience (2 days)
└── 6.5: Testing & docs (2 days)

Phase 7: vibe CLI Integration [1 week]
├── 7.1: CLI commands (3 days)
├── 7.2: Integration examples (2 days)
└── 7.3: Documentation & release (2 days)

Total: ~10 weeks, ~70 atomic tasks
```

## Next Actions (Phase 2)

### 2.1.1: Multipart Form Support [2 days]
**Priority**: High
**File**: `gosdk/discord/webhook/multipart.go`

**Tasks**:
1. Create FileAttachment type
2. Implement SendWithFiles method
3. Handle multipart/form-data encoding
4. Validate file size limits
5. Write tests with mock files
6. Add example with image upload

**Entry Point**:
```go
type FileAttachment struct {
    Name        string
    ContentType string
    Reader      io.Reader
}

func (c *Client) SendWithFiles(ctx context.Context, msg *WebhookMessage, files []FileAttachment) error
```

## Agentic Workflow Features

Every phase includes:
- ✅ **Observability**: Request/response tracing
- ✅ **Recording**: Capture operations for replay
- ✅ **Dry-run**: Validate without executing
- ✅ **Templates**: Common operation patterns
- ✅ **Declarative**: YAML-based configuration
- ✅ **Metrics**: Performance monitoring

## Package Architecture

```
gosdk/
├── discord/
│   ├── types/          # Core types, errors, models
│   ├── webhook/        # Webhook client [PHASE 1-2]
│   ├── client/         # Bot API client [PHASE 3]
│   ├── interactions/   # Slash commands [PHASE 4]
│   ├── gateway/        # WebSocket events [PHASE 5]
│   ├── permissions/    # Permission system [PHASE 6]
│   ├── embeds/         # Embed builder [PHASE 6]
│   └── utils/          # Utilities [PHASE 6]
├── config/             # Configuration [PHASE 1]
├── logger/             # Structured logging [PHASE 1]
├── ratelimit/          # Rate limiting [PHASE 2]
├── observability/      # Tracing [PHASE 6]
├── recorder/           # Request recording [PHASE 6]
├── templates/          # Operation templates [PHASE 6]
├── declarative/        # Declarative specs [PHASE 6]
├── metrics/            # Metrics collection [PHASE 6]
├── workflow/           # Workflow automation [PHASE 6]
├── cmd/discord/        # CLI commands [PHASE 7]
└── examples/           # Usage examples [ALL PHASES]
```

## Success Metrics

| Metric | Target | Phase |
|--------|--------|-------|
| Test Coverage | >80% | 2-5 |
| Test Coverage | >85% | 6-7 |
| API Latency (p95) | <100ms | 6 |
| Rate Limit Errors | <1% | 2 |
| Gateway Uptime | >90% | 5 |
| Documentation | 100% | 7 |

## Decision Points

| After Phase | Decision | Document |
|-------------|----------|----------|
| 2 | Rate limiting strategy effectiveness | OPEN_QUESTIONS.md |
| 3 | Bot API completeness | OPEN_QUESTIONS.md |
| 4 | Component patterns | OPEN_QUESTIONS.md |
| 5 | Gateway priority for vibe CLI | ROADMAP.md |
| 6 | Performance targets met | STATUS.md |
| 7 | v1.0 readiness | STATUS.md |

## Common Patterns

### Task Structure
Each task follows:
1. **Define types** → Write structs, interfaces, constants
2. **Implement logic** → Core functionality
3. **Add validation** → Input checking, error handling
4. **Write tests** → Unit tests, table-driven tests
5. **Create examples** → Usage demonstrations
6. **Document** → Godoc, guides, AGENTS.md update

### Testing Approach
- **Unit tests**: Mock external dependencies
- **Integration tests**: Optional, build tag `integration`
- **Golden tests**: JSON serialization
- **Benchmarks**: Performance-critical paths

### Error Handling
- Define typed errors in `types/errors.go`
- Wrap errors with context using `fmt.Errorf(..., %w, err)`
- Use `errors.Is` and `errors.As` for checking
- Provide actionable error messages

### Logging Pattern
```go
logger.Info("operation starting",
    "param1", value1,
    "param2", value2,
)

// ... operation

if err != nil {
    logger.Error("operation failed",
        "error", err,
        "param1", value1,
        "elapsed_ms", elapsed.Milliseconds(),
    )
}
```

## Quick Start for Agents

### Starting a New Task
1. Read full task description in IMPLEMENTATION_PLAN.md
2. Check dependencies are complete
3. Create feature branch: `feat/phase-X-task-Y-Z`
4. Implement following pattern above
5. Run tests: `go test -v ./...`
6. Update STATUS.md
7. Commit with conventional format: `feat(package): description`

### Blocked on Question?
1. Add entry to OPEN_QUESTIONS.md
2. Tag with phase and package
3. Propose experiment or options
4. Continue with next independent task

### Ready for Review?
1. Verify >80% coverage: `go test -cover ./...`
2. Verify no lint errors: `go vet ./...`
3. Verify formatting: `gofmt -l .`
4. Update documentation
5. Create PR with checklist from IMPLEMENTATION_PLAN.md

## Resources

- **Full Implementation Plan**: [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)
- **Development Roadmap**: [ROADMAP.md](ROADMAP.md)
- **Current Status**: [../progress/STATUS.md](../progress/STATUS.md)
- **Design Principles**: [../design/CLI_DESIGN_PRINCIPLES.md](../design/CLI_DESIGN_PRINCIPLES.md)
- **Patterns Cookbook**: [../design/CLI_PATTERNS_COOKBOOK.md](../design/CLI_PATTERNS_COOKBOOK.md)
- **Open Questions**: [../OPEN_QUESTIONS.md](../OPEN_QUESTIONS.md)
- **Agent Guide**: [../../AGENTS.md](../../AGENTS.md)
