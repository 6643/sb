# Changelog

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
