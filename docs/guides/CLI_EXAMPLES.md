# CLI Integration Examples

This guide shows how to wire the `discord` CLI commands into real workflows.

## Webhook Notifications

Use the CLI to send webhook notifications before pushing or after deployments:

```bash
discord webhook --webhook https://discord.com/api/webhooks/... --output table
```

Behind the scenes, the `webhook` command uses `config.Config` (config discovery + flags) to populate `cfg.Discord.Webhooks["default"]`. Use `--output json` to capture structured responses for automation.

## Bot Operations

The `message` command exercises the bot client:

```bash
discord message --token "$DISCORD_BOT_TOKEN" --output json
```

This command prints token metadata but can be extended to call `gosdk/discord/client` operations (send message, create channel, manage guild) by reusing the injected `config.Config`.

## Slash Command Puzzles

`interaction` currently reports configured webhooks; extend it to call `gosdk/discord/interactions` helpers for registration and handling:

```bash
discord interaction --config examples/cli/interaction.yaml
```

Add your own YAML describing slash commands in `gosdk/examples/interaction.yaml` and parse it via the `interactions` package when wiring the command.

## Event Listener

The CLI can start a gateway listener (future work) by building on `gosdk/discord/gateway`:

1. Use `newRootCommand` to load config, then instantiate `gateway.NewClient`.
2. Register `OnMessageCreate` or `OnInteraction` handlers that log to the CLI output formatter.
3. Run the client loop while streaming events through structured logging (JSON by default).

## Integration Patterns

- Share the `config.Config` returned by `loadConfig` across subcommands to keep CLI behavior consistent.
- Use `formatter := output.NewFormatter(--output)` to write JSON, YAML, or table output for automation and human users.
- Document new subcommands in `docs/guides/CLI_EXAMPLES.md` and update `AGENTS.md` quick links.
