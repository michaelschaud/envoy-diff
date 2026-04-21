# Environment Variable Override Filter

The `ApplyEnvOverride` filter allows values loaded from `.env` files to be
overridden by variables set in the **current process environment** before the
diff is computed or displayed.

## Use Case

This is useful when you want to inject secrets or local overrides without
modifying the source `.env` files — for example in CI pipelines.

## Behaviour

- For every key in `Added`, `Removed`, `Same`, and `Changed` sections, if a
  matching environment variable exists in the process, its value replaces the
  loaded value.
- For `Changed` entries, only the **new** value is overridden; the **old**
  value is preserved as-is from the file.
- Keys not present in the process environment are left unchanged.

## Example

```bash
# Override a production secret before diffing
DB_PASSWORD=supersecret envoy-diff staging.env production.env
```

## `EnvVarNames` Helper

The `EnvVarNames(prefix string)` function returns all environment variable
names currently set in the process, optionally filtered by a case-insensitive
prefix. This is used internally and can assist with shell completion or
debugging.

```go
names := filter.EnvVarNames("AWS_")
// returns all env var names starting with AWS_ (case-insensitive)
```

## Related Filters

- [`docs/filters.md`](./filters.md) — overview of all available filters
- `ApplyKeyFilter` — filter result by key substring
- `ApplyValueFilter` — filter result by value substring
