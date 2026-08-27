# Adapters Capability Specification

## Purpose
Provides bi-directional translation between Antimatter's canonical policy representation and specific AI coding assistant environments (Google Antigravity, Claude Code, Claude Desktop).

## Requirements

### Requirement: Google Antigravity Adapter
The adapter MUST translate canonical whitelist and configuration rules into Antigravity target stores (`~/.gemini/antigravity-cli/`).

#### Scenario: Export Rules to Antigravity
- **GIVEN** a canonical policy containing whitelisted CLI commands and MCP tools
- **WHEN** the Antigravity adapter compiles and writes the configuration
- **THEN** Antigravity CLI whitelist configs are updated atomically with prior backup creation.

### Requirement: Claude Code / Claude Desktop Adapter
The adapter MUST translate canonical whitelist rules into Claude permission specifications (`~/.claude.json` / Claude settings).

#### Scenario: Export Rules to Claude
- **GIVEN** a canonical policy containing whitelisted shell patterns
- **WHEN** the Claude adapter compiles and writes the configuration
- **THEN** Claude Code allowed commands / tools are updated with accurate syntax.

### Requirement: MCP Tool Configuration Adapter
The adapter MUST register and synchronize shared MCP servers across all supported AI assistants.

#### Scenario: Multi-Client MCP Registration
- **GIVEN** an active MCP server configuration
- **WHEN** adapter installation is triggered
- **THEN** target client configuration files (Antigravity, Claude Desktop, Claude Code, Cursor) are safely merged.
