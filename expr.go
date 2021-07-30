package main

// Expr is an expression of the parser.
type Expr interface {
	Accept(Visitor) (interface{}, *Error)
}

type Stmt interface {
	Accept(Visitor) (interface{}, *Error)
}

type Visitor interface {
	Visit(interface{}) (interface{}, *Error)
}

type Callable interface {
	Call(i *Interpreter, args []interface{}) (interface{}, *Error)
	Arity() int
}
