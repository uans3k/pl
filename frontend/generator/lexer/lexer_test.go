package lexer

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/uans3k/pl/frontend/generator/lexer/backend/golang"
	"github.com/uans3k/pl/frontend/generator/lexer/generator"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	Convey("test lexer", t, func() {
		gen := generator.NewGenerator(golang.NewGolangBackend())
		fmt.Printf("config :\n %+v", gen.ConfigDescription())
		in, err := os.Open(".output/doc/lex/t1.lex")
		So(err, ShouldBeNil)
		out, err := os.Create(".output/gen/lexer.go")
		So(err, ShouldBeNil)
		tokenOut, err := os.Create(".output/gen/lexer.token")
		So(err, ShouldBeNil)
		err = gen.Generate(in, out, tokenOut, map[string]string{
			golang.PackageConfig.Key:    "gen",
			golang.RuntimeKeyConfig.Key: "github.com/uans3k/pl/frontend/runtime/lexer",
		})
		So(err, ShouldBeNil)
	})
}
