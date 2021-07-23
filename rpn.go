package main

import "fmt"

type RPN struct {
	stack []Token
}

func (r *RPN) Print(e Expr) string {
	return e.Accept(r).(string)
}

func (r *RPN) Visit(l interface{}) interface{} {
	switch e := l.(type) {
	case *Binary:
		s := e.left.Accept(r).(string) + " "
		s += e.right.Accept(r).(string) + " "
		return s + string(e.op.Lexeme)
	case *Grouping:
		return fmt.Sprint(e.expr.Accept(r))
	case *Literal:
		return fmt.Sprint(e.val)
	case *Unary:
		return e.right.Accept(r).(string) + " " + string(e.op.Lexeme)
	default:
		return nil
	}
}

func rpnrpn(text []byte, es ...Expr) interface{} {
	s := "(" + string(text)
	for _, v := range es {
		s += " "
		s += v.Accept(&AstPrinter{}).(string)
	}
	s += ")"
	return s
}
