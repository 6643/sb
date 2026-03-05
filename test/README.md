# Test Suite

该目录提供面向强正确性和鲁棒性的联测资产, 覆盖序列化与 RPC 两条链路.

## 目录结构

- `robust_core.sb`: 主联测 schema, 覆盖 enum/struct/nil 返回/尾随字节/非法枚举等场景.
- `robust_size_limits.sb`: 大小边界 schema, 配套 Go+Bun 语义联测.
- `robust_optional_matrix.sb`: 可空与嵌套矩阵 schema, 配套 Go+Bun 语义联测.
- `harness/robust_core_go_test.go`: Go 端集成测试.
- `harness/robust_size_limits_go_test.go`: `robust_size_limits` Go 语义测试.
- `harness/robust_optional_matrix_go_test.go`: `robust_optional_matrix` Go 语义测试.
- `run_go_bun_tests.sh`: 一键生成并执行 Go+Bun 联测.

## 运行方式

在仓库根目录执行:

```bash
./test/run_go_bun_tests.sh
```

脚本行为:

1. 用 `go run .` 基于 `test/*.sb` 生成临时代码到 `test/.work`.
2. 对 `robust_core.sb` 同时执行:
   - Go 集成测试: `go test ./...`
   - Bun 自动生成 smoke: `bun test rpc_smoke.test.ts`
3. 对 `robust_size_limits.sb` 与 `robust_optional_matrix.sb` 执行:
   - Go 语义测试: `go test ./...`
   - Bun 自动生成 smoke: `bun test rpc_smoke.test.ts`

## 覆盖重点

- 序列化 round-trip 正确性.
- 非法枚举值拒绝.
- struct `nil` 入参拒绝.
- `nil` 返回接口非空 body 拒绝.
- 返回尾随垃圾字节拒绝.
- 响应体上限保护.
- RPC 请求前参数合法性校验.
- 可空 struct 字段 presence 语义.
- `nil` 返回接口空 body 约束.
- TS 侧仅作为客户端 smoke 验证, 测试文件由生成器自动产出.
