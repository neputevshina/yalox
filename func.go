package main

type Func struct {
	declaration *Function
}

func (f *Func) Call(i *Interpreter, args []interface{}) interface{} {
	env := NewEnvironment(i.globals)
	for i := range f.declaration.Params {
		env.Define(string(f.declaration.Params[i].Lexeme), args[i])
	}
	ni := NewInterpreter(env)
	ni.Interpret(append(f.declaration.Body, &Return{}))
	return <-ni.ret
}

func (f *Func) Arity() int {
	return len(f.declaration.Params)
}

func (f *Func) String() string {
	return "<fn " + string(f.declaration.Name.Lexeme) + ">"
}
