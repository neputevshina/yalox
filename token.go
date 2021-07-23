package main

import "fmt"

// Token is a lexical token of an interpreter.
type Token struct {
	Type    int
	Lexeme  []byte
	Literal interface{}
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", t.Type, t.Lexeme, t.Literal)
}
