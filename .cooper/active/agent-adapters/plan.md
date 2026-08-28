# Implementation Plan: Live Agent Adapters for Antigravity & Claude Code

## Phase 1: Adapter Core Framework & Backup Engine (TDD)

- [x] Task 1.1: Adapter Interface, Types & Backup Engine (1637fd6)
  - [x] Sub-task: Write unit tests for backup engine and adapter manager (Red)
  - [x] Sub-task: Implement `pkg/adapters/adapter.go`, `pkg/adapters/manager.go`, `pkg/adapters/backup.go` (Green)
  - [x] Sub-task: Refactor & verify test coverage >80% (Refactor)

- [x] Task 1.2: Antigravity Adapter Implementation (TDD) (3bd40ce)
  - [x] Sub-task: Write unit tests for Antigravity export & import with mock filesystems (Red)
  - [x] Sub-task: Implement `pkg/adapters/antigravity/adapter.go` (Green)
  - [x] Sub-task: Refactor & verify test coverage >80% (Refactor)

- [~] Task 1.3: Claude Code & Desktop Adapter Implementation (TDD)
  - [ ] Sub-task: Write unit tests for Claude Code export & import with mock filesystems (Red)
  - [ ] Sub-task: Implement `pkg/adapters/claude/adapter.go` (Green)
  - [ ] Sub-task: Refactor & verify test coverage >80% (Refactor)

- [ ] Task 1.4: Phase 1 Verification & Checkpoint
  - [ ] Sync workflow & living specs (`git fetch origin main`)
  - [ ] Run full test suite with coverage report (`go test -v -coverprofile=coverage.out ./...`)
  - [ ] Checkpoint commit & Git Notes summary

---

## Phase 2: CLI Integration & Sync Commands

- [ ] Task 2.1: CLI `export` and `import` Commands
  - [ ] Sub-task: Write CLI unit tests for `antimatter export` and `antimatter import` (Red)
  - [ ] Sub-task: Implement `cmd/sync_adapters.go` (Green)
  - [ ] Sub-task: Refactor & test CLI error handling (Refactor)

- [ ] Task 2.2: CLI `status` Command
  - [ ] Sub-task: Write CLI unit tests for `antimatter status` (Red)
  - [ ] Sub-task: Implement `cmd/status.go` displaying adapter sync state (Green)
  - [ ] Sub-task: Refactor output formatting (Refactor)

- [ ] Task 2.3: Phase 2 Verification & Checkpoint
  - [ ] Sync workflow & living specs (`git fetch origin main`)
  - [ ] Run full test suite and verify >80% coverage
  - [ ] Manual test of CLI binary execution (`./bin/antimatter status`, `./bin/antimatter export --help`)
  - [ ] Checkpoint commit & Git Notes summary
