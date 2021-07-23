package main

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Scanner is a lexical analyzer class.
type Scanner struct {
	Source []byte
	Tokens []Token

	start, current, line int
}

// NewScanner is a constructor for Scanner.
func NewScanner(source []byte) *Scanner {
	return &Scanner{
		Source: source,
		line:   1,
	}
}

// ScanTokens lexes the input and returns a slice of tokens.
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.Tokens = append(s.Tokens, Token{Type: tokenEOF})
	return s.Tokens
}

// Characters that don't require special treatment.
var singles = map[rune]int{
	'(': tokenLeftParen,
	')': tokenRightParen,
	'{': tokenLeftBrace,
	'}': tokenRightBrace,
	',': tokenComma,
	'.': tokenDot,
	'-': tokenMinus,
	'+': tokenPlus,
	';': tokenSemicolon,
	'*': tokenStar,
}

func (s *Scanner) scanToken() {
	r := s.advance()
	// Trying to compress repeating code.
	match1 := func(prim, alt int) {
		if s.match('=') {
			s.addToken(prim, nil)
		} else {
			s.addToken(alt, nil)
		}
	}
	switch r {
	case '!':
		match1(tokenBangEqual, tokenBang)
	case '=':
		match1(tokenEqualEqual, tokenEqual)
	case '<':
		match1(tokenLessEqual, tokenLess)
	case '>':
		match1(tokenGreaterEqual, tokenGreater)
	case '/':
		if s.match('/') {
			// Skip comments.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') { // Challenge 4.4: add multiline comments.
			// Intentionally won't work in REPL mode.
			for !(s.peek() == '*' && s.peekNext() == '/') {
				s.advance()
			}
		} else {
			s.addToken(tokenSlash, nil)
		}

	case ' ':
	case '\r':
	case '\t': // skip
		break

	case '\n':
		s.line++

	case '"':
		s.str()

	default:
		if isdigit(r) {
			s.number()
		} else if unicode.IsLetter(r) {
			// I'm cheating here, but it will allow e.g. cyrillic idents.
			s.ident()
		} else if v, k := singles[r]; k {
			s.addToken(v, nil)
		} else {
			loxerr(s.line, "unexpected character")
		}
	}
}

var keywords = map[string]int{
	"and":    tokenAnd,
	"class":  tokenClass,
	"else":   tokenElse,
	"false":  tokenFalse,
	"for":    tokenFor,
	"fun":    tokenFun,
	"if":     tokenIf,
	"nil":    tokenNil,
	"or":     tokenOr,
	"print":  tokenPrint,
	"return": tokenReturn,
	"super":  tokenSuper,
	"this":   tokenThis,
	"true":   tokenTrue,
	"var":    tokenVar,
	"while":  tokenWhile,
}

// Identifier parser.
func (s *Scanner) ident() {
	r := s.peek()
	for unicode.IsNumber(r) || unicode.IsLetter(r) {
		r = s.advance()
	}
	text := s.Source[s.start:s.current]
	typ, k := keywords[string(text)]
	if !k {
		typ = tokenIdent
	}
	s.addToken(typ, nil)
}

// Number parser.
func (s *Scanner) number() {
	for isdigit(s.peek()) {
		s.advance()
	}
	// Look for fractional part.
	if s.peek() == '.' && isdigit(s.peekNext()) {
		// Consume the dot.
		s.advance()
		for isdigit(s.peek()) {
			s.advance()
		}
	}
	n, err := strconv.ParseFloat(string(s.Source[s.start:s.current]), 64)
	if err != nil {
		loxerr(s.line, "incorrect number")
	}
	s.addToken(tokenNumber, n)
}

func (s *Scanner) peekNext() rune {
	// Pretty cursed.
	_, sz1 := utf8.DecodeRune(s.Source[s.current:])
	r, sz2 := utf8.DecodeRune(s.Source[s.current+sz1:])
	if sz2 >= len(s.Source) {
		return 0
	}
	return r
}

// String parser.
func (s *Scanner) str() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		loxerr(s.line, "unterminated string")
	}
	// The closing ".
	s.advance()
	// Trim the surrounding quotes and add token.
	s.addToken(tokenString, s.Source[1:len(s.Source)-1])
}

func (s *Scanner) addToken(tokentype int, literal interface{}) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, Token{
		Type:    tokentype,
		Lexeme:  text,
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func isdigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Get a current rune.
func (s *Scanner) getr() rune {
	r, _ := utf8.DecodeRune(s.Source[s.current:])
	return r
}

// Advance current position by one rune.
func (s *Scanner) skip() {
	_, sz := utf8.DecodeRune(s.Source[s.current:])
	s.current += sz
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.getr() != expected {
		return false
	}
	s.skip()
	return true
}

func (s *Scanner) advance() rune {
	r := s.getr()
	s.skip()
	return r
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.getr()
}
