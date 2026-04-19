# envoy-diff

A CLI tool to diff environment variable sets across staging and production configs.

## Installation

```bash
go install github.com/yourorg/envoy-diff@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/envoy-diff.git
cd envoy-diff && go build -o envoy-diff .
```

## Usage

```bash
envoy-diff --staging staging.env --production production.env
```

**Example output:**

```
~ API_TIMEOUT        30s  →  10s
+ NEW_RELIC_KEY      (missing in staging)
- DEBUG_MODE         true (missing in production)
```

### Flags

| Flag           | Description                        | Default  |
|----------------|------------------------------------|----------|
| `--staging`    | Path to staging env file           | required |
| `--production` | Path to production env file        | required |
| `--format`     | Output format: `text`, `json`      | `text`   |
| `--only-diff`  | Show only differing variables      | `false`  |

### Supported File Formats

- `.env` (dotenv)
- `.yaml` / `.yml`
- `.json`

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

## License

[MIT](LICENSE)