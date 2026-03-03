package parser

import (
	"testing"

	"sb/internal/lexer"
)

func TestParseSchema(t *testing.T) {
	input := `
		Status = Ok(0) | Err(1)
		User { id u32, name text }
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
