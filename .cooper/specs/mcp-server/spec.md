# MCP Server Capability Specification

## Purpose
Provides an embedded Model Context Protocol (MCP) server over standard I/O (JSON-RPC 2.0) powered by `twoBoots/bender/pkg/mcp` to allow AI agents to directly query, inspect, and register whitelist policies during pairing sessions.

## Requirements

### Requirement: JSON-RPC 2.0 stdio MCP Engine
The MCP server MUST serve protocol-compliant requests over `os.Stdin` and `os.Stdout`.

#### Scenario: Agent Initializes MCP Session
- **GIVEN** Antigravity or Claude initiates an MCP handshake via `antimatter mcp`
- **WHEN** client sends `initialize` JSON-RPC request
- **THEN** server responds with protocol capabilities, server info, and registered tools.

### Requirement: Dynamic Whitelist MCP Tools
The MCP server MUST expose tools for inspecting and updating whitelist policies:
- `antimatter_get_whitelist`: Return active whitelisted commands and tools.
- `antimatter_add_whitelist_rule`: Register an approved command pattern or tool call.
- `antimatter_sync_status`: Inspect sync status and remote branch drift.

#### Scenario: AI Agent Adds Command Whitelist via MCP Tool Call
- **GIVEN** an active MCP session with Antigravity or Claude
- **WHEN** agent calls `antimatter_add_whitelist_rule` with command pattern
- **THEN** rule is validated, appended to policy store, and synced to target agent adapters.
