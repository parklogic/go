# Workspace builds this commit

CI and local tests use a committed `go.work` that lists every module except `graceful/v2`. `go.work` and `go.work.sum` are not ignored. `api/router.go` imports `github.com/go-chi/chi/v5` so `api` compiles against the `chi` in this tree.

Without the workspace, `go test` in `api` builds the published `github.com/parklogic/go/chi` pseudo-version required by `api/go.mod`. That published copy does not compile. The `chi` directory in this commit does.

## Considered options

- **`replace` directives in `api/go.mod`.** Rejected: a workspace covers every sibling module with one file, and downstream users of a published module ignore a `go.work` they never see.
- **Test `api` against the published `chi`.** Rejected: that module does not compile, so `api` could not have a test.

## Consequences

A test in this repo builds sibling modules from the same commit. Publishing one module still uses the versions in its `go.mod`. `graceful/v2` stays out of the workspace until it compiles.
