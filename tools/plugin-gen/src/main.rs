use clap::Parser;
use reqwest::blocking::Client;
use reqwest::header::{ACCEPT, USER_AGENT};
use serde::Deserialize;
use std::collections::HashMap;

#[derive(Parser)]
#[command(name = "plugin-gen")]
#[command(about = "Generate proto TOML plugin values from GitHub releases")]
struct Cli {
    /// GitHub owner/organization (e.g. "hashicorp")
    #[arg(short, long)]
    owner: String,

    /// GitHub repository name (e.g. "terraform")
    #[arg(short, long)]
    repo: String,

    /// GitHub token for API authentication (optional, increases rate limit)
    #[arg(short, long, env = "GITHUB_TOKEN")]
    token: Option<String>,

    /// Output format: "moon" for moon generate args, "json" for raw JSON
    #[arg(short, long, default_value = "moon")]
    format: String,
}

#[derive(Deserialize, Debug)]
struct Release {
    tag_name: String,
    assets: Vec<Asset>,
}

#[derive(Deserialize, Debug)]
struct Asset {
    name: String,
}

#[derive(Debug)]
struct PluginInfo {
    name: String,
    github_owner: String,
    github_repo: String,
    binary_name: String,
    linux_download_file: String,
    macos_download_file: String,
    windows_download_file: String,
    has_checksum: bool,
    checksum_file: String,
    arch_aarch64: String,
    arch_x86_64: String,
    needs_unpack: bool,
    version: String,
}

fn main() {
    let cli = Cli::parse();

    let client = Client::new();
    let url = format!(
        "https://api.github.com/repos/{}/{}/releases/latest",
        cli.owner, cli.repo
    );

    let mut request = client
        .get(&url)
        .header(USER_AGENT, "plugin-gen/0.1.0")
        .header(ACCEPT, "application/vnd.github+json");

    if let Some(ref token) = cli.token {
        request = request.header("Authorization", format!("Bearer {}", token));
    }

    let response = request.send().expect("Failed to fetch release");
    if !response.status().is_success() {
        eprintln!(
            "Error: GitHub API returned status {}",
            response.status()
        );
        std::process::exit(1);
    }

    let release: Release = response.json().expect("Failed to parse release JSON");
    let info = analyze_release(&cli.owner, &cli.repo, &release);

    match cli.format.as_str() {
        "moon" => print_moon_args(&info),
        "json" => print_json(&info),
        _ => {
            eprintln!("Unknown format: {}", cli.format);
            std::process::exit(1);
        }
    }
}

fn analyze_release(owner: &str, repo: &str, release: &Release) -> PluginInfo {
    let version = release
        .tag_name
        .strip_prefix('v')
        .unwrap_or(&release.tag_name)
        .to_string();

    let asset_names: Vec<&str> = release.assets.iter().map(|a| a.name.as_str()).collect();

    let binary_name = guess_binary_name(repo);
    let arch_mappings = detect_arch_mappings(&asset_names);
    let has_checksum = detect_checksum(&asset_names);
    let checksum_file = find_checksum_pattern(&asset_names, &version);
    let needs_unpack = detect_needs_unpack(&asset_names);

    let linux_download = find_platform_asset(&asset_names, &version, "linux", &arch_mappings);
    let macos_download = find_platform_asset(&asset_names, &version, "macos", &arch_mappings);
    let windows_download = find_platform_asset(&asset_names, &version, "windows", &arch_mappings);

    PluginInfo {
        name: binary_name.clone(),
        github_owner: owner.to_string(),
        github_repo: repo.to_string(),
        binary_name,
        linux_download_file: linux_download,
        macos_download_file: macos_download,
        windows_download_file: windows_download,
        has_checksum,
        checksum_file,
        arch_aarch64: arch_mappings
            .get("aarch64")
            .cloned()
            .unwrap_or_else(|| "arm64".to_string()),
        arch_x86_64: arch_mappings
            .get("x86_64")
            .cloned()
            .unwrap_or_else(|| "amd64".to_string()),
        needs_unpack,
        version,
    }
}

fn guess_binary_name(repo: &str) -> String {
    repo.to_lowercase()
}

fn detect_arch_mappings(assets: &[&str]) -> HashMap<String, String> {
    let mut mappings = HashMap::new();

    let aarch64_patterns = ["arm64", "aarch64", "ARM64"];
    let x86_64_patterns = ["amd64", "x86_64", "x64", "64bit"];

    for asset in assets {
        for pattern in &aarch64_patterns {
            if asset.contains(pattern) {
                mappings.insert("aarch64".to_string(), pattern.to_lowercase());
                break;
            }
        }
        for pattern in &x86_64_patterns {
            if asset.contains(pattern) {
                mappings.insert("x86_64".to_string(), pattern.to_lowercase());
                break;
            }
        }
    }

    if !mappings.contains_key("aarch64") {
        mappings.insert("aarch64".to_string(), "arm64".to_string());
    }
    if !mappings.contains_key("x86_64") {
        mappings.insert("x86_64".to_string(), "amd64".to_string());
    }

    mappings
}

fn detect_checksum(assets: &[&str]) -> bool {
    assets.iter().any(|a| {
        a.contains("checksum")
            || a.contains("sha256")
            || a.contains("SHA256")
            || a.ends_with(".sha256")
            || a.ends_with(".sha512")
    })
}

fn find_checksum_pattern(assets: &[&str], version: &str) -> String {
    for asset in assets {
        if asset.contains("checksum")
            || asset.contains("sha256")
            || asset.contains("SHA256")
            || asset.ends_with(".sha256")
        {
            return asset.replace(version, "{version}");
        }
    }
    String::new()
}

fn detect_needs_unpack(assets: &[&str]) -> bool {
    let archive_exts = [".tar.gz", ".tgz", ".zip", ".tar.xz", ".tar.bz2"];
    let linux_assets: Vec<&&str> = assets
        .iter()
        .filter(|a| a.to_lowercase().contains("linux"))
        .collect();

    if linux_assets.is_empty() {
        return assets
            .iter()
            .any(|a| archive_exts.iter().any(|ext| a.ends_with(ext)));
    }

    linux_assets
        .iter()
        .any(|a| archive_exts.iter().any(|ext| a.ends_with(ext)))
}

fn find_platform_asset(
    assets: &[&str],
    version: &str,
    platform: &str,
    arch_mappings: &HashMap<String, String>,
) -> String {
    let os_patterns: Vec<&str> = match platform {
        "linux" => vec!["linux", "Linux"],
        "macos" => vec!["darwin", "Darwin", "macos", "macOS", "apple"],
        "windows" => vec!["windows", "Windows", "win64", "win"],
        _ => vec![],
    };

    let archive_exts = [".tar.gz", ".tgz", ".zip", ".tar.xz"];
    let binary_exts = ["", ".exe"];
    let checksum_patterns = ["checksum", "sha256", "SHA256", ".sha"];

    let x86_64_arch = arch_mappings
        .get("x86_64")
        .map(|s| s.as_str())
        .unwrap_or("amd64");

    for asset in assets {
        let lower = asset.to_lowercase();
        if checksum_patterns.iter().any(|p| lower.contains(p)) {
            continue;
        }

        let matches_os = os_patterns.iter().any(|p| asset.contains(p));
        if !matches_os {
            continue;
        }

        let has_archive_ext = archive_exts.iter().any(|ext| asset.ends_with(ext));
        let has_binary_ext = if platform == "windows" {
            asset.ends_with(".exe") || asset.ends_with(".zip")
        } else {
            binary_exts.iter().any(|ext| asset.ends_with(ext))
        };

        if has_archive_ext || has_binary_ext {
            let contains_arch = asset.contains(x86_64_arch)
                || asset.contains("x86_64")
                || asset.contains("amd64")
                || asset.contains("64");

            if contains_arch || !has_multiple_arches(assets, &os_patterns) {
                let pattern = asset
                    .replace(version, "{version}")
                    .replace(x86_64_arch, "{arch}")
                    .replace("x86_64", "{arch}")
                    .replace("amd64", "{arch}");
                return pattern;
            }
        }
    }

    format!("MANUAL_ENTRY_NEEDED_{}", platform)
}

fn has_multiple_arches(assets: &[&str], os_patterns: &[&str]) -> bool {
    let os_assets: Vec<&&str> = assets
        .iter()
        .filter(|a| os_patterns.iter().any(|p| a.contains(p)))
        .collect();

    let arch_indicators = ["arm64", "aarch64", "amd64", "x86_64", "x64", "arm"];
    let mut found_arches = 0;
    for arch in &arch_indicators {
        if os_assets.iter().any(|a| a.contains(arch)) {
            found_arches += 1;
        }
    }
    found_arches > 1
}

fn print_moon_args(info: &PluginInfo) {
    println!("# moon generate plugin -- \\");
    println!("#   --name \"{}\" \\", info.name);
    println!("#   --github_owner \"{}\" \\", info.github_owner);
    println!("#   --github_repo \"{}\" \\", info.github_repo);
    println!("#   --binary_name \"{}\" \\", info.binary_name);
    println!("#   --version_command \"--version\" \\");
    println!("#   --linux_download_file \"{}\" \\", info.linux_download_file);
    println!("#   --macos_download_file \"{}\" \\", info.macos_download_file);
    println!("#   --windows_download_file \"{}\" \\", info.windows_download_file);
    println!("#   --has_checksum {} \\", info.has_checksum);
    println!("#   --checksum_file \"{}\" \\", info.checksum_file);
    println!("#   --arch_aarch64 \"{}\" \\", info.arch_aarch64);
    println!("#   --arch_x86_64 \"{}\" \\", info.arch_x86_64);
    println!("#   --needs_unpack {}", info.needs_unpack);
    println!();
    println!("# Detected version: {}", info.version);
    println!("# Add to .prototools: {} = \"{}\"", info.name, info.version);
}

fn print_json(info: &PluginInfo) {
    let json = serde_json::json!({
        "name": info.name,
        "github_owner": info.github_owner,
        "github_repo": info.github_repo,
        "binary_name": info.binary_name,
        "version_command": "--version",
        "linux_download_file": info.linux_download_file,
        "macos_download_file": info.macos_download_file,
        "windows_download_file": info.windows_download_file,
        "has_checksum": info.has_checksum,
        "checksum_file": info.checksum_file,
        "arch_aarch64": info.arch_aarch64,
        "arch_x86_64": info.arch_x86_64,
        "needs_unpack": info.needs_unpack,
        "version": info.version,
    });
    println!("{}", serde_json::to_string_pretty(&json).unwrap());
}
