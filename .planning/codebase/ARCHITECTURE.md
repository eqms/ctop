<!-- refreshed: 2026-06-11 -->
# Architecture

**Analysis Date:** 2026-06-11

## System Overview

```text
┌─────────────────────────────────────────────────────────────────┐
│                        main package                              │
│  `main.go` · `grid.go` · `cursor.go` · `menus.go` · `keys.go`  │
└─────┬───────────────────┬──────────────────────┬────────────────┘
      │                   │                      │
      ▼                   ▼                      ▼
┌──────────────┐  ┌───────────────┐   ┌──────────────────────────┐
│   connector  │  │   config      │   │   cwidgets / widgets     │
│  `connector/`│  │  `config/`    │   │  `cwidgets/` `widgets/`  │
│  ─────────── │  │  ─────────── │   │  ──────────────────────  │
│  Docker      │  │  Singleton   │   │  CompactGrid, header,    │
│  runC        │  │  TOML file   │   │  menus, status, single   │
│  Mock        │  │  CLI overrides│   │  view                    │
└──────┬───────┘  └───────────────┘   └──────────────────────────┘
       │
       ▼
┌───────────────────────────────────────────────────────────────┐
│            connector sub-packages                              │
│  `connector/collector/` — metrics stream (Collector iface)    │
│  `connector/manager/`   — lifecycle ops (Manager iface)       │
└───────────────────────────────────┬───────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────┐
│  container · models                                            │
│  `container/main.go`  — Container struct (owns Collector,     │
│                          Manager, CompactRow, Meta, Metrics)   │
│  `models/main.go`     — Metrics struct, Meta map, Log type     │
└───────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

| Component | Responsibility | File |
|-----------|----------------|------|
| main | Init, CLI flags, global vars, render loop entry, shutdown | `main.go` |
| Display / grid | termui event loop, 1 s timer, keyboard dispatch | `grid.go` |
| GridCursor | Filtered container list, scroll/pagination, highlight | `cursor.go` |
| menus | Overlay menu windows (help, filter, sort, log, exec, columns) | `menus.go` |
| keys | Common key binding map, `HandleKeys` helper | `keys.go` |
| colors | Color map + invert logic | `colors.go` |
| connector.ConnectorSuper | Wraps backend, reconnect goroutine, thread-safe Get() | `connector/main.go` |
| connector.Docker | Docker event watcher, refresh goroutine, `All()`/`Get()` | `connector/docker.go` |
| connector.runC | runC libcontainer backend (Linux only) | `connector/runc.go` |
| connector.Mock | Synthetic random data for development | `connector/mock.go` |
| collector.Collector | Streams `chan models.Metrics`; provides `LogCollector` | `connector/collector/main.go` |
| manager.Manager | Container lifecycle: Start/Stop/Pause/Remove/Exec/Commit | `connector/manager/main.go` |
| container.Container | Aggregate model: ties Meta, Metrics, Widgets, Collector, Manager | `container/main.go` |
| container.Containers | Sortable/filterable slice; `Sort()`, `Filter()`, `Sorters` map | `container/sort.go` |
| models | Value types: `Metrics`, `Meta`, `Log` | `models/main.go` |
| config | Global singleton: params, switches, columns; TOML read/write | `config/main.go`, `config/columns.go` |
| cwidgets.WidgetUpdater | Interface: `SetMeta(Meta)` / `SetMetrics(Metrics)` | `cwidgets/main.go` |
| compact.CompactGrid | Paged grid of CompactRows with header | `cwidgets/compact/grid.go` |
| compact.CompactRow | One container row: slice of `CompactCol` + row background | `cwidgets/compact/row.go` |
| compact.CompactCol | Column interface + `allCols` registry of 16 column types | `cwidgets/compact/column.go` |
| widgets | Header bar, status line, error overlay, single-container view | `widgets/` |
| logging | Ring-buffer logger, status queue, debug HTTP server | `logging/` |

## Pattern Overview

**Overall:** Event-driven TUI with a pull-refresh loop, layered connectors, and a widget registry

**Key Characteristics:**
- Package-level global singletons for all UI state (`cursor`, `cGrid`, `header`, `status`, `errView`, `log` in `main.go`)
- Connector backends register via `init()` into a `map[string]ConnectorFn`; the UI only talks to `ConnectorSuper`
- Metrics flow one-way: backend → `chan models.Metrics` → `Container.Read()` goroutine → `WidgetUpdater.SetMetrics()`
- Column set is runtime-configurable; `allCols` registry in `cwidgets/compact/column.go` maps name → constructor

## Layers

**UI / Presentation Layer:**
- Purpose: Render terminal widgets, handle keyboard events, display menus
- Location: `main.go`, `grid.go`, `cursor.go`, `menus.go`, `keys.go`, `colors.go`, `cwidgets/`, `widgets/`
- Contains: termui widget wrappers, event handlers, layout math
- Depends on: `config`, `connector`, `container`, `cwidgets`, `models`
- Used by: nothing (top of stack)

**Configuration Layer:**
- Purpose: Single source of truth for params, switches, column order; reads/writes TOML
- Location: `config/`
- Contains: `GlobalParams`, `GlobalSwitches`, `GlobalColumns`; `Read()`, `Write()`, `Update()`, `Toggle()`
- Depends on: `logging`
- Used by: all layers

**Container Model Layer:**
- Purpose: Aggregate container state; owns lifecycle, metrics routing, widget reference
- Location: `container/`
- Contains: `Container` struct, `Containers` slice type, `Sorters` map, `Filter()`
- Depends on: `connector/collector`, `connector/manager`, `cwidgets`, `cwidgets/compact`, `models`
- Used by: connector backends, UI layer

**Connector Layer:**
- Purpose: Abstract backend (Docker/runC/Mock) with reconnect resilience
- Location: `connector/`
- Contains: `Connector` interface, `ConnectorSuper`, backend impls, `enabled` registry
- Depends on: `container`, `connector/collector`, `connector/manager`, `models`
- Used by: UI layer (`cursor.go` holds `*ConnectorSuper`)

**Collector Sub-layer:**
- Purpose: Per-container metrics and log streaming
- Location: `connector/collector/`
- Contains: `Collector` interface, Docker/runC/Mock/proc implementations, `LogCollector` interface
- Depends on: `models`, external cgroup/Docker APIs
- Used by: `container.Container`, connector backends

**Manager Sub-layer:**
- Purpose: Container lifecycle operations
- Location: `connector/manager/`
- Contains: `Manager` interface, Docker/runC/Mock implementations
- Depends on: Docker API, libcontainer
- Used by: `container.Container`

**Models Layer:**
- Purpose: Pure value types shared across all layers
- Location: `models/`
- Contains: `Metrics`, `Meta`, `Log`
- Depends on: nothing
- Used by: all layers

## Data Flow

### Primary Render Path (1 s tick)

1. termui fires `/timer/1s` → `Display()` in `grid.go`
2. `RefreshDisplay()` calls `cursor.RefreshContainers()`
3. `GridCursor.RefreshContainers()` calls `ConnectorSuper.Get()` → `Connector.All()`
4. `All()` returns sorted+filtered `container.Containers` slice
5. `RedrawRows()` iterates `cursor.filtered`, adds `c.Widgets` (each a `*compact.CompactRow`) to `cGrid`
6. `ui.Render(header, cGrid)` pushes all buffers to termbox

### Metrics Update Path (background goroutine per container)

1. When `Container.SetState("running")` is called → `collector.Start()`
2. `Container.Read(collector.Stream())` spawns goroutine consuming `chan models.Metrics`
3. Each `models.Metrics` value is delivered to `updater.SetMetrics(metrics)` (the `CompactRow`)
4. `CompactRow.SetMetrics()` fans out to all `CompactCol.SetMetrics()` — columns update their internal state
5. Next render tick picks up updated column state via `CompactRow.Buffer()` → `CompactCol.Buffer()`

### Docker Event Path

1. `connector.Docker.watchEvents()` goroutine receives Docker API events
2. `create` / lifecycle events → `cm.needsRefresh` channel → `Loop()` goroutine calls `cm.refresh(c)`
3. `health_status` / state changes → `cm.statuses` channel → `LoopStatuses()` calls `c.SetMeta` / `c.SetState`

**State Management:**
- Container metadata (`Meta`) is a `map[string]string` mutated via `Container.SetMeta()`
- Container metrics (`Metrics`) are replaced atomically per tick via `Container.Read()` goroutine
- Display config is in the `config` singleton behind a `sync.RWMutex`
- `GridCursor.filtered` is rebuilt from scratch on every `RefreshContainers()` call

## Key Abstractions

**`Connector` interface:**
- Purpose: Backend-agnostic container source
- Defined: `connector/main.go:23`
- Implementations: `connector/docker.go`, `connector/runc.go`, `connector/mock.go`
- Pattern: `init()` registration into `enabled map[string]ConnectorFn`

**`Collector` interface:**
- Purpose: Per-container metrics and log stream
- Defined: `connector/collector/main.go:14`
- Implementations: `connector/collector/docker.go`, `connector/collector/runc.go`, `connector/collector/mock.go`, `connector/collector/proc.go`
- Pattern: `Start()` begins streaming to `chan models.Metrics`; callers call `Stream()` to get channel

**`Manager` interface:**
- Purpose: Container lifecycle control
- Defined: `connector/manager/main.go:5`
- Implementations: `connector/manager/docker.go`, `connector/manager/runc.go`, `connector/manager/mock.go`

**`CompactCol` interface:**
- Purpose: One column in the compact container list
- Defined: `cwidgets/compact/column.go:50`
- Registry: `allCols map[string]NewCompactColFn` — 16 registered column types
- Pattern: Add new column → implement `CompactCol`, register in `allCols`, add to `config/columns.go` defaults

**`WidgetUpdater` interface:**
- Purpose: Decouples metric delivery from widget implementation
- Defined: `cwidgets/main.go`
- Used by: `container.Container.updater` — swapped between `CompactRow` (list view) and `single.Single` (detail view)

## Entry Points

**Binary entry:**
- Location: `main.go:main()`
- Triggers: CLI invocation `ctop [flags]`
- Responsibilities: parse flags, init config, init termui, create ConnectorSuper, enter `Display()` loop

**Render loop:**
- Location: `grid.go:Display()`
- Triggers: called from `main()` in a `for` loop; re-entered after each menu/action returns
- Responsibilities: register all keyboard handlers, start termui event loop, dispatch to menus

## Architectural Constraints

- **Threading:** Main goroutine owns all termui renders. Background goroutines per connector (`Loop`, `LoopStatuses`, `watchEvents`) and per container (`Read` metrics goroutine). Race-free via `sync.RWMutex` in `connector.Docker` and `config` packages.
- **Global state:** Package-level globals in `main.go` (`cursor`, `cGrid`, `header`, `status`, `errView`, `log`). Config singleton in `config/main.go` (`GlobalParams`, `GlobalSwitches`, `GlobalColumns`). Logger singleton in `logging/main.go` (`Log`). Connector registry `enabled` in `connector/main.go`.
- **Circular imports:** None detected. Dependency direction is strictly: `main` → `connector`/`cwidgets`/`widgets`/`config` → `container` → `models`. Config is read by all but depends only on `logging`.
- **Linux-only files:** `connector/runc.go` and `connector/collector/runc.go` use `//go:build linux`. Always verify with `GOOS=linux go build ./...` after editing them.
- **UI libraries:** Both `github.com/gizak/termui` and `github.com/nsf/termbox-go` are required together — termui is the widget layer, termbox-go is the raw terminal backend. Both are abandoned upstream. A Bubble Tea migration plan exists at `_docs/bubbletea-migration-plan.md` but has not started; do not introduce Bubble Tea piecemeal.

## Anti-Patterns

### Mixing termui and termbox directly

**What happens:** Both `ui "github.com/gizak/termui"` and `tm "github.com/nsf/termbox-go"` are imported in the same files (`main.go`, `grid.go`, `menus.go`). `tm.Clear()` is called before every menu overlay; `ui.Render()` is used for widget output.
**Why it's wrong:** The two libraries share terminal state; using them inconsistently can cause flicker and rendering artifacts.
**Do this instead:** Follow the existing pattern — call `tm.Clear(tm.ColorDefault, tm.ColorDefault)` before overlay menus (as in `grid.go:ShowConnError`, `menus.go`), then `ui.Render()` for widgets. Do not call termbox functions inside widget `Buffer()` methods.

### Package-level global UI singletons

**What happens:** `cursor`, `cGrid`, `header`, `status`, `errView` are `var` globals in `main.go`, mutated by `grid.go` and `menus.go`.
**Why it's wrong:** Makes unit testing impossible and hides dependencies; any file in `package main` can mutate them.
**Do this instead:** Until the Bubble Tea migration, continue using the established pattern — pass these values as parameters if extracting new functions, rather than adding more package-level globals.

## Error Handling

**Strategy:** Errors surface to the UI via the status line or the error overlay; the connector layer retries automatically.

**Patterns:**
- Connector errors: `ConnectorSuper.loop()` retries every 3 s; `Get()` returns `(nil, err)` which propagates to `RefreshDisplay()` → `ShowConnError()` overlay in `grid.go`
- Container operation errors (start/stop/remove etc.): logged via `log.Warningf` and pushed to the status queue via `log.StatusErr(err)` for display in the status line
- Panic recovery: `main()` defers `panicExit()` which calls `Shutdown()` then `os.Exit(1)`

## Cross-Cutting Concerns

**Logging:** `logging.Init()` returns the `*CTopLogger` singleton. All packages call `logging.Init()` and store the result in a package-level `var log`. The logger uses a ring buffer (1024 entries) and a status message queue for TUI-visible messages. Debug HTTP server enabled when `CTOP_DEBUG=1`.

**Validation:** Container sort field validated at startup in `main.go:validSort()`. Column names validated at row-widget construction time in `cwidgets/compact/column.go:newRowWidgets()` (panic on unknown name).

**Authentication:** None — ctop connects to local Docker socket or runC root dir; no auth layer.

---

*Architecture analysis: 2026-06-11*
