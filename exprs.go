package main

// Binary
func Binary(left Expr, op *Token, right Expr) Expr {
	return &binary{
		left:  left,
		op:    op,
		right: right,
	}
}

type binary struct {
	left  Expr
	op    *Token
	right Expr
}

func (b *binary) Left() Expr  { return nil }
func (b *binary) Op() *Token  { return &Token{} }
func (b *binary) Right() Expr { return nil }

func (b *binary) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(b)
}

// Grouping
func Grouping(expr Expr) Expr {
	return &grouping{
		expr: expr,
	}
}

type grouping struct {
	expr Expr
}

func (g *grouping) Left() Expr  { return nil }
func (g *grouping) Op() *Token  { return &Token{} }
func (g *grouping) Right() Expr { return nil }

func (g *grouping) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(g)
}

// Literal
func Literal(val interface{}) Expr {
	return &literal{
		val: val,
	}
}

type literal struct {
	val interface{}
}

func (l *literal) Left() Expr  { return nil }
func (l *literal) Op() *Token  { return &Token{} }
func (l *literal) Right() Expr { return nil }

func (l *literal) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(l)
}

// Unary
func Unary(op *Token, right Expr) Expr {
	return &unary{
		op:    op,
		right: right,
	}
}

type unary struct {
	op    *Token
	right Expr
}

func (u *unary) Left() Expr  { return nil }
func (u *unary) Op() *Token  { return &Token{} }
func (u *unary) Right() Expr { return nil }

func (u *unary) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(u)
}
