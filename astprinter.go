package main

import "fmt"

// AstPrinter is an s-expression pretty printer
type AstPrinter struct{}

// Print converts an Expr to the s-expression representation of it.
func (ap *AstPrinter) Print(e Expr) string {
	return e.Accept(ap).(string)
}

// Visit conforms ExprVisitor.
// Since Go does have type switches, we don't need to write visitors in their pure form.
func (ap *AstPrinter) Visit(l interface{}) interface{} {
	switch e := l.(type) {
	case *Binary:
		return parenth(e.Op.Lexeme, e.Left, e.Right)
	case *Grouping:
		return parenth([]byte(`group`), e.Expr)
	case *Literal:
		if e.Val == nil {
			return nil
		}
		return fmt.Sprint(e.Val)
	case *Unary:
		return parenth(e.Op.Lexeme, e.Right)
	}
	panic("unreachable")
}

func parenth(text []byte, es ...Expr) interface{} {
	s := "(" + string(text)
	for _, v := range es {
		s += " "
		s += v.Accept(&AstPrinter{}).(string)
	}
	s += ")"
	return s
}

func apmain() {
	expr :=
		&Binary{
			Left: &Unary{
				Op:    Token{Type: tokenMinus, Lexeme: []byte{'-'}},
				Right: &Literal{123},
			},
			Op: Token{Type: tokenStar, Lexeme: []byte{'*'}},
			Right: &Grouping{
				Expr: &Literal{46.67},
			},
		}
	expr2 := &Binary{
		&Grouping{
			&Binary{
				&Literal{1},
				Token{Type: tokenPlus, Lexeme: []byte{'+'}},
				&Literal{2},
			},
		},
		Token{Type: tokenStar, Lexeme: []byte{'*'}},
		&Grouping{
			&Binary{
				&Literal{4},
				Token{Type: tokenMinus, Lexeme: []byte{'-'}},
				&Literal{3},
			},
		},
	}
	fmt.Println((&AstPrinter{}).Print(expr))
	fmt.Println((&RPN{}).Print(expr2))
}
