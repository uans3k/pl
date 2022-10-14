package lexer

import "github.com/pkg/errors"

var (
	InvalidToken = errors.New("invalid token")
)

type TokenType interface {
	String() string
}

type TokenValue interface {
	Chars() []rune
	String() string
}

type Token struct {
	Value  TokenValue
	Type   TokenType
	Row    int
	Column int
}

type Lexer interface {
	NextToken() (token *Token, err error)
}
