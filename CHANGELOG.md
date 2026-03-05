# Changelog

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
