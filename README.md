# SB (Simple Binary) 代码生成器

SB 是一个面向 Go 和 TypeScript 的二进制协议与 RPC 代码生成工具. 输入是自定义的 `.sb` 描述文件, 输出是对等的类型定义, 编解码逻辑, RPC 客户端, Go API handler 和 Markdown 文档.

## 1. 项目定位

- 目标: 用固定 schema 生成强类型二进制协议与 RPC 代码.
- 优先级: 正确性 > 可维护性 > 性能 > 开发速度.
- 当前实现: 生成主链路已统一到 `sb/internal` 根包, 代码生成不再依赖模板文件.
- 适用场景: 内部服务通信, 前后端对等协议, 小到中等体积的结构化二进制数据.

## 2. 快速开始

### 2.1 运行命令

确保本机已安装 Go, 然后在项目根目录执行:

```bash
go run . <input.sb> [flags]
```

示例:

```bash
go run . -go ./go -ts ./ts -tag bson,json aaa.sb
```

### 2.2 命令行参数

- `-go`: Go 代码输出目录, 默认 `./go`
- `-ts`: TypeScript 代码输出目录, 默认 `./ts`
- `-tag`: 为 Go 结构体追加 tag, 例如 `bson,json`

### 2.3 生成目录约定

- `go/sb`: Go 协议类型, 运行时, 编解码, RPC 客户端, API handler, `api.*.go`
- `ts/sb`: TypeScript 类型, 运行时, 编解码, RPC 客户端
- `go/sb/DOC.md`: 自动生成的 Markdown API 文档
- `api.*.go`: 可手改逻辑文件, 生成器使用指纹保护已改文件不被覆盖

## 3. `.sb` 语法

### 3.1 基础类型与长度限制

| 类型 | 说明 | 长度或范围 |
| :--- | :--- | :--- |
| `u8-u64` | 无符号整数 | 8 到 64 位 |
| `i8-i64` | 有符号整数 | 8 到 64 位 |
| `f32`, `f64` | 浮点数 | - |
| `bool` | 布尔值 | - |
| `text` | UTF-8 字符串 | 最大 `65535` 字节, 使用 `u16` 前缀 |
| `bin` | 二进制数据 | 最大 `4294967295` 字节, 使用 `u32` 前缀 |
| `[T]` | 数组或切片 | 最大 `65535` 个元素, 使用 `u16` 前缀 |

说明:

- `text` 的上限按 UTF-8 字节数计算, 不是字符数.
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
- Go 端会生成逻辑接口和 HTTP handler
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
- 路由注册支持可选 middleware, 内置 `TimeoutMiddleware` 可为服务端请求注入超时 deadline
- 推荐服务端采用两层超时: `http.Server` 的连接级超时配合 `TimeoutMiddleware` 的请求级超时
- 运行时提供 `GetBin(buf)` 和 `GetBinView(buf)`
- `GetBinView(buf)` 是可选借用接口, 返回切片只在底层 `buf` 未继续消费或复用前有效

### 5.2 TypeScript

- 不使用 `throw` 作为主错误流
- API 和序列化函数统一返回 `[data, err]`
- 当 `err` 非空时, `data` 仍返回该类型的安全零值
- 生成代码兼容严格模式下的显式类型检查

## 6. 验证与测试

常用命令:

```bash
go test ./...
go run . -go ./go -ts ./ts -tag bson,json aaa.sb
```

关键回归测试入口:

- [lexer_scanner_test.go](/._/tool/sb/internal/lexer_scanner_test.go): 词法规则与非法字符边界
- [parser_schema_test.go](/._/tool/sb/internal/parser_schema_test.go): `.sb` 语法规则, 注释规则, 嵌入规则, 枚举规则
- [semantic_resolver_test.go](/._/tool/sb/internal/semantic_resolver_test.go): 语义展开与枚举值冲突
- [tpl_template_regression_test.go](/._/tool/sb/internal/tpl_template_regression_test.go): Go/TS/DOC 生成回归
- [legacy_go_generator_test.go](/._/tool/sb/internal/legacy_go_generator_test.go): 旧 Go 生成器兼容回归
- [legacy_ts_generator_test.go](/._/tool/sb/internal/legacy_ts_generator_test.go): 旧 TS 生成器兼容回归

## 7. 已知限制

- 当前协议没有 `field id`, 字段顺序变更会直接改变 wire format
- 当前不支持 `map` 类型
- `enum` 不支持成员注释
- `text` 仍按字节长度限制, 不是按字符数限制
- `bin` 虽然使用 `u32` 长度前缀, 但默认 RPC 包体上限仍受运行时配置约束

## 8. 维护约定

- 结构变更, 协议变更, 生成器修复都必须更新 [CHANGELOG.md](/._/tool/sb/CHANGELOG.md)
- 新记录插入 `CHANGELOG.md` 标题后, 禁止追加到底部
- 主文档以本文件为准, 不再维护独立的 PEG 说明文档
