# 027: Build zapp package

## Objective
Implement the `zapp` package — an application lifecycle toolkit for zarlcorp tools. Provides resource tracking with ordered cleanup, signal-based context cancellation, and functional options. Toolkit, not framework — consumer owns `main`.

## Context
zapp is the top of the package layering in zarlcorp/core. Every tool will import it for lifecycle management. The design deliberately avoids framework patterns (no `App.Run(callback)`, no hidden goroutines). Consumer wires things together explicitly.

The `go.mod` already exists at `pkg/zapp/go.mod` with module path `github.com/zarlcorp/core/pkg/zapp`. The module is already listed in `go.work`. The zoptions dependency is at `github.com/zarlcorp/core/pkg/zoptions` and provides `Option[T any] func(*T)`.

## Requirements

### App struct
```go
type App struct {
    name    string
    closers []io.Closer
}
```

`New(opts ...zoptions.Option[App]) *App` — creates an App with defaults. Default name is the binary name from `os.Args[0]` (base only, no path).

### Track / Close
- `Track(c io.Closer)` — registers a closer for cleanup. Thread-safe (mutex-protected).
- `Close() error` — closes all tracked resources in LIFO order (last tracked, first closed). Sequential, not concurrent. Collects all errors via `errors.Join`. Safe to call multiple times — subsequent calls return the same error without re-closing.

### SignalContext
```go
func SignalContext(parent context.Context) (context.Context, context.CancelFunc)
```
Returns a context that is cancelled when SIGINT or SIGTERM is received, or when the returned cancel func is called. Uses `signal.NotifyContext` under the hood.

### CloserFunc
```go
type CloserFunc func() error
func (f CloserFunc) Close() error { return f() }
```
Adapter that wraps a `func() error` as `io.Closer`. Allows tracking cleanup functions that aren't already io.Closer.

### Functional options
- `WithName(name string) zoptions.Option[App]` — overrides the default app name.

Options live in `options.go`.

### What zapp does NOT do
- No logging setup — consumer configures slog
- No `Run` method or callback pattern — consumer owns the control flow
- No concurrent service management (no errgroup)
- No timeout on close — consumer can wrap with context if needed

### Package doc
Brief package comment on `zapp.go` explaining it's a lifecycle toolkit. Follow existing zoptions style — short, with a usage example showing `New`, `Track`, `SignalContext`, and `Close`.

### Example usage (for package doc)
```go
func main() {
    app := zapp.New(zapp.WithName("myservice"))

    ctx, cancel := zapp.SignalContext(context.Background())
    defer cancel()

    db := openDB()
    app.Track(db)

    srv := startServer(ctx, db)
    app.Track(zapp.CloserFunc(func() error {
        return srv.Shutdown(context.Background())
    }))

    <-ctx.Done()

    if err := app.Close(); err != nil {
        slog.Error("shutdown", "err", err)
        os.Exit(1)
    }
}
```

## Acceptance Criteria
1. `New` returns an `*App` with name defaulting to binary basename
2. `WithName` overrides the default name
3. `Track` registers closers, is safe for concurrent use
4. `Close` tears down in LIFO order, returns joined errors, is idempotent
5. `SignalContext` returns a context cancelled on SIGINT/SIGTERM
6. `CloserFunc` satisfies `io.Closer`
7. All tests pass (`go test ./...` from `pkg/zapp/`)
8. `go.mod` has the zoptions dependency with a local replace directive for workspace compatibility

## Target Repo
zarlcorp/core

## Agent Role
backend

## Files to Create/Modify
- `pkg/zapp/zapp.go` — App struct, New, Track, Close, SignalContext, CloserFunc, package doc
- `pkg/zapp/options.go` — WithName option
- `pkg/zapp/zapp_test.go` — table-driven tests for all public API
- `pkg/zapp/go.mod` — add zoptions dependency

## Dependencies
None — this is the first zarlcorp/core package build.

## Notes
- zoptions module path: `github.com/zarlcorp/core/pkg/zoptions`
- Use workspace replace directive in go.mod: `replace github.com/zarlcorp/core/pkg/zoptions => ../zoptions`
- The `name` field is unexported but accessible via option. No getter needed unless a future spec requires it.
- Keep Close idempotent using `sync.Once`.
- Tests should NOT send real signals — test LIFO ordering and error collection directly. SignalContext can be tested by calling the cancel func.
