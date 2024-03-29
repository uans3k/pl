// Code generated by u3k_lexer. DO NOT EDIT.
package gen

import (
	"github.com/pkg/errors"
	"io"
	runtime "github.com/uans3k/pl/frontend/runtime/lexer"
)

var(
	
	TokenType_White runtime.TokenType = 0
	
	TokenType_Line runtime.TokenType = 1
	
	TokenType_Left_b runtime.TokenType = 2
	
	TokenType_Right_b runtime.TokenType = 3
	
	TokenType_String runtime.TokenType = 4
	
	TokenType_Number runtime.TokenType = 5
	
	acceptState2TokenType = map[int]runtime.TokenType{
		
			1 : TokenType_String,
		
			2 : TokenType_Left_b,
		
			3 : TokenType_Number,
		
			4 : TokenType_Number,
		
			5 : TokenType_White,
		
			6 : TokenType_Line,
		
			7 : TokenType_Right_b,
		
	}
	tokenType2FuncCalls = [][]string{
		[]string{  "Hidden" },
		[]string{  "Hidden" ,"Row" },
		[]string{  },
		[]string{  },
		[]string{  },
		[]string{  },
		
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

s0:
	if accept(0){
		states = nil
	}
    states = append(states,0)
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	
	if char >= 97 && char <= 122 || char >= 65 && char <= 90 || char == 95 {
		goto s1
	}else if char == 123 {
		goto s2
	}else if char >= 49 && char <= 57 {
		goto s3
	}else if char == 32 || char == 9 {
		goto s5
	}else if char == 48 {
		goto s4
	}else if char == 10 {
		goto s6
	}else if char == 125 {
		goto s7
	}else{
		goto sEnd
	}

s1:
	if accept(1){
		states = nil
	}
    states = append(states,1)
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	
	if char == 95 || char >= 97 && char <= 122 || char >= 65 && char <= 90 || char >= 48 && char <= 57 {
		goto s1
	}else{
		goto sEnd
	}

s2:
	if accept(2){
		states = nil
	}
    states = append(states,2)
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	goto sEnd

s3:
	if accept(3){
		states = nil
	}
    states = append(states,3)
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	
	if char >= 48 && char <= 57 {
		goto s3
	}else{
		goto sEnd
	}

s4:
	if accept(4){
		states = nil
	}
    states = append(states,4)
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	goto sEnd

s5:
	if accept(5){
		states = nil
	}
    states = append(states,5)
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	goto sEnd

s6:
	if accept(6){
		states = nil
	}
    states = append(states,6)
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	goto sEnd

s7:
	if accept(7){
		states = nil
	}
    states = append(states,7)
	char,err = l.stream.NextChar()
	if err!=nil{
		goto sEnd
	}
	goto sEnd


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
