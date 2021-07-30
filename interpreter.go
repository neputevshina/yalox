package main

import (
	"fmt"
	"time"
)

type Error struct {
	Token   Token
	Message string
}

type Interpreter struct {
	raise   chan *Error
	env     *Environment
	globals *Environment
	ret     chan interface{}
}

type nf_clock struct{}

func (*nf_clock) Call(i *Interpreter, args []interface{}) interface{} {
	return time.Now().Unix()
}

func (*nf_clock) Arity() int {
	return 0
}

func (*nf_clock) String() string {
	return "<native fn>"
}

func NewInterpreter(env *Environment) *Interpreter {
	i := &Interpreter{
		raise:   make(chan *Error),
		globals: env,
		ret:     make(chan interface{}, 0),
	}
	i.env = i.globals
	i.globals.Define("clock", &nf_clock{})
	return i
}

func (i *Interpreter) Interpret(stmts []Stmt) {
	for _, s := range stmts {
		err := i.exec(s)
		if err != nil {
			loxerr2(err)
		}
	}
}

func (i *Interpreter) exec(s Stmt) *Error {
	_, err := s.Accept(i)
	return err
}

func (i *Interpreter) eval(e Expr) (interface{}, *Error) {
	return e.Accept(i)
}

func (i *Interpreter) Visit(v interface{}) (interface{}, *Error) {
	switch a := v.(type) {
	case *If:
		v, err := i.eval(a.Cond)
		if err != nil {
			return nil, err
		}
		if istruthy(v) {
			err = i.exec(a.Then)
		} else if a.Else != nil {
			err = i.exec(a.Else)
		}
		return nil, err
	case *Return:
		val := interface{}(nil)
		var err *Error
		if a.Value != nil {
			val, err = i.eval(a.Value)
		}
		if err != nil {
			err = &Error{Token{Type: returnMe, Literal: val}, ""}
		}
		return nil, err
	case *Block:
		i.executeBlock(a.Stmts, NewEnvironment(i.env))
		return nil, nil
	case *Expression:
		_, err := i.eval(a.Expr)
		return nil, err
	case *Function:
		fn := &Func{a}
		i.env.Define(string(a.Name.Lexeme), fn)
		return nil, nil
	case *Print:
		v, err := i.eval(a.Expr)
		if err == nil {
			fmt.Println(stringify(v))
		}
		return nil, err
	case *Var:
		var val interface{}
		var err *Error
		if a.Init != nil {
			val, err = i.eval(a.Init)
		}
		i.env.Define(string(a.Name.Lexeme), val)
		return nil, err
	case *While:
		v, err := i.eval(a.Cond)
		if err != nil {
			return nil, err
		}
		for istruthy(v) {
			err := i.exec(a.Body)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
		//
	case *Literal:
		return a.Val, nil
	case *Logical:
		l, err := i.eval(a.Left)
		if err != nil {
			return nil, err
		}
		switch a.Op.Type {
		case tokenOr:
			if istruthy(l) {
				return l, nil
			}
		case tokenAnd:
			if !istruthy(l) {
				return l, nil
			}
		default:
			panic("unreachable")
		}
		return i.eval(a.Right)
	case *Grouping:
		return i.eval(a.Expr)
	case *Unary:
		r, err := i.eval(a.Right)
		switch a.Op.Type {
		case tokenMinus:
			v, err := i.maybefloat(a.Op, r)
			return -v, err
		case tokenBang:
			return !istruthy(r), err
		}
		panic("unreachable")
	case *Binary:
		l, err := i.eval(a.Left)
		if err != nil {
			return nil, err
		}
		r, err := i.eval(a.Right)
		if err != nil {
			return nil, err
		}
		switch a.Op.Type {
		case tokenMinus:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl - fr, err
		case tokenSlash:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl / fr, err
		case tokenStar:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl * fr, err
		case tokenPlus:
			switch le := l.(type) {
			case float64:
				if re, k := r.(float64); k {
					return le + re, nil
				}
				goto fail
			case []byte:
				if re, k := r.([]byte); k {
					return append(le, re...), nil
				}
				goto fail
			}
		fail:
			return nil, &Error{a.Op, "both operands must be either strings or numbers"}
		case tokenGreater:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl > fr, err
		case tokenGreaterEqual:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl >= fr, err
		case tokenLess:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl < fr, err
		case tokenLessEqual:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl <= fr, err
		case tokenEqualEqual:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl == fr, err
		case tokenBangEqual:
			fl, fr, err := i.maybefloats(a.Op, l, r)
			return fl != fr, err
		}
	case *Call:
		callee, err := i.eval(a.Callee)
		if err != nil {
			return nil, err
		}

		args := make([]interface{}, 0, 10)
		for _, ar := range a.Args {
			v, err := i.eval(ar)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		fn, k := callee.(Callable)
		if !k {
			return nil, &Error{a.Paren, "can only call functions and classes"}
		}
		if len(args) != fn.Arity() {
			return nil, &Error{a.Paren, fmt.Sprintf("expected %d arguments but got %d", len(args), fn.Arity())}
		}
		return fn.Call(i, args)

	case *Variable:
		v, err := i.env.Get(a.Name)
		return v, err
	case *Assign:
		value, err := i.eval(a.Val)
		if err != nil {
			return nil, err
		}
		err = i.env.Assign(a.Name, value)
		return value, err
	}
	panic("unreachable")
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) *Error {
	// i cross my fingers
	prev := i.env
	defer func() { i.env = prev }()
	i.env = env
	for _, s := range stmts {
		err := i.exec(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) maybefloat(t Token, v interface{}) (float64, *Error) {
	var err *Error
	f, k := v.(float64)
	if !k {
		err = &Error{t, "operand must be a number"}
	}
	return f, err
}

func (i *Interpreter) maybefloats(t Token, lv interface{}, rv interface{}) (float64, float64, *Error) {
	var err *Error
	l, kl := lv.(float64)
	r, kr := rv.(float64)
	if !(kl && kr) {
		err = &Error{t, "operands must be numbers"}
	}
	return l, r, err
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
