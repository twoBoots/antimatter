# Track Proposal: Upgrade Bender to v1.1.0

## Rationale & Goal
Upgrade the foundational dependency `github.com/twoBoots/bender` from `v1.0.1` to `v1.1.0`. Bender v1.1.0 introduces standardized SemVer 2.0.0 parsing and comparison in `pkg/updater` (powered by `golang.org/x/mod/semver`), addressing edge cases in version comparison (pre-releases, build metadata, partial version strings). Additionally, promote `bender` from an `// indirect` dependency to an explicit direct dependency in `go.mod`.

## User & Developer Benefits
- **Robust Update Engine**: Self-update checks accurately handle SemVer 2.0.0 specs including pre-releases (`-rc.1`, `-beta.2`) and partial version tags.
- **Dependency Hygiene**: Clean, direct dependency tracking in `go.mod` without indirect marker ambiguity for `pkg/updater` and `pkg/mcp`.
- **Ecosystem Alignment**: Aligns Antimatter with the latest release across the twoBoots tooling suite.

## Scope Boundaries
- **In-Scope**:
  - Upgrade `github.com/twoBoots/bender` to `v1.1.0` in `go.mod` and promote to direct requirement.
  - Verify and test `pkg/updater` integration in `cmd/update.go` and `cmd/update_test.go` against Bender v1.1.0 SemVer engine.
  - Verify `pkg/mcp` integration in `cmd/mcp.go` and `cmd/mcp_test.go`.
  - Full test suite passing with race detector and >80% code coverage.
- **Out-of-Scope**:
  - Introducing new MCP tools (future track).
  - Modifying policy engine core logic.
