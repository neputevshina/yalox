package main

// Expr is an expression of the parser.
type Expr interface {
	Accept(ExprVisitor) interface{}
}

type ExprVisitor interface {
	Visit(interface{}) interface{}
}
