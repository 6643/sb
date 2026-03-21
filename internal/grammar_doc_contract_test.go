package internal

import (
	"os"
	"testing"
)

func TestGrammarDocContract(t *testing.T) {
	grammar, err := os.ReadFile("../grammar.peg")
	if err != nil {
		t.Fatalf("read grammar.peg: %v", err)
	}
	implementation, err := os.ReadFile("../IMPLEMENTATION.md")
	if err != nil {
		t.Fatalf("read IMPLEMENTATION.md: %v", err)
	}

	assertContains(t, string(grammar), "仅作设计参考")
	assertContains(t, string(grammar), "不参与当前解析实现")
	assertContains(t, string(implementation), "grammar.peg")
	assertContains(t, string(implementation), "internal/lexer_scanner.go")
	assertContains(t, string(implementation), "internal/parser_schema.go")
	assertContains(t, string(implementation), "internal/parser_schema_test.go")
}
