export interface RunReport {
    actions: Action[];
}

interface Action {
    node: ActionNode;
    operations: Operation[];
    status: string;
}

interface ActionNode {
    action: string;
    params: ActionParams;
}

interface ActionParams {
    target: string;
}

interface Operation {
    meta: OperationMeta;
}

interface OperationMeta {
    type: string;
    command?: string;
}

interface TargetIdentity {
    project: string;
    task: string;
}

interface LogReader {
    (project: string, task: string): { stdout: string; stderr: string };
}

const STATUS_BADGES: Record<string, string> = {
    running: bgGreen(" RUNNING "),
    passed: bgGreen(" PASS "),
    failed: bgRed(" FAIL "),
    "timed-out": bgRed(" TIMED OUT "),
    aborted: bgRed(" ABORTED "),
    invalid: bgRed(" INVALID "),
    "failed-and-abort": bgRed(" FAILED AND ABORT "),
    skipped: bgBlue(" SKIP "),
    cached: bgBlue(" CACHED "),
    "cached-from-remote": bgBlue(" REMOTE CACHED "),
};

const STD_BADGES = {
    out: bgDarkGray(`　${green("⏺")} STDOUT　`),
    err: bgDarkGray(`　${red("⏺")} STDERR　`),
};

export function parseTarget(target: string): TargetIdentity {
    const parts = target.split(":");
    if (parts.length !== 2) {
        return { project: "unknown", task: "unknown" };
    }
    return { project: parts[0], task: parts[1] };
}

export function formatBadge(status: string): string {
    return STATUS_BADGES[status] ?? status;
}

export function formatOutput(report: RunReport, readLogs: LogReader): string {
    const lines: string[] = [];

    for (const action of report.actions) {
        if (action.node.action !== "run-task") continue;

        const target = parseTarget(action.node.params.target);
        const targetStr = `${target.project}:${target.task}`;
        const command = commandOf(action);
        const { stdout, stderr } = readLogs(target.project, target.task);

        const badge = formatBadge(action.status);
        lines.push(`::group::${badge} ${bold(targetStr)}`);

        if (command) {
            lines.push(blue(`$ ${command}`));
        }

        if (stdout.trim()) {
            lines.push(STD_BADGES.out);
            lines.push(stdout);
        }

        if (stderr.trim()) {
            lines.push(STD_BADGES.err);
            lines.push(stderr);
        }

        lines.push("::endgroup::");
    }

    return lines.join("\n");
}

export async function loadReport(workspaceRoot: string): Promise<RunReport | null> {
    const candidates = ["ciReport.json", "runReport.json"];

    for (const fileName of candidates) {
        const reportPath = `${workspaceRoot}/.moon/cache/${fileName}`;
        try {
            const data = await Deno.readTextFile(reportPath);
            return JSON.parse(data) as RunReport;
        } catch {
            continue;
        }
    }

    return null;
}

function commandOf(action: Action): string {
    for (const op of action.operations) {
        if (op.meta.type === "task-execution") {
            return op.meta.command ?? "";
        }
    }
    return "";
}

function readLogsFromDisk(workspaceRoot: string): LogReader {
    return (project: string, task: string) => {
        const stateDir = `${workspaceRoot}/.moon/cache/states/${project}/${task}`;
        let stdout = "";
        let stderr = "";

        try {
            stdout = Deno.readTextFileSync(`${stateDir}/stdout.log`);
        } catch {
            // file may not exist
        }

        try {
            stderr = Deno.readTextFileSync(`${stateDir}/stderr.log`);
        } catch {
            // file may not exist
        }

        return { stdout, stderr };
    };
}

async function main() {
    const root = Deno.cwd();
    const report = await loadReport(root);

    if (!report) {
        console.log("::warning::Run report does not exist, has `moon ci` or `moon run` ran?");
        return;
    }

    const output = formatOutput(report, readLogsFromDisk(root));
    console.log(output);
}

function bgGreen(text: string): string {
    return `\x1b[42m${text}\x1b[49m`;
}
function bgRed(text: string): string {
    return `\x1b[41m${text}\x1b[49m`;
}
function bgBlue(text: string): string {
    return `\x1b[44m${text}\x1b[49m`;
}
function bgDarkGray(text: string): string {
    return `\x1b[48;5;236m${text}\x1b[49m`;
}
function bold(text: string): string {
    return `\x1b[1m${text}\x1b[22m`;
}
function green(text: string): string {
    return `\x1b[32m${text}\x1b[39m`;
}
function red(text: string): string {
    return `\x1b[31m${text}\x1b[39m`;
}
function blue(text: string): string {
    return `\x1b[34m${text}\x1b[39m`;
}

if (import.meta.main) {
    main();
}
