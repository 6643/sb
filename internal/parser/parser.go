package parser

import (
	"fmt"
	"sb/internal/ast"
	"sb/internal/lexer"
	"strconv"
	"strings"
)

// Parser 语法分析器
// 采用递归下降 (Recursive Descent) 策略
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token

	// 符号表: 用于快速校验类型引用有效性
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

// ParseSchema 解析完整的 Schema 文件
// 包含两个阶段:
// 1. 语法解析: 构建 AST 结构, 收集所有定义
// 2. 类型语义分析 (resolveTypes): 校验类型引用, 展开嵌入结构体
func (p *Parser) ParseSchema() (*ast.Schema, error) {
	schema := &ast.Schema{}
	var lastNote string

	for p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenError {
			return nil, fmt.Errorf("lexing error: %s", p.curToken.Value)
		}

		// 收集注释作为下一个定义的文档
		if p.curToken.Type == lexer.TokenComment {
			if lastNote != "" {
				lastNote += "\n"
			}
			lastNote += p.curToken.Value
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIdent {
			if err := p.parseDefinition(schema, &lastNote); err != nil {
				return nil, err
			}
			continue
		}

		return nil, fmt.Errorf("line %d: unexpected token %q", p.curToken.Line, p.curToken.Value)
	}

	// 语义分析阶段
	if err := p.resolveTypes(schema); err != nil {
		return nil, err
	}
	return schema, nil
}

func (p *Parser) parseDefinition(schema *ast.Schema, lastNote *string) error {

	defer func() { *lastNote = "" }()

	note := *lastNote



	if p.peekToken.Type == lexer.TokenLBrace {

		return p.parseAndAddStruct(schema, note)

	}



	if p.isEnumDefinition() {

		return p.parseAndAddEnum(schema, note)

	}



	if p.isApiDefinition() {

		return p.parseAndAddApi(schema, note)

	}



	return fmt.Errorf("行 %d: 未预期标识符 %q", p.curToken.Line, p.curToken.Value)

}



func (p *Parser) parseAndAddStruct(schema *ast.Schema, note string) error {

	if p.isDefined(p.curToken.Value) {

		return fmt.Errorf("行 %d: %s 重复定义", p.curToken.Line, p.curToken.Value)

	}

	s, err := p.parseStruct(note)

	if err != nil {

		return err

	}

	schema.Structs = append(schema.Structs, s)

	p.structNames[s.Name] = true

	return nil

}



func (p *Parser) parseAndAddEnum(schema *ast.Schema, note string) error {

	if p.isDefined(p.curToken.Value) {

		return fmt.Errorf("行 %d: %s 重复定义", p.curToken.Line, p.curToken.Value)

	}

	e, err := p.parseEnum(note)

	if err != nil {

		return err

	}

	schema.Enums = append(schema.Enums, e)

	p.enumNames[e.Name] = true

	return nil

}



func (p *Parser) parseAndAddApi(schema *ast.Schema, note string) error {

	api, err := p.parseApi(note)

	if err != nil {

		return err

	}

	schema.Apis = append(schema.Apis, api)

	return nil

}



func (p *Parser) isDefined(name string) bool {

	return p.structNames[name] || p.enumNames[name]

}



func (p *Parser) isEnumDefinition() bool {

	return p.peekToken.Type == lexer.TokenAssign || p.peekToken.Type == lexer.TokenPipe

}



func (p *Parser) isApiDefinition() bool {

	return p.curToken.Type == lexer.TokenIdent && (p.peekToken.Type == lexer.TokenLParen || p.peekToken.Type == lexer.TokenDot)

}



func (p *Parser) parseStruct(note string) (ast.Struct, error) {

	s := ast.Struct{Name: p.curToken.Value, Note: note}

	p.nextToken() // 名称

	p.nextToken() // {



	for p.curToken.Type != lexer.TokenRBrace && p.curToken.Type != lexer.TokenEOF {

		if p.curToken.Type == lexer.TokenComment {

			p.nextToken()

			continue

		}

		if p.curToken.Type == lexer.TokenComma {

			p.nextToken() // 跳过可选逗号

			continue

		}



		field, err := p.parseStructField()

		if err != nil {

			return s, err

		}

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



	// 嵌入结构体情况: 名称实际上是类型

	if p.curToken.Line != startLine {

		f.Type = ast.Type{Name: f.Name}

		f.Name = ""

		return f, nil

	}

	

	// 普通字段情况

	if p.curToken.Type == lexer.TokenIdent || p.curToken.Type == lexer.TokenLBracket {

		f.Type = p.parseType()

		if p.curToken.Type == lexer.TokenIdent && (strings.HasPrefix(p.curToken.Value, "\"") || strings.HasPrefix(p.curToken.Value, "`")) {

			f.Tag = strings.Trim(p.curToken.Value, "\"`")

			p.nextToken()

		}

	} else {

		// 嵌入兜底

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

	if p.curToken.Type != lexer.TokenLBracket {

		t.Name = p.curToken.Value

		p.nextToken()

		return t

	}



	t.IsList = true

	p.nextToken() // [

	t.Name = p.curToken.Value

	p.nextToken() // 名称

	p.nextToken() // ]

	return t

}



func (p *Parser) parseEnum(note string) (ast.Enum, error) {

	e := ast.Enum{Name: p.curToken.Value, Note: note}

	p.nextToken() // 名称

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

		if p.curToken.Type != lexer.TokenIdent {

			break

		}



		child, err := p.parseEnumChild(&lastID, &isFirst)

		if err != nil {

			return e, err

		}

		e.Children = append(e.Children, child)



		if p.curToken.Type != lexer.TokenPipe {

			break

		}

	}

	return e, nil

}



func (p *Parser) parseEnumChild(lastID *uint8, isFirst *bool) (ast.EnumChild, error) {

	child := ast.EnumChild{Name: p.curToken.Value}

	childLine := p.curToken.Line

	p.nextToken()



	if p.curToken.Type == lexer.TokenLParen {

		p.nextToken() // (

		id, err := strconv.ParseUint(p.curToken.Value, 10, 8)

		if err != nil {

			return child, fmt.Errorf("行 %d: 无效枚举值 %q: %w", p.curToken.Line, p.curToken.Value, err)

		}

		child.ID = uint8(id)

		*lastID = child.ID

		p.nextToken() // 数值

		p.nextToken() // )

		*isFirst = false

	} else {

		if *isFirst {

			child.ID = 0

			*isFirst = false

		} else {

			if *lastID == 255 {

				return child, fmt.Errorf("行 %d: 枚举值溢出", childLine)

			}

			*lastID++

			child.ID = *lastID

		}

	}



	if p.curToken.Type == lexer.TokenComment && p.curToken.Line == childLine {

		child.Note = p.curToken.Value

		p.nextToken()

	}

	return child, nil

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

		if p.curToken.Type == lexer.TokenComma {

			p.nextToken()

		}

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

	if err := p.resolveStructFields(s); err != nil {

		return err

	}

	if err := p.resolveApiArgs(s); err != nil {

		return err

	}

	return p.expandEmbeddedStructs(s)

}



func (p *Parser) resolveStructFields(s *ast.Schema) error {

	for i := range s.Structs {

		for j := range s.Structs[i].Fields {

			if err := p.resolveType(&s.Structs[i].Fields[j].Type); err != nil {

				return fmt.Errorf("结构体 %s 字段 %s: %w", s.Structs[i].Name, s.Structs[i].Fields[j].Name, err)

			}

		}

	}

	return nil

}



func (p *Parser) resolveApiArgs(s *ast.Schema) error {

	for i := range s.Apis {

		for j := range s.Apis[i].Args {

			if err := p.resolveType(&s.Apis[i].Args[j].Type); err != nil {

				return fmt.Errorf("api %s 参数 %s: %w", s.Apis[i].Name, s.Apis[i].Args[j].Name, err)

			}

		}

		if err := p.resolveType(&s.Apis[i].Result); err != nil {

			return fmt.Errorf("api %s 结果: %w", s.Apis[i].Name, err)

		}

	}

	return nil

}



func (p *Parser) resolveType(t *ast.Type) error {

	if t.Name == "nil" || isBaseType(t.Name) {

		t.Kind = ast.KindBase

		return nil

	}

	if p.structNames[t.Name] {

		t.Kind = ast.KindStruct

		return nil

	}

	if p.enumNames[t.Name] {

		t.Kind = ast.KindEnum

		return nil

	}

	return fmt.Errorf("未定义类型: %s", t.Name)

}



func isBaseType(name string) bool {

	switch name {

	case "i8", "u8", "i16", "u16", "i32", "u32", "i64", "u64", 

		 "f32", "f64", "bool", "text", "bin":

		return true

	}

	return false

}



func (p *Parser) expandEmbeddedStructs(s *ast.Schema) error {

	structMap := make(map[string]ast.Struct)

	for _, st := range s.Structs {

		structMap[st.Name] = st

	}



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

		return nil, fmt.Errorf("检测到循环嵌入: %s", rootName)

	}

	visited[rootName] = true

	defer func() { visited[rootName] = false }()



	var result []ast.StructField

	for _, f := range fields {

		if f.Name != "" {

			result = append(result, f)

			continue

		}



		base, ok := structMap[f.Type.Name]

		if !ok {

			// 由于 resolveTypes 理论上不可达

			return nil, fmt.Errorf("未找到嵌入结构体 %s", f.Type.Name)

		}



		expanded, err := p.expandFields(base.Fields, structMap, visited, f.Type.Name)

		if err != nil {

			return nil, err

		}

		result = append(result, expanded...)

	}

	return result, nil

}
