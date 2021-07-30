// Generated, DO NOT EDIT. Name is intentional.
package main

// Accept is an auto-generated acceptor method for While
func (w *While) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(w)
}

// Accept is an auto-generated acceptor method for Binary
func (b *Binary) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(b)
}

// Accept is an auto-generated acceptor method for Grouping
func (g *Grouping) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(g)
}

// Accept is an auto-generated acceptor method for Variable
func (v *Variable) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(v)
}

// Accept is an auto-generated acceptor method for If
func (i *If) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(i)
}

// Accept is an auto-generated acceptor method for Assign
func (a *Assign) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(a)
}

// Accept is an auto-generated acceptor method for Return
func (r *Return) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(r)
}

// Accept is an auto-generated acceptor method for Expression
func (e *Expression) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(e)
}

// Accept is an auto-generated acceptor method for Print
func (p *Print) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(p)
}

// Accept is an auto-generated acceptor method for Block
func (b *Block) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(b)
}

// Accept is an auto-generated acceptor method for Call
func (c *Call) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(c)
}

// Accept is an auto-generated acceptor method for Function
func (f *Function) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(f)
}

// Accept is an auto-generated acceptor method for Literal
func (l *Literal) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(l)
}

// Accept is an auto-generated acceptor method for Unary
func (u *Unary) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(u)
}

// Accept is an auto-generated acceptor method for Var
func (v *Var) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(v)
}

// Accept is an auto-generated acceptor method for Logical
func (l *Logical) Accept(vis Visitor) (interface{}, *Error) {
	return vis.Visit(l)
}
