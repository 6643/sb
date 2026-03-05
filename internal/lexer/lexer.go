package lexer

import (
	"fmt"
	"strings"
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
	l.skipSpaces()
	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF, Line: l.line}
	}

	ch := l.input[l.pos]
	if ch == '\n' || (ch == '\r' && l.peek() == '\n') {
		return l.readNewLine()
	}
	if isIdentStart(ch) {
		return l.readIdent()
	}
	if ch >= '0' && ch <= '9' {
		return l.readNumber()
	}
	if ch == '"' {
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

func (l *Lexer) skipSpaces() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' {
			l.pos++
			continue
		}
		break
	}
}

func (l *Lexer) readNewLine() Token {
	line := l.line
	if l.input[l.pos] == '\r' {
		l.pos += 2
		l.line++
		return Token{Type: TokenNewLine, Literal: "\r\n", Line: line}
	}
	l.pos++
	l.line++
	return Token{Type: TokenNewLine, Literal: "\n", Line: line}
}

func isIdentStart(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isIdentPart(ch rune) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
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
	for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
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
	for l.pos < len(l.input) && l.input[l.pos] != '\n' && l.input[l.pos] != '\r' {
		l.pos++
	}
	return Token{Type: TokenComment, Literal: strings.TrimSpace(string(l.input[start:l.pos])), Line: line}
}
