# Policy Engine Capability Specification

## Purpose
Defines the canonical schema, validation, merging, and conflict resolution rules for agent tool permissions and command-line execution whitelists.

## Requirements

### Requirement: Canonical Whitelist Schema
The policy engine MUST define a portable specification (YAML/JSON) supporting:
- Command whitelist entries (exact match, glob patterns, regexes).
- Tool call whitelist entries (MCP tools, native agent tools).
- Scoping (global vs device-specific overrides).

#### Scenario: Validate Whitelist Rule
- **GIVEN** a raw rule string with valid regex or command prefix
- **WHEN** policy engine parses and validates the rule
- **THEN** it produces a structured `PolicyRule` object or returns a descriptive validation error if syntax is unsafe.

### Requirement: Intelligent Merging & Conflict Resolution
The policy engine MUST merge multi-device policy updates deterministically without duplicate entries.

#### Scenario: Merging Disjoint Whitelists
- **GIVEN** local policy with rules `[A, B]` and remote policy with rules `[B, C]`
- **WHEN** merge is evaluated
- **THEN** the resulting policy contains unified rules `[A, B, C]` preserving metadata and order.
