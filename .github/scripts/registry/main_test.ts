import {
    buildRegistryEntry,
    extractBins,
    extractHomepageUrl,
    parsePluginToml,
    type PluginToml,
    type RegistryEntry,
} from "./main.ts";

function assertEquals<T>(actual: T, expected: T): void {
    const a = JSON.stringify(actual);
    const e = JSON.stringify(expected);
    if (a !== e) {
        throw new Error(`Assertion failed:\n  actual:   ${a}\n  expected: ${e}`);
    }
}

const SAMPLE_TOML = `
name = "pinact"
type = "cli"

[platform.linux]
download-file = "pinact_linux_{arch}.tar.gz"
checksum-file = "pinact_{version}_checksums.txt"
bin-path = "pinact"

[platform.macos]
download-file = "pinact_darwin_{arch}.tar.gz"
checksum-file = "pinact_{version}_checksums.txt"
bin-path = "pinact"

[platform.windows]
download-file = "pinact_windows_{arch}.zip"
checksum-file = "pinact_{version}_checksums.txt"
bin-path = "pinact.exe"

[install]
download-url = "https://github.com/suzuki-shunsuke/pinact/releases/download/v{version}/{download_file}"
checksum-url = "https://github.com/suzuki-shunsuke/pinact/releases/download/v{version}/{checksum_file}"
unpack = false

[install.arch]
aarch64 = "arm64"
x86_64 = "amd64"

[resolve]
git-url = "https://github.com/suzuki-shunsuke/pinact"
`;

Deno.test("parsePluginToml - parses name and type", () => {
    const result = parsePluginToml(SAMPLE_TOML);
    assertEquals(result.name, "pinact");
    assertEquals(result.type, "cli");
});

Deno.test("parsePluginToml - parses platforms", () => {
    const result = parsePluginToml(SAMPLE_TOML);
    assertEquals(Object.keys(result.platform).sort(), ["linux", "macos", "windows"]);
});

Deno.test("parsePluginToml - parses resolve git-url", () => {
    const result = parsePluginToml(SAMPLE_TOML);
    assertEquals(result.resolve?.["git-url"], "https://github.com/suzuki-shunsuke/pinact");
});

Deno.test("extractBins - extracts unique bin names without extension", () => {
    const plugin = parsePluginToml(SAMPLE_TOML);
    const bins = extractBins(plugin);
    assertEquals(bins, ["pinact"]);
});

Deno.test("extractBins - handles multiple different binaries", () => {
    const plugin: PluginToml = {
        name: "test",
        type: "cli",
        platform: {
            linux: { "bin-path": "mytool" },
            windows: { "bin-path": "mytool.exe" },
        },
        resolve: {},
    };
    const bins = extractBins(plugin);
    assertEquals(bins, ["mytool"]);
});

Deno.test("extractHomepageUrl - extracts from git-url", () => {
    const plugin = parsePluginToml(SAMPLE_TOML);
    const url = extractHomepageUrl(plugin);
    assertEquals(url, "https://github.com/suzuki-shunsuke/pinact");
});

Deno.test("extractHomepageUrl - returns null when no resolve", () => {
    const plugin: PluginToml = {
        name: "test",
        type: "cli",
        platform: {},
        resolve: {},
    };
    const url = extractHomepageUrl(plugin);
    assertEquals(url, null);
});

Deno.test("buildRegistryEntry - generates correct entry", () => {
    const plugin = parsePluginToml(SAMPLE_TOML);
    const entry = buildRegistryEntry("pinact", plugin);
    const expected: RegistryEntry = {
        id: "pinact",
        locator: "https://raw.githubusercontent.com/ageha734/proto-plugins/master/toml/pinact.toml",
        format: "toml",
        name: "pinact",
        description: "",
        author: "ageha734",
        homepageUrl: "https://github.com/suzuki-shunsuke/pinact",
        repositoryUrl: "https://github.com/ageha734/proto-plugins/blob/master/toml/pinact.toml",
        bins: ["pinact"],
    };
    assertEquals(entry, expected);
});

Deno.test("buildRegistryEntry - uses plugin name field for display name", () => {
    const toml = `
name = "Action Lint"
type = "cli"

[platform.linux]
bin-path = "actionlint"

[resolve]
git-url = "https://github.com/rhysd/actionlint"
`;
    const plugin = parsePluginToml(toml);
    const entry = buildRegistryEntry("actionlint", plugin);
    assertEquals(entry.id, "actionlint");
    assertEquals(entry.name, "Action Lint");
});
