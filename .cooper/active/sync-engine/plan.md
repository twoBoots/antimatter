# Implementation Plan: Git-Backed Remote Sync Engine

## Phase 1: Git Backend, Secret Sanitizer & Syncer Domain Core (TDD)

- [~] Task 1.1: Secret Sanitization Engine (TDD)
  - [ ] Sub-task: Write unit tests for API key / token pattern detection and masking (Red)
  - [ ] Sub-task: Implement `pkg/sync/sanitizer.go` (Green)
  - [ ] Sub-task: Refactor & verify test coverage >80% (Refactor)

- [ ] Task 1.2: Git Backend Driver (TDD)
  - [ ] Sub-task: Write unit tests for local repo init, commit, diff, status, and remotes (Red)
  - [ ] Sub-task: Implement `pkg/sync/git.go` (Green)
  - [ ] Sub-task: Refactor & verify test coverage >80% (Refactor)

- [ ] Task 1.3: Sync Orchestrator (TDD)
  - [ ] Sub-task: Write unit tests for full push/pull/merge sync cycles (Red)
  - [ ] Sub-task: Implement `pkg/sync/syncer.go` coordinating policy merger and adapters (Green)
  - [ ] Sub-task: Refactor & verify test coverage >80% (Refactor)

- [ ] Task 1.4: Phase 1 Verification & Checkpoint
  - [ ] Sync workflow & living specs (`git fetch origin main`)
  - [ ] Run full test suite with coverage report (`go test -v -coverprofile=coverage.out ./...`)
  - [ ] Checkpoint commit & Git Notes summary

---

## Phase 2: CLI Commands & End-to-End Workflows

- [ ] Task 2.1: CLI Sync & Remote Commands (`sync`, `push`, `pull`, `diff`, `remote`)
  - [ ] Sub-task: Write CLI tests for sync command hierarchy (Red)
  - [ ] Sub-task: Implement `cmd/sync.go` (Green)
  - [ ] Sub-task: Refactor CLI flags and output formatting (Refactor)

- [ ] Task 2.2: Phase 2 Verification & Checkpoint
  - [ ] Sync workflow & living specs (`git fetch origin main`)
  - [ ] Run full test suite and verify >80% coverage
  - [ ] Manual test of CLI binary execution (`./bin/antimatter sync --help`, `./bin/antimatter remote --help`)
  - [ ] Checkpoint commit & Git Notes summary
