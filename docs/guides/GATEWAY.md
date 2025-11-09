# Gateway Guide

This guide documents the Discord gateway stack (`gosdk/discord/gateway`) so agents understand how to connect, consume events, and scale.

## Architecture

- **Connection** (`connection.go`): wraps `websocket` dialing, heartbeat scheduling, sequence tracking, and reconnection helpers. Use it through `Client` or the shard manager; it enforces context timeouts, logs through `logger`, and exposes JSON payload observability.
- **Client** (`client.go`): coordinates a connection, dispatcher, intents, and presence management. The read loop decodes dispatch payloads (`Ready`, `MESSAGE_CREATE`, component interactions) and passes typed events to the dispatcher for handling.
- **Dispatcher** (`dispatcher.go`): thread-safe registry that supports generic handlers plus typed helpers (`OnMessageCreate`, `OnInteraction`). It logs failures and returns aggregated errors so callers can surface multi-handler issues.
- **Intents** (`intents.go`): bitmask helpers (`Intent`, `AllIntents`, `DefaultIntents`, `Has`) that gate which payloads Discord delivers. Use `DefaultIntents()` for bots without privileged access and `AllIntents()` for internal tooling (request `DISCORD_GATEWAY_INTENTS` from env/flags as needed).
- **Cache** (`cache.go`): in-memory TTL-bounded cache for guilds, channels, and members with hit/miss metrics. Inject via helper utilities when building higher-level features that need quick lookups from gateway events.
- **Sharding** (`shard.go`): `ShardManager` spins up multiple clients, shares the dispatcher, and can autoscale via `/gateway/bot` recommendations or fixed strategies. It includes broadcast helpers, config hooks (`WithShardGatewayBotURL`, `WithShardGatewayHTTPClient`), and `AutoScale` logic.

## Getting Started

1. **Create a client**:
   ```go
   client, err := gateway.NewClient(os.Getenv("DISCORD_BOT_TOKEN"), gateway.DefaultIntents())
   if err != nil {
       log.Fatal(err)
   }
   client.OnMessageCreate(func(ctx context.Context, event *gateway.MessageCreateEvent) error {
       fmt.Println("received", event.ID)
       return nil
   })
   client.Connect(context.Background())
   ```
2. **Handle sharding (optional)**:
   ```go
   manager := gateway.NewShardManager(token, recommendedShards, gateway.DefaultIntents())
   manager.OnInteraction(func(ctx context.Context, event *gateway.InteractionCreateEvent) error {
       // reuse interaction response builders from interactions package
       return nil
   })
   manager.AutoScale(context.Background(), guildCount, &gateway.RecommendedSharding{})
   manager.Connect(context.Background())
   ```
3. **Track cache stats** when necessary:
   ```go
   cache := gateway.NewMemoryCache(5 * time.Minute)
   cache.SetGuild(guild)
   stats := cache.Stats()
   fmt.Printf("cache hits %d", stats.GuildHits)
   ```

## Testing & Validation

- Run unit tests: `cd gosdk && go test ./discord/gateway`
- Gateway coverage includes connection heartbeats, dispatcher behaviors, shard manager scaling, and cache expiration.
- Add integration smoke tests later with `//go:build integration` when a Discord bot/token is available; gate them behind env vars such as `DISCORD_GATEWAY_TOKEN`.

## Observability

- Enable structured logging by passing `WithGatewayLogger(logger.New(...))` or via `ShardManager` options; log entries include opcodes, event types, and rate-limit waits.
- Use `cache.Stats()` for live hit/miss metrics and expose them to whichever measurement system the CLI uses.

## Troubleshooting

- `401 Unauthorized` when calling `/gateway/bot`: verify the `Authorization: Bot <token>` header is present (the shard manager adds it automatically).
- Heartbeat timeouts: adjust `WithHeartbeatInterval` when debugging or when Discord reports mismatched values (the client reconfigures when it receives `Hello`).
- Missing events: ensure your intents include the categories you expect (`IntentGuildMessages`, `IntentMessageContent`, etc.).

## References

- `docs/plans/IMPLEMENTATION_PLAN.md` Phase 5 for roadmap.
- `docs/guides/INTERACTIONS.md` for building slash commands and component handlers that consume gateway events.
