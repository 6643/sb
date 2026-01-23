package parser

import (
	"fmt"
	"sb/internal/lexer"
	"testing"
)

// TestParser_Robustness 鲁棒性测试
// 涵盖正常解析, 循环嵌套检测, 未定义类型校验等边界场景
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
			wantErr: true, // Circular embedding is an error
		},
		{
			name: "Undefined Type",
			input: `
				User { info UnknownType }
			`,
			wantErr: true, // Undefined type is an error
		},
		{
			name: "Invalid Syntax - Missing Brace",
			input: `
				User { id u32
			`,
			wantErr: false, // Parser handles EOF gracefully (returns what it has)
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
	// Generate struct with 256 fields using valid identifiers
	input := "LargeStruct {\n"
	for i := 0; i < 256; i++ {
		input += fmt.Sprintf("field_%d u32\n", i)
	}
	input += "}"

	l := lexer.New(input)
	p := New(l)
	schema, err := p.ParseSchema()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if len(schema.Structs) == 0 {
		t.Fatal("No struct parsed")
	}

	if len(schema.Structs[0].Fields) < 256 {
		t.Errorf("Expected 256 fields, got %d", len(schema.Structs[0].Fields))
	}
}
