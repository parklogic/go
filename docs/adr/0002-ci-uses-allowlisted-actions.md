# CI uses allowlisted actions

The test workflow installs `golangci-lint` and `govulncheck` itself and uses only `actions/checkout` and `actions/setup-go`, each pinned to a full commit SHA. The ParkLogic GitHub org runs an action only when it is `parklogic/*`, GitHub-owned, or an approved `docker/*` action, and only when the reference is a full SHA. A major tag such as `@v5` fails at startup with no jobs.

## Considered options

- **`golangci/golangci-lint-action` and `codecov/codecov-action`.** Rejected: neither is on the allowlist.
- **Pin GitHub-owned actions to a floating major tag.** Rejected: the org requires a full commit SHA.

## Consequences

Bumping the linter or `govulncheck` is a change to the install step. Coverage is not uploaded.
