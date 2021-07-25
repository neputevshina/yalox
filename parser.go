package main

type Parser struct {
	Tokens  []Token
	current int
	raise   chan struct{}
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens: tokens,
		raise:  make(chan struct{}),
	}
}

func (p *Parser) Parse() Expr {
	// Somewhat retarded way to do it, but it would probably work.
	var expr Expr
	done := make(chan struct{})
	go func() {
		expr = p.expression()
		done <- struct{}{}
	}()
	select {
	case <-p.raise:
		return NewLiteral("<error>")
	case <-done:
		return expr
	}
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	e := p.comparison()
	for p.match(tokenBangEqual, tokenEqualEqual) {
		op := p.previous()
		r := p.comparison()
		e = NewBinary(e, op, r)
	}
	return e
}

func (p *Parser) comparison() Expr {
	e := p.term()
	for p.match(tokenGreater, tokenGreaterEqual, tokenLess, tokenLessEqual) {
		op := p.previous()
		r := p.term()
		e = NewBinary(e, op, r)
	}
	return e
}

func (p *Parser) term() Expr {
	e := p.factor()
	for p.match(tokenMinus, tokenPlus) {
		op := p.previous()
		r := p.factor()
		e = NewBinary(e, op, r)
	}
	return e
}

func (p *Parser) factor() Expr {
	e := p.unary()
	for p.match(tokenSlash, tokenStar) {
		op := p.previous()
		r := p.unary()
		e = NewBinary(e, op, r)
	}
	return e
}

func (p *Parser) unary() Expr {
	if p.match(tokenBang, tokenMinus) {
		op := p.previous()
		r := p.unary()
		return NewUnary(op, r)
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(tokenFalse):
		return NewLiteral(false)
	case p.match(tokenTrue):
		return NewLiteral(true)
	case p.match(tokenNil):
		return NewLiteral(nil)
	}

	if p.match(tokenNumber, tokenString) {
		return NewLiteral(p.previous().Literal)
	}
	if p.match(tokenLeftParen) {
		e := p.expression()
		p.consume(tokenRightParen, "expect ')' after expression")
		return NewGrouping(e)
	}
	p.raise <- struct{}{}
	return nil
}

func (p *Parser) match(types ...int) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(typ int, message string) Token {
	if p.check(typ) {
		return p.advance()
	}
	loxparseerr(p.peek(), message)
	p.raise <- struct{}{}
	return Token{} // todo
}

func (p *Parser) check(typ int) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == typ
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tokenEOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.current-1]
}
