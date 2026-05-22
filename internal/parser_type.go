package internal

import "fmt"

func (p *Parser) parseType(allowNil bool) (TypeRef, error) {
	var t TypeRef

	if p.curToken.Type == TokenIdent {
		if allowNil && p.curToken.Literal == "nil" {
			t.Name = "nil"
			p.nextToken()
			return t, nil
		}
		if !allowNil && p.curToken.Literal == "nil" {
			return t, fmt.Errorf("行 %d: nil 仅允许作为 API 返回类型", p.curToken.Line)
		}
		t.Name = p.curToken.Literal
		p.nextToken()
		return t, nil
	}

	if p.curToken.Type != TokenLBracket {
		return t, fmt.Errorf("行 %d: 非法类型 %q", p.curToken.Line, p.curToken.Literal)
	}

	t.IsList = true
	p.nextToken()
	p.skipNewLines()
	if p.curToken.Type != TokenIdent {
		return t, fmt.Errorf("行 %d: 数组类型缺少元素类型", p.curToken.Line)
	}
	if p.curToken.Literal == "nil" {
		return t, fmt.Errorf("行 %d: nil 不能作为数组元素类型", p.curToken.Line)
	}

	t.Name = p.curToken.Literal
	p.nextToken()
	p.skipNewLines()
	if err := p.expectCurrent(TokenRBracket, "]"); err != nil {
		return t, err
	}
	p.nextToken()
	return t, nil
}
