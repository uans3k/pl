package generator

import (
	"github.com/pkg/errors"
	"github.com/uans3k/pl/frontend/generator/lexer/fa"
	"io"
)

var (
	InvalidConfig = errors.New("invalid config")
)

type ConfigDescription struct {
	Key         string
	Description string
}

type Backend interface {
	ConfigDescription() []*ConfigDescription
	Generate(minDFA *fa.MinDFA, writer io.Writer, config map[string]string) error
}

type Generator struct {
	Backend
}

func NewGenerator(backend Backend) *Generator {
	return &Generator{Backend: backend}
}

func (g *Generator) Generate(in io.Reader, out io.Writer, config map[string]string) error {
	nfa, err := fa.ParseNFA(in)
	if err != nil {
		return err
	}
	return g.Backend.Generate(fa.DFA2MinDFA(fa.NFA2DFA(nfa)), out, config)
}
