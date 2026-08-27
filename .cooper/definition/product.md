# Product Definition

## Vision
**Antimatter** is a high-performance Go CLI utility and embedded Model Context Protocol (MCP) server designed to persist, manage, version control, and synchronize developer AI agent configurations—including whitelisted tool calls, approved CLI commands, and MCP settings—across multiple devices and between different AI assistants (Google Antigravity and Claude Code / Claude Desktop).

## Target Audience
- Developers and engineers using Google Antigravity and/or Claude Code across multiple workstations (macOS, Linux, Windows).
- Engineering teams seeking standardized, reproducible, and version-controlled agent permission policies and tool whitelists.

## Core Capabilities & Initial Scope

1. **Canonical Policy & Whitelist Engine (`policy-engine`)**:
   - Portable, human-readable specification (YAML/JSON) for command whitelist rules, glob/regex patterns, and tool permissions.
   - Intelligent rule merging, deduplication, and conflict resolution across devices.

2. **Agent Adapters (`adapters`)**:
   - **Antigravity Adapter**: Target compilation and sync to `~/.gemini/antigravity-cli/` configuration and whitelist stores.
   - **Claude Adapter**: Target compilation and sync to Claude Code (`~/.claude.json`, settings) and Claude Desktop configurations.
   - **MCP Tool Whitelist Adapter**: Synchronize shared Model Context Protocol server tools and permissions.

3. **Storage & Sync Backend (`git-backend`)**:
   - Dotfile / Git-backed synchronization backend (push, pull, diff, status, merge).
   - Local caching and device-specific override support.

4. **Embedded MCP Server Engine (`mcp-server`)**:
   - Stdio JSON-RPC 2.0 MCP server powered by `twoBoots/bender/pkg/mcp`.
   - Allows AI agents (Antigravity, Claude) to query policies and register approved tool calls directly during pairing sessions.

5. **CLI Ergonomics & Self-Updater (`cli`, `updater`)**:
   - Unified Cobra CLI hierarchy (`antimatter sync`, `antimatter push`, `antimatter pull`, `antimatter whitelist`, `antimatter mcp`).
   - Atomic in-place self-updater powered by `twoBoots/bender/pkg/updater`.

## Quality & Non-Functional Goals
- **Safety & Non-Destructive Operation**: Safe atomic writes and automatic configuration backups prior to applying sync changes.
- **Secrets Sanitization**: Automatically strip, sanitize, or mask API tokens and machine-specific secrets prior to synchronizing policies.
- **Strict TDD & Coverage**: Comprehensive test coverage (>80%) across all modules and adapters.
- **Zero-Dependency Static Binary**: Compiled with Go for instant execution across macOS, Linux, and Windows.
