package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenIdent
	TokenNumber
	TokenLBrace   // {
	TokenRBrace   // }
	TokenLParen   // (
	TokenRParen   // )
	TokenLBracket // [
	TokenRBracket // ]
	TokenAssign   // =
	TokenPipe     // |
	TokenComma    // ,
	TokenDot      // .
	TokenArrow    // =>
	TokenComment
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}

type Lexer struct {
	input []rune
	pos   int
	line  int
}

func New(input string) *Lexer {
	return &Lexer{input: []rune(input), line: 1}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF, Line: l.line}
	}

	ch := l.input[l.pos]

	if isIdentStart(ch) {
		return l.readIdent()
	}

	if unicode.IsDigit(ch) || (ch == '-' && unicode.IsDigit(l.peek())) {
		return l.readNumber()
	}

	if ch == '"' || ch == '`' {
		return l.readTag(ch)
	}

	// Double char tokens
	if ch == '/' && l.peek() == '/' {
		return l.readComment()
	}
	if ch == '=' && l.peek() == '>' {
		l.pos += 2
		return Token{Type: TokenArrow, Value: "=>", Line: l.line}
	}

	// Single char tokens
	switch ch {
	case '{':
		return l.advanceAndMakeToken(TokenLBrace, "{")
	case '}':
		return l.advanceAndMakeToken(TokenRBrace, "}")
	case '(':
		return l.advanceAndMakeToken(TokenLParen, "(")
	case ')':
		return l.advanceAndMakeToken(TokenRParen, ")")
	case '[':
		return l.advanceAndMakeToken(TokenLBracket, "[")
	case ']':
		return l.advanceAndMakeToken(TokenRBracket, "]")
	case '=':
		return l.advanceAndMakeToken(TokenAssign, "=")
	case '|':
		return l.advanceAndMakeToken(TokenPipe, "|")
	case ',':
		return l.advanceAndMakeToken(TokenComma, ",")
	case '.':
		return l.advanceAndMakeToken(TokenDot, ".")
	}

	// Error handling - MUST advance to avoid infinite loop
	l.pos++
	return Token{Type: TokenError, Value: fmt.Sprintf("unexpected character %q", ch), Line: l.line}
}

func (l *Lexer) advanceAndMakeToken(t TokenType, val string) Token {
	l.pos++
	return Token{Type: t, Value: val, Line: l.line}
}

func (l *Lexer) peek() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\n' {
			l.line++
			l.pos++
			continue
		}
		if unicode.IsSpace(ch) {
			l.pos++
			continue
		}
		break
	}
}

func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isIdentPart(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '.'
}

func (l *Lexer) readIdent() Token {
	start := l.pos
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	return Token{Type: TokenIdent, Value: string(l.input[start:l.pos]), Line: l.line}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	if l.input[l.pos] == '-' {
		l.pos++
	}
	for l.pos < len(l.input) && unicode.IsDigit(l.input[l.pos]) {
		l.pos++
	}
	return Token{Type: TokenNumber, Value: string(l.input[start:l.pos]), Line: l.line}
}

func (l *Lexer) readTag(quote rune) Token {
	start := l.pos
	l.pos++ // skip start quote
	for l.pos < len(l.input) && l.input[l.pos] != quote {
		if l.input[l.pos] == '\n' {
			l.line++
		}
		l.pos++
	}
	if l.pos < len(l.input) {
		l.pos++ // skip end quote
	}
	return Token{Type: TokenIdent, Value: string(l.input[start:l.pos]), Line: l.line}
}

func (l *Lexer) readComment() Token {
	l.pos += 2 // skip //
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.pos++
	}
	return Token{Type: TokenComment, Value: strings.TrimSpace(string(l.input[start:l.pos])), Line: l.line}
}