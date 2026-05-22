package internal

import (
	"fmt"
	"strconv"
)

func (p *Parser) parseEnum(note string) (Enum, error) {
	e := Enum{Name: p.curToken.Literal, Note: note}
	p.nextToken()
	if err := p.expectCurrent(TokenAssign, "="); err != nil {
		return e, fmt.Errorf("行 %d: 枚举定义必须包含 =", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != TokenEOF {
		p.skipNewLines()
		if p.curToken.Type == TokenComment {
			return e, fmt.Errorf("行 %d: 枚举仅支持定义前整体注释, 不支持成员注释", p.curToken.Line)
		}
		if p.curToken.Type != TokenIdent {
			break
		}
		memberLine := p.curToken.Line
		member, err := p.parseEnumMember()
		if err != nil {
			return e, err
		}
		e.Members = append(e.Members, member)

		if p.curToken.Type == TokenComment && p.curToken.Line == memberLine {
			return e, fmt.Errorf("行 %d: 枚举仅支持定义前整体注释, 不支持成员注释", p.curToken.Line)
		}
		if p.curToken.Type == TokenIdent && p.curToken.Line == memberLine {
			return e, fmt.Errorf("行 %d: 枚举 %s 成员之间缺少 |", p.curToken.Line, e.Name)
		}
		if p.curToken.Type == TokenComma {
			return e, fmt.Errorf("行 %d: 枚举 %s 成员之间仅允许 |", p.curToken.Line, e.Name)
		}

		p.skipNewLines()
		if p.curToken.Type != TokenPipe {
			break
		}
		p.nextToken()
		p.skipNewLines()
		if p.curToken.Type == TokenComment {
			return e, fmt.Errorf("行 %d: 枚举仅支持定义前整体注释, 不支持成员注释", p.curToken.Line)
		}
		if p.curToken.Type != TokenIdent {
			return e, fmt.Errorf("行 %d: 枚举 %s 的 | 后缺少成员", p.curToken.Line, e.Name)
		}
	}

	if len(e.Members) == 0 {
		return e, fmt.Errorf("行 %d: 枚举 %s 至少需要一个成员", p.curToken.Line, e.Name)
	}

	return e, nil
}

func (p *Parser) parseEnumMember() (EnumMemberRaw, error) {
	m := EnumMemberRaw{Name: p.curToken.Literal}
	p.nextToken()

	if p.curToken.Type == TokenLParen {
		p.nextToken()
		if p.curToken.Type == TokenError {
			return m, fmt.Errorf("行 %d: %s", p.curToken.Line, p.curToken.Literal)
		}
		if p.curToken.Type != TokenNumber {
			return m, fmt.Errorf("行 %d: 枚举值必须为数字", p.curToken.Line)
		}
		if len(p.curToken.Literal) > 1 && p.curToken.Literal[0] == '0' {
			return m, fmt.Errorf("行 %d: 无效枚举值 %q", p.curToken.Line, p.curToken.Literal)
		}
		parsed, err := strconv.ParseUint(p.curToken.Literal, 10, 8)
		if err != nil {
			return m, fmt.Errorf("行 %d: 无效枚举值 %q", p.curToken.Line, p.curToken.Literal)
		}
		v := uint8(parsed)
		m.Value = &v

		p.nextToken()
		if err := p.expectCurrent(TokenRParen, ")"); err != nil {
			return m, err
		}
		p.nextToken()
	}

	return m, nil
}
