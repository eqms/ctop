# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`ctop` is a Go terminal UI that shows real-time container metrics (CPU / mem / net / I/O) for Docker and runC. This repo is the maintained fork of the upstream `bcicen/ctop`, living at `github.com/eqms/ctop`. The module path in `go.mod` is `github.com/eqms/ctop` — not `eqms/ctop` or `bcicen/ctop`.

## Build, test, lint, run

```bash
make build           # cgo-free binary with version+build ldflags from VERSION + git HEAD
make build-all       # 5-platform cross-compile into _build/ with sha256sums.txt
make run-dev         # build with CTOP_DEBUG=1 and run — enables the debug log server
make image           # docker build

go test ./...        # unit tests. widgets/view_test.go and redact/redact_test.go exist today
golangci-lint run --timeout=5m   # local config is v2 (version: "2" at the top of .golangci.yml)
GOOS=linux go build ./...        # REQUIRED before touching connector/runc.go — see below
```

**Linux-only files**: `connector/runc.go` and `connector/collector/runc.go` have `//go:build linux`. They are invisible to `go build` on macOS. If you edit them, verify with `GOOS=linux go build ./...` — `go test` on macOS will not compile them.

**Local lint ≠ CI lint version**: CI pins `golangci-lint v2.11` via the v8 action. Keep both in sync when bumping.

## Release workflow — tag-driven

There is no manual release. Everything fires from a `v*` tag on master.

```bash
# 1. bump file + notes
echo "0.8.7" > VERSION
# edit RELEASE_NOTES.md — bilingual, DE section first, then EN, each under its own H2
# 2. commit + push
git commit -am "[CHG] Release v0.8.7: ..."
git push origin master
# 3. tag triggers everything
git tag v0.8.7
git push origin v0.8.7
```

The `v*` tag on push triggers `.github/workflows/ci.yml`:
- `build` matrix → 5 binaries uploaded as artifacts
- `release` → `softprops/action-gh-release@v2` creates the GitHub Release with `generate_release_notes: true` and all binaries + `sha256sums.txt` attached
- `docker` → builds `ghcr.io/eqms/ctop:<version>` and `:latest` for linux/amd64+arm64
- `homebrew` → computes SHAs of the 4 unix binaries, patches `homebrew-tap/Formula/ctop.rb`, and auto-commits `[CHG] Update Homebrew formula to <version>` back to master as `github-actions[bot]`

Consequence: **do not hand-edit `homebrew-tap/Formula/ctop.rb`** — the next release overwrites it. Only bump the file if you are shipping a formula-only change without a release.

The `VERSION` file is the single source of truth. `main.go` reads its value via `-X main.version=...` ldflags, not at runtime.

## Architecture — what requires reading multiple files

### The render loop is termui + termbox, not Bubble Tea (yet)

`main.go` holds package-level globals: `cursor`, `cGrid`, `header`, `status`, `errView`, `log`. They are initialised once and mutated in place. Every file in this repo that renders to the terminal imports both `ui "github.com/gizak/termui"` and `tm "github.com/nsf/termbox-go"` — they are not interchangeable, one is the widget lib and the other is the raw terminal backend. Event handling is `ui.Handle("/sys/kbd/...", fn)` wired up in `grid.go`, `menus.go`, `cursor.go`.

**Both UI libraries are dead upstream.** `_docs/bubbletea-migration-plan.md` is the agreed 6-phase plan to replace them with Charm's Bubble Tea + Lipgloss + Bubbles. Until that migration starts, keep the termui patterns — do not introduce Bubble Tea piecemeal.

### Connector pattern — ConnectorSuper wraps the backend

`connector/main.go` defines a registry: each backend calls `init() { enabled["name"] = NewX }`. `ByName(name)` wraps the constructor in a `ConnectorSuper` which handles reconnect loops — the TUI layer only sees `ConnectorSuper.Get()` returning `(Connector, error)`. A connector must implement `All() container.Containers`, `Get(id) (*container.Container, bool)`, `Wait() struct{}`.

Backends:
- `connector/docker.go` — Docker via `fsouza/go-dockerclient` (+ its transitive `docker/docker`). Cross-platform.
- `connector/runc.go` — runC via `github.com/opencontainers/runc/libcontainer` **v1.4+ API** (`libcontainer.Load(root, id)`, no Factory). Linux-only. cgroup stats come from `github.com/opencontainers/cgroups` (extracted from runc in v1.2+).
- `connector/mock.go` — synthetic data for `run-dev`. Used when `CTOP_CONNECTOR=mock` or similar.

### Compact columns are a registry, not hardcoded

`cwidgets/compact/column.go` exposes `allCols map[string]NewCompactColFn` — the set of column types ctop knows how to render. The **order and enabled set** comes from `config.EnabledColumns()` (`config/columns.go`), which is populated from the TOML config file. To add a new column: (1) implement `CompactCol` in a new file under `cwidgets/compact/`, (2) register its constructor in `allCols`, (3) add the name to the default column list in `config/columns.go`.

`CompactCol` interface lives in `column.go:50`. Each column is responsible for its own `FixedWidth()` vs. dynamic sizing, `Highlight()`/`UnHighlight()`, and receiving `SetMeta(Meta)` / `SetMetrics(Metrics)` per tick.

### Global config singleton

`config/main.go` exposes `Init()`, `Read()`, `Update(key, value)`, `Toggle(key)`, `EnabledColumns()` as package-level functions with an implicit singleton. CLI flags in `main.go` call `config.Update(...)` and `config.Toggle(...)` to override file defaults. Don't build a second config object — read/write through this package.

## Conventions

- Commit prefixes (enforced by CI's release-notes parsing): `[ADD]` new feature, `[CHG]` modification, `[FIX]` bugfix. Use `[CHG]` for dep upgrades and CI changes.
- `RELEASE_NOTES.md` is bilingual: `## Deutsche Release Notes` first, `## English Release Notes` second, each with `### v<x.y.z> (YYYY-MM-DD)` sections. Both sides must be updated for every release.
- Go version in `go.mod` follows the actively supported release line. Bump both `go.mod` and all three `go-version:` entries in `.github/workflows/ci.yml` together.
- The `// +build` legacy constraint coexists with `//go:build` on existing Linux-only files — leave the legacy line in place when editing.
