# Annotate Filter

The `--annotate` flag lets you wrap or prefix environment variable values with a custom template string. This is useful for adding context, labels, or environment markers to specific keys in the diff output.

## Usage

```bash
envoy-diff --annotate DB_HOST="[prod] {{value}}" --annotate API_URL="endpoint:{{value}}"
```

## Rule Format

Each rule is a `key=template` pair:

- **key**: The exact environment variable name to annotate (case-insensitive match).
- **template**: A string that may include `{{value}}` as a placeholder for the original value.

## Examples

### Prefix a value

```
--annotate DB_HOST="[prod] {{value}}"
```

If `DB_HOST` has value `db.internal`, the output will show `[prod] db.internal`.

### Wrap a value

```
--annotate SECRET_KEY="***{{value}}***"
```

## Behaviour

- Rules are applied to **added**, **removed**, **changed** (both old and new), and **same** values.
- If a key matches multiple rules, all matching rules are applied in order.
- Keys not matching any rule are left unchanged.
- An empty rule list is a no-op.

## Parsing Errors

The CLI will return an error if:
- A rule is missing the `=` separator.
- The key portion is empty.
- The template portion is empty.
