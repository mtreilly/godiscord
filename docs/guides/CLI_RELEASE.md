# CLI Release Playbook

Before shipping Phase 7, follow these steps:

1. **Tagging strategy** – use `git tag v0.x.y` for releases that bundle CLI + SDK improvements. Include the phase name (`phase7`) in the tag message or changelog entry.
2. **Changelog generation** – summarize new commands (`discord webhook`, `message`, etc.), config discovery, formatter options, and health/resilience coverage. Append to `docs/plans/QUICK_REFERENCE.md` or a dedicated `CHANGELOG.md`.
3. **Release notes** – explain configuration options, highlight CLI examples from `docs/guides/CLI_EXAMPLES.md`, and mention the migration guide or Phase 6 docs for context.
4. **Packaging** – ensure `gosdk/cmd/discord` builds cleanly (`go build ./cmd/discord`). Bundle binaries along with examples if needed.
5. **Publishing** – update `README.md` and `AGENTS.md` release references, then push tags and changelog to the repo. Automate tarball generation if required by vibe release process.

Keep this playbook updated whenever CLI commands or config patterns change.
