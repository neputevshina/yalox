package main

import "fmt"

type Error struct {
	Token   Token
	Message string
}

type Interpreter struct {
	raise chan Error
}

func NewInterpreter() *Interpreter {
	return &Interpreter{raise: make(chan Error)}
}

func (i *Interpreter) Interpret(e Expr) {
	done := make(chan struct{})
	var value interface{}
	go func() {
		value = i.evaluate(e)
		done <- struct{}{}
	}()
	select {
	case <-done:
		fmt.Println(stringify(value))
	case err := <-i.raise:
		loxerr2(err)
	}
}

func (i *Interpreter) evaluate(e Expr) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) Visit(v interface{}) interface{} {
	switch e := v.(type) {
	case *Literal:
		return e.Val
	case *Grouping:
		return i.evaluate(e.Expr)
	case *Unary:
		r := i.evaluate(e.Right)
		switch e.Op.Type {
		case tokenMinus:
			return -i.maybefloat(e.Op, r)
		case tokenBang:
			return !istruthy(r)
		}
		panic("unreachable")
	case *Binary:
		l := i.evaluate(e.Left)
		r := i.evaluate(e.Right)
		switch e.Op.Type {
		case tokenMinus:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fr - fl
		case tokenSlash:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fr / fl
		case tokenStar:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fr * fl
		case tokenPlus:
			switch le := l.(type) {
			case float64:
				if re, k := r.(float64); k {
					return le + re
				}
				goto fail
			case []byte:
				if re, k := r.([]byte); k {
					return append(le, re...)
				}
				goto fail
			}
		fail:
			i.raise <- Error{e.Op, "both operands must be either strings or numbers"}
			return ""
		case tokenGreater:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fl > fr
		case tokenGreaterEqual:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fl >= fr
		case tokenLess:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fl < fr
		case tokenLessEqual:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fl <= fr
		case tokenEqual:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fl == fr
		case tokenBangEqual:
			fl, fr := i.maybefloats(e.Op, l, r)
			return fl != fr
		}
	}
	panic("unreachable")
}

func (i *Interpreter) maybefloat(t Token, v interface{}) float64 {
	f, k := v.(float64)
	if !k {
		i.raise <- Error{t, "operand must be a number"}
	}
	return f
}

func (i *Interpreter) maybefloats(t Token, lv interface{}, rv interface{}) (float64, float64) {
	l, kl := lv.(float64)
	r, kr := rv.(float64)
	if !(kl && kr) {
		i.raise <- Error{t, "operands must be numbers"}
	}
	return l, r
}

func istruthy(v interface{}) bool {
	if b, k := v.(bool); k {
		return b
	} else if v == nil {
		return false
	}
	return true
}

func isequal(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	}
	return a == b
}

func stringify(v interface{}) string {
	if v == nil {
		return "nil"
	}
	// Go doesn't print decimal automatically, but byte slices must be treated
	if s, k := v.([]byte); k {
		return string(s)
	}
	return fmt.Sprint(v)
}

var _ = ExprVisitor(&Interpreter{})
