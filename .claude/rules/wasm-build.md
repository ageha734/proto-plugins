---
paths:
  - "wasm/**/*.rs"
  - "wasm/**/Cargo.toml"
---

# WASM Plugin Build Rules

- Compile with: `cargo build --target wasm32-wasip1 --release`
- Use `proto_pdk` crate for hook implementations
- Use `extism_pdk` for host function bindings
- Do not use packages that require filesystem or unrestricted network access
- Test with: `proto plugin add <name> source:./target/wasm32-wasip1/release/<name>.wasm`
- Every hook function must be annotated with `#[plugin_fn]`
- Return types must be wrapped in `FnResult<Json<...>>`
