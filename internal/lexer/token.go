package lexer

// TokenType 表示词法单元类型。
type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenIdent
	TokenNumber
	TokenString
	TokenComment
	TokenNewLine  // \n or \r\n
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
)

// Token 是词法结果。
type Token struct {
	Type    TokenType
	Literal string
	Line    int
}
