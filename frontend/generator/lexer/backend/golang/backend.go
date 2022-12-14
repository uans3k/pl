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
	TokenTypes            map[string]bool
	AcceptState2TokenType map[int]*fa.TokenType
	TokenType2FuncCalls   map[string][]string
	State2Edges           [][]*fa.MinDFAEdge
	Package               string
	Runtime               string
}

type golangBackend struct {
}

func (g *golangBackend) ConfigDescription() []*generator.ConfigDescription {
	return []*generator.ConfigDescription{
		PackageConfig, RuntimeKeyConfig,
	}
}

func (g *golangBackend) Generate(minDFA *fa.MinDFA, writer io.Writer, config map[string]string) error {
	packagePath, ok := config[PackageConfig.Key]
	if !ok {
		return errors.Wrapf(generator.InvalidConfig, "%s must fill", PackageConfig.Key)
	}
	runtimePath, ok := config[RuntimeKeyConfig.Key]
	if !ok {
		return errors.Wrapf(generator.InvalidConfig, "%s must fill", RuntimeKeyConfig.Key)
	}

	tmpl, err := template.New("goLexerTemplate").Funcs(template.FuncMap{
		"CharCompare": func(char fa.Char) string {
			switch v := char.(type) {
			case *fa.CharSingle:
				return fmt.Sprintf("char == %d", v.Char)
			case *fa.CharRange:
				return fmt.Sprintf("char >= %d && char <= %d", v.LeftChar, v.RightChar)
			default:
				panic("unexpect char type")
			}
		},
	}).Parse(goLexerTemplate)
	if err != nil {
		return err
	}
	tokenType2FuncCalls := map[string][]string{}
	for _, tokenType := range minDFA.AcceptState2TokenType {

		tokenType2FuncCalls[tokenType.Name] = tokenType.FuncCalls
	}
	return tmpl.Execute(writer, &lexerData{
		TokenTypes:            minDFA.TokenTypes,
		AcceptState2TokenType: minDFA.AcceptState2TokenType,
		TokenType2FuncCalls:   tokenType2FuncCalls,
		State2Edges:           minDFA.State2Edges,
		Package:               packagePath,
		Runtime:               runtimePath,
	})
}

func NewGolangBackend() generator.Backend {
	return &golangBackend{}
}
