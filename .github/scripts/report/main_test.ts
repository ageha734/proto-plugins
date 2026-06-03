import {
    formatBadge,
    formatOutput,
    loadReport,
    parseTarget,
    type RunReport,
} from "./main.ts";

function assertEquals<T>(actual: T, expected: T): void {
    const a = JSON.stringify(actual);
    const e = JSON.stringify(expected);
    if (a !== e) {
        throw new Error(`Assertion failed:\n  actual:   ${a}\n  expected: ${e}`);
    }
}

const FIXTURE_REPORT: RunReport = {
    actions: [
        {
            node: {
                action: "run-task",
                params: { target: "toml:lint" },
            },
            operations: [
                {
                    meta: {
                        type: "task-execution",
                        command: "golangci-lint run",
                    },
                },
            ],
            status: "passed",
        },
        {
            node: {
                action: "run-task",
                params: { target: "toml:test" },
            },
            operations: [
                {
                    meta: {
                        type: "task-execution",
                        command: "go test",
                    },
                },
            ],
            status: "failed",
        },
        {
            node: {
                action: "run-task",
                params: { target: "workspace:lint" },
            },
            operations: [
                {
                    meta: {
                        type: "task-execution",
                        command: "dprint check",
                    },
                },
            ],
            status: "cached",
        },
        {
            node: {
                action: "sync-workspace",
                params: { target: "" },
            },
            operations: [],
            status: "passed",
        },
    ],
};

Deno.test("parseTarget - splits project and task", () => {
    const result = parseTarget("toml:lint");
    assertEquals(result, { project: "toml", task: "lint" });
});

Deno.test("parseTarget - handles unknown format", () => {
    const result = parseTarget("invalid");
    assertEquals(result, { project: "unknown", task: "unknown" });
});

Deno.test("formatBadge - returns correct badge for known statuses", () => {
    const passed = formatBadge("passed");
    assertEquals(passed.includes("PASS"), true);

    const failed = formatBadge("failed");
    assertEquals(failed.includes("FAIL"), true);

    const cached = formatBadge("cached");
    assertEquals(cached.includes("CACHED"), true);

    const skipped = formatBadge("skipped");
    assertEquals(skipped.includes("SKIP"), true);
});

Deno.test("formatBadge - returns raw status for unknown", () => {
    const result = formatBadge("unknown-status");
    assertEquals(result, "unknown-status");
});

Deno.test("formatOutput - filters only run-task actions", () => {
    const output = formatOutput(FIXTURE_REPORT, () => ({ stdout: "", stderr: "" }));
    const lines = output.split("\n");

    const groups = lines.filter((l) => l.startsWith("::group::"));
    assertEquals(groups.length, 3);

    const hasSync = lines.some((l) => l.includes("sync-workspace"));
    assertEquals(hasSync, false);
});

Deno.test("formatOutput - includes command in output", () => {
    const output = formatOutput(FIXTURE_REPORT, () => ({ stdout: "ok", stderr: "" }));
    assertEquals(output.includes("$ golangci-lint run"), true);
    assertEquals(output.includes("$ go test"), true);
});

Deno.test("formatOutput - includes stdout and stderr when present", () => {
    const output = formatOutput(FIXTURE_REPORT, (project, task) => {
        if (project === "toml" && task === "test") {
            return { stdout: "FAIL: TestFoo", stderr: "error occurred" };
        }
        return { stdout: "", stderr: "" };
    });
    assertEquals(output.includes("FAIL: TestFoo"), true);
    assertEquals(output.includes("error occurred"), true);
});

Deno.test("formatOutput - skips empty stdout/stderr", () => {
    const output = formatOutput(FIXTURE_REPORT, () => ({ stdout: "", stderr: "" }));
    assertEquals(output.includes("STDOUT"), false);
    assertEquals(output.includes("STDERR"), false);
});

Deno.test("loadReport - returns null for non-existent file", async () => {
    const result = await loadReport("/non/existent/path");
    assertEquals(result, null);
});
