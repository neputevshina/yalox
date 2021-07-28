package main

type Parser struct {
	Tokens  []Token
	current int
	raise   chan *Error
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens: tokens,
		raise:  make(chan *Error),
	}
}

func (p *Parser) Parse() ([]Stmt, *Error) {
	// Moved it, probably prematurely...
	statements := make([]Stmt, 0, 10)
	for !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}
	return statements, nil
}

func (p *Parser) declaration() (Stmt, *Error) {
	// Please kill me...
	done := make(chan struct{})
	var stmt Stmt
	go func() {
		if p.match(tokenVar) {
			stmt = p.vardeclaration()
		} else {
			stmt = p.statement()
		}
		done <- struct{}{}
	}()
	select {
	case err := <-p.raise:
		p.synchronize()
		return nil, err
	case <-done:
		return stmt, nil
	}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == tokenSemicolon {
			return
		}
	}
	switch p.peek().Type {
	case tokenClass:
		fallthrough
	case tokenFun:
		fallthrough
	case tokenVar:
		fallthrough
	case tokenFor:
		fallthrough
	case tokenIf:
		fallthrough
	case tokenWhile:
		fallthrough
	case tokenPrint:
		fallthrough
	case tokenReturn:
		return
	}
	p.advance()
}

func (p *Parser) vardeclaration() Stmt {
	name := p.consume(tokenIdent, "expect variable name")
	var init Expr
	if p.match(tokenEqual) {
		init = p.expression()
	}
	p.consume(tokenSemicolon, "expect ';' after variable declaration")
	return &Var{name, init}
}

func (p *Parser) statement() Stmt {
	if p.match(tokenPrint) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(tokenSemicolon, "expect ';' after value")
	return &Print{Expr: value}
}

func (p *Parser) expressionStatement() Stmt {
	value := p.expression()
	p.consume(tokenSemicolon, "expect ';' after expression")
	return &Expression{Expr: value}
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	e := p.comparison()
	for p.match(tokenBangEqual, tokenEqualEqual) {
		op := p.previous()
		r := p.comparison()
		e = &Binary{e, op, r}
	}
	return e
}

func (p *Parser) comparison() Expr {
	e := p.term()
	for p.match(tokenGreater, tokenGreaterEqual, tokenLess, tokenLessEqual) {
		op := p.previous()
		r := p.term()
		e = &Binary{e, op, r}
	}
	return e
}

func (p *Parser) term() Expr {
	e := p.factor()
	for p.match(tokenMinus, tokenPlus) {
		op := p.previous()
		r := p.factor()
		e = &Binary{e, op, r}
	}
	return e
}

func (p *Parser) factor() Expr {
	e := p.unary()
	for p.match(tokenSlash, tokenStar) {
		op := p.previous()
		r := p.unary()
		e = &Binary{e, op, r}
	}
	return e
}

func (p *Parser) unary() Expr {
	if p.match(tokenBang, tokenMinus) {
		op := p.previous()
		r := p.unary()
		return &Unary{op, r}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(tokenFalse):
		return &Literal{false}
	case p.match(tokenTrue):
		return &Literal{true}
	case p.match(tokenNil):
		return &Literal{nil}
	}

	if p.match(tokenNumber, tokenString) {
		return &Literal{p.previous().Literal}
	}
	if p.match(tokenIdent) {
		return &Variable{p.previous()}
	}
	if p.match(tokenLeftParen) {
		e := p.expression()
		p.consume(tokenRightParen, "expect ')' after expression")
		return &Grouping{e}
	}
	panic("debug")
	p.raise <- &Error{p.peek(), "expect expression or value"}
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
	p.raise <- &Error{p.peek(), message}
	return Token{}
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
