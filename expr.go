package main

// Expr is an expression of the parser.
type Expr interface {
	Left() Expr
	Right() Expr
	Op() *Token
	Accept(ExprVisitor) interface{}
}

type ExprVisitor interface {
	Visit(interface{}) interface{}
}
