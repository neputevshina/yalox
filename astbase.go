package main

type Binary struct {
	Left  Expr
	Op    Token
	Right Expr
}

func (b *Binary) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(b)
}

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(g)
}

type Literal struct {
	Val interface{}
}

func (l *Literal) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(l)
}

type Unary struct {
	Op    Token
	Right Expr
}

func (u *Unary) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(u)
}
