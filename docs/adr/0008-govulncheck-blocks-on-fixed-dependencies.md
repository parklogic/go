# Govulncheck blocks on fixed dependencies

The called third-party findings are bumped to the versions that clear them. `golang.org/x/net` v0.60.0 declares `go 1.26.0` and requires `golang.org/x/text` v0.42.0, so `go work sync` raises `go.work` and every module that picks up that build list from `go 1.25.4` to `go 1.26.0`. CI still installs Go 1.26.11. Module code stays within the Go 1.25 language.

`govulncheck` fails the workflow. `github.com/osrg/gobgp/v4` v4.11.0 still calls GO-2026-4736, a NEXT_HOP denial of service with no fixed version in the Go vuln database, so only the `bgp` job stays `continue-on-error`.

The bumps are `github.com/go-chi/chi/v5` v5.3.0, `github.com/jackc/pgx/v5` v5.9.2, `golang.org/x/text` v0.42.0, `golang.org/x/net` v0.60.0, `google.golang.org/grpc` v1.83.2, `github.com/osrg/gobgp/v4` v4.11.0, and `golang.org/x/sys` v0.48.0. `grpc` is v1.83.2 because that is what `gobgp` v4.11.0 requires, above the v1.83.1 fix from the scan. `x/text` is v0.42.0 because v0.60.0 of `x/net` rejects v0.41.0.

## Considered options

- **Keep `go 1.25.4` and an older `x/net`.** Rejected: the `bgp` findings in `x/net` are fixed in v0.60.0, and tidy will not leave the `go` line at 1.25.4 once that module is required.
- **Fail `bgp` while GO-2026-4736 has no fixed version.** Rejected: the vuln database marks every release affected, and `server.fsmHandler.recvMessageloop` is still present in v4.11.0.

## Consequences

This supersedes the `go 1.25.4` line and the non-bump, non-blocking `govulncheck` decision in [ADR 0005](0005-ci-uses-go-1-26-11.md). A module whose selected dependencies still declare `go 1.25` keeps that line. A later `gobgp` release that the vuln database marks fixed can make the `bgp` job blocking.
