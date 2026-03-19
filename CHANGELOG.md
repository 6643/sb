# Changelog

## [2026-03-19 00:00:00] [type: fix] [scope: rpc,go,ts,docs,test]
- [修复 Go RPC 客户端在 `Do(req)` 返回后立即 `cancel()` 导致的响应体读取误判问题, 将 `Client.Timeout` 明确为覆盖单次 attempt 全生命周期, 并补齐 body timeout / retry / cancel 传输测试]
- [修复 Go handler 直接把内部 `RpcErrCode` 写成 HTTP 状态码的问题, 为 `RpcNoConn` / `RpcRespErr` / 未知内部状态统一映射 `500`, 保持 `400/401/404/408` 等已知状态稳定]
- [修复 TS RPC 客户端对 `maxRespBytes` 仅在 `arrayBuffer()` 后校验的问题, 改为优先流式累计读取并在超限时立即失败; 同时将未知 HTTP 状态保留为原始数字并更新模板断言与文档说明]

## [2026-03-08 07:45:00] [type: refactor] [scope: ts,runtime,rpc,typing]
- [将 TS 侧剩余的 `class` 全部收口为接口：`Buffer`、`BitWriter`、`BitReader`、`RpcClient` 现在都以 `interface` 暴露类型面，并改为无 `class` 关键字的构造函数值实现]
- [同步更新 TS 运行时模板、RPC 模板、模板断言与生成产物，并重新执行生成、Go 全量测试、Bun 全量测试和 `bunx tsc` 类型检查]

## [2026-03-08 07:31:00] [type: refactor] [scope: ts,runtime,typing]
- [将 TS 运行时的 `Buffer` 从 `export class Buffer` 调整为 `export interface Buffer`，并保留同名构造值 `export const Buffer = class ...`，使类型面改为接口而不扩大调用面改动]
- [重新生成 TS 产物，并通过 Go 全量测试、Bun 全量测试与 `bunx tsc` 类型检查验证这次类型面调整]

## [2026-03-08 07:18:00] [type: refactor] [scope: naming,runtime,tplgen]
- [继续收口 helper 命名以保持一致：Go runtime 去掉 `Sized` / `Validated` 这类仅反映内部实现方式的后缀，统一为 `SetDefaultList`、`SetZeroList`、`SetDefaultPtrList`、`SizeDefaultPtrList`]
- [同步更新 Go/TS 模板、生成产物与断言，并重新执行生成、Go 全量测试、Bun 全量测试与 `bunx tsc` 类型检查]

## [2026-03-08 07:02:00] [type: refactor] [scope: naming,runtime,tplgen]
- [收口当前协议 runtime/helper 的命名：去掉 `Compact`、`Bitmap*List*`、`ZeroValue/ZeroFixed` 这类历史性和实现细节过重的后缀，统一为 `GetText/SetText`、`GetDefaultListInto/SetDefaultListSized`、`GetZeroListInto/SetZeroListSized`、`GetDefaultPtrListInto/SetDefaultPtrListValidated` 等更直接的名称]
- [同步更新 Go/TS 生成模板、生成产物、跨语言一致性脚本与模板断言，并重新执行生成、Go 全量测试、Bun 全量测试与 `bunx tsc` 类型检查]

## [2026-03-08 06:24:00] [type: test] [scope: go,ts,cross]
- [为双端字节级一致性回归新增随机样本覆盖：`go/sb/cross_consistency_test.go` 现在基于固定 seed 生成全量随机 struct / API 载荷，并新增 `TestCrossLanguageWireConsistencyRandom`]
- [`ts/sb/cross_consistency.test.ts` 同步增加随机对拍入口，并为跨语言对拍测试单独放宽超时，继续通过 Bun 驱动 Go↔TS 字节级 roundtrip 校验]

## [2026-03-08 06:08:00] [type: refactor] [scope: naming,docs,tplgen]
- [将当前主链路中的版本化命名收口为正式命名：协议文档更名为 `PROTOCOL.md`，内部模板文件更名为 `tpl_go_render.go`、`tpl_ts_render.go`、`tpl_ts_runtime_render.go` 等通用名称]
- [同步清理 README 与协议文档中的迁移期措辞，删除当前实现里对并行协议链路的引用，并把生成器/运行时/一致性回归的实际落点更新到文档]

## [2026-03-08 05:41:00] [type: refactor] [scope: ts,runtime,tplgen]
- [为 TS runtime 新增 `getZeroValueListCompact` / `setZeroValueListCompact`，把 number/bigint 这类零默认值 list 的容器级重复包装下沉到 runtime]
- [TS 生成器改为让基础标量 list 直接调用新 helper，减少生成产物中的重复 lambda 与模板分支；随后重生成产物并通过 Go/Bun/`bunx tsc`/双端一致性回归]

## [2026-03-08 05:32:00] [type: refactor] [scope: ts,runtime,tplgen,test]
- [为 TS runtime 新增 `readHeader` / `writeHeader`，统一处理 struct header 的状态打包与解包；TS 生成器改为为每个 struct 生成 `*HeaderWidths` 常量，不再在产物里手写 `BitReader` / `BitWriter` 的逐字段 header 位运算]
- [同步更新 TS 生成器回归断言，并在产物重生成后继续执行 Go 全量测试、Bun 全量测试、`bunx tsc` 与双端数据一致性回归]

## [2026-03-08 05:05:00] [type: refactor] [scope: go,runtime,tplgen]
- [为 Go runtime 新增 `ReadHeader` / `WriteHeader`，统一处理 struct header 的位图/状态打包与解包，并补充 `go/sb/rt/runtime_test.go` 覆盖跨字节与非法 padding/范围校验]
- [Go 生成器改为为每个 struct 生成 `*HeaderWidths` 常量，生成代码不再手写逐 bit 的 header pack/unpack 位运算，减少模板重复与产物代码量]

## [2026-03-08 04:46:00] [type: refactor] [scope: go,ts,test,docs]
- [将 Go runtime 路径从 `go/sb/v2` 收口为 `go/sb/rt`，同步更新生成器、生成产物导入路径与文档说明]
- [新增 `ts/sb/cross_consistency.test.ts`，将双端字节级一致性校验接入正式 `bun test` 入口，并复用 Go 侧权威对拍测试]

## [2026-03-08 04:31:00] [type: test] [scope: go,ts]
- [新增 `go/sb/cross_consistency_test.go` 与 `ts/sb/cross_consistency.ts`，用同一批结构体与 API 载荷样本做 Go↔TS 字节级 roundtrip 对拍，覆盖空值、常规值与关键边界值]

## [2026-03-08 04:02:00] [type: breaking] [scope: tplgen,go,ts,docs]
- [删除旧版 v1 生成链路与对应产物目录，生成器现在只保留当前紧凑协议实现；CLI 移除 `-go-v2` / `-ts-v2`，默认直接生成到 `go/sb` 与 `ts/sb`]
- [将活跃生成器入口与 helper 命名去掉 `v2` 前缀，并同步更新 README / `PROTOCOL_V2.md`，明确当前仓库只保留这一套默认协议]

## [2026-03-08 03:16:00] [type: refactor] [scope: tplgen,ts]
- [将 TS v1 runtime 的通用 `getList` / `setList` / `into` helper 收口为 `Err` 风格返回值，同时兼容回调返回 `Error | null | undefined`，避免高层生成代码与 runtime helper 的错误类型继续分裂]

## [2026-03-08 03:05:00] [type: fix] [scope: tplgen,ts,v2]
- [修正 TS v2 `bitmap list` runtime 与生成器之间的错误类型签名不一致：`getBitmapListCompact` / `setBitmapListCompact` 现在接收并返回 `Err`，避免 `readX` / `setX` 回调返回 `Error | undefined` 时与旧的 `Error | null` 签名冲突]
- [同步移除 TS v2 `struct/enum list` 包装层中多余的 `resultU` / `errU` 包装，直接透传 `Err`，使生成产物与 runtime 签名一致]
- [补齐 TS v2 primitive bitmap-list 生成：对 `rt.getU32/getI32/...`、`rt.setU32/setI32/...` 这类低层 `Error | null` API 在生成器侧做 `Err` 适配，并把对应判错条件从 `!== null` 收口为 `!== undefined`]

## [2026-03-08 02:12:00] [type: refactor] [scope: tplgen,ts]
- [将 TS 面向业务侧的生成函数统一收口为 `Error | undefined` / `[T, Err]` 风格；新增 `Err`、`errU`、`resultU` 适配器，避免高层继续暴露 `Error | null`]
- [保持底层 TS primitive runtime 的 `Error | null` 不变，仅通过适配层改动 `struct/list/v2 validate` 等高层生成签名，降低改动面并保持现有低层实现稳定]

## [2026-03-08 01:36:00] [type: refactor] [scope: tplgen,go,ts,v2]
- [移除 Go/TS v2 API 参数与返回值的自动 wrapper struct 组装，改为在 handler / RPC client 中按参数顺序直接对 `buf` 调用基础 getter/setter 与 struct codec，显著减少生成代码量]
- [停止生成 TS v2 的 `struct_api_*_{req,resp}.ts` 包装体文件，并在生成时清理旧 wrapper 产物，保持 Go server / Go client / TS client 的 v2 API wire format 一致]

## [2026-03-08 00:52:10] [type: refactor] [scope: tplgen,go,v2]
- [继续收紧 Go v2 的 API 导出边界：将 `api._.go` 中请求/响应包装体的主 codec 函数 `Size/Get/Set/Read/Eq` 改为包内私有，仅保留类型、handler、路由注册与 RPC 客户端入口导出]
- [同步更新 `rpc.go` / `api._.go` 生成调用与生成器回归断言，确认私有 codec 仍能在同包内正常协作]

## [2026-03-08 00:41:28] [type: refactor] [scope: tplgen,go,v2]
- [继续收紧 Go v2 的导出边界：把枚举校验函数 `IsX` 也改为包内私有，仅保留结构体主编解码 API 和 RPC/路由相关导出]
- [同步更新 `go/sbv2` 生成产物、`gopls check` 与生成器回归，确认导出面收口后仍可正常编译]

## [2026-03-08 00:33:27] [type: refactor] [scope: tplgen,go,v2]
- [收紧 Go v2 生成代码的导出边界：保留 `Get/Set/Read/Size/Eq` 这类主编解码 API，把 `Default/Normalize/IsDefault/IsAssignable/Reset/IsZero/Validate` 等 helper 改为包内私有，减少不必要的导出符号]
- [同步更新 `go/sbv2` 生成产物与生成器回归断言，保持同包文件间调用不变]

## [2026-03-08 00:09:14] [type: refactor] [scope: runtime,tplgen,go,v2]
- [补回缺失的 `go/sb/v2` 最小 runtime，实现 `HeaderSize`、`BitReader`、紧凑 `text/bin/list`、bitmap list 与 ptr-list 泛型 helper，恢复 `go/sbv2` 的真实可编译与可验证状态]
- [继续收敛 Go v2 的基础数值 list 模板：移除 `get/set/sizeXListBody` 这类重复生成包装，改为生成代码直接调用 `rt.GetZeroFixedListCompactInto`、`rt.SetZeroFixedListCompactSized`、`rt.SizeZeroFixedListCompact`]

## [2026-03-07 23:56:11] [type: refactor] [scope: runtime,tplgen,go,v2]
- [继续下沉 Go v2 的 `[]*Struct` 泛型复用：新增 `SizeBitmapPtrListCompactValidated` 与 `SetBitmapPtrListCompactValidated`，让生成器与手写 `GameList` 统一复用 ptr-list 验证/编码逻辑]
- [同步更新生成器回归与 `GameList` 路径，保持 nil-item 语义与现有测试一致]

## [2026-03-07 23:44:02] [type: refactor] [scope: runtime,tplgen,go,v2]
- [为 Go v2 runtime 新增 `GetBitmapPtrListCompactInto` 泛型 helper，并让生成器用它收敛 `[]*Struct` / `[]*ApiReqResp` 的 bitmap list 解码逻辑，减少重复生成代码]
- [补充 `go/sbv2` 回归，固定 `getSimListBodyReuse` 会复用同索引已有对象指针]

## [2026-03-07 23:30:18] [type: fix] [scope: tplgen,ts,v2]
- [修正 TS v2 生成器在嵌套 struct 字段编码时的空值收窄，`isZeroX` 分支内改用非空断言调用 `setX`，消除 `SimInfo | null` 传给 `setSimInfo` 的类型错误]

## [2026-03-07 23:19:44] [type: chore] [scope: tplgen,go,v2]
- [移除 Go v2 生成代码里可由编译器推断的泛型实参，清掉 `go/sbv2/api._.go` 等文件中的 `unnecessary type arguments` 诊断]

## [2026-03-07 23:10:21] [type: feat] [scope: protocol,runtime,v2]
- [将 v2 `list` 的外层 `count` 收紧为仅支持 `empty/u8`，`>255` 直接报错，`10/11` 状态对 `list count` 视为非法]
- [同步更新 Go/TS v2 runtime、协议草案和一致性边界样本，把 `256` 元素 list 的旧边界改为非法路径]

## [2026-03-07 22:55:33] [type: test] [scope: go,ts,v2]
- [为 `go/sbv2` 与 `ts/sbv2` 新增 `runtime_consistency` 回归，固定 `Sim` 与 `Sim list body` 在零值 / 常规 / 边界态样本上的 `encode -> decode -> re-encode` 数据一致性]
- [一致性回归同时校验对象 roundtrip 与字节级 re-encode 稳定性，不再把重点放在跨语言互通脚本上]

## [2026-03-07 22:36:54] [type: perf] [scope: runtime,tplgen,go,v2]
- [将 Go v2 struct header 编码改为生成期直接位写入，去掉 `BitWriter` 的头部临时对象，使 `SetSimV2` 和 `SetSimListV2` 降到 `0 allocs/op`]
- [将 `text/bin/bitmap list` 的读路径改为直接遍历 state/bitmap block，并把 `struct list` 解码改成 backing slice + 指针视图，去掉中间 `[]bool` / `[]uint8` 临时切片和每元素单独 struct 分配，降低 `GetSimV2` 与 `GetSimListV2` 的分配]
- [为 Go v2 生成器拆出 `size/set<Struct>Validated` 内部函数，去掉 `Set -> Size -> Validate` 和 `struct list` 内的重复校验；同时移除 `Get` 路径末尾的重复 `Validate` 与 enum 默认值重复赋值]
- [补充 Go v2 复用解码回归与 benchmark：重复向同一 `Sim` / `[]*Sim` 目标解码时会复用已有 struct/list 容量，`GetSimV2Reuse` 较普通路径少 5 次分配]

## [2026-03-07 20:44:31] [type: feat] [scope: bench,v1,v2]
- [新增 `bench_v1_v2.sh`，统一串起 `go/sb`、`go/sbv2`、`ts/sb` 与 `ts/sbv2` 的对照 benchmark，并支持通过环境变量调节 schema、迭代次数和是否跳过生成]
- [在 `README.md` 补充 `v1/v2` 基准对比入口和常用环境变量，避免每次手工拼 benchmark 命令]

## [2026-03-07 20:36:52] [type: test] [scope: go,ts,v2]
- [为 `go/sbv2` 新增 `runtime_smoke_test.go`，固定 `Sim` 与 `Sim list body` 的 roundtrip 回归]
- [为 `ts/sbv2` 新增 `runtime_smoke.test.ts`，通过 Bun 固定 `Sim` 与 `Sim list body` 的 roundtrip 回归]

## [2026-03-07 20:28:36] [type: feat] [scope: tplgen,ts,v2]
- [新增 `-ts-v2` 开关，并接入 `ts/sbv2` 生成流程，当前会生成 TypeScript `v2` 的 runtime、`enum/struct`、API 请求/响应 wrapper 和 RPC 客户端]
- [新增 `internal/tpl_ts_v2_generator.go` 与 `internal/tpl_ts_v2_generator_test.go`，并通过 Bun 对 `ts/sbv2` 的 smoke test 与 `Sim` roundtrip 做真实验证]

## [2026-03-07 20:03:41] [type: docs] [scope: readme,protocol,cli]
- [更新 `-go-v2` 帮助文本，去掉“暂不包含 API”的过期描述，并明确当前会生成 `struct/enum/API/RPC`]
- [同步修正 `README.md` 与 `PROTOCOL_V2.md` 的状态说明，补充 `go/sbv2` 与 `go/sb/v2` 的当前落地范围和 `-go-v2` 用法]

## [2026-03-07 19:42:18] [type: feat] [scope: tplgen,go,v2,api]
- [为 Go v2 生成器补上 API/RPC：生成 `go/sbv2/api._.go`、`go/sbv2/rpc.go` 和带 fingerprint 的 `api.<name>.go` stub]
- [v2 API 采用“请求/响应合成 struct”方案复用现有 v2 struct 编码，不再单独定义顶层裸值 wire format]

## [2026-03-07 19:14:27] [type: feat] [scope: tplgen,go,v2]
- [新增 `-go-v2` 开关，并把 Go v2 生成接入主流程；开启后会额外落盘到 `go/sbv2/schema.gen.go`，当前范围为 `enum/struct`，暂不生成 API]
- [新增 `internal/tpl_go_v2_generator.go` 与回归测试，基于 `go/sb/v2` runtime 生成 `header + fixed block + var tail` 风格的 Go v2 struct/enum 代码]

## [2026-03-07 18:39:11] [type: feat] [scope: runtime,go,v2]
- [为 `go/sb/v2` 增加三类 list body 原型：`[bool]` 的 `count + bitset`、`[text]/[bin]` 的 `count + item header block + item tail`，以及 `[struct]` 的 `count + bitmap + non-default bodies`]
- [补充 bitmap/state block 的 padding 零位校验，并为 `GameList`、`text/bin list` 与非法 padding/非法状态增加回归测试]

## [2026-03-07 18:24:18] [type: feat] [scope: runtime,go,v2]
- [为 `go/sb/v2` 增加紧凑字段 helper：`text/bin/list count` 的状态选择、读写与 canonical 校验，以及 `u16/u32` 小端读写辅助]
- [新增手写 `Game` 的 `v2` 编解码原型，落地 `header + fixed block + var tail`，并用文档中的 `0x20 / 0x80 / 0xA0` 示例做回归验证]

## [2026-03-07 18:01:32] [type: feat] [scope: runtime,go,v2]
- [新增独立 Go `v2` 子包 `go/sb/v2`，提供最小运行时骨架：连续 bitstream 的 `BitWriter/BitReader`、`u24` 读写、紧凑长度状态选择、bitmap/state block helper]
- [新增 `go/sb/v2/runtime_test.go`，覆盖 `Game` header、跨字节 header、`item header block`、`u24` 与长度状态选择等文档示例]

## [2026-03-07 17:34:42] [type: docs] [scope: protocol]
- [将 `PROTOCOL_V2.md` 的“待确认项”收敛为正式规则，明确 v2 采用严格解码并拒绝所有 non-canonical encoding]
- [补充 v2 运行时错误信息建议字段，要求至少包含路径、阶段、偏移、状态和值域信息]

## [2026-03-07 17:27:26] [type: docs] [scope: protocol]
- [为 `PROTOCOL_V2.md` 新增“实现任务清单”，按运行时、代码生成器、辅助函数、测试与落地顺序拆出可执行步骤]
- [收紧 `PROTOCOL_V2.md` 的待确认项，去掉已在正文中定死的 `enum` 默认值问题]

## [2026-03-07 17:16:11] [type: docs] [scope: protocol]
- [将 `PROTOCOL_V2.md` 中 `[struct]` 也并入统一的 `count + value bitmap + non-default bodies` 模型，把除 `[bool]` 外的 `1-bit default-state` list 收敛为同一套规则]
- [补充 `[struct]` list 的语义示例，并明确 `struct` 默认值按“所有字段都为默认值”判断]

## [2026-03-07 17:07:55] [type: docs] [scope: protocol]
- [将 `PROTOCOL_V2.md` 中 `[i8/u8/enum]` 也统一到 `count + value bitmap + non-default bodies`，把非 `bool` 的定宽标量 list 收敛为一套规则]
- [补充 `[i8] = [1, 2, 0, 0, 3]` 的位图示例，并明确 `enum` list 的默认值按 schema 第一项判断]

## [2026-03-07 17:02:08] [type: docs] [scope: protocol]
- [在 `PROTOCOL_V2.md` 中把 list 编码细化为三类：`[bool]` 继续使用 `count + bitset`，宽定宽标量 list 使用 `count + value bitmap + non-default bodies`，`[i8/u8/enum]` 维持 dense 编码]
- [补充宽定宽标量 list 的规范编码要求，并新增 `[u32] = [1, 2, 0, 0, 3]` 的位图示例]

## [2026-03-07 16:44:11] [type: docs] [scope: protocol]
- [补充 `PROTOCOL_V2.md` 对 `[text]` / `[bin]` 的 `item header block` 说明，明确该块整体会跨多个字节，但单个 `2 bit` 元素状态不会跨字节]
- [新增 `item header block` 示例，使用 5 个元素状态展示 `49 00` 的打包结果]

## [2026-03-07 16:37:18] [type: docs] [scope: protocol]
- [在 `PROTOCOL_V2.md` 增加跨字节 header 示例，使用 9-bit 的 `Demo` 结构展示 `2 bit` tag 被拆到两个字节时的具体落位]

## [2026-03-07 16:31:44] [type: docs] [scope: protocol]
- [将 `PROTOCOL_V2.md` 的 header 位打包规则改为“连续 bitstream + MSB -> LSB + 可跨字节”，并明确 `2 bit` tag 按高位在前写入]
- [同步更新 `Game` 示例中的 header 字节值，将示例从原先的低位在前写法改为 `0x20 / 0x80 / 0xA0`]

## [2026-03-07 16:18:57] [type: docs] [scope: protocol]
- [将 `PROTOCOL_V2.md` 中 `[text]` / `[bin]` 的列表元素编码也收敛到 v2 紧凑长度规则，定义为 `count + item header block + item tail`]
- [补充列表元素的规范编码要求，明确 `[text]` 使用 `u8/u16`，`[bin]` 使用 `u8/u16/u24` 的最短长度编码]

## [2026-03-07 16:08:44] [type: docs] [scope: protocol]
- [将 `PROTOCOL_V2.md` 中 `bin` 的紧凑长度状态从 `u32` 收敛为 `u24`，把单字段上限明确为 `16777215` (`16MB-1`)]
- [保留 `README.md` 中当前实现的 `bin=u32` 描述不变，避免把协议草案与现有实现混淆]

## [2026-03-07 15:20:36] [type: docs] [scope: protocol,readme]
- [新增 `PROTOCOL_V2.md`，整理 `SB v2` 的非兼容协议草案，明确零值省略、`text/bin/list` 的紧凑长度状态和 `header + fixed block + var tail` 布局]
- [在 `README.md` 增加协议草案入口，明确当前实现仍以 README 描述为准，`PROTOCOL_V2.md` 仅用于设计讨论]

## [2026-03-07 14:24:51] [type: perf] [scope: benchmark,go,ts]
- [新增 `go/sb/runtime_bench_test.go`，为 `Sim` / `SimList` 提供标准 `testing.B` 编码与解码基准，覆盖单结构与结构体列表的完整编解码热路径]
- [新增 `ts/sb/runtime_bench.ts`，提供基于 Bun 的可执行基准脚本，并在脚本启动阶段先做 `Sim` / `SimList` roundtrip 校验]
- [确认当前 Bun 环境不支持 `bun:test` 的 `bench` 导出，因此选择无需额外依赖的可执行脚本路径]
- [修复 TypeScript 枚举校验函数的 `switch case` 生成错误，避免只有最后一个枚举值会被识别为合法；并新增 `ts/sb/enum_smoke.test.ts` 做运行时校验]

## [2026-03-07 13:54:14] [type: perf] [scope: runtime,tplgen,ts]
- [为生成的 TypeScript 运行时新增 `Buffer.writeUnsafe` / `Buffer.rewindWrite`，将已预扩容的 `u8/text/bin` 及相关 list 写路径切到无重复容量检查的直接写入]
- [将生成的 TypeScript 结构体 setter 回滚改为调用 `buf.rewindWrite(startOffset)`，避免继续使用内部字段写回]
- [修复生成的 TypeScript RPC 客户端里请求编码变量与响应解码变量同作用域重名的问题，并通过 `bun test ts/sb/rpc_smoke.test.ts` 验证]

## [2026-03-07 13:38:30] [type: perf] [scope: tplgen,ts]
- [将生成的 TypeScript 结构体 setter 改为直接写入目标 `Buffer`，去掉中间 `body` 缓冲与尾部拷贝]
- [为直接写入路径增加 `startOffset` 回滚，保持 setter 出错时不会在目标 `Buffer` 中留下部分写入]
- [补充模板回归断言，并通过真实生成命令与 `go test ./...` 验证]

## [2026-03-07 13:27:14] [type: perf] [scope: tplgen,go,ts]
- [将生成的 Go `struct list` 解码改为直接顺序循环，去掉 `getList(ReadStruct)` 泛型分发]
- [将生成的 TypeScript `struct list` 编解码改为直接顺序循环，去掉 `_.getList` / `_.setList` 分发]
- [补充对应模板回归断言，并通过真实生成命令与 `go test ./...` 验证]

## [2026-03-07 13:17:42] [type: perf] [scope: runtime,go,ts]
- [为 Go 运行时的 `GetBinList` / `GetTextList` 增加专用循环路径，去掉泛型 `getList` + getter 闭包分发]
- [为 TypeScript 运行时的 `get/setBinList` 与 `get/setTextList` 增加专用循环与预扩容路径，减少通用 `getList` / `setList` 分发与重复扩容]
- [补充模板回归断言，并通过真实生成命令与 `go test ./...` 验证]

## [2026-03-07 13:07:10] [type: perf] [scope: runtime,tplgen,go]
- [优化 Go 运行时 `bool list` 编码，改为直接向目标缓冲追加长度和位图字节，去掉中间 `bits` 临时切片]
- [为生成的 Go `struct list` setter 增加总大小预估与一次性 `Grow`，避免列表编码退回到通用 `setList` 的逐步扩容路径]
- [优化 Go RPC 请求头写入路径，直接在 `RLock` 下把 `c.headers` 写入请求头，去掉每次调用的中间 map 拷贝]

## [2026-03-07 12:59:26] [type: perf] [scope: runtime,tplgen,go]
- [优化 Go 运行时的 `text` / `bin` 列表编码：`SetText` 改用 `WriteString` 避免字符串转 `[]byte` 分配，`SetTextList` / `SetBinList` 改为基于总大小预估的一次性扩容写入]
- [为生成的 Go RPC 请求缓冲和 API handler 响应缓冲增加大小预估与 `Grow`，减少多参数/大响应场景下的重复扩容]
- [补充对应模板回归断言，并通过真实生成命令与 `go test ./...` 验证]

## [2026-03-07 12:50:02] [type: perf] [scope: tplgen,go]
- [将生成的 Go API handler 改为直接顺序解码请求参数并直接编码响应，去掉 `parseRequest(...Getter)` / `sendResponse(...Setter)` 这层闭包与批处理分发]
- [保留现有请求体大小限制、空请求校验与状态码语义，同时将响应写入简化为原始 `[]byte` 输出]
- [补充模板回归断言，并通过真实生成命令与 `go test ./...` 验证]

## [2026-03-07 12:36:24] [type: perf] [scope: tplgen,go]
- [将生成的 Go 结构体 setter 改为在完成 bitmask/body 大小预估后直接写入目标 `bytes.Buffer`，去掉中间 `body` 缓冲与二次拷贝]
- [保持现有完整编解码语义不变，继续复用 `size<Struct>Body` 预估结果来一次性扩容目标缓冲]
- [补充模板回归断言，并再次通过真实生成命令与 `go test ./...` 验证]

## [2026-03-07 12:26:49] [type: perf] [scope: runtime,tplgen,go,ts]
- [为 Go 生成运行时新增定宽 primitive list 批量编解码快路径, 覆盖 `i8/i16/u16/i32/u32/i64/u64/f32/f64`，减少逐元素 getter/setter 调用与重复 buffer 操作]
- [为 TypeScript 生成运行时新增 `_getFixedWidthList` / `_setFixedWidthList`，将定宽 number/bigint list 切换为批量 DataView 读写路径]
- [补充模板回归断言并再次通过真实生成命令与 `go test ./...` 验证]

## [2026-03-07 11:44:22] [type: perf] [scope: runtime,tplgen,go,ts]
- [优化生成运行时的编码热路径: Go 定宽数值改为基于 `bytes.Buffer.AvailableBuffer` 的追加写入, `U8List`/`text`/`bin`/`bool list` 写入减少额外分配]
- [为生成的 Go 结构体 setter 增加常量 bitmask 大小与 body 预估扩容, 避免 `math.Ceil` 和动态 body 扩容带来的额外开销]
- [优化生成的 Go/TS RPC 客户端参数编码路径, 去掉 `SetAll` 闭包包装; TypeScript 运行时复用 `TextEncoder`/`TextDecoder` 并为 `U8List` 增加直接读写快路径]

## [2026-03-07 10:29:29] [type: docs] [scope: docs,readme]
- [在 `README.md` 的序列化对比章节补充性能对比说明, 新增定性性能表与 `SB` 相对 `JSON`、`Protocol Buffers`、`SBE`、`FlatBuffers`、`MessagePack` 的性能取舍说明]

## [2026-03-07 10:21:19] [type: docs] [scope: docs,readme]
- [在 `README.md` 的序列化对比章节补充精简选型建议表，便于快速判断 `SB`、`JSON`、`Protocol Buffers`、`SBE`、`FlatBuffers` 与 `MessagePack` 的适用场景]

## [2026-03-07 10:17:54] [type: docs] [scope: docs,readme]
- [在 `README.md` 的序列化选型对比中补充 `SBE`，并新增 `SB vs SBE` 的场景判断]

## [2026-03-07 10:08:19] [type: docs] [scope: docs,readme]
- [扩展 `README.md` 序列化对比章节, 新增 `FlatBuffers` 对照, 并将说明重写为以 `SB` 为中心的选型比较]

## [2026-03-07 10:00:01] [type: docs] [scope: docs,readme]
- [在 `README.md` 中新增 SB 与 `JSON`、`Protocol Buffers`、`MessagePack` 的优缺点对比, 明确当前二进制格式的适用边界]

## [2026-03-07 00:58:43] [type: docs] [scope: docs,readme]
- [将 Go 服务端示例更新为“两层超时”推荐配置: `http.Server` 连接级超时 + `TimeoutMiddleware` 请求级超时]
- [补充模板回归断言，固化 `DOC.md` 中 `http.Server` 超时字段、`rpcTimeout` 变量和 `ListenAndServe` 示例]

## [2026-03-07 00:44:00] [type: feat] [scope: tplgen,go,docs]
- [为生成的 Go 服务端路由新增 `TimeoutMiddleware`，通过 request context 为 handler 与业务函数注入可配置超时 deadline]
- [handler 在业务返回 `RpcOk` 但 request context 已超时时统一改写为 `RpcTimeout`，避免超时后仍返回成功状态]
- [更新模板回归、运行时行为测试与文档示例，补充服务端超时控制的用法说明]

## [2026-03-07 00:08:48] [type: refactor] [scope: internal,legacy,tplgen]
- [新增统一的生成产物目录写入器，收敛目录创建、按内容写盘、仅首次写入与指纹受控覆盖逻辑，减少 Go/TS 新旧生成器中的重复文件策略代码]
- [将模板生成器、legacy 生成器与文档输出统一切到公共写入器，保持 API stub 的人工修改保护语义不变]
- [补充公共写入器回归测试，固化“自动建目录、无变化不重写、已有手改文件不被指纹覆盖”的行为]

## [2026-03-06 23:51:39] [type: perf] [scope: internal,tplgen]
- [为模板生成与文档输出新增“内容未变化则跳过写盘”的统一写文件 helper，减少重复生成时的无效 IO 与时间戳抖动]
- [将当前 Go/TS 生成路径及旧 legacy Go/TS 生成器的稳定产物写入切换为按内容写盘，保留 API stub 的人工保护语义]
- [新增模板回归用例，固化重复生成不会改动无变化产物 `mtime` 的行为]

## [2026-03-06 20:24:48] [type: docs] [scope: readme,docs]
- 重组 README.md 为单一主手册, 合并原 peg 一致性说明, 统一快速开始, 语法规则, 规则边界, 生成约定与验证入口
- 删除 docs/peg_consistency.md, 修正文档中的旧路径与过时环境说明, 以当前 internal 根包结构为准

## [2026-03-06 19:57:59] [type: fix] [scope: internal,tplgen]
- [修正 `Tpl*` 批量重命名残留的字符串副作用, 恢复 `Content-Type` 头名、文档表头 `Type` 和 `Register*Api` 函数名]
- [同步更新根包 `tpl_*` 回归断言, 并再次通过 `go test ./...` 与真实生成命令验证]

## [2026-03-06 19:53:28] [type: refactor] [scope: internal,tplgen]
- [将 `internal/tplgen` 全量迁移到 `internal/` 根包, 新文件统一改为 `tpl_*` 命名]
- [将模板后端模型统一改为 `Tpl*` 前缀, 消除与 AST、IR 根类型冲突, 并更新 `main.go` 直接调用根包 `Generate`]

## [2026-03-06 19:39:57] [type: refactor] [scope: internal,legacy-ts]
- [将 `internal/gen/tsgen` 迁移到 `internal/` 根包, 新文件名为 `legacy_ts_generator.go` 与 `legacy_ts_generator_test.go`]
- [将旧 TypeScript 生成器入口和辅助函数统一改为 `LegacyTs` 前缀, 并在迁移完成后清理空目录 `internal/gen`]

## [2026-03-06 19:37:23] [type: refactor] [scope: internal,legacy-go]
- [将 `internal/gen/gogen` 迁移到 `internal/` 根包, 新文件名为 `legacy_go_generator.go` 与 `legacy_go_generator_test.go`]
- [将旧 Go 生成器入口和辅助函数统一改为 `LegacyGo` 前缀, 消除与后续 TS 旧生成器打平后的命名冲突]

## [2026-03-06 19:34:40] [type: refactor] [scope: internal,gen]
- [将 `internal/gen/config.go` 与 `internal/gen/note.go` 迁移到 `internal/` 根包, 新文件名为 `gen_config.go` 与 `gen_note.go`]
- [更新 `tplgen`、`main.go` 以及旧 `gogen/tsgen` 对 `Config` 和注释渲染 helper 的引用路径]

## [2026-03-06 19:32:17] [type: refactor] [scope: internal,semantic]
- [将 `internal/semantic` 迁移到 `internal/` 根包, 新文件名为 `semantic_resolver.go` 与 `semantic_resolver_test.go`]
- [更新 `main.go` 调用路径为根包 `Resolve`, 保持语义解析和展开逻辑不变]

## [2026-03-06 19:23:03] [type: refactor] [scope: internal,parser]
- [将 `internal/parser` 迁移到 `internal/` 根包, 新文件名为 `parser_schema.go` 与 `parser_schema_test.go`]
- [将解析器构造函数从 `New` 重命名为 `NewParser`, 解决与词法器 `New` 的根包命名冲突, 同步更新 `main.go`]

## [2026-03-06 19:16:44] [type: refactor] [scope: internal,lexer]
- [将 `internal/lexer` 的 `Token`、`Lexer` 与测试迁移到 `internal/` 根包, 新文件名分别为 `lexer_tokens.go`、`lexer_scanner.go`、`lexer_scanner_test.go`]
- [更新 `parser`、`parser_test` 与 `main.go` 对词法器的 import 路径, 保持词法行为与测试覆盖不变]

## [2026-03-06 18:31:06] [type: refactor] [scope: internal,ir]
- [将 `internal/ir/ir.go` 迁移到 `internal/ir_types.go`, 并将 IR 类型统一重命名为 `IR*` 前缀以避免与 AST 根类型冲突]
- [更新 `semantic`、`tplgen`、旧 `gen` 与 `main.go` 对 IR 模型的 import 路径和类型引用]

## [2026-03-06 18:31:06] [type: refactor] [scope: internal,ast]
- [将 `internal/ast/ast.go` 迁移到 `internal/ast_types.go`, 继续推进 `internal/` 根包平铺]
- [更新 `parser` 与 `semantic` 对 AST 类型的 import 路径, 保持 AST 结构与解析语义不变]

## [2026-03-06 18:21:27] [type: refactor] [scope: internal,util]
- [将 `internal/util/naming.go` 迁移到 `internal/util_naming.go`, 建立 `internal/` 根包的首个平铺文件]
- [更新 `tplgen` 与旧 `gen` 路径对 `SnakeCase`、`PascalCase`、`CamelCase`、`APIGroup` 的 import 引用]

## [2026-03-06 18:01:18] [type: feat] [scope: tplgen,go]
- [为 Go 运行时新增可选 `GetBinView` 借用接口, 保持现有 `GetBin` 继续返回独立拷贝]
- [补充 `type.go` 生成回归断言并在 README 中明确借用切片的生命周期约束]

## [2026-03-06 17:45:06] [type: refactor] [scope: protocol,tplgen]
- [调整 wire format 长度前缀: `list` 改为 `u16`, `text` 保持 `u16`, `bin` 升为 `u32`, 同步 Go/TS 运行时生成逻辑]
- [补充 `tplgen` 回归断言并将 `robust_size_limits` 样例列表长度提升到 `256`, 覆盖旧 `u8` 上限之外的场景]
- [更新 `README.md` 协议长度说明并通过真实生成命令与 `go test ./...` 全量回归]

## [2026-03-06 17:07:17] [type: refactor] [scope: tplgen]
- [删除 `runtime_sources.go` 中的 Go/TS 运行时大字符串常量, 改为 `renderGoRuntimeSource` 与 `renderTsRuntimeSource` builder 输出]
- [保持 `type.go` 与 `type.ts` 生成内容稳定, 并继续通过真实生成命令与全量 Go 回归]

## [2026-03-06 16:47:01] [type: refactor] [scope: tplgen]
- [移除 `internal/tplgen/_tpl/*`、`text/template` 与 `embed` 依赖, 将 Go/TS/DOC 生成统一重构为字符串拼接实现]
- [保留 API 指纹覆盖策略与文档输出能力, 新增 `source_writer`、静态运行时源码与纯代码渲染器]
- [重生成 `go/sb` 与 `ts/sb` 产物并通过 `go test ./...` 与真实 `go run . -go ./go -ts ./ts -tag bson,json aaa.sb` 回归]

## [2026-03-06 14:59:16] [type: fix] [scope: tplgen,ts]
- [修正 `ts.rpc.tpl` 与 `ts.enum.tpl` 的空白控制, 消除生成的 TS 注释首行多缩进、语句拼接同一行和枚举闭括号前多余空行]
- [补充模板回归断言, 固化 `rpc.ts` 与 `enum.ts` 的稳定输出格式]

## [2026-03-06 14:46:13] [type: fix] [scope: tplgen,ts]
- [修正 `ts.struct.tpl` 对可空 struct 字段生成的 `eq` 逻辑, 先做空值一致性判断再做深比较]
- [修复 `struct_recharge_b.ts` 等生成文件将 `Struct | null` 直接传给 `eqStruct` 导致的 TypeScript `TS2345` 报错]

## [2026-03-06 14:18:27] [type: fix] [scope: tplgen,ts]
- [为 `ts.rpc.tpl` 生成的 `_.setAll` 闭包参数补显式类型 `_.Buffer`, 消除 `TS7006` 隐式 `any` 报错]
- [补充模板回归断言并重生成 `ts/sb/rpc.ts`, 对齐 TypeScript 严格模式要求]

## [2026-03-06 14:02:19] [type: fix] [scope: internal/gen]
- [为旧的 `gogen` 与 `tsgen` 路径补充多行 note 安全渲染 helper, 避免多行说明直接打坏生成代码]
- [将 enum、struct、rpc 注释输出改为逐行安全注释, 并新增旧生成器回归测试]

## [2026-03-06 13:08:26] [type: fix] [scope: tplgen]
- [为 Go、TypeScript 与文档模板新增多行 note 安全渲染函数, 统一将多行说明展开为合法注释或 Markdown 内容]
- [修复 `enum.go`、`rpc.go`、`enum.ts`、`rpc.ts` 与 `DOC.md` 在多行 note 下生成非法内容的问题, 并补充模板回归测试]

## [2026-03-06 12:58:06] [type: fix] [scope: aaa.sb]
- [将 `aaa.sb` 中 struct 字段的行尾注释前移为字段头顶注释, 对齐当前语法约束]
- [保持字段顺序与业务含义不变, 仅调整注释位置以恢复 schema 可解析性]

## [2026-03-06 12:50:57] [type: fix] [scope: lexer,parser,grammar,docs]
- [将 struct 嵌入语法从隐式裸类型切换为显式 `...TypeName`, 新增 `Ellipsis` 词法单元与对应 PEG 规则]
- [收敛 parser: 仅接受 `...TypeName` 嵌入, 旧写法 `User` 在 struct 中改为按缺少类型报错]
- [补齐 lexer/parser 回归测试与文档示例, 同步 `README.md`、`docs/peg_consistency.md` 与 `aaa.sb`]

## [2026-03-06 12:36:05] [type: docs] [scope: docs]
- [补齐 `docs/peg_consistency.md` 回归测试索引遗漏项: `TestParseEnumCommentAfterPipeRejected`]
- [对齐 enum 成员注释禁用规则的分支测试与文档索引]

## [2026-03-06 12:06:19] [type: fix] [scope: parser,test]
- [新增 `TestParseEnumCommentAfterPipeRejected`, 固化 `|` 后注释同样按成员注释非法处理]
- [补齐 enum 成员注释禁用规则的分支回归覆盖]

## [2026-03-06 12:03:35] [type: docs] [scope: readme]
- [在 README 的枚举章节补充约束说明: enum 仅支持定义前整体注释, 不支持成员注释]
- [降低新用户按旧成员注释写法编写 schema 的误用风险]

## [2026-03-06 12:00:08] [type: docs] [scope: docs]
- [在 `docs/peg_consistency.md` 回归测试索引补充 `TestParseEnumHeaderCommentAccepted`]
- [对齐文档索引与 parser 最新正向回归用例]

## [2026-03-06 11:58:14] [type: fix] [scope: parser,test]
- [新增 `TestParseEnumHeaderCommentAccepted`, 固化 enum 定义前整体注释可解析并绑定到 `Enum.Note`]
- [补齐“禁成员注释”后的正向回归覆盖, 避免后续误伤整体注释能力]

## [2026-03-06 11:56:50] [type: docs] [scope: readme]
- [修正 README 枚举示例, 移除成员注释写法, 对齐“enum 仅支持定义前整体注释”规则]
- [消除 README 示例与 grammar/parser 行为冲突]

## [2026-03-06 11:55:33] [type: docs] [scope: docs]
- [修正 `docs/peg_consistency.md` 枚举注释规则描述为“仅支持定义前整体注释, 不支持成员注释”]
- [消除文档规则与最新 grammar/parser 行为冲突]

## [2026-03-06 11:54:21] [type: docs] [scope: docs]
- [同步 `docs/peg_consistency.md` 回归测试索引名, 将 enum 成员注释用例更新为 `TestParseEnumMemberHeadCommentRejected`]
- [消除文档索引与 parser 实际测试函数名不一致问题]

## [2026-03-06 11:51:48] [type: fix] [scope: parser]
- [收敛 `parseEnum` 注释策略: 枚举仅支持定义前整体注释, 禁止成员头顶注释与成员行尾注释]
- [统一枚举成员注释错误提示为 `枚举仅支持定义前整体注释, 不支持成员注释`]
- [同步 parser 用例: 成员头顶注释与成员行尾注释均应被拒绝]

## [2026-03-06 11:45:07] [type: docs] [scope: readme,docs,grammar]
- [更新 README 的 struct/enum 示例, 改为头顶注释行写法, 移除行尾注释示例]
- [更新 `docs/peg_consistency.md` 已对齐规则与回归测试索引, 对齐最新注释约束与测试名]
- [更新 `grammar.peg` 末尾语义约束注释, 移除过期的 inline 绑定描述]

## [2026-03-06 11:41:51] [type: fix] [scope: semantic]
- [新增语义用例: 同一 enum 中显式成员值重复应报 `成员 ID 重复`]
- [新增语义用例: 不同 enum 之间允许使用相同成员值]
- [将 enum 值唯一性范围固化为“单 enum 作用域”]

## [2026-03-06 11:09:20] [type: fix] [scope: parser,grammar]
- [枚举成员注释改为仅支持头顶注释行绑定, 支持连续多行 `//` 合并为成员注释]
- [禁止枚举成员行尾注释, 命中时返回 `枚举成员注释必须独占一行并写在成员前`]
- [同步 `grammar.peg` enum 成员注释约束与 parser 测试用例]

## [2026-03-06 11:04:52] [type: fix] [scope: parser,grammar]
- [结构体字段注释改为仅支持头顶注释行绑定, 支持连续多行 `//` 合并为字段注释]
- [禁止结构体字段行尾注释, 命中时返回 `字段注释必须独占一行并写在字段前`]
- [同步 `grammar.peg` 字段注释约束与 parser 测试用例]

## [2026-03-06 03:34:42] [type: fix] [scope: test]
- [将 `test/harness` 下 3 个 Go 联测源文件改为 `.go.txt` 模板, 避免被根目录 `go test ./...` 直接编译]
- [更新 `test/run_go_bun_tests.sh` 复制路径, 运行时仍生成 `*_integration_test.go` 参与联测]
- [更新 `test/README.md` 目录说明, 明确 harness 文件为测试模板]

## [2026-03-06 03:28:33] [type: fix] [scope: test]
- [补齐 `test/harness/robust_core_go_test.go`, 恢复 core 场景 Go 联测注入源文件]
- [补齐 `test/harness/robust_size_limits_go_test.go` 与 `test/harness/robust_optional_matrix_go_test.go`]
- [修复 `run_go_bun_tests.sh` 复制 harness 阶段 `cp: cannot stat` 失败]

## [2026-03-06 03:25:05] [type: fix] [scope: parser,tplgen]
- [在 `parseEnumMember` 增加 `TokenError` 直通分支, 枚举值位置遇到非法字符时保留词法错误文本]
- [修复 Go RPC 模板中 `enum list` 参数签名类型, 统一生成为 `TypeNameList` 而非切片逻辑类型]
- [对齐 `internal/parser` 与 `internal/tplgen` 回归测试预期, 消除当前 `go test ./...` 两处失败]

## [2026-03-06 03:20:01] [type: docs] [scope: docs]
- [在 `docs/peg_consistency.md` 新增“完成判定”章节]
- [补齐交付闭环核对项: 验证结论, 测试清单, 风险清单, 变更记录]
- [本次仅文档收口, 不涉及实现变更]

## [2026-03-06 03:17:00] [type: docs] [scope: docs]
- [在 `docs/peg_consistency.md` 新增“下一步验证建议”章节]
- [补充 Go 工具链恢复后的标准验证动作: `go test ./...` 与 `./test/run_go_bun_tests.sh`]
- [新增失败定位路径说明, 指向文档中的关键回归测试索引]

## [2026-03-06 03:14:18] [type: docs] [scope: docs]
- [在 `docs/peg_consistency.md` 新增“验证现状”章节, 明确当前环境无法执行 `go test`]
- [补充 `go/` 目录用途说明, 避免将其误解为 Go 工具链]
- [标注当前一致性结论基于规则对照与静态审查]

## [2026-03-06 03:09:55] [type: docs] [scope: docs]
- [补齐 `docs/peg_consistency.md` 的测试索引遗漏项]
- [新增 `TestParseEnumNegativeValueRejected` 与 `TestParseStructCommentLineNotInlineFieldNote` 条目]
- [本次仅文档索引完善, 不涉及实现与语义变更]

## [2026-03-06 03:08:06] [type: docs] [scope: docs]
- [将 `docs/peg_consistency.md` 的嵌套列表改为单层扁平列表, 统一文档结构风格]
- [保持原有信息不变, 仅调整呈现方式以降低维护和评审成本]
- [本次仅文档格式变更, 不涉及 lexer/parser/grammar 语义]

## [2026-03-06 03:05:03] [type: docs] [scope: docs]
- [新增 `docs/peg_consistency.md`, 汇总 PEG 与 lexer/parser 的一致性现状]
- [列出仍需语义阶段处理的约束边界: 唯一性, 类型存在性, 枚举值冲突, 依赖关系]
- [补充关键回归测试索引, 便于后续验证和评审]

## [2026-03-06 03:02:28] [type: fix] [scope: grammar.peg]
- [将 `CommentLine` 前导空白从 `WS*` 收敛为 `HWS*`, 明确仅允许行内空白]
- [消除 `WS*` 可跨行吞并导致的注释行入口歧义, 与 parser 的换行分支模型保持一致]
- [本次仅文法规约收敛, 不涉及 lexer/parser 实现]

## [2026-03-06 03:00:27] [type: fix] [scope: parser]
- [为 `parseStruct` 注释分支补充 EOF 守卫, 统一执行“注释行必须以换行结束”约束]
- [修复结构体内末尾注释无换行时的错误语义偏差, 不再退化为 `结构体缺少 }`]
- [新增 parser 用例覆盖 `User {\\n// only comment` 非法场景]

## [2026-03-06 02:56:44] [type: fix] [scope: lexer]
- [移除 `internal/lexer/lexer.go` 中未使用的 `unicode` 导入]
- [消除潜在编译错误 `imported and not used`, 恢复词法模块基础可编译性]
- [本次仅清理构建阻塞点, 不涉及词法语义变更]

## [2026-03-06 02:54:33] [type: fix] [scope: parser]
- [同步修正 `TestParseEnumNegativeValueRejected` 断言, 匹配负数在 lexer 阶段失败的新行为]
- [将期望从 `无效枚举值` 调整为 `未预期字符 '-'`, 避免测试语义与实现脱节]
- [本次仅更新测试断言与记录, 不涉及运行时逻辑]

## [2026-03-06 02:52:47] [type: fix] [scope: lexer]
- [将数字词法从 `unicode.IsDigit` 收敛为 ASCII `0-9`, 与 `grammar.peg` 的 `[0-9]` 规则一致]
- [`NextToken/readNumber` 统一改为 ASCII 数字判定, 杜绝非 ASCII 数字被误识别]
- [新增 lexer 用例覆盖全角数字输入 `１２`, 锁定返回 `TokenError`]

## [2026-03-06 02:51:37] [type: fix] [scope: parser]
- [移除 `parseAPI` 中未使用的局部变量 `line`, 消除编译期未使用变量风险]
- [本次仅清理编译阻塞点, 不涉及语义行为变更]
- [为后续 `go test` 恢复基础可编译条件]

## [2026-03-06 02:48:04] [type: fix] [scope: lexer]
- [收紧数字词法规则, 移除 `-` 前缀数字识别, 与 `grammar.peg` 的非负 `Number` 约束一致]
- [`readNumber` 改为仅消费 `[0-9]+`, 负数字面量前移到词法阶段失败]
- [新增 lexer 用例覆盖 `-1` 输入, 锁定返回 `TokenError`]

## [2026-03-06 02:43:04] [type: fix] [scope: lexer]
- [将 `skipSpaces` 收敛为仅跳过 ASCII `space/tab`, 不再使用 `unicode.IsSpace`]
- [消除非 ASCII 空白字符被静默吞掉的宽松行为, 与 `grammar.peg` 的 `HWS` 约束对齐]
- [新增 lexer 用例覆盖 `NBSP`(`\\u00a0`) 不应被跳过, 应返回 `TokenError`]

## [2026-03-06 02:39:55] [type: fix] [scope: lexer]
- [将标识符词法规则收紧为 ASCII: `[A-Za-z_][A-Za-z0-9_]*`, 与 `grammar.peg` 一致]
- [移除 `unicode.IsLetter/IsDigit` 在标识符判定中的宽松行为, 避免非 ASCII 标识符绕过 PEG 约束]
- [新增 lexer 用例覆盖中文标识符输入, 锁定返回 `TokenError`]

## [2026-03-06 02:38:27] [type: fix] [scope: parser]
- [按 PEG 收紧注释行规则, 末尾无换行的 `//` 注释不再被接受]
- [在 `ParseSchema` 的注释分支增加 EOF 守卫, 返回 `注释行必须以换行结束`]
- [新增 parser 用例覆盖 `// only comment`(无尾换行) 非法场景]

## [2026-03-06 02:34:31] [type: fix] [scope: parser]
- [收紧枚举显式值格式, 拒绝带前导零的数字(如 `01`), 与 `grammar.peg` 的 `Number` 规则一致]
- [在 `parseEnumMember` 中新增前导零守卫, 统一复用 `无效枚举值` 报错语义]
- [新增 parser 用例覆盖 `Status = Ok(01)` 非法场景]

## [2026-03-06 02:32:53] [type: fix] [scope: grammar.peg,parser]
- [修复 `TypeRef <- ... / Ident` 对 `nil` 的误匹配, 新增 `UserTypeName <- !NilType Ident`]
- [将 `TypeRef/NonNilTypeName` 的自定义类型分支统一改为 `UserTypeName`, 明确排除关键字 `nil`]
- [新增 parser 用例覆盖 `User { v nil }` 非法场景, 锁定 `nil 仅允许作为 API 返回类型` 报错]

## [2026-03-06 02:29:30] [type: fix] [scope: grammar.peg,parser]
- [将 `InlineComment` 从 `WS*` 收敛为 `HWS*`, 明确 inline 注释只能与声明同一行]
- [消除 `WS*` 可跨行导致的文法歧义, 与 parser 的同一行绑定行为对齐]
- [新增 parser 用例覆盖 struct 字段后跨行注释不应绑定为字段 note]

## [2026-03-06 02:26:08] [type: fix] [scope: grammar.peg,parser]
- [将 `Number` 规则进一步收敛为 `0-255`, 与枚举值 `uint8` 实现保持一致]
- [语法层显式约束上限, 避免 `256+` 延后到语义阶段才失败]
- [新增 parser 用例覆盖 `Status = Ok(256)` 非法场景, 锁定 `无效枚举值` 报错]

## [2026-03-06 02:24:16] [type: fix] [scope: parser]
- [修复 `parseType(false)` 对 `[nil]` 的误放行, 字段/参数位置不再接受 `nil` 作为数组元素类型]
- [将数组元素 `nil` 拒绝逻辑改为与 `allowNil` 无关, 全场景统一返回 `nil 不能作为数组元素类型`]
- [新增 parser 用例覆盖 `User { items [nil] }` 非法场景]

## [2026-03-06 02:22:45] [type: fix] [scope: grammar.peg,parser]
- [将 `Number` 语法从可选负号收敛为非负整数, 与枚举值 `uint8` 语义一致]
- [新增 parser 回归用例覆盖 `Status = Ok(-1)` 非法场景, 锁定 `无效枚举值` 报错]
- [本次仅修正文法约束与测试覆盖, 不改解析实现]

## [2026-03-06 02:19:25] [type: fix] [scope: parser]
- [收紧 `parseType` 对 `nil` 的接受范围, 仅允许在 `API` 返回类型位置使用]
- [参数/字段位置命中 `nil` 时返回专用错误 `nil 仅允许作为 API 返回类型`]
- [新增 parser 用例覆盖 `user.set(v nil) => u8` 非法场景]

## [2026-03-06 02:16:18] [type: fix] [scope: lexer]
- [修复 `readComment` 在单独 `\\r` 上的越界吞并问题, 注释读取在 `\\r` 或 `\\n` 都会停止]
- [新增 `TestCommentStopsAtCR`, 锁定 `// note\\rabc` 不会把 `\\rabc` 吞进注释]
- [与已收敛的换行规范保持一致, 单独 `\\r` 继续按非法输入处理]

## [2026-03-06 02:12:47] [type: fix] [scope: lexer,grammar.peg]
- [收敛换行规范为仅支持 `\\n` 与 `\\r\\n`, 删除单独 `\\r` 兼容路径]
- [lexer 仅在 `\\n` 或 `\\r\\n` 产出 `TokenNewLine`, 单独 `\\r` 返回 `TokenError`]
- [grammar.peg 同步移除 `NewLine` 的 `\\r` 分支, 并将 `WS` 改为 `(HWS / NewLine)+` 防止绕过]

## [2026-03-06 02:09:10] [type: docs] [scope: grammar.peg]
- [补齐 `NewLine` 规则, 增加对单独 `\\r` 的匹配, 与 lexer 当前行为对齐]
- [保持 `\\r\\n` 和 `\\n` 规则不变, 仅修正文法覆盖面]
- [本次仅文档语法变更, 不涉及 parser/lexer 实现]

## [2026-03-06 02:07:53] [type: fix] [scope: parser]
- [修正 `TestParseSchema` 基线样例, 移除已废弃的结构体字段行尾逗号写法]
- [将示例输入改为多行字段合法形态, 与当前 `逗号独占行` 规则一致]
- [本次仅调整测试数据与记录, 不涉及解析逻辑]

## [2026-03-06 02:05:14] [type: fix] [scope: parser]
- [新增结构体逗号独占行正向回归用例, 明确 `id u32` + `,` + `name text` 为合法写法]
- [与已收紧的“字段行尾逗号非法”规则形成正反对照, 防止后续误收敛]
- [本次仅补测试与记录, 不改解析实现]

## [2026-03-06 02:01:30] [type: fix] [scope: parser,grammar.peg]
- [将结构体逗号规则收敛为仅允许独占行, 禁止字段行尾逗号 `id u32,`]
- [parser 在字段行尾命中 `,` 时统一返回 `逗号必须独占一行`, 与 `StructCommaLine` 语义对齐]
- [grammar.peg 移除 `StructFieldLineEnd` 中 `Comma?` 并新增回归用例覆盖字段行尾逗号非法场景]

## [2026-03-06 01:56:21] [type: fix] [scope: parser]
- [修复 `parseAPI` 尾注释判定副作用, 不再把 API 结束后的独立注释行误判为尾注释]
- [保留结果类型同一行尾注释拒绝策略, 继续要求 API 注释写在定义前]
- [新增 parser 用例覆盖 `api -> 换行注释 -> 下一个 api` 的头注释继承场景]

## [2026-03-06 01:51:54] [type: fix] [scope: parser]
- [收紧 `parseAPI` 尾部注释判定, 不再依赖 API 名称行号]
- [结果类型后只要紧跟注释即报错 `API 注释必须写在定义前`, 与 `grammar.peg` 的 `ApiDef` 对齐]
- [新增 parser 用例覆盖 `=>` 换行后结果类型同一行尾注释非法场景]

## [2026-03-06 01:46:25] [type: fix] [scope: parser]
- [补充结构体逗号行级约束回归用例, 覆盖 `User { id u32, name text }` 非法场景]
- [锁定错误文案 `逗号必须独占一行`, 防止后续回退为泛化报错]
- [与 parser 已落盘的 `StructCommaLine` 行模型保持一致]

## [2026-03-06 01:42:47] [type: fix] [scope: parser]
- [收紧 `parseEnum` 与 `parseAPI` 头部换行容忍度, 禁止名称与关键分隔符跨行]
- [`enum` 名称与 `=` 必须同一行, `api` 名称与 `(` 必须同一行, 与 `grammar.peg` 的 `HWS*` 约束对齐]
- [新增 parser 内部用例覆盖 `Status\\n= Ok` 与 `user\\n(id u8) => u8` 非法场景]

## [2026-03-06 01:39:22] [type: docs] [scope: grammar.peg]
- [将 `StructDef/EnumDef/ApiDef/ApiName` 头部间隔从 `WS*` 收敛为 `HWS*`, 明确同一行约束]
- [将 `StructFieldLineEnd` 与 `StructCommaLine` 的行尾空白约束收敛为 `HWS*`, 避免换行被空白规则吞并]
- [统一 `grammar.peg` 对行内空白与换行边界的表达, 与 parser 当前行为保持一致]

## [2026-03-06 01:35:06] [type: docs] [scope: grammar.peg]
- [重写 `StructBody` 语法为显式行模型, 新增 `StructFieldLine/StructFieldLineEnd` 规则]
- [在 PEG 中明确“结构体字段每行最多一个”, 字段行仅允许可选逗号与行内注释]
- [新增 `NewLine` 或 `&(RBrace)` 行尾约束, 与 parser 当前“一字段一行”行为保持一致]

## [2026-03-06 01:30:09] [type: fix] [scope: parser]
- [新增结构体字段行尾守卫, 强制 `struct` 每个字段独占一行]
- [当同一行出现连续字段定义时返回专用错误: `结构体 <Name> 字段必须独占一行`]
- [新增 parser 用例覆盖 `User { id u32 name text }` 非法场景]

## [2026-03-06 01:27:34] [type: fix] [scope: parser]
- [新增 parser 侧 `TokenNewLine` 消费逻辑, 恢复顶层/struct/enum/api 在显式换行 token 下的可解析性]
- [修正 `parseField` 对换行结束字段的判定, 避免 `TokenNewLine` 导致嵌入字段误报缺少类型]
- [补充 parser 回归用例, 覆盖 enum 多行成员, struct 多行字段, api 多行参数场景]

## [2026-03-06 01:22:33] [type: fix] [scope: lexer]
- [新增 `TokenNewLine`, 将换行从被动吞掉改为显式 token 输出, 支持 `\\n/\\r/\\r\\n`]
- [`skipWhitespace` 收敛为仅跳过非换行空白, 为 PEG 行级规则提供词法基础]
- [新增 lexer 用例覆盖 LF 与 CRLF 的 `TokenNewLine` 产出及行号推进]

## [2026-03-06 01:17:14] [type: fix] [scope: parser]
- [新增枚举成员逗号分隔专用报错, `Status = Ok, Err` 明确提示 `成员之间仅允许 |`]
- [在 `parseEnum` 中前置识别 `TokenComma`, 避免退化为后续泛化错误]
- [新增 parser 用例覆盖逗号分隔非法场景并锁定错误文案]

## [2026-03-06 01:13:14] [type: fix] [scope: parser]
- [将 `parseEnum` 内部逻辑改为强制要求 `=` 分隔符, 移除可选兼容分支]
- [即使未来绕过顶层入口直接调用 `parseEnum`, 仍会稳定返回 `枚举定义必须包含 =`]
- [新增 parser 内部回归用例 `Status Ok`, 锁定无 `=` 时的错误语义]

## [2026-03-06 01:10:24] [type: fix] [scope: parser]
- [新增结构体字段 `tag` 误用场景专用报错: 字段名后直接字符串时提示 `缺少类型, 不允许直接写 tag`]
- [在 `parseField` 中前置识别 `TokenString` 分支, 避免退化为泛化缺少类型错误]
- [新增 parser 用例覆盖 `User { id \"_id\" }` 非法定义]

## [2026-03-06 01:08:17] [type: fix] [scope: parser]
- [新增枚举成员缺少分隔符 `|` 的专用报错, 覆盖同一行 `Status = Ok Err` 场景]
- [在 `parseEnum` 中前移错误识别, 避免后续退化为泛化顶层定义错误]
- [新增 parser 用例锁定 `成员之间缺少 |` 报错文案]

## [2026-03-06 01:05:20] [type: fix] [scope: parser,grammar.peg]
- [将 API 名称规则收敛为 1-2 段, 仅支持 `a` 或 `a.b`]
- [parser 新增三段名称拒绝逻辑, 对 `a.b.c` 返回 `API 名称仅支持 1-2 段`]
- [新增 parser 用例覆盖单段通过, 双段通过, 三段失败]

## [2026-03-06 00:59:53] [type: fix] [scope: parser]
- [新增 API 缺少 `=>` 的专用报错, 从通用期望提示改为 `API 定义缺少 =>`]
- [补充 parser 用例覆盖 `user.get_count(page u8) u8` 非法定义场景]
- [保持 PEG 约束不变, 仅提升错误定位准确性]

## [2026-03-06 00:55:41] [type: fix] [scope: parser]
- [新增枚举省略 `=` 的专用报错, 从泛化提示改为 `枚举定义必须包含 =`]
- [保留严格 PEG 约束行为, 仅改进错误可读性与定位效率]
- [更新 parser 测试断言, 锁定错误文案避免回退]

## [2026-03-06 00:50:50] [type: fix] [scope: lexer]
- [按 PEG 约束移除反引号字符串支持, lexer 仅将双引号识别为 `TokenString`]
- [新增 lexer 用例: 双引号字符串解析通过, 反引号字符串返回 `TokenError`]
- [统一实现与 `grammar.peg` 的字段 tag 字符串规则, 消除文法与实现偏差]

## [2026-03-06 00:40:07] [type: fix] [scope: parser]
- [顶层枚举定义改为严格要求 `=` 连接符, 移除 `Name | Member...` 兼容入口]
- [新增 parser 用例: `Status = Ok | Err` 通过, `Status | Ok | Err` 拒绝]
- [使解析行为与 `grammar.peg` 中 `EnumDef <- Ident WS* Assign ...` 规则保持一致]

## [2026-03-06 00:10:30] [type: fix] [scope: parser,grammar.peg]
- [将 API 语法改为仅接受头部注释, `ApiDef` 去除尾部 `InlineComment`]
- [解析阶段新增 API 尾部注释拒绝逻辑, 发现同一行注释时直接报错]
- [新增 parser 用例覆盖: API 头部多行注释通过, API 尾部注释失败]

## [2026-03-05 23:58:54] [type: docs] [scope: grammar.peg]
- [恢复 `BaseType` 为框架完整基础类型集合, 不再仅限 `aaa.sb` 示例中出现的子集]
- [新增 `i8/i16/i32/i64/u64/f32/f64` 规则分支, 与现有 `u8/u16/u32/bool/text/bin` 合并]
- [本次仅修正语法文档覆盖范围, 不涉及 parser/lexer 逻辑改动]

## [2026-03-05 23:55:24] [type: docs] [scope: grammar.peg]
- [以 `aaa.sb` 为唯一语法来源重整 `grammar.peg`, 移除 `SchemaPure` 与 pure 系列规则]
- [统一规则排版为 `<-` 对齐与逐行 `/` 分支风格, 提升规则一致性]
- [收敛类型与字符串规则到 `aaa.sb` 覆盖范围, 并保留语义阶段约束说明]

## [2026-03-05 23:44:50] [type: docs] [scope: grammar.peg]
- [按统一对齐风格重排 `grammar.peg`, 使 `<-` 纵向对齐并统一空白格式]
- [将含 alternatives 的规则统一改为逐行 `/` 分支写法, 贴合约定示例风格]
- [将过长规则拆分为可读块并保持语义不变, 仅做格式层改动]

## [2026-03-05 23:34:36] [type: docs] [scope: grammar.peg]
- [新增 `SchemaPure` 纯 PEG 入口, 与兼容入口 `Schema` 并存]
- [新增 `StructDefPure/EnumDefPure/ApiDefPure/TypeRefPure` 规则, 去除跨定义引用与嵌入字段依赖]
- [将"无法纯 PEG 约束"说明限定到兼容模式, 并补充 `SchemaPure` 约束边界说明]

## [2026-03-05 23:26:58] [type: docs] [scope: grammar.peg]
- [将 `grammar.peg` 全部规则改为单行定义风格, 消除换行 alternatives 写法]
- [新增"无法仅用 PEG 完整约束"示例清单, 明确名称唯一性与跨定义引用等边界]
- [新增"完全 PEG 约束"改造建议, 说明需削减语义能力或采用 PEG+语义分析两阶段]

## [2026-03-05 23:15:10] [type: docs] [scope: grammar.peg]
- [将 `grammar.peg` 中注释与说明文本统一改为中文, 语法规则本体保持不变]
- [保留 PEG 难表达语义约束说明, 改为中文表述以降低阅读门槛]
- [不涉及 lexer/parser 实现变更, 仅文档可读性优化]

## [2026-03-05 23:08:19] [type: docs] [scope: grammar.peg]
- [新增 `grammar.peg`, 给出 `sb` schema 的 PEG 语法定义, 覆盖 enum/struct/api/type/comment/string 等核心结构]
- [规则内容对齐当前 `internal/lexer` 与 `internal/parser` 行为, 便于语法对照与扩展评审]
- [补充 PEG 难表达的行级语义约束注释, 明确 trailing comma 与 inline 绑定等边界]

## [2026-03-05 22:51:40] [type: refactor] [scope: test/run_go_bun_tests.sh]
- [移除 `run_smoke_case` 中已废弃的 `smoke_import.test.ts` 兜底分支, 删除动态写入与执行逻辑]
- [将未知 smoke case 改为显式失败(`unsupported case`), 防止静默走弱校验路径]
- [收敛脚本执行路径到三类已定义 schema, 降低维护与排障复杂度]

## [2026-03-05 22:43:41] [type: chore] [scope: test/harness]
- [删除已废弃的 Bun 手写测试与 Go 测试服务端 harness 文件, 收敛到“TS 测试由生成器产出”路径]
- [清理 `robust_core/robust_size_limits/robust_optional_matrix` 对应 `*_bun.test.ts` 与 `*_go_server_main.go` 冗余资产]
- [同步更新 `test/README.md` 文件清单, 移除已删除条目]

## [2026-03-05 22:37:54] [type: fix] [scope: test/robust_size_limits]
- [修正 `robust_size_limits` Go 用例中错误的 `ReqErr` 断言场景: `level=0` 在当前语义中表示未设置, 非非法请求]
- [将请求错误断言改为 `CallRound(..., nil)` 空结构体入参路径, 与 RPC 参数校验规则保持一致]
- [消除 `TestSizeRPCRoundTripAndRequestGuard` 的伪失败, 提升联测结果可信度]

## [2026-03-05 22:34:53] [type: refactor] [scope: test/run_go_bun_tests.sh]
- [脚本改为直接执行生成产物 `ts/sb/rpc_smoke.test.ts`, 删除 Bun 手工测试文件复制与 Go 测试服务端拉起逻辑]
- [保留 Go 端语义测试作为主正确性验证, TS 侧收敛为自动生成 smoke 校验]
- [降低脚本复杂度与维护成本, 对齐“测试文件由生成器产出”目标]

## [2026-03-05 22:29:52] [type: feat] [scope: tplgen/ts.test]
- [新增模板 `ts.test.tpl`, 生成 `ts/sb/rpc_smoke.test.ts` 以自动产出 Bun 客户端 smoke 测试文件]
- [`TsGenerator` 在生成 `rpc.ts` 后追加生成 `rpc_smoke.test.ts`, 统一输出路径为 `ts/sb`]
- [补充模板回归断言, 防止 `rpc_smoke.test.ts` 生成流程与关键内容回退]

## [2026-03-05 22:12:30] [type: fix] [scope: test/go+bun]
- [将 Bun 联测从 TS 自建服务端切换为 Go 测试服务端, 保持架构为服务端 Go + 客户端 Go/TS]
- [新增 3 个 Go 服务端 harness(`robust_core/size_limits/optional_matrix`), 由脚本后台拉起并注入 `SB_BASE_URL` 给 Bun]
- [修复 Bun 用例中 `const [_, ...]` 造成的 `_` 命名遮蔽与 TDZ 问题, 消除 `Cannot access '_' before initialization`]

## [2026-03-05 21:47:59] [type: chore] [scope: test/robust_optional_matrix]
- [新增 `robust_optional_matrix` Go+Bun 语义联测, 覆盖序列化 round-trip, 空值/presence 语义, 非法枚举拒绝]
- [覆盖 RPC `submit/fetch` 关键健壮性分支: 请求非法参数 `ReqErr`, 响应尾随字节与畸形体 `RespErr`, 正常路径 `Ok`]
- [更新 `run_go_bun_tests.sh` 与 `test/README.md`, 将该 case 从 smoke 升级为双端语义测试并对齐文档]

## [2026-03-05 21:37:54] [type: chore] [scope: test/robust_size_limits]
- [新增 `robust_size_limits` 的 Go 语义联测, 覆盖边界序列化 round-trip, 请求非法枚举拒绝, 响应超限与畸形响应拒绝]
- [新增 `robust_size_limits` 的 Bun 语义联测, 与 Go 侧对齐 `ReqErr/RespErr/Ok` 断言路径]
- [调整 `run_go_bun_tests.sh`: 该 case 改为执行 Go+Bun 语义测试, 非仅编译 smoke]

## [2026-03-05 21:21:40] [type: chore] [scope: test/go+bun]
- [为 `run_smoke_case` 增加 Bun 侧 `smoke_import.test.ts`, 校验 TS 生成代码可导入与客户端可构造]
- [将 `robust_size_limits` 与 `robust_optional_matrix` 从仅 Go 编译 smoke 扩展为 Go+Bun 双端 smoke]
- [前置拦截 TS 生成层面的语法或导出回退, 减少跨语言联测盲区]

## [2026-03-05 20:57:18] [type: chore] [scope: test/harness/bun]
- [在 `robust_core_bun.test.ts` 新增 `rpc transport errors` 用例, 覆盖 `NoConn/Timeout/401->NotAuth`]
- [补充 Bun 测试服启停辅助函数, 保证超时与鉴权分支测试可重复且可回收]
- [对齐 Go 端传输层错误覆盖矩阵, 缩小跨语言联测盲区]

## [2026-03-05 20:49:24] [type: chore] [scope: test/harness/go]
- [在 `robust_core_go_test.go` 新增 `TestRPCTransportErrors`, 覆盖 `RpcNoConn/RpcTimeout/401->RpcNotAuth` 传输分支]
- [增加超时场景断言: `Timeout=1ms` + 服务端延迟响应, 校验客户端超时映射稳定性]
- [补齐网络层错误路径验证, 降低仅覆盖业务负载路径导致的鲁棒性盲区]

## [2026-03-05 20:16:31] [type: fix] [scope: test/go+bun]
- [为 `test/run_go_bun_tests.sh` 增加生成产物静态守卫, 扫描 `api.*.go` 中 nilable 返回类型的错误默认值]
- [当检测到 `result *T` 或 `result []T` 仍生成 `return 0, RpcRespErr` 时立即失败并输出定位行]
- [在 core/smoke 两类 case 生成后统一执行守卫, 将模板回退前置到测试入口阶段阻断]

## [2026-03-05 20:01:36] [type: fix] [scope: tplgen/go.api]
- [修复 Go API 逻辑模板在结构体或列表返回类型下错误生成 `return 0, RpcRespErr` 的问题]
- [按返回类型区分默认值: 可空类型返回 `nil`, 标量维持零值]
- [新增模板回归用例, 覆盖 struct 指针, struct 列表, 标量三种返回默认值分支]

## [2026-03-05 19:47:41] [type: feat] [scope: test/go+bun]
- [新增 `test/*.sb` 联测 schema, 覆盖 core/size-limit/optional-matrix 三类协议场景]
- [新增 Go+Bun 双端集成测试资产, 验证序列化与 RPC 在正确性和鲁棒性上的关键边界]
- [新增 `test/run_go_bun_tests.sh` 一键脚本, 自动完成生成、Go 测试与 Bun 测试执行]

## [2026-03-05 19:10:10] [type: fix] [scope: tplgen/go.rpc,ts.rpc]
- [统一对 Go/TS RPC 客户端 `host/baseURL` 执行尾部斜杠收敛, 避免请求路径出现双斜杠]
- [Go `NewClient` 增加 `strings.TrimRight(baseURL, \"/\")`, TS 构造阶段归一 `config.host`]
- [新增模板回归断言, 防止 URL 归一逻辑回退]

## [2026-03-05 19:07:12] [type: fix] [scope: tplgen/go.struct]
- [将 Go `SetStruct(nil)` 从静默成功改为显式报错, 返回 `SetXxx: nil value`]
- [阻断 struct 列表中 `nil` 元素经 `setList` 编码时生成畸形 payload 的风险]
- [新增模板回归断言, 防止 `SetStruct` 空指针校验逻辑回退]

## [2026-03-05 19:00:19] [type: fix] [scope: tplgen/ts.rpc]
- [为 TS RPC 有返回值分支增加响应尾随字节校验, 解码后 `respBuf.len !== 0` 时返回 `RpcErrCode.RespErr`]
- [将结果解码缓冲改为复用 `respBuf`, 避免仅解析前缀数据后静默放行异常尾部]
- [新增模板回归断言, 防止返回值分支尾随字节校验逻辑回退]

## [2026-03-05 18:54:54] [type: fix] [scope: tplgen/go.rpc]
- [为 Go RPC 非列表 struct 入参增加 `nil` 前置校验, 空指针直接返回 `RpcReqErr`]
- [将参数防御前移到请求编码前, 避免无效 `nil` 结构体请求进入网络链路]
- [新增模板回归用例, 防止 struct 入参空指针校验逻辑回退]

## [2026-03-05 18:47:48] [type: fix] [scope: tplgen/ts.rpc]
- [为 TS RPC `nil` 返回类型增加非空响应体校验, `bytes.byteLength !== 0` 时返回 `RpcErrCode.RespErr`]
- [收敛 `200 + 非空body` 的异常响应行为, 避免 `nil` 返回接口被静默判定成功]
- [新增模板回归断言, 防止 `nil` 返回响应校验逻辑回退]

## [2026-03-05 18:39:51] [type: fix] [scope: tplgen/go.rpc]
- [为 Go RPC `nil` 返回类型增加非空响应体校验, `len(body) != 0` 时返回 `RpcRespErr`]
- [调整 `CallXxx` 模板保留 `body` 变量, 防止 `nil` 返回接口静默放行异常 payload]
- [新增模板回归断言, 防止 `nil` 返回响应校验逻辑回退]

## [2026-03-05 18:36:47] [type: fix] [scope: tplgen/go.rpc]
- [为 Go RPC 非空响应解码增加尾随字节校验, 若缓冲区未完全消费则返回 `RpcRespErr`]
- [将响应解码缓冲改为显式 `respBuf`, 解码后统一执行 `respBuf.Len() == 0` 断言]
- [新增模板回归断言, 防止响应尾随字节校验逻辑回退]

## [2026-03-05 18:30:14] [type: fix] [scope: tplgen/go.rpc]
- [为 Go RPC 客户端增加枚举入参本地校验, 单值与列表参数均在发请求前执行 `IsEnum` 校验]
- [非法枚举入参直接返回 `RpcReqErr`, 避免无效请求进入网络链路]
- [新增模板回归用例, 防止枚举入参校验逻辑回退]

## [2026-03-05 18:26:11] [type: fix] [scope: tplgen/ts.rpc]
- [为 TS RPC `timeout` 增加上限收敛 `2147483647ms`, 避免超大值触发定时器异常折叠]
- [构造阶段将 `timeout` 归一为 `Math.min(Math.floor(cfgTimeout), maxTimeoutMs)`]
- [新增模板回归断言, 防止 `timeout` 上限收敛逻辑回退]

## [2026-03-05 18:21:10] [type: fix] [scope: tplgen/go.struct]
- [将 Go struct 列表校验中的 `nil` 元素从静默跳过改为显式报错]
- [阻断 `[]*Struct{nil}` 在编码链路中的脏数据透传, 提升请求数据一致性]
- [新增模板回归断言, 防止 `nil item` 校验逻辑回退]

## [2026-03-05 18:17:55] [type: fix] [scope: tplgen/go.rpc]
- [为 Go RPC `doClient` 增加 `nil context` 兜底, 统一回落到 `context.Background()`]
- [避免调用方误传空上下文时在 `ctx.Done()`/请求构建路径触发异常]
- [新增模板回归断言, 防止 `nil context` 防御逻辑回退]

## [2026-03-05 18:12:25] [type: fix] [scope: tplgen/go.rpc]
- [修复 Go RPC `doClient` 在并发请求下对 `c.HTTP` 的共享写入风险]
- [将 HTTP 客户端兜底改为局部变量 `httpClient`, 避免请求路径修改共享状态]
- [新增模板回归断言, 防止 `c.HTTP` 回写逻辑回退]

## [2026-03-05 18:01:46] [type: fix] [scope: tplgen/struct]
- [收紧 Go/TS struct 单枚举字段解码校验: 仅当 presence-bit 置位时必须通过 `IsEnum`]
- [移除 TS struct 解码对单枚举值 `0` 的豁免, 与 RPC 严格校验语义保持一致]
- [新增模板回归断言, 防止已置位单枚举字段校验逻辑回退]

## [2026-03-05 17:52:07] [type: fix] [scope: tplgen/go.enum]
- [为 `go.enum.tpl` 增加条件导入, 仅在存在枚举时生成 `bytes/slices/unsafe` 导入块]
- [修复空枚举 schema 下 `enum.go` 可能触发未使用导入导致编译失败的问题]
- [新增模板回归用例, 防止空枚举导入逻辑回退]

## [2026-03-05 17:46:36] [type: fix] [scope: tplgen/ts.rpc]
- [将 TS RPC `timeout` 配置归一为有限非负整数, 非法值回落默认 5000ms]
- [支持 `timeout=0` 显式禁用请求超时定时器, 避免被 `||` 误改为 5000ms]
- [新增模板回归断言, 防止超时归一和定时器启停逻辑回退]

## [2026-03-05 17:41:17] [type: fix] [scope: tplgen/type.go]
- [修复 `GetBin` 长度校验的类型截断风险, 将 `uint16(buf.Len()) < l` 改为 `buf.Len() < int(l)`]
- [避免缓冲区长度超过 65535 时发生误判, 提升大包场景解码正确性]
- [新增模板回归断言, 防止长度校验逻辑回退]

## [2026-03-05 17:28:28] [type: fix] [scope: tplgen/ts.rpc]
- [移除 TS RPC 单枚举响应值对 `0` 的豁免, 统一改为严格枚举合法性校验]
- [当响应枚举值非法时统一返回 `RpcErrCode.RespErr`, 与 Go 客户端行为对齐]
- [新增模板回归断言, 防止单枚举响应校验逻辑回退]

## [2026-03-05 17:12:26] [type: fix] [scope: tplgen/go.rpc]
- [为 Go RPC 枚举返回值增加合法性校验, 覆盖单值和列表结果]
- [当服务端返回非法枚举值时, 客户端统一返回 `RpcRespErr`]
- [新增模板回归用例, 防止枚举返回校验逻辑回退]

## [2026-03-05 16:56:09] [type: fix] [scope: tplgen/ts.rpc]
- [将 TS RPC 非列表枚举参数校验改为严格模式, 移除对 `0` 的豁免]
- [对齐 Go 服务端枚举参数校验规则, 提前返回 `ReqErr`]
- [新增模板回归断言, 防止枚举参数校验逻辑回退]

## [2026-03-05 16:49:57] [type: fix] [scope: tplgen/ts.rpc]
- [将 TS `retries` 配置归一为有限非负整数, 非法值回落默认 3]
- [构造函数改为 `cfgRetries` 显式校验并执行 `Math.floor`]
- [新增模板回归断言, 防止重试数值归一逻辑回退]

## [2026-03-05 16:44:06] [type: fix] [scope: tplgen/go.api]
- [为无参数 API 增加空请求体校验, 禁止携带多余 payload]
- [`parseRequest` 在 `len(args)==0` 分支读取至多 1 字节并严格判空]
- [新增模板回归断言, 防止无参 body 校验逻辑回退]

## [2026-03-05 16:35:06] [type: fix] [scope: tplgen/go.api]
- [为 `sendResponse` 增加默认 `Content-Type: application/octet-stream` 响应头]
- [补充 `w.Write` 错误分支处理, 失败时返回 `500`]
- [新增模板回归断言, 防止响应头和写失败处理逻辑回退]

## [2026-03-05 16:11:04] [type: fix] [scope: tplgen/go.rpc]
- [为 Go 客户端请求默认补充 `Content-Type: application/octet-stream`]
- [仅在调用方未显式设置 `Content-Type` 时生效, 降低行为变更风险]
- [新增模板回归断言, 防止请求头默认值逻辑回退]

## [2026-03-05 16:04:51] [type: fix] [scope: tplgen/go.api]
- [为 `parseRequest` 增加尾随字节校验, 避免前缀可解析请求被误接受]
- [将请求解码缓冲改为显式 `buf` 复用, 解码后要求 `buf.Len() == 0`]
- [新增模板回归断言, 防止尾随字节校验逻辑回退]

## [2026-03-05 15:58:22] [type: fix] [scope: tplgen/ts.rpc]
- [为 TS `maxRespBytes` 增加有限值校验, 拒绝 `Infinity/NaN` 等非法配置]
- [对有效值执行 `Math.floor` 和 `Number.MAX_SAFE_INTEGER` 上限收敛]
- [新增模板回归断言, 防止数值归一逻辑回退]

## [2026-03-05 15:33:44] [type: fix] [scope: tplgen/ts.rpc]
- [将 TS 客户端 `retries` 负值归一为 0, 统一与 Go 重试边界语义]
- [避免非法重试配置导致 `_fetch` 循环不执行的隐式分支]
- [新增模板回归断言, 防止重试下限逻辑回退]

## [2026-03-05 15:25:17] [type: fix] [scope: tplgen/ts.rpc]
- [为 TS RpcConfig 增加 `maxRespBytes` 并设置默认值 4MB]
- [在 TS `_fetch` 中增加响应大小校验, 超限返回 `RpcErrCode.RespErr`]
- [新增模板回归断言, 防止响应上限逻辑回退]

## [2026-03-05 15:20:04] [type: fix] [scope: tplgen/go.api]
- [为 Go 服务端 `parseRequest` 增加请求体默认上限 4MB]
- [请求解析改为 `io.LimitReader` 限流读取, 超限返回 `413 Request Entity Too Large`]
- [新增模板回归断言, 防止请求体上限逻辑回退]

## [2026-03-05 15:15:49] [type: fix] [scope: tplgen/go.rpc]
- [为 Go RPC 客户端增加 `MaxRespBytes` 响应大小上限配置, 默认 4MB]
- [将响应读取改为 `io.LimitReader` 限流读取, 超限返回 `RpcRespErr`]
- [新增模板回归断言, 防止响应上限逻辑回退]

## [2026-03-05 15:02:53] [type: fix] [scope: tplgen/go.rpc]
- [为 Go Client 头部读写增加 `sync.RWMutex` 保护, 消除并发读写 map 风险]
- [发送请求前复制 header 快照, 避免持锁跨网络调用]
- [新增模板回归断言, 防止并发保护逻辑回退]

## [2026-03-05 14:53:41] [type: fix] [scope: tplgen/go.rpc]
- [为 Go RPC 重试次数增加下限保护, 将负值归一为 0]
- [补充 `resp == nil` 防御分支, 避免非法配置路径触发空指针]
- [新增模板回归断言, 防止重试边界与空指针防御回退]

## [2026-03-05 14:38:42] [type: fix] [scope: tplgen/ts.struct]
- [将 TS 非列表 struct 字段改为 `Type | null`, 与 Go 指针可空语义对齐]
- [将 TS 非列表 struct 字段默认值改为 `null`, 并在编码时仅对非空字段置 presence-bit]
- [新增模板回归测试, 防止嵌套 struct 可空语义回退]

## [2026-03-05 14:34:11] [type: fix] [scope: tplgen/go.struct]
- [移除 Go struct 解码对空 body 的静默成功分支, 统一按 bitmask 严格校验]
- [新增模板回归断言, 防止 `if buf.Len() == 0 { return nil }` 回归]
