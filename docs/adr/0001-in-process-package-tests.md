# In-process package tests

Each Go package is tested in-process, through the code as it stands. A test may use `httptest`, a temporary directory, environment variables, and a context value the package already reads. Production code does not grow interfaces, clocks, or listeners for the sake of tests. Functions that only dial, listen, ping, or serve stay untested, including the ACME directory calls, `linode.NewClient`, `bgp.Start`, `database.New`, successful query paths in `database/mysql`, `database/postgres`, and `database/transaction`, `pprof.Start` when profiling is enabled, `http/server.Start`, and real OS signal delivery in `graceful`. Branches that depend on the machine stay untested too: stderr being a terminal, the default public-suffix file, and `cli/loader` searching `/etc` and the real user config directory. Explicit config paths, `ForceColor`, and a list file created in a temporary directory stay in scope. `graceful/v2` is out of this effort, including CI, because that module does not compile.

## Considered options

- **Seams in production code.** Inject the HTTP client, database, listener, and clock so those functions can run in-process. Rejected: the wrappers are thin, and the seams would exist only for tests.
- **Live integration suite.** Run MySQL, Postgres, and an ACME directory under test. Rejected: that is a separate effort, with different setup and failure modes.

## Consequences

CI compiles every module except `graceful/v2`. A coverage percentage would mostly measure the untested dial, listen, ping, and serve paths. Those paths are a gap by decision.
