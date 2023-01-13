package golang

import (
	"github.com/uans3k/pl/frontend/generator/paser/generator"
	"github.com/uans3k/pl/frontend/generator/paser/ll"
	"io"
)

type golangBackend struct {
}

func (g golangBackend) ConfigDescription() []*generator.ConfigDescription {
	//TODO implement me
	panic("implement me")
}

func (g golangBackend) Generate(firstPlusSet *ll.FirstPlusSet, out io.Writer, config map[string]string) error {
	//TODO implement me
	panic("implement me")
}

func NewGolangBackend() ll.Backend {
	return &golangBackend{}
}
