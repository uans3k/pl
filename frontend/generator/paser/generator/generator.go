package generator

import (
	"github.com/pkg/errors"
	"io"
)

var (
	InvalidConfig = errors.New("invalid config")
)

type ConfigDescription struct {
	Key         string
	Description string
}

type Generator interface {
	ConfigDescription() []*ConfigDescription
	Generate(in io.Reader, out io.Writer, config map[string]string) error
}
