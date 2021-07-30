package main

// Binary is a binary expression node
type Binary struct {
	Left  Expr
	Op    Token
	Right Expr
}

// Grouping is an expression inside parentheses
type Grouping struct {
	Expr Expr
}

// Literal is literal value in code
type Literal struct {
	Val interface{}
}

// Unary is an unary operation node
type Unary struct {
	Op    Token
	Right Expr
}

// Variable is a variable name in expression
type Variable struct {
	Name Token
}

// Expression is autological
type Expression struct {
	Expr Expr
}

// Print is a print statement node
type Print struct {
	Expr Expr
}

// Var is a variable declaration
type Var struct {
	Name Token
	Init Expr
}

// Assign is an assignment statement
type Assign struct {
	Name Token
	Val  Expr
}

// Block is a block statement: a statement comprising a list of statements
type Block struct {
	Stmts []Stmt
}

// If statement
type If struct {
	Cond Expr
	Then Stmt
	Else Stmt
}

// Logical operators: “and” and “or”
type Logical struct {
	Left  Expr
	Op    Token
	Right Expr
}

// While loop
type While struct {
	Cond Expr
	Body Stmt
}

// Call inside expression
type Call struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

// Function declaration
type Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

// Return statement
type Return struct {
	Keyword Token
	Value   Expr
}
