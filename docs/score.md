# Score Filter

The `score` filter reorders diff result entries by a **relevance score** computed from how many user-supplied substrings appear in each key name. Keys that match more substrings bubble to the top of their category (Added, Removed, Changed, Same). Ties are broken lexicographically.

This is useful when you want to surface the most relevant changes first in large environment diffs.

## Usage

```
envoy-diff --score DB,CACHE staging.env production.env
```

## How scoring works

For each key, the score is the count of substrings (case-insensitive) found anywhere in the key name.

| Key            | Substrings     | Score |
|----------------|----------------|-------|
| `DB_HOST`      | `DB`, `CACHE`  | 1     |
| `DB_CACHE_KEY` | `DB`, `CACHE`  | 2     |
| `APP_NAME`     | `DB`, `CACHE`  | 0     |

Keys with score `0` are still included — they appear after higher-scored keys.

## API

```go
import "github.com/yourorg/envoy-diff/internal/filter"

out := filter.ApplyScore(result, []string{"DB", "CACHE"})
```

`ComputeScore(key string, substrings []string) int` is also exported for use in custom reporters or formatters that want to highlight matched keys.

## Notes

- Scoring does **not** remove any keys; it only reorders them.
- Because Go maps are unordered, ordered output requires passing the result through a formatter that calls `SortedKeys` or equivalent after scoring.
- Combine with `--only-changed` or prefix filters to narrow the result before scoring.
