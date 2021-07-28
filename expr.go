package main

// Expr is an expression of the parser.
type Expr interface {
	Accept(Visitor) interface{}
}

type Stmt interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	Visit(interface{}) interface{}
}
