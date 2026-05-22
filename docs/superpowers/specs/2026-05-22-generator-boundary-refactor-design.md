# Generator Boundary Refactor Design

## 1. 目标

第一阶段优化生成器实现边界; 第二阶段继续收敛 parser 和文档渲染边界。两阶段都不改变 `.sb` 语法、wire format、生成代码语义或 public API。

当前核心问题是 `internal/tpl_go_render.go` 和 `internal/tpl_ts_render.go` 文件过大, 且混合了生成入口、文件同步、schema/enum/struct/RPC 渲染、类型表达式和清理策略。重构目标是把这些职责按生命周期拆开, 保持同一个 `internal` 包, 降低一次变更的风险。

## 2. 非目标

- 不拆 Go package 到 `parse/resolve/gen/write` 子包。
- 不重写 resolver、runtime 模板或 fixture 结构。
- 不改变 parser 接受/拒绝的语法、错误文案和注释绑定规则。
- 不改变文档生成内容, 只拆分章节渲染边界。
- 不调整协议规范、字段顺序、默认值、RPC 状态码和错误语义。
- 不修改生成产物的文本输出, 除非 gofmt 或现有生成器必然产生完全一致结果。
- 不回退或覆盖当前工作区已有的生成器清理修复。

## 3. 设计

### 3.1 边界

继续保留 `GoGenerator` 和 `TsGenerator` 作为配置承载类型, 但按职责拆分当前大文件。拆分只在文件层面完成, 所有代码仍属于 `package internal`, 这样测试和调用路径不需要跨包迁移。

Go 侧建议拆分:

- `tpl_go_emit.go`: `GoGenerator.Generate`、文件写入和旧产物清理。
- `tpl_go_render_schema.go`: schema 文件、shared helper、enum 和 struct 入口。
- `tpl_go_render_api.go`: API handler、路由和服务端 helper。
- `tpl_go_render_api_handler.go`: 单个 API handler 渲染。
- `tpl_go_render_api_middleware.go`: Go API middleware 渲染。
- `tpl_go_render_api_routes.go`: Go API route registration 渲染。
- `tpl_go_render_api_http.go`: Go API HTTP helper 渲染。
- `tpl_go_render_rpc.go`: RPC client 渲染。
- `tpl_go_render_shared.go`: Go schema 共享 helper 渲染。
- `tpl_go_render_enum.go`: enum 和 enum list helper 渲染。
- `tpl_go_render_struct.go`: struct 渲染调度。
- `tpl_go_render_struct_meta.go`: struct type/meta 渲染。
- `tpl_go_render_struct_size.go`: struct size 渲染。
- `tpl_go_render_struct_get.go`: struct get 渲染。
- `tpl_go_render_struct_set.go`: struct set 渲染。
- `tpl_go_render_struct_read.go`: struct read/eq 渲染。
- `tpl_go_render_struct_list.go`: struct list helper 渲染。
- `tpl_go_names.go`: enum/struct 命名和 header width 名称。
- `tpl_go_type_expr.go`: tag width、默认值判断、状态初始化和 primitive getter/setter。
- `tpl_go_list_expr.go`: list size/get/set、bitmap callback、struct field meta 和 list 校验/比较表达式。
- `tpl_go_rpc_expr.go`: API/RPC 参数、返回值和 request/response wrapper helper。
- `tpl_go_direct_expr.go`: direct read/write 和 direct list 表达式。

TS 侧建议拆分:

- `tpl_ts_emit.go`: `TsGenerator.Generate`、文件写入和旧产物清理。
- `tpl_ts_render_schema.go`: enum、struct、index 和 smoke test 入口。
- `tpl_ts_render_rpc.go`: RPC client 与 RPC smoke test。
- `tpl_ts_render_struct.go`: struct 编解码和 list helper。
- `tpl_ts_type_expr.go`: TS 类型、默认值、tag width、primitive getter/setter/eq。
- `tpl_ts_imports.go`: struct/RPC import 汇总和排序。
- `tpl_ts_struct_expr.go`: struct 字段、meta、list、bitmap 和 eq 表达式。
- `tpl_ts_rpc_expr.go`: RPC 参数、返回值和 request/response wrapper helper。
- `tpl_ts_direct_expr.go`: direct read/write、direct list 表达式和 read error suffix。

Parser 侧拆分:

- `parser_core.go`: `Parser` 状态、token 推进、换行跳过和 current token 校验。
- `parser_schema.go`: schema 顶层分发和 pending note 绑定。
- `parser_struct.go`: struct、field 和字段行尾约束。
- `parser_enum.go`: enum 和 enum member。
- `parser_api.go`: API 名称、参数和返回类型链路。
- `parser_type.go`: 标量/list/nil 类型解析。

文档渲染侧拆分:

- `tpl_doc_render.go`: 文档生成入口和章节调度。
- `tpl_doc_sections.go`: API list、RPC 状态、用法示例和类型章节。
- `tpl_doc_expr.go`: 参数/返回值和 Go/TS 示例值表达式。

Go runtime 模板侧拆分:

- `tpl_go_runtime_core.go.txt`: package/import、state 常量和 header helper。
- `tpl_go_runtime_enum.go.txt`: enum meta 与 enum list helper。
- `tpl_go_runtime_struct.go.txt`: struct meta helper。
- `tpl_go_runtime_bits.go.txt`: bit reader 和 padding 校验。
- `tpl_go_runtime_base.go.txt`: 数值基础类型编解码。
- `tpl_go_runtime_text.go.txt`: text/bin 状态和 text 编解码。
- `tpl_go_runtime_bool_list.go.txt`: bool list 编解码。
- `tpl_go_runtime_bin.go.txt`: bin 和 bin list 编解码。
- `tpl_go_runtime_list.go.txt`: 通用 list/default/zero/pointer list helper。
- `tpl_go_runtime_text_list.go.txt`: text list 编解码。

Resolver 侧拆分:

- `semantic_resolver.go`: resolver 入口和主调度。
- `semantic_defs.go`: 顶层定义收集和重名检查。
- `semantic_enum.go`: enum 成员 ID 分配与校验。
- `semantic_struct.go`: struct 展开、嵌入和字段重名校验。
- `semantic_api.go`: API 参数、返回值和重名校验。
- `semantic_type.go`: base/struct/enum/nil 类型解析。

### 3.2 风格

- 保持扁平控制流, 先处理错误和空分支, 主干直行。
- 无状态逻辑优先做纯 helper, 输入通过参数显式传入。
- 不引入 builder、class 式封装或链式 API。
- 不新增全局状态。
- 不把“拆文件”扩大成行为重写。

## 4. 数据流

主链路保持不变:

```text
main.go -> parseAndResolve -> internal.Generate -> GoGenerator.Generate / TsGenerator.Generate -> generated files
```

重构后只是把 `GoGenerator.Generate` 和 `TsGenerator.Generate` 调用到的 helper 分散到更小的文件中, 不改变参数、返回值或写盘顺序。

## 5. 错误处理

错误处理保持当前显式返回 `error` 的方式。文件清理和写入仍按现有策略处理:

- Go 手写 API stub 只在 fingerprint 可验证时删除。
- TS 不删除 `*.test.ts` 中的非生成测试文件。
- 文件系统错误继续原样返回上层, 不吞错。

## 6. 测试

实施后必须通过:

```bash
go test ./...
go run . -go ./fixtures/go -ts ./fixtures/ts -tag bson,json ./fixtures/schema.sb
bun test fixtures/ts/sb
git diff --check -- .
```

还要检查 fixture 重生成没有非预期 diff。若拆文件后生成产物出现文本变化, 必须先定位原因, 不能把生成物变化当作“顺手更新”合并。

## 7. 回滚

本阶段是文件边界重排, 回滚方式是恢复新增/移动的生成器源码文件。由于不改变协议和生成语义, 回滚不需要数据迁移或兼容策略。
