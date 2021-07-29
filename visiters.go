// Generated, DO NOT EDIT. Name is intentional.
package main

func (p *Print) Accept(vis Visitor) interface{} {
	return vis.Visit(p)
}

func (g *Grouping) Accept(vis Visitor) interface{} {
	return vis.Visit(g)
}

func (l *Literal) Accept(vis Visitor) interface{} {
	return vis.Visit(l)
}

func (v *Variable) Accept(vis Visitor) interface{} {
	return vis.Visit(v)
}

func (e *Expression) Accept(vis Visitor) interface{} {
	return vis.Visit(e)
}

func (b *Block) Accept(vis Visitor) interface{} {
	return vis.Visit(b)
}

func (b *Binary) Accept(vis Visitor) interface{} {
	return vis.Visit(b)
}

func (u *Unary) Accept(vis Visitor) interface{} {
	return vis.Visit(u)
}

func (v *Var) Accept(vis Visitor) interface{} {
	return vis.Visit(v)
}

func (a *Assign) Accept(vis Visitor) interface{} {
	return vis.Visit(a)
}
