package fa

import (
	"github.com/pkg/errors"
)

var (
	UnknownError         = errors.New("unknown error")
	UnexpectEOF          = errors.New("unexpect a EOF")
	InvalidTokenType     = errors.New("invalid TokenType . must be [a-zA-Z_][0-9a-zA-Z_]*")
	InvalidLineColon     = errors.New("line expect a ':' after TokenType")
	InvalidLineSemicolon = errors.New("line expect a ';' after Expression")
	InvalidFactorRParen  = errors.New("factor expect a ')' after Expression")
	InvalidStringQuote   = errors.New("string expect a ''' around Char")
	InvalidStringChar    = errors.New("string must have a Char at least")
	InvalidESCChar       = errors.New("invalid ESCChar . must be ['\\bfnrt]")
	InvalidFunCalls      = errors.New("invalid FuncCalls. must be '{{' NAME (',' NAME)* '}}'")
	UnknownFuncCall      = errors.New("unknown FuncCall")
)

func Assert(b bool, elseError error) {
	if !b {
		panic(elseError)
	}
}

func Catch(e *error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			*e = err
		} else {
			*e = errors.Wrapf(UnknownError, "value :%+v", r)
		}
	}
}
