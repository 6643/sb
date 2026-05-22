package internal

import "fmt"

func (p *Parser) parseAPI(note string) (API, error) {
	api := API{Note: note}
	name := p.curToken.Literal
	p.nextToken()

	if p.curToken.Type == TokenDot {
		p.nextToken()
		if p.curToken.Type != TokenIdent {
			return api, fmt.Errorf("行 %d: API 名称缺少段", p.curToken.Line)
		}
		name += "." + p.curToken.Literal
		p.nextToken()
	}
	if p.curToken.Type == TokenDot {
		return api, fmt.Errorf("行 %d: API 名称仅支持 1-2 段", p.curToken.Line)
	}
	api.Name = name

	if err := p.expectCurrent(TokenLParen, "("); err != nil {
		return api, err
	}
	p.nextToken()
	p.skipNewLines()

	for p.curToken.Type != TokenRParen {
		if p.curToken.Type == TokenEOF {
			return api, fmt.Errorf("行 %d: API %s 缺少 )", p.curToken.Line, api.Name)
		}
		p.skipNewLines()
		if p.curToken.Type == TokenRParen {
			break
		}
		if p.curToken.Type != TokenIdent {
			return api, fmt.Errorf("行 %d: API 参数名非法 %q", p.curToken.Line, p.curToken.Literal)
		}
		arg := APIArg{Name: p.curToken.Literal}
		p.nextToken()

		t, err := p.parseType(false)
		if err != nil {
			return api, err
		}
		arg.Type = t
		api.Args = append(api.Args, arg)
		p.skipNewLines()

		if p.curToken.Type == TokenComma {
			p.nextToken()
			p.skipNewLines()
			if p.curToken.Type == TokenRParen {
				return api, fmt.Errorf("行 %d: API 参数列表尾逗号非法", p.curToken.Line)
			}
			continue
		}
		if p.curToken.Type != TokenRParen {
			return api, fmt.Errorf("行 %d: API 参数后需要 , 或 )", p.curToken.Line)
		}
	}

	p.nextToken()
	p.skipNewLines()
	if p.curToken.Type != TokenArrow {
		return api, fmt.Errorf("行 %d: API 定义缺少 =>", p.curToken.Line)
	}
	p.nextToken()
	p.skipNewLines()

	result, err := p.parseType(true)
	if err != nil {
		return api, err
	}
	api.Result = result

	if p.curToken.Type == TokenComment {
		return api, fmt.Errorf("行 %d: API 注释必须写在定义前, 不支持尾部注释", p.curToken.Line)
	}

	return api, nil
}
