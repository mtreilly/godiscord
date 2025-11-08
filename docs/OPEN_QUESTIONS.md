OPEN_QUESTIONS.md — Discord Go SDK

Purpose
- Track unresolved questions, assumptions, and decisions needed across SDK design, vibe CLI integration, and Discord API patterns. This is a living document for agents to update continuously.

How To Use
- Add an entry whenever you encounter uncertainty, a blocked decision, or a design tradeoff.
- Keep entries concise, include context and a proposed next step or experiment.
- Close questions by linking to commits/PRs/docs that resolved them; keep historical entries for traceability.

Entry Template
```
## [ID]: [Short title]
Scope: [SDK|Integration|Gateway|Testing|Docs|Research]
Owner: [name or 'unassigned']
Last Updated: YYYY-MM-DD

Context
- One paragraph describing the situation and why it matters.

Open Question(s)
- Bullet(s) with the exact question(s) we need answered.

Hypotheses / Options
- Option A: ... (pros/cons)
- Option B: ... (pros/cons)

Proposed Experiment(s)
- Step 1: ...
- Step 2: ...

Signals / Success Criteria
- What we will measure or observe to decide.

Links
- Related code (paths) or PRs
- Related docs (plans/manuals)
- Discord API docs

Status
- open | in_progress | answered (with link)
```

Seed Questions

## Q1: Rate limiting strategy for SDK
Scope: SDK | Owner: unassigned | Last Updated: 2025-11-08

Context
- Discord has per-route rate limits and global rate limits. SDK needs to handle these gracefully without exposing complexity to vibe CLI users. We need to decide between client-side tracking vs reactive handling.

Open Question(s)
- Should we implement proactive rate limit tracking or reactive backoff?
- How do we expose rate limit information to CLI users?
- Should we queue requests or return errors immediately?

Hypotheses / Options
- A) Reactive: Respect 429 responses, exponential backoff, retry automatically (simpler, robust)
- B) Proactive: Track rate limits from headers, queue requests, prevent 429s (complex, optimal)
- C) Hybrid: Reactive with optional proactive mode via config flag

2025-11-08 Update
- Implemented route-alias-aware tracker so proactive/adaptive strategies share bucket state across endpoints (`gosdk/ratelimit/tracker.go`).
- Added `docs/guides/RATE_LIMITS.md` to document strategy/telemetry knobs exposed via config/env.
- Need to validate real-world alias coverage once webhook + bot clients share trackers; keep telemetry hooks in Task 2.2.3.

Proposed Experiment(s)
- Implement reactive approach first in webhook package
- Measure 429 frequency in real usage (vibe CLI integration)
- Add proactive tracking if 429s are frequent (>5% of requests)

Signals / Success Criteria
- <1% of requests result in unrecoverable rate limit errors
- Webhook sends complete within 2x theoretical minimum time
- Users don't need to understand rate limits to use SDK

Links
- Discord rate limits: https://discord.com/developers/docs/topics/rate-limits
- gosdk/discord/client (when implemented)

Status
- open

## Q2: Configuration management approach
Scope: SDK|Integration | Owner: unassigned | Last Updated: 2025-11-08

Context
- SDK needs configuration (tokens, timeouts, retry counts). vibe CLI has its own config system. We need to decide how gosdk config integrates with vibe's config without tight coupling.

Open Question(s)
- Should gosdk have its own config package or rely on vibe's?
- How do we handle config precedence (env vars, config files, params)?
- Should we support multiple config formats (YAML, JSON, TOML)?

Hypotheses / Options
- A) SDK has minimal config package, vibe CLI maps vibe config → SDK options (loose coupling, flexible)
- B) SDK uses vibe config directly via import (tight coupling, simpler for vibe users)
- C) SDK accepts functional options only, no config files (code-first, explicit)

Proposed Experiment(s)
- Implement option A: gosdk/config with YAML support + functional options
- Create example showing vibe CLI integration with config mapping
- Evaluate ergonomics and maintainability

Signals / Success Criteria
- Clear separation between SDK and vibe CLI concerns
- Configuration is easy to understand and override
- No config duplication or drift between SDK and vibe CLI

Links
- gosdk/config (to be implemented)
- vibe CLI config system (check ~/vibe-engineering)

Status
- open

## Q3: Testing strategy for Discord API interactions
Scope: Testing | Owner: unassigned | Last Updated: 2025-11-08

Context
- Discord SDK makes HTTP requests to Discord API. We need comprehensive testing without hitting real API in CI. Options include mocking, recording, or test Discord servers.

Open Question(s)
- Should we use HTTP mocking, recorded interactions, or real test Discord server?
- How do we test rate limiting and error scenarios?
- What coverage target is appropriate for SDK code?

Hypotheses / Options
- A) Mock HTTP client: Fast, deterministic, requires manual scenario crafting (good for unit tests)
- B) Record/replay: Real responses, requires one-time recording, can drift from API (good for integration)
- C) Test Discord server: Most realistic, requires setup, may have rate limits (good for end-to-end)
- D) Combination: Mock for unit tests, test server for integration tests

Proposed Experiment(s)
- Implement HTTP client interface for mockability
- Create mock client for unit tests in webhook package
- Evaluate need for integration tests as SDK matures

Signals / Success Criteria
- >80% test coverage on core packages
- All error paths are testable
- Tests run in <10s without network dependencies
- Integration tests (if needed) are optional and clearly marked

Links
- gosdk/discord/webhook/*_test.go (to be implemented)
- httptest package: https://pkg.go.dev/net/http/httptest

Status
- open

## Q4: Gateway (WebSocket) implementation priority
Scope: Gateway | Owner: unassigned | Last Updated: 2025-11-08

Context
- Discord Gateway provides real-time events via WebSocket. Webhooks and REST API cover many use cases. Gateway is complex and may not be needed for initial vibe CLI integration.

Open Question(s)
- Is Gateway support required for v1.0 of the SDK?
- What use cases in vibe CLI would require Gateway?
- Should we design SDK structure anticipating Gateway, even if not implemented yet?

Hypotheses / Options
- A) Defer Gateway to post-v1.0, focus on webhooks + REST API (faster delivery, simpler)
- B) Implement basic Gateway in v1.0 for event-driven use cases (more complete, slower)
- C) Design package structure for Gateway, implement stub/placeholder (future-proof, minimal cost)

Proposed Experiment(s)
- Survey vibe CLI use cases: which need real-time events vs request/response?
- Estimate Gateway implementation complexity (time + testing burden)
- Design gosdk/discord/gateway package structure without implementation

Signals / Success Criteria
- Clear understanding of vibe CLI Gateway requirements
- Package structure supports Gateway addition without breaking changes
- Decision documented and communicated to vibe CLI team

Links
- Discord Gateway docs: https://discord.com/developers/docs/topics/gateway
- Old Python bot implementation: discord-bot/coordinator/discord/bot.py
- docs/plans/ROADMAP.md (to be created)

Status
- open

## Q5: Error handling and typed errors
Scope: SDK | Owner: unassigned | Last Updated: 2025-11-08

Context
- SDK will encounter various errors: network failures, API errors (4xx, 5xx), rate limits, invalid tokens, etc. We need consistent error handling that allows vibe CLI to handle errors programmatically.

Open Question(s)
- Should we use sentinel errors, error types, or both?
- How granular should error types be (per API endpoint vs per error category)?
- Should errors be wrapped with context using %w or custom error types?

Hypotheses / Options
- A) Sentinel errors (errors.New) + wrapping: Simple, idiomatic, limited type safety
- B) Custom error types with fields: Rich context, type-safe, more boilerplate
- C) Combination: Custom types for major categories, wrapping for context

Proposed Experiment(s)
- Define error categories: RateLimitError, AuthError, NetworkError, ValidationError, APIError
- Implement custom error types in gosdk/discord/types/errors.go
- Use errors.Is and errors.As in examples to demonstrate handling

Signals / Success Criteria
- Errors are easy to handle programmatically (errors.Is, errors.As work)
- Error messages are actionable for developers
- Vibe CLI can distinguish rate limits from auth errors from network failures

Links
- Go errors package: https://pkg.go.dev/errors
- gosdk/discord/types/errors.go (to be implemented)

Status
- open
