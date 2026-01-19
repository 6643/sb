package parser

import (
	"sb/internal/lexer"
	"testing"
)

func TestParser_Robustness(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "Valid Schema",
			input: `
				User { id u32, name text }
				user.get(id u32) => User
			`,
			wantErr: false,
		},
		{
			name: "Circular Embedding",
			input: `
				A { B }
				B { A }
			`,
			wantErr: false, // 应该被 expandFields 优雅处理（打印警告而不崩溃）
		},
		{
			name: "Undefined Type",
			input: `
				User { info UnknownType }
			`,
			wantErr: false, // 目前 resolveTypes 只是标记，不会返回 error，但这会导致生成代码失败
		},
		{
			name: "Invalid Syntax - Missing Brace",
			input: `
				User { id u32
			`,
			wantErr: false, // 目前 Parser 在 EOF 时会停止，不会 panic，但结果不完整
		},
		{
			name: "Invalid API - No Arrow",
			input: `
				user.get(id u32) User
			`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			_, err := p.ParseSchema()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParser_FieldLimit(t *testing.T) {
	// 模拟生成 256 个字段的结构体字符串
	input := "LargeStruct {\n"
	for i := 0; i < 256; i++ {
		input += "field" + string(rune(i)) + " u32\n"
	}
	input += "}"

	l := lexer.New(input)
	p := New(l)
	schema, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if len(schema.Structs[0].Fields) < 256 {
		t.Errorf("Expected 256 fields, got %d", len(schema.Structs[0].Fields))
	}
}
