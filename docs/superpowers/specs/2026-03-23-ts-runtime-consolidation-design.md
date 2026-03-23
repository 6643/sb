# TS Runtime 收口归并设计

## 1. 结论
- 本轮不再继续细拆 TypeScript runtime.
- 本轮目标改为把当前已拆散的 runtime 重新收口到粗粒度文件集.
- 最终生成结果以粗文件为主, 保留 `type.ts` 作为超薄兼容 barrel, 允许少量中间层文件存在.

## 2. 边界
- 输入:
  - 当前 TS generator 仍会写出大量细粒度 `runtime_*.ts` 文件.
  - 当前 `type.ts` 已经降成轻量 re-export 文件.
  - 当前 `enum.ts` / `struct_*.ts` / `rpc.ts` 已经从 `./type` 切到较窄的 runtime 入口.
- 输出:
  - TS 生成结果收口到粗粒度 runtime 文件集.
  - 细粒度 runtime 实现文件不再作为最终生成产物落盘.
  - `enum.ts` / `struct_*.ts` / `rpc.ts` 继续依赖粗粒度 runtime 入口.
- 约束:
  - 不改变 wire 规则.
  - 不改变 enum/struct/rpc 对外 API 语义.
  - 不引入新依赖.
  - 本轮优先处理 TS, 不强制 Go 与 TS 完全同构.
- 非目标:
  - 不继续做更细粒度拆分.
  - 不把所有实现重新并回 `type.ts`.
  - 不在本轮处理 PR 或发布流程.

## 3. 依据
- 事实:
  - `internal/tpl_ts_render.go` 当前会显式写出几十个 `runtime_*.ts` 文件.
  - `fixtures/ts/sb/runtime_text.ts`、`fixtures/ts/sb/runtime_bin.ts`、`fixtures/ts/sb/runtime_list.ts`、`fixtures/ts/sb/runtime_struct.ts` 当前仍主要是 barrel.
  - `fixtures/ts/sb/type.ts` 当前已经只做 re-export.
- 假设:
  - 用户要求的“收口”不仅是入口收口, 还包括最终生成目录的物理文件数收口.
  - 允许保留少量稳定中间层, 但不允许继续保留大批细粒度实现文件.
- 推断:
  - 下一步应回并模板与写盘清单, 而不是继续把 text/bin/struct/list 向更细处拆分.

## 4. 方案比较

### 方案 A: 只改 barrel, 不收物理文件
- 做法:
  - 保留现有细粒度 runtime 文件.
  - 仅新增或保留粗粒度 barrel.
- 优点:
  - 表面改动小.
- 缺点:
  - 不能真正减少最终生成结果的碎片度.
  - 与本轮用户确认的目标不一致.

### 方案 B: 回并粗文件直出, 保留少量中间层
- 做法:
  - 把细粒度实现重新并回粗粒度 runtime 文件.
  - 仅保留 `runtime_core.ts` 与 `runtime_meta.ts` 这类稳定聚合层.
  - `type.ts` 继续保留为超薄兼容 barrel.
- 优点:
  - 最符合本轮目标.
  - 生成结果真实收口.
  - 导入边界清晰, 后续维护稳定.
- 缺点:
  - 这轮模板和测试改动面较大.

### 方案 C: 半收口
- 做法:
  - 只回并 `text/bin/list/struct`.
  - 继续保留较多底层细文件.
- 优点:
  - 风险更低.
- 缺点:
  - 后续仍需继续收口.
  - 不能一次性达到目标状态.

## 5. 推荐方案
- 采用方案 B.
- 原因:
  - 与用户确认的验收边界一致.
  - 改完后最终生成目录与对外入口同时收敛.
  - 可以保留必要的中间层, 避免把所有实现重新堆成一个新巨石.

## 6. 最终文件边界

### 6.1 最终允许落盘的 runtime 主文件
- `type.ts`
- `runtime_base.ts`
- `runtime_header.ts`
- `runtime_text.ts`
- `runtime_bin.ts`
- `runtime_list.ts`
- `runtime_struct.ts`
- `runtime_enum.ts`

### 6.2 最终允许落盘的少量中间层
- `runtime_core.ts`
- `runtime_meta.ts`

### 6.3 继续存在的 schema 生成文件
- `enum.ts`
- `struct_*.ts`
- `rpc.ts`
- `_.ts`
- `rpc_smoke.test.ts`
- 以及 schema 派生的 smoke test / cross-language 辅助文件

### 6.4 不再落盘的细粒度 runtime 文件
- 例如:
  - `runtime_text_state.ts`
  - `runtime_text_io.ts`
  - `runtime_bin_length.ts`
  - `runtime_bin_io.ts`
  - `runtime_struct_scalar_field.ts`
  - `runtime_struct_default_list_field.ts`
  - 以及同类细粒度 runtime 实现文件

## 7. 模块职责设计

### 7.1 `type.ts`
- 角色: 兼容入口.
- 规则:
  - 只保留 `export * from`.
  - 不承载真实实现.

### 7.2 `runtime_base.ts`
- 角色: 最底层基础能力.
- 内容:
  - buffer
  - bit
  - integer
  - float
  - equal
  - result

### 7.3 `runtime_header.ts`
- 角色: header/state 编解码域.
- 内容:
  - bitmap
  - state block
  - header codec

### 7.4 `runtime_text.ts`
- 角色: text 完整实现域.
- 内容:
  - utf8
  - state
  - io
  - codec

### 7.5 `runtime_bin.ts`
- 角色: bin 完整实现域.
- 内容:
  - state
  - length
  - io
  - codec

### 7.6 `runtime_list.ts`
- 角色: list 完整实现域.
- 内容:
  - count/header 逻辑
  - bool list
  - default list
  - text list
  - bin list

### 7.7 `runtime_enum.ts`
- 角色: enum 元信息与 list helper.
- 内容:
  - enum meta
  - normalize
  - assignable
  - eq
  - enum list helper

### 7.8 `runtime_struct.ts`
- 角色: struct 元信息、编解码与 field builder.
- 内容:
  - struct meta
  - validate / eq
  - read / write
  - value field
  - list field

### 7.9 `runtime_core.ts`
- 角色: 基础聚合入口.
- 规则:
  - 只聚合 `runtime_base.ts`、`runtime_header.ts`、`runtime_text.ts`、`runtime_bin.ts`、`runtime_list.ts`.
  - 不承载真实实现.

### 7.10 `runtime_meta.ts`
- 角色: 元信息聚合入口.
- 规则:
  - 只聚合 `runtime_enum.ts`、`runtime_struct.ts`.
  - 不承载真实实现.

## 8. 依赖方向
- `type.ts` -> 只依赖粗文件与中间层.
- `runtime_core.ts` -> 只聚合基础运行时域.
- `runtime_meta.ts` -> 只聚合 enum/struct 元信息域.
- `runtime_enum.ts` / `runtime_struct.ts` 可以依赖 `runtime_core.ts`.
- `runtime_base.ts` / `runtime_header.ts` / `runtime_text.ts` / `runtime_bin.ts` / `runtime_list.ts` 不依赖 `runtime_meta.ts`.
- `enum.ts` / `struct_*.ts` / `rpc.ts` 不再依赖 `./type`.

## 9. 实施顺序

### 阶段 1: 先锁红灯
- 修改 `internal/tpl_ts_render_test.go`.
- 增加以下断言:
  - 最终仅允许粗文件与少量中间层文件存在.
  - 细粒度 runtime 文件不存在.
  - `type.ts` 仍是超薄 barrel.
  - `type.ts` 仍 re-export 代表性的公开 runtime 符号, 例如 `Buffer`、`writeHeader`、`getText`、`setBin`.
  - `enum.ts` / `struct_*.ts` / `rpc.ts` 只依赖粗入口.
- 补或收紧 TS 侧公开语义 smoke:
  - `fixtures/ts/sb/enum_smoke.test.ts` 继续验证生成 enum API 用法不变
  - `fixtures/ts/sb/rpc_smoke.test.ts` 继续验证 `RpcClient` 与方法入口存在
  - 若当前缺少对 `struct_*.ts` 公开用法的 TS 侧验证, 先补一个 focused smoke test 再进入模板回并
- 先运行 focused test, 确认现状不满足该目标.

### 阶段 2: 回并模板
- 把细粒度 runtime 模板并回粗模板.
- 目标映射表:
  - `runtime_buffer*` + `runtime_bit*` + `runtime_int*` + `runtime_u24` + `runtime_float` + `runtime_equal` + `runtime_result` -> `runtime_base.ts`
  - `runtime_bitmap` + `runtime_state_header` + `runtime_state_block` + `runtime_header_codec` -> `runtime_header.ts`
  - `runtime_text_utf8` + `runtime_text_state` + `runtime_text_io` + `runtime_text_codec` -> `runtime_text.ts`
  - `runtime_bin_state` + `runtime_bin_length` + `runtime_bin_io` + `runtime_bin_codec` -> `runtime_bin.ts`
  - `runtime_bitmap_list` + `runtime_list_count` + `runtime_list_default` + `runtime_item_list` + `runtime_text_list` + `runtime_bin_list` -> `runtime_list.ts`
  - `runtime_struct_core` + `runtime_struct_meta` + `runtime_struct_codec` + `runtime_struct_read` + `runtime_struct_write` + `runtime_struct_field` + `runtime_struct_value_field` + `runtime_struct_scalar_field` + `runtime_struct_primitive_field` + `runtime_struct_blob_field` + `runtime_struct_ref_field` + `runtime_struct_list_field` + `runtime_struct_bitmap_field` + `runtime_struct_bool_list_field` + `runtime_struct_default_list_field` + `runtime_struct_item_field` -> `runtime_struct.ts`
  - `runtime_enum` 保留为单文件真实实现, 仅整理 import 边界
- 优先处理:
  - `runtime_text`
  - `runtime_bin`
  - `runtime_list`
  - `runtime_struct`
  - `runtime_header`
  - `runtime_base`
  - `runtime_enum`
- `runtime_core.ts` 与 `runtime_meta.ts` 只允许保留 barrel 语句, 不允许承载真实实现.

### 阶段 3: 收写盘清单
- 修改 `internal/tpl_ts_runtime_render.go`.
- 修改 `internal/tpl_ts_render.go`.
- 在这一阶段统一收掉不再需要的 embed、render 函数与写盘清单项.
- 让 generator 只写最终允许落盘的文件.
- 明确清理契约:
  - generator 必须先计算允许存在的 runtime 文件集合
  - 对目标 TS 输出目录中, 任何“已生成但不在允许集合内”的 `runtime*.ts` 文件, 在本次生成期间自动删除
  - 不允许依赖手工删除旧细文件来达成收口结果
- 确认生成目录会清掉旧的细粒度 runtime 文件, 不留残留.

### 阶段 4: 重生成与回归
- 对现有 fixture 目录执行真实生成, 不手工删除旧细文件.
- 生成命令:
  - `go run . -go ./fixtures/go -ts ./fixtures/ts -tag bson,json ./fixtures/schema.sb`
- 该步骤必须在已有 `fixtures/ts/sb` 上执行, 用于验证 generator 会自动清理旧细粒度 runtime 文件.
- 先跑 focused regression.
- 再跑全量 `go test ./...`.
- 再跑 `bun test`.
- 最后检查生成目录的实际文件集是否符合设计.

## 10. 测试与验收

### 10.1 验收标准
- `fixtures/ts/sb/` 中不再出现细粒度 runtime 实现文件.
- `type.ts` 保持轻量 re-export.
- 经 `./type` 导入的兼容 runtime API 仍可编译并运行.
- 仅允许 `runtime_core.ts` 与 `runtime_meta.ts` 作为中间层 runtime 文件存在.
- `runtime_core.ts` 与 `runtime_meta.ts` 都只包含 re-export, 不承载真实实现.
- `enum.ts` / `struct_*.ts` / `rpc.ts` 不再从 `./type` 导入.
- TS 侧生成 API 公开用法不变:
  - `enum.ts` 相关 helper 仍按原方式可用
  - `struct_*.ts` 暴露的 `new*` / `get*` / `set*` / `eq*` 仍按原方式可用
  - `rpc.ts` 暴露的 `RpcClient` / `RpcErrCode` / 生成方法仍按原方式可用
- TS runtime 协议 golden vector 与 round-trip 行为不变.
- Go <-> TS cross-language wire consistency 行为不变.
- RPC 传输与返回值语义不变.
- `go test ./...` 通过.
- `bun test` 通过.

### 10.2 关键验证命令
- `go test ./internal -run TestTsGeneratorWritesSchemaFile -count=1`
- `go run . -go ./fixtures/go -ts ./fixtures/ts -tag bson,json ./fixtures/schema.sb`
- `go test ./fixtures/go/sb -run 'TestCrossLanguageWireConsistency|TestCrossLanguageWireRejectsMalformedInputs|TestProtocolGoldenVectors' -count=1`
- `go test ./...`
- `bun test`

### 10.3 语义验证落点
- TS runtime 语义:
  - 依赖 `fixtures/ts/sb/runtime.test.ts` 中的 header、bitmap、state block、text/bin、数值编解码断言
  - `fixtures/ts/sb/runtime.test.ts` 继续通过 `./type` 导入公开 runtime API, 作为兼容 barrel 的运行时验证
  - 若回并后这些断言不足以覆盖新风险, 先补 focused test 再改模板
- TS 公开 API 语义:
  - 依赖 `fixtures/ts/sb/enum_smoke.test.ts`
  - 依赖 `fixtures/ts/sb/rpc_smoke.test.ts`
  - 若当前缺少对 `struct_*.ts` 公开 API 的 TS 侧覆盖, 增加一个 focused smoke test, 验证代表性 struct 的 `new*` / `validate*` / `get*` / `set*` / `eq*` 使用方式不变
- Go/TS 一致性:
  - 依赖 `fixtures/go/sb/cross_consistency_test.go`
  - 必须验证固定样例与随机样例均通过
- Go golden / round-trip:
  - 依赖 `fixtures/go/sb/runtime_test.go`
  - 必须覆盖 protocol vector、header round-trip、state block 与 list bitmap 语义
- RPC 行为:
  - 依赖 `fixtures/go/sb/rpc_transport_test.go`
  - 若 TS 侧 `rpc.ts` 导入变化影响 smoke 流程, 补 focused smoke test 后再回并实现

### 10.4 TDD 执行要求
- 先写或改失败断言, 再改模板与写盘清单.
- 若发现现有语义测试未能直接证明“行为不变”, 必须先补 focused red test, 再进入实现.
- Focused red test 必须显式断言:
  - 只剩粗 runtime 文件 + `runtime_core.ts` + `runtime_meta.ts`
  - `runtime_core.ts` / `runtime_meta.ts` 仅为薄 barrel
  - 其他中间层 `runtime_*.ts` 文件不存在

## 11. 风险与缓解
- 风险: `runtime_struct.ts` 回并后重新变大.
  - 缓解: 接受其作为当前粒度下的单域文件, 但不再继续物理碎片化.
- 风险: 测试当前大量绑定细文件名, 回并时会出现一轮大面积断言改写.
  - 缓解: 先把允许落盘文件名单写成明确断言, 再回并实现.
- 风险: embed 与写盘清单漏删, 导致目录里残留旧文件.
  - 缓解: focused regression 中显式断言细文件不存在.
- 风险: `runtime_core.ts` / `runtime_meta.ts` 退化为新的大而全文件.
  - 缓解: 明确规定其只做聚合, 不承载真实实现.

## 12. 回滚路径
- 若回并过程中出现大范围行为偏差, 可回滚到当前半拆分状态:
  - 恢复 `internal/tpl_ts_runtime_render.go` 的细文件 embed 清单.
  - 恢复 `internal/tpl_ts_render.go` 的细文件写盘清单.
  - 恢复 `internal/tpl_ts_render_test.go` 对细文件结构的断言.
- 回滚后重新生成 fixtures 并执行同一套 Go/Bun 回归.

## 13. 生效边界与残余风险
- 生效边界:
  - 仅影响 TS generator 与 TS generated fixtures.
  - 不改变协议线上的编码规则.
- 残余风险:
  - `runtime_struct.ts` 和 `runtime_list.ts` 的文件体积仍可能偏大.
  - 若后续继续追求更细 tree-shaking, 可能需要重新评估“物理收口”和“构建优化”之间的平衡.
