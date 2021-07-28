package main

import "fmt"

type Environment struct {
	values map[string]interface{}
}

func (e *Environment) Define(name string, val interface{}) {
	e.values[name] = val
}

func (e *Environment) Get(name Token) (interface{}, *Error) {
	fmt.Println(e.values)
	val, ok := e.values[string(name.Lexeme)]
	if !ok {
		return nil, &Error{name, "undefined variable '" + string(name.Lexeme) + "'"}
	}
	return val, nil
}
