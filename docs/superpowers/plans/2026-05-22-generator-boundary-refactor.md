# Generator Boundary Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [x]`) syntax for tracking.

**Goal:** Split the large Go and TypeScript generator render files into smaller responsibility-focused files without changing generated output or protocol behavior.

**Architecture:** Keep all code in `package internal`. Move existing functions by responsibility only: emit/file cleanup, API/RPC rendering, schema/enum/struct rendering, expression helpers, parser phases, and documentation sections. Preserve current function signatures and call order so this is a mechanical boundary refactor.

**Tech Stack:** Go 1.26, Bun 1.3, existing `go test`, `go run`, and `bun test` verification.

---

## File Structure

- Modify: `internal/tpl_go_render.go`
  - Keep only schema-level render entry if useful, or leave as a thin compatibility file after moving functions.
- Create: `internal/tpl_go_emit.go`
  - Own `GoGenerator.Generate`, stale Go API cleanup, fingerprint-managed removal, and legacy Go layout cleanup.
- Create: `internal/tpl_go_render_api.go`
  - Own `renderAPIStub`, `renderAPIFile`, `renderAPIRoutesAndHelpers`.
- Create: `internal/tpl_go_render_api_handler.go`
  - Own single API handler rendering.
- Create: `internal/tpl_go_render_api_middleware.go`
  - Own API middleware rendering.
- Create: `internal/tpl_go_render_api_routes.go`
  - Own route registration rendering.
- Create: `internal/tpl_go_render_api_http.go`
  - Own HTTP helper rendering.
- Create: `internal/tpl_go_render_rpc.go`
  - Own `renderRPCFile`.
- Create: `internal/tpl_go_render_struct.go`
  - Own `renderStruct` orchestration.
- Create: `internal/tpl_go_render_struct_meta.go`
  - Own struct type/meta body render section.
- Create: `internal/tpl_go_render_struct_size.go`
  - Own struct size body render section.
- Create: `internal/tpl_go_render_struct_get.go`
  - Own struct get body render section.
- Create: `internal/tpl_go_render_struct_set.go`
  - Own struct set body render section.
- Create: `internal/tpl_go_render_struct_read.go`
  - Own struct read/equality body render section.
- Create: `internal/tpl_go_render_shared.go`
  - Own `renderSharedHelpers`.
- Create: `internal/tpl_go_render_enum.go`
  - Own `renderEnum`.
- Create: `internal/tpl_go_render_struct_list.go`
  - Own `renderStructListBodyHelpers`.
- Create: `internal/tpl_go_names.go`
  - Own Go enum/struct naming helpers and header width names.
- Create: `internal/tpl_go_type_expr.go`
  - Own tag width, default checks, state init, and primitive getter/setter helpers.
- Create: `internal/tpl_go_list_expr.go`
  - Own list size/get/set, bitmap callbacks, field meta, list validation, and equality helpers.
- Create: `internal/tpl_go_rpc_expr.go`
  - Own API/RPC parameter, return, and request/response wrapper helpers.
- Create: `internal/tpl_go_direct_expr.go`
  - Own direct read/write and direct list expression helpers.
- Modify: `internal/tpl_ts_render.go`
  - Keep only schema-level render entry if useful, or leave as a thin compatibility file after moving functions.
- Create: `internal/tpl_ts_emit.go`
  - Own `TsGenerator.Generate`, stale TS cleanup, runtime file cleanup.
- Create: `internal/tpl_ts_render_rpc.go`
  - Own `renderRPCFile` and `renderSmokeTest`.
- Create: `internal/tpl_ts_render_struct.go`
  - Own `renderStructFile`, `renderStructListBodyHelpers`, and struct file helpers.
- Create: `internal/tpl_ts_type_expr.go`
  - Own TS type, default, tag width, and primitive expression helpers.
- Create: `internal/tpl_ts_imports.go`
  - Own generated import aggregation and sorting helpers.
- Create: `internal/tpl_ts_struct_expr.go`
  - Own struct field, meta, list, bitmap, and equality expression helpers.
- Create: `internal/tpl_ts_rpc_expr.go`
  - Own RPC parameter, return, and request/response wrapper helpers.
- Create: `internal/tpl_ts_direct_expr.go`
  - Own direct read/write, direct list, and read error suffix helpers.
- Test: existing `internal/tpl_go_render_test.go`, `internal/tpl_ts_render_test.go`, generated fixture tests, and cross-consistency tests.
- Modify: `internal/parser_schema.go`
  - Own only top-level schema dispatch.
- Create: `internal/parser_core.go`
  - Own `Parser` state, token advancement, newline skipping, and token expectations.
- Create: `internal/parser_struct.go`
  - Own struct, field, and struct field line-ending parsing.
- Create: `internal/parser_enum.go`
  - Own enum and enum member parsing.
- Create: `internal/parser_api.go`
  - Own API name, argument, and return parsing.
- Create: `internal/parser_type.go`
  - Own scalar/list/nil type parsing.
- Modify: `internal/tpl_doc_render.go`
  - Keep doc generation entry and section orchestration.
- Create: `internal/tpl_doc_sections.go`
  - Own API list, RPC status, usage demos, and types sections.
- Create: `internal/tpl_doc_expr.go`
  - Own doc argument/return formatting and Go/TS example value helpers.
- Test: existing parser, resolver, doc render, full generator, and generated fixture tests.
- Modify: `internal/tpl_go_runtime_render.go`
  - Embed and concatenate Go runtime template fragments.
- Delete: `internal/tpl_go_runtime.go.txt`
  - Replaced by responsibility-focused fragments.
- Create: `internal/tpl_go_runtime_core.go.txt`
- Create: `internal/tpl_go_runtime_enum.go.txt`
- Create: `internal/tpl_go_runtime_struct.go.txt`
- Create: `internal/tpl_go_runtime_bits.go.txt`
- Create: `internal/tpl_go_runtime_base.go.txt`
- Create: `internal/tpl_go_runtime_text.go.txt`
- Create: `internal/tpl_go_runtime_bool_list.go.txt`
- Create: `internal/tpl_go_runtime_bin.go.txt`
- Create: `internal/tpl_go_runtime_list.go.txt`
- Create: `internal/tpl_go_runtime_text_list.go.txt`
- Modify: `internal/semantic_resolver.go`
  - Keep resolver entry and main orchestration.
- Create: `internal/semantic_defs.go`
  - Own definition collection.
- Create: `internal/semantic_enum.go`
  - Own enum resolution.
- Create: `internal/semantic_struct.go`
  - Own struct expansion and duplicate field checks.
- Create: `internal/semantic_api.go`
  - Own API resolution.
- Create: `internal/semantic_type.go`
  - Own type resolution.

## Task 1: Baseline Guard

**Files:**
- Read: `internal/tpl_go_render.go`
- Read: `internal/tpl_ts_render.go`
- Test: full repository test commands

- [x] **Step 1: Confirm dirty baseline**

Run:

```bash
git status --short
```

Expected: the four existing generator cleanup files may be modified, and this plan/spec may be untracked. Do not revert them.

- [x] **Step 2: Run baseline Go tests**

Run:

```bash
go test ./...
```

Expected: exit 0.

- [x] **Step 3: Run baseline TS tests**

Run:

```bash
bun test fixtures/ts/sb
```

Expected: exit 0 with 85 passing tests.

## Task 2: Move Go Emit And Cleanup Logic

**Files:**
- Create: `internal/tpl_go_emit.go`
- Modify: `internal/tpl_go_render.go`

- [x] **Step 1: Create `internal/tpl_go_emit.go`**

Move these exact functions from `internal/tpl_go_render.go` into the new file:

```go
func (g *GoGenerator) Generate(schema *TplSchema) error
func removeStaleGoAPIFiles(targetDir string, apis []TplApi) error
func removeFingerprintManagedFile(path string) error
func removeLegacyGoLayoutFiles(targetDir string) error
```

Use this import block in `internal/tpl_go_emit.go`:

```go
package internal

import (
	"os"
	"path/filepath"
)
```

- [x] **Step 2: Remove moved imports from `internal/tpl_go_render.go`**

After moving the emit functions, remove imports no longer used by `internal/tpl_go_render.go`. The remaining import block should include only packages still referenced in the file, likely:

```go
import "fmt"
```

- [x] **Step 3: Format and test Go generator package**

Run:

```bash
gofmt -w internal/tpl_go_emit.go internal/tpl_go_render.go
go test ./internal -run 'TestGoGenerator|TestGeneratedDir' -count=1
```

Expected: exit 0.

## Task 3: Move Go API And RPC Render Logic

**Files:**
- Create: `internal/tpl_go_render_api.go`
- Create: `internal/tpl_go_render_rpc.go`
- Modify: `internal/tpl_go_render.go`

- [x] **Step 1: Move API render functions**

Move these functions from `internal/tpl_go_render.go` to `internal/tpl_go_render_api.go`:

```go
func (g *GoGenerator) renderAPIStub(api TplApi) []byte
func (g *GoGenerator) renderAPIFile(apis []TplApi) string
func (g *GoGenerator) renderAPIRoutesAndHelpers(w *sourceWriter, apis []TplApi)
```

Use:

```go
package internal

import "fmt"
```

Keep the function bodies unchanged.

- [x] **Step 2: Move RPC render function**

Move this function from `internal/tpl_go_render.go` to `internal/tpl_go_render_rpc.go`:

```go
func (g *GoGenerator) renderRPCFile(apis []TplApi) string
```

Use:

```go
package internal

import "fmt"
```

Keep the function body unchanged.

- [x] **Step 3: Format and run API/RPC render tests**

Run:

```bash
gofmt -w internal/tpl_go_render.go internal/tpl_go_render_api.go internal/tpl_go_render_rpc.go
go test ./internal -run TestGoGeneratorWritesSchemaFile -count=1
```

Expected: exit 0.

## Task 4: Move Go Struct And Expression Helpers

**Files:**
- Create: `internal/tpl_go_render_struct.go`
- Create: `internal/tpl_go_names.go`
- Create: `internal/tpl_go_type_expr.go`
- Create: `internal/tpl_go_list_expr.go`
- Create: `internal/tpl_go_rpc_expr.go`
- Create: `internal/tpl_go_direct_expr.go`
- Modify: `internal/tpl_go_render.go`

- [x] **Step 1: Move struct render functions**

Move these functions from `internal/tpl_go_render.go` to `internal/tpl_go_render_struct.go`:

```go
func (g *GoGenerator) renderSharedHelpers(w *sourceWriter)
func (g *GoGenerator) renderEnum(w *sourceWriter, enum TplEnum)
func (g *GoGenerator) renderStruct(w *sourceWriter, st TplStruct)
func (g *GoGenerator) renderStructListBodyHelpers(w *sourceWriter, st TplStruct)
```

Use:

```go
package internal
```

No import block should be needed unless the moved code references a package directly.

- [x] **Step 2: Move expression helper functions**

Move all helper functions from `enumDefaultName` through `directListSetExpr` from `internal/tpl_go_render.go` into these responsibility files:

- `internal/tpl_go_names.go`: `enumDefaultName` through `isAPIWrapper`.
- `internal/tpl_go_type_expr.go`: `headerBits` through `primitiveSetter`.
- `internal/tpl_go_list_expr.go`: `listSizeExpr` through `goListEqExpr`.
- `internal/tpl_go_rpc_expr.go`: `logicType` through `callArgNames`.
- `internal/tpl_go_direct_expr.go`: `directRead` through `directListSetExpr`.

- [x] **Step 3: Confirm `tpl_go_render.go` is schema-only**

After the moves, `internal/tpl_go_render.go` should keep:

```go
func (g *GoGenerator) renderSchemaFile(schema *TplSchema) string
```

It may keep `package internal` and a minimal import block. Do not change the body except for gofmt.

- [x] **Step 4: Format and run Go tests**

Run:

```bash
gofmt -w internal/tpl_go_render.go internal/tpl_go_render_struct.go internal/tpl_go_names.go internal/tpl_go_type_expr.go internal/tpl_go_list_expr.go internal/tpl_go_rpc_expr.go internal/tpl_go_direct_expr.go
go test ./internal -run TestGoGeneratorWritesSchemaFile -count=1
```

Expected: exit 0.

## Task 5: Move TS Emit And Runtime Cleanup Logic

**Files:**
- Create: `internal/tpl_ts_emit.go`
- Modify: `internal/tpl_ts_render.go`

- [x] **Step 1: Move TS emit and cleanup functions**

Move these functions from `internal/tpl_ts_render.go` into `internal/tpl_ts_emit.go`:

```go
func (g *TsGenerator) Generate(schema *TplSchema) error
func (g *TsGenerator) removeLegacyAPIWrapperStructFiles(targetDir string) error
func removeStaleTSGeneratedFiles(targetDir string, schema *TplSchema) error
func removeStaleTSStructFiles(targetDir string, active map[string]struct{}) error
func removeTSFiles(targetDir string, names ...string) error
func (g *TsGenerator) tsGeneratedFiles(schema *TplSchema) []generatedFile
func tsAllowedRuntimeFiles(files []generatedFile) map[string]struct{}
func removeLegacyTSRuntimeFiles(targetDir string, allowed map[string]struct{}) error
```

Use:

```go
package internal

import (
	"os"
	"path/filepath"
	"strings"
)
```

- [x] **Step 2: Remove moved imports from `internal/tpl_ts_render.go`**

After the move, remove `os`, `path/filepath`, and any other unused imports from `internal/tpl_ts_render.go`.

- [x] **Step 3: Format and run TS generator tests**

Run:

```bash
gofmt -w internal/tpl_ts_emit.go internal/tpl_ts_render.go
go test ./internal -run 'TestTsGenerator|TestGeneratedDir' -count=1
```

Expected: exit 0.

## Task 6: Move TS RPC, Struct, And Expression Helpers

**Files:**
- Create: `internal/tpl_ts_render_rpc.go`
- Create: `internal/tpl_ts_render_struct.go`
- Create: `internal/tpl_ts_type_expr.go`
- Create: `internal/tpl_ts_imports.go`
- Create: `internal/tpl_ts_struct_expr.go`
- Create: `internal/tpl_ts_rpc_expr.go`
- Create: `internal/tpl_ts_direct_expr.go`
- Modify: `internal/tpl_ts_render.go`

- [x] **Step 1: Move RPC functions**

Move these functions from `internal/tpl_ts_render.go` to `internal/tpl_ts_render_rpc.go`:

```go
func (g *TsGenerator) renderRPCFile(apis []TplApi) string
func (g *TsGenerator) renderSmokeTest(apis []TplApi) string
```

Use:

```go
package internal

import (
	"fmt"
	"strings"
)
```

- [x] **Step 2: Move struct render functions**

Move these functions from `internal/tpl_ts_render.go` to `internal/tpl_ts_render_struct.go`:

```go
func (g *TsGenerator) renderStructFile(st TplStruct) string
func (g *TsGenerator) renderStructListBodyHelpers(w *sourceWriter, st TplStruct)
```

Use:

```go
package internal

import "fmt"
```

- [x] **Step 3: Move expression and import helpers**

Move these functions from `fieldType` through `readErrSuffix` from `internal/tpl_ts_render.go` into these responsibility files:

- `internal/tpl_ts_type_expr.go`: `fieldType` through `primitiveEq`.
- `internal/tpl_ts_imports.go`: `tsStructImportLines`, `tsRpcImportLines`, and `joinImportNames`.
- `internal/tpl_ts_struct_expr.go`: `structFieldType` through `eqExpr`.
- `internal/tpl_ts_rpc_expr.go`: `rpcArgList` through `apiRespTypeName`.
- `internal/tpl_ts_direct_expr.go`: `tsDirectRead` through `readErrSuffix`.

- [x] **Step 4: Leave schema render in `tpl_ts_render.go`**

After the moves, `internal/tpl_ts_render.go` should keep:

```go
func (g *TsGenerator) renderEnumFile(enums []TplEnum) string
func (g *TsGenerator) renderIndexFile(files []string) string
func (g *TsGenerator) structHeaderWidthsName(st TplStruct) string
func (g *TsGenerator) structHeaderWidthsLiteral(st TplStruct) string
func (g *TsGenerator) renderEnumSmokeTest(enums []TplEnum) string
```

Keep function bodies unchanged.

- [x] **Step 5: Format and run TS generator tests**

Run:

```bash
gofmt -w internal/tpl_ts_render.go internal/tpl_ts_render_rpc.go internal/tpl_ts_render_struct.go internal/tpl_ts_type_expr.go internal/tpl_ts_imports.go internal/tpl_ts_struct_expr.go internal/tpl_ts_rpc_expr.go internal/tpl_ts_direct_expr.go
go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1
```

Expected: exit 0.

## Task 7: Full Regeneration And Verification

**Files:**
- Verify: all modified `internal/tpl_*` files
- Verify: `fixtures/go/sb/*`
- Verify: `fixtures/ts/sb/*`

- [x] **Step 1: Run full Go tests**

Run:

```bash
go test ./...
```

Expected: exit 0.

- [x] **Step 2: Regenerate fixtures**

Run:

```bash
go run . -go ./fixtures/go -ts ./fixtures/ts -tag bson,json ./fixtures/schema.sb
```

Expected: `代码生成成功。`

- [x] **Step 3: Check fixture diff**

Run:

```bash
git diff -- fixtures/go/sb fixtures/ts/sb
```

Expected: no output. If output appears, stop and inspect before continuing.

- [x] **Step 4: Run Bun tests**

Run:

```bash
bun test fixtures/ts/sb
```

Expected: exit 0 with 85 passing tests.

- [x] **Step 5: Run whitespace check**

Run:

```bash
git diff --check -- .
```

Expected: no output, exit 0.

- [x] **Step 6: Review final file sizes**

Run:

```bash
wc -l internal/tpl_go_*.go internal/tpl_ts_*.go | sort -n
```

Expected: no single moved generator file should remain near the old 1000+ line size. This is an informational check; do not chase arbitrary line-count targets by changing behavior.

## Task 8: Split Parser And Documentation Boundaries

**Files:**
- Create: `internal/parser_core.go`
- Modify: `internal/parser_schema.go`
- Create: `internal/parser_struct.go`
- Create: `internal/parser_enum.go`
- Create: `internal/parser_api.go`
- Create: `internal/parser_type.go`
- Modify: `internal/tpl_doc_render.go`
- Create: `internal/tpl_doc_sections.go`
- Create: `internal/tpl_doc_expr.go`

- [x] **Step 1: Split parser by parse phase**

Move existing parser functions without changing bodies:

- `parser_core.go`: `Parser`, `NewParser`, `nextToken`, `skipNewLines`, `expectCurrent`.
- `parser_schema.go`: `ParseSchema`.
- `parser_struct.go`: `parseStruct`, `ensureStructFieldLineEnd`, `parseField`.
- `parser_enum.go`: `parseEnum`, `parseEnumMember`.
- `parser_api.go`: `parseAPI`.
- `parser_type.go`: `parseType`.

- [x] **Step 2: Verify parser behavior**

Run:

```bash
gofmt -w internal/parser_core.go internal/parser_schema.go internal/parser_struct.go internal/parser_enum.go internal/parser_api.go internal/parser_type.go
go test ./internal -run 'TestParse|TestResolve|TestGrammarDocContract' -count=1
```

Expected: exit 0.

## Task 9: Split Go Struct Render Body

**Files:**
- Modify: `internal/tpl_go_render_struct.go`
- Create: `internal/tpl_go_render_struct_meta.go`
- Create: `internal/tpl_go_render_struct_size.go`
- Create: `internal/tpl_go_render_struct_get.go`
- Create: `internal/tpl_go_render_struct_set.go`
- Create: `internal/tpl_go_render_struct_read.go`
- Create: `internal/tpl_go_render_shared.go`
- Create: `internal/tpl_go_render_enum.go`
- Create: `internal/tpl_go_render_struct_list.go`

- [x] **Step 1: Move shared, enum, and list renderers**

Move these functions out of `internal/tpl_go_render_struct.go`:

- `renderSharedHelpers` to `internal/tpl_go_render_shared.go`.
- `renderEnum` to `internal/tpl_go_render_enum.go`.
- `renderStructListBodyHelpers` to `internal/tpl_go_render_struct_list.go`.

- [x] **Step 2: Split struct body rendering**

Keep `renderStruct` as the orchestration function in `internal/tpl_go_render_struct.go`. Move its body sections to separate files:

- `renderStructTypeAndMeta`
- `renderStructSize`
- `renderStructGet`
- `renderStructSet`
- `renderStructReadAndEq`

- [x] **Step 3: Verify Go struct rendering**

Run:

```bash
gofmt -w internal/tpl_go_render_shared.go internal/tpl_go_render_enum.go internal/tpl_go_render_struct.go internal/tpl_go_render_struct_meta.go internal/tpl_go_render_struct_size.go internal/tpl_go_render_struct_get.go internal/tpl_go_render_struct_set.go internal/tpl_go_render_struct_read.go internal/tpl_go_render_struct_list.go
go test ./internal -run TestGoGeneratorWritesSchemaFile -count=1
```

Expected: exit 0.

## Task 10: Split Resolver And Go API Render Boundaries

**Files:**
- Modify: `internal/semantic_resolver.go`
- Create: `internal/semantic_defs.go`
- Create: `internal/semantic_enum.go`
- Create: `internal/semantic_struct.go`
- Create: `internal/semantic_api.go`
- Create: `internal/semantic_type.go`
- Modify: `internal/tpl_go_render_api.go`
- Create: `internal/tpl_go_render_api_handler.go`
- Create: `internal/tpl_go_render_api_middleware.go`
- Create: `internal/tpl_go_render_api_routes.go`
- Create: `internal/tpl_go_render_api_http.go`

- [x] **Step 1: Split resolver by semantic phase**

Keep `Resolve`, `resolver`, and `resolve` in `internal/semantic_resolver.go`. Move definition, enum, struct, API, and type resolution helpers into the matching `semantic_*.go` files.

- [x] **Step 2: Verify resolver behavior**

Run:

```bash
gofmt -w internal/semantic_resolver.go internal/semantic_defs.go internal/semantic_enum.go internal/semantic_struct.go internal/semantic_api.go internal/semantic_type.go
go test ./internal -run TestResolve -count=1
```

Expected: exit 0.

## Task 11: Split Go Runtime Template

**Files:**
- Modify: `internal/tpl_go_runtime_render.go`
- Delete: `internal/tpl_go_runtime.go.txt`
- Create: `internal/tpl_go_runtime_core.go.txt`
- Create: `internal/tpl_go_runtime_enum.go.txt`
- Create: `internal/tpl_go_runtime_struct.go.txt`
- Create: `internal/tpl_go_runtime_bits.go.txt`
- Create: `internal/tpl_go_runtime_base.go.txt`
- Create: `internal/tpl_go_runtime_text.go.txt`
- Create: `internal/tpl_go_runtime_bool_list.go.txt`
- Create: `internal/tpl_go_runtime_bin.go.txt`
- Create: `internal/tpl_go_runtime_list.go.txt`
- Create: `internal/tpl_go_runtime_text_list.go.txt`

- [x] **Step 1: Split runtime template by protocol domain**

Split the previous single Go runtime template into core/header, enum, struct meta, bits, base numeric, text, bool list, bin, generic list, and text list fragments.

- [x] **Step 2: Concatenate fragments in render order**

Update `renderGoRuntimeSource` to embed every fragment and concatenate them in the original source order.

- [x] **Step 3: Verify generated runtime is unchanged**

Run:

```bash
gofmt -w internal/tpl_go_runtime_render.go
go test ./internal -run TestGoGeneratorWritesSchemaFile -count=1
```

Expected: exit 0.

- [x] **Step 3: Split Go API render sections**

Keep `renderAPIStub`, `renderAPIFile`, and `renderAPIRoutesAndHelpers` orchestration in `internal/tpl_go_render_api.go`. Move handler, middleware, route registration, and HTTP helper rendering into separate files.

- [x] **Step 4: Verify Go API rendering**

Run:

```bash
gofmt -w internal/tpl_go_render_api.go internal/tpl_go_render_api_handler.go internal/tpl_go_render_api_middleware.go internal/tpl_go_render_api_routes.go internal/tpl_go_render_api_http.go
go test ./internal -run TestGoGeneratorWritesSchemaFile -count=1
```

Expected: exit 0.

- [x] **Step 3: Split doc rendering by section**

Keep `generateDoc` and `renderDoc` in `internal/tpl_doc_render.go`. Move section rendering to `internal/tpl_doc_sections.go` and argument/example helpers to `internal/tpl_doc_expr.go`.

- [x] **Step 4: Verify doc rendering**

Run:

```bash
gofmt -w internal/tpl_doc_render.go internal/tpl_doc_sections.go internal/tpl_doc_expr.go
go test ./internal -run TestRenderDocUsesSbPublicImports -count=1
```

Expected: exit 0.
