package lexer

import "github.com/pkg/errors"

var (
	InvalidToken = errors.New("invalid token")
)

type TokenType int

const (
	TokenType_EOF TokenType = 1
)

type Token interface {
	Text() string
	Type() TokenType
	Row() int
	Column() int
}

type token struct {
	text      []rune
	tokenType TokenType
	row       int
	column    int
}

func (t token) Text() string {
	return string(t.text)
}

func (t token) Type() TokenType {
	return t.tokenType
}

func (t token) Row() int {
	return t.row
}

func (t token) Column() int {
	return t.column
}

func NewToken(text []rune, tokenType TokenType, row int, column int) Token {
	return &token{
		text:      text,
		tokenType: tokenType,
		row:       row,
		column:    column,
	}
}

type Lexer interface {
	NextToken() (token Token, err error)
}
