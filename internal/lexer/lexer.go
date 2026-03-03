package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// Lexer 负责把输入文本转换为 token 流。
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
		return l.readString(ch)
	}
	if ch == '/' && l.peek() == '/' {
		return l.readComment()
	}
	if ch == '=' && l.peek() == '>' {
		l.pos += 2
		return Token{Type: TokenArrow, Literal: "=>", Line: l.line}
	}

	switch ch {
	case '{':
		return l.advance(TokenLBrace, "{")
	case '}':
		return l.advance(TokenRBrace, "}")
	case '(':
		return l.advance(TokenLParen, "(")
	case ')':
		return l.advance(TokenRParen, ")")
	case '[':
		return l.advance(TokenLBracket, "[")
	case ']':
		return l.advance(TokenRBracket, "]")
	case '=':
		return l.advance(TokenAssign, "=")
	case '|':
		return l.advance(TokenPipe, "|")
	case ',':
		return l.advance(TokenComma, ",")
	case '.':
		return l.advance(TokenDot, ".")
	}

	l.pos++
	return Token{Type: TokenError, Literal: fmt.Sprintf("未预期字符 %q", ch), Line: l.line}
}

func (l *Lexer) advance(tt TokenType, lit string) Token {
	l.pos++
	return Token{Type: tt, Literal: lit, Line: l.line}
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
			l.pos++
			l.line++
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
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}

func (l *Lexer) readIdent() Token {
	start := l.pos
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	return Token{Type: TokenIdent, Literal: string(l.input[start:l.pos]), Line: l.line}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	if l.input[l.pos] == '-' {
		l.pos++
	}
	for l.pos < len(l.input) && unicode.IsDigit(l.input[l.pos]) {
		l.pos++
	}
	return Token{Type: TokenNumber, Literal: string(l.input[start:l.pos]), Line: l.line}
}

func (l *Lexer) readString(quote rune) Token {
	line := l.line
	l.pos++
	start := l.pos

	for l.pos < len(l.input) && l.input[l.pos] != quote {
		if l.input[l.pos] == '\n' {
			l.line++
		}
		l.pos++
	}

	if l.pos >= len(l.input) {
		return Token{Type: TokenError, Literal: "字符串未闭合", Line: line}
	}

	lit := string(l.input[start:l.pos])
	l.pos++
	return Token{Type: TokenString, Literal: lit, Line: line}
}

func (l *Lexer) readComment() Token {
	line := l.line
	l.pos += 2
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.pos++
	}
	return Token{Type: TokenComment, Literal: strings.TrimSpace(string(l.input[start:l.pos])), Line: line}
}
