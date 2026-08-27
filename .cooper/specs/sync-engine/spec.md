# Sync Engine Capability Specification

## Purpose
Manages version control, dotfile repository integration, remote sync (push/pull), local backups, and change diffing for agent configurations.

## Requirements

### Requirement: Git Backend Synchronization
The sync engine MUST support synchronizing settings with a dedicated Git repository backend.

#### Scenario: Pull Remote Settings
- **GIVEN** a configured remote sync repository
- **WHEN** user or background routine executes sync pull
- **THEN** remote changes are fetched, merged with local rules, and applied to active agent adapters.

#### Scenario: Push Local Changes
- **GIVEN** local changes to whitelists or profiles
- **WHEN** user executes sync push
- **THEN** changes are sanitized of secrets, committed to the sync repo, and pushed to remote.

### Requirement: Automatic Backup & Safety
The sync engine MUST create non-destructive backups before applying any updates to live agent configurations.

#### Scenario: Automatic Backup Creation
- **GIVEN** an active Antigravity or Claude configuration
- **WHEN** incoming sync changes are applied
- **THEN** a timestamped backup directory is preserved under `~/.antimatter/backups/`.
