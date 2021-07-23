package main

import "fmt"

type AstPrinter struct{}

func (ap *AstPrinter) Print(e Expr) string {
	return e.Accept(ap).(string)
}

// Visit conforms ExprVisitor.
// Since Go does have type switches, we don't need to write visitors in their pure form.
func (ap *AstPrinter) Visit(l interface{}) interface{} {
	switch e := l.(type) {
	case *binary:
		return parenth(e.op.Lexeme, e.left, e.right)
	case *grouping:
		return parenth([]byte(`group`), e.expr)
	case *literal:
		if e.val == nil {
			return nil
		}
		return fmt.Sprint(e.val)
	case *unary:
		return parenth(e.op.Lexeme, e.right)
	default:
		return nil
	}
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
		Binary(
			Unary(
				&Token{Type: tokenMinus, Lexeme: []byte{'-'}},
				Literal(123),
			),
			&Token{Type: tokenStar, Lexeme: []byte{'*'}},
			Grouping(
				Literal(45.67),
			),
		)
	fmt.Println((&AstPrinter{}).Print(expr))
}
