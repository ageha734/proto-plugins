# Project Conventions

## Architecture

- `toml/` — Non-WASM (TOML) proto plugin definitions + Go integration tests
- `wasm/` — WASM proto plugins (Rust, compiled to wasm32-wasip1)
- `.moon/` — Moon workspace config, task definitions, generators
- `.github/scripts/` — Deno TypeScript scripts (CI report, registry)
- `.github/workflows/` — GitHub Actions (reusable, push, PR, cron, release)
- `tools/plugin-gen/` — Rust tool for GitHub API release analysis
- `docs/` — Documentation

## Rules

- Directory structure MUST NOT change
- All GitHub Actions pinned by full commit SHA
- Custom CI scripts in Deno (TypeScript)
- TDD: write tests before implementation
- Prefer `uses:` actions in workflows; custom scripts only when no action exists
- All tool versions in `.prototools` must be exact (no `>=` ranges)
- Go formatting via dprint (exec plugin with gofumpt)
- Go linting via dprint (dprint-plugin-golangci, lint only, no formatters)

## Plugin Development

### Non-WASM (TOML) — Simple CLI tools

```bash
moon generate plugin
# or
cargo run --manifest-path tools/plugin-gen/Cargo.toml -- --owner OWNER --repo REPO
```

### WASM — Complex tools (gcloud, aws-cli, etc.)

Use `/wasm-plugin` skill.

## Commands

- `moon ci` — Full CI for all projects
- `moon run toml:test` — Plugin integration tests
- `moon run workspace:lint` — dprint check
- `moon generate plugin` — Generate new TOML plugin
- `deno test --allow-read .github/scripts/` — Script tests
