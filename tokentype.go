package main

// Enum TokenType in the book.
const (
	returnMe = -1
	// One character
	_ = iota
	tokenLeftParen
	tokenRightParen
	tokenLeftBrace
	tokenRightBrace
	tokenComma
	tokenDot
	tokenMinus
	tokenPlus
	tokenSemicolon
	tokenSlash
	tokenStar

	// Don't move these!
	// One
	tokenBang
	tokenEqual
	tokenGreater
	tokenLess
	// ...or two
	tokenBangEqual
	tokenEqualEqual
	tokenGreaterEqual
	tokenLessEqual

	// Literals
	tokenIdent
	tokenString
	tokenNumber

	// KWs
	tokenAnd
	tokenOr
	tokenIf
	tokenElse
	tokenFor
	tokenTrue
	tokenFalse
	tokenNil
	tokenVar
	tokenWhile
	tokenThis
	tokenFun
	tokenClass
	tokenReturn
	tokenSuper
	tokenPrint

	// EOF because it's handy
	tokenEOF
)
