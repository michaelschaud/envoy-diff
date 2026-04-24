# Value Masking (`--mask`)

The `--mask` flag redacts sensitive environment variable values before they
are written to any output format. This is useful when sharing diff reports
with teammates who should not see raw secret values.

## Usage

```bash
envoy-diff --mask password --mask secret --mask token \
  staging.env production.env
```

Multiple `--mask` flags may be supplied. Each value is treated as a
case-insensitive **substring** that is matched against the **key name**.
When a key matches, its value is replaced with `***REDACTED***` in every
output category (`added`, `removed`, `same`, `changed`).

## Example

Given the following diff:

| Key          | Staging         | Production      |
|--------------|-----------------|-----------------|
| DB_PASSWORD  | `hunter2`       | `c0rrectH0rse`  |
| APP_PORT     | `8080`          | `443`           |

Running with `--mask password` produces:

| Key          | Staging           | Production        |
|--------------|-------------------|-------------------|
| DB_PASSWORD  | `***REDACTED***`  | `***REDACTED***`  |
| APP_PORT     | `8080`            | `443`             |

## Notes

- Matching is **case-insensitive**: `--mask PASSWORD` and `--mask password`
  are equivalent.
- Masking is applied **after** all other filters (prefix, glob, regex, etc.)
  so the full key set is visible but sensitive values are hidden.
- The original input files are never modified.
