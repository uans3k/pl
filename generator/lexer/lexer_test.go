package lexer

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/uans3k/pl/generator/lexer/backend/golang"
	"github.com/uans3k/pl/generator/lexer/generator"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	Convey("test lexer", t, func() {
		gen := generator.NewGenerator(golang.NewGolangBackend())
		fmt.Printf("config :\n %+v", gen.ConfigDescription())
		in, err := os.Open("./doc/lex/t1.lex")
		So(err, ShouldBeNil)
		out, err := os.Create("./out/lexer.go")
		So(err, ShouldBeNil)
		err = gen.Generate(in, out, map[string]string{
			golang.PackageConfig.Key:    "out",
			golang.RuntimeKeyConfig.Key: "github.com/uans3k/pl/runtime/lexer",
		})
		So(err, ShouldBeNil)
	})
}
