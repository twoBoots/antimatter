# Technology Stack Definition

## Primary Languages & Runtime
- **Language**: Go (`go 1.27+`)
- **Runtime / Engine**: Static Go binary compiled with `CGO_ENABLED=0` for cross-platform portability (`darwin`, `linux`, `windows`)

## Architecture & Frameworks
- **CLI Framework**: `github.com/spf13/cobra`
- **Archetype & Core Libraries**: `github.com/twoBoots/bender`
  - `pkg/updater`: In-place binary self-updater with GitHub Releases API integration
  - `pkg/mcp`: Stdio JSON-RPC 2.0 embedded MCP server engine and multi-client configuration installer
- **Configuration & Serialization**: `gopkg.in/yaml.v3` (Canonical policies) & `encoding/json`
- **Git Backend**: Pure Go Git implementation (`github.com/go-git/go-git/v5`) or OS `git` subprocess wrapper for robust dotfile synchronizations

## Testing & Quality Control
- **Test Runner / Framework**: Go standard `testing` package with `github.com/stretchr/testify` (`assert`, `require`)
- **Coverage Target**: >80% code coverage across all packages
- **Linter & Formatter**: `gofmt`, `go vet`, `golangci-lint`

## Build & CI/CD
- **Package Manager**: Go Modules (`go.mod`)
- **CI System**: GitHub Actions (multi-platform matrix build and SemVer release tagging)
- **Installer**: 3-tier zero-prompt `install.sh` (Local compile -> GitHub release download -> Fallback)
