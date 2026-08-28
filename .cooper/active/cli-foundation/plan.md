# Implementation Plan: CLI Foundation & Policy Engine Scaffold

## Phase 1: Go Module Setup & Policy Engine Domain Core (TDD)

- [x] Task 1.1: Go Module Initialization & Dependencies (8dcb4d7)
  - [x] Sub-task: Initialize `go.mod` for `github.com/twoBoots/antimatter`
  - [x] Sub-task: Import `github.com/spf13/cobra`, `gopkg.in/yaml.v3`, `github.com/stretchr/testify`
  - [x] Sub-task: Set up `pkg/version` package

- [x] Task 1.2: Policy Data Models & Codec (TDD) (a278711)
  - [x] Sub-task: Write unit tests for YAML/JSON policy parsing & serialization (Red)
  - [x] Sub-task: Implement `pkg/policy` data models and codec (Green)
  - [x] Sub-task: Refactor & verify test coverage >80% (Refactor)

- [~] Task 1.3: Policy Evaluator & Match Engine (TDD)
  - [ ] Sub-task: Write unit tests for exact, prefix, glob, and regex command & tool matching (Red)
  - [ ] Sub-task: Implement `pkg/policy` matcher logic (Green)
  - [ ] Sub-task: Implement rule deduplication and merge resolver (Green)
  - [ ] Sub-task: Refactor & verify test coverage >80% (Refactor)

- [ ] Task 1.4: Phase 1 Verification & Checkpoint
  - [ ] Sync workflow & living specs (`git fetch origin main`)
  - [ ] Run full test suite with coverage report (`go test -v -coverprofile=coverage.out ./...`)
  - [ ] Checkpoint commit & Git Notes summary

---

## Phase 2: CLI Command Scaffolding & Bender Integration

- [ ] Task 2.1: Cobra Command Hierarchy Scaffolding
  - [ ] Sub-task: Write CLI unit/integration tests for root and version commands (Red)
  - [ ] Sub-task: Implement `cmd/root.go`, `cmd/version.go`, `main.go` (Green)
  - [ ] Sub-task: Refactor & test help text output (Refactor)

- [ ] Task 2.2: Whitelist CLI Commands (`list`, `add`, `check`)
  - [ ] Sub-task: Write CLI tests for `antimatter whitelist` subcommands (Red)
  - [ ] Sub-task: Implement `cmd/whitelist.go` invoking `pkg/policy` (Green)
  - [ ] Sub-task: Refactor & verify output formatting (Refactor)

- [ ] Task 2.3: Bender Integration (`update` & `mcp`)
  - [ ] Sub-task: Implement `cmd/update.go` integrated with `github.com/twoBoots/bender/pkg/updater`
  - [ ] Sub-task: Implement `cmd/mcp.go` integrated with `github.com/twoBoots/bender/pkg/mcp`
  - [ ] Sub-task: Register baseline Antimatter MCP tools (`antimatter_get_whitelist`, `antimatter_check_command`)

- [ ] Task 2.4: Phase 2 Verification & Checkpoint
  - [ ] Sync workflow & living specs (`git fetch origin main`)
  - [ ] Run full test suite and verify >80% coverage
  - [ ] Manual test of CLI binary execution (`./bin/antimatter version`, `./bin/antimatter whitelist --help`)
  - [ ] Checkpoint commit & Git Notes summary
