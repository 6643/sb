# PEG 一致性清单

## 目标
- 记录 `grammar.peg` 与 `lexer/parser` 的已对齐规则.
- 明确仍需语义阶段处理的约束边界.
- 提供关键回归测试索引, 便于后续回归与评审.

## 已对齐规则
- 标识符: 仅 ASCII, `[A-Za-z_][A-Za-z0-9_]*`.
- 空白: 行内空白仅 `space/tab`, 换行仅 `\n` 或 `\r\n`.
- 字符串: 仅支持双引号字符串.
- 枚举: 定义必须包含 `=`.
- 枚举: 成员分隔仅允许 `|`.
- 枚举: 显式值范围 `0-255`, 禁止负数和前导零.
- 结构体: 字段必须一行一个.
- 结构体: 逗号必须独占一行.
- 结构体: `tag` 仅可出现在字段类型之后.
- API: 名称仅支持 `a` 或 `a.b`.
- API: 注释仅允许头部注释, 不允许尾部注释.
- 类型: `nil` 仅允许作为 API 返回类型.
- 类型: `[nil]` 非法.

## 仍需语义阶段处理
- 名称唯一性: 结构体字段名, 枚举成员名, API 名称冲突检查.
- 类型存在性: 自定义类型是否已定义.
- 枚举值冲突: 显式值重复, 显式值与隐式分配冲突.
- 跨定义依赖关系: 引用图完整性与循环依赖策略.

## 关键回归测试索引
- `internal/lexer/lexer_test.go`: `TestBacktickStringRejected`.
- `internal/lexer/lexer_test.go`: `TestStandaloneCRRejected`.
- `internal/lexer/lexer_test.go`: `TestNonASCIIIdentRejected`.
- `internal/lexer/lexer_test.go`: `TestNonASCIISpaceNotSkipped`.
- `internal/lexer/lexer_test.go`: `TestNegativeNumberRejectedAtLexer`.
- `internal/lexer/lexer_test.go`: `TestNonASCIIDigitRejected`.
- `internal/parser/parser_test.go`: `TestParseAPITailCommentRejected`.
- `internal/parser/parser_test.go`: `TestParseAPINamedThreeSegmentsRejected`.
- `internal/parser/parser_test.go`: `TestParseStructFieldsMustBeOnePerLine`.
- `internal/parser/parser_test.go`: `TestParseStructFieldLineTrailingCommaRejected`.
- `internal/parser/parser_test.go`: `TestParseStructStandaloneCommaLineAccepted`.
- `internal/parser/parser_test.go`: `TestParseAPIArgNilRejected`.
- `internal/parser/parser_test.go`: `TestParseFieldListNilRejected`.
- `internal/parser/parser_test.go`: `TestParseFieldNilRejected`.
- `internal/parser/parser_test.go`: `TestParseStructCommentLineNotInlineFieldNote`.
- `internal/parser/parser_test.go`: `TestParseEnumNegativeValueRejected`.
- `internal/parser/parser_test.go`: `TestParseEnumValueOverflowRejected`.
- `internal/parser/parser_test.go`: `TestParseEnumValueLeadingZeroRejected`.
- `internal/parser/parser_test.go`: `TestParseCommentLineWithoutTrailingNewlineRejected`.
- `internal/parser/parser_test.go`: `TestParseStructCommentLineWithoutTrailingNewlineRejected`.

## 验证现状
- 当前工作环境缺少 `go` 可执行程序, 无法在本地直接运行 `go test`.
- 仓库内 `go/` 目录为生成代码输出目录, 不包含 Go 工具链二进制.
- 在补齐 Go 运行环境前, 本文档中的一致性结论基于规则对照与静态变更审查.
