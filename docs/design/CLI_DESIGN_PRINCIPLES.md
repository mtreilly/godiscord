# CLI Design Principles (Discord SDK)

Source inspiration: `~/vibe-engineering/docs/design/CLI_DESIGN_PRINCIPLES.md` and `../agent-mobile/docs/design/CLI_DESIGN_PRINCIPLES.md`

## Goals
- Fast, ergonomic Discord SDK that is safe to automate and easy to reason about
- Consistent API: same patterns, types, and error semantics across packages
- Agent-friendly: JSON-first, deterministic outputs, excellent logs
- CLI-ready: designed to integrate seamlessly into the vibe CLI

## Package Model
- Clear hierarchy: `gosdk/<feature>/<package>`
  - Example: `gosdk/discord/client`, `gosdk/discord/webhook`, `gosdk/discord/interactions`
- Predictable defaults; require configuration only when essential
- Global patterns reserved (e.g., logging, error handling, context propagation)

## Outputs
- Structured logging with levels (debug, info, warn, error)
- JSON-serializable types for all data structures
- Never mix human and machine output; use proper logging
- Write artifacts to configurable paths (default: `artifacts/`)

## Errors
- Fail fast with actionable error messages (what failed, next step)
- Use typed errors for programmatic handling
- Include context in error chains (wrap errors with additional information)
- Provide short hints for common recovery patterns

## Configuration
- Config files override defaults
- Environment variables override config files
- Explicit parameters override all (precedence: params > env > config > defaults)
- Use explicit targeting with sensible auto-detection fallback
- Timeouts and retries are explicit; log actual values in debug mode

## Performance
- Avoid unnecessary goroutines; use connection pooling for HTTP
- Cache cheap metadata (auth tokens, rate limits) appropriately
- Provide configurable timeouts for all network operations
- Support context cancellation throughout

## Testing
- Unit test pure logic (parsers, formatters, validators)
- Table-driven tests for comprehensive coverage
- Golden tests for JSON serialization/deserialization
- Integration tests against Discord test servers (when available)

## Logging
- Structured logging using standard library or minimal dependencies
- Debug mode surfaces timing, retries, rate limits
- Logs saved with clear levels and context
- Summaries call out next steps and potential improvements

## Extensibility
- Compose small packages; avoid "god" packages
- Keep schemas versioned and documented
- Maintain backward compatibility or provide clear migration paths
- Use interfaces for mockability and testing

## Anti-Patterns
- Blocking operations without context support
- Silent failures or hidden retries without logging
- Inconsistent error types between similar operations
- Global state or singletons (prefer dependency injection)

## Checklist (per package)
- [ ] Clear godoc with examples
- [ ] JSON-serializable types where applicable
- [ ] Proper error wrapping and typing
- [ ] Context support for cancellation/timeouts
- [ ] Comprehensive unit tests (>80% coverage)
- [ ] Integration examples in examples/ directory
