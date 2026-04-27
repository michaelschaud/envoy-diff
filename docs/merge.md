# Merge Filter

The `merge` filter injects an overlay map of key/value pairs into a diff result,
resolving conflicts according to a configurable strategy.

## Strategies

| Strategy | Behaviour |
|----------|-----------|
| `left`   | Existing values are always preserved; overlay is ignored for conflicting keys. |
| `right`  | Overlay values overwrite existing values for matching keys (default). |
| `union`  | Overlay values overwrite matching keys **and** new keys from the overlay are inserted. |

## CLI flags

```
--merge-overlay KEY=VALUE   Repeat to supply multiple pairs.
--merge-strategy left|right|union   Default: right.
```

## Example

```bash
envoy-diff staging.env production.env \
  --merge-overlay FEATURE_FLAG=true \
  --merge-overlay MAX_RETRIES=5 \
  --merge-strategy union
```

This will:
1. Load and diff the two env files.
2. Overwrite any existing `FEATURE_FLAG` or `MAX_RETRIES` entries with the
   supplied values.
3. Insert those keys into the result even if they were not present in either
   file (because `union` strategy is used).

## Notes

- The overlay is applied **after** the initial diff, so it does not affect
  which keys are classified as added, removed, or changed — it only updates
  their values in the output.
- Combine with `--only-changed` to focus on keys that differ between
  environments even after the overlay is applied.
