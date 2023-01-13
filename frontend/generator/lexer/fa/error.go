package fa

import (
	"github.com/pkg/errors"
)

var (
	UnexpectEOF          = errors.New("unexpect a EOF")
	InvalidTokenType     = errors.New("invalid TokenType . must be [A-Z][0-9a-zA-Z_]*")
	InvalidLineColon     = errors.New("line expect a ':' after TokenType")
	InvalidLineSemicolon = errors.New("line expect a ';' after Expression")
	InvalidFactorRParen  = errors.New("factor expect a ')' after Expression")
	InvalidStringQuote   = errors.New("string expect a ''' around Char")
	InvalidStringChar    = errors.New("string must have a Char at least")
	InvalidESCChar       = errors.New("invalid ESCChar . must be ['\\bfnrt]")
	InvalidTerm          = errors.New("invalid Terml")
	InvalidCharRange     = errors.New("invalid CharRange. left char  must <= right")
	InvalidComposed      = errors.New("invalid composed. must be '[' (char '-' char | char)+']' ")
	InvalidFunCalls      = errors.New("invalid FuncCalls. must be '{{' NAME (',' NAME)* '}}'")
	UnknownFuncCall      = errors.New("unknown FuncCall")
)
