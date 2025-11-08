# Discord Bot (extracted)

Standalone Discord bot extracted from the historical `agent-coordinator` module.

## Setup

- Copy `.env.example` to `.env` and set `DISCORD_BOT_TOKEN`.
- Optionally set `DISCORD_CHANNEL_ID` in `.env` and update `config.yaml`.

## Install and Run

```
cd discord-bot
uv sync
uv run python -m coordinator.discord.bot

# or via script entrypoint
uv run discord-bot
```

## Files

- `coordinator/core/config.py` – loads `.env` and `config.yaml`
- `coordinator/discord/bot.py` – slash commands and connection
- `coordinator/discord/formatters.py` – message formatting helpers
- `coordinator/discord/reporter.py` – one-shot reporting helpers
- `coordinator/tasks/queue.py` – JSON-backed task queue
- `coordinator/registries/agent_registry.py` – reads agent/model registries
- `registries/*.yaml` – agents and models

