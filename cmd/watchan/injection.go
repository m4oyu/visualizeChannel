package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

func injection(expr *ast.File) {
	// Inject channel operation into the file content
	ast.Inspect(expr, func(n ast.Node) bool {

		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			fmt.Println("modify import")
		}

		if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name.Name == "main" {
			for i, _ := range funcDecl.Body.List {
				InspectBlockStmt(i, funcDecl.Body)
			}
		}
		return true
	})
}

func InspectBlockStmt(index int, blockStmt *ast.BlockStmt) bool {
	ast.Inspect(blockStmt.List[index], func(n ast.Node) bool {
		// block roop
		if blockStmt, ok := n.(*ast.BlockStmt); ok {
			for i, _ := range blockStmt.List {
				InspectBlockStmt(i, blockStmt)
			}
		}

		// make
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if funIdent, ok := callExpr.Fun.(*ast.Ident); ok && funIdent.Name == "make" {
				if _, ok := callExpr.Args[0].(*ast.ChanType); ok {
					chanSize := "1"
					if len(callExpr.Args) > 1 {
						chanSize = callExpr.Args[1].(*ast.BasicLit).Value
					}
					*callExpr = ast.CallExpr{
						Fun: &ast.Ident{Name: "chanx.Make"},
						Args: []ast.Expr{
							&ast.Ident{Name: callExpr.Args[0].(*ast.ChanType).Value.(*ast.Ident).Name},
							&ast.Ident{Name: chanSize},
						},
					}
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

			blockStmt.List[index] = newSend
			return true
		}

		// recv
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

	return true

}
