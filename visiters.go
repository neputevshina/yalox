package main

func (e *Expression) Accept(vis Visitor) interface{} {
	return vis.Visit(e)
}

func (p *Print) Accept(vis Visitor) interface{} {
	return vis.Visit(p)
}

func (v *Var) Accept(vis Visitor) interface{} {
	return vis.Visit(v)
}

func (b *Binary) Accept(vis Visitor) interface{} {
	return vis.Visit(b)
}

func (g *Grouping) Accept(vis Visitor) interface{} {
	return vis.Visit(g)
}

func (l *Literal) Accept(vis Visitor) interface{} {
	return vis.Visit(l)
}

func (u *Unary) Accept(vis Visitor) interface{} {
	return vis.Visit(u)
}

func (v *Variable) Accept(vis Visitor) interface{} {
	return vis.Visit(v)
}
