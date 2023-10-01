package main

import (
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
	type args struct {
		expr *ast.File
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injection(tt.args.expr)
		})
	}
}
