# TS Runtime Consolidation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 TS generator 的 runtime 输出从当前细粒度 `runtime_*.ts` 集合收口为粗粒度文件集, 同时保留 `type.ts` 兼容 barrel 与 `runtime_core.ts` / `runtime_meta.ts` 两个薄中间层。

**Architecture:** 先用 focused red tests 锁定最终允许落盘文件、`type.ts` 兼容语义与 TS/Go 行为不变, 再把 text/bin/list/struct/header/base 的细模板回并到粗模板中。随后统一收掉 embed 与写盘清单里的细文件项, 为 generator 增加确定性的 stale runtime 清理规则, 最后重生成 fixtures 并跑 Go/Bun 全量回归。

**Tech Stack:** Go generator templates, embedded TS runtime assets, Bun tests, generated fixtures under `fixtures/ts/sb`, Go tests under `internal` and `fixtures/go/sb`

---

## File Map

- `docs/superpowers/specs/2026-03-23-ts-runtime-consolidation-design.md` - 已批准设计, 执行前先对齐边界与验收标准。
- `internal/tpl_ts_render_test.go` - generator focused regression; 锁定允许落盘文件、barrel 兼容、导入方向与中间层约束。
- `fixtures/ts/sb/runtime.test.ts` - 通过 `./type` 导入 runtime API 的兼容与协议语义验证。
- `fixtures/ts/sb/enum_smoke.test.ts` - enum 公开 API smoke。
- `fixtures/ts/sb/rpc_smoke.test.ts` - `RpcClient` 与生成方法入口 smoke。
- `fixtures/ts/sb/rpc_fetch.test.ts` - RPC 传输行为守卫。
- `fixtures/ts/sb/struct_smoke.test.ts` - 新增的 struct 公开 API smoke, 验证 `new*` / `validate*` / `get*` / `set*` / `eq*` 用法稳定。
- `internal/tpl_ts_runtime.ts.txt` - `type.ts` 兼容 barrel 模板, 只允许 re-export。
- `internal/tpl_ts_runtime_core.ts.txt` - `runtime_core.ts` barrel 模板, 只聚合基础 runtime 域。
- `internal/tpl_ts_runtime_meta.ts.txt` - `runtime_meta.ts` barrel 模板, 只聚合 enum/struct 域。
- `internal/tpl_ts_runtime_base.ts.txt` - 合并 buffer/bit/int/float/equal/result 后的基础 runtime 实现。
- `internal/tpl_ts_runtime_header.ts.txt` - 合并 bitmap/state block/header codec 后的 header runtime 实现。
- `internal/tpl_ts_runtime_text.ts.txt` - 合并 utf8/state/io/codec 后的 text runtime 实现。
- `internal/tpl_ts_runtime_bin.ts.txt` - 合并 state/length/io/codec 后的 bin runtime 实现。
- `internal/tpl_ts_runtime_list.ts.txt` - 合并 count/default/text/bin/bool list 后的 list runtime 实现。
- `internal/tpl_ts_runtime_enum.ts.txt` - enum runtime 实现, 保持单文件真实实现。
- `internal/tpl_ts_runtime_struct.ts.txt` - 合并 meta/read/write/field builder 后的 struct runtime 实现。
- `internal/tpl_ts_runtime_render.go` - embed 清单与 render 函数; 最终只能暴露粗模板与两个薄中间层。
- `internal/tpl_ts_render.go` - TS generator 写盘清单、runtime 文件允许集合、旧细文件清理逻辑。
- `internal/generated_file_writer.go` - 若需要通用删除 helper, 在这里扩展生成目录写入能力。
- `fixtures/ts/sb/*` - 重生成后的目标产物; 用来验证旧细文件自动清理。

## Implementation Notes

- 全程遵循 @test-driven-development: 每个任务先写 red test, 确认失败原因正确, 再写最小实现。
- 最终完成前遵循 @verification-before-completion: 必须保留 focused regression、真实 generator 重生成、Go/Bun 全量结果。
- 不要手工删除 `fixtures/ts/sb` 中的旧细粒度 runtime 文件来制造“通过”; 必须让 generator 在真实生成时自动清理。
- `runtime_core.ts` 与 `runtime_meta.ts` 只能是薄 barrel. 一旦出现真实实现代码, 视为偏离设计。

### Task 1: Lock the coarse runtime output contract

**Files:**
- Modify: `internal/tpl_ts_render_test.go`
- Test: `internal/tpl_ts_render_test.go`

- [ ] **Step 1: Write the failing generator assertions**

Add/replace assertions so the temp TS output only permits:

```go
allowed := []string{
    "type.ts",
    "runtime_base.ts",
    "runtime_header.ts",
    "runtime_text.ts",
    "runtime_bin.ts",
    "runtime_list.ts",
    "runtime_struct.ts",
    "runtime_enum.ts",
    "runtime_core.ts",
    "runtime_meta.ts",
}
```

Also assert:
- disallowed files like `runtime_text_io.ts`, `runtime_bin_length.ts`, `runtime_struct_scalar_field.ts` do not exist
- `runtime_core.ts` and `runtime_meta.ts` contain only `export * from`
- `type.ts` still re-exports representative public symbols through the coarse modules
- generated `enum.ts`, one representative `struct_*.ts`, and `rpc.ts` do not contain `from "./type"`
- generated `enum.ts` / `struct_*.ts` / `rpc.ts` do contain the expected coarse imports such as `from "./runtime_core"` and `from "./runtime_meta"`

- [ ] **Step 2: Run the focused generator test and verify RED**

Run: `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`

Expected: FAIL because the generator still writes fine-grained runtime files and current content assertions still target split files.

- [ ] **Step 3: Add the minimal file-list/content helpers needed by the test**

If the current test lacks file-existence helpers, add the smallest possible helper in the same file, for example:

```go
func assertFileMissing(t *testing.T, path string) {
    t.Helper()
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        t.Fatalf("expected file missing: %s err=%v", path, err)
    }
}
```

- [ ] **Step 4: Re-run the focused generator test and keep it RED for the correct reason**

Run: `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`

Expected: still FAIL, now specifically because coarse-file expectations are not yet implemented.

- [ ] **Step 5: Commit the red test scaffold**

```bash
git add internal/tpl_ts_render_test.go
git commit -m "test: lock ts runtime consolidation output contract"
```

### Task 2: Lock TS public API compatibility before merging templates

**Files:**
- Modify: `fixtures/ts/sb/runtime.test.ts`
- Modify: `fixtures/ts/sb/enum_smoke.test.ts`
- Modify: `fixtures/ts/sb/rpc_smoke.test.ts`
- Modify: `fixtures/ts/sb/rpc_fetch.test.ts`
- Create: `fixtures/ts/sb/struct_smoke.test.ts`
- Test: `fixtures/ts/sb/runtime.test.ts`
- Test: `fixtures/ts/sb/enum_smoke.test.ts`
- Test: `fixtures/ts/sb/struct_smoke.test.ts`
- Test: `fixtures/ts/sb/rpc_smoke.test.ts`
- Test: `fixtures/ts/sb/rpc_fetch.test.ts`

- [ ] **Step 1: Write the missing struct public API smoke test**

Create a focused Bun test that uses representative generated APIs without peeking into internals, for example:

```ts
import { describe, expect, test } from "bun:test";
import { newSimInfo, validateSimInfo, eqSimInfo, setSimInfo, getSimInfo } from "./struct_sim_info";
import { Buffer } from "./type";

test("struct public api remains usable", () => {
  const value = newSimInfo();
  expect(validateSimInfo(value)).toBeUndefined();
  const buf = new Buffer();
  expect(setSimInfo(buf, value)).toBeUndefined();
  const [decoded, err] = getSimInfo(new Buffer(buf.bytes));
  expect(err).toBeUndefined();
  expect(eqSimInfo(value, decoded)).toBe(true);
});
```

- [ ] **Step 2: Tighten `./type` compatibility assertions**

In `fixtures/ts/sb/runtime.test.ts`, add a minimal barrel-compatibility test that exercises representative imports like `Buffer`, `writeHeader`, `getText`, and `setBin` via `./type`.

- [ ] **Step 3: Run the focused Bun smoke suite and verify RED/GREEN expectations**

Run: `bun test fixtures/ts/sb/runtime.test.ts fixtures/ts/sb/enum_smoke.test.ts fixtures/ts/sb/struct_smoke.test.ts fixtures/ts/sb/rpc_smoke.test.ts fixtures/ts/sb/rpc_fetch.test.ts`

Expected:
- any new assertion written for missing coverage fails first
- after minimal test-only fixes, the focused suite passes on the current baseline

- [ ] **Step 4: Keep the TS smoke suite green before touching runtime templates**

Re-run: `bun test fixtures/ts/sb/runtime.test.ts fixtures/ts/sb/enum_smoke.test.ts fixtures/ts/sb/struct_smoke.test.ts fixtures/ts/sb/rpc_smoke.test.ts fixtures/ts/sb/rpc_fetch.test.ts`

Expected: PASS.

- [ ] **Step 5: Commit the compatibility tests**

```bash
git add fixtures/ts/sb/runtime.test.ts fixtures/ts/sb/enum_smoke.test.ts fixtures/ts/sb/struct_smoke.test.ts fixtures/ts/sb/rpc_smoke.test.ts fixtures/ts/sb/rpc_fetch.test.ts
git commit -m "test: cover ts runtime compatibility barrels"
```

### Task 3: Merge base, header, text, bin, and list implementations into coarse templates

**Files:**
- Modify: `internal/tpl_ts_runtime_base.ts.txt`
- Modify: `internal/tpl_ts_runtime_header.ts.txt`
- Modify: `internal/tpl_ts_runtime_text.ts.txt`
- Modify: `internal/tpl_ts_runtime_bin.ts.txt`
- Modify: `internal/tpl_ts_runtime_list.ts.txt`
- Modify: `internal/tpl_ts_render_test.go`
- Test: `internal/tpl_ts_render_test.go`

- [ ] **Step 1: Extend the generator test with domain-specific coarse-file assertions**

Add exact RED assertions that the coarse files now contain real implementation, for example:

```go
assertContains(t, textText, "export const textState = (v: string): [number, Error | null] => {")
assertContains(t, binText, "export const getBinLength = (buf: Buffer, state: number, max: number, kind: string): [number, Error | null] => {")
assertContains(t, listText, "export const getDefaultList = <T>(")
assertNotContains(t, textText, "export * from \"./runtime_text_io\"")
```

- [ ] **Step 2: Run the focused generator test and verify RED**

Run: `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`

Expected: FAIL because the coarse templates are still thin barrels.

- [ ] **Step 3: Merge the fine template content into the coarse template files**

Use the approved mapping:
- `runtime_buffer*`, `runtime_bit*`, `runtime_int*`, `runtime_u24`, `runtime_float`, `runtime_equal`, `runtime_result` -> `internal/tpl_ts_runtime_base.ts.txt`
- `runtime_bitmap`, `runtime_state_header`, `runtime_state_block`, `runtime_header_codec` -> `internal/tpl_ts_runtime_header.ts.txt`
- `runtime_text_utf8`, `runtime_text_state`, `runtime_text_io`, `runtime_text_codec` -> `internal/tpl_ts_runtime_text.ts.txt`
- `runtime_bin_state`, `runtime_bin_length`, `runtime_bin_io`, `runtime_bin_codec` -> `internal/tpl_ts_runtime_bin.ts.txt`
- `runtime_bitmap_list`, `runtime_list_count`, `runtime_list_default`, `runtime_item_list`, `runtime_text_list`, `runtime_bin_list` -> `internal/tpl_ts_runtime_list.ts.txt`

Do not change behavior while merging; copy bodies intact and only normalize imports/export order as needed.

- [ ] **Step 4: Re-run the focused generator test and get GREEN for these domains**

Run: `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`

Expected: either PASS for these domains or fail only on the remaining enum/struct/embed cleanup expectations.

- [ ] **Step 5: Commit the coarse runtime domain merge**

```bash
git add internal/tpl_ts_runtime_base.ts.txt internal/tpl_ts_runtime_header.ts.txt internal/tpl_ts_runtime_text.ts.txt internal/tpl_ts_runtime_bin.ts.txt internal/tpl_ts_runtime_list.ts.txt internal/tpl_ts_render_test.go
git commit -m "refactor: merge coarse ts runtime domains"
```

### Task 4: Merge enum and struct implementations into their coarse templates

**Files:**
- Modify: `internal/tpl_ts_runtime_enum.ts.txt`
- Modify: `internal/tpl_ts_runtime_struct.ts.txt`
- Modify: `internal/tpl_ts_runtime_meta.ts.txt`
- Modify: `internal/tpl_ts_render_test.go`
- Test: `internal/tpl_ts_render_test.go`
- Test: `fixtures/ts/sb/struct_smoke.test.ts`

- [ ] **Step 1: Add RED assertions for enum/struct coarse implementations**

Add focused assertions such as:

```go
assertContains(t, enumRuntimeText, "export interface EnumMeta<T extends number>")
assertContains(t, structRuntimeText, "export interface StructMeta<T>")
assertContains(t, structRuntimeText, "export const setStruct = <T>(meta: StructMeta<T>, buf: Buffer, value: T | null | undefined): Err => {")
assertNotContains(t, metaText, "export interface StructMeta<T>")
```

- [ ] **Step 2: Run the focused generator test and verify RED**

Run: `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`

Expected: FAIL because `runtime_struct.ts` and `runtime_meta.ts` still reflect the split layout.

- [ ] **Step 3: Merge the fine struct implementation into `runtime_struct.ts` and keep `runtime_meta.ts` thin**

Merge this mapping without changing semantics:
- `runtime_struct_core`, `runtime_struct_meta`, `runtime_struct_codec`, `runtime_struct_read`, `runtime_struct_write`, `runtime_struct_field`, `runtime_struct_value_field`, `runtime_struct_scalar_field`, `runtime_struct_primitive_field`, `runtime_struct_blob_field`, `runtime_struct_ref_field`, `runtime_struct_list_field`, `runtime_struct_bitmap_field`, `runtime_struct_bool_list_field`, `runtime_struct_default_list_field`, `runtime_struct_item_field` -> `internal/tpl_ts_runtime_struct.ts.txt`
- keep `internal/tpl_ts_runtime_enum.ts.txt` as the real enum implementation file
- keep `internal/tpl_ts_runtime_meta.ts.txt` as a re-export-only barrel of `runtime_enum.ts` and `runtime_struct.ts`

- [ ] **Step 4: Run focused Go and Bun tests to get GREEN**

Run:
- `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`
- `bun test fixtures/ts/sb/struct_smoke.test.ts fixtures/ts/sb/enum_smoke.test.ts`

Expected: PASS.

- [ ] **Step 5: Commit the enum/struct consolidation**

```bash
git add internal/tpl_ts_runtime_enum.ts.txt internal/tpl_ts_runtime_struct.ts.txt internal/tpl_ts_runtime_meta.ts.txt internal/tpl_ts_render_test.go fixtures/ts/sb/struct_smoke.test.ts fixtures/ts/sb/enum_smoke.test.ts
git commit -m "refactor: consolidate ts runtime meta domains"
```

### Task 5: Prune embed/output lists and enforce deterministic runtime cleanup

**Files:**
- Modify: `internal/tpl_ts_runtime.ts.txt`
- Modify: `internal/tpl_ts_runtime_core.ts.txt`
- Modify: `internal/tpl_ts_runtime_meta.ts.txt`
- Modify: `internal/tpl_ts_runtime_render.go`
- Modify: `internal/tpl_ts_render.go`
- Modify: `internal/generated_file_writer.go`
- Modify: `internal/tpl_ts_render_test.go`
- Test: `internal/tpl_ts_render_test.go`

- [ ] **Step 1: Write the failing cleanup-contract assertions**

Add RED assertions that a real generation pass removes stale split files from an existing `sb` directory. The test should set up at least one stale file before calling `Generate`, for example:

```go
stale := filepath.Join(tmpDir, "sb", "runtime_text_io.ts")
if err := os.WriteFile(stale, []byte("stale"), 0o644); err != nil { t.Fatal(err) }
if err := g.Generate(schema); err != nil { t.Fatal(err) }
assertFileMissing(t, stale)
```

- [ ] **Step 2: Run the focused generator test and verify RED**

Run: `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`

Expected: FAIL because the current generator does not compute an allowed runtime set or delete stale split files deterministically.

- [ ] **Step 3: Implement the minimal cleanup contract**

Implement one of these minimal designs:

```go
func removeLegacyTSRuntimeFiles(targetDir string, allowed map[string]struct{}) error
```

Requirements:
- compute the exact allowed runtime filename set inside `internal/tpl_ts_render.go`
- remove any existing generated `runtime_*.ts` file not in that set before/while writing outputs
- never treat `fixtures/ts/sb/runtime.test.ts` or any `*.test.ts` file as a cleanup target
- shrink `internal/tpl_ts_runtime_render.go` so it only embeds the surviving coarse templates and thin barrel templates
- keep `internal/tpl_ts_runtime.ts.txt`, `internal/tpl_ts_runtime_core.ts.txt`, and `internal/tpl_ts_runtime_meta.ts.txt` as re-export-only templates

- [ ] **Step 4: Re-run the focused generator regression and get GREEN**

Run: `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`

Expected: PASS, including stale-file cleanup and thin-barrel assertions.

- [ ] **Step 5: Commit the generator cleanup contract**

```bash
git add internal/tpl_ts_runtime.ts.txt internal/tpl_ts_runtime_core.ts.txt internal/tpl_ts_runtime_meta.ts.txt internal/tpl_ts_runtime_render.go internal/tpl_ts_render.go internal/generated_file_writer.go internal/tpl_ts_render_test.go
git commit -m "refactor: clean up consolidated ts runtime outputs"
```

### Task 6: Regenerate fixtures and verify end-to-end behavior

**Files:**
- Modify: `fixtures/ts/sb/*`
- Modify: `fixtures/go/sb/*` (only if generator output changes require it)
- Test: `fixtures/go/sb/cross_consistency_test.go`
- Test: `fixtures/go/sb/runtime_test.go`
- Test: `fixtures/go/sb/rpc_transport_test.go`
- Test: `fixtures/ts/sb/runtime.test.ts`
- Test: `fixtures/ts/sb/enum_smoke.test.ts`
- Test: `fixtures/ts/sb/struct_smoke.test.ts`
- Test: `fixtures/ts/sb/rpc_smoke.test.ts`
- Test: `fixtures/ts/sb/rpc_fetch.test.ts`

- [ ] **Step 1: Run a real generator pass against the fixture directories**

Run: `go run . -go ./fixtures/go -ts ./fixtures/ts -tag bson,json ./fixtures/schema.sb`

Expected:
- `fixtures/ts/sb/` keeps `type.ts` as the compatibility barrel
- `fixtures/ts/sb/` only contains the coarse runtime files plus `runtime_core.ts` and `runtime_meta.ts`
- old fine-grained runtime files are gone without manual deletion

- [ ] **Step 2: Run focused semantic regressions first**

Run:
- `go test ./fixtures/go/sb -run 'TestCrossLanguageWireConsistency|TestCrossLanguageWireConsistencyRandom|TestCrossLanguageWireRejectsMalformedInputs|TestProtocolGoldenVectors|TestHeaderReadWriteRoundTrip|TestHeaderPaddingAndRangeValidation' -count=1`
- `go test ./fixtures/go/sb -run 'TestDoClientSlowBodyWithinTimeout|TestDoClientBodyReadTimeout|TestDoClientTruncatedBodyReturnsRespErr|TestDoClientDoesNotRetryWithoutEnableRetries|TestDoClientBodyReadTimeoutThenRetrySuccess|TestDoClientRetriesOnHTTP408|TestDoClientCancelDuringBackoff' -count=1`
- `! rg 'from "\./type"' fixtures/ts/sb/enum.ts fixtures/ts/sb/struct_*.ts fixtures/ts/sb/rpc.ts`
- `rg 'from "\./runtime_(core|meta)"' fixtures/ts/sb/enum.ts fixtures/ts/sb/struct_*.ts fixtures/ts/sb/rpc.ts`
- `bun test fixtures/ts/sb/runtime.test.ts fixtures/ts/sb/enum_smoke.test.ts fixtures/ts/sb/struct_smoke.test.ts fixtures/ts/sb/rpc_smoke.test.ts fixtures/ts/sb/rpc_fetch.test.ts`

Expected:
- both `go test` commands PASS
- the first `rg` command returns no matches
- the second `rg` command returns matches
- the Bun suite PASSes

- [ ] **Step 3: Run full repository verification**

Run:
- `go test ./...`
- `bun test`

Expected: PASS, with no new red tests.

- [ ] **Step 4: Inspect the final generated runtime footprint**

Run:

```bash
rg --files fixtures/ts/sb -g 'runtime_*.ts'
```

Also verify separately that `fixtures/ts/sb/type.ts` still exists as a thin compatibility barrel.
Also verify separately that `fixtures/ts/sb/runtime.test.ts` still exists and was not treated as stale generated output.

Expected `rg` output should only include:
- `fixtures/ts/sb/runtime_base.ts`
- `fixtures/ts/sb/runtime_header.ts`
- `fixtures/ts/sb/runtime_text.ts`
- `fixtures/ts/sb/runtime_bin.ts`
- `fixtures/ts/sb/runtime_list.ts`
- `fixtures/ts/sb/runtime_struct.ts`
- `fixtures/ts/sb/runtime_enum.ts`
- `fixtures/ts/sb/runtime_core.ts`
- `fixtures/ts/sb/runtime_meta.ts`

- [ ] **Step 5: Commit the regenerated fixtures and final consolidation**

```bash
git add fixtures/ts/sb fixtures/go/sb
git commit -m "refactor: consolidate generated ts runtime"
```

## Final Verification Checklist

- [ ] `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`
- [ ] `go run . -go ./fixtures/go -ts ./fixtures/ts -tag bson,json ./fixtures/schema.sb`
- [ ] `go test ./fixtures/go/sb -run 'TestCrossLanguageWireConsistency|TestCrossLanguageWireConsistencyRandom|TestCrossLanguageWireRejectsMalformedInputs|TestProtocolGoldenVectors|TestHeaderReadWriteRoundTrip|TestHeaderPaddingAndRangeValidation' -count=1`
- [ ] `go test ./fixtures/go/sb -run 'TestDoClientSlowBodyWithinTimeout|TestDoClientBodyReadTimeout|TestDoClientTruncatedBodyReturnsRespErr|TestDoClientDoesNotRetryWithoutEnableRetries|TestDoClientBodyReadTimeoutThenRetrySuccess|TestDoClientRetriesOnHTTP408|TestDoClientCancelDuringBackoff' -count=1`
- [ ] `! rg 'from "\./type"' fixtures/ts/sb/enum.ts fixtures/ts/sb/struct_*.ts fixtures/ts/sb/rpc.ts`
- [ ] `rg 'from "\./runtime_(core|meta)"' fixtures/ts/sb/enum.ts fixtures/ts/sb/struct_*.ts fixtures/ts/sb/rpc.ts`
- [ ] `bun test fixtures/ts/sb/runtime.test.ts fixtures/ts/sb/enum_smoke.test.ts fixtures/ts/sb/struct_smoke.test.ts fixtures/ts/sb/rpc_smoke.test.ts fixtures/ts/sb/rpc_fetch.test.ts`
- [ ] `go test ./...`
- [ ] `bun test`
- [ ] `rg --files fixtures/ts/sb -g 'runtime_*.ts'`
- [ ] Verify `type.ts`, `runtime_core.ts`, and `runtime_meta.ts` are thin barrels
- [ ] Verify `runtime.test.ts` still exists and was not cleaned as generated runtime output
- [ ] Verify no hand-deleted fixture files were required to make the build green
