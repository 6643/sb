# Demo Rules

这个目录现在只用于放规则样例, 不再保存生成后的 Go / TypeScript 产物.

## 文件

- `all_rules.sb`: 当前 `.sb` 语法和主要结构约束的正例总览

## 覆盖范围

`all_rules.sb` 覆盖这些正例:

- 文件级注释
- 枚举
- 显式枚举值
- 混合显式 / 隐式枚举值
- 结构体字段
- 字段注释
- 字段 tag
- 嵌入字段 `...TypeName`
- 基础类型
- `text`
- `bin`
- `[T]` 列表
- 单段 API 名称
- 双段 API 名称
- `nil` 返回值

## 不在本文件中表达的约束

以下规则主要通过 parser / lexer / 测试表达, 不适合写进一个必须可解析的正例文件:

- 尾部注释非法
- `[nil]` 非法
- `nil` 不能作为参数或字段类型
- API 名称超过两段非法
- 枚举成员之间缺少 `|`
- struct 字段必须独占一行
- 逗号必须独占一行

这些边界以 `internal/lexer_scanner.go`、`internal/parser_schema.go` 和对应测试为准.
