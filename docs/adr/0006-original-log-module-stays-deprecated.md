# The original log module stays deprecated

`log` is deprecated and stays as it is. That includes process-wide zerolog setup, the slog bridge (`GetSlogLevel`, `GetSlogger`), and the caller-path fields `ShortenPath` and `RemovePrefix`. This line of work does not change that module.

## Considered options

- **Deepen `log`.** One caller-path rule, and move the slog bridge into `bgp`, its only caller in this repo. Rejected: a deprecated module gets no further changes.
- **Fold `log` into `log/v2`.** Rejected: that is a change to `log`, and `log/v2` has no caller here except `graceful/v2`, which is out of scope.

## Consequences

Future architecture reviews must not propose changes to `log`. Existing callers stay on it. New logging work goes in `log/v2`.
