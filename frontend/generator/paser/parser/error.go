package parser

import (
	"github.com/pkg/errors"
)

var (
	UnexpectEOF          = errors.New("unexpect a EOF")
	InvalidNT            = errors.New("Invalid NT . NT must be [a-z][0-9a-zA-Z_]*")
	InvalidT             = errors.New("Invalid T . T must be [A-Z][0-9a-zA-Z_]*")
	InvalidLineColon     = errors.New("Invalid LineColon. line expect a ':' after TokenType")
	InvalidLineSemicolon = errors.New("line expect a ';' after Expression")
	InvalidEOF           = errors.New("Invalid EOF. EOF must be 'EOF'")
	InvalidEpsilon       = errors.New("Invalid Epsilon. Term can contains a '#' only")
	InvalidSymbol        = errors.New("Invalid Symbol")
)
