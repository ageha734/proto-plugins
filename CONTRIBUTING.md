# Contributing

## Prerequisites

- [proto](https://moonrepo.dev/proto) >= 0.57.3
- [moon](https://moonrepo.dev/moon) >= 2.3.0

Both are installed automatically via `.prototools` when you run `proto use`.

## Development Setup

```bash
git clone https://github.com/ageha734/proto-plugins.git
cd proto-plugins
proto use
```

## Adding a New TOML Plugin

### Option 1: Moon Generator (Recommended)

```bash
moon generate plugin
```

Follow the interactive prompts to provide plugin metadata.

### Option 2: Rust Tool (Auto-detect from GitHub)

```bash
cargo run --manifest-path tools/plugin-gen/Cargo.toml -- --owner OWNER --repo REPO
```

This analyzes the GitHub release assets and outputs `moon generate` arguments.

### Option 3: Manual

1. Create `toml/<name>.toml` following existing plugins as reference
2. Create `toml/<name>_test.go` with the standard test pattern
3. Add version to `.prototools`
4. Add plugin entry to `[plugins]` in `.prototools`

## Testing

```bash
# Run all tests
moon run :ci

# Run plugin integration tests only
moon run toml:test

# Run a specific plugin test
cd toml && go test -run TestPinact -v

# Run Deno script tests
deno test --allow-read .github/scripts/
```

## Code Style

- Go formatting: `dprint fmt` (via exec plugin with gofumpt)
- Go linting: `golangci-lint run`
- Config files: `dprint check` (JSON, TOML, YAML, Markdown)
- Deno scripts: `deno fmt` / `deno lint`

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` — New plugin or feature
- `fix:` — Bug fix
- `chore:` — Maintenance (version bumps, CI changes)
- `docs:` — Documentation updates

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Add your plugin with tests
4. Ensure `moon run :ci` passes
5. Submit a PR targeting `master`

## Plugin Guidelines

- All plugins must include integration tests
- Checksum verification is strongly recommended
- Use `[resolve] git-url` for version resolution
- Pin exact versions in `.prototools`
- Platform support: Linux (x86_64, aarch64), macOS (aarch64), Windows (x86_64) at minimum
