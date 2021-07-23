package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var definitions = map[string]string{
	`Binary`:   `left Expr, op *Token, right Expr`,
	`Grouping`: `expr Expr`,
	`Literal`:  `val interface{}`,
	`Unary`:    `op *Token, right Expr`,
}

func main() {
	f, err := os.Create(os.Args[1] + ".go")
	if err != nil {
		panic(err)
	}
	defast(f, "Expr", definitions)
}

func defast(where *os.File, iface string, types map[string]string) {
	fmt.Fprintln(where, `package main`)
	// defvisi(where, types)
	for n, f := range types {
		deftype(where, n, f)
	}
}

func deftype(w io.Writer, name string, fields string) {
	newfs := strings.Split(fields, ",")
	defcons(w, name, newfs)
	defstruct(w, name, newfs)
	defgetters(w, name, newfs)
	defmyvisi(w, name)
}

func defstruct(w io.Writer, name string, fields []string) {
	templ := `
type %s struct {
	%s
}
	`
	fmt.Fprintf(w, templ,
		strings.ToLower(name),
		strings.Join(fields, `;`),
	)
}

func defcons(w io.Writer, name string, fields []string) {
	templ := `
// %s
func %s(%s) Expr {
	return &%s{
		%s	
	}
}
	`
	fmt.Fprintf(w, templ, name,
		name, strings.Join(fields, ","),
		strings.ToLower(name),
		func() string {
			s := ""
			for _, v := range fields {
				b := strings.Fields(v)[0]
				s += b + `:` + b + ",\n"
			}
			return s
		}(),
	)

}

func defgetters(w io.Writer, name string, newfs []string) {
	templ := `
func (%c *%s) Left() Expr {return nil}
func (%c *%s) Op() *Token {return &Token{}}
func (%c *%s) Right() Expr {return nil}
`
	lowname := strings.ToLower(name)
	fmt.Fprintf(w, templ,
		lowname[0], lowname,
		lowname[0], lowname,
		lowname[0], lowname,
	)
}

// func defvisi(w io.Writer, types map[string]string) {
// 	fmt.Fprintln(w, `type ExprVisitor interface {`)
// 	for v := range types {
// 		fmt.Fprintf(w, "Visit%sExpr(*%s) interface{}\n", v, strings.ToLower(v))
// 	}
// 	fmt.Fprintln(w, `}`)
// }

func defmyvisi(w io.Writer, name string) {
	templ := `
func (%c *%s) Accept(vis ExprVisitor) interface{} {
	return vis.Visit(%c)
}
`
	lowname := strings.ToLower(name)
	fmt.Fprintf(w, templ,
		lowname[0], lowname,
		//name,
		lowname[0],
	)

}
