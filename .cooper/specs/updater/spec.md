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

### Requirement: SemVer 2.0.0 Standardized Comparison
The self-updater version comparison MUST adhere to SemVer 2.0.0 specification via `golang.org/x/mod/semver`, supporting pre-release tags, build metadata, and partial version strings.

#### Scenario: Evaluate Pre-Release and Partial Target Versions
- **GIVEN** a running binary at version `v1.0.0`
- **WHEN** user checks for updates against candidate pre-release or partial versions
- **THEN** updater correctly orders versions according to SemVer 2.0.0 precedence.

