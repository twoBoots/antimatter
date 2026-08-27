# Product Guidelines

## Brand Voice & Tone
- **Tone**: Precise, deterministic, safe, and developer-focused.
- **Prose Style**: Concise and actionable. Output should be machine-parseable when requested (JSON/YAML) and clean, human-friendly by default.

## User Experience & CLI Ergonomics
- **Zero-Surprise Operations**: Whitelist updates and sync executions must clearly display diffs before making modifications to live agent settings.
- **Safety First & Backups**: Always create timestamped backups before overwriting any active agent configuration files (`~/.gemini/antigravity-cli/`, `~/.claude/`).
- **Predictable Behavior**: Commands should exit with standard POSIX exit codes and provide clear next-step remediation on failures.
- **Agent Interoperability**: Treat Antigravity, Claude Code, and Claude Desktop as first-class citizens with equal priority for tool/command whitelists.

## Security & Secrets Management
- **Secret Isolation**: Never sync API tokens, credentials, session IDs, or unmasked environment secrets into version-controlled sync repositories.
- **Strict Validation**: Validate glob and regex expressions before persisting them into whitelist rules to prevent unintentional wide-open execution risks.

## Documentation & Markdown Standards
- Use GitHub Flavored Markdown.
- Ensure relative links between documents are valid and clickable.
- Structure complex information with tables and Mermaid diagrams where beneficial.
