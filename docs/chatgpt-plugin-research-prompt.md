# ChatGPT プラグイン調査プロンプト

GitHub API でリリースアセット情報が取得できない場合に、このプロンプトを ChatGPT に渡して調査してもらう。

## 用途

Non-WASM (TOML) proto plugin を作成するために必要な情報を調査する。

## プロンプト

```text
以下のCLIツールについて、proto TOML pluginを作成するために必要な情報を調査してください。

ツール名: {TOOL_NAME}
リポジトリ: https://github.com/{OWNER}/{REPO}

以下の情報を正確に調べてください:

1. **バイナリ名**: インストール後に実行するコマンド名
2. **バージョン確認コマンド**: `--version` or `version` のどちらか
3. **リリースアセットのパターン**: GitHubのReleasesページから、各プラットフォーム向けのダウンロードファイル名パターンを特定
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)
4. **アーカイブ形式**: tar.gz / zip / 単一バイナリ のどれか
5. **アーカイブ内のプレフィックス**: 展開後にディレクトリが含まれる場合のパス
6. **チェックサム**: checksumファイルが提供されているか、ファイル名パターン
7. **アーキテクチャ表記**: aarch64/arm64/ARM64 のどれが使われているか、x86_64/amd64/x64 のどれか

回答は以下のJSON形式で提供してください:

{
  "name": "ツール名(kebab-case)",
  "github_owner": "オーナー",
  "github_repo": "リポジトリ名",
  "binary_name": "バイナリ名",
  "version_command": "--version or version",
  "linux_download_file": "パターン ({version}, {arch} を使用)",
  "macos_download_file": "パターン",
  "windows_download_file": "パターン",
  "has_checksum": true/false,
  "checksum_file": "チェックサムファイル名パターン",
  "arch_aarch64": "arm64 or aarch64",
  "arch_x86_64": "amd64 or x86_64",
  "needs_unpack": true/false,
  "archive_prefix_linux": "展開時のプレフィックス (あれば)",
  "archive_prefix_macos": "展開時のプレフィックス (あれば)"
}

重要な注意点:
- {version} はバージョン番号のプレースホルダー（vプレフィックスなし）
- {arch} はアーキテクチャのプレースホルダー
- 実際のリリースページを確認して正確なファイル名を調べてください
- checksum_fileには{version}プレースホルダーを使ってください
```

## 取得後の使い方

1. ChatGPT からJSON回答を取得
2. JSON の値を使って `moon generate plugin` を実行:

    ```bash
    moon generate plugin -- \
      --name "{name}" \
      --github_owner "{github_owner}" \
      --github_repo "{github_repo}" \
      --binary_name "{binary_name}" \
      --version_command "{version_command}" \
      --linux_download_file "{linux_download_file}" \
      --macos_download_file "{macos_download_file}" \
      --windows_download_file "{windows_download_file}" \
      --has_checksum {has_checksum} \
      --checksum_file "{checksum_file}" \
      --arch_aarch64 "{arch_aarch64}" \
      --arch_x86_64 "{arch_x86_64}" \
      --needs_unpack {needs_unpack}
    ```

3. 生成された `toml/{name}.toml` と `toml/{name}_test.go` を確認
4. `.prototools` にバージョンエントリを追加
5. テスト実行: `moon run toml:test`

## 判定基準: TOML vs WASM

以下に該当する場合は **WASM plugin が必要** (このプロンプトでは対応不可):

- インストーラースクリプトの実行が必要（gcloud, aws-cli）
- インストール後に環境変数の設定が必要
- バージョン解決にカスタムロジックが必要（npm レジストリ等）
- シェルプロファイルの変更が必要
- ダウンロード後にビルドが必要
- 複雑なアーカイブ構造（条件分岐が必要）

WASM plugin が必要な場合は Claude Code の `/wasm-plugin` スキルを使用してください。
