package main

type Func struct {
	declaration *Function
}

func (f *Func) Call(i *Interpreter, args []interface{}) interface{} {
	env := NewEnvironment(i.globals)
	for i := range f.declaration.Params {
		env.Define(string(f.declaration.Params[i].Lexeme), args[i])
	}
	// cause omitting this will cause deadlock
	body := append(f.declaration.Body, &Return{})
	go i.executeBlock(body, env)
	return <-i.ret
}

func (f *Func) Arity() int {
	return len(f.declaration.Params)
}

func (f *Func) String() string {
	return "<fn " + string(f.declaration.Name.Lexeme) + ">"
}
