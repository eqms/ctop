# External Integrations

**Analysis Date:** 2026-06-11

## Container Runtimes

**Docker:**
- Purpose: Primary container backend — list containers, stream stats, subscribe to events
- SDK/Client: `github.com/fsouza/go-dockerclient v1.13.1` (wraps `github.com/moby/moby/client`)
- Auth/Connection: `api.NewClientFromEnv()` reads `DOCKER_HOST`, `DOCKER_CERT_PATH`, `DOCKER_TLS_VERIFY` from environment; defaults to Unix socket `/var/run/docker.sock`
- Implementation: `connector/docker.go`, collectors in `connector/collector/docker.go` and `connector/collector/docker_logs.go`
- Manager: `connector/manager/docker.go`

**runC (OCI):**
- Purpose: Alternative backend for raw runC containers (no Docker daemon required)
- SDK/Client: `github.com/opencontainers/runc v1.4.0` via `libcontainer.Load(root, id)` API; cgroup stats from `github.com/opencontainers/cgroups v0.0.6`
- Auth/Connection: Filesystem access to `RUNC_ROOT` (default `/run/runc`); `RUNC_SYSTEMD_CGROUP=1` for systemd cgroup variant
- Implementation: `connector/runc.go` — **Linux-only** (`//go:build linux`), `connector/collector/runc.go`
- Manager: `connector/manager/runc.go`

**Mock Connector:**
- Purpose: Synthetic data for development and testing without a live container runtime
- Activated via: `CTOP_CONNECTOR=mock` or `-connector mock` CLI flag
- Implementation: `connector/mock.go`, `connector/collector/mock.go`, `connector/collector/mock_logs.go`

## Data Storage

**Databases:**
- None. ctop is a stateless read-only TUI; it holds no persistent database.

**File Storage:**
- User config file: TOML at `~/.config/ctop/config` (XDG) or `~/.ctop`
  - Written by `config/file.go` `Write()` function
  - Read at startup by `config/file.go` `Read()` function
- Debug log file: optional, path from `CTOP_DEBUG_FILE` env var
  - Written by `logging/main.go` `ctopHandler`

**Caching:**
- In-process only: `connector/docker.go` maintains `containers map[string]*container.Container` protected by `sync.RWMutex`; no external cache

## Authentication & Identity

**Auth Provider:**
- None — ctop reads credentials from Docker's standard environment variables (`DOCKER_HOST`, `DOCKER_CERT_PATH`, `DOCKER_TLS_VERIFY`) and defers all authentication to the Docker client library
- No built-in login flow; mutual TLS for Docker is handled transparently by `go-dockerclient`

## Monitoring & Observability

**Error Tracking:**
- None. Errors surface in the TUI status bar and the internal ring-buffer log.

**Logs:**
- Custom structured logger wrapping `log/slog` (stdlib): `logging/main.go`
- Ring buffer (1024 messages) with multiple subscriber channels
- Output destinations: in-memory ring (always), optional file (`CTOP_DEBUG_FILE`), optional debug HTTP server (Unix socket or TCP when `CTOP_DEBUG=1`)
- Debug HTTP server: `logging/server.go` — serves a log tail endpoint; opened automatically when `CTOP_DEBUG=1`; URL opened in browser via `github.com/pkg/browser`
- Log level: INFO by default; DEBUG when `CTOP_DEBUG=1`

## CI/CD & Deployment

**Hosting:**
- GitHub: `github.com/eqms/ctop` (primary)
- Container registry: `ghcr.io/eqms/ctop` (GitHub Container Registry)
- Package distribution: GitHub Releases (binaries + `sha256sums.txt`)
- Homebrew: `homebrew-tap/Formula/ctop.rb` distributed from the same repo

**CI Pipeline:**
- GitHub Actions: `.github/workflows/ci.yml`
- Jobs (all triggered on push to master/main or `v*` tag):
  - `test` — `go test ./...` + `govulncheck` (continue-on-error)
  - `lint` — `golangci/golangci-lint-action@v8` at `v2.11`
  - `build` — 5-platform matrix cross-compile, artifacts uploaded via `actions/upload-artifact@v4`
  - `release` (tag-triggered) — `softprops/action-gh-release@v2`, attaches all binaries + checksums, `generate_release_notes: true`
  - `docker` (tag-triggered) — multi-platform build (`linux/amd64`, `linux/arm64`) via `docker/build-push-action@v6`, pushed to `ghcr.io`
  - `homebrew` (tag-triggered) — auto-patches `homebrew-tap/Formula/ctop.rb` with new version + SHA256 checksums, commits back to master as `github-actions[bot]`
- Node.js 24 opt-in: `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true` env var set globally in CI
- Go version pinned: `go-version: '1.26'` in all three action steps

## Environment Configuration

**Required environment variables at runtime:**
- `HOME` — mandatory; used to locate config file (`config/file.go`)
- Docker connector also reads: `DOCKER_HOST`, `DOCKER_CERT_PATH`, `DOCKER_TLS_VERIFY` (all optional; fall back to default socket)

**Optional environment variables:**
- `CTOP_DEBUG=1` — debug mode + log server
- `CTOP_DEBUG_TCP=1` — TCP debug server instead of Unix socket
- `CTOP_DEBUG_FILE=<path>` — write logs to file
- `CTOP_SHELL=<path>` — shell used for `exec` into containers
- `RUNC_ROOT=<path>` — runC state root (default: `/run/runc`)
- `RUNC_SYSTEMD_CGROUP=1` — use systemd cgroup driver
- `XDG_CONFIG_HOME` — XDG config dir (detected by scanning env for `XDG_*` prefix)

**Secrets location:**
- No application-level secrets. Docker TLS credentials live in the path referenced by `DOCKER_CERT_PATH`. CI uses `secrets.GITHUB_TOKEN` (injected by GitHub Actions) for GHCR push and release creation.

## Webhooks & Callbacks

**Incoming:**
- None. ctop has no HTTP server in normal operation.
- Debug mode only: `logging/server.go` opens a local Unix socket (or TCP port when `CTOP_DEBUG_TCP=1`) to stream log output; not a webhook endpoint.

**Outgoing:**
- None. ctop is a read-only metrics viewer; it does not call external HTTP endpoints.
- `github.com/pkg/browser` opens the default browser to the debug server URL when `CTOP_DEBUG=1` — this is a local URL, not an external callback.

---

*Integration audit: 2026-06-11*
