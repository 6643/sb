package internal

import "fmt"

// Parser 负责把 token 流转换为 AST。
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) skipNewLines() {
	for p.curToken.Type == TokenNewLine {
		p.nextToken()
	}
}

func (p *Parser) expectCurrent(tt TokenType, name string) error {
	if p.curToken.Type == tt {
		return nil
	}
	return fmt.Errorf("行 %d: 期望 %s, 实际 %q", p.curToken.Line, name, p.curToken.Literal)
}
