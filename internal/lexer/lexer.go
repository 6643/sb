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
		TokenComment  // 注释
	)
	
	type Token struct {
		Type  TokenType
		Value string
		Line  int
	}
	
	// Lexer 词法分析器状态
	type Lexer struct {
		input []rune // 完整的输入字符流
		pos   int    // 当前处理的字符位置
		line  int    // 当前行号 (用于错误报告)
	}
	
	func New(input string) *Lexer {
		return &Lexer{input: []rune(input), line: 1}
	}
	
	// NextToken 获取下一个 Token (核心状态机)
	// 自动跳过空白, 处理标识符, 数字, 字符串及特殊符号
	func (l *Lexer) NextToken() Token {
		l.skipWhitespace()
	
		if l.pos >= len(l.input) {
			return Token{Type: TokenEOF, Line: l.line}
		}
	
		ch := l.input[l.pos]
	
		// 标识符 (Identifier)
		if isIdentStart(ch) {
			return l.readIdent()
		}
	
		// 数字 (Number) - 支持负数
		if unicode.IsDigit(ch) || (ch == '-' && unicode.IsDigit(l.peek())) {
			return l.readNumber()
		}
	
		// Tag 字符串 ("..." or `...`)
		if ch == '"' || ch == '`' {
			return l.readTag(ch)
		}
	
		// 双字符 Token
		if ch == '/' && l.peek() == '/' {
			return l.readComment()
		}
		if ch == '=' && l.peek() == '>' {
			l.pos += 2
			return Token{Type: TokenArrow, Value: "=>", Line: l.line}
		}
	
		// 单字符 Token
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
	
		// 错误处理: 遇到非法字符必须推进指针, 防止死循环
		l.pos++
		return Token{Type: TokenError, Value: fmt.Sprintf("未预期字符 %q", ch), Line: l.line}
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
		l.pos++ // 跳过起始引号
		for l.pos < len(l.input) && l.input[l.pos] != quote {
			if l.input[l.pos] == '\n' {
				l.line++
			}
			l.pos++
		}
		if l.pos < len(l.input) {
			l.pos++ // 跳过结束引号
		}
		return Token{Type: TokenIdent, Value: string(l.input[start:l.pos]), Line: l.line}
	}
	
	func (l *Lexer) readComment() Token {
		l.pos += 2 // 跳过 //
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '\n' {
			l.pos++
		}
		return Token{Type: TokenComment, Value: strings.TrimSpace(string(l.input[start:l.pos])), Line: l.line}
	}
	