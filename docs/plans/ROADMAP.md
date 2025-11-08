# Discord Go SDK Roadmap

## Vision

Build a production-ready Go SDK for Discord interactions that integrates seamlessly with the vibe CLI, following best practices and agent-friendly design patterns.

## Goals

1. **Developer Experience**: Easy to use, well-documented, idiomatic Go
2. **Reliability**: Comprehensive error handling, retries, rate limiting
3. **Testing**: >80% test coverage with integration and unit tests
4. **Integration**: Seamless vibe CLI integration with minimal friction
5. **Maintainability**: Clear architecture, good documentation, open design process

## Phases

### Phase 1: Foundation ✅ COMPLETE (2025-11-08)

**Goal**: Establish project structure, core types, and basic webhook functionality

Deliverables:
- [x] Project documentation (AGENTS.md, design docs, OPEN_QUESTIONS.md)
- [x] Core types package (Message, Embed, User, errors)
- [x] Webhook client with retry logic
- [x] Configuration management
- [x] Structured logging
- [x] Basic tests and examples
- [x] Project README and setup files

**Duration**: 1 day
**Status**: ✅ Complete

---

### Phase 2: Core Features (Next)

**Goal**: Expand webhook functionality and implement bot API client

Deliverables:
- [ ] Full webhook API
  - File uploads and attachments
  - Thread support
  - Edit/delete webhook messages
  - Webhook info retrieval
- [ ] Bot API client package
  - Channel operations (get, create, modify, delete)
  - Message CRUD (create, get, edit, delete)
  - Guild information
  - User information
- [ ] Enhanced rate limiting
  - Per-route rate limit tracking
  - Global rate limit handling
  - Configurable strategies (reactive, proactive, hybrid)
- [ ] Comprehensive error handling
  - More granular error types
  - Better error messages
  - Recovery suggestions
- [ ] Expanded test coverage
  - Target: >80% coverage
  - Integration test examples
  - Golden tests for JSON

**Duration**: 1-2 weeks
**Status**: Planned

---

### Phase 3: Advanced Features

**Goal**: Implement slash commands and component interactions

Deliverables:
- [ ] Slash commands
  - Command registration (global, guild-specific)
  - Command updating and deletion
  - Interaction handling
  - Response types (immediate, deferred, ephemeral)
  - Autocomplete support
- [ ] Component interactions
  - Buttons (primary, secondary, success, danger, link)
  - Select menus (string, user, role, channel, mentionable)
  - Modals and text inputs
  - Component state management
- [ ] Embed builder
  - Fluent API for building embeds
  - Validation and limits
  - Templates for common use cases
  - Color helpers
- [ ] Permissions handling
  - Permission calculation
  - Permission checks
  - Role management

**Duration**: 2-3 weeks
**Status**: Planned

---

### Phase 4: Integration & Polish

**Goal**: Finalize vibe CLI integration and prepare for v1.0 release

Deliverables:
- [ ] vibe CLI integration
  - Integration guide
  - Config mapping examples
  - Command implementations
  - Error handling patterns
- [ ] Documentation completion
  - Complete godoc for all packages
  - Usage guides
  - Migration guide (from old Python bot)
  - Troubleshooting guide
- [ ] Performance optimization
  - Benchmarks for critical paths
  - HTTP connection pooling
  - Caching strategies
  - Memory optimization
- [ ] API stability review
  - Review all public APIs
  - Breaking change assessment
  - Deprecation strategy
  - Versioning plan
- [ ] Examples and templates
  - Common use case examples
  - Best practice templates
  - Anti-pattern warnings

**Duration**: 2-3 weeks
**Status**: Planned

---

### Phase 5: Gateway (Future)

**Goal**: Add WebSocket gateway support for real-time events

Deliverables:
- [ ] Gateway connection
  - WebSocket connection management
  - Heartbeat and reconnection
  - Session resumption
  - Compression support
- [ ] Event handling
  - Event dispatcher
  - Event filtering
  - Custom event handlers
  - Event replay/logging
- [ ] Presence management
  - Set/update presence
  - Status (online, idle, dnd, offline)
  - Activities (playing, streaming, listening, watching)
- [ ] Sharding (if needed)
  - Shard management
  - Shard identification
  - Load balancing
- [ ] Voice support (optional)
  - Voice connection
  - Audio streaming
  - Voice state management

**Duration**: 3-4 weeks
**Status**: Future consideration

---

## Milestones

### v0.1.0 - Foundation (DONE ✅)
- Core types
- Basic webhook client
- Configuration and logging
- Basic tests

### v0.2.0 - Core Features (Next)
- Full webhook API
- Bot API client
- Enhanced rate limiting
- >80% test coverage

### v0.3.0 - Advanced Features
- Slash commands
- Component interactions
- Embed builder
- Permissions

### v0.4.0 - Integration Ready
- vibe CLI integration
- Complete documentation
- Performance optimization
- API stability

### v1.0.0 - Production Ready
- Stable API
- Comprehensive docs
- High test coverage
- vibe CLI integrated

### v2.0.0 - Gateway (Future)
- WebSocket gateway
- Event handling
- Presence management

## Success Metrics

- **Code Quality**: >80% test coverage, no critical lint issues
- **Documentation**: All public APIs documented, comprehensive guides
- **Performance**: <100ms p95 for webhook sends, <50ms for message operations
- **Reliability**: <1% rate limit errors, automatic retry success >95%
- **Adoption**: Integrated into vibe CLI, positive agent feedback

## Dependencies

- Go 1.21+ (for latest stdlib improvements)
- gopkg.in/yaml.v3 (YAML config parsing)
- Potential additions:
  - gorilla/websocket (for Gateway)
  - go-chi/chi (if building example HTTP server)

## Decision Points

1. **After Phase 2**: Evaluate rate limiting strategy effectiveness
2. **After Phase 3**: Decide on Gateway implementation priority based on vibe CLI needs
3. **After Phase 4**: Determine v1.0 readiness and API stability
4. **Ongoing**: Track OPEN_QUESTIONS.md and resolve design decisions

## Notes

- Phases may overlap based on agent availability and priorities
- vibe CLI integration may influence feature prioritization
- Gateway support can be deferred if REST API + webhooks suffice
- Regular sync with vibe CLI team to ensure alignment

## References

- Discord API: https://discord.com/developers/docs
- Old Python bot: `discord-bot/` (reference only)
- Design docs: `../design/`
- Open questions: `../OPEN_QUESTIONS.md`
