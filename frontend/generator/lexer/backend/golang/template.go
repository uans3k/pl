package golang

var goLexerTemplate = `// Code generated by u3k_lexer. DO NOT EDIT.
package {{.Package}}

import (
	"github.com/pkg/errors"
	"io"
	runtime "{{.Runtime}}"
)

var(
	{{ range $i,$tokenType := .TokenTypes }}
	TokenType_{{$tokenType.Name}} runtime.TokenType = {{$i}}
	{{ end }}
	acceptState2TokenType = map[int]runtime.TokenType{
		{{ range $state,$tokenType := .AcceptState2TokenType }}
			{{$state}} : TokenType_{{$tokenType.Name}},
		{{ end }}
	}
	tokenType2FuncCalls = [][]string{
		{{ range $i,$tokenType := .TokenTypes }}[]string{ {{ range $j,$funcCall := $tokenType.FuncCalls }} {{if ne $j 0}},{{end}}"{{$funcCall}}"{{ end }} },
		{{ end }}
	}
)

func accept(state int) (ok bool){
	_,ok = acceptState2TokenType[state]
	return
}

type lexer struct{
	stream runtime.CharStream
	row  int
	col  int
	latestToken []rune
}

func NewLexer(stream runtime.CharStream) runtime.Lexer{
	return &lexer{
		stream : stream,
	}
}

func (l *lexer) handleFuncCalls(funcCalls []string) (skip bool){
	for _,funcCall:=range funcCalls{
		switch funcCall{
		case "Hidden":
			skip = true
		case "Row":
			l.row++
			l.col=0
		}
	}
	return 
}

func (l *lexer) NextToken() (runtime.Token,error){
sInit:
	var(
    	badState 	= -1
		states   	= []int{badState}
		curState 	= 0
		char rune 	= 0
		err  error 	= nil
	)
	goto s0
{{ range $state,$edges := .State2Edges }}
s{{$state}}:
	if accept({{$state}}){
		states = nil
	}
    states = append(states,{{$state}})
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	{{if eq (len $edges) 0 }}goto sEnd{{ else }}
	{{ range $i,$edge := $edges }}{{ if eq $i 0 }}if{{ else }}else if{{ end }}{{ range $j,$char := $edge.Chars }}{{ if ne $j 0 }}||{{ end }} {{$.CharCompare $char}} {{ end }}{
		goto s{{$edge.ToState}}
	}{{ end }}else{
		goto sEnd
	}{{ end }}
{{ end }}

sEnd:
	if err!=nil && err != io.EOF{
		return nil,err
	}
	for curState!=badState && !accept(curState) {
		curState = states[len(states)-1]
		states = states[0:len(states)-1]
		l.stream.Rollback()
	}
	if acceptTokenType,ok:= acceptState2TokenType[curState];ok{
		tokenStr:=l.stream.Consume()
		if funcCalls:=  tokenType2FuncCalls[acceptTokenType];len(funcCalls)!=0{
			if l.handleFuncCalls(funcCalls){
				goto sInit
			}
		}
		col := l.col
		l.col++
		l.latestToken = tokenStr
		return runtime.NewToken(tokenStr,acceptTokenType,l.row,col),nil
	}else if err==io.EOF{
		return runtime.NewToken([]rune("EOF"),runtime.TokenType_EOF,l.row,l.col),nil
	}else{
		return nil,errors.Wrapf(runtime.InvalidToken,"[row : %d ,col :%d ,latest token :%s]",l.row,l.col,string(l.latestToken))
	}
}
`

var tokenTemplate = `{{ range $i,$tokenType := .TokenTypes }}{{$tokenType.Name}}={{$i}}
{{ end }}
`
