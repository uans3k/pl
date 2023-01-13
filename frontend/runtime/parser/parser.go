package parser

type Type int

const (
	Terminal   Type = 0
	NoTerminal Type = 1
	Epsilon    Type = 2
	EOF        Type = 3
)

type NodeType interface {
	String() int
	Type() Type
}

type NodeValue interface {
	Chars() []rune
	String() string
	Row() int
	Column() int
}

type Node struct {
	Type     NodeType
	Value    NodeValue
	Children []*Node
}

type Parser interface {
	Parser() (*Node, error)
}
