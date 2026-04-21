# Bubbletea Migration Plan

Replace the dead TUI dependency stack (`gizak/termui` v2 pre-release + `nsf/termbox-go` archived) with the actively maintained Charm ecosystem.

## Goals

- Keep ctop UX identical (key bindings, layout, colors, menus behave the same).
- End up with only maintained UI dependencies.
- Add unit tests for the new components while we rewrite them (current coverage: 1/60 files).

## Target stack

| Current | Target | Role |
|---|---|---|
| `nsf/termbox-go` | `charmbracelet/bubbletea` built-in (Bubble Tea handles its own terminal I/O via `charmbracelet/x/term`) | Terminal backend |
| `gizak/termui` primitives (`List`, `Par`, `Gauge`, `Block`) | `charmbracelet/lipgloss` + `charmbracelet/bubbles` (list, viewport, textinput, help, table, progress) | Widgets |
| `termui.Handle(...)` event dispatch | `tea.Model.Update(msg tea.Msg)` | Event loop |

## Scope — files to rewrite

Rendering / event loop:
- `main.go` — replace `ui.Init()/ui.Close()`, `ui.Render(...)`, event subscriptions with `tea.NewProgram(...).Run()`.
- `grid.go` — becomes the root `tea.Model` composing header, rows, status bar.
- `cursor.go` — becomes state on the root model (`selectedIdx int`).
- `menus.go` — each menu becomes a sub-model (`tea.Model`) pushed onto a stack.

Widgets:
- `widgets/header.go`, `widgets/status.go`, `widgets/error.go`, `widgets/input.go` — rewrite with `lipgloss` styles.
- `widgets/menu/` — replace with `bubbles/list` + `bubbles/help`.

Container widgets:
- `cwidgets/compact/` — one row per container; replace termui `Par/Gauge` composition with a `lipgloss` styled string per column. Column set is config-driven, do not hardcode.
- `cwidgets/single/` — detail/inspect view; replace with `bubbles/viewport` for scrolling, `lipgloss` for layout.

Not in scope:
- `connector/`, `container/`, `config/`, `models/`, `logging/` — no TUI imports, stay untouched.

## Migration strategy — incremental behind a build tag

Big-bang rewrite is risky; the current UI cannot be half-migrated because both systems grab the terminal.

Approach: introduce `internal/tui/` for the new implementation and keep the old code paths compiling under `-tags legacy_tui` until the rewrite is complete. `main.go` picks one at build time. Once `internal/tui/` is feature-complete, drop the build-tag, delete old files and `termui`/`termbox-go` from `go.mod`.

## Phased plan

### Phase A — skeleton (0.5 day)
1. `go get github.com/charmbracelet/bubbletea bubbles lipgloss`.
2. Create `internal/tui/app.go` with minimal `tea.Model` (Init/Update/View) that renders a blank screen and quits on `q`/Ctrl-C.
3. Wire `main.go` behind a build tag switch.
4. Verify `go build -tags experimental_tui ./...` on macOS + `GOOS=linux`.

### Phase B — layout shell (0.5 day)
Header with version, connector, container count; status bar showing filter + sort + mode; empty body in between. Use `lipgloss.JoinVertical` for the three rows, compute row heights from `tea.WindowSizeMsg`.

### Phase C — container list (1–1.5 days)
Most important piece. Read container state from `Connector.All()` on a 1s tick (`tea.Tick`). Render each container as one `lipgloss`-styled line with columns coming from `config.GlobalColumns`. Reuse the column ordering/width config already in `config/columns.go`. Cursor navigation: up/down/home/end, page up/down. **Test**: snapshot-test the rendering of a static container list (use `lipgloss.NewRenderer(&bytes.Buffer{})`).

### Phase D — menus, one at a time (2–3 days)
Order by complexity, simplest first so patterns solidify:
1. Help menu (static list)
2. Sort menu (`bubbles/list`)
3. Filter menu (`bubbles/textinput`)
4. Column config menu (checkbox list — `bubbles/list` with custom delegate)
5. Inspect menu (`bubbles/viewport`, scrollable)
6. Exec menu (text input + shell launcher — careful: `menus.go:437` `OpenInBrowser` port validation from the audit needs to be ported over)
7. Logs menu if it exists in the current build

Each menu is its own `tea.Model`; push/pop handled by a simple stack in the root model.

### Phase E — single/detail view (1 day)
Per-container metric view. Current implementation composes several termui widgets (gauges + text). With lipgloss: build a small dashboard using `bubbles/progress` for gauges and a viewport for log tail.

### Phase F — cleanup + drop legacy (0.5 day)
1. Remove build tag, delete old `grid.go`, `cursor.go`, `menus.go`, `widgets/`, `cwidgets/` (keep the Go files that are pure data/logic; only UI code goes).
2. `go mod edit -droprequire github.com/gizak/termui github.com/nsf/termbox-go github.com/mattn/go-runewidth` (runewidth may still be needed by lipgloss transitively — check `go mod tidy` output).
3. Remove `termbox-go` entries from `.golangci.yml` errcheck exclusions; adjust `composites` / `unkeyed fields` suppression — likely no longer needed.
4. Update `widgets/view_test.go` (the one existing test) or port it.

## Estimate

5–7 working days focused work. Smaller chunks possible but each phase should land as its own PR so UX regressions are easy to bisect.

## Risks

- **Color palette drift** — termui uses attribute codes, lipgloss uses ANSI 256/truecolor. Some colors may need manual re-tuning on low-color terminals.
- **Keyboard handling differences** — termbox is raw, Bubble Tea normalises some keys. Special keys like Enter in filter input need testing in tmux/screen.
- **Terminal resize behaviour** — verify under tmux split, iTerm, plain Linux console, and Windows Terminal (Windows build is in the release matrix).
- **Unicode rendering** — container names with non-ASCII chars. lipgloss uses `runewidth` for width calc; should be at least as correct as termui.

## Out of scope for Phase 3

- Column config file migration (existing `.ctop` config format stays).
- New features. Only UX parity.
- `widgets/view_test.go` — port in Phase F, do not expand.

## Open questions (resolve before starting)

1. Should the windows build stay in the release matrix? Bubble Tea supports Windows, but termbox-go's Windows path was always flaky; check the current Windows release for real users.
2. Keep `htop`-style cursor highlight or switch to the cleaner bubbles/list inverse-video selection?
3. Minimum terminal width — currently ~80 cols. Stay there or relax.
