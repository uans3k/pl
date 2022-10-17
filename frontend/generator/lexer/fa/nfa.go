package fa

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/uans3k/pl/infra"
	"io"
	"strings"
)

type CharType int64

var (
	builtinInFuncCalls = infra.NewSet("Hidden", "Row")
)

type Char interface {
	Equal(t Char) bool
	Key() string
}

var CharEliminate Char = &charEliminate{}

type charEliminate struct {
}

func (c *charEliminate) Key() string {
	return ""
}

func (c *charEliminate) Equal(t Char) bool {
	_, ok := t.(*charEliminate)
	return ok
}

type CharSingle struct {
	Char rune
}

func (c *CharSingle) Key() string {
	return string(c.Char)
}

func NewCharSingle(char rune) Char {
	return &CharSingle{Char: char}
}

func (c *CharSingle) Equal(t Char) bool {
	if e, ok := t.(*CharSingle); ok {
		return e.Char == c.Char
	}
	return false
}

type CharRange struct {
	LeftChar  rune
	RightChar rune
}

func (c *CharRange) Key() string {
	return fmt.Sprintf("%c-%c", c.LeftChar, c.RightChar)
}

func NewCharRange(l, r rune) Char {
	return &CharRange{LeftChar: l, RightChar: r}
}

func (c *CharRange) Equal(t Char) bool {
	if e, ok := t.(*CharRange); ok {
		return e.LeftChar == c.LeftChar && e.RightChar == c.RightChar
	}
	return false
}

type TokenType struct {
	Name      string
	Order     int
	FuncCalls []string
}

func (t *TokenType) Equal(f *TokenType) bool {
	return t.Name == f.Name
}

type NFAEdge struct {
	Char      Char
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
	Assert(builtinInFuncCalls.ContainsAll(funcCalls), errors.Wrapf(UnknownFuncCall, "actual %+v", funcCalls))

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
		Char:      CharEliminate,
		FromState: startState,
		ToState:   endState,
	})
}

func (n *NFA) addCharSingleEdge(startState, endState int, char rune) {
	n.State2Edges[startState] = append(n.State2Edges[startState], &NFAEdge{
		Char:      NewCharSingle(char),
		FromState: startState,
		ToState:   endState,
	})
}

func (n *NFA) addCharRangeEdge(startState, endState int, leftChar, rightChar rune) {
	Assert(leftChar <= rightChar, n.errorWithLocalChar(InvalidCharRange))
	n.State2Edges[startState] = append(n.State2Edges[startState], &NFAEdge{
		Char:      NewCharRange(leftChar, rightChar),
		FromState: startState,
		ToState:   endState,
	})
}

// term ::= piece+
func (n *NFA) parseTerm() (startState, endState int) {
	var curState int
	Assert(IsTermStart(n.curChar), n.errorWithLocalChar(InvalidTerm))
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

// factor ::= string | composed |'(' expr ')'
func (n *NFA) parseFactor() (startState, endState int) {
	if n.curChar == '(' {
		n.nextChar()
		startState, endState = n.parseExpr()
		Assert(n.curChar == ')', n.errorWithLocalChar(InvalidFactorRParen))
		n.nextChar()
	} else if IsComposedStart(n.curChar) {
		startState, endState = n.parseComposed()
	} else {
		startState, endState = n.parseString()
	}
	return
}

// composed ::= '[' (char '-' char | char)+']'
func (n *NFA) parseComposed() (startState, endState int) {
	Assert(IsComposedStart(n.curChar), n.errorWithLocalChar(InvalidComposed))
	n.nextCharWithWhite()
	startState = n.nextEnableState()
	endState = n.nextEnableState()
	Assert(n.curChar != ']', n.errorWithLocalChar(InvalidComposed))
	for n.curChar != ']' {
		leftChar := n.parseChar()
		n.nextCharWithWhite()
		if n.curChar == '-' {
			n.nextCharWithWhite()
			rightChar := n.parseChar()
			n.addCharRangeEdge(startState, endState, leftChar, rightChar)
			n.nextCharWithWhite()
		} else {
			n.addCharSingleEdge(startState, endState, leftChar)
			if n.curChar != ']' {
				n.addCharSingleEdge(startState, endState, n.curChar)
				n.nextCharWithWhite()
			}
		}
	}
	Assert(n.curChar == ']', n.errorWithLocalChar(InvalidComposed))
	n.nextChar()
	return
}

// string ::= "'" Char+ "'"

func (n *NFA) parseString() (startState, endState int) {
	Assert(n.curChar == '\'', n.errorWithLocalChar(InvalidStringQuote))
	startState = n.nextEnableState()
	curState := startState
	n.nextCharWithWhite()
	Assert(n.curChar != '\'', n.errorWithLocalChar(InvalidStringChar))
	for n.curChar != '\'' {
		n.curChar = n.parseChar()
		nextState := n.nextEnableState()
		n.addCharSingleEdge(curState, nextState, n.curChar)
		curState = nextState
		n.nextCharWithWhite()
	}
	// '\''
	endState = curState
	n.nextChar()
	return
}

// Char   ::=  SAFE_CHAR | ESC_CHAR
// SAFE_CHAR ::= ~['\\]
// ESC_CHAR ::= '\' ['\\bfnrt]
func (n *NFA) parseChar() rune {
	if n.curChar == '\\' { //esc Char
		n.nextCharWithWhite()
		return n.escChar(n.curChar)
	}
	return n.curChar
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

func IsComposedStart(char rune) bool {
	return char == '['
}

func IsTermStart(char rune) bool {
	return char == '\'' || char == '(' || char == '['
}

// [A-Z]
func IsNameStart(char rune) bool {
	return char >= 'A' && char <= 'Z'
}

// [0-9a-zA-Z_]*
func IsNameContinue(char rune) bool {
	return char == '_' ||
		(char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z')
}

func IsWhite(char rune) bool {
	return char == ' ' || char == '\t' || char == '\v' || char == '\r' || char == '\n' || char == '\f'
}
