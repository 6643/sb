package internal

import "fmt"

func (p *Parser) parseStruct(note string) (Struct, error) {
	st := Struct{Name: p.curToken.Literal, Note: note}
	pendingFieldNote := ""

	p.nextToken() // name 后一个 token
	if err := p.expectCurrent(TokenLBrace, "{"); err != nil {
		return st, err
	}

	p.nextToken()
	p.skipNewLines()
	for p.curToken.Type != TokenRBrace {
		if p.curToken.Type == TokenEOF {
			return st, fmt.Errorf("行 %d: 结构体 %s 缺少 }", p.curToken.Line, st.Name)
		}
		if p.curToken.Type == TokenNewLine {
			p.nextToken()
			p.skipNewLines()
			continue
		}
		if p.curToken.Type == TokenComment {
			if pendingFieldNote == "" {
				pendingFieldNote = p.curToken.Literal
			} else {
				pendingFieldNote += "\n" + p.curToken.Literal
			}
			p.nextToken()
			if p.curToken.Type == TokenEOF {
				return st, fmt.Errorf("行 %d: 注释行必须以换行结束", p.curToken.Line)
			}
			p.skipNewLines()
			continue
		}
		if p.curToken.Type == TokenComma {
			commaLine := p.curToken.Line
			p.nextToken()
			if p.curToken.Type == TokenNewLine {
				p.skipNewLines()
				continue
			}
			if p.curToken.Line == commaLine {
				return st, fmt.Errorf("行 %d: 结构体 %s 逗号必须独占一行", p.curToken.Line, st.Name)
			}
			continue
		}
		if p.curToken.Type != TokenIdent && p.curToken.Type != TokenEllipsis {
			return st, fmt.Errorf("行 %d: 结构体字段名非法 %q", p.curToken.Line, p.curToken.Literal)
		}

		fieldLine := p.curToken.Line
		field, err := p.parseField()
		if err != nil {
			return st, err
		}
		if pendingFieldNote != "" {
			field.Note = pendingFieldNote
			pendingFieldNote = ""
		}
		st.Fields = append(st.Fields, field)

		if err := p.ensureStructFieldLineEnd(st.Name, fieldLine); err != nil {
			return st, err
		}
		p.skipNewLines()
	}

	p.nextToken() // 消费 }
	return st, nil
}

func (p *Parser) ensureStructFieldLineEnd(structName string, fieldLine int) error {
	if p.curToken.Type == TokenComma {
		if p.curToken.Line == fieldLine {
			return fmt.Errorf("行 %d: 结构体 %s 逗号必须独占一行", p.curToken.Line, structName)
		}
		return nil
	}
	if p.curToken.Type == TokenComment && p.curToken.Line == fieldLine {
		return fmt.Errorf("行 %d: 结构体字段注释必须独占一行并写在字段前", p.curToken.Line)
	}
	if p.curToken.Type == TokenNewLine || p.curToken.Type == TokenRBrace || p.curToken.Type == TokenEOF {
		return nil
	}
	if p.curToken.Line == fieldLine {
		return fmt.Errorf("行 %d: 结构体 %s 字段必须独占一行", p.curToken.Line, structName)
	}
	return nil
}

func (p *Parser) parseField() (Field, error) {
	if p.curToken.Type == TokenEllipsis {
		line := p.curToken.Line
		p.nextToken()
		if p.curToken.Type != TokenIdent || p.curToken.Line != line {
			return Field{}, fmt.Errorf("行 %d: 嵌入字段必须写为 ...TypeName", line)
		}
		f := Field{Type: TypeRef{Name: p.curToken.Literal}, Embedded: true}
		p.nextToken()
		return f, nil
	}

	f := Field{Name: p.curToken.Literal}
	line := p.curToken.Line
	p.nextToken()

	if p.curToken.Type == TokenNewLine || p.curToken.Line != line {
		return f, fmt.Errorf("行 %d: 字段 %s 缺少类型", line, f.Name)
	}

	if p.curToken.Type == TokenIdent || p.curToken.Type == TokenLBracket {
		t, err := p.parseType(false)
		if err != nil {
			return f, err
		}
		f.Type = t
	} else if p.curToken.Type == TokenString && p.curToken.Line == line {
		return f, fmt.Errorf("行 %d: 字段 %s 缺少类型, 不允许直接写 tag", p.curToken.Line, f.Name)
	} else if p.curToken.Type == TokenComment && p.curToken.Line == line {
		return f, fmt.Errorf("行 %d: 结构体字段注释必须独占一行并写在字段前", p.curToken.Line)
	} else if p.curToken.Type == TokenComma || p.curToken.Type == TokenRBrace {
		return f, fmt.Errorf("行 %d: 字段 %s 缺少类型", line, f.Name)
	} else {
		return f, fmt.Errorf("行 %d: 字段 %s 缺少类型", p.curToken.Line, f.Name)
	}

	if p.curToken.Type == TokenString && p.curToken.Line == line {
		f.Tag = p.curToken.Literal
		p.nextToken()
	}

	return f, nil
}
