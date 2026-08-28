# Track Proposal: Live Agent Adapters for Antigravity & Claude Code

## Rationale & Goal
Implement bidirectional configuration adapters for Google Antigravity and Claude Code / Claude Desktop, enabling seamless translation between Antimatter's canonical policy representation and the real-world settings/whitelists on developers' machines.

## User & Developer Benefits
- **Antigravity Synchronization**: Export canonical policies to Antigravity (`~/.gemini/antigravity-cli/`) and import live user-whitelisted commands and tools into Antimatter.
- **Claude Cross-Compatibility**: Export canonical policies to Claude Code (`~/.claude.json`) and Claude Desktop so tool approvals and allowed shell commands are synchronized between assistants.
- **Automated Backups & Safety**: Every export creates timestamped backups in `~/.antimatter/backups/<agent>/` prior to writing, guaranteeing zero data loss.
- **CLI Commands**: `antimatter export`, `antimatter import`, and `antimatter status`.

## Scope Boundaries
- **In-Scope**:
  - `pkg/adapters/` interface and registry (`Adapter`, `Manager`).
  - Google Antigravity adapter implementation and tests with mock environments.
  - Claude Code / Desktop adapter implementation and tests with mock environments.
  - Automatic backup creation mechanism.
  - CLI subcommands `export`, `import`, and `status`.
- **Out-of-Scope (Next Tracks)**:
  - Remote Git synchronization (pushing to a remote GitHub dotfiles repo; handled in `track_sync_engine`).
