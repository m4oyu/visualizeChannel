package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mytool <go-file>")
		return
	}

	goFilePath := os.Args[1]
	fmt.Println(goFilePath)

	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, goFilePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing source code:", err)
		return
	}

	ast.Print(fset, expr)

	// Inject channel operation into the file content
	ast.Inspect(expr, func(n ast.Node) bool {

		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			fmt.Println("modify import")
		}

		if funcDecl, ok := n.(*ast.FuncDecl); ok {

			for i := 0; i < len(funcDecl.Body.List); i++ {
				InspectFunc(i, funcDecl)
			}

		}
		return true
	})

	var buf bytes.Buffer
	err = format.Node(&buf, fset, expr)
	if err != nil {
		fmt.Println("Error generating code:", err)
		return
	}
	fmt.Println(buf.String())

	modifiedFilePath := "modified_" + goFilePath
	err = os.WriteFile(modifiedFilePath, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing modified file:", err)
		return
	}

	fmt.Println("Modified file written to:", modifiedFilePath)

}

func InspectFunc(index int, funcDecl *ast.FuncDecl) {
	ast.Inspect(funcDecl.Body.List[index], func(n ast.Node) bool {

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
			return true
		}

		// send
		if sendStmt, ok := n.(*ast.SendStmt); ok {
			newSend := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{
						Name: "chanx.Send",
					},
					Args: []ast.Expr{
						sendStmt.Value,
					},
				},
			}
			funcDecl.Body.List[index] = newSend
			return true
		}

		// recv

		// AssignStmt [recv := <-ch]
		// TODO
		// ExprStmt
		//  -> UnaryExpr
		//  -> CallExpr.Args
		// Select文での使用,
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			for i, v := range assignStmt.Rhs {
				if unaryExpr, ok := v.(*ast.UnaryExpr); ok && unaryExpr.Op == token.ARROW {
					newRecv := &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "chanx"},
							Sel: &ast.Ident{Name: "Recv"},
						},
					}
					assignStmt.Rhs[i] = newRecv
				}
			}
			return true
		}
		if exprStmt, ok := n.(*ast.ExprStmt); ok {
			if unaryExpr, ok := exprStmt.X.(*ast.UnaryExpr); ok && unaryExpr.Op == token.ARROW {
				newRecv := &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "chanx"},
						Sel: &ast.Ident{Name: "Recv"},
					},
				}
				exprStmt.X = newRecv
			}
			return true
		}

		return true
	})
}
