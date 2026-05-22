package internal

import "fmt"

func (p *Parser) ParseSchema() (*Schema, error) {
	s := &Schema{}
	pendingNote := ""

	for p.curToken.Type != TokenEOF {
		if p.curToken.Type == TokenNewLine {
			p.nextToken()
			continue
		}

		if p.curToken.Type == TokenError {
			return nil, fmt.Errorf("行 %d: %s", p.curToken.Line, p.curToken.Literal)
		}

		if p.curToken.Type == TokenComment {
			if pendingNote == "" {
				pendingNote = p.curToken.Literal
			} else {
				pendingNote += "\n" + p.curToken.Literal
			}
			p.nextToken()
			if p.curToken.Type == TokenEOF {
				return nil, fmt.Errorf("行 %d: 注释行必须以换行结束", p.curToken.Line)
			}
			continue
		}

		if p.curToken.Type != TokenIdent {
			return nil, fmt.Errorf("行 %d: 未预期 token %q", p.curToken.Line, p.curToken.Literal)
		}

		note := pendingNote
		pendingNote = ""
		if p.peekToken.Type == TokenLBrace {
			st, err := p.parseStruct(note)
			if err != nil {
				return nil, err
			}
			s.Structs = append(s.Structs, st)
			continue
		}

		if p.peekToken.Type == TokenAssign {
			enum, err := p.parseEnum(note)
			if err != nil {
				return nil, err
			}
			s.Enums = append(s.Enums, enum)
			continue
		}

		if p.peekToken.Type == TokenDot || p.peekToken.Type == TokenLParen {
			api, err := p.parseAPI(note)
			if err != nil {
				return nil, err
			}
			s.APIs = append(s.APIs, api)
			continue
		}

		if p.peekToken.Type == TokenPipe {
			return nil, fmt.Errorf("行 %d: 枚举定义必须包含 =", p.curToken.Line)
		}

		return nil, fmt.Errorf("行 %d: 无法识别定义 %q", p.curToken.Line, p.curToken.Literal)
	}

	s.Note = pendingNote
	return s, nil
}
