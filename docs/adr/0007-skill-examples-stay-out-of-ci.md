# Skill examples stay out of CI

CI used to list modules by finding every `go.mod` except `graceful/v2`. That list included `.agents/skills/golang-cli/assets/examples`, whose module path is `github.com/you/myapp/cmd/myapp`. `go test` and `golangci-lint` then failed, because [ADR 0004](0004-workspace-builds-this-commit.md) does not put that directory in `go.work`.

The workspace and the CI module list cover library modules this repo ships. A skill example under `.agents` stays out of both. This narrows the "every module except `graceful/v2`" lines in ADR 0001 and ADR 0004: a skill example is not one of those modules. `graceful/v2` stays out for the reason in ADR 0004.

## Considered options

- **Add the example to `go.work`.** Rejected: CI would build and lint a tutorial module that is not a ParkLogic package.
- **Delete the example `go.mod`.** Rejected: the skill ships that fixture so the example stands alone.

## Consequences

The module finder skips `./.agents/*`. A new library module still belongs in `go.work` and is picked up by that finder. A skill example with its own `go.mod` is not a CI target.
