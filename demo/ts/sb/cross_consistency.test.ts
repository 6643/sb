import { describe, expect, test } from "bun:test";

const goRunPattern = (name: string): string => `^${name}$`;
const hasGo = typeof Bun.which("go") === "string";

const runGoCrossTest = (name: string): void => {
    if (!hasGo) return;
    const repoRoot = decodeURIComponent(new URL("../../..", import.meta.url).pathname);
    const proc = Bun.spawnSync({
        cmd: ["go", "test", "./demo/go/sb", "-run", goRunPattern(name), "-count=1"],
        cwd: repoRoot,
        env: { ...process.env, GOCACHE: "/tmp/go-build" },
        stdout: "pipe",
        stderr: "pipe",
    });
    const stdout = new TextDecoder().decode(proc.stdout).trim();
    const stderr = new TextDecoder().decode(proc.stderr).trim();
    expect(proc.exitCode, `${stdout}\n${stderr}`.trim()).toBe(0);
};

describe("cross consistency", () => {
    test("go test pattern is anchored", () => {
        expect(goRunPattern("TestCrossLanguageWireConsistency")).toBe("^TestCrossLanguageWireConsistency$");
    });

    (hasGo ? test : test.skip)("go-ts wire consistency passes", () => {
        runGoCrossTest("TestCrossLanguageWireConsistency");
    }, 20000);

    (hasGo ? test : test.skip)("go-ts random wire consistency passes", () => {
        runGoCrossTest("TestCrossLanguageWireConsistencyRandom");
    }, 20000);

    (hasGo ? test : test.skip)("go-ts malformed wire cases reject", () => {
        runGoCrossTest("TestCrossLanguageWireRejectsMalformedInputs");
    }, 20000);
});
