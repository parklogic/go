# Lint gates new code only

The CI lint job runs the existing `.golangci.yml` with `--new-from-rev`, against the pull-request base and against the parent commit on `main`. Findings already in the tree do not fail the job, and this effort does not clear them. `testpackage` and `paralleltest` are disabled. `exhaustruct` does not run on `_test.go`.

`testpackage` forbids same-package tests, and `paralleltest` requires `t.Parallel` in every test. Both contradict the testing decisions. `exhaustruct` flags a table case that sets only the field under test. A full-tree run of this config reports 1,389 findings, so a job that failed on the whole tree would be red before any new test is added.

## Considered options

- **Clear every existing finding in this effort.** Rejected: that cleanup is separate from adding tests.
- **Keep `testpackage` and `paralleltest` enabled.** Rejected: new tests are same-package, and `t.Parallel` is only for tests that do not touch process state.

## Consequences

A pull request fails lint only for issues on lines it introduces. The older findings stay until a deliberate cleanup.
