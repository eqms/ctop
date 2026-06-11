# Testing Patterns

**Analysis Date:** 2026-06-11

## Test Framework

**Runner:**
- Go standard `testing` package — no third-party test framework
- No test configuration file (no `testify`, no `gomock`, no `gocheck`)

**Assertion Library:**
- None — raw `t.Errorf(...)` used for all assertions

**Run Commands:**
```bash
go test ./...        # Run all tests
go test -v ./...     # Verbose output
go test -run TestName ./...  # Run specific test
```

CI runs `go test ./...` on `ubuntu-latest` with Go 1.26 (see `.github/workflows/ci.yml`).

## Test File Organization

**Location:**
- Co-located with the package under test: `widgets/view_test.go` tests `widgets/view.go`

**Naming:**
- Files: `<subject>_test.go`
- Test functions: `Test<Subject><Scenario>` — e.g., `TestSplitEmptyLine`, `TestSplitLineLongerThanLimit`

**Current test coverage:**
```
widgets/view_test.go   # Only test file in the repo
```

## Test Structure

**Suite Organization:**
```go
package widgets   // same package as code under test (whitebox)

import "testing"

func TestSplitEmptyLine(t *testing.T) {
    result := splitLine("", 5)
    if len(result) != 0 {
        t.Errorf("expected: 0 lines, got: %d", len(result))
    }
}
```

**Patterns:**
- No setup or teardown (`TestMain`, `t.Cleanup`, etc.) — tests are self-contained
- No table-driven tests in the existing file (each scenario is its own function)
- Assertions use `t.Errorf` (non-fatal, continues test) rather than `t.Fatalf`
- No subtests (`t.Run(...)`) in current code

**Example — full test file (`widgets/view_test.go`):**
```go
package widgets

import "testing"

func TestSplitEmptyLine(t *testing.T) {
    result := splitLine("", 5)
    if len(result) != 0 {
        t.Errorf("expected: 0 lines, got: %d", len(result))
    }
}

func TestSplitLineShorterThanLimit(t *testing.T) {
    result := splitLine("hello", 7)
    if len(result) != 1 {
        t.Errorf("expected: 0 lines, got: %d", len(result))
    }
}

func TestSplitLineLongerThanLimit(t *testing.T) {
    result := splitLine("hello", 3)
    if len(result) != 2 {
        t.Errorf("expected: 0 lines, got: %d", len(result))
    }
}

func TestSplitLineSameAsLimit(t *testing.T) {
    result := splitLine("hello", 5)
    if len(result) != 1 {
        t.Errorf("expected: 0 lines, got: %d", len(result))
    }
}
```

## Mocking

**Framework:** None — the codebase uses a `mock` build tag (`//go:build !release`) to swap in fake implementations at compile time rather than using a mocking library.

**Mock backends:**
- `connector/mock.go` — synthetic `MockConnector` used with `CTOP_CONNECTOR=mock` or in `run-dev` mode
- `connector/collector/mock.go` — `MockCollector` that generates fake metrics
- `connector/manager/mock.go` — `MockManager` stub for all manager operations
- `connector/collector/mock_logs.go` — fake log stream

**Build-tag mocking pattern:**
```go
//go:build !release
// +build !release

package connector

func init() { enabled["mock"] = NewMock }
```

The mock connector is excluded from release binaries (built with `-tags release`) but available in development and tests.

**What to mock:**
- Connector backends — swap in `mock` connector for integration-style testing without a running Docker daemon

**What NOT to mock:**
- Pure functions (sorting, string splitting, config parsing) — test directly

## Fixtures and Factories

**Test Data:**
- No fixtures directory, no factory helpers
- Test inputs are inline literals: `splitLine("hello", 5)`

**Models:**
- `models.NewMetrics()` and `models.NewMeta(kvs ...string)` serve as factories for test data in production code — they can be used directly in tests when needed

## Coverage

**Requirements:** None enforced — no coverage threshold in CI or Makefile

**View Coverage:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## Test Types

**Unit Tests:**
- Scope: pure functions only (currently just `widgets.splitLine`)
- Location: `widgets/view_test.go`
- Approach: whitebox, same package, raw `testing` package

**Integration Tests:**
- Not present in test suite
- The mock connector (`connector/mock.go`, build tag `!release`) enables manual integration-style testing via `make run-dev`, but this is not automated

**E2E Tests:**
- Not implemented — TUI applications are not tested end-to-end

## Vulnerability Scanning

CI runs `govulncheck ./...` with `continue-on-error: true` — failures produce a `::warning::` annotation but do not block the build. See `.github/workflows/ci.yml`.

```yaml
- name: Run govulncheck
  continue-on-error: true
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./... || echo "::warning::govulncheck found vulnerabilities in transitive dependencies"
```

## Adding New Tests

**Where to place:** Co-locate with the package — `<package>/<file>_test.go`.

**Package declaration:** Use the same package name as the file under test (whitebox testing, no `_test` suffix) unless you need blackbox testing.

**Pattern to follow:**
```go
package <packagename>

import "testing"

func Test<Subject><Scenario>(t *testing.T) {
    result := functionUnderTest(input)
    if result != expected {
        t.Errorf("expected: %v, got: %v", expected, result)
    }
}
```

**Linux-only code:** Use build tag to restrict tests:
```go
//go:build linux
// +build linux
```
Always verify with `GOOS=linux go test ./...` after editing `connector/runc.go` or `connector/collector/runc.go`.

---

*Testing analysis: 2026-06-11*
