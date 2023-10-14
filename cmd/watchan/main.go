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

	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, goFilePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing source code:", err)
		return
	}

	ast.Print(fset, expr)

	injection(fset, expr)

	// check output
	var buf bytes.Buffer
	err = format.Node(&buf, fset, expr)
	if err != nil {
		fmt.Println("Error generating code:", err)
		return
	}
	fmt.Println(buf.String())

	// write file
	modifiedFilePath := "modified_" + goFilePath
	err = os.WriteFile(modifiedFilePath, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing modified file:", err)
		return
	}

	fmt.Println("Modified file written to:", modifiedFilePath)
}
