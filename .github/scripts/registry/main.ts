export interface PluginToml {
    name: string;
    type: string;
    platform: Record<string, PlatformConfig>;
    resolve?: ResolveConfig;
}

interface PlatformConfig {
    "download-file"?: string;
    "checksum-file"?: string;
    "bin-path"?: string;
}

interface ResolveConfig {
    "git-url"?: string;
    "github-repo"?: string;
}

export interface RegistryEntry {
    id: string;
    locator: string;
    format: string;
    name: string;
    description: string;
    author: string;
    homepageUrl: string | null;
    repositoryUrl: string;
    bins: string[];
}

interface RegistryDocument {
    $schema: string;
    version: number;
    plugins: RegistryEntry[];
}

const REPO_OWNER = "ageha734";
const REPO_NAME = "proto-plugins";
const BRANCH = "master";

export function parsePluginToml(content: string): PluginToml {
    const result: PluginToml = {
        name: "",
        type: "",
        platform: {},
        resolve: {},
    };

    let currentSection = "";
    let currentPlatform = "";

    for (const line of content.split("\n")) {
        const trimmed = line.trim();
        if (!trimmed || trimmed.startsWith("#")) continue;

        const sectionMatch = trimmed.match(/^\[(.+)\]$/);
        if (sectionMatch) {
            currentSection = sectionMatch[1];
            const platformMatch = currentSection.match(/^platform\.(.+)$/);
            if (platformMatch) {
                currentPlatform = platformMatch[1];
                result.platform[currentPlatform] = {};
            }
            continue;
        }

        const kvMatch = trimmed.match(/^([^=]+)\s*=\s*"([^"]*)"$/);
        if (!kvMatch) continue;

        const key = kvMatch[1].trim();
        const value = kvMatch[2];

        if (currentSection === "") {
            if (key === "name") result.name = value;
            if (key === "type") result.type = value;
        } else if (currentSection.startsWith("platform.") && currentPlatform) {
            const platform = result.platform[currentPlatform];
            if (platform) {
                (platform as Record<string, string>)[key] = value;
            }
        } else if (currentSection === "resolve") {
            if (!result.resolve) result.resolve = {};
            (result.resolve as Record<string, string>)[key] = value;
        }
    }

    return result;
}

export function extractBins(plugin: PluginToml): string[] {
    const bins = new Set<string>();

    for (const config of Object.values(plugin.platform)) {
        const binPath = config["bin-path"];
        if (binPath) {
            const name = binPath.replace(/\.exe$/, "");
            bins.add(name);
        }
    }

    return [...bins].sort();
}

export function extractHomepageUrl(plugin: PluginToml): string | null {
    return plugin.resolve?.["git-url"] ?? plugin.resolve?.["github-repo"] ?? null;
}

export function buildRegistryEntry(id: string, plugin: PluginToml): RegistryEntry {
    return {
        id,
        locator: `https://raw.githubusercontent.com/${REPO_OWNER}/${REPO_NAME}/${BRANCH}/toml/${id}.toml`,
        format: "toml",
        name: plugin.name || id,
        description: "",
        author: REPO_OWNER,
        homepageUrl: extractHomepageUrl(plugin),
        repositoryUrl: `https://github.com/${REPO_OWNER}/${REPO_NAME}/blob/${BRANCH}/toml/${id}.toml`,
        bins: extractBins(plugin),
    };
}

export async function generateAllEntries(tomlDir: string): Promise<RegistryEntry[]> {
    const entries: RegistryEntry[] = [];

    for await (const entry of Deno.readDir(tomlDir)) {
        if (!entry.isFile || !entry.name.endsWith(".toml")) continue;

        const id = entry.name.replace(/\.toml$/, "");
        const content = await Deno.readTextFile(`${tomlDir}/${entry.name}`);
        const plugin = parsePluginToml(content);
        entries.push(buildRegistryEntry(id, plugin));
    }

    return entries.sort((a, b) => a.id.localeCompare(b.id));
}

export function mergeIntoRegistry(
    existing: RegistryDocument,
    newEntries: RegistryEntry[],
): RegistryDocument {
    const existingMap = new Map(existing.plugins.map((p) => [p.id, p]));

    for (const entry of newEntries) {
        existingMap.set(entry.id, entry);
    }

    const plugins = [...existingMap.values()].sort((a, b) => a.id.localeCompare(b.id));

    return {
        ...existing,
        plugins,
    };
}

async function fetchUpstreamRegistry(): Promise<RegistryDocument> {
    const url =
        "https://raw.githubusercontent.com/moonrepo/proto/master/registry/data/third-party.json";
    const resp = await fetch(url);
    if (!resp.ok) {
        throw new Error(`Failed to fetch upstream registry: ${resp.status}`);
    }
    return await resp.json();
}

async function main() {
    const root = Deno.cwd();
    const tomlDir = `${root}/toml`;

    console.log("Generating registry entries from TOML plugins...");
    const entries = await generateAllEntries(tomlDir);
    console.log(`Generated ${entries.length} entries`);

    console.log("Fetching upstream registry...");
    const upstream = await fetchUpstreamRegistry();

    const merged = mergeIntoRegistry(upstream, entries);
    const outputPath = `${root}/registry-output.json`;
    await Deno.writeTextFile(outputPath, JSON.stringify(merged, null, 2) + "\n");
    console.log(`Written merged registry to ${outputPath}`);

    console.log("\nEntries added/updated:");
    for (const entry of entries) {
        console.log(`  - ${entry.id}: ${entry.locator}`);
    }
}

if (import.meta.main) {
    main();
}
