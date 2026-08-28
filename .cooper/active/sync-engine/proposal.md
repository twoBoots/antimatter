# Track Proposal: Git-Backed Remote Sync Engine

## Rationale & Goal
Implement a Git-backed dotfile repository synchronization backend for Antimatter, enabling developers to persist, version control, and seamlessly synchronize agent whitelist policies and settings across all their machines via private or team Git repositories (e.g. GitHub/GitLab).

## User & Developer Benefits
- **Multi-Device Cohesion**: Sync approved commands and tool whitelists across developer laptops, desktops, and cloud workstations.
- **Git Native**: Full audit history of policy changes, rollbacks, and team branch merging.
- **Secrets Sanitization**: Automatically strip sensitive tokens (API keys, auth credentials) before pushing policy diffs to remote repositories.
- **Unified Sync Commands**: Single-command synchronization (`antimatter sync`, `antimatter push`, `antimatter pull`, `antimatter diff`, `antimatter remote`).

## Scope Boundaries
- **In-Scope**:
  - `pkg/sync/` Git backend (Init, Clone, Status, Diff, Pull, Push, and Merge).
  - Secret sanitization engine scanning token patterns before git commits.
  - End-to-end sync orchestrator integrating `pkg/policy` and `pkg/adapters`.
  - CLI subcommands: `antimatter sync`, `push`, `pull`, `diff`, and `remote`.
- **Out-of-Scope (Future Enhancements)**:
  - Background daemon / recurring cron auto-sync (can be added as an optional flag or future track).
