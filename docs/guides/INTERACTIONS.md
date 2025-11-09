# Discord Interactions Guide

This guide collects the current best practices for building slash commands, components, modals, and the HTTP interaction server that powers them. Follow it for agent-ready handlers, predictable responses, and tests that keep the SDK reliable.

## Slash Commands

1. Define your commands via the fluent builders in `gosdk/discord/interactions/builder.go`. The builders enforce the Discord limits on names, descriptions, options, and choice counts before you call `Build()`.
2. Register handlers through the interaction server:
   ```go
   server.RegisterCommand("hello", func(ctx context.Context, i *types.Interaction) (*types.InteractionResponse, error) {
       // Build responses via `NewMessageResponse` or `NewDeferredResponse`
       return interactions.NewMessageResponse("world").Build()
   })
   ```
3. Commands are normalized to lowercase before matching, and validation occurs in builders so deployment time errors are rare.
4. Use `go test ./discord/interactions` frequently—tests already cover command routing, middleware order, and error paths to make sure your handlers behave deterministically.

## Components & Responses

- Components are always returned from an action row (`types.ComponentTypeActionRow`). Use the typed builders (`NewButton`, `NewActionRow`, `NewStringSelect`, etc.) to keep your payloads valid.
- Response builders live in `gosdk/discord/interactions/response_builder.go`. Common workflows look like this:
  ```go
  response := interactions.NewMessageResponse("Done!").
      SetEphemeral(true).
      AddComponentRow(row).
      Build()
  ```
- Builders validate that you only add action rows at the top level and only text inputs when building a modal. Our unit tests assert these guards (`response_builder_test.go`).
- Call `SetComponents` or `SetModalComponents` when you need to replace rows, and rely on the helpers to convert the typed components into the raw `types.MessageComponent` structure.

## Modals

1. Create modal responses with `NewModalResponse(customID, title)` and compose inputs via `NewTextInput`:
   ```go
   modal := interactions.NewModalResponse("modal_submit", "Feedback")
   modal.SetModalComponents(
       interactions.NewActionRow().AddComponent(
           interactions.NewTextInput("feedback", "Tell us", types.TextInputStyleParagraph).Build(),
       ).Build(),
   )
   ```
2. Modals validate title length, custom ID length, and text input constraints (length, placeholder) before the request is sent.
3. Modal submissions land as `types.InteractionTypeModalSubmit` on the server. Register them with `RegisterModal` and build a response just like a command.

## Interaction Server

- Start the server with your application's Ed25519 public key and optionally inject a custom router or logger:
  ```go
  server, err := interactions.NewServer(pubKey, interactions.WithRouter(router), interactions.WithDryRun(false))
  ```
- `HandleInteraction` automatically checks HTTP method, verifies the Discord signature, and routes the payload. Pings reply with a `PONG`, and unknown interactions return `404`.
- You can register component handlers and middleware via `RegisterComponent`, `RegisterModal`, or by using `NewRouter()` to handle regex patterns and shared middleware chains. Middleware order is preserved and tested (`server_test.go`).
- When `dryRun` is enabled, the server skips signature verification—handy for local dev but never enable it in production.

## Testing & Troubleshooting

- The package ships with router/server integration tests (`server_test.go`) plus router-specific coverage (`router_test.go`). Run `go test ./discord/interactions` to exercise all scenarios.
- Use `examples/` (future) to prototype slash commands, then copy the builders into production handlers for consistent behavior.
- If you see `401 Unauthorized`, verify the timestamp/signature headers are forwarded by your reverse proxy; the tests simulate signed requests with `newSignedRequest` for reference.

## References
- See `docs/guides/WEBHOOKS.md` for webhook flow comparisons and `docs/guides/RATE_LIMITS.md` for adapting these handlers under load.
