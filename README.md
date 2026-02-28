<p align="center"><img width="200px" src="/_docs/img/logo.png" alt="ctop"/></p>

# ctop

> **Language / Sprache**: [Deutsch](#deutsche-dokumentation) | [English](#english-documentation)

Top-like interface for container metrics / Top-artige Oberfläche für Container-Metriken

<p align="center"><img src="_docs/img/grid.gif" alt="ctop"/></p>

---

## Deutsche Dokumentation

### Projektübersicht

`ctop` ist ein Terminal-basiertes Monitoring-Tool für Container. Es zeigt Echtzeit-Metriken (CPU, Speicher, Netzwerk, I/O) für mehrere Container in einer kompakten Übersicht sowie einer [Einzelansicht](_docs/single.md) für detaillierte Container-Informationen.

Unterstützte Container-Runtimes: **Docker** und **runC**.

#### Fork-Informationen

Dies ist ein gepflegter Fork von [bcicen/ctop](https://github.com/bcicen/ctop), ursprünglich erstellt von [VektorLab](https://github.com/vektorlab).

`ctop` ist für mich ein unverzichtbares Werkzeug im täglichen Umgang mit Containern — kompakt, schnell und auf den Punkt. Leider wird das Original-Projekt seit längerer Zeit nicht mehr aktiv weiterentwickelt. Da ich das Tool regelmäßig nutze und es nicht aufgeben möchte, pflege ich es hier als eigenständigen Fork weiter: mit aktuellen Abhängigkeiten, Security-Fixes und modernem Go.

Änderungen in diesem Fork:
- Security-Fixes (Shell-Injection, Dateiberechtigungen, Debug-Server-Binding)
- Aktualisierung auf Go 1.22 mit modernen stdlib-Paketen (`log/slog`, `os.ReadDir`)
- Ersetzen nicht gepflegter Dependencies (`op/go-logging`, `nu7hatch/gouuid`, `pkg/errors`)
- GitHub Actions CI/CD mit Multi-Plattform-Builds (Linux amd64/arm64, macOS amd64/arm64, Windows amd64)
- Docker-Images über GHCR (`ghcr.io/eqms/ctop`)

### Installation

#### Binaries herunterladen

Die neuesten Binaries gibt es auf der [Releases-Seite](https://github.com/eqms/ctop/releases):

| Plattform | Architektur | Datei |
|-----------|-------------|-------|
| Linux | amd64 | `ctop-<version>-linux-amd64` |
| Linux | arm64 | `ctop-<version>-linux-arm64` |
| macOS | Intel | `ctop-<version>-darwin-amd64` |
| macOS | Apple Silicon | `ctop-<version>-darwin-arm64` |
| Windows | amd64 | `ctop-<version>-windows-amd64.exe` |

```bash
# Linux (amd64)
sudo curl -Lo /usr/local/bin/ctop https://github.com/eqms/ctop/releases/latest/download/ctop-0.8.0-linux-amd64
sudo chmod +x /usr/local/bin/ctop

# macOS (Apple Silicon)
curl -Lo ctop https://github.com/eqms/ctop/releases/latest/download/ctop-0.8.0-darwin-arm64
chmod +x ctop
sudo mv ctop /usr/local/bin/

# macOS (Intel)
curl -Lo ctop https://github.com/eqms/ctop/releases/latest/download/ctop-0.8.0-darwin-amd64
chmod +x ctop
sudo mv ctop /usr/local/bin/
```

> **macOS Hinweis**: Beim Download über einen Browser setzt macOS das Quarantine-Flag, wodurch Gatekeeper die Ausführung blockiert. Abhilfe:
> ```bash
> # Quarantine-Flag entfernen
> xattr -d com.apple.quarantine ./ctop-0.8.0-darwin-arm64
> chmod +x ./ctop-0.8.0-darwin-arm64
> ```
> Bei Installation per `curl` (wie oben) wird kein Quarantine-Flag gesetzt.

#### Homebrew (macOS & Linux)

```bash
brew install eqms/ctop/ctop
```

#### Docker

```bash
docker run --rm -ti \
  --name=ctop \
  --volume /var/run/docker.sock:/var/run/docker.sock:ro \
  ghcr.io/eqms/ctop:latest
```

### Verwendung

`ctop` benötigt keine Argumente und nutzt standardmäßig die Docker-Host-Variablen. Weitere Konfigurationsoptionen unter [Connectors](_docs/connectors.md).

#### Optionen

| Option | Beschreibung |
|--------|-------------|
| `-a` | Nur aktive Container anzeigen |
| `-f <string>` | Initiale Filterzeichenkette setzen |
| `-h` | Hilfe anzeigen |
| `-i` | Standardfarben invertieren |
| `-r` | Container-Sortierreihenfolge umkehren |
| `-s` | Initiales Sortierfeld wählen |
| `-v` | Versionsinformationen ausgeben |
| `--connector` | Container-Connector wählen (`docker`, `runc`) |

#### Tastenbelegung

| Taste | Aktion |
|:-----:|--------|
| <kbd>Enter</kbd> | Container-Menü öffnen |
| <kbd>a</kbd> | Alle Container ein-/ausblenden |
| <kbd>f</kbd> | Container filtern (<kbd>Esc</kbd> zum Löschen) |
| <kbd>H</kbd> | Header ein-/ausblenden |
| <kbd>h</kbd> | Hilfe anzeigen |
| <kbd>s</kbd> | Sortierfeld wählen |
| <kbd>r</kbd> | Sortierreihenfolge umkehren |
| <kbd>o</kbd> | Einzelansicht öffnen |
| <kbd>l</kbd> | Container-Logs anzeigen (<kbd>t</kbd> für Zeitstempel) |
| <kbd>e</kbd> | Shell im Container öffnen |
| <kbd>c</kbd> | Spalten konfigurieren |
| <kbd>S</kbd> | Konfiguration in Datei speichern |
| <kbd>q</kbd> | Beenden |

#### Konfigurationsdatei

Während der Ausführung kann mit <kbd>S</kbd> die aktuelle Konfiguration gespeichert werden:
- XDG-Systeme: `~/.config/ctop/config`
- Fallback: `~/.ctop`

### Bauen aus Quellcode

```bash
git clone https://github.com/eqms/ctop.git
cd ctop
make build
```

#### Alle Plattformen

```bash
make build-all    # Binaries in _build/
```

#### Docker-Image

```bash
make image        # Erstellt lokales ctop:latest Image
```

### Architektur

```
ctop/
├── main.go                  # Einstiegspunkt, CLI-Flags
├── grid.go                  # Haupt-UI-Grid-Steuerung
├── cursor.go                # Grid-Cursor / Navigation
├── menus.go                 # TUI-Menüs (Container, Filter, Sortierung)
├── config/                  # Konfigurationsverwaltung (TOML)
├── connector/               # Container-Runtime-Abstraktionen
│   ├── docker.go            # Docker-Connector
│   ├── runc.go              # runC-Connector (Linux)
│   ├── mock.go              # Mock-Connector (Entwicklung)
│   ├── collector/           # Metrik-Sammler
│   └── manager/             # Container-Lifecycle-Management
├── container/               # Container-Modell und Sortierung
├── cwidgets/                # TUI-Widgets
│   ├── compact/             # Kompakte Grid-Ansicht
│   └── single/              # Einzelansicht (CPU, Mem, Net, I/O)
├── logging/                 # Logging (log/slog + Ring Buffer)
├── widgets/                 # Basis-UI-Widgets
├── Makefile                 # Build-Targets
├── Dockerfile               # Multi-Stage Docker Build
└── .github/workflows/ci.yml # GitHub Actions CI/CD
```

### Entwicklung

**Voraussetzungen**: Go 1.22+

```bash
# Entwicklungsmodus mit Debug-Server
make run-dev

# Debug-Optionen (Umgebungsvariablen)
CTOP_DEBUG=1          # Debug-Logging und Unix-Socket-Server aktivieren
CTOP_DEBUG_TCP=1      # TCP Debug-Server (127.0.0.1:9000) statt Unix-Socket
CTOP_DEBUG_FILE=/path # Log-Ausgabe zusätzlich in Datei schreiben
```

#### Connectors

| Connector | Umgebungsvariable | Standardwert |
|-----------|-------------------|--------------|
| Docker | `DOCKER_HOST` | `unix://var/run/docker.sock` |
| runC | `RUNC_ROOT` | `/run/runc` |
| runC | `RUNC_SYSTEMD_CGROUP` | (deaktiviert) |

### CI/CD

GitHub Actions Pipeline mit 4 Jobs:

1. **test** — `go test ./...`
2. **build** — Cross-Compile für 5 Plattformen
3. **release** — GitHub Release mit Binaries + SHA256 Checksums (bei Tag-Push `v*`)
4. **docker** — Multi-Arch Docker Image auf GHCR (bei Tag-Push `v*`)

### Lizenz

MIT License — Copyright (c) 2017 VektorLab. Siehe [LICENSE](LICENSE).

---

## English Documentation

### Project Overview

`ctop` is a terminal-based monitoring tool for containers. It displays real-time metrics (CPU, memory, network, I/O) for multiple containers in a compact overview as well as a [single container view](_docs/single.md) for detailed inspection.

Supported container runtimes: **Docker** and **runC**.

#### Fork Information

This is a maintained fork of [bcicen/ctop](https://github.com/bcicen/ctop), originally created by [VektorLab](https://github.com/vektorlab).

`ctop` is an indispensable tool in my daily work with containers — compact, fast and to the point. Unfortunately, the original project is no longer actively maintained. Since I use this tool on a regular basis and don't want to let it go, I maintain it here as an independent fork: with up-to-date dependencies, security fixes and modern Go.

Changes in this fork:
- Security fixes (shell injection, file permissions, debug server binding)
- Updated to Go 1.22 with modern stdlib packages (`log/slog`, `os.ReadDir`)
- Replaced unmaintained dependencies (`op/go-logging`, `nu7hatch/gouuid`, `pkg/errors`)
- GitHub Actions CI/CD with multi-platform builds (Linux amd64/arm64, macOS amd64/arm64, Windows amd64)
- Docker images published to GHCR (`ghcr.io/eqms/ctop`)

### Installation

#### Download Binaries

Get the latest binaries from the [Releases page](https://github.com/eqms/ctop/releases):

| Platform | Architecture | File |
|----------|-------------|------|
| Linux | amd64 | `ctop-<version>-linux-amd64` |
| Linux | arm64 | `ctop-<version>-linux-arm64` |
| macOS | Intel | `ctop-<version>-darwin-amd64` |
| macOS | Apple Silicon | `ctop-<version>-darwin-arm64` |
| Windows | amd64 | `ctop-<version>-windows-amd64.exe` |

```bash
# Linux (amd64)
sudo curl -Lo /usr/local/bin/ctop https://github.com/eqms/ctop/releases/latest/download/ctop-0.8.0-linux-amd64
sudo chmod +x /usr/local/bin/ctop

# macOS (Apple Silicon)
curl -Lo ctop https://github.com/eqms/ctop/releases/latest/download/ctop-0.8.0-darwin-arm64
chmod +x ctop
sudo mv ctop /usr/local/bin/

# macOS (Intel)
curl -Lo ctop https://github.com/eqms/ctop/releases/latest/download/ctop-0.8.0-darwin-amd64
chmod +x ctop
sudo mv ctop /usr/local/bin/
```

> **macOS Note**: When downloading via a browser, macOS sets the quarantine flag and Gatekeeper blocks execution. Fix:
> ```bash
> # Remove quarantine flag
> xattr -d com.apple.quarantine ./ctop-0.8.0-darwin-arm64
> chmod +x ./ctop-0.8.0-darwin-arm64
> ```
> When installing via `curl` (as shown above), no quarantine flag is set.

#### Homebrew (macOS & Linux)

```bash
brew install eqms/ctop/ctop
```

#### Docker

```bash
docker run --rm -ti \
  --name=ctop \
  --volume /var/run/docker.sock:/var/run/docker.sock:ro \
  ghcr.io/eqms/ctop:latest
```

### Usage

`ctop` requires no arguments and uses Docker host variables by default. See [connectors](_docs/connectors.md) for further configuration options.

#### Options

| Option | Description |
|--------|------------|
| `-a` | Show active containers only |
| `-f <string>` | Set an initial filter string |
| `-h` | Display help dialog |
| `-i` | Invert default colors |
| `-r` | Reverse container sort order |
| `-s` | Select initial container sort field |
| `-v` | Output version information and exit |
| `--connector` | Select container connector (`docker`, `runc`) |

#### Keybindings

| Key | Action |
|:---:|--------|
| <kbd>Enter</kbd> | Open container menu |
| <kbd>a</kbd> | Toggle display of all containers |
| <kbd>f</kbd> | Filter displayed containers (<kbd>Esc</kbd> to clear) |
| <kbd>H</kbd> | Toggle ctop header |
| <kbd>h</kbd> | Open help dialog |
| <kbd>s</kbd> | Select container sort field |
| <kbd>r</kbd> | Reverse container sort order |
| <kbd>o</kbd> | Open single view |
| <kbd>l</kbd> | View container logs (<kbd>t</kbd> to toggle timestamp) |
| <kbd>e</kbd> | Exec shell |
| <kbd>c</kbd> | Configure columns |
| <kbd>S</kbd> | Save current configuration to file |
| <kbd>q</kbd> | Quit ctop |

#### Config File

While running, press <kbd>S</kbd> to save the current configuration:
- XDG systems: `~/.config/ctop/config`
- Fallback: `~/.ctop`

### Building from Source

```bash
git clone https://github.com/eqms/ctop.git
cd ctop
make build
```

#### All Platforms

```bash
make build-all    # Binaries in _build/
```

#### Docker Image

```bash
make image        # Builds local ctop:latest image
```

### Architecture

```
ctop/
├── main.go                  # Entry point, CLI flags
├── grid.go                  # Main UI grid controller
├── cursor.go                # Grid cursor / navigation
├── menus.go                 # TUI menus (container, filter, sort)
├── config/                  # Configuration management (TOML)
├── connector/               # Container runtime abstractions
│   ├── docker.go            # Docker connector
│   ├── runc.go              # runC connector (Linux only)
│   ├── mock.go              # Mock connector (development)
│   ├── collector/           # Metrics collectors
│   └── manager/             # Container lifecycle management
├── container/               # Container model and sorting
├── cwidgets/                # TUI widgets
│   ├── compact/             # Compact grid view
│   └── single/              # Single view (CPU, Mem, Net, I/O)
├── logging/                 # Logging (log/slog + ring buffer)
├── widgets/                 # Base UI widgets
├── Makefile                 # Build targets
├── Dockerfile               # Multi-stage Docker build
└── .github/workflows/ci.yml # GitHub Actions CI/CD
```

### Development

**Prerequisites**: Go 1.22+

```bash
# Development mode with debug server
make run-dev

# Debug options (environment variables)
CTOP_DEBUG=1          # Enable debug logging and Unix socket server
CTOP_DEBUG_TCP=1      # Use TCP debug server (127.0.0.1:9000) instead of Unix socket
CTOP_DEBUG_FILE=/path # Additionally write log output to file
```

#### Connectors

| Connector | Environment Variable | Default |
|-----------|---------------------|---------|
| Docker | `DOCKER_HOST` | `unix://var/run/docker.sock` |
| runC | `RUNC_ROOT` | `/run/runc` |
| runC | `RUNC_SYSTEMD_CGROUP` | (disabled) |

### CI/CD

GitHub Actions pipeline with 4 jobs:

1. **test** — `go test ./...`
2. **build** — Cross-compile for 5 platforms
3. **release** — GitHub Release with binaries + SHA256 checksums (on tag push `v*`)
4. **docker** — Multi-arch Docker image to GHCR (on tag push `v*`)

### License

MIT License — Copyright (c) 2017 VektorLab. See [LICENSE](LICENSE).
