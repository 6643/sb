# SB 协议规范

> 状态: 当前默认协议实现, Go / TypeScript 生成代码落地到 `go/sb` 与 `ts/sb`, Go runtime 位于 `go/sb/rt`  
> 兼容性: 与仓库早期实现不兼容, 旧版生成链路已删除  
> 目的: 描述当前默认生成器的 wire format 与约束

## 1. 范围

本草案只讨论结构体与基础类型的二进制编码规则, 不讨论:

- 向后兼容
- pointer / offset / 随机访问
- 部分读取
- 新增类型系统能力

约束前提:

- 永远完整编解码
- 仍然是固定 schema
- 结构字段顺序仍然由 schema 决定
- 变长字段放在尾部顺序编码, 不引入 pointer

## 2. 目标

- 对零值字段做到“尽量不传”
- 对 `text` / `bin` / `list` 的短长度做更紧凑编码
- 保留 `SB` 当前“顺序读写、实现简单、生成代码直接”的特点

## 3. 非目标

- 不追求自描述
- 不支持未知字段保留
- 不引入字段编号
- 不做 FlatBuffers 式 offset 访问

## 4. 默认值

字段默认值定义如下:

- `bool`: `false`
- `i8/u8/i16/u16/i32/u32/i64/u64`: `0`
- `f32/f64`: `0.0`
- `enum`: schema 中第一个枚举项
- `text`: `""`
- `bin`: 空字节串
- `[T]`: 空列表
- `struct`: 所有字段都为默认值

说明:

- `f32/f64` 接受 `-0.0` 在编码时被视为零值并被规范化为 `0.0`
- Go 中 `nil slice` 与空 slice 在 wire 上等价

## 5. 总体布局

一个结构体编码为:

```text
header + fixed block + var tail
```

其中:

- `header`: 逐字段写入 tag bit
- `fixed block`: 所有“已出现的定宽字段”按 schema 顺序写入
- `var tail`: 所有“非空变长字段”按 schema 顺序写入

不使用 pointer / offset.

## 6. 字段分类与 tag 宽度

### 6.1 `bool`

`bool` 直接用 `1 bit` 表示值, 不进入 body:

- `0 = false`
- `1 = true`

这不是“零值省略 + body”的模式, 而是直接把值编码在 header 中.

### 6.2 定宽基础字段

包括:

- `i8/u8/i16/u16/i32/u32/i64/u64`
- `f32/f64`
- `enum`

tag 宽度都是 `1 bit`:

- `0 = 默认值, 不传 body`
- `1 = 非默认值, 在 fixed block 中传实际值`

### 6.3 `text`

tag 宽度 `2 bit`:

- `00 = 空串, 不传 body`
- `01 = 长度使用 u8`
- `10 = 长度使用 u16`
- `11 = 非法`

body 形式:

- `01`: `u8(len) + bytes`
- `10`: `u16(len) + bytes`

### 6.4 `bin`

tag 宽度 `2 bit`:

- `00 = 空, 不传 body`
- `01 = 长度使用 u8`
- `10 = 长度使用 u16`
- `11 = 长度使用 u24`

body 形式:

- `01`: `u8(len) + bytes`
- `10`: `u16(len) + bytes`
- `11`: `u24(len) + bytes`

约束:

- 单个 `bin` 字段最大长度为 `16777215` (`0xFFFFFF`)
- 超过该上限直接视为编码错误

### 6.5 `list`

tag 宽度 `2 bit`:

- `00 = 空列表, 不传 body`
- `01 = count 使用 u8`
- `10 = 非法`
- `11 = 非法`

这里的长度表示的是**元素个数**, 不是字节数.

body 形式:

- `01`: `u8(count) + elements`

### 6.6 `struct`

嵌套 `struct` 使用 `1 bit`:

- `0 = 零值 struct, 不传 body`
- `1 = 递归编码整个子 struct`

## 7. Header 打包规则

- 按 schema 字段顺序依次写入每个字段的 tag
- `1 bit` 字段占 1 位
- `2 bit` 字段占 2 位
- header 视为**连续 bitstream**
- 字节内按 **MSB -> LSB** 填充
- `2 bit` tag 按**高位在前**写入
- 允许单个 tag 跨字节
- `header_size = ceil(total_tag_bits / 8)`
- 最后一个字节未使用的 padding bit 必须为 `0`

全局 bit offset 定义:

- `byteIndex = bitOffset / 8`
- `bitIndex = 7 - (bitOffset % 8)`

也就是:

- 第 1 个 bit 写到第 1 个字节的 bit7
- 第 8 个 bit 写到第 1 个字节的 bit0
- 第 9 个 bit 写到第 2 个字节的 bit7

示意:

```text
字段1(1bit) + 字段2(2bit) + 字段3(1bit) + ...
```

说明:

- header 内部不为字段做额外对齐
- 只有整个 header 末尾允许 padding
- `list` 的 `item header block` 也复用这一套 bit 打包规则

## 8. Body 编码规则

### 8.1 `fixed block`

按 schema 顺序写入所有 tag 标记为“有值”的定宽字段:

- `i/u/f`
- `enum`

编码方式继续沿用当前 `SB`:

- 小端
- 固定宽度

### 8.2 `var tail`

按 schema 顺序写入所有非空变长字段:

- `text`
- `bin`
- `list`
- `struct`

### 8.3 `list` 元素编码

外层 list 的 tag 只决定 count 的宽度.

元素按元素类型继续编码:

- `[bool]`: count 后跟 bit-packed bool 数据
- `[i8/u8/i16/u16/i32/u32/i64/u64/f32/f64/enum/struct]`: count 后跟“value bitmap + non-default bodies”
- `[text]`: count 后跟“元素状态头 + 元素数据尾部”
- `[bin]`: count 后跟“元素状态头 + 元素数据尾部”

#### 8.3.1 `[bool]`

`[bool]` 继续作为特殊类型处理:

```text
count + bitset
```

其中:

- bit=`0` 表示 `false`
- bit=`1` 表示 `true`

说明:

- 这已经等价于“零值省略”的最优形式
- `[bool]` 不再引入额外的 presence bitmap 和 body

#### 8.3.2 `1-bit default-state` list

对以下 list:

- `[i8/u8]`
- `[i16/u16]`
- `[i32/u32/f32]`
- `[i64/u64/f64]`
- `[enum]`
- `[struct]`

编码布局为:

```text
count + value bitmap + non-default bodies
```

其中:

- `value bitmap` 每个元素占 `1 bit`
- bit=`0` 表示该元素等于默认值, 不写 body
- bit=`1` 表示该元素不等于默认值, 在 `non-default bodies` 中按顺序写入
- `bitmap_size = ceil(count / 8)`

默认值定义:

- 整数默认值为 `0`
- 浮点默认值为 `0.0`
- `enum` 默认值为 schema 第一个枚举项
- `struct` 默认值为“所有字段都为默认值”

说明:

- `[bool]` 仍然单独走 bitset, 不包含在这里
- 这条规则统一了绝大部分 `1-bit default-state` list 的解码路径
- 对 `i8/u8/enum` 这类 `1 byte` 元素, 当默认值比例不高时, 包体可能比 dense 更大, 这是这版设计接受的取舍
- 对 `[struct]`, 当 bitmap 中某一位为 `1` 时, 对应元素在 `non-default bodies` 中递归按 struct 规则编码

示例:

```text
[u32] = [1, 2, 0, 0, 3]
```

若 `count` 使用 `u8`, 则:

- `count = 0x05`
- `value bitmap = 11001000 = 0xC8`
- `non-default bodies = 01 00 00 00 | 02 00 00 00 | 03 00 00 00`

因此 body 为:

```text
05 C8 01 00 00 00 02 00 00 00 03 00 00 00
```

再看 `1 byte` 元素的例子:

```text
[i8] = [1, 2, 0, 0, 3]
```

若 `count` 使用 `u8`, 则:

- `count = 0x05`
- `value bitmap = 11001000 = 0xC8`
- `non-default bodies = 01 02 03`

因此 body 为:

```text
05 C8 01 02 03
```

说明:

- 这比 dense 编码 `05 01 02 00 00 03` 更小 `1 byte`
- 但如果默认值很少, 同样可能比 dense 更大

再看 `[struct]` 的语义:

```text
[Game] = [zeroGame, gameA, zeroGame, gameB]
```

若 `count` 使用 `u8`, 则:

- `count = 0x04`
- `value bitmap = 01010000 = 0x50`
- `non-default bodies = encode(gameA) | encode(gameB)`

也就是:

```text
04 50 encode(gameA) encode(gameB)
```

其中 `zeroGame` 表示“所有字段都为默认值”的 `Game`.

其中 `[text]` / `[bin]` 的列表内布局为:

```text
count + item header block + item tail
```

- `item header block`: 每个元素占 `2 bit`
- `item tail`: 按元素顺序写入所有非空元素的数据

`[text]` 的元素状态:

- `00 = 空串`
- `01 = 长度使用 u8`
- `10 = 长度使用 u16`
- `11 = 非法`

`[bin]` 的元素状态:

- `00 = 空`
- `01 = 长度使用 u8`
- `10 = 长度使用 u16`
- `11 = 长度使用 u24`

说明:

- 列表内的 `text` / `bin` 元素与字段级 `text` / `bin` 使用同一套紧凑长度语义
- `item_header_size = ceil(count * 2 / 8)`
- 空元素只在 `item header block` 中占状态位, 不在 `item tail` 中占 body
- `item header block` 也按“连续 bitstream + MSB -> LSB”打包
- 由于 `item header block` 总是从 `count` 之后的字节边界开始, 且每个元素固定占 `2 bit`, 所以**单个元素状态不会跨字节**
- 每 `4` 个元素正好占满 `1` 个字节; 超过 `4` 个元素时, 是整个 `item header block` 跨到下一个字节, 不是某个元素状态被拆开

## 9. 规范编码要求

必须使用最短合法编码.

### 9.1 定宽字段

- 默认值必须写 tag=`0`
- 非默认值必须写 tag=`1` 并在 body 传实际值

### 9.2 `text`

- `len=0` 只能写 `00`
- `1..255` 必须写 `01`
- `256..65535` 必须写 `10`

### 9.3 `bin`

- `len=0` 只能写 `00`
- `1..255` 必须写 `01`
- `256..65535` 必须写 `10`
- `65536..16777215` 必须写 `11`
- `>16777215` 非法

### 9.4 `list`

- `count=0` 只能写 `00`
- `1..255` 必须写 `01`
- `>255` 非法

### 9.5 `[bool]`

- `[bool]` 必须编码为 `count + bitset`
- 不允许再叠加 presence bitmap 或额外 body

### 9.6 `1-bit default-state` list

- `[i8/u8/i16/u16/i32/u32/i64/u64/f32/f64/enum/struct]` 必须编码为 `count + value bitmap + non-default bodies`
- bitmap 中 bit=`0` 表示默认值元素, 不写 body
- bitmap 中 bit=`1` 表示非默认值元素, 按 schema 顺序写入 body
- 默认值元素不得写入 body
- 非默认值元素不得省略 body
- `enum` 的默认值按 schema 第一个枚举项判断, 不是按数值 `0` 判断
- `struct` 的默认值按“所有字段都为默认值”判断
- `struct` 的非默认值元素 body 必须递归使用 struct 编码规则

### 9.7 `[text]` / `[bin]` 列表元素

- `[text]` 列表中的空字符串元素必须写 `00`
- `[text]` 列表中的非空元素必须按 `u8/u16` 最短长度编码
- `[bin]` 列表中的空元素必须写 `00`
- `[bin]` 列表中的非空元素必须按 `u8/u16/u24` 最短长度编码

解码器可以对非规范编码直接报错.

## 10. `Game` 示例

Schema:

```sb
Game{
    id u32
    name text
}
```

字段 tag 宽度:

- `id`: `1 bit`
- `name`: `2 bit`

总共 `3 bit`, 所以 header 占 `1 byte`.

### 10.1 `Game{id:0, name:"lol"}`

逻辑 tag:

- `id = 0`
- `name = 01`

header 位流:

```text
0 | 01 | 00000
```

按连续 bitstream 且 `MSB -> LSB` 填充后:

```text
00100000 = 0x20
```

body:

- `fixed block`: 空
- `var tail`: `03 6C 6F 6C`

最终 wire:

```text
20 03 6C 6F 6C
```

### 10.2 `Game{id:7, name:""}`

逻辑 tag:

- `id = 1`
- `name = 00`

header 位流:

```text
1 | 00 | 00000
```

按连续 bitstream 且 `MSB -> LSB` 填充后:

```text
10000000 = 0x80
```

body:

- `fixed block`: `07 00 00 00`
- `var tail`: 空

最终 wire:

```text
80 07 00 00 00
```

### 10.3 `Game{id:7, name:"lol"}`

逻辑 tag:

- `id = 1`
- `name = 01`

header 位流:

```text
1 | 01 | 00000
```

按连续 bitstream 且 `MSB -> LSB` 填充后:

```text
10100000 = 0xA0
```

body:

- `fixed block`: `07 00 00 00`
- `var tail`: `03 6C 6F 6C`

最终 wire:

```text
A0 07 00 00 00 03 6C 6F 6C
```

### 10.4 跨字节 header 示例

Schema:

```sb
Demo{
    id u32
    title text
    desc text
    blob bin
    items [u32]
}
```

字段 tag 宽度:

- `id`: `1 bit`
- `title`: `2 bit`
- `desc`: `2 bit`
- `blob`: `2 bit`
- `items`: `2 bit`

总共 `9 bit`, 所以 header 占 `2 byte`.

假设:

- `id = 7`，所以 `id = 1`
- `title = "a"`，长度可用 `u8`，所以 `title = 01`
- `desc` 长度为 `300`，需用 `u16`，所以 `desc = 10`
- `blob` 为空，所以 `blob = 00`
- `items = [1, 2, 3]`，count 可用 `u8`，所以 `items = 01`

header 位流:

```text
1 | 01 | 10 | 00 | 01
```

把它当作连续 bitstream, 按 `MSB -> LSB` 填充后:

```text
10110000 10000000
   0xB0     0x80
```

这里最后一个字段 `items = 01` 正好跨字节:

- `items` 的第 1 个 bit(`0`) 落在第 1 个字节的 bit0
- `items` 的第 2 个 bit(`1`) 落在第 2 个字节的 bit7

这就是“允许 2-bit tag 跨字节”的具体效果.

### 10.5 `[text]` / `[bin]` 列表的 `item header block` 示例

对于 `[text]` / `[bin]` 列表, `item header block` 虽然也是 bitstream, 但因为它从 `count` 之后的字节边界开始, 且每个元素状态固定占 `2 bit`, 所以**单个元素状态不会跨字节**.

假设:

- 有一个 `[text]` 列表, `count = 5`
- 5 个元素状态依次为: `01 | 00 | 10 | 01 | 00`

那么前 `4` 个元素正好占满第 `1` 个字节:

```text
01 | 00 | 10 | 01
= 01001001
= 0x49
```

第 `5` 个元素从下一个字节重新开始:

```text
00 | 000000
= 00000000
= 0x00
```

因此 `item header block` 是:

```text
49 00
```

这里跨字节的是整个 `item header block`, 不是某个单独元素状态.

## 11. 设计取舍

这版草案的核心取舍是:

- 用更复杂的 header 换更小的零值/短字符串包体
- 保留顺序读取, 不引入 pointer
- 不为部分读取优化

它适合:

- 完整编解码
- 固定 schema
- 短 `text` / 短 `bin` / 小 list 很常见

它不适合:

- 需要强兼容演进
- 需要随机访问
- 希望 runtime 规则尽量极简到只有“presence bit + 固定长度前缀”

## 12. 错误与规范化策略

当前协议采用**严格解码**策略, 不提供“尽量兼容”的规范模式.

### 12.1 非规范编码

以下情况统一视为解码错误:

- 本应使用更短编码, 却使用了更长编码
- 默认值元素仍然写入了 body
- 非默认值元素缺少 body
- `text` / `bin` / `list` 使用了非法状态值
- 长度或 count 超出该状态允许的范围
- header、bitmap、item header block 与 body 实际内容不一致

也就是:

- 解码器应拒绝所有 non-canonical encoding
- 不应在规范实现中自动“纠正”或“兼容”这些输入

### 12.2 错误信息

运行时错误信息建议至少包含以下字段:

- 路径
  - 例如 `Game.Name`
  - 或 `Games[3].Title`
- 阶段
  - `header`
  - `fixed_block`
  - `var_tail`
  - `value_bitmap`
  - `item_header_block`
- 偏移
  - `byte offset`
  - 如果发生在 bitstream 中, 再补 `bit offset`
- 状态
  - 当前读到的 `1 bit` / `2 bit` 状态值
- 长度信息
  - 实际 `len` / `count`
  - 以及对应的允许范围

错误文本不必强制完全一致, 但应能回答这几个问题:

- 出错的是哪个字段
- 出错发生在解码的哪一阶段
- 读到了什么状态或长度
- 为什么这在当前规则下非法

## 13. 实现与维护要点

当前仓库已经默认使用这套协议, Go / TypeScript 两端都按这里的规则生成代码.

### 13.1 运行时

Go 与 TypeScript 运行时都需要持续保持以下基础能力一致:

- bitstream header 写入与读取
  - 支持 `1 bit` 与 `2 bit` 状态
  - 支持 `MSB -> LSB`
  - 支持字段级 header 可跨字节
- `u24` 读写
  - `GetU24`
  - `SetU24`
  - 上限检查与错误信息
- 紧凑长度读写
  - `text: empty/u8/u16`
  - `bin: empty/u8/u16/u24`
  - `list count: empty/u8`
- list 基础读写
  - `[bool] -> count + bitset`
  - `[scalar/enum/struct] -> count + value bitmap + non-default bodies`
  - `[text]/[bin] -> count + item header block + item tail`

### 13.2 代码生成器

生成器需要始终按当前规则生成:

- struct field 分类
  - `bool`
  - `1-bit default-state`: `scalar/enum/struct`
  - `2-bit state`: `text/bin/list`
- struct 编码布局
  - 先写 header
  - 再写 `fixed block`
  - 最后写 `var tail`
- struct 解码布局
  - 先读 header
  - 再按 schema 顺序读 `fixed block`
  - 最后按 schema 顺序读 `var tail`
- 默认值判定生成
  - `enum` 默认值 = schema 第一项
  - `struct` 默认值 = 所有字段都为默认值
- list 生成
  - `[bool]` 生成 bitset 路径
  - `[scalar/enum/struct]` 生成 value bitmap 路径
  - `[text]/[bin]` 生成 item header block 路径

### 13.3 辅助函数

为了让生成代码保持简单, 运行时最好提供统一 helper:

- `writeHeaderBits` / `readHeaderBits`
- `writeState1` / `readState1`
- `writeState2` / `readState2`
- `sizeTextCompact` / `sizeBinCompact`
- `sizeScalarListBitmap`
- `sizeTextListCompact`
- `sizeBinListCompact`
- `isZeroStruct`

其中 `isZeroStruct` 更适合由生成器按具体 struct 内联生成, 不建议做成反射式通用函数.

### 13.4 测试

建议长期保留以下回归样例:

- `Game{id:0, name:"lol"} -> 20 03 6C 6F 6C`
- `Game{id:7, name:""} -> 80 07 00 00 00`
- `Game{id:7, name:"lol"} -> A0 07 00 00 00 03 6C 6F 6C`
- `Demo` 的跨字节 header 示例
- `[text]` / `[bin]` 的 `item header block = 49 00` 示例
- `[u32] = [1, 2, 0, 0, 3]` 位图示例
- `[i8] = [1, 2, 0, 0, 3]` 位图示例
- `[struct]` list 的 bitmap 语义示例

还需要补边界测试:

- `u8/u16/u24` 长度边界
- `count=0/255/256`
- 非规范编码
- `-0.0` 被规范化为 `0.0`

### 13.5 当前工程落点

当前仓库中的主要落点如下:

- Go runtime: `go/sb/rt`
- TypeScript runtime: `ts/sb/type.ts`
- Go 生成器: `internal/tpl_go_render.go`
- TypeScript 生成器: `internal/tpl_ts_render.go`
- 双端一致性回归:
  - `go/sb/cross_consistency_test.go`
  - `ts/sb/cross_consistency.test.ts`

后续维护应直接在这套主链路上继续演进, 不再保留并行协议实现.
