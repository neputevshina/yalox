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

type If struct {
	Cond Expr
	Then Expr
	Else Expr
}

type Logical struct {
	Left  Expr
	Op    Token
	Right Expr
}

type While struct {
	Cond Expr
	Body Stmt
}
