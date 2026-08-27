# CLI Capability Specification

## Purpose
Provides a standardized, developer-friendly Command-Line Interface (CLI) for managing, inspecting, and synchronizing Antigravity and Claude Code agent policies and configurations.

## Requirements

### Requirement: Command Hierarchy & Subcommands
The CLI MUST expose structured subcommands for core operations:
- `antimatter sync`: Perform bi-directional or directional synchronization between local agent configs and the sync store.
- `antimatter push`: Push local whitelist rules and settings to the remote sync backend.
- `antimatter pull`: Pull remote whitelist rules and apply to local agent settings.
- `antimatter diff`: Compare active agent configurations with the canonical sync store.
- `antimatter whitelist`: Inspect, add, remove, and validate whitelisted tool calls and CLI execution commands.
- `antimatter mcp`: Launch stdio Model Context Protocol (MCP) server or install MCP client configs.
- `antimatter update`: In-place self-updater to latest GitHub release.

#### Scenario: Display Help and Hierarchy
- **GIVEN** `antimatter` binary is installed
- **WHEN** user executes `antimatter --help` or `antimatter <subcommand> --help`
- **THEN** it outputs structured usage, flags, available subcommands, and examples with exit code `0`.

#### Scenario: Add Whitelist Rule
- **GIVEN** a command pattern (e.g. `git status` or `npm test`)
- **WHEN** user executes `antimatter whitelist add --command "git status"`
- **THEN** the rule is validated, recorded in canonical format, and reflected in Antigravity and Claude adapter targets.
