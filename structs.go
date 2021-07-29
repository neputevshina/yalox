package main

type Binary struct {
	Left  Expr
	Op    Token
	Right Expr
}

type Grouping struct {
	Expr Expr
}

type Literal struct {
	Val interface{}
}

type Unary struct {
	Op    Token
	Right Expr
}

type Variable struct {
	Name Token
}

type Expression struct {
	Expr Expr
}

type Print struct {
	Expr Expr
}

type Var struct {
	Name Token
	Init Expr
}

type Assign struct {
	Name Token
	Val  Expr
}

type Block struct {
	Stmts []Stmt
}
