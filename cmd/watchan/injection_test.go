package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func TestInspectBlockStmt(t *testing.T) {
	type args struct {
		index     int
		blockStmt *ast.BlockStmt
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InspectBlockStmt(tt.args.index, tt.args.blockStmt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InspectBlockStmt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_injection(t *testing.T) {

	tests := []struct {
		name         string
		gotFilePath  string
		wantFilePath string
	}{
		{
			name:         "normal",
			gotFilePath:  "testdata/example_got.go",
			wantFilePath: "testdata/example_want.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fset := token.NewFileSet()
			expr, err := parser.ParseFile(fset, tt.gotFilePath, nil, parser.ParseComments)
			if err != nil {
				fmt.Println("Error parsing source code:", err)
				return
			}

			injection(expr)

			fset2 := token.NewFileSet()
			expr2, err := parser.ParseFile(fset2, tt.wantFilePath, nil, parser.ParseComments)
			if err != nil {
				fmt.Println("Error parsing source code:", err)
				return
			}

			if expr != expr2 {
				t.Error("error")
			}
		})
	}
}
