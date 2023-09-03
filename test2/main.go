package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

func main() {
	sourceCode := `
		package main

		import "fmt"

		func main() {
			ch := make(chan int, 10)
			ch <- 5
			fmt.Println(<-ch)
		}
	`

	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, "example.go", sourceCode, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing source code:", err)
		return
	}

	ast.Print(fset, expr)

	ast.Inspect(expr, func(n ast.Node) bool {

		// make
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if funIdent, ok := callExpr.Fun.(*ast.Ident); ok && funIdent.Name == "make" {
				if _, ok := callExpr.Args[0].(*ast.ChanType); ok {

					chanSize := "1"
					if len(callExpr.Args) > 1 {
						chanSize = callExpr.Args[1].(*ast.BasicLit).Value
					}

					newCall := &ast.CallExpr{
						Fun: &ast.Ident{Name: "chanx.Make"},
						Args: []ast.Expr{
							&ast.Ident{Name: callExpr.Args[0].(*ast.ChanType).Value.(*ast.Ident).Name},
							&ast.Ident{Name: chanSize},
						},
					}
					*callExpr = *newCall
				}
			}
		}

		// send
		if _, ok := n.(*ast.SendStmt); ok {

			// newSend := &ast.CallExpr{
			// 	Fun: &ast.Ident{
			// 		Name: "chanx.Send",
			// 	},
			// 	Args: []ast.Expr{
			// 		&ast.BasicLit{Value: sendStmt.Value.(*ast.BasicLit).Value},
			// 	},
			// }

			// parentBlock := n.
			fmt.Printf("%+v\n", n)
		}

		// recv

		return true
	})

	var buf bytes.Buffer
	err = format.Node(&buf, fset, expr)
	if err != nil {
		fmt.Println("Error generating code:", err)
		return
	}
	fmt.Println(buf.String())

}
