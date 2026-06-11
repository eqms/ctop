# Technology Stack

**Analysis Date:** 2026-06-11

## Languages

**Primary:**
- Go 1.26 - Entire application: TUI, connectors, config, logging

**Secondary:**
- Ruby - Homebrew formula only (`homebrew-tap/Formula/ctop.rb`)
- YAML - CI/CD pipeline (`.github/workflows/ci.yml`), lint config (`.golangci.yml`)
- TOML - User config file format (`config/file.go`)

## Runtime

**Environment:**
- Go 1.26 (module: `github.com/eqms/ctop`)
- CGO_ENABLED=0 for all release builds (fully static binary, no libc dependency)
- Linux-only files: `connector/runc.go`, `connector/collector/runc.go` use `//go:build linux` + `// +build linux`

**Package Manager:**
- Go modules (`go.mod` / `go.sum`)
- Lockfile: `go.sum` present

## Frameworks

**TUI — two dead-upstream libraries (migration planned):**
- `github.com/gizak/termui v2.3.1-0.20180817033724-8d4faad06196+incompatible` — widget layer (layouts, rendering, event routing via `ui.Handle`)
- `github.com/nsf/termbox-go v0.0.0-20190121233118-02980233997d` — raw terminal backend (input mode, cell drawing, `tm.Clear`/`tm.Sync`)
- See `_docs/bubbletea-migration-plan.md` for the agreed 6-phase Bubble Tea + Lipgloss replacement plan; do not introduce Bubble Tea piecemeal

**Config:**
- `github.com/BurntSushi/toml v1.5.0` — TOML encode/decode for `~/.config/ctop/config` (or `~/.ctop`)

**Testing:**
- Go standard `testing` package only
- Single test file: `widgets/view_test.go`
- Run: `go test ./...`

**Build/Dev:**
- `make build` — CGO_ENABLED=0 static binary with `-X main.version` / `-X main.build` ldflags
- `make build-all` — 5-platform cross-compile into `_build/` with `sha256sums.txt`
- `make run-dev` — build without `-tags release`, set `CTOP_DEBUG=1`, runs locally
- `make image` — Docker build via `Dockerfile`

**Lint:**
- `golangci-lint v2` config: `.golangci.yml` (version: "2" header)
- CI pins `golangci-lint v2.11` via `golangci/golangci-lint-action@v8`

## Key Dependencies

**Critical:**
- `github.com/fsouza/go-dockerclient v1.13.1` — Docker connector; wraps Moby API; cross-platform
- `github.com/opencontainers/runc v1.4.0` — runC connector via `libcontainer.Load(root, id)` v1.4+ API; Linux-only
- `github.com/opencontainers/cgroups v0.0.6` — cgroup stats for runC connector (extracted from runc in v1.2+)
- `github.com/gizak/termui` + `github.com/nsf/termbox-go` — TUI rendering (see above)

**Utility:**
- `github.com/BurntSushi/toml v1.5.0` — config persistence
- `github.com/google/uuid v1.6.0` — container ID generation (mock connector)
- `github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b` — human-readable duration display
- `github.com/jgautheron/codename-generator v0.0.0-20150829203204-16d037c7cc3c` — mock container name generation
- `github.com/mattn/go-runewidth v0.0.16` — unicode-aware column width calculation
- `github.com/c9s/goprocinfo v0.0.0-20170609001544-b34328d6e0cd` — `/proc` filesystem parsing for host metrics
- `github.com/pkg/browser v0.0.0-20201207095918-0426ae3fba23` — open browser for debug log server URL

**Indirect (Moby/Docker ecosystem):**
- `github.com/moby/moby/api v1.54.0`, `github.com/moby/moby/client v0.3.0` — underlying Docker API client
- `github.com/sirupsen/logrus v1.9.4` — pulled in by go-dockerclient
- `google.golang.org/protobuf v1.36.8` — pulled in by runc/cgroups

## Configuration

**Environment variables:**
- `CTOP_DEBUG=1` — enables debug level logging + starts HTTP log server (Unix socket or TCP)
- `CTOP_DEBUG_TCP=1` — force TCP mode for debug server
- `CTOP_DEBUG_FILE=<path>` — write log output to file
- `CTOP_SHELL=<shell>` — override default exec shell (also settable via `-shell` flag)
- `DOCKER_HOST`, `DOCKER_CERT_PATH`, `DOCKER_TLS_VERIFY` — passed through to `api.NewClientFromEnv()` for Docker connector
- `RUNC_ROOT=<path>` — runC root path (default: `/run/runc`)
- `RUNC_SYSTEMD_CGROUP=1` — enable systemd cgroups for runC connector
- `HOME` — required for config file path resolution
- `XDG_CONFIG_HOME` — XDG config dir override (detected via `XDG_*` env prefix)

**User config file:**
- TOML format: `~/.config/ctop/config` (XDG) or `~/.ctop` (non-XDG)
- Sections: `[options]` (string key-value) and `[toggles]` (bool key-value)
- Columns order/enable state persisted as comma-separated string

**Build config:**
- `VERSION` — plain-text version file; source of truth for `make` ldflags and CI
- `-tags release` — build tag applied in all release builds (controls debug-only code paths)
- `go.mod` `go 1.26.0` line governs minimum Go toolchain

## Platform Requirements

**Development:**
- Go 1.26+
- `golangci-lint v2.11` for lint
- `govulncheck` for vulnerability scanning (CI: `continue-on-error: true`)
- Linux required to compile `connector/runc.go` (verify cross-compilation with `GOOS=linux go build ./...`)
- Docker daemon accessible for Docker connector testing

**Production / Distribution:**
- Static binary, no runtime dependencies (CGO_ENABLED=0)
- Platforms: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
- Docker image: `ghcr.io/eqms/ctop:<version>` and `:latest` for linux/amd64 + linux/arm64
- Homebrew tap: `homebrew-tap/Formula/ctop.rb` (auto-updated by CI on tag push)

---

*Stack analysis: 2026-06-11*
