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
		return parenth(e.op.Lexeme, e.left, e.right)
	case *Grouping:
		return parenth([]byte(`group`), e.expr)
	case *Literal:
		if e.val == nil {
			return nil
		}
		return fmt.Sprint(e.val)
	case *Unary:
		return parenth(e.op.Lexeme, e.right)
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
		NewBinary(
			NewUnary(
				Token{Type: tokenMinus, Lexeme: []byte{'-'}},
				NewLiteral(123),
			),
			Token{Type: tokenStar, Lexeme: []byte{'*'}},
			NewGrouping(
				NewLiteral(45.67),
			),
		)
	expr2 := NewBinary(
		NewGrouping(
			NewBinary(
				NewLiteral(1),
				Token{Type: tokenPlus, Lexeme: []byte{'+'}},
				NewLiteral(2),
			),
		),
		Token{Type: tokenStar, Lexeme: []byte{'*'}},
		NewGrouping(
			NewBinary(
				NewLiteral(4),
				Token{Type: tokenMinus, Lexeme: []byte{'-'}},
				NewLiteral(3),
			),
		),
	)
	fmt.Println((&AstPrinter{}).Print(expr))
	fmt.Println((&RPN{}).Print(expr2))
}
