# Release Notes

> **Language / Sprache**: [Deutsch](#deutsche-release-notes) | [English](#english-release-notes)

---

## Deutsche Release Notes

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
