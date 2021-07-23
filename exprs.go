package main

func NewBinary(left Expr, op *Token, right Expr) Expr {
	return &Binary{
		left:  left,
		op:    op,
		right: right,
	}
}

type Binary struct {
	left  Expr
	op    *Token
	right Expr
}

func (b *Binary) Left() Expr  { return nil }
func (b *Binary) Op() *Token  { return &Token{} }
func (b *Binary) Right() Expr { return nil }

func (b *Binary) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(b)
}

func NewGrouping(expr Expr) Expr {
	return &Grouping{
		expr: expr,
	}
}

type Grouping struct {
	expr Expr
}

func (g *Grouping) Left() Expr  { return nil }
func (g *Grouping) Op() *Token  { return &Token{} }
func (g *Grouping) Right() Expr { return nil }

func (g *Grouping) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(g)
}

func NewLiteral(val interface{}) Expr {
	return &Literal{
		val: val,
	}
}

type Literal struct {
	val interface{}
}

func (l *Literal) Left() Expr  { return nil }
func (l *Literal) Op() *Token  { return &Token{} }
func (l *Literal) Right() Expr { return nil }

func (l *Literal) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(l)
}

func NewUnary(op *Token, right Expr) Expr {
	return &Unary{
		op:    op,
		right: right,
	}
}

type Unary struct {
	op    *Token
	right Expr
}

func (u *Unary) Left() Expr  { return nil }
func (u *Unary) Op() *Token  { return &Token{} }
func (u *Unary) Right() Expr { return nil }

func (u *Unary) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(u)
}
