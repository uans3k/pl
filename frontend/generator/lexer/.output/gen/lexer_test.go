package gen

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	runtime "github.com/uans3k/pl/frontend/runtime/lexer"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	Convey("test lexer.go", t, func() {
		in, err := os.Open("../doc/file/t1_file.u")
		So(err, ShouldBeNil)
		charStream := runtime.NewCycleCharStream(1024, in)
		aLexer := NewLexer(charStream)
		var tokens []string
		for {
			token, err := aLexer.NextToken()
			So(err, ShouldBeNil)
			tokens = append(tokens, token.Text())
			if token.Type() == runtime.TokenType_EOF {
				break
			}
		}
		fmt.Printf("%+v \n", tokens)
		So(len(tokens), ShouldEqual, 8)
	})
}
