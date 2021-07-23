package main

import "fmt"

type AstPrinter struct{}

func (ap *AstPrinter) Print(e Expr) string {
	return e.Accept(ap).(string)
}

func (ap *AstPrinter) VisitBinaryExpr(e *binary) interface{} {
	return parenth(e.op.Lexeme, e.left, e.right)
}

func (ap *AstPrinter) VisitGroupingExpr(e *grouping) interface{} {
	return parenth([]byte(`group`), e.expr)
}

func (ap *AstPrinter) VisitLiteralExpr(e *literal) interface{} {
	if e.val == nil {
		return nil
	}
	return fmt.Sprint(e.val)
}

func (ap *AstPrinter) VisitUnaryExpr(e *unary) interface{} {
	return parenth(e.op.Lexeme, e.right)
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
