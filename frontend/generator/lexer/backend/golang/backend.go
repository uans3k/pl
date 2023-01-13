package golang

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/uans3k/pl/frontend/generator/lexer/fa"
	"github.com/uans3k/pl/frontend/generator/lexer/generator"
	"io"
	"text/template"
)

var (
	PackageConfig = &generator.ConfigDescription{
		Key:         "package",
		Description: "generate file package",
	}
	RuntimeKeyConfig = &generator.ConfigDescription{
		Key:         "runtime",
		Description: "generate file import runtime",
	}
)

type lexerData struct {
	TokenTypes            []*fa.TokenType
	AcceptState2TokenType map[int]*fa.TokenType
	State2Edges           [][]*fa.MinDFAEdge
	Package               string
	Runtime               string
}

func (l *lexerData) CharCompare(char fa.Char) string {
	switch v := char.(type) {
	case *fa.CharSingle:
		return fmt.Sprintf("char == %d", v.Char)
	case *fa.CharRange:
		return fmt.Sprintf("char >= %d && char <= %d", v.LeftChar, v.RightChar)
	default:
		panic("unexpect char type")
	}
}

type golangBackend struct {
}

func (g *golangBackend) ConfigDescription() []*generator.ConfigDescription {
	return []*generator.ConfigDescription{
		PackageConfig, RuntimeKeyConfig,
	}
}

func (g *golangBackend) Generate(minDFA *fa.MinDFA, lexerWriter io.Writer, lexerTokenWriter io.Writer, config map[string]string) error {
	packagePath, ok := config[PackageConfig.Key]
	if !ok {
		return errors.Wrapf(generator.InvalidConfig, "%s must fill", PackageConfig.Key)
	}
	runtimePath, ok := config[RuntimeKeyConfig.Key]
	if !ok {
		return errors.Wrapf(generator.InvalidConfig, "%s must fill", RuntimeKeyConfig.Key)
	}
	data := &lexerData{
		TokenTypes:            minDFA.TokenTypes,
		AcceptState2TokenType: minDFA.AcceptState2TokenType,
		State2Edges:           minDFA.State2Edges,
		Package:               packagePath,
		Runtime:               runtimePath,
	}
	tmpl, err := template.New("goLexerTemplate").Parse(goLexerTemplate)
	if err != nil {
		return err
	}
	err = tmpl.Execute(lexerWriter, data)
	if err != nil {
		return err
	}
	tokenTmpl, err := tmpl.New("tokenTemplate").Parse(tokenTemplate)
	if err != nil {
		return err
	}
	return tokenTmpl.Execute(lexerTokenWriter, data)
}

func NewGolangBackend() generator.Backend {
	return &golangBackend{}
}
