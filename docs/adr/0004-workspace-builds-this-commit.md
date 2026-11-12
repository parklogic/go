# Workspace builds this commit

CI and local tests use a committed `go.work`. `go.work` and `go.work.sum` are not ignored. `go work sync` writes the workspace build list back into each module's `go.mod`, including that module's `go` line when the build list requires a newer Go. `api/router.go` imports `github.com/go-chi/chi/v5` so `api` compiles against the `chi` in this tree.

Without the workspace, `go test` in `api` builds the published `github.com/parklogic/go/chi` pseudo-version required by `api/go.mod`. That published copy does not compile. The `chi` directory in this commit does.

## Considered options

- **`replace` directives in `api/go.mod`.** Rejected: a workspace covers every sibling module with one file, and downstream users of a published module ignore a `go.work` they never see.
- **Test `api` against the published `chi`.** Rejected: that module does not compile, so `api` could not have a test.

## Consequences

A test in this repo builds sibling modules from the same commit. Publishing one module still uses the versions in its `go.mod`. Which modules belong in the workspace is decided in [ADR 0007](0007-skill-examples-stay-out-of-ci.md). `graceful/v2` stays out until it compiles.

## History

This ADR first said `go.work` lists every module except `graceful/v2`. [ADR 0007](0007-skill-examples-stay-out-of-ci.md) narrowed that: a skill example under `.agents` is not one of those modules.
