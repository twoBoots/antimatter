# Implementation Plan: Upgrade Bender to v1.1.0

## Phase 1: Dependency Upgrade & Integration Verification (TDD)

- [x] Task 1.1: Dependency Promotion & Upgrade (TDD) (da97f5c)
  - [x] Sub-task: Write unit test validating version comparison and update flags (Red)
  - [x] Sub-task: Upgrade `github.com/twoBoots/bender` to `v1.1.0` as direct dependency in `go.mod` (Green)
  - [x] Sub-task: Verify tests pass and maintain >80% coverage (Refactor)

- [~] Task 1.2: MCP Server & Updater Regression Tests (TDD)
  - [ ] Sub-task: Write integration assertions for MCP server build and updater options (Red)
  - [ ] Sub-task: Ensure seamless integration with Bender v1.1.0 MCP engine (Green)
  - [ ] Sub-task: Refactor test helpers and ensure clean outputs (Refactor)

- [ ] Task 1.3: Phase 1 Verification & Checkpoint
  - [ ] Sync workflow & living specs (`git fetch origin main`)
  - [ ] Run full test suite with race detection (`go test -race -v -coverprofile=coverage.out ./...`)
  - [ ] Verify test coverage >80%
  - [ ] Checkpoint commit & Git Notes summary [checkpoint: pending]
