# Result Limiting

The `--limit` flags allow you to cap the number of entries shown per diff category. This is useful when working with large environment files where you only need a quick overview.

## Flags

| Flag              | Description                                      | Default |
|-------------------|--------------------------------------------------|---------|
| `--limit-added`   | Maximum number of added keys to display          | 0 (all) |
| `--limit-removed` | Maximum number of removed keys to display        | 0 (all) |
| `--limit-changed` | Maximum number of changed keys to display        | 0 (all) |
| `--limit-same`    | Maximum number of identical keys to display      | 0 (all) |

A value of `0` means no limit is applied for that category.

## Example

```bash
envoy-diff staging.env production.env --limit-added 5 --limit-changed 10
```

This will show at most 5 added keys and 10 changed keys, while showing all removed and identical keys.

## Truncation Notice

When a limit is applied and entries are truncated, the output will include a notice indicating how many entries were omitted. For example:

```
... and 42 more added keys (use --limit-added 0 to show all)
```

This helps you understand when you are seeing a partial view of the diff.

## Notes

- Limits are applied **after** all other filters (prefix, glob, regex, value, key).
- The selection of which entries are shown within a limit is not guaranteed to be sorted; use output formats like `--format=csv` if stable ordering is required.
- Combining `--limit-same 0` with `--only-changed` is equivalent — neither will show identical keys.
- Limit values must be non-negative integers; passing a negative value will result in an error.
