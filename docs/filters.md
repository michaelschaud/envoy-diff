# envoy-diff Filters

envoy-diff supports several filtering mechanisms to narrow down diff output.

## Prefix Filter (`--prefix`)

Keep only keys that start with a given prefix (case-insensitive).

```bash
envoy-diff --prefix APP_ staging.env production.env
```

Multiple prefixes can be supplied:

```bash
envoy-diff --prefix APP_ --prefix DB_ staging.env production.env
```

## Only Changed (`--only-changed`)

Suppress keys that are identical between environments.

```bash
envoy-diff --only-changed staging.env production.env
```

## Glob Exclude (`--exclude-glob`)

Exclude keys matching a shell-style glob pattern (case-insensitive).

```bash
envoy-diff --exclude-glob '*_SECRET' --exclude-glob '*_TOKEN' staging.env production.env
```

## Regex Exclude (`--exclude-regex`)

Exclude keys matching a regular expression.

```bash
envoy-diff --exclude-regex '^LEGACY_' staging.env production.env
```

## Value Filter (`--exclude-value`)

Exclude any key whose value (in either environment) contains the given substring.
Matching is case-insensitive. Useful for hiding placeholder or redacted values.

```bash
envoy-diff --exclude-value REDACTED --exclude-value TODO staging.env production.env
```

This applies to all categories: added, removed, changed, and same keys.
For **changed** keys, the key is excluded if *either* the old or the new value matches.

## Combining Filters

All filters are applied in sequence:

1. Prefix filter (keep matching keys)
2. Glob exclude (remove matching keys)
3. Regex exclude (remove matching keys)
4. Value exclude (remove matching keys)
5. Only-changed (remove identical keys)
