import { describe, expect, test } from "bun:test";

describe("cross consistency", () => {
    test("go-ts wire consistency passes", () => {
        const repoRoot = decodeURIComponent(new URL("../..", import.meta.url).pathname);
        const proc = Bun.spawnSync({
            cmd: ["go", "test", "./go/sb", "-run", "TestCrossLanguageWireConsistency", "-count=1"],
            cwd: repoRoot,
            env: { ...process.env, GOCACHE: "/tmp/go-build" },
            stdout: "pipe",
            stderr: "pipe",
        });
        const stdout = new TextDecoder().decode(proc.stdout).trim();
        const stderr = new TextDecoder().decode(proc.stderr).trim();
        expect(proc.exitCode, `${stdout}\n${stderr}`.trim()).toBe(0);
    });
});
