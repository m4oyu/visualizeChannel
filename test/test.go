package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {

	expr, err := parser.ParseExpr(`func(x, y int){}(10, 20)`)
	if err != nil {
		log.Fatal(err)
	}
	n := astutil.Apply(expr, func(cr *astutil.Cursor) bool {
		if cr.Name() == "Args" && cr.Index() == 0 {
			cr.InsertBefore(&ast.BasicLit{
				Kind:  token.STRING,
				Value: "hi",
			})
			cr.InsertAfter(&ast.BasicLit{
				Kind:  token.STRING,
				Value: "gopher",
			})
		}
		return true
	}, nil)

	if err := format.Node(os.Stdout, token.NewFileSet(), n); err != nil {
		log.Fatalln("Error:", err)
	}
	fmt.Println()

}
