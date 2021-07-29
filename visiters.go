// Generated, DO NOT EDIT. Name is intentional.
package main

func (l *Logical) Accept(vis Visitor) interface{} {
	return vis.Visit(l)
}

func (l *Literal) Accept(vis Visitor) interface{} {
	return vis.Visit(l)
}

func (v *Variable) Accept(vis Visitor) interface{} {
	return vis.Visit(v)
}

func (p *Print) Accept(vis Visitor) interface{} {
	return vis.Visit(p)
}

func (v *Var) Accept(vis Visitor) interface{} {
	return vis.Visit(v)
}

func (a *Assign) Accept(vis Visitor) interface{} {
	return vis.Visit(a)
}

func (b *Block) Accept(vis Visitor) interface{} {
	return vis.Visit(b)
}

func (i *If) Accept(vis Visitor) interface{} {
	return vis.Visit(i)
}

func (w *While) Accept(vis Visitor) interface{} {
	return vis.Visit(w)
}

func (b *Binary) Accept(vis Visitor) interface{} {
	return vis.Visit(b)
}

func (g *Grouping) Accept(vis Visitor) interface{} {
	return vis.Visit(g)
}

func (u *Unary) Accept(vis Visitor) interface{} {
	return vis.Visit(u)
}

func (e *Expression) Accept(vis Visitor) interface{} {
	return vis.Visit(e)
}
