package main

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

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

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.next()
	}
	s.Tokens = append(s.Tokens, Token{Type: tokenEOF, Lexeme: []byte{}})
	return s.Tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func isdigit(r rune) bool {
	return r >= '0' && r <= '9'
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

func (s *Scanner) next() {
	r := s.advance()
	match1 := func(prim, alt int) {
		if s.match('=') {
			s.addToken(prim, nil)
		} else {
			s.addToken(alt, nil)
		}
	}
	switch r {
	case '\n':
		s.line++
	case '!':
		match1(tokenBangEqual, tokenBang)
	case '=':
		match1(tokenEqualEqual, tokenEqual)
	case '<':
		match1(tokenLessEqual, tokenLess)
	case '>':
		match1(tokenGreaterEqual, tokenGreater)
	case '"':
		s.string()
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tokenSlash, nil)
		}
	default:
		if isdigit(r) {
			s.number()
		} else if c, ok := singles[r]; ok {
			s.addToken(c, nil)
		} else if unicode.IsLetter(r) {
			s.ident()
		} else {
			loxerr(s.line, "unexpected character")
		}
	case ' ':
	case '\r':
	case '\t':
	}
}

func (s *Scanner) ident() {
	for unicode.IsDigit(s.peek()) || unicode.IsLetter(s.peek()) {
		s.advance()
	}

	text := s.Source[s.start:s.current]
	if typ, ok := keywords[string(text)]; ok {
		s.addToken(typ, nil)
		return
	}
	s.addToken(tokenIdent, nil)
}

func (s *Scanner) number() {
	for isdigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isdigit(s.peekNext()) {
		s.advance()
		for isdigit(s.peek()) {
			s.advance()
		}
	}
	f, _ := strconv.ParseFloat(string(s.Source[s.start:s.current]), 64)
	s.addToken(tokenNumber, f)

}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		loxerr(s.line, "unterminated string")
		return
	}
	s.advance()

	val := s.Source[1+s.start : s.current-1]
	s.addToken(tokenString, val)
}

func (s *Scanner) peekNext() rune {
	_, sz1 := utf8.DecodeRune(s.Source[s.current:])
	r, _ := utf8.DecodeRune(s.Source[s.current+sz1:])
	if s.current+sz1 >= len(s.Source) {
		return 0
	}
	return r
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.cur()
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.cur() != expected {
		return false
	}
	s.adv()
	return true
}

func (s *Scanner) cur() rune {
	r, _ := utf8.DecodeRune(s.Source[s.current:])
	return r
}

func (s *Scanner) adv() {
	_, sz := utf8.DecodeRune(s.Source[s.current:])
	s.current += sz
}

func (s *Scanner) advance() rune {
	r := s.cur()
	s.adv()
	return r
}

func (s *Scanner) addToken(typ int, literal interface{}) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, Token{typ, text, literal, s.line})
}
