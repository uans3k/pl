package parser

type Quantifier int

const (
	Quantifier_Absence Quantifier = 0
	Quantifier_Star    Quantifier = 1
	Quantifier_Plus    Quantifier = 2
	Quantifier_Quesion Quantifier = 3
)

type Parser struct {
	StartRule  StartRule
	Rules      []Rule
	State2Rule []Rule
}

type StartRule interface {
	RuleRef() RuleRef
	NT() NT
	EOF() EOF
}

type Rule interface {
	RuleRef() RuleRef
	Exprs() []Expr
}

type RuleRef interface {
	NT() NT
}

type Expr interface {
	Terms() []Term
}

type Term interface {
	Pieces() (pieces []Piece, ok bool)
	Epsilon() (e Epsilon, ok bool)
}

type Epsilon interface {
	Text() string
}

type Piece interface {
	Quantifier() Quantifier
	Factor() Factor
}

type Factor interface {
	Symbol() (expr Expr, ok bool)
	Expr() (expr Expr, ok bool)
}

type Symbol interface {
	T() (t T, ok bool)
	NT() (nt NT, ok bool)
}

type T interface {
	Text() string
	Type()
}

type NT interface {
	Text() string
}

type EOF interface {
	Text() string
}

type factor struct {
}

func (f factor) Expr() (expr Expr, ok bool) {
	return nil, false
}

func (f factor) NT() (nt NT, ok bool) {
	return nil, false
}

func (f factor) T() (t T, ok bool) {
	return nil, false
}

type parser struct {
	startRule StartRule
	rules     []Rule
}

func (p *parser) StartRule() StartRule {
	return p.startRule
}

func (p *parser) Rules() []Rule {
	return p.rules
}

type rule struct {
}
