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
	var stmt Stmt
	var err *Error
	switch {
	case p.match(tokenFun):
		stmt, err = p.function("function")
	case p.match(tokenVar):
		stmt, err = p.varDeclaration()
	default:
		stmt, err = p.statement()
	}
	if err != nil {
		p.synchronize()
		return nil, err
	}
	return stmt, nil

}

func (p *Parser) function(kind string) (Stmt, *Error) {
	var body []Stmt
	params := make([]Token, 0, 10)

	name, err := p.consume(tokenIdent, "expect "+kind+" name")
	if err != nil {
		goto fail
	}
	if _, err = p.consume(tokenLeftParen, "expect '(' after "+kind+" name"); err != nil {
		goto fail
	}
	if !p.check(tokenRightParen) {
		var p2 Token
		p2, err = p.consume(tokenIdent, "expect parameter name")
		if err != nil {
			goto fail
		}
		params = append(params, p2)
		for p.match(tokenComma) {
			if len(params) >= 255 {
				loxerr2(&Error{p.peek(), "can't have more than 255 arguments"})
			}
			p2, err = p.consume(tokenIdent, "expect parameter name")
			if err != nil {
				goto fail
			}
			params = append(params, p2)
		}
	}

	if _, err = p.consume(tokenLeftBrace, "expect { before "+kind+" body"); err != nil {
		goto fail
	}
	if _, err = p.consume(tokenRightParen, "expect ')' after parameters"); err != nil {
		goto fail
	}
	body, err = p.block()
	if err != nil {
		return nil, err
	}
	return &Function{name, params, body}, nil
fail:
	return nil, err
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

func (p *Parser) varDeclaration() (Stmt, *Error) {
	name, err := p.consume(tokenIdent, "expect variable name")
	if err != nil {
		return nil, err
	}
	var init Expr
	if p.match(tokenEqual) {
		init, err = p.expression()
	}
	if _, err := p.consume(tokenSemicolon, "expect ';' after variable declaration"); err != nil {
		return nil, err
	}
	return &Var{name, init}, nil
}

func (p *Parser) statement() (Stmt, *Error) {
	switch {
	case p.match(tokenReturn):
		return p.returnStatement()
	case p.match(tokenFor):
		return p.forStatement()
	case p.match(tokenIf):
		return p.ifStatement()
	case p.match(tokenPrint):
		return p.printStatement()
	case p.match(tokenWhile):
		return p.whileStatement()
	case p.match(tokenLeftBrace):
		b, e := p.block()
		return &Block{b}, e
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) returnStatement() (s Stmt, err *Error) {
	kw := p.previous()
	val := Expr(nil)
	if !p.check(tokenSemicolon) {
		val, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokenSemicolon, "expect ';' after return value")
	return &Return{kw, val}, err
}

func (p *Parser) forStatement() (Stmt, *Error) {
	var err *Error
	_, err = p.consume(tokenLeftParen, "expect '(' after 'for'")
	var init Stmt
	if p.match(tokenSemicolon) {
		init = nil
	} else if p.match(tokenVar) {
		init, err = p.varDeclaration()
	} else {
		init, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}
	var cond Expr
	if !p.check(tokenSemicolon) {
		cond, err = p.expression()
	}
	_, err = p.consume(tokenSemicolon, "expect ';' after loop condition")
	if err != nil {
		return nil, err
	}
	var incr Expr
	if !p.check(tokenRightParen) {
		incr, err = p.expression()
	}
	_, err = p.consume(tokenRightParen, "expect ')' after loop condition")
	body, err := p.statement()
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
	return body, err
}

func (p *Parser) whileStatement() (Stmt, *Error) {
	if _, err := p.consume(tokenLeftParen, "expect '(' after 'while'"); err != nil {
		return nil, err
	}
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(tokenRightParen, "expect ')' after condition"); err != nil {
		return nil, err
	}
	body, err := p.statement()
	return &While{cond, body}, err
}

func (p *Parser) ifStatement() (Stmt, *Error) {
	if _, err := p.consume(tokenLeftParen, "expect '(' after 'if'"); err != nil {
		return nil, err
	}
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(tokenRightParen, "expect ')' after condition"); err != nil {
		return nil, err
	}
	then, err := p.statement()
	if err != nil {
		return nil, err
	}
	els := Stmt(nil)
	if p.match(tokenElse) {
		els, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &If{cond, then, els}, nil
}

func (p *Parser) block() ([]Stmt, *Error) {
	stmts := make([]Stmt, 0, 10)
	for !(p.check(tokenRightBrace) || p.isAtEnd()) {
		if s, err := p.declaration(); err != nil {
			return nil, err
		} else {
			stmts = append(stmts, s)
		}
	}
	if _, err := p.consume(tokenRightBrace, "expect '}' after block"); err != nil {
		return nil, err
	}
	return stmts, nil
}

func (p *Parser) printStatement() (Expr, *Error) {
	value, err := p.expression()
	_, err = p.consume(tokenSemicolon, "expect ';' after value")
	return &Print{Expr: value}, err
}

func (p *Parser) expressionStatement() (Expr, *Error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokenSemicolon, "expect ';' after expression")
	return &Expression{Expr: value}, err
}

func (p *Parser) expression() (Expr, *Error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, *Error) {
	expr, err := p.or()
	if p.match(tokenEqual) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if e, k := expr.(*Variable); k {
			name := e.Name
			return &Assign{name, value}, nil
		}
		loxerr2(&Error{equals, "invalid assignment target"})
	}
	return expr, err
}

func (p *Parser) or() (Expr, *Error) {
	expr, err := p.and()
	for p.match(tokenOr) {
		op := p.previous()
		var r Expr
		r, err = p.and()
		expr = &Logical{expr, op, r}
	}
	return expr, err
}

func (p *Parser) and() (Expr, *Error) {
	expr, err := p.equality()
	for p.match(tokenAnd) {
		op := p.previous()
		var r Expr
		r, err = p.and()
		expr = &Logical{expr, op, r}
	}
	return expr, err
}

func (p *Parser) equality() (Expr, *Error) {
	e, err := p.comparison()
	for p.match(tokenBangEqual, tokenEqualEqual) {
		op := p.previous()
		var r Expr
		r, err = p.comparison()
		e = &Binary{e, op, r}
	}
	return e, err
}

func (p *Parser) comparison() (Expr, *Error) {
	e, err := p.term()
	for p.match(tokenGreater, tokenGreaterEqual, tokenLess, tokenLessEqual) {
		op := p.previous()
		var r Expr
		r, err = p.term()
		e = &Binary{e, op, r}
	}
	return e, err
}

func (p *Parser) term() (Expr, *Error) {
	e, err := p.factor()
	for p.match(tokenMinus, tokenPlus) {
		op := p.previous()
		var r Expr
		r, err = p.factor()
		e = &Binary{e, op, r}
	}
	return e, err
}

func (p *Parser) factor() (Expr, *Error) {
	e, err := p.unary()
	for p.match(tokenSlash, tokenStar) {
		op := p.previous()
		var r Expr
		r, err = p.unary()
		e = &Binary{e, op, r}
	}
	return e, err
}

func (p *Parser) unary() (Expr, *Error) {
	if p.match(tokenBang, tokenMinus) {
		op := p.previous()
		r, err := p.unary()
		return &Unary{op, r}, err
	}
	return p.call()
}

func (p *Parser) call() (Expr, *Error) {
	e, err := p.primary()
	for {
		if p.match(tokenLeftParen) {
			e, err = p.finishCall(e)
		} else {
			break
		}
	}
	return e, err
}

func (p *Parser) finishCall(callee Expr) (Expr, *Error) {
	args := make([]Expr, 0, 10)
	if !p.check(tokenRightParen) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}
		args = append(args, e)
		for p.match(tokenComma) {
			if len(args) >= 255 {
				loxerr2(&Error{p.peek(), "can't have more than 255 arguments"})
			}
			e, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, e)
		}
	}
	paren, err := p.consume(tokenRightParen, "expect ')' after arguments.")
	return &Call{callee, paren, args}, err
}

func (p *Parser) primary() (Expr, *Error) {
	switch {
	case p.match(tokenFalse):
		return &Literal{false}, nil
	case p.match(tokenTrue):
		return &Literal{true}, nil
	case p.match(tokenNil):
		return &Literal{nil}, nil
	}

	if p.match(tokenNumber, tokenString) {
		return &Literal{p.previous().Literal}, nil
	}
	if p.match(tokenIdent) {
		return &Variable{p.previous()}, nil
	}
	if p.match(tokenLeftParen) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tokenRightParen, "expect ')' after expression")
		return &Grouping{e}, err
	}
	return nil, &Error{p.peek(), "expect expression or value"}
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

func (p *Parser) consume(typ int, message string) (Token, *Error) {
	if p.check(typ) {
		return p.advance(), nil
	}
	return Token{}, &Error{p.peek(), message}
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
