# Migration from Python Bot

This migration guide helps port workflows from the legacy `discord-bot` Python implementation to the new Go SDK.

## Configuration

- Python bot used `discord-bot/config.yaml`. The Go CLI reads similar YAML files via `gosdk/config` and supports `DISCORD_*` env overrides plus CLI flags.
- Mirror webhook URLs, bot tokens, and logging settings directly; the Go config loads defaults (`30s timeout`, `adaptive rate limits`).

## Webhooks

- Python `bot/webhook.py` posted JSON payloads with attachments. Use `gosdk/discord/webhook` clients or the CLI `discord webhook --webhook ...` to send the same embeds/files.
- Structured logging (JSON) and retry policies are built into the Go client while Python relied on manual try/except blocks.

## Slash Commands & Interactions

- Python command registration is replaced by `gosdk/discord/interactions` builders (`NewSlashCommand`, `NewMessageResponse`, etc.).
- Implement handlers using CLI command context to access configs, and respond via `interaction` command templates before wiring real HTTP handlers.

## Gateway & Events

- Python gateway listeners used `discord.py`. The Go counterpart is `gosdk/discord/gateway` with retries, caching, and sharding baked in.
- Use the CLI to demonstrate event flows before deploying the full gateway client for production (see `docs/guides/GATEWAY.md`).

## Testing & Resilience

- Python unit tests live under `discord-bot/tests`. Go tests moved to `gosdk/..._test.go` files and rely on `go test ./...`.
- Retry/circuit breaker behaviors now live in `gosdk/discord/types/resilience.go` and health checks in `gosdk/discord/health`. Document them here so your migration covers resilience semantics.
