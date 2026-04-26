# Diff Mode Filter

The `--mode` flag lets you restrict the output to specific diff categories,
so you can focus on only what matters for a given review.

## Supported modes

| Mode      | Description                                      |
|-----------|--------------------------------------------------|
| `added`   | Keys present in production but not in staging    |
| `removed` | Keys present in staging but not in production    |
| `changed` | Keys present in both but with different values   |
| `same`    | Keys present in both with identical values       |

## Usage

```bash
# Show only added and removed keys
envoy-diff --mode added,removed staging.env production.env

# Show only changed keys
envoy-diff --mode changed staging.env production.env

# Show everything (default — same as omitting --mode)
envoy-diff staging.env production.env
```

## Combining with other filters

`--mode` composes naturally with other filters such as `--prefix`, `--only-changed`,
`--exclude-glob`, and `--mask`. It is applied after the diff is computed and before
formatting, so all reporters (text, JSON, CSV, Markdown) respect it.

```bash
# Show only changed DATABASE_* keys, masking passwords
envoy-diff --mode changed --prefix DATABASE_ --mask PASSWORD staging.env production.env
```

## Notes

- Modes are **case-insensitive** and **comma-separated**.
- Unknown mode names are silently ignored.
- Specifying no `--mode` flag returns all categories (added, removed, changed, same).
