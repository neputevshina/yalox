package main

import "fmt"

type Error struct {
	Token   Token
	Message string
}

type Interpreter struct {
	raise chan *Error
	env   *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		raise: make(chan *Error),
		env:   NewEnvironment(nil),
	}
}

func (i *Interpreter) Interpret(stmts []Stmt) {
	done := make(chan struct{})
	go func() {
		for _, s := range stmts {
			i.execute(s)
		}
		done <- struct{}{}
	}()
	select {
	case err := <-i.raise:
		loxerr2(err)
	case <-done:
		return
	}
}

func (i *Interpreter) execute(s Stmt) {
	_ = s.Accept(i)
}

func (i *Interpreter) evaluate(e Expr) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) Visit(v interface{}) interface{} {
	switch a := v.(type) {
	case *Block:
		i.executeBlock(a.Stmts, NewEnvironment(i.env))
		return nil
	case *Expression:
		i.evaluate(a.Expr)
		return nil
	case *Print:
		v := i.evaluate(a.Expr)
		fmt.Println(stringify(v))
		return nil
	case *Var:
		var val interface{}
		if a.Init != nil {
			val = i.evaluate(a.Init)
		}
		i.env.Define(string(a.Name.Lexeme), val)
		return nil
		//
	case *Literal:
		return a.Val
	case *Grouping:
		return i.evaluate(a.Expr)
	case *Unary:
		r := i.evaluate(a.Right)
		switch a.Op.Type {
		case tokenMinus:
			return -i.maybefloat(a.Op, r)
		case tokenBang:
			return !istruthy(r)
		}
		panic("unreachable")
	case *Binary:
		l := i.evaluate(a.Left)
		r := i.evaluate(a.Right)
		switch a.Op.Type {
		case tokenMinus:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl - fr
		case tokenSlash:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl / fr
		case tokenStar:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl * fr
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
			i.raise <- &Error{a.Op, "both operands must be either strings or numbers"}
			return ""
		case tokenGreater:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl > fr
		case tokenGreaterEqual:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl >= fr
		case tokenLess:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl < fr
		case tokenLessEqual:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl <= fr
		case tokenEqualEqual:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl == fr
		case tokenBangEqual:
			fl, fr := i.maybefloats(a.Op, l, r)
			return fl != fr
		}
	case *Variable:
		v, err := i.env.Get(a.Name)
		if err != nil {
			i.raise <- err
		}
		return v
	case *Assign:
		value := i.evaluate(a.Val)
		i.env.Assign(a.Name, value)
		return value
	}
	panic("unreachable")
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) {
	// i cross my fingers
	prev := i.env
	defer func() { i.env = prev }()
	i.env = env
	for _, s := range stmts {
		i.execute(s)
	}
}

func (i *Interpreter) maybefloat(t Token, v interface{}) float64 {
	f, k := v.(float64)
	if !k {
		i.raise <- &Error{t, "operand must be a number"}
	}
	return f
}

func (i *Interpreter) maybefloats(t Token, lv interface{}, rv interface{}) (float64, float64) {
	l, kl := lv.(float64)
	r, kr := rv.(float64)
	if !(kl && kr) {
		i.raise <- &Error{t, "operands must be numbers"}
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

var _ = Visitor(&Interpreter{})
