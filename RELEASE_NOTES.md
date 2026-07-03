# Release Notes

> **Language / Sprache**: [Deutsch](#deutsche-release-notes) | [English](#english-release-notes)

---

## Deutsche Release Notes

### v0.8.8 (2026-07-03)

#### Sicherheit
- **Env-Redaction für Debug-Dump** (`redact/`, `debug.go`): Neues gemeinsames `redact`-Package — der Container-State-Dump (Taste `D`) maskiert Credential-artige Env-Variablen jetzt ebenfalls als `[REDACTED]`; bisher landeten sie unredacted im Log-Ring-Buffer, in der Debug-Datei (`CTOP_DEBUG_FILE`) und auf dem Debug-Socket. Das Pattern wurde um `AUTH`, `DSN`, `CERT`, `COOKIE`, `JWT` erweitert und erkennt zusätzlich Credentials im Wert selbst (Connection-Strings wie `postgres://user:pass@host/db`), unabhängig vom Variablennamen
- **CI-Supply-Chain-Härtung** (`.github/workflows/ci.yml`): Alle GitHub Actions auf volle Commit-SHAs gepinnt (mit Versions-Kommentar), `govulncheck` auf v1.5.0 gepinnt statt `@latest`, `timeout-minutes` für alle sechs Jobs
- **Dependabot** (`.github/dependabot.yml`): Wöchentliche automatische Updates für Go-Module, GitHub Actions (inkl. SHA-Pins) und das Docker-Base-Image
- **Docker-Härtung**: Base-Image `golang:1.26-alpine` per SHA256-Digest gepinnt; root-by-design für den `docker.sock`-Zugriff im Dockerfile dokumentiert; neue `.dockerignore` verschlankt den Build-Context

#### Intern
- `cwidgets/single/env.go` nutzt das gemeinsame `redact`-Package (Duplikat-Logik entfernt); neue Unit-Tests in `redact/redact_test.go`
- `.gitignore` um `*.sock`, `.env`, `*.pem`, `*.key` ergänzt

---

### v0.8.7 (2026-06-11)

#### Sicherheit
- **runc-Upgrade v1.4.0 → v1.4.2**: Schließt drei weitere High-Severity Container-Escape-Schwachstellen (GO-2025-4096, GO-2025-4097, GO-2025-4098), veröffentlicht im November 2025 und erst ab runc v1.4.1 gefixt
- **CI-Härtung** (`.github/workflows/ci.yml`): Shell-Injection über den `workflow_dispatch`-Tag-Input behoben (Übergabe via `env:` statt direkter `${{ }}`-Interpolation); Workflow-Permissions per Job gescoped (Top-Level nur noch `contents: read`); `govulncheck` blockiert jetzt den Build statt nur zu warnen
- **Debug-Socket abgesichert** (`logging/server.go`): Unix-Socket liegt jetzt in einem privaten Runtime-Verzeichnis (`$XDG_RUNTIME_DIR/ctop.sock` bzw. `$TMPDIR/ctop-<uid>/`, Verzeichnis 0700, Socket 0600) statt im Arbeitsverzeichnis; stale Sockets werden beim Start entfernt
- **Env-Variablen-Redaction** (`cwidgets/single/env.go`): Werte von Credential-artigen Variablen (`SECRET`, `PASSWORD`, `TOKEN`, `API_KEY`, …) werden in der Single-Container-Ansicht als `[REDACTED]` maskiert
- **RUNC_ROOT-Validierung** (`connector/runc.go`): Pfade außerhalb von `/run` und `/var/run` werden abgelehnt — eine manipulierte Env-Variable kann den Connector nicht mehr auf `/etc` oder `/proc` richten

#### Intern
- **Dependency-Bumps**: `BurntSushi/toml` v1.5.0 → v1.6.0, `fsouza/go-dockerclient` v1.13.1 → v1.13.2, Go-Direktive 1.26.0 → 1.26.4

---

### v0.8.6 (2026-04-21)

#### Sicherheit
- **runc-Upgrade v1.1.14 → v1.4.0**: Schließt drei Container-Escape-Schwachstellen (CVE-2025-31133, CVE-2025-52565, CVE-2025-52881), die in der runc 1.1.x-Reihe unpatched sind — 1.1.x ist end-of-life

#### Bugfixes
- **Race Conditions im runC-Connector** (`connector/runc.go`): Map-Zugriffe auf `libContainers` und `containers` werden jetzt korrekt mit `RLock` geschützt; `refreshAll` erstellt einen ID-Snapshot unter Lock vor dem Channel-Send
- **nil-Channel Deadlock**: `needsRefresh`-Channel wurde im Original-Code nie initialisiert — Sends und Receives wären stillschweigend blockiert
- **libcontainer-API-Migration**: Factory-Pattern durch direkten `libcontainer.Load(root, id)` ersetzt; `cgroups`-Paket aus `opencontainers/cgroups` bezogen (runc v1.2+ Änderung)

#### Intern
- **Go 1.24 → 1.26**: Go 1.24 EOL seit 11.02.2026
- **golangci-lint v1 → v2**: Config-Migration und CI-Action v6 → v8
- **Dependency Cleanup**: `cilium/ebpf` v0.7.0 → v0.17.3 entpinnt (der 1.1.14-Kompatibilitätspin ist obsolet); `fsouza/go-dockerclient` v1.7.0 → v1.13.1

#### Dokumentation
- **Bubbletea-Migrationsplan** (`_docs/bubbletea-migration-plan.md`): 6-Phasen-Plan für die Ablösung des toten TUI-Stacks (`gizak/termui` + `nsf/termbox-go`) durch den Charm-Stack

---

### v0.8.5 (2026-03-29)

#### Bugfixes
- **Menü-Redraw Fix**: Nach Rückkehr aus der Container-Shell wird das Hintergrund-Grid jetzt korrekt neu gezeichnet und das Container-Menü wieder fokussiert

---

### v0.8.4 (2026-02-28)

#### Bugfixes
- **Input-Widget Width-Fix**: Commit-Dialog Eingabefeld war zu schmal (28 statt 68 Zeichen) — `SetMaxLen()` berechnet jetzt die Widget-Breite korrekt neu
- **Exec Shell Fix**: Shell-Ausführung im Container überlagerte sich mit der TUI — Terminal wird jetzt vor Shell-Start freigegeben (`ui.Close()`) und nach Beendigung sauber wiederhergestellt (`ui.Init()`)

---

### v0.8.3 (2026-02-28)

#### Neue Features
- **Flicker-Fix**: Bildschirmflackern beim 1-Sekunden-Refresh behoben — Header und Grid werden jetzt in einem einzigen Render-Aufruf gezeichnet. Screen-Clear nur noch bei schrumpfender Container-Anzahl
- **Stop + Remove**: Neuer kombinierter Menüeintrag `[X] stop + remove` im Container-Menü für laufende Container — stoppt und entfernt in einem Schritt mit Sicherheitsabfrage
- **Docker Commit**: Container als lokale Images speichern via `[C] commit to image` im Container-Menü — Input-Dialog mit Namensvorschlag (`<name>-snapshot`) und optionalem Tag (Standard: `latest`)

---

### v0.8.2 (2026-02-28)

#### Neue Features
- **Versionsanzeige im Header**: Version wird rechts oben im TUI-Header angezeigt (z.B. `v0.8.2`)
- **Konfigurierbare Shell**: Shell für `exec` ist konfigurierbar via `--shell` Flag, `CTOP_SHELL` Umgebungsvariable oder Config-Datei (Priorität: CLI > Env > Config > Default `/bin/sh`)
- **Health-Status Spalte**: Neue optionale Spalte zeigt Health-Check-Status (`healthy`/`unhealthy`/`starting`) mit Farbkodierung
- **Restart-Count Spalte**: Neue optionale Spalte zeigt Container-Neustarts mit Farbwarnung (>0 gelb, >5 rot)

#### Bugfixes
- Versionsanzeige im Header korrekt positioniert (nicht mehr am äußersten rechten Rand)

#### CI/CD
- **Homebrew-Tap Auto-Update**: CI-Job aktualisiert automatisch die Homebrew-Formula nach einem Release
- Neues Makefile-Target `make update-homebrew` zum Synchronisieren des Homebrew-Tap Repos

---

### v0.8.0 (2025-06-24)

Erster Release des gepflegten Forks von [bcicen/ctop](https://github.com/bcicen/ctop).

#### Security-Fixes
- Shell-Injection-Schwachstelle in `container.Exec()` behoben — Verwendung von `exec.Command()` statt `sh -c`
- Dateiberechtigungen für Config-Datei auf `0600` gesetzt
- Debug-Server bindet jetzt auf `localhost` statt `0.0.0.0`

#### Modernisierung
- Aktualisierung auf **Go 1.22** mit modernen stdlib-Paketen
- Migration von `op/go-logging` auf `log/slog` (strukturiertes Logging)
- Ersetzung von `nu7hatch/gouuid` durch `google/uuid`
- Ersetzung von `pkg/errors` durch `fmt.Errorf` mit `%w`-Wrapping
- Ersetzung von `io/ioutil` durch `os`/`io`-Funktionen
- Verwendung von `os.ReadDir` statt `ioutil.ReadDir`

#### CI/CD & Infrastruktur
- GitHub Actions Workflow für automatische Multi-Plattform-Builds (Linux, macOS, Windows — AMD64 + ARM64)
- Automatische GitHub Releases mit Binaries und Checksummen
- Docker-Image via `ghcr.io/eqms/ctop`

#### Dokumentation
- Zweisprachige README (Deutsch/Englisch)
- Fork-Motivation und Änderungsübersicht dokumentiert

---

## English Release Notes

### v0.8.8 (2026-07-03)

#### Security
- **Env redaction for debug dump** (`redact/`, `debug.go`): new shared `redact` package — the container state dump (key `D`) now masks credential-like env variables as `[REDACTED]` too; previously they were logged unredacted to the log ring buffer, the debug file (`CTOP_DEBUG_FILE`) and the debug socket. The pattern was extended with `AUTH`, `DSN`, `CERT`, `COOKIE`, `JWT` and additionally detects credentials embedded in values (connection strings like `postgres://user:pass@host/db`) regardless of the variable name
- **CI supply-chain hardening** (`.github/workflows/ci.yml`): all GitHub Actions pinned to full commit SHAs (with version comments), `govulncheck` pinned to v1.5.0 instead of `@latest`, `timeout-minutes` on all six jobs
- **Dependabot** (`.github/dependabot.yml`): weekly automated updates for Go modules, GitHub Actions (incl. SHA pins) and the Docker base image
- **Docker hardening**: base image `golang:1.26-alpine` pinned by SHA256 digest; root-by-design for `docker.sock` access documented in the Dockerfile; new `.dockerignore` slims down the build context

#### Internal
- `cwidgets/single/env.go` now uses the shared `redact` package (duplicate logic removed); new unit tests in `redact/redact_test.go`
- `.gitignore` extended with `*.sock`, `.env`, `*.pem`, `*.key`

---

### v0.8.7 (2026-06-11)

#### Security
- **runc upgrade v1.4.0 → v1.4.2**: Patches three more high-severity container-escape vulnerabilities (GO-2025-4096, GO-2025-4097, GO-2025-4098), published November 2025 and only fixed from runc v1.4.1 onwards
- **CI hardening** (`.github/workflows/ci.yml`): fixed shell injection via the `workflow_dispatch` tag input (passed through `env:` instead of direct `${{ }}` interpolation); workflow permissions scoped per job (top level reduced to `contents: read`); `govulncheck` now blocks the build instead of only warning
- **Debug socket secured** (`logging/server.go`): the Unix socket now lives in a private runtime directory (`$XDG_RUNTIME_DIR/ctop.sock` or `$TMPDIR/ctop-<uid>/`, directory 0700, socket 0600) instead of the working directory; stale sockets are removed on startup
- **Env variable redaction** (`cwidgets/single/env.go`): values of credential-like variables (`SECRET`, `PASSWORD`, `TOKEN`, `API_KEY`, …) are masked as `[REDACTED]` in the single-container view
- **RUNC_ROOT validation** (`connector/runc.go`): paths outside `/run` and `/var/run` are rejected — a poisoned env variable can no longer point the connector at `/etc` or `/proc`

#### Internal
- **Dependency bumps**: `BurntSushi/toml` v1.5.0 → v1.6.0, `fsouza/go-dockerclient` v1.13.1 → v1.13.2, Go directive 1.26.0 → 1.26.4

---

### v0.8.6 (2026-04-21)

#### Security
- **runc upgrade v1.1.14 → v1.4.0**: Patches three container-escape vulnerabilities (CVE-2025-31133, CVE-2025-52565, CVE-2025-52881) that remain unpatched in the runc 1.1.x line — 1.1.x is end-of-life

#### Bugfixes
- **Race conditions in the runC connector** (`connector/runc.go`): map accesses on `libContainers` and `containers` are now correctly guarded with `RLock`; `refreshAll` snapshots IDs under the lock before sending on the channel
- **Nil-channel deadlock**: `needsRefresh` channel was never initialised in the original code — any send or receive would have silently blocked
- **libcontainer API migration**: replaced Factory pattern with direct `libcontainer.Load(root, id)`; `cgroups` package now pulled from `opencontainers/cgroups` (runc v1.2+ split)

#### Internal
- **Go 1.24 → 1.26**: Go 1.24 reached EOL on 2026-02-11
- **golangci-lint v1 → v2**: config migrated and CI action v6 → v8
- **Dependency cleanup**: `cilium/ebpf` unpinned v0.7.0 → v0.17.3 (the 1.1.14-compat pin is now obsolete); `fsouza/go-dockerclient` v1.7.0 → v1.13.1

#### Documentation
- **Bubble Tea migration plan** (`_docs/bubbletea-migration-plan.md`): 6-phase plan to replace the dead TUI stack (`gizak/termui` + `nsf/termbox-go`) with the Charm stack

---

### v0.8.5 (2026-03-29)

#### Bugfixes
- **Menu redraw fix**: After returning from container shell, the background grid is now correctly redrawn and the container menu is refocused

---

### v0.8.4 (2026-02-28)

#### Bugfixes
- **Input widget width fix**: Commit dialog input field was too narrow (28 instead of 68 characters) — `SetMaxLen()` now correctly recalculates widget width
- **Exec shell fix**: Shell execution inside container overlapped with the TUI — terminal is now released before shell start (`ui.Close()`) and cleanly restored after exit (`ui.Init()`)

---

### v0.8.3 (2026-02-28)

#### New Features
- **Flicker fix**: Eliminated screen flicker during 1-second refresh — header and grid are now drawn in a single render call. Screen clear only triggered when container count shrinks
- **Stop + Remove**: New combined menu entry `[X] stop + remove` in container menu for running containers — stops and removes in one step with confirmation dialog
- **Docker Commit**: Save containers as local images via `[C] commit to image` in container menu — input dialog with name suggestion (`<name>-snapshot`) and optional tag (default: `latest`)

---

### v0.8.2 (2026-02-28)

#### New Features
- **Version display in header**: Version is shown in the top-right corner of the TUI header (e.g. `v0.8.2`)
- **Configurable shell**: Shell for `exec` is configurable via `--shell` flag, `CTOP_SHELL` environment variable, or config file (priority: CLI > Env > Config > Default `/bin/sh`)
- **Health status column**: New optional column shows health check status (`healthy`/`unhealthy`/`starting`) with color coding
- **Restart count column**: New optional column shows container restart count with color warnings (>0 yellow, >5 red)

#### Bugfixes
- Fixed header version positioning (no longer at the very edge of the terminal)

#### CI/CD
- **Homebrew tap auto-update**: CI job automatically updates the Homebrew formula after a release
- New Makefile target `make update-homebrew` to sync the Homebrew tap repo

---

### v0.8.0 (2025-06-24)

First release of the maintained fork of [bcicen/ctop](https://github.com/bcicen/ctop).

#### Security Fixes
- Fixed shell injection vulnerability in `container.Exec()` — using `exec.Command()` instead of `sh -c`
- Set config file permissions to `0600`
- Debug server now binds to `localhost` instead of `0.0.0.0`

#### Modernization
- Updated to **Go 1.22** with modern stdlib packages
- Migrated from `op/go-logging` to `log/slog` (structured logging)
- Replaced `nu7hatch/gouuid` with `google/uuid`
- Replaced `pkg/errors` with `fmt.Errorf` using `%w` wrapping
- Replaced `io/ioutil` with `os`/`io` functions
- Using `os.ReadDir` instead of `ioutil.ReadDir`

#### CI/CD & Infrastructure
- GitHub Actions workflow for automatic multi-platform builds (Linux, macOS, Windows — AMD64 + ARM64)
- Automatic GitHub Releases with binaries and checksums
- Docker image via `ghcr.io/eqms/ctop`

#### Documentation
- Bilingual README (German/English)
- Fork motivation and change overview documented
