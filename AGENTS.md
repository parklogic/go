# Agent instructions

## Agent skills

### Issue tracker

GitHub Issues in `parklogic/go`, accessed via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

Default labels: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: one `GLOSSARY.md` and `docs/adr/` at the repo root. See `docs/agents/domain.md`. When an ADR names a successor, obey the successor.

### Checks

Lint with `scripts/lint-new` before pushing. Run `govulncheck` with `GOTOOLCHAIN=go<version>`, where `<version>` is the `go-version` in `.github/workflows/ci.yml`.
