package fa

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type CharType int64

const (
	CharTypeNormal    CharType = 0
	CharTypeEliminate CharType = 1
)

type TokenType struct {
	Name      string
	Order     int
	FuncCalls []string
}

func (t *TokenType) Equal(f *TokenType) bool {
	return t.Name == f.Name
}

type NFAEdge struct {
	CharType  CharType
	Char      rune
	FromState int
	ToState   int
}

type NFA struct {
	reader                *bufio.Reader
	row                   int
	col                   int
	curChar               rune
	enableState           int
	enableTokenTypeOrder  int
	StartState            int
	AcceptState2TokenType map[int]*TokenType
	State2Edges           [][]*NFAEdge
}

func ParseNFA(reader io.Reader) (n *NFA, e error) {
	defer Catch(&e)
	n = &NFA{
		reader:                bufio.NewReader(reader),
		AcceptState2TokenType: map[int]*TokenType{},
	}
	n.parseLex()
	return
}

func (n *NFA) errorWithLocal(err error) error {
	return errors.Wrapf(err, "[row : %d,col : %d]", n.row, n.col)
}
func (n *NFA) errorWithLocalChar(err error) error {
	return errors.Wrapf(err, "[row : %d,col : %d , Char : %c]", n.row, n.col, n.curChar)
}

func (n *NFA) nextEnableState() int {
	nextEnableState := n.enableState
	n.State2Edges = append(n.State2Edges, nil)
	n.enableState++
	return nextEnableState
}

func (n *NFA) nextEnableTokenTypeOrder() int {
	nextEnableTokenTypeOrder := n.enableTokenTypeOrder
	n.enableTokenTypeOrder++
	return nextEnableTokenTypeOrder
}

// lex ::= (line)*  EOF
func (n *NFA) parseLex() {
	var (
		err            error
		lineStartState int
	)
	n.StartState = n.nextEnableState()
	if err = n.nextCharWithEOF(); err != nil {
		return
	}
	for {
		lineStartState, _, err = n.parseLine()
		n.addEliminateEdge(n.StartState, lineStartState)
		if err == io.EOF {
			return
		}
	}
}

func (n *NFA) nextChar() {
	n.nextCharWithWhite()
	for IsWhite(n.curChar) {
		n.nextCharWithWhite()
	}
}

func (n *NFA) nextCharWithEOF() (err error) {
	n.curChar, _, err = n.reader.ReadRune()
	if err == io.EOF {
		return
	}
	Assert(err == nil, n.errorWithLocal(err))
	for IsWhite(n.curChar) {
		n.curChar, _, err = n.reader.ReadRune()
		if err == io.EOF {
			return err
		}
		Assert(err == nil, n.errorWithLocal(err))
	}
	return
}

func (n *NFA) nextCharWithWhite() {
	var err error
	n.curChar, _, err = n.reader.ReadRune()
	Assert(err != io.EOF, n.errorWithLocal(UnexpectEOF))
	Assert(err == nil, n.errorWithLocal(err))
}

// line ::= TokenType ':'  expr ';' funcCalls
func (n *NFA) parseLine() (startState, endState int, err error) {
	tokenType := n.parseTokenType()

	// ':'
	Assert(n.curChar == ':', n.errorWithLocalChar(InvalidLineColon))
	n.nextChar()

	// expr
	startState, endState = n.parseExpr()

	// ';'
	Assert(n.curChar == ';', n.errorWithLocalChar(InvalidLineSemicolon))

	var funcCalls []string
	if err = n.nextCharWithEOF(); err == nil {
		// funcCalls
		funcCalls, err = n.parseFuncCalls()
	}

	n.row++
	n.col = 0
	n.AcceptState2TokenType[endState] = &TokenType{
		Name:      tokenType,
		Order:     n.nextEnableTokenTypeOrder(),
		FuncCalls: funcCalls,
	}
	return
}

// tokenType ::= NAME
func (n *NFA) parseTokenType() (tokenType string) {
	tokenType = n.parseName()
	n.col++
	return
}

// NAME ::= [a-zA-Z_][0-9a-zA-Z_]*
func (n *NFA) parseName() string {
	Assert(IsNameStart(n.curChar), errors.Wrapf(InvalidTokenType, ""))
	sb := &strings.Builder{}
	sb.WriteRune(n.curChar)
	n.nextChar()
	for IsNameContinue(n.curChar) {
		sb.WriteRune(n.curChar)
		n.nextChar()
	}
	return sb.String()
}

// expr ::= term ('|' term)*
func (n *NFA) parseExpr() (startState, entState int) {
	startState = n.nextEnableState()
	entState = n.nextEnableState()

	termStartState, termEndState := n.parseTerm()
	n.addEliminateEdge(startState, termStartState)
	n.addEliminateEdge(termEndState, entState)

	for n.curChar == '|' {
		n.nextChar()
		termStartState, termEndState = n.parseTerm()
		n.addEliminateEdge(startState, termStartState)
		n.addEliminateEdge(termEndState, entState)
	}
	return
}

func (n *NFA) addEliminateEdge(startState, endState int) {
	n.State2Edges[startState] = append(n.State2Edges[startState], &NFAEdge{
		CharType:  CharTypeEliminate,
		Char:      0,
		FromState: startState,
		ToState:   endState,
	})
}

func (n *NFA) addNormalEdge(startState, endState int, char rune) {
	n.State2Edges[startState] = append(n.State2Edges[startState], &NFAEdge{
		CharType:  CharTypeNormal,
		Char:      char,
		FromState: startState,
		ToState:   endState,
	})
}

// term ::= piece+
func (n *NFA) parseTerm() (startState, endState int) {
	var curState int
	startState, curState = n.parsePiece()
	for IsTermStart(n.curChar) {
		pieceStartState, pieceEndState := n.parsePiece()
		n.addEliminateEdge(curState, pieceStartState)
		curState = pieceEndState
	}
	endState = curState
	return
}

// piece ::= factor QUANTIFIER?
// QUANTIFIER ::= [*?+]
func (n *NFA) parsePiece() (startState, endState int) {
	startState = n.nextEnableState()
	endState = n.nextEnableState()
	factorStartState, factorEndState := n.parseFactor()
	n.addEliminateEdge(startState, factorStartState)
	n.addEliminateEdge(factorEndState, endState)
	switch n.curChar {
	case '+':
		n.addEliminateEdge(factorEndState, factorStartState)
		n.nextChar()
	case '?':
		n.addEliminateEdge(startState, endState)
		n.nextChar()
	case '*':
		n.addEliminateEdge(factorEndState, factorStartState)
		n.addEliminateEdge(startState, endState)
		n.nextChar()
	}
	return
}

// factor ::= string | '(' expr ')'
func (n *NFA) parseFactor() (startState, endState int) {
	if n.curChar == '(' {
		n.nextChar()
		startState, endState = n.parseExpr()
		Assert(n.curChar == ')', n.errorWithLocalChar(InvalidFactorRParen))
		n.nextChar()
	} else {
		startState, endState = n.parseString()
	}
	return
}

// string ::= "'" Char+ "'"
// Char   ::=  SAFE_CHAR | ESC_CHAR
// SAFE_CHAR ::= ~['\\]
// ESC_CHAR ::= '\' ['\\bfnrt]
func (n *NFA) parseString() (startState, endState int) {
	Assert(n.curChar == '\'', n.errorWithLocalChar(InvalidStringQuote))
	startState = n.nextEnableState()
	curState := startState
	n.nextCharWithWhite()
	Assert(n.curChar != '\'', n.errorWithLocalChar(InvalidStringChar))
	for n.curChar != '\'' {
		if n.curChar == '\\' { //esc Char
			n.nextCharWithWhite()
			n.curChar = n.escChar(n.curChar)
		}
		nextState := n.nextEnableState()
		n.addNormalEdge(curState, nextState, n.curChar)
		curState = nextState
		n.nextCharWithWhite()
	}
	// '\''
	endState = curState
	n.nextChar()
	return startState, endState
}

// ESC_CHAR ::= '\' ['\\bfnrt]
func (n *NFA) escChar(char rune) rune {
	//return char == '\'' || char == '\\' || char == '/' ||
	switch char {
	case '\'':
		return 39
	case '\\':
		return 92
	case 'b':
		return 8
	case 'f':
		return 12
	case 'n':
		return 10
	case 'r':
		return 13
	case 't':
		return 9
	default:
		Assert(false, n.errorWithLocalChar(InvalidESCChar))
		return 0
	}
}

// funcCalls ::= ('{{' NAME (',' NAME)* '}}')?
func (n *NFA) parseFuncCalls() (funcCalls []string, err error) {
	if n.curChar != '{' {
		return
	}
	n.nextCharWithWhite()
	Assert('{' == n.curChar, n.errorWithLocalChar(InvalidFunCalls))

	n.nextChar()

	name := n.parseName()
	funcCalls = append(funcCalls, name)
	for n.curChar == ',' {
		n.nextChar()
		name = n.parseName()
		funcCalls = append(funcCalls, name)
	}
	Assert('}' == n.curChar, n.errorWithLocalChar(InvalidFunCalls))
	n.nextCharWithWhite()
	Assert('}' == n.curChar, n.errorWithLocalChar(InvalidFunCalls))

	err = n.nextCharWithEOF()
	return
}

func IsTermStart(char rune) bool {
	return char == '\'' || char == '('
}

// [a-zA-Z_]
func IsNameStart(char rune) bool {
	return char == '_' ||
		char > 'a' || char < 'z' ||
		char > 'A' || char < 'Z'
}

// [0-9a-zA-Z_]*
func IsNameContinue(char rune) bool {
	return char == '_' ||
		(char > '0' && char < '9') ||
		(char > 'a' && char < 'z') ||
		(char > 'A' && char < 'Z')
}

func IsWhite(char rune) bool {
	return char == ' ' || char == '\t' || char == '\v' || char == '\r' || char == '\n' || char == '\f'
}
