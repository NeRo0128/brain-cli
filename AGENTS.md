# AGENTS.md - Brain CLI v1.0.0

## Quick Start

```bash
# Build
make build          # → build/brain-cli

# Run
make run            # → go run cmd/brain-cli/main.go

# Test
make test           # → go test ./... -cover
make test-verbose   # → go test ./... -v -cover

# Lint
make vet            # → go vet ./...
go fmt ./...
goimports -w .
```

## Commands

| Action | Command |
|--------|---------|
| Build binary | `make build` |
| Run app | `make run` |
| All tests | `make test` |
| Verbose tests | `make test-verbose` |
| Go vet | `make vet` |
| Cross-compile | `make release` |

## Architecture

- **Domain** (`internal/core/`): Entities (Task, Provider, Tool, Execution) + repository interfaces
- **Use Cases** (`internal/usecases/`): Application logic (task execution, AI, provider mgmt)
- **Adapters** (`internal/adapters/`): Concrete implementations (SQLite, executor, AI providers, config)
- **UI** (`internal/ui/`): Bubble Tea TUI with Lipgloss styles
- **Entry point**: `cmd/brain-cli/main.go` — wires deps and starts TUI

## Configuration

- File: `configs/config.yaml`
- Env vars prefix: `BRAIN_` (e.g. `BRAIN_DB_PATH`, `BRAIN_THEME`, `BRAIN_LOG_LEVEL`)
- Config loaded at startup via `config.EnsureConfig()`; overrides applied by settings manager
- Defaults: DB `data/brain.db`, log level `debug`, theme `brain`

## Build/Run Notes

- `CGO_ENABLED=0` in Makefile builds — static binary
- LDFLAGS set `Version`, `Commit`, `BuildDate` via `-X main.*`
- DB auto-migrates on startup if `database.auto_migrate: true`
- `paths.DevMode()` affects which config path is used (dev: `configs/config.yaml`; prod: `paths.ConfigFile()`)

## Key Directories

```
internal/core/      → entities + repo interfaces
internal/usecases/ → task execution, AI, provider mgmt
internal/adapters/  → SQLite, executor, AI providers, config, email
internal/ui/        → Bubble Tea TUI (model, screens, components, styles)
cmd/brain-cli/      → main entrypoint
configs/            → config.yaml + .env.example
docker/             → Dockerfile + compose
scripts/            → user-defined scripts
```

## Testing

- Most internal packages have unit tests (`go test ./...` passes)
- No tests in `cmd/` or `internal/ui/`
- Run coverage: `go test -cover ./...`
- Generate report: `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`

## Git & Release

- Conventional commits: `feat:`, `fix:`, `refactor:`
- Cross-compile: `make release` → 6 platforms (linux/arm64, darwin, windows)
- Checksums: `make checksums` → SHA256SUMS in release dir