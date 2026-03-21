# sb

sb 是一个面向 Go 和 TypeScript 的二进制协议与 RPC 代码生成工具. 输入是自定义的 `.sb` 描述文件, 输出是对等的类型定义, 编解码逻辑, RPC 客户端, Go API handler 和 Markdown 文档.

## 1. 项目定位

- 目标: 用固定 schema 生成强类型二进制协议与 RPC 代码.
- 优先级: 正确性 > 可维护性 > 性能 > 开发速度.
- 当前实现: 生成主链路已统一到 `sb/internal` 根包, 代码生成不再依赖模板文件.
- 适用场景: 内部服务通信, 前后端对等协议, 小到中等体积的结构化二进制数据.

### 协议规范

- README 只保留项目概览, 协议细节与 wire 语义以 [PROTOCOL.md](PROTOCOL.md) 为准.
- 旧版生成链路已删除, 当前仓库只保留这套紧凑编码协议实现.

### 1.1 与常见序列化的优缺点对比

sb 当前生成的序列化格式有几个核心特征:

- 固定 schema, 通过代码生成完成 Go 与 TypeScript 的强类型编解码.
- 数值类型使用固定宽度的小端编码.
- `text`、`bin` 和 `[T]` 使用长度前缀, 不携带字段名.
- 结构体字段顺序就是 wire format, 当前没有 `field id`.

如果把问题收敛成“为什么选 `sb`, 而不是其他常见序列化”, 核心取舍其实很直接:

- `sb` 优先的是: 固定协议下的简单性, 生成代码的直接性, 以及 Go/TS 双端行为的一致性.
- `sb` 放弃的是: 自描述能力, 强兼容演进机制, 以及超大生态下的通用工具链红利.

## 2. 快速开始

### 2.1 运行命令

确保本机已安装 Go, 然后在项目根目录执行:

```bash
go run . <input.sb> [flags]
```

示例:

```bash
go build -o sb
cd demo
../sb -go ./go -ts ./ts -tag bson,json aaa.sb
```

### 2.2 命令行参数

- `-go`: Go 代码输出目录, 默认 `./go`
- `-ts`: TypeScript 代码输出目录, 默认 `./ts`
- `-tag`: 为 Go 结构体追加 tag, 例如 `bson,json`

### 2.3 生成目录约定

- `go/sb`: Go 协议类型, 编解码, RPC 客户端, API handler, `api.*.go`
- `ts/sb`: TypeScript 类型, 运行时, 编解码, RPC 客户端
- `go/sb/runtime.go`: Go runtime
- `go/sb/DOC.md`, `ts/sb/DOC.md`: 自动生成的 Markdown API 文档
- `api.*.go`: 可手改逻辑文件, 生成器使用指纹保护已改文件不被覆盖

## 3. `.sb` 语法

### 3.1 基础类型与长度限制

| 类型 | 说明 | 长度或范围 |
| :--- | :--- | :--- |
| `u8-u64` | 无符号整数 | 8 到 64 位 |
| `i8-i64` | 有符号整数 | 8 到 64 位 |
| `f32`, `f64` | 浮点数 | - |
| `bool` | 布尔值 | - |
| `text` | UTF-8 字符串 | 变长类型, 具体 wire 编码见 `PROTOCOL.md` |
| `bin` | 二进制数据 | 变长类型, 具体 wire 编码见 `PROTOCOL.md` |
| `[T]` | 数组或切片 | 变长类型, 具体 wire 编码见 `PROTOCOL.md` |

说明:

- `text` 的长度按 UTF-8 字节数计算, 不是字符数.
- 默认 RPC 请求和响应上限仍由运行时控制, 当前默认值是 `4MB`.

### 3.2 通用规则

- 标识符只允许 ASCII, 规则是 `[A-Za-z_][A-Za-z0-9_]*`
- 行内空白只支持 `space` 和 `tab`
- 换行只支持 `\n` 或 `\r\n`
- 字符串只支持双引号 `"..."`
- `nil` 只允许作为 API 返回类型
- `[nil]` 非法

### 3.3 注释规则

- 只支持 `//` 单行注释
- 文件级, 枚举级, 结构体级, API 级注释都写在定义前
- 结构体字段注释只允许写在字段前的独立注释行, 支持多行合并
- 枚举只支持定义前整体注释, 不支持成员注释
- API 不支持尾部注释

示例:

```sb
// 账户状态
Status = Offline | Online | Deleted

User {
  // 用户唯一 ID
  // 会映射到 Go tag
  id u32 "_id"
  name text
}
```

### 3.4 枚举

```sb
// 带显式值的枚举
Status = Ok(0) | Err(1) | Forbidden(255)
```

约束:

- 枚举定义必须包含 `=`
- 成员分隔只允许 `|`
- 显式值范围是 `0-255`
- 禁止负数和前导零
- 同一个枚举内, 显式值不能重复

### 3.5 结构体

```sb
User {
  id u32
  name text
}

Admin {
  ...User
  role text
}
```

约束:

- 字段必须一行一个
- 逗号必须独占一行
- `tag` 只允许写在字段类型之后
- 嵌入字段必须显式写为 `...TypeName`
- 结构体字段顺序就是当前协议编码顺序

### 3.6 API

```sb
// 命名空间.方法名(参数) => 返回类型
user.get_info(id u32) => UserInfo
set_flag(enabled bool) => nil
```

约束:

- API 名称只支持 `a` 或 `a.b`
- 返回 `nil` 表示无 body 返回
- Go 端会生成逻辑 stub 和 HTTP handler
- TypeScript 端会生成对应的 RPC 调用方法

## 4. 规则边界

### 4.1 语法阶段已覆盖规则

当前 `lexer + parser` 已对齐的核心规则如下:

- 标识符, 空白, 换行, 双引号字符串
- 枚举必须包含 `=`
- 枚举成员分隔只允许 `|`
- 枚举显式值范围 `0-255`, 禁止负数和前导零
- 枚举不支持成员注释
- 结构体字段必须一行一个
- 结构体逗号必须独占一行
- 结构体嵌入必须写成 `...TypeName`
- 结构体字段注释必须写在字段前
- API 名称只支持 `a` 或 `a.b`
- API 不支持尾部注释
- `nil` 只允许作为 API 返回类型
- `[nil]` 非法

### 4.2 仍在语义阶段处理的规则

以下规则不属于 PEG 或 parser 语法本身, 由语义阶段处理:

- 名称唯一性: 结构体字段名, 枚举成员名, API 名称冲突
- 类型存在性: 自定义类型是否已声明
- 枚举值冲突: 显式值重复, 显式值与隐式分配冲突
- 跨定义依赖关系: 引用完整性和循环依赖策略

## 5. 生成约定

### 5.1 Go

- 错误类型使用原生 `error`
- 生成 RPC handler, 请求解析和响应编码逻辑
- 路由注册支持可选 middleware, 内置 `TimeoutMiddleware` 会在超时后返回 HTTP `408`, 同时把 deadline 传递给业务逻辑
- 推荐服务端采用两层超时: `http.Server` 的连接级超时配合 `TimeoutMiddleware` 的请求级超时
- Go 客户端 `Client.Timeout` 覆盖单次请求的完整 attempt, 包括发请求和读取响应体
- Go 客户端默认关闭自动重试; 只有显式开启 `EnableRetries=true` 且请求具备幂等保障时才会按 `Retries` 执行重试
- handler 返回 `RpcNoConn`、`RpcRespErr` 或未知内部状态时, 会统一映射为合法 HTTP `500`
- Go runtime 提供基础 header、定宽类型和紧凑变长字段读写 helper

### 5.2 TypeScript

- 不使用 `throw` 作为主错误流
- 解码函数返回 `[data, err]`, 编码函数返回 `err`, RPC 调用返回 `[data, status]` 或 `status`
- 当 `err` 非空时, `data` 仍返回该类型的安全零值
- TS 客户端 `maxRespBytes` 是传输中限制, chunked 响应超限会在读取过程中立即失败
- TS RPC 已知状态保留 `RpcErrCode`, 未知 HTTP 状态可能直接透传为 `number`
- 生成代码兼容严格模式下的显式类型检查

## 6. 验证与测试

常用命令:

```bash
go test ./...
go run . -go ./go -ts ./ts -tag bson,json aaa.sb
```

关键回归测试入口:

- [lexer_scanner_test.go](internal/lexer_scanner_test.go): 词法规则与非法字符边界
- [parser_schema_test.go](internal/parser_schema_test.go): `.sb` 语法规则, 注释规则, 嵌入规则, 枚举规则
- [semantic_resolver_test.go](internal/semantic_resolver_test.go): 语义展开与枚举值冲突
- [tpl_go_render_test.go](internal/tpl_go_render_test.go): Go 生成代码回归
- [tpl_ts_render_test.go](internal/tpl_ts_render_test.go): TS 生成代码回归
- [go/sb/runtime_test.go](go/sb/runtime_test.go): Go runtime 协议细节与 canonical 规则
- [go/sb/cross_consistency_test.go](go/sb/cross_consistency_test.go): Go 端协议一致性回归
- [ts/sb/runtime.test.ts](ts/sb/runtime.test.ts): TS runtime 协议细节与 canonical 规则

## 7. 已知限制

- 当前协议没有 `field id`, 字段顺序变更会直接改变 wire format
- 当前不支持 `map` 类型
- `enum` 不支持成员注释
- `text` 仍按字节长度限制, 不是按字符数限制
- `bin` 等变长字段的具体长度状态与上限以 `PROTOCOL.md` 和 runtime 实现为准, 默认 RPC 包体上限仍受运行时配置约束

## 8. 维护约定

- 结构变更, 协议变更, 生成器修复都必须更新 [CHANGELOG.md](CHANGELOG.md)
- 新记录插入 `CHANGELOG.md` 标题后, 禁止追加到底部
- 协议细节以 `PROTOCOL.md` 为准, README 只保留使用与维护概览
