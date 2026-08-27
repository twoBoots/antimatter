# Track Proposal: CLI Foundation & Policy Engine Scaffold

## Rationale & Goal
Initialize the core Go project structure for **`antimatter`**, import foundational packages from `github.com/twoBoots/bender`, and establish the baseline command hierarchy alongside the canonical Policy Engine domain models and validation rules using strict Test-Driven Development (TDD).

## User & Developer Benefits
- **Executable CLI Binary**: Developers can run `antimatter --help`, `antimatter version`, `antimatter whitelist`, `antimatter mcp`, and `antimatter update`.
- **Standardized Policy Engine**: A robust, validated schema for whitelisted shell commands and tool calls (exact, glob, regex) with unit tests.
- **Embedded MCP Engine**: Ready-to-use JSON-RPC 2.0 stdio MCP server for agent tool calls.
- **In-Place Self-Updating**: Instant binary update capability via GitHub Releases.

## Scope Boundaries
- **In-Scope**:
  - `go.mod` module setup (`github.com/twoBoots/antimatter`) and dependency vendoring/integration with `github.com/twoBoots/bender`.
  - Cobra command scaffolding (`cmd/root.go`, `cmd/version.go`, `cmd/whitelist.go`, `cmd/sync.go`, `cmd/mcp.go`, `cmd/update.go`).
  - Canonical Policy Engine data structures (`Policy`, `Rule`, `RuleType`, `CommandPattern`, `ToolPattern`) and YAML/JSON serialization.
  - Policy validation and command/tool match evaluation algorithms with >80% test coverage.
- **Out-of-Scope (Next Tracks)**:
  - Full Git remote synchronization network routines (handled in `track_sync_engine`).
  - Direct live file modification of `~/.gemini/antigravity-cli/` and `~/.claude.json` (handled in `track_agent_adapters`).
