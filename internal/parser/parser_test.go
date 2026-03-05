package parser

import (
	"strings"
	"testing"

	"sb/internal/lexer"
)

func TestParseSchema(t *testing.T) {
	input := `
		Status = Ok(0) | Err(1)
		User {
			id u32
			name text
		}
		user.get(id u32) => User
	`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.Enums) != 1 || len(s.Structs) != 1 || len(s.APIs) != 1 {
		t.Fatalf("unexpected schema counts: enums=%d structs=%d apis=%d", len(s.Enums), len(s.Structs), len(s.APIs))
	}
}

func TestParseSchemaMissingBrace(t *testing.T) {
	input := `User { id u32`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestParseCommentLineWithoutTrailingNewlineRejected(t *testing.T) {
	input := `// only comment`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "注释行必须以换行结束") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseStructCommentLineWithoutTrailingNewlineRejected(t *testing.T) {
	input := `User {
// only comment`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "注释行必须以换行结束") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseAPIHeaderMultiLineNote(t *testing.T) {
	input := `
	// 第一行
	// 第二行
	user.get_count(page u8) => u8
	`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.APIs) != 1 {
		t.Fatalf("unexpected api count: %d", len(s.APIs))
	}
	if s.APIs[0].Note != "第一行\n第二行" {
		t.Fatalf("unexpected api note: %q", s.APIs[0].Note)
	}
}

func TestParseAPITailCommentRejected(t *testing.T) {
	input := `user.get_count(page u8) => u8 // 尾部注释`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "API 注释必须写在定义前") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseAPITailCommentRejectedOnResultNextLine(t *testing.T) {
	input := `user.get_count(page u8) =>
u8 // 尾部注释`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "API 注释必须写在定义前") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseAPIAllowsCommentLineAfterDefinition(t *testing.T) {
	input := `get_count(page u8) => u8
// 第二个 API 注释
get_bin(page u8) => bin`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.APIs) != 2 {
		t.Fatalf("unexpected api count: %d", len(s.APIs))
	}
	if s.APIs[1].Note != "第二个 API 注释" {
		t.Fatalf("unexpected api note: %q", s.APIs[1].Note)
	}
}

func TestParseEnumRequiresAssign(t *testing.T) {
	input := `Status = Ok | Err`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.Enums) != 1 {
		t.Fatalf("unexpected enum count: %d", len(s.Enums))
	}
}

func TestParseEnumWithoutAssignRejected(t *testing.T) {
	input := `Status | Ok | Err`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "枚举定义必须包含 =") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseEnumInternalRequiresAssign(t *testing.T) {
	l := lexer.New(`Status Ok`)
	p := New(l)
	_, err := p.parseEnum("")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "枚举定义必须包含 =") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseAPIMissingArrowRejected(t *testing.T) {
	input := `user.get_count(page u8) u8`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "API 定义缺少 =>") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseAPIArgNilRejected(t *testing.T) {
	input := `user.set(v nil) => u8`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "nil 仅允许作为 API 返回类型") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseFieldListNilRejected(t *testing.T) {
	input := `User {
items [nil]
}`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "nil 不能作为数组元素类型") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseFieldNilRejected(t *testing.T) {
	input := `User {
v nil
}`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "nil 仅允许作为 API 返回类型") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseStructCommentLineNotInlineFieldNote(t *testing.T) {
	input := `User {
id u32
// not inline
}`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.Structs) != 1 || len(s.Structs[0].Fields) != 1 {
		t.Fatalf("unexpected schema: %+v", s)
	}
	if s.Structs[0].Fields[0].Note != "" {
		t.Fatalf("unexpected field note: %q", s.Structs[0].Fields[0].Note)
	}
}

func TestParseAPINamedSingleSegmentAccepted(t *testing.T) {
	input := `get_count(page u8) => u8`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.APIs) != 1 || s.APIs[0].Name != "get_count" {
		t.Fatalf("unexpected apis: %+v", s.APIs)
	}
}

func TestParseAPINamedTwoSegmentsAccepted(t *testing.T) {
	input := `user.get_count(page u8) => u8`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.APIs) != 1 || s.APIs[0].Name != "user.get_count" {
		t.Fatalf("unexpected apis: %+v", s.APIs)
	}
}

func TestParseAPINamedThreeSegmentsRejected(t *testing.T) {
	input := `a.b.c(page u8) => u8`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "API 名称仅支持 1-2 段") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseEnumMissingPipeBetweenMembersRejected(t *testing.T) {
	input := `Status = Ok Err`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "成员之间缺少 |") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseStructFieldTagWithoutTypeRejected(t *testing.T) {
	input := `User { id "_id" }`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "缺少类型, 不允许直接写 tag") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseEnumCommaSeparatorRejected(t *testing.T) {
	input := `Status = Ok, Err`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "成员之间仅允许 |") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseEnumNegativeValueRejected(t *testing.T) {
	input := `Status = Ok(-1)`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "未预期字符 '-'") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseEnumValueOverflowRejected(t *testing.T) {
	input := `Status = Ok(256)`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "无效枚举值") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseEnumValueLeadingZeroRejected(t *testing.T) {
	input := `Status = Ok(01)`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "无效枚举值") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseSchemaWithNewLines(t *testing.T) {
	input := `
Status = Ok(0)
    | Err(1)

User {
    id u32
    name text
}

user.get(
    id u32
) => User
`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.Enums) != 1 || len(s.Structs) != 1 || len(s.APIs) != 1 {
		t.Fatalf("unexpected schema counts: enums=%d structs=%d apis=%d", len(s.Enums), len(s.Structs), len(s.APIs))
	}
	if len(s.Enums[0].Members) != 2 {
		t.Fatalf("unexpected enum member count: %d", len(s.Enums[0].Members))
	}
	if len(s.Structs[0].Fields) != 2 {
		t.Fatalf("unexpected struct field count: %d", len(s.Structs[0].Fields))
	}
}

func TestParseStructFieldsMustBeOnePerLine(t *testing.T) {
	input := `User { id u32 name text }`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "字段必须独占一行") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseStructCommaMustBeStandaloneLine(t *testing.T) {
	input := `User { id u32, name text }`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "逗号必须独占一行") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseStructFieldLineTrailingCommaRejected(t *testing.T) {
	input := `User {
id u32,
}`
	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseSchema()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "逗号必须独占一行") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseStructStandaloneCommaLineAccepted(t *testing.T) {
	input := `User {
id u32
,
name text
}`
	l := lexer.New(input)
	p := New(l)
	s, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("ParseSchema failed: %v", err)
	}
	if len(s.Structs) != 1 {
		t.Fatalf("unexpected struct count: %d", len(s.Structs))
	}
	if len(s.Structs[0].Fields) != 2 {
		t.Fatalf("unexpected field count: %d", len(s.Structs[0].Fields))
	}
	if s.Structs[0].Fields[0].Name != "id" || s.Structs[0].Fields[1].Name != "name" {
		t.Fatalf("unexpected fields: %+v", s.Structs[0].Fields)
	}
}

func TestParseEnumAssignMustBeSameLine(t *testing.T) {
	l := lexer.New("Status\n= Ok")
	p := New(l)
	_, err := p.parseEnum("")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "枚举定义必须包含 =") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseAPINameAndLParenMustBeSameLine(t *testing.T) {
	l := lexer.New("user\n(id u8) => u8")
	p := New(l)
	_, err := p.parseAPI("")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "期望 (") {
		t.Fatalf("unexpected error: %v", err)
	}
}
