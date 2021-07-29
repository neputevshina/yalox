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
	statements := make([]Stmt, 0, 10)
	done := make(chan struct{})
	go func() {
		for !p.isAtEnd() {
			decl := p.declaration()
			statements = append(statements, decl)
		}
		done <- struct{}{}
	}()
	select {
	case err := <-p.raise:
		return nil, err
	case <-done:
		return statements, nil
	}
}

func (p *Parser) declaration() Stmt {
	done := make(chan struct{})
	var stmt Stmt
	go func() {
		if p.match(tokenVar) {
			stmt = p.varDeclaration()
		} else {
			stmt = p.statement()
		}
		done <- struct{}{}
	}()
	select {
	case err := <-p.raise:
		p.synchronize()
		// Please kill me...
		p.raise <- err
		return nil
	case <-done:
		return stmt
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

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(tokenIdent, "expect variable name")
	var init Expr
	if p.match(tokenEqual) {
		init = p.expression()
	}
	p.consume(tokenSemicolon, "expect ';' after variable declaration")
	return &Var{name, init}
}

func (p *Parser) statement() Stmt {
	switch {
	case p.match(tokenFor):
		return p.forStatement()
	case p.match(tokenIf):
		return p.ifStatement()
	case p.match(tokenPrint):
		return p.printStatement()
	case p.match(tokenWhile):
		return p.whileStatement()
	case p.match(tokenLeftBrace):
		return &Block{p.block()}
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) forStatement() Stmt {
	p.consume(tokenLeftParen, "expect '(' after 'for'")
	var init Stmt
	if p.match(tokenSemicolon) {
		init = nil
	} else if p.match(tokenVar) {
		init = p.varDeclaration()
	} else {
		init = p.expressionStatement()
	}
	var cond Expr
	if !p.check(tokenSemicolon) {
		cond = p.expression()
	}
	p.consume(tokenSemicolon, "expect ';' after loop condition")
	var incr Expr
	if !p.check(tokenRightParen) {
		incr = p.expression()
	}
	p.consume(tokenRightParen, "expect ')' after loop condition")
	body := p.statement()
	if incr != nil {
		body = &Block{Stmts: []Stmt{body, &Expression{incr}}}
	}
	if cond == nil {
		cond = &Literal{true}
	}
	body = &While{cond, body}
	if init != nil {
		body = &Block{[]Stmt{init, body}}
	}

	return body
}

func (p *Parser) whileStatement() Stmt {
	p.consume(tokenLeftParen, "expect '(' after 'while'")
	cond := p.expression()
	p.consume(tokenRightParen, "expect ')' after condition")
	body := p.statement()
	return &While{cond, body}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(tokenLeftParen, "expect '(' after 'if'")
	cond := p.expression()
	p.consume(tokenRightParen, "expect ')' after condition")

	then := p.statement()
	els := Stmt(nil)
	if p.match(tokenElse) {
		els = p.statement()
	}
	return &If{cond, then, els}
}

func (p *Parser) block() []Stmt {
	stmts := make([]Stmt, 0, 10)
	for !(p.check(tokenRightBrace) || p.isAtEnd()) {
		stmts = append(stmts, p.declaration())
	}
	p.consume(tokenRightBrace, "expect '}' after block")
	return stmts
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
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(tokenEqual) {
		equals := p.previous()
		value := p.assignment()
		if e, k := expr.(*Variable); k {
			name := e.Name
			return &Assign{name, value}
		}
		loxerr2(&Error{equals, "invalid assignment target"})
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(tokenOr) {
		op := p.previous()
		r := p.and()
		expr = &Logical{expr, op, r}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(tokenAnd) {
		op := p.previous()
		r := p.and()
		expr = &Logical{expr, op, r}
	}
	return expr
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
