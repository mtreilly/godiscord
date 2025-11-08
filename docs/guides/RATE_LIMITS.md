# Rate Limit Guide

Discord enforces per-route and global API rate limits. The Go SDK ships with a tracker and three strategies so agents can choose deterministic behavior for their workflows.

## Strategies

| Strategy   | Behavior | When to use |
|------------|----------|-------------|
| `reactive` | Waits only after Discord returns `429` | Simple jobs or low traffic |
| `proactive` | Starts waiting when a bucket approaches its limit (threshold + safety margin) | Build systems that need to avoid `429`s entirely |
| `adaptive` (default) | Learns from recent traffic and adjusts thresholds automatically | Long-running agents with unpredictable workloads |

Switch strategies at runtime using either configuration files, environment variables, or webhook options.

```yaml
client:
  rate_limit:
    strategy: adaptive   # reactive | proactive | adaptive
    backoff_base: 1s     # retry backoff floor when Discord sends retry_after
    backoff_max: 60s     # retry backoff ceiling
```

Environment variable override:

```bash
export DISCORD_RATE_LIMIT_STRATEGY=proactive
```

Webhook options override everything else:

```go
client, err := webhook.NewClient(
    webhookURL,
    webhook.WithStrategyName("adaptive"),
)
```

## Tracker Behavior

- `ratelimit.MemoryTracker` stores buckets by Discord's `X-RateLimit-Bucket` and maps every route to that bucket, so concurrent endpoints share the same counters.
- `Client.waitForRateLimit` performs proactive waits (strategy-driven) before falling back to the tracker's blocking `Wait`.
- Structured logs surface every proactive/reactive wait plus `429` warnings, making it easy to trace latency spikes.

## Observability Tips

- Enable debug logging (`DISCORD_LOG_LEVEL=debug`) to see proactive and reactive waits with durations.
- Inspect adaptive stats via `ratelimit.AdaptiveStrategy.GetStats()` if you inject a shared strategy instance.
- Record rate limit events in metrics by wrapping the tracker interface (expose `Update` hooks).

## Troubleshooting

| Symptom | Possible Cause | Fix |
|---------|----------------|-----|
| Frequent `429` warnings | Strategy too aggressive | Switch to `proactive` or `adaptive`, or increase safety margin |
| Requests pause unexpectedly | Tracker carried over an expired bucket | Call `Tracker.Clear()` before re-using clients between test runs |
| Large file uploads fail early | Aggregate attachment limit exceeded | Adjust file sizes or split payloads; total counter enforces 25 MB cap |

Keep this guide close when tuning agents; documenting chosen strategy + rationale in `docs/OPEN_QUESTIONS.md` helps future contributors understand trade-offs.
