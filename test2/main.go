package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	sourceCode := `
		package main

		import "fmt"

		func main() {
			ch := make(chan int, 10)
			fmt.Println(ch)
		}
	`

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "example.go", sourceCode, parser.ParseComments)
	if err != nil {

		fmt.Println("Error parsing source code:", err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok || len(callExpr.Args) != 2 {
			return true
		}

		fun, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || fun.Sel.Name != "make" {
			return true
		}

		chanType, ok := callExpr.Args[0].(*ast.ChanType)
		if !ok {
			return true
		}

		chanName := ""
		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}
			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok || len(valueSpec.Names) != 1 {
					continue
				}
				chanName = valueSpec.Names[0].Name
			}
		}

		fmt.Println("Channel Name:", chanName)
		fmt.Println("Channel Type:", chanType.Value)
		fmt.Println("Channel Size:", callExpr.Args[1])

		return true
	})
}
