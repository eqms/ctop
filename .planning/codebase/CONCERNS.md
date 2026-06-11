# Codebase Concerns

**Analysis Date:** 2026-06-11

## Tech Debt

**Dead UI library dependency stack:**
- Issue: Both TUI libraries are abandoned upstream. `github.com/gizak/termui` is pinned to a 2018 pre-release pseudo-version (`v2.3.1-0.20180817033724-8d4faad06196+incompatible`). `github.com/nsf/termbox-go` is pinned to a 2019 snapshot. Neither repository receives security fixes or Go version compatibility updates.
- Files: `go.mod:7-12`, `main.go`, `grid.go`, `cursor.go`, `menus.go`, `cwidgets/compact/row.go`, `widgets/*.go`
- Impact: Any future Go version incompatibility or CVE in the terminal backend has no upstream fix path. The entire TUI layer — roughly 20 files — is coupled to these two libraries.
- Fix approach: Documented migration plan exists at `_docs/bubbletea-migration-plan.md`. 6 phases (A–F), estimated 5–7 working days. Uses build tag `experimental_tui` to migrate incrementally without breaking the existing build.

**Stale small dependencies without upstream activity:**
- Issue: `github.com/c9s/goprocinfo` pinned to a 2017 commit (`v0.0.0-20170609001544-b34328d6e0cd`); `github.com/jgautheron/codename-generator` pinned to a 2015 commit; `github.com/hako/durafmt` last touched 2021; `github.com/pkg/browser` pinned to 2020.
- Files: `go.mod:5,9,10,15`, `connector/collector/proc.go`, `connector/mock.go`, `connector/docker.go`, `menus.go`
- Impact: `goprocinfo` reads `/proc/stat` and `/proc/meminfo` on Linux using an unmaintained parser — any kernel-side format changes will silently produce wrong values. The others carry lower risk but increase supply-chain surface.
- Fix approach: Replace `goprocinfo` with direct `/proc` parsing or `golang.org/x/sys/unix`; evaluate whether `durafmt` can be replaced with stdlib `time.Duration.String()`; replace `codename-generator` with a short local implementation.

**`Container.Meta` map has no lock:**
- Issue: `models.Meta` is a plain `map[string]string`. `container.SetMeta` writes to it from goroutines (Docker event handler via `LoopStatuses`, refresh loop via `Loop`), while the TUI main goroutine reads `c.Meta["state"]`, `c.Meta["Web Port"]` etc. directly at `menus.go:222-288` and `container/main.go:102-172`. There is no `sync.RWMutex` protecting this map.
- Files: `models/main.go:10`, `container/main.go:61`, `menus.go:222`, `connector/docker.go:256-262`
- Impact: Data race. Detectable with `go test -race`. In practice it is latent because the TUI runs in a tight 1-second timer loop, but concurrent map writes will cause a panic under Go's race detector and occasionally in production.
- Fix approach: Add a `sync.RWMutex` to `container.Container`, use it in `SetMeta`/`GetMeta`, and replace bare `c.Meta["key"]` reads in `menus.go` with `c.GetMeta("key")`.

**`container.Display` field written without synchronisation:**
- Issue: `container.Display` is a plain `bool` on the `Container` struct. It is written by `Containers.Filter()` (called from `All()` which runs in the connector goroutine) and read from the TUI goroutine in `cursor.go:42`. No lock protects it.
- Files: `container/sort.go:118-125`, `cursor.go:42`
- Impact: Benign data race in most cases, but triggers the Go race detector and is technically undefined behaviour.
- Fix approach: Protect with the same mutex added for `Meta`, or change `Filter()` to return a new filtered slice rather than mutating `Display` in place (preferred: avoids shared mutable state entirely).

**`menus.go` is a 600-line mega-file:**
- Issue: All menu implementations (Help, Filter, Sort, Columns, Container actions, Log, ExecShell, Commit, Confirm, OpenInBrowser) plus helper functions live in a single 600-line file. This makes navigating and testing individual menus difficult.
- Files: `menus.go`
- Impact: Low immediate risk; slows down the Bubble Tea migration (each menu must be extracted in Phase D).
- Fix approach: Split into `menus_container.go`, `menus_filter.go`, `menus_sort.go`, etc. before starting Phase D of the migration plan.

## Known Bugs

**`regexp.MustCompile` panics on invalid filter input:**
- Symptoms: If a user types an invalid regex character sequence (e.g., `[unclosed`) in the filter box, `MustCompile` panics. The `panicExit()` recovery in `main.go` catches it, but terminates the program.
- Files: `container/sort.go:115`
- Trigger: Open filter menu (`f`), type an invalid regex, press Enter.
- Workaround: Use `regexp.Compile` with an error return; fall back to literal string match on error.

**runc connector `ReadIO` is never called:**
- Symptoms: I/O bytes read/write metrics are always 0 (or initialised to -1) for runc containers. The `ReadIO` method exists at `connector/collector/runc.go:122` but `run()` at line 62–85 only calls `ReadCPU`, `ReadMem`, and `ReadNet`.
- Files: `connector/collector/runc.go:62-85`
- Trigger: Run ctop against a runc backend; observe I/O column always empty.
- Workaround: None. Fix: add `c.ReadIO(stats.CgroupStats)` after `c.ReadMem` in the `run()` loop.

**runc `GetLibc` does not persist newly loaded containers to `libContainers` map:**
- Symptoms: Every call to `GetLibc` for a container that is not yet in `libContainers` calls `libcontainer.Load()` again on the next invocation (after the container is added to `containers` via `MustGet`), because `GetLibc` returns early before the `libContainers[id] = libc` assignment that only happens inside `MustGet`. If `GetLibc` is called on the same ID from `refresh()` and `refreshAll()` concurrently it loads the container twice.
- Files: `connector/runc.go:90-109`
- Trigger: High-frequency runc container churn on large hosts.
- Workaround: None. Fix: store the loaded container in `libContainers` inside `GetLibc` when the load succeeds.

**`Docker.LoopStatuses` can silently drop status updates:**
- Symptoms: When a status update arrives for a container ID not yet in the map, `cm.Get(statusUpdate.Cid)` returns `nil` and the update is discarded. If events arrive before `refreshAll()` completes (race between `watchEvents` goroutine and `refreshAll` in `NewDocker`), the container's initial status is wrong.
- Files: `connector/docker.go:252-268`
- Trigger: Container starts very quickly; Docker fires a `start` event before `refreshAll` has inserted the container.
- Workaround: None. Fix: call `cm.MustGet` instead of `cm.Get` in `LoopStatuses` so the container is created if missing.

**Log pipe goroutine never closed on viewer exit:**
- Symptoms: When the log viewer is closed (`LogMenu` exits), the goroutine running `bufio.NewScanner(r)` in `docker_logs.go:49-58` blocks on `scanner.Scan()` indefinitely because the `io.Pipe` read end is not explicitly closed.
- Files: `connector/collector/docker_logs.go:47-60`
- Trigger: Open log view, close it; one goroutine leaks per log view session.
- Workaround: None. Fix: close the `io.Pipe` reader when the context is cancelled.

## Security Considerations

**Shell whitelist is hardcoded and excludes common non-standard paths:**
- Risk: `isAllowedShell` at `menus.go:591-599` rejects shells at paths like `/usr/local/bin/nu`, `/opt/homebrew/bin/fish`, or container-specific paths. Users with those shells configured see a confusing error with no fallback.
- Files: `menus.go:590-599`
- Current mitigation: Whitelist prevents arbitrary command injection via the `--shell` flag.
- Recommendations: Replace the whitelist with path validation (must be absolute path, no shell metacharacters, optionally verify the binary exists inside the container). The `menus.go:66` comment in the migration plan notes this validation must be ported to Phase D.

**`OpenInBrowser` constructs a URL from unvalidated container metadata:**
- Risk: `webPort` value comes from Docker port bindings (`connector/docker.go:147-162`). If a container publishes a crafted `HostIP:HostPort` value (e.g., containing `@attacker.com`), the `"http://" + webPort + "/"` concatenation produces a malformed or unintended URL that is passed to `browser.OpenURL`.
- Files: `menus.go:437`, `connector/docker.go:147-162`
- Current mitigation: Only running containers under user control are shown; threat model is local.
- Recommendations: Validate `webPort` against `host:port` format before constructing the URL (use `net.SplitHostPort`).

**`govulncheck` runs with `continue-on-error: true`:**
- Risk: Vulnerability scan failures are logged as CI warnings only — they do not block merges or releases.
- Files: `.github/workflows/ci.yml:40`
- Current mitigation: Scan does run; results are visible in CI annotations.
- Recommendations: Promote to a blocking step or route results to a dedicated security review workflow. At minimum, track known-acceptable failures in a `govulncheck` ignore list.

## Performance Bottlenecks

**`Containers.Filter()` recompiles the regex on every 1-second tick:**
- Problem: `regexp.MustCompile` is called inside `Filter()` which is called by `Connector.All()` which is called by `cursor.RefreshContainers()` every 1 second.
- Files: `container/sort.go:115`, `cursor.go:37`, `grid.go:117`
- Cause: No caching of the compiled regex; every tick allocates and compiles.
- Improvement path: Cache the compiled regex in `config` (invalidate on filter string change), or pre-compile in `cursor.RefreshContainers`.

**`cursor.Idx()` is O(n) on every scroll event:**
- Problem: `Idx()` iterates the full filtered container slice to find the selected container index on every Up/Down/PgUp/PgDown call.
- Files: `cursor.go:76-84`
- Cause: Uses ID string lookup in a slice rather than maintaining an integer index.
- Improvement path: Track `selectedIdx int` directly; update on Up/Down instead of recomputing.

## Fragile Areas

**`CompactRow.Highlight()` assumes `Cols[1]` always exists:**
- Files: `cwidgets/compact/row.go:98`
- Why fragile: `row.Cols[1].Highlight()` panics with an index-out-of-range if the user disables all columns except one via the Columns menu. There is no bounds check.
- Safe modification: Add `if len(row.Cols) > 1` guard before indexing, or always keep a minimum of two columns enabled in `config/columns.go`.
- Test coverage: No test for column enable/disable state in rendering.

**`Runc.refreshAll()` closes `closed` channel on `ReadDir` error:**
- Files: `connector/runc.go:147-151`
- Why fragile: Any transient permission or I/O error on `/run/runc` causes `close(cm.closed)` to be called, which signals `Wait()` and triggers a full reconnect loop in `ConnectorSuper`. The channel is also closed from the refresh goroutine, so a second error (race between the ticker and `refreshAll`) calls `close` on an already-closed channel, which panics.
- Safe modification: Use `sync.Once` to guard `close(cm.closed)`.
- Test coverage: None.

**`runc.Loop()` goroutine never exits:**
- Files: `connector/runc.go:181-185`
- Why fragile: `Loop` iterates `range cm.needsRefresh`. The channel is never closed, so even after `closed` fires and `ConnectorSuper` creates a new `Runc` instance, the old goroutine continues to run until the process exits. Multiple reconnects leak goroutines.
- Safe modification: Use a select with `cm.closed` in the loop body, and close `cm.needsRefresh` on shutdown.
- Test coverage: None.

**Global package-level logger singletons called at `init` time:**
- Files: `cwidgets/compact/row.go:13`, `connector/collector/main.go:10`, `config/main.go:15`, `connector/main.go:14`
- Why fragile: `logging.Init()` is called at package-init time. If the singleton `Log` is nil at that point and multiple packages init concurrently, there is a data race on `logging.Log`. In practice Go serialises package init, but the pattern is brittle and prevents writing unit tests without a full logger initialisation.
- Safe modification: Pass the logger as a dependency rather than calling `logging.Init()` in `var` declarations.
- Test coverage: None.

## Scaling Limits

**`needsRefresh` channel capacity is 60 (both Docker and runc):**
- Current capacity: Buffer of 60 container IDs.
- Limit: On hosts with more than 60 containers, sending to `needsRefresh` blocks the caller (Docker event goroutine or runc refresh ticker). For Docker this causes `watchEvents` to stall and miss events.
- Scaling path: Use an unbuffered refresh with a coalescing mechanism (e.g., a `map[string]struct{}` + mutex + ticker drain) instead of a fixed-size channel.

## Dependencies at Risk

**`github.com/gizak/termui` (archived, v2 pre-release):**
- Risk: No upstream maintenance since 2018. No CVE path. No Go 1.26+ compatibility guarantee beyond what the current pseudo-version delivers.
- Impact: All rendering code breaks if Go deprecates any API termui depends on.
- Migration plan: `_docs/bubbletea-migration-plan.md` — replace with `charmbracelet/bubbletea` + `lipgloss`.

**`github.com/nsf/termbox-go` (archived, 2019 snapshot):**
- Risk: Repository is archived. No upstream maintenance. Used for `tm.Clear`, `tm.SetInputMode`, `tm.Sync`, raw terminal control.
- Impact: Removed or broken in a future Go version would require patching the vendored code.
- Migration plan: Bubble Tea migration removes this dependency entirely.

**`github.com/c9s/goprocinfo` (2017 snapshot, no recent activity):**
- Risk: Pinned to a 9-year-old commit. Parses `/proc/stat` and `/proc/meminfo` with a format that may diverge from current kernels. Any new field added by the kernel could cause parsing errors.
- Impact: Silent wrong CPU/memory values on modern kernels.
- Migration plan: Replace with direct syscall or `golang.org/x/sys/unix`.

## Missing Critical Features

**No integration tests or coverage for connectors:**
- Problem: Only 1 test file exists (`widgets/view_test.go`) covering 4 trivial string-split cases. Zero test coverage for connector logic, collector math, or cursor navigation.
- Blocks: Refactoring the data race in `Container.Meta`, the Bubble Tea migration, and runc I/O fix cannot be verified without a regression safety net.

**runc connector has no Logs support:**
- Problem: `connector/collector/runc.go:58-60` returns `nil` from `Logs()`. The log viewer in `menus.go:351` will receive a nil `LogCollector` and panic when calling `Stream()`.
- Files: `connector/collector/runc.go:58-60`, `menus.go:362`
- Blocks: Using the log viewer ([l] key) against a runc backend.

## Test Coverage Gaps

**Connector data-race paths:**
- What's not tested: Concurrent writes to `Container.Meta` from event goroutines and reads from the TUI goroutine.
- Files: `container/main.go`, `connector/docker.go:252-268`
- Risk: Silent data corruption; panic in production under the race detector.
- Priority: High

**runc-specific logic:**
- What's not tested: `GetLibc` caching, `refreshAll` error path that closes the channel, `Loop` goroutine lifecycle.
- Files: `connector/runc.go`
- Risk: Double-close panic, goroutine leak on reconnect.
- Priority: High

**Filter regex compilation:**
- What's not tested: Invalid regex input to the filter box.
- Files: `container/sort.go:115`
- Risk: Program exits on user typo.
- Priority: Medium

**Column widget rendering with edge-case column counts:**
- What's not tested: `CompactRow.Highlight()` with fewer than 2 enabled columns.
- Files: `cwidgets/compact/row.go:98`
- Risk: Index panic.
- Priority: Medium

---

*Concerns audit: 2026-06-11*
