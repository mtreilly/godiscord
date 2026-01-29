# Design Documentation Index

This directory contains design principles and patterns for the Discord Go SDK.

## Documents

### [DESIGN_PRINCIPLES.md](DESIGN_PRINCIPLES.md)
Core design principles for building the Discord SDK:
- Goals and philosophy
- Package organization
- Error handling patterns
- Configuration management
- Testing strategy
- Performance considerations
- Anti-patterns to avoid

### [PATTERNS_COOKBOOK.md](PATTERNS_COOKBOOK.md)
Practical patterns and recipes:
- Package organization examples
- Client initialization patterns
- Error handling examples
- Context propagation
- Configuration file formats
- JSON serialization
- Retry and rate limiting logic
- Testing patterns
- Usage examples

## How to Use

1. **Starting a new package?** Read DESIGN_PRINCIPLES.md first
2. **Implementing a feature?** Check PATTERNS_COOKBOOK.md for examples
3. **Unsure about a design decision?** Add to [../OPEN_QUESTIONS.md](../OPEN_QUESTIONS.md)

## Principles Summary

- **Context everywhere**: All operations accept `context.Context`
- **Error wrapping**: Use `fmt.Errorf` with `%w` to wrap errors
- **Typed errors**: Define sentinel errors and error types for programmatic handling
- **Interfaces**: Use interfaces for testability and dependency injection
- **Table-driven tests**: Comprehensive test coverage with clear test cases
- **No global state**: Avoid global variables and singletons
- **JSON-first**: All data structures should be JSON-serializable
- **Godoc comments**: Document all exported symbols

## References

- Discord API: https://discord.com/developers/docs
- Go best practices: https://go.dev/doc/effective_go
