# vibe CLI Integration

This guide explains how the `discord` CLI integrates with vibe’s configuration system, how commands map to SDK packages, and what artifacts to ship during release.

## CLI Overview

- `discord webhook` — posts via configured webhook (uses `gosdk/discord/webhook`, respects `config.Discord.Webhooks`).
- `discord message` — placeholder for bot message workflows (future: wrap `gosdk/discord/client` message send).
- `discord channel` — shows channel/webhook metadata and can expand into channel operations.
- `discord guild` — surfaces application ID and can tie into guild CRUD helpers.
- `discord interaction` — demonstrates interaction configs and will evolve into slash command and component workflows.

Each command shares the loaded `config.Config` (auto-discovered from `discord-config.yaml`, `config/discord.yaml`, or env vars) and writes output via `--output` formatters (JSON/YAML/Table).

## Configuration and Flags

| Flag | Purpose |
|------|---------|
| `--config` | Path to YAML config file. |
| `--token` | Override bot token from env/config. |
| `--webhook` | Override default webhook URL. |
| `--output` | Choose `json`, `yaml`, or `table` output. |

## Output Formats

- `json`: Indented JSON (default). Good for automation.
- `yaml`: YAML documents for config-heavy tools.
- `table`: Human-friendly key/value table.

## Release Considerations

1. Update `AGENTS.md` and docs when new commands arrive.
2. Capture CLI usage examples for pipeline docs under `docs/guides/CLI_EXAMPLES.md`.
3. Tag releases with `v0.n.m` after bundling CLI + SDK changes; include changelog references to Phase 7 docs.
