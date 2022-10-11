package fa

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestNFA(t *testing.T) {
	Convey("test nfa.go", t, func() {
		f, err := os.Open("../doc/t0.lex")
		So(err, ShouldBeNil)
		nfa, err := ParseNFA(f)
		So(err, ShouldBeNil)
		fmt.Printf("%+v", nfa)
	})
}
