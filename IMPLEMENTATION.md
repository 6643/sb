# Implementation

本文档描述当前仓库实现落点、目录结构和验证入口.  
协议规则、wire format 和规范编码要求以 [PROTOCOL.md](PROTOCOL.md) 为准.  
本文件是实现说明, 不构成协议约束.

## 1. 目标

- 说明当前仓库如何承载 `sb` 协议实现.
- 说明哪些目录是源码, 哪些目录是示例生成结果.
- 给出当前主链路的验证入口.

## 2. 目录结构

### 2.1 手写源码

- `main.go`: CLI 入口
- `internal/`: 解析器、语义分析、生成器模板和相关测试
- `grammar.peg`: 语法设计参考文档, 不参与当前解析实现
- `PROTOCOL.md`: 协议规范
- `README.md`: 使用概览

### 2.2 仓库内示例生成结果

- `demo/all_rules.sb`: 仓库内保留的规则总览样例
- `demo/README.md`: 规则样例覆盖说明

### 2.3 仓库内 fixture

- `fixtures/go/sb/`: Go fixture 与测试资产
- `fixtures/ts/sb/`: TypeScript fixture 与测试资产

### 2.4 关键实现落点

- Go runtime: `fixtures/go/sb/runtime.go`
- TypeScript runtime: `fixtures/ts/sb/type.ts`
- Go 生成器: `internal/tpl_go_render.go`
- TypeScript 生成器: `internal/tpl_ts_render.go`
- 文档生成器: `internal/tpl_doc_render.go`

## 3. 生成链路

当前生成命令:

```bash
go run . -go ./out/go -ts ./out/ts -tag bson,json ./demo/all_rules.sb
```

说明:

- 生成器输出目录由 `-go` 和 `-ts` 参数决定.
- `demo/` 只用于规则样例说明.
- 当前仓库保留的生成产物与测试资产位于 `fixtures/`.
- 根目录不再保留额外的 `go/`、`ts/`、`aaa.sb` fixture.

## 4. 验证入口

### 4.1 Go

- `go test ./internal`
- `go test ./fixtures/go/sb`

### 4.2 TypeScript

- `bun test fixtures/ts/sb`

### 4.3 关键回归

- `fixtures/go/sb/runtime_test.go`: Go runtime 协议细节与 canonical 规则
- `fixtures/go/sb/cross_consistency_test.go`: Go 端双端一致性回归
- `fixtures/ts/sb/runtime.test.ts`: TS runtime 协议细节与 canonical 规则
- `fixtures/ts/sb/cross_consistency.test.ts`: TS 侧驱动的双端一致性回归

## 5. 文档边界

- [PROTOCOL.md](PROTOCOL.md) 只描述协议规则和规范编码要求.
- 本文档只描述当前仓库中的实现映射、目录落点和验证入口.
- 当目录结构、测试入口、生成命令变化时, 优先更新本文档.

## 6. 语法真相源

- `grammar.peg` 仅作设计参考, 不参与当前解析实现.
- 当前语法真相源是:
  - `internal/lexer_scanner.go`
  - `internal/parser_schema.go`
  - `internal/parser_schema_test.go`
- 当语法边界发生变化时, 应先更新实现与测试, 再决定是否同步更新 `grammar.peg`.
