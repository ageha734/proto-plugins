# Proto Plugins

Proto TOML plugin registry — 25 CLI tools for moonrepo/proto.

## Quick Commands

- `moon run :ci` — Full CI (lint + test)
- `moon run toml:test` — Plugin integration tests
- `moon generate plugin` — Generate new TOML plugin from template
- `deno test --allow-read .github/scripts/` — Script tests
- `cargo run --manifest-path tools/plugin-gen/Cargo.toml -- -o OWNER -r REPO` — Auto-detect plugin values

## Plugin Types

| Type | Use Case                                 | How                    |
|------|------------------------------------------|------------------------|
| TOML | Simple CLI (download binary from GitHub) | `moon generate plugin` |
| WASM | Complex install (gcloud, aws-cli)        | `/wasm-plugin` skill   |

## Key Directories

- `toml/` — Plugin definitions + Go tests
- `.github/scripts/` — Deno scripts (report, registry)
- `.github/workflows/` — GitHub Actions
- `tools/plugin-gen/` — Rust tool for GitHub API analysis
- `.moon/generators/plugin/` — Moon template
