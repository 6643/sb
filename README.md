# SB (Simple Binary) 代码生成器

SB 是一个高性能的二进制序列化协议和 RPC 代码生成工具。它通过解析自定义的 `.sb` 描述文件，为 **Go** 和 **TypeScript** 生成对等的、强类型的编解码逻辑和 API 调用接口。

## 1. 核心特性

*   **极致性能**: 采用紧凑的二进制格式进行序列化，比 JSON 更小、更快。
*   **统一的错误处理**: 
    *   **TypeScript 无异常设计**: 完全移除 `throw`，采用与 Go 一致的 `[data, err]` 多返回值模式。
    *   **HTTP 状态码驱动**: 废弃泛型包裹类，直接通过标准的 HTTP 状态码映射业务错误。
*   **语义化继承**: 结构体支持嵌入（Embedding）方式实现字段复用，保持代码简洁。
*   **现代化 TS 规范**: 全量采用箭头函数，Buffer 操作配合严格的边界检查，确保前端环境下的安全性。
*   **文档自动生成**: 在生成代码的同时，自动输出 Markdown 格式的 API 文档。

## 2. 快速开始

### 安装与运行
确保你已经安装了 Go 环境，在项目根目录下执行：

```bash
go run . <input.sb> [flags]
```

### 命令行参数
*   `-go`: Go 代码输出目录（默认 `./go`）。
*   `-ts`: TypeScript 代码输出目录（默认 `./ts`）。
*   `-tag`: 为 Go 结构体生成的额外 Tag（例如 `bson,json`）。

**示例命令：**
```bash
go run . -go ./go -ts ./ts -tag bson,json aaa.sb
```

## 3. `.sb` 语法规范

### 3.1 基础类型与长度限制
| 类型 | 说明 | 长度/范围限制 |
| :--- | :--- | :--- |
| `u8-u64` | 无符号整数 | 8 到 64 位 |
| `i8-i64` | 有符号整数 | 8 到 64 位 |
| `f32, f64` | 浮点数 | - |
| `bool` | 布尔值 | - |
| `text` | 字符串 | 最大 **65535** 字节 (u16) |
| `bin` | 二进制数据 | 最大 **65535** 字节 (u16) |
| `[T]` | 数组/切片 | 最大 **255** 个元素 (u8) |

### 3.2 枚举 (Enums)
支持简单的枚举值或带指定数值的变体：
```sb
// 账户状态
AccountStatus = Offline | Online | Deleted

// 带数值的错误码定义
Status = Ok(0) | Err(1) | Forbidden(403)
```

### 3.3 结构体 (Structs)
支持字段注释、Tag 定义以及结构体嵌入：
```sb
User {
    id   u32 "_id" // 映射到 Go 的 tag
    name text
}

// 继承 User 字段
Admin {
    User
    role text
}
```

### 3.4 API 定义
API 支持命名空间，并映射为不同语言的 Handler 或 Method：
```sb
// 命名空间.方法名(参数) => 返回类型
user.get_info(id u32) => UserInfo
```
*   返回 `nil` 表示无 Body 返回。
*   Go 端会生成逻辑接口 `user_get_info` 和 HTTP 处理函数 `UserGetInfo`。

## 4. 跨语言开发规范

### Go 语言
*   **标准错误**: 返回原生的 `error` 接口。
*   **自动化 Handler**: 生成的 RPC 代码会自动处理参数的反序列化和结果的序列化。

### TypeScript 语言
*   **[data, err] 模式**: 调用任何 API 或序列化函数都返回元组。
    ```typescript
    const [info, err] = await api.user_get_info(123);
    if (err) {
        console.error("请求失败", err.code);
        return;
    }
    console.log(info.name);
    ```
*   **零值保证**: 当 `err` 不为空时，`data` 永远是该类型的安全零值（如 `0`, `""`, `[]`）。
