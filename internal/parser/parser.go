package parser

import (
	"fmt"
	"strconv"

	"sb/internal/ast"
	"sb/internal/lexer"
)

// Parser 负责把 token 流转换为 AST。
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

func New(l *lexer.Lexer) *Parser {
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
	for p.curToken.Type == lexer.TokenNewLine {
		p.nextToken()
	}
}

func (p *Parser) expectCurrent(tt lexer.TokenType, name string) error {
	if p.curToken.Type == tt {
		return nil
	}
	return fmt.Errorf("行 %d: 期望 %s, 实际 %q", p.curToken.Line, name, p.curToken.Literal)
}

func (p *Parser) ParseSchema() (*ast.Schema, error) {
	s := &ast.Schema{}
	pendingNote := ""

	for p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewLine {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenError {
			return nil, fmt.Errorf("行 %d: %s", p.curToken.Line, p.curToken.Literal)
		}

		if p.curToken.Type == lexer.TokenComment {
			if pendingNote == "" {
				pendingNote = p.curToken.Literal
			} else {
				pendingNote += "\n" + p.curToken.Literal
			}
			p.nextToken()
			if p.curToken.Type == lexer.TokenEOF {
				return nil, fmt.Errorf("行 %d: 注释行必须以换行结束", p.curToken.Line)
			}
			continue
		}

		if p.curToken.Type != lexer.TokenIdent {
			return nil, fmt.Errorf("行 %d: 未预期 token %q", p.curToken.Line, p.curToken.Literal)
		}

		note := pendingNote
		pendingNote = ""
		if p.peekToken.Type == lexer.TokenLBrace {
			st, err := p.parseStruct(note)
			if err != nil {
				return nil, err
			}
			s.Structs = append(s.Structs, st)
			continue
		}

		if p.peekToken.Type == lexer.TokenAssign {
			enum, err := p.parseEnum(note)
			if err != nil {
				return nil, err
			}
			s.Enums = append(s.Enums, enum)
			continue
		}

		if p.peekToken.Type == lexer.TokenDot || p.peekToken.Type == lexer.TokenLParen {
			api, err := p.parseAPI(note)
			if err != nil {
				return nil, err
			}
			s.APIs = append(s.APIs, api)
			continue
		}

		if p.peekToken.Type == lexer.TokenPipe {
			return nil, fmt.Errorf("行 %d: 枚举定义必须包含 =", p.curToken.Line)
		}

		return nil, fmt.Errorf("行 %d: 无法识别定义 %q", p.curToken.Line, p.curToken.Literal)
	}

	s.Note = pendingNote
	return s, nil
}

func (p *Parser) parseStruct(note string) (ast.Struct, error) {
	st := ast.Struct{Name: p.curToken.Literal, Note: note}

	p.nextToken() // name 后一个 token
	if err := p.expectCurrent(lexer.TokenLBrace, "{"); err != nil {
		return st, err
	}

	p.nextToken()
	p.skipNewLines()
	for p.curToken.Type != lexer.TokenRBrace {
		if p.curToken.Type == lexer.TokenEOF {
			return st, fmt.Errorf("行 %d: 结构体 %s 缺少 }", p.curToken.Line, st.Name)
		}
		if p.curToken.Type == lexer.TokenNewLine {
			p.nextToken()
			p.skipNewLines()
			continue
		}
		if p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			if p.curToken.Type == lexer.TokenEOF {
				return st, fmt.Errorf("行 %d: 注释行必须以换行结束", p.curToken.Line)
			}
			p.skipNewLines()
			continue
		}
		if p.curToken.Type == lexer.TokenComma {
			commaLine := p.curToken.Line
			p.nextToken()
			if p.curToken.Type == lexer.TokenNewLine {
				p.skipNewLines()
				continue
			}
			if p.curToken.Line == commaLine {
				return st, fmt.Errorf("行 %d: 结构体 %s 逗号必须独占一行", p.curToken.Line, st.Name)
			}
			continue
		}
		if p.curToken.Type != lexer.TokenIdent {
			return st, fmt.Errorf("行 %d: 结构体字段名非法 %q", p.curToken.Line, p.curToken.Literal)
		}

		fieldLine := p.curToken.Line
		field, err := p.parseField()
		if err != nil {
			return st, err
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
	if p.curToken.Type == lexer.TokenComma {
		if p.curToken.Line == fieldLine {
			return fmt.Errorf("行 %d: 结构体 %s 逗号必须独占一行", p.curToken.Line, structName)
		}
		return nil
	}
	if p.curToken.Type == lexer.TokenComment && p.curToken.Line == fieldLine {
		p.nextToken()
	}
	if p.curToken.Type == lexer.TokenNewLine || p.curToken.Type == lexer.TokenRBrace || p.curToken.Type == lexer.TokenEOF {
		return nil
	}
	if p.curToken.Line == fieldLine {
		return fmt.Errorf("行 %d: 结构体 %s 字段必须独占一行", p.curToken.Line, structName)
	}
	return nil
}

func (p *Parser) parseField() (ast.Field, error) {
	f := ast.Field{Name: p.curToken.Literal}
	line := p.curToken.Line
	p.nextToken()

	if p.curToken.Type == lexer.TokenNewLine || p.curToken.Line != line {
		f.Type = ast.TypeRef{Name: f.Name}
		f.Name = ""
		f.Embedded = true
		return f, nil
	}

	if p.curToken.Type == lexer.TokenIdent || p.curToken.Type == lexer.TokenLBracket {
		t, err := p.parseType(false)
		if err != nil {
			return f, err
		}
		f.Type = t
	} else if p.curToken.Type == lexer.TokenString && p.curToken.Line == line {
		return f, fmt.Errorf("行 %d: 字段 %s 缺少类型, 不允许直接写 tag", p.curToken.Line, f.Name)
	} else if p.curToken.Type == lexer.TokenComma || p.curToken.Type == lexer.TokenComment || p.curToken.Type == lexer.TokenRBrace {
		f.Type = ast.TypeRef{Name: f.Name}
		f.Name = ""
		f.Embedded = true
	} else {
		return f, fmt.Errorf("行 %d: 字段 %s 缺少类型", p.curToken.Line, f.Name)
	}

	if p.curToken.Type == lexer.TokenString && p.curToken.Line == line {
		f.Tag = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == lexer.TokenComment && p.curToken.Line == line {
		f.Note = p.curToken.Literal
		p.nextToken()
	}

	return f, nil
}

func (p *Parser) parseEnum(note string) (ast.Enum, error) {
	e := ast.Enum{Name: p.curToken.Literal, Note: note}
	p.nextToken()
	if err := p.expectCurrent(lexer.TokenAssign, "="); err != nil {
		return e, fmt.Errorf("行 %d: 枚举定义必须包含 =", p.curToken.Line)
	}
	p.nextToken()

	if p.curToken.Type != lexer.TokenIdent {
		return e, fmt.Errorf("行 %d: 枚举 %s 至少需要一个成员", p.curToken.Line, e.Name)
	}

	for p.curToken.Type != lexer.TokenEOF {
		p.skipNewLines()
		if p.curToken.Type != lexer.TokenIdent {
			break
		}
		memberLine := p.curToken.Line
		member, err := p.parseEnumMember()
		if err != nil {
			return e, err
		}
		e.Members = append(e.Members, member)

		if p.curToken.Type == lexer.TokenIdent && p.curToken.Line == memberLine {
			return e, fmt.Errorf("行 %d: 枚举 %s 成员之间缺少 |", p.curToken.Line, e.Name)
		}
		if p.curToken.Type == lexer.TokenComma {
			return e, fmt.Errorf("行 %d: 枚举 %s 成员之间仅允许 |", p.curToken.Line, e.Name)
		}

		p.skipNewLines()
		if p.curToken.Type != lexer.TokenPipe {
			break
		}
		p.nextToken()
		p.skipNewLines()
		if p.curToken.Type != lexer.TokenIdent {
			return e, fmt.Errorf("行 %d: 枚举 %s 的 | 后缺少成员", p.curToken.Line, e.Name)
		}
	}

	if len(e.Members) == 0 {
		return e, fmt.Errorf("行 %d: 枚举 %s 至少需要一个成员", p.curToken.Line, e.Name)
	}

	return e, nil
}

func (p *Parser) parseEnumMember() (ast.EnumMemberRaw, error) {
	m := ast.EnumMemberRaw{Name: p.curToken.Literal}
	line := p.curToken.Line
	p.nextToken()

	if p.curToken.Type == lexer.TokenLParen {
		p.nextToken()
		if p.curToken.Type != lexer.TokenNumber {
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
		if err := p.expectCurrent(lexer.TokenRParen, ")"); err != nil {
			return m, err
		}
		p.nextToken()
	}

	if p.curToken.Type == lexer.TokenComment && p.curToken.Line == line {
		m.Note = p.curToken.Literal
		p.nextToken()
	}

	return m, nil
}

func (p *Parser) parseAPI(note string) (ast.API, error) {
	api := ast.API{Note: note}
	name := p.curToken.Literal
	p.nextToken()

	if p.curToken.Type == lexer.TokenDot {
		p.nextToken()
		if p.curToken.Type != lexer.TokenIdent {
			return api, fmt.Errorf("行 %d: API 名称缺少段", p.curToken.Line)
		}
		name += "." + p.curToken.Literal
		p.nextToken()
	}
	if p.curToken.Type == lexer.TokenDot {
		return api, fmt.Errorf("行 %d: API 名称仅支持 1-2 段", p.curToken.Line)
	}
	api.Name = name

	if err := p.expectCurrent(lexer.TokenLParen, "("); err != nil {
		return api, err
	}
	p.nextToken()
	p.skipNewLines()

	for p.curToken.Type != lexer.TokenRParen {
		if p.curToken.Type == lexer.TokenEOF {
			return api, fmt.Errorf("行 %d: API %s 缺少 )", p.curToken.Line, api.Name)
		}
		p.skipNewLines()
		if p.curToken.Type == lexer.TokenRParen {
			break
		}
		if p.curToken.Type != lexer.TokenIdent {
			return api, fmt.Errorf("行 %d: API 参数名非法 %q", p.curToken.Line, p.curToken.Literal)
		}
		arg := ast.APIArg{Name: p.curToken.Literal}
		p.nextToken()

		t, err := p.parseType(false)
		if err != nil {
			return api, err
		}
		arg.Type = t
		api.Args = append(api.Args, arg)
		p.skipNewLines()

		if p.curToken.Type == lexer.TokenComma {
			p.nextToken()
			p.skipNewLines()
			if p.curToken.Type == lexer.TokenRParen {
				return api, fmt.Errorf("行 %d: API 参数列表尾逗号非法", p.curToken.Line)
			}
			continue
		}
		if p.curToken.Type != lexer.TokenRParen {
			return api, fmt.Errorf("行 %d: API 参数后需要 , 或 )", p.curToken.Line)
		}
	}

	p.nextToken()
	p.skipNewLines()
	if p.curToken.Type != lexer.TokenArrow {
		return api, fmt.Errorf("行 %d: API 定义缺少 =>", p.curToken.Line)
	}
	p.nextToken()
	p.skipNewLines()

	result, err := p.parseType(true)
	if err != nil {
		return api, err
	}
	api.Result = result

	if p.curToken.Type == lexer.TokenComment {
		return api, fmt.Errorf("行 %d: API 注释必须写在定义前, 不支持尾部注释", p.curToken.Line)
	}

	return api, nil
}

func (p *Parser) parseType(allowNil bool) (ast.TypeRef, error) {
	var t ast.TypeRef

	if p.curToken.Type == lexer.TokenIdent {
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

	if p.curToken.Type != lexer.TokenLBracket {
		return t, fmt.Errorf("行 %d: 非法类型 %q", p.curToken.Line, p.curToken.Literal)
	}

	t.IsList = true
	p.nextToken()
	p.skipNewLines()
	if p.curToken.Type != lexer.TokenIdent {
		return t, fmt.Errorf("行 %d: 数组类型缺少元素类型", p.curToken.Line)
	}
	if p.curToken.Literal == "nil" {
		return t, fmt.Errorf("行 %d: nil 不能作为数组元素类型", p.curToken.Line)
	}

	t.Name = p.curToken.Literal
	p.nextToken()
	p.skipNewLines()
	if err := p.expectCurrent(lexer.TokenRBracket, "]"); err != nil {
		return t, err
	}
	p.nextToken()
	return t, nil
}
