# Webhook Guide

This guide walks through the Go SDK’s webhook features end-to-end so agents can ship production-ready Discord integrations quickly.

## 1. Configure Credentials

```bash
export DISCORD_WEBHOOK="https://discord.com/api/webhooks/ID/TOKEN"
export DISCORD_RATE_LIMIT_STRATEGY=adaptive   # optional (reactive|proactive|adaptive)
export DISCORD_WEBHOOK_THREAD_ID=1234567890   # optional, for thread sends
```

For persistent environments, mirror these values in `.env` or `config.yaml`:

```yaml
client:
  timeout: 30s
  retries: 3
  rate_limit:
    strategy: adaptive
    backoff_base: 1s
    backoff_max: 60s
```

## 2. Create a Webhook Client

```go
client, err := webhook.NewClient(
    webhookURL,
    webhook.WithMaxRetries(3),
    webhook.WithTimeout(30*time.Second),
    webhook.WithStrategyName("adaptive"),
)
```

Inject a custom `ratelimit.Tracker` or logger with the existing option helpers when integrating into the vibe CLI.

## 3. Send Messages

```go
msg := &types.WebhookMessage{
    Content: "Build finished!",
    Embeds: []types.Embed{
        {
            Title:       "✅ Success",
            Description: "All tests passed",
        },
    },
}
if err := client.Send(ctx, msg); err != nil {
    // Handle typed API/validation errors
}
```

Use `SendSimple` for quick text messages, or `SendWithFiles` with `FileAttachment` for multipart uploads. Attachments are validated against Discord’s limits and streamed through a counting reader to prevent oversize payloads.

## 4. Work with Threads and Forums

- `SendToThread(ctx, threadID, msg)` routes into an existing thread (set `ThreadID` or provide the parameter).
- `CreateThread(ctx, threadName, msg)` automatically sets `WebhookMessage.ThreadName` and lets Discord create a new forum thread. Validation ensures you never set both `ThreadID` and `ThreadName`.
- Example: `gosdk/examples/webhook-thread`.

## 5. Rate Limiting & Observability

- The default adaptive strategy learns from traffic and minimizes 429s.
- Set `DISCORD_LOG_LEVEL=debug` to see proactive vs reactive waits in logs.
- `docs/guides/RATE_LIMITS.md` covers strategy selection, overrides, and troubleshooting.
- Optional integration test (`go test -tags integration ./discord/webhook`) sends a real webhook when `DISCORD_WEBHOOK` is set.

## 6. Testing Workflow

```
cd gosdk
go test ./...
go test -race ./discord/webhook
go test ./discord/webhook -run Golden   # JSON fixtures
go test -tags integration ./discord/webhook   # requires DISCORD_WEBHOOK
```

Golden tests live under `discord/webhook/testdata/golden` and guard serialization semantics. Benchmarks (`BenchmarkClientSend`) and race tests (`TestClientSendConcurrent`) ensure ongoing performance and safety.

## 7. Common Patterns

| Use Case               | Helper                              |
|------------------------|-------------------------------------|
| Simple message         | `client.SendSimple(ctx, content)`   |
| Complex embeds         | `types.WebhookMessage{Embeds: …}`   |
| File uploads           | `client.SendWithFiles(...)`         |
| Thread updates         | `client.SendToThread`               |
| Forum posts            | `client.CreateThread`               |
| Rate limit override    | `WithStrategyName("proactive")`     |

Keep AGENTS.md handy for broader project conventions, and update `docs/OPEN_QUESTIONS.md` whenever new webhook decisions surface.
