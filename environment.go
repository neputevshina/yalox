package main

import (
	"fmt"
	"sync"
)

type Environment struct {
	sync.Mutex
	enclosing *Environment
	values    map[string]interface{}
}

func (e *Environment) Define(name string, val interface{}) {
	e.values[name] = val
}

func (e *Environment) Get(name Token) (interface{}, *Error) {
	val, ok := e.values[string(name.Lexeme)]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}
		return nil, &Error{name, "undefined variable '" + string(name.Lexeme) + "'"}
	}
	return val, nil
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) Assign(name Token, value interface{}) *Error {
	lex := string(name.Lexeme)
	if _, k := e.values[lex]; k {
		e.values[lex] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return &Error{name, fmt.Sprintf("undefined variable '%s'", lex)}
}
