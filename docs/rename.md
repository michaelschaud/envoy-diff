# Key Renaming

The `--rename` flag lets you rewrite environment variable key names in the diff output before they are displayed or exported. This is useful when comparing configs that use different naming conventions across environments.

## Usage

```bash
envoy-diff --rename APP=SVC --rename DB=DATABASE staging.env production.env
```

## Rule Format

Each rule is specified as `FROM=TO`:

- `FROM` — the key or key prefix to match (case-insensitive)
- `TO` — the replacement name or prefix

## Matching Behaviour

| Input Key     | Rule       | Output Key        |
|---------------|------------|-------------------|
| `APP_HOST`    | `APP=SVC`  | `SVC_HOST`        |
| `APP_TIMEOUT` | `APP=SVC`  | `SVC_TIMEOUT`     |
| `APP`         | `APP=SVC`  | `SVC`             |
| `APPNAME`     | `APP=SVC`  | `APPNAME` (no match — requires `_` boundary or exact) |

Matching uses a **prefix boundary** rule: a rule `FROM=TO` matches a key if the key equals `FROM` exactly, or if the key starts with `FROM_`. This prevents unintended partial-word matches.

## Multiple Rules

Multiple `--rename` flags may be provided. Rules are evaluated **in order** and the **first match wins**.

```bash
envoy-diff --rename APP=SERVICE --rename APP=OTHER staging.env production.env
# APP_HOST → SERVICE_HOST  (first rule wins)
```

## Notes

- Renaming is applied **after** filtering (prefix, glob, regex, value, key) and **before** output formatting.
- Renaming does not affect which keys are included or excluded — only how they are labelled in output.
- Renamed keys are reflected in all output formats: text, JSON, CSV, and Markdown.
