# envdiff

Compare `.env` files across environments and surface missing or conflicting variables with structured output.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git && cd envdiff && go build ./...
```

---

## Usage

```bash
envdiff [flags] <base> <target> [target...]
```

### Example

```bash
envdiff .env.development .env.staging .env.production
```

**Output:**

```
MISSING in .env.staging:
  - DATABASE_URL
  - REDIS_URL

CONFLICT in .env.production:
  - API_BASE_URL  (dev: "http://localhost" | prod: "https://api.example.com")

OK: 12 variables match across all environments
```

### Flags

| Flag | Description |
|------|-------------|
| `--format` | Output format: `text` (default), `json`, `yaml` |
| `--ignore` | Comma-separated list of keys to ignore |
| `--strict` | Exit with non-zero status if any diff is found |

---

## Why envdiff?

Manually comparing `.env` files across multiple environments is error-prone. `envdiff` automates the process and integrates cleanly into CI pipelines to catch configuration drift before it reaches production.

---

## License

MIT © [yourusername](https://github.com/yourusername)