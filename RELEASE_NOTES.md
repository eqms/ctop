# Release Notes

> **Language / Sprache**: [Deutsch](#deutsche-release-notes) | [English](#english-release-notes)

---

## Deutsche Release Notes

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
