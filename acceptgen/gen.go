package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	out, err := os.Create(os.Args[2] + ".go")
	out.WriteString(`
// Generated, DO NOT EDIT. Name is intentional.
package main
`)
	if err != nil {
		panic(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, os.Args[1]+".go", nil, parser.DeclarationErrors)
	if err != nil {
		panic(err)
	}

	for _, d := range f.Scope.Objects {
		if d.Kind == ast.Typ {
			// I have tried to do and tired doing some fuckery with Go's AST package and has decided to use brute force.
			name := d.Name
			i0, _ := utf8.DecodeRuneInString(d.Name)
			ini := unicode.ToLower(i0)
			out.WriteString(fmt.Sprintf(`
func (%c *%s) Accept(vis Visitor) interface{} {
	return vis.Visit(%c)
}
`,
				ini, name, ini),
			)
		}
	}
}
