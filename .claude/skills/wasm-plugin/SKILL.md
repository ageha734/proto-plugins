---
name: wasm-plugin
description: Create a WASM proto plugin for tools that cannot be implemented as TOML plugins (gcloud, aws-cli, nvm, etc). Use when the user mentions WASM plugin, complex installation, or tools that need custom install logic.
allowed-tools: Bash(cargo *) Bash(proto *) Read Write Edit
agent: wasm-plugin-dev
context: fork
---

# Create WASM Proto Plugin

## Trigger Conditions

Use this skill when the target tool requires any of:

- Running an installer script
- Environment variable configuration during installation
- Non-standard version resolution (npm registry, custom API)
- Post-install processing steps
- Complex archive structures with conditional logic
- Shell profile modification

## Workflow

1. **Research the tool's installation process**
   - How is it typically installed?
   - What platforms/architectures are supported?
   - How are versions resolved?
   - What environment variables are needed?

2. **Scaffold the plugin**

   ```bash
   mkdir -p wasm/<name>/src
   ```

3. **Implement required hooks**
   - `register_tool` — Always required
   - `resolve_version` — If custom version resolution needed
   - `download_prebuilt` — For downloading binaries
   - `native_install` — For running install scripts
   - `locate_executables` — For finding binaries post-install

4. **Build and test**

   ```bash
   cd wasm/<name>
   cargo build --target wasm32-wasip1 --release
   proto plugin add <name> source:./target/wasm32-wasip1/release/<name>.wasm
   proto install <name> latest
   <name> --version
   ```

5. **Integration**
   - Add to `.prototools`
   - Create integration test

## Examples of WASM-only tools

| Tool | Reason |
| ------ | -------- |
| gcloud | Installer script + env vars |
| aws-cli | Bundled installer execution |
| nvm/pyenv/rbenv | Shell profile modification + version management |
| volta | Custom archive + env vars |
