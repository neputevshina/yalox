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
		s := e.Left.Accept(r).(string) + " "
		s += e.Right.Accept(r).(string) + " "
		return s + string(e.Op.Lexeme)
	case *Grouping:
		return fmt.Sprint(e.Expr.Accept(r))
	case *Literal:
		return fmt.Sprint(e.Val)
	case *Unary:
		return e.Right.Accept(r).(string) + " " + string(e.Op.Lexeme)
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
