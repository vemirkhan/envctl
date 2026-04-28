# envctl

Manage and sync environment variable sets across multiple projects and deployment targets from a single config.

---

## Installation

```bash
go install github.com/yourname/envctl@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourname/envctl/releases).

---

## Usage

Define your environments in a single `envctl.yaml` config file:

```yaml
projects:
  my-api:
    targets:
      staging:
        DATABASE_URL: postgres://staging-host/mydb
        LOG_LEVEL: debug
      production:
        DATABASE_URL: postgres://prod-host/mydb
        LOG_LEVEL: info
```

Then sync variables to your desired target:

```bash
# List all projects and targets
envctl list

# Sync environment variables to a target
envctl sync my-api --target staging

# Export variables to a .env file
envctl export my-api --target production --out .env

# Diff two targets
envctl diff my-api --from staging --to production
```

---

## Why envctl?

- **Single source of truth** for all environment configs across projects
- **Target-aware** — manage staging, production, and any custom environments
- **Portable** — works with `.env` files, CI pipelines, and shell exports
- **Simple config** — plain YAML, easy to version control

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)