package internal

import "testing"

func TestDoubleQuoteStringToken(t *testing.T) {
	l := New(`"abc"`)
	tok := l.NextToken()
	if tok.Type != TokenString {
		t.Fatalf("expected TokenString, got %v (%q)", tok.Type, tok.Literal)
	}
	if tok.Literal != "abc" {
		t.Fatalf("expected literal abc, got %q", tok.Literal)
	}
}

func TestBacktickStringRejected(t *testing.T) {
	l := New("`abc`")
	tok := l.NextToken()
	if tok.Type != TokenError {
		t.Fatalf("expected TokenError, got %v (%q)", tok.Type, tok.Literal)
	}
}

func TestNewLineTokenLF(t *testing.T) {
	l := New("\n")
	tok := l.NextToken()
	if tok.Type != TokenNewLine || tok.Literal != "\n" || tok.Line != 1 {
		t.Fatalf("unexpected token: %+v", tok)
	}
}

func TestNewLineTokenCRLFAndLineAdvance(t *testing.T) {
	l := New(" \t\r\nabc")
	tok := l.NextToken()
	if tok.Type != TokenNewLine || tok.Literal != "\r\n" || tok.Line != 1 {
		t.Fatalf("unexpected newline token: %+v", tok)
	}
	next := l.NextToken()
	if next.Type != TokenIdent || next.Literal != "abc" || next.Line != 2 {
		t.Fatalf("unexpected next token: %+v", next)
	}
}

func TestEllipsisToken(t *testing.T) {
	l := New("...User")
	tok := l.NextToken()
	if tok.Type != TokenEllipsis || tok.Literal != "..." {
		t.Fatalf("unexpected token: %+v", tok)
	}
	next := l.NextToken()
	if next.Type != TokenIdent || next.Literal != "User" {
		t.Fatalf("unexpected next token: %+v", next)
	}
}

func TestStandaloneCRRejected(t *testing.T) {
	l := New("\rabc")
	tok := l.NextToken()
	if tok.Type != TokenError {
		t.Fatalf("expected TokenError, got %v (%q)", tok.Type, tok.Literal)
	}
}

func TestCommentStopsAtCR(t *testing.T) {
	l := New("// note\rabc")
	tok := l.NextToken()
	if tok.Type != TokenComment || tok.Literal != "note" {
		t.Fatalf("unexpected comment token: %+v", tok)
	}
	next := l.NextToken()
	if next.Type != TokenError {
		t.Fatalf("expected TokenError after comment, got %+v", next)
	}
}

func TestNonASCIIIdentRejected(t *testing.T) {
	l := New("用户")
	tok := l.NextToken()
	if tok.Type != TokenError {
		t.Fatalf("expected TokenError, got %v (%q)", tok.Type, tok.Literal)
	}
}

func TestNonASCIISpaceNotSkipped(t *testing.T) {
	l := New("\u00a0abc")
	tok := l.NextToken()
	if tok.Type != TokenError {
		t.Fatalf("expected TokenError, got %v (%q)", tok.Type, tok.Literal)
	}
}

func TestNegativeNumberRejectedAtLexer(t *testing.T) {
	l := New("-1")
	tok := l.NextToken()
	if tok.Type != TokenError {
		t.Fatalf("expected TokenError, got %v (%q)", tok.Type, tok.Literal)
	}
}

func TestNonASCIIDigitRejected(t *testing.T) {
	l := New("１２")
	tok := l.NextToken()
	if tok.Type != TokenError {
		t.Fatalf("expected TokenError, got %v (%q)", tok.Type, tok.Literal)
	}
}
