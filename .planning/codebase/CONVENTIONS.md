# Coding Conventions

**Analysis Date:** 2026-06-11

## Naming Patterns

**Files:**
- Flat names: `main.go`, `sort.go`, `file.go`, `param.go` — one concern per file
- Named after the type they export: `column.go` exports `CompactCol` interface, `row.go` exports `CompactRow`
- Test files co-located with suffix `_test.go`: `widgets/view_test.go`

**Functions / Methods:**
- Constructors always named `New<Type>`: `NewDocker()`, `NewCompactRow()`, `NewConnectorSuper()`, `NewMeta()`
- Exported functions use PascalCase: `RefreshContainers()`, `EnabledColumns()`, `SortFields()`
- Unexported functions use camelCase: `swapCols()`, `colIndex()`, `sumNet()`, `panicExit()`
- Method receivers use two-letter abbreviations matching the type: `(cm *Docker)`, `(gc *GridCursor)`, `(cs *ConnectorSuper)`, `(dc *Docker)`, `(c *Container)`, `(rb *ringBuffer)`

**Variables:**
- Package-level singletons use unexported camelCase: `log`, `enabled`, `lock`
- Exported singletons use PascalCase: `GlobalParams`, `GlobalSwitches`, `GlobalColumns`, `Log`, `Sorters`, `ErrActionNotImpl`
- Short, contextual names in methods: `idx`, `c`, `cs`, `err`, `n`

**Types:**
- Interfaces named after capability, not implementation: `Connector`, `Collector`, `Manager`, `LogCollector`, `CompactCol`, `ToggleText`
- Concrete types match their backing concept: `Docker`, `ConnectorSuper`, `GridCursor`, `CompactRow`, `CTopLogger`
- Type aliases for slices: `type Containers []*Container`, `type Meta map[string]string`
- Function types: `type sortMethod func(c1, c2 *Container) bool`, `type NewCompactColFn func() CompactCol`, `type ConnectorFn func() (Connector, error)`

**Constants:**
- Unexported constants use camelCase: `running = "running"`, `rowPadding = 1`, `ringSize = 1024`

## Code Style

**Formatting:**
- Standard `gofmt` — enforced implicitly by `golangci-lint`
- No additional formatter configured

**Linting:**
- `golangci-lint` v2.11 (v2 config: `version: "2"` in `.golangci.yml`)
- Config file: `.golangci.yml`
- Key suppressed linters:
  - `errcheck` suppressed for specific termbox functions (`Clear`, `Sync`) and browser URL opening
  - `govet/composites` disabled globally
  - `staticcheck QF1008` suppressed globally
  - `errcheck` suppressed on `(grid|menus|cursor).go` for UI refresh calls
- Run locally: `golangci-lint run --timeout=5m`

## Import Organization

**Order (idiomatic Go):**
1. Standard library
2. Blank line
3. External modules
4. Blank line (optional)
5. Internal packages (`github.com/eqms/ctop/...`)

**Import Aliases:**
- `ui "github.com/gizak/termui"` — widget library aliased as `ui`
- `tm "github.com/nsf/termbox-go"` — terminal backend aliased as `tm`
- `api "github.com/fsouza/go-dockerclient"` — Docker client aliased as `api`

**Example from `main.go`:**
```go
import (
    "flag"
    "fmt"
    "os"

    "github.com/eqms/ctop/config"
    "github.com/eqms/ctop/connector"
    ui "github.com/gizak/termui"
    tm "github.com/nsf/termbox-go"
)
```

## Build Tags

**Dual constraint format** — both modern and legacy forms are present on Linux-only files (leave both when editing):
```go
//go:build linux
// +build linux
```

**Tags in use:**
- `linux` — `connector/runc.go`, `connector/collector/runc.go`, `connector/collector/proc.go`
- `!release` — `connector/mock.go`, `connector/collector/mock.go` (mock backends excluded from release binary)
- `release` — passed via `CGO_ENABLED=0 go build -tags release` in `Makefile`

## Error Handling

**Strategies used (from most to least common):**

1. **Log and return** — used in `container/main.go` for container operations:
   ```go
   if err := c.manager.Start(); err != nil {
       log.Warningf("container %s: %v", c.Id, err)
       log.StatusErr(err)
       return
   }
   ```

2. **Return wrapped error** — used in `config/file.go` and `connector/manager/docker.go`:
   ```go
   return fmt.Errorf("cannot start container: %v", err)
   ```

3. **Panic for unrecoverable startup** — used in `main.go` and for programmer errors:
   ```go
   if err := ui.Init(); err != nil {
       panic(err)
   }
   panic(fmt.Sprintf("no such widget name: %s", name)) // invariant violation
   ```
   All panics in `main()` are caught by `defer panicExit()` which calls `Shutdown()` before printing to stderr and calling `os.Exit(1)`.

4. **Log and continue** — Docker event loop uses `log.Errorf(...)` without returning or panicking.

**Sentinel errors:**
- `connector/manager/main.go` defines `var ErrActionNotImpl = errors.New("action not implemented")` for unimplemented manager methods.

## Logging

**Framework:** Custom `CTopLogger` wrapping `log/slog` (stdlib) — package at `logging/`

**Logger initialization pattern** — package-level singleton obtained via `logging.Init()`:
```go
var log = logging.Init() // at top of each package needing logs
```
This is used in: `connector/main.go`, `connector/collector/main.go`, `cwidgets/compact/row.go`, `container/main.go`, `config/main.go`

**Log levels used:**
- `log.Debugf(...)` — verbose diagnostic, only emitted when `CTOP_DEBUG=1`
- `log.Infof(...)` / `log.Info(...)` — lifecycle events (start/stop collectors, connector state)
- `log.Noticef(...)` — config changes
- `log.Warningf(...)` — recoverable errors on container operations
- `log.Errorf(...)` — connector/collector failures

**Status messages** (displayed in the TUI status bar):
```go
log.Status("committed foo as repo:tag")  // informational
log.StatusErr(err)                        // error in TUI
log.Statusf("committed %s as %s:%s", ...) // formatted
```

**Output:** Ring buffer (1024 entries) viewable via debug log server; optionally written to file via `CTOP_DEBUG_FILE` env var.

## Comments

**Style:**
- Doc comments on exported functions and types: single-line `//` above the declaration
- Inline comments for non-obvious logic, on the same line or the line above
- No JSDoc-style block comments; no `/* */` except possibly in legacy code

**Examples:**
```go
// Refresh containers from source, returning whether the quantity of
// containers has changed and any error
func (gc *GridCursor) RefreshContainers() (bool, error) {

// return index of column with given name, if any
func colIndex(name string) int {

Field  string // "status" or "health"
```

## Function Design

**Size:** Functions are concise — most are 10–30 lines. Long functions (e.g., `main()`, Docker event loop) are broken up with inline comments as section markers.

**Parameters:** Prefer few parameters; pass context via receiver or package-level config singleton rather than threading state through function args.

**Return Values:**
- Errors returned as last return value: `(Connector, error)`, `(bool, error)`, `(string, error)`
- Boolean + error for "did state change" + reason: `RefreshContainers() (bool, error)`
- Methods with no meaningful return often use side effects + logging

## Module Design

**Module path:** `github.com/eqms/ctop`

**Package structure:** One package per directory. Subpackages for sub-concerns:
- `connector/collector` — metric collection per backend
- `connector/manager` — lifecycle management per backend
- `cwidgets/compact` — compact row widget family
- `cwidgets/single` — single-container view widgets

**Exports:** Each package exports only the interface and types needed by callers; constructors (`New*`) and interface types are always exported. Implementation details remain unexported.

**Barrel files:** Not used. Each file exports directly from its package.

**Backend registration pattern via `init()`:**
```go
// connector/docker.go
func init() { enabled["docker"] = NewDocker }

// connector/runc.go
func init() { enabled["runc"] = NewRunC }
```

## Commit Message Convention

Prefix required, enforced by release-notes parsing:
- `[ADD]` — new feature
- `[CHG]` — modification or change
- `[FIX]` — bug fix

---

*Convention analysis: 2026-06-11*
