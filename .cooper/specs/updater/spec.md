# Self-Updater Capability Specification

## Purpose
Provides atomic in-place binary self-updating for `antimatter` powered by `twoBoots/bender/pkg/updater` via GitHub Releases.

## Requirements

### Requirement: SemVer Check & Self-Update
The updater MUST query the GitHub Releases API, compare semantic versions, and replace the running binary in-place.

#### Scenario: Check for Newer Version
- **GIVEN** running binary at version `v1.0.0` and latest GitHub release at `v1.1.0`
- **WHEN** user executes `antimatter update --check`
- **THEN** it reports that an update is available without modifying the local binary.

#### Scenario: In-Place Atomic Binary Replacement
- **GIVEN** an available update
- **WHEN** user executes `antimatter update`
- **THEN** matching platform asset is downloaded, checksum verified, quarantine stripped on macOS (`codesign -s - --force`), and atomically replaced.
