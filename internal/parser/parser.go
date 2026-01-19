package parser

import (
	"fmt"
	"sb/internal/ast"
	"sb/internal/lexer"
	"strconv"
	"strings"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	
	structNames map[string]bool
	enumNames   map[string]bool
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:           l,
		structNames: make(map[string]bool),
		enumNames:   make(map[string]bool),
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseSchema() (*ast.Schema, error) {
	schema := &ast.Schema{}
	var lastNote string

	for p.curToken.Type != lexer.TokenEOF {
		switch p.curToken.Type {
		case lexer.TokenComment:
			if lastNote != "" {
				lastNote += "\n"
			}
			lastNote += p.curToken.Value
			p.nextToken()
		case lexer.TokenIdent:
			if p.peekToken.Type == lexer.TokenLBrace {
				if p.structNames[p.curToken.Value] || p.enumNames[p.curToken.Value] {
					return nil, fmt.Errorf("line %d: %s redefined", p.curToken.Line, p.curToken.Value)
				}
				s, err := p.parseStruct(lastNote)
				if err != nil { return nil, err }
				schema.Structs = append(schema.Structs, s)
				p.structNames[s.Name] = true
				lastNote = ""
			} else if p.isEnumDefinition() {
				if p.structNames[p.curToken.Value] || p.enumNames[p.curToken.Value] {
					return nil, fmt.Errorf("line %d: %s redefined", p.curToken.Line, p.curToken.Value)
				}
				e, err := p.parseEnum(lastNote)
				if err != nil { return nil, err }
				schema.Enums = append(schema.Enums, e)
				p.enumNames[e.Name] = true
				lastNote = ""
			} else if p.isApiDefinition() {
				api, err := p.parseApi(lastNote)
				if err != nil { return nil, err }
				schema.Apis = append(schema.Apis, api)
				lastNote = ""
			} else {
				return nil, fmt.Errorf("line %d: unexpected identifier %q", p.curToken.Line, p.curToken.Value)
			}
		default:
			p.nextToken()
		}
	}

	if err := p.resolveTypes(schema); err != nil {
		return nil, err
	}
	return schema, nil
}

func (p *Parser) isEnumDefinition() bool {
	return p.peekToken.Type == lexer.TokenAssign || p.peekToken.Type == lexer.TokenPipe
}

func (p *Parser) isApiDefinition() bool {
	return p.curToken.Type == lexer.TokenIdent && (p.peekToken.Type == lexer.TokenLParen || p.peekToken.Type == lexer.TokenDot)
}

func (p *Parser) parseStruct(note string) (ast.Struct, error) {
	s := ast.Struct{Name: p.curToken.Value, Note: note}
	p.nextToken() // name
	p.nextToken() // {

	for p.curToken.Type != lexer.TokenRBrace && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}
		field, err := p.parseStructField()
		if err != nil { return s, err }
		s.Fields = append(s.Fields, field)
	}
	p.nextToken() // }
	return s, nil
}

func (p *Parser) parseStructField() (ast.StructField, error) {
	var f ast.StructField
	startLine := p.curToken.Line

	f.Name = p.curToken.Value
	p.nextToken()

	if (p.curToken.Type == lexer.TokenIdent || p.curToken.Type == lexer.TokenLBracket) && p.curToken.Line == startLine {
		// Normal field: Name Type
		f.Type = p.parseType()

		// Handle tag (simple string/ident check)
		if p.curToken.Type == lexer.TokenIdent && (strings.HasPrefix(p.curToken.Value, "\"") || strings.HasPrefix(p.curToken.Value, "`")) {
			f.Tag = strings.Trim(p.curToken.Value, "\"`")
			p.nextToken()
		}
	} else {
		// Embedded struct: Name is actually the Type
		f.Type = ast.Type{Name: f.Name}
		f.Name = ""
	}

	if p.curToken.Type == lexer.TokenComment && p.curToken.Line == startLine {
		f.Note = p.curToken.Value
		p.nextToken()
	}

	return f, nil
}

func (p *Parser) parseType() ast.Type {
	var t ast.Type
	if p.curToken.Type == lexer.TokenLBracket {
		t.IsList = true
		p.nextToken() // [
		t.Name = p.curToken.Value
		p.nextToken() // name
		p.nextToken() // ]
	} else {
		t.Name = p.curToken.Value
		p.nextToken()
	}
	return t
}

func (p *Parser) parseEnum(note string) (ast.Enum, error) {
	e := ast.Enum{Name: p.curToken.Value, Note: note}
	p.nextToken() // name
	if p.curToken.Type == lexer.TokenAssign {
		p.nextToken() // =
	}

	var lastID uint8 = 0
	isFirst := true

	for p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenPipe {
			p.nextToken()
			continue
		}
		if p.curToken.Type != lexer.TokenIdent { break }

		child := ast.EnumChild{Name: p.curToken.Value}
		childLine := p.curToken.Line
		p.nextToken()

		if p.curToken.Type == lexer.TokenLParen {
			p.nextToken() // (
			id, err := strconv.ParseUint(p.curToken.Value, 10, 8)
			if err != nil {
				return e, fmt.Errorf("line %d: invalid enum value %q: %w", p.curToken.Line, p.curToken.Value, err)
			}
			child.ID = uint8(id)
			lastID = child.ID
			p.nextToken() // num
			p.nextToken() // )
			isFirst = false
		} else {
			if isFirst {
				child.ID = 0
				isFirst = false
			} else {
				if lastID == 255 {
					return e, fmt.Errorf("line %d: enum value overflow", p.curToken.Line)
				}
				lastID++
				child.ID = lastID
			}
		}

		if p.curToken.Type == lexer.TokenComment && p.curToken.Line == childLine {
			child.Note = p.curToken.Value
			p.nextToken()
		}
		e.Children = append(e.Children, child)
		if p.curToken.Type != lexer.TokenPipe { break }
	}
	return e, nil
}

func (p *Parser) parseApi(note string) (ast.Api, error) {
	api := ast.Api{Note: note}
	apiLine := p.curToken.Line
	name := p.curToken.Value
	p.nextToken()
	for p.curToken.Type == lexer.TokenDot {
		p.nextToken()
		name += "." + p.curToken.Value
		p.nextToken()
	}
	api.Name = name

	p.nextToken() // (
	for p.curToken.Type != lexer.TokenRParen && p.curToken.Type != lexer.TokenEOF {
		arg := ast.ApiArg{Name: p.curToken.Value}
		p.nextToken()
		arg.Type = p.parseType()
		api.Args = append(api.Args, arg)
		if p.curToken.Type == lexer.TokenComma { p.nextToken() }
	}
	p.nextToken() // )
	p.nextToken() // =>

	if p.curToken.Type == lexer.TokenIdent || p.curToken.Type == lexer.TokenLBracket {
		api.Result = p.parseType()
	} else if p.curToken.Value == "nil" {
		api.Result = ast.Type{Name: "nil"}
		p.nextToken()
	}

	if p.curToken.Type == lexer.TokenComment && p.curToken.Line == apiLine {
		api.Note = p.curToken.Value
		p.nextToken()
	}

	return api, nil
}

func (p *Parser) resolveTypes(s *ast.Schema) error {
	baseTypes := map[string]bool{
		"i8": true, "u8": true, "i16": true, "u16": true,
		"i32": true, "u32": true, "i64": true, "u64": true,
		"f32": true, "f64": true, "bool": true, "text": true, "bin": true,
	}

	resolve := func(t *ast.Type) error {
		if t.Name == "nil" {
			t.Kind = ast.KindBase
			return nil
		}
		if baseTypes[t.Name] {
			t.Kind = ast.KindBase
		} else if p.structNames[t.Name] {
			t.Kind = ast.KindStruct
		} else if p.enumNames[t.Name] {
			t.Kind = ast.KindEnum
		} else {
			return fmt.Errorf("undefined type: %s", t.Name)
		}
		return nil
	}

	for i := range s.Structs {
		for j := range s.Structs[i].Fields {
			if err := resolve(&s.Structs[i].Fields[j].Type); err != nil {
				return fmt.Errorf("struct %s field %s: %w", s.Structs[i].Name, s.Structs[i].Fields[j].Name, err)
			}
		}
	}
	for i := range s.Apis {
		for j := range s.Apis[i].Args {
			if err := resolve(&s.Apis[i].Args[j].Type); err != nil {
				return fmt.Errorf("api %s arg %s: %w", s.Apis[i].Name, s.Apis[i].Args[j].Name, err)
			}
		}
		if err := resolve(&s.Apis[i].Result); err != nil {
			return fmt.Errorf("api %s result: %w", s.Apis[i].Name, err)
		}
	}

	// Resolve embedded structs
	structMap := make(map[string]ast.Struct)
	for _, st := range s.Structs {
		structMap[st.Name] = st
	}

	// Simple recursive expansion for each struct
	visited := make(map[string]bool)
	for i := range s.Structs {
		clear(visited)
		expanded, err := p.expandFields(s.Structs[i].Fields, structMap, visited, s.Structs[i].Name)
		if err != nil {
			return err
		}
		s.Structs[i].Fields = expanded
	}
	return nil
}

func (p *Parser) expandFields(fields []ast.StructField, structMap map[string]ast.Struct, visited map[string]bool, rootName string) ([]ast.StructField, error) {
	if visited[rootName] {
		return nil, fmt.Errorf("circular embedding detected: %s", rootName)
	}
	visited[rootName] = true
	defer func() { visited[rootName] = false }()

	var result []ast.StructField
	for _, f := range fields {
		if f.Name == "" {
			if base, ok := structMap[f.Type.Name]; ok {
				// Recursively expand base fields too
				expanded, err := p.expandFields(base.Fields, structMap, visited, f.Type.Name)
				if err != nil {
					return nil, err
				}
				result = append(result, expanded...)
			}
		} else {
			result = append(result, f)
		}
	}
	return result, nil
}