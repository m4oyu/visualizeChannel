package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"testing"
)

// func TestInspectBlockStmt(t *testing.T) {
// 	type args struct {
// 		index     int
// 		blockStmt *ast.BlockStmt
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := InspectBlockStmt(tt.args.index, tt.args.blockStmt); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("InspectBlockStmt() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

var fset = token.NewFileSet()

func parse(t *testing.T, name, in string) *ast.File {
	file, err := parser.ParseFile(fset, name, in, parser.ParseComments)
	if err != nil {
		t.Fatalf("%s parse: %v", name, err)
	}
	return file
}

func print(t *testing.T, name string, f *ast.File) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		t.Fatalf("%s gofmt: %v", name, err)
	}
	return buf.String()
}

type test struct {
	name string
	in   string
	out  string
}

var tests = []test{
	{
		name: "add import",
		in: `package main

import (
	"fmt"
	"time"
)
`,
		out: `package main

import (
	"fmt"
	"github.com/m4oyu/visualizeChannel/chanx"
	"time"
)
`,
	},
	{
		name: "replace chanx.Make(), size 1",
		in: `package main

func main() {
	ch := make(chan int)
}
`,
		out: `package main

import "github.com/m4oyu/visualizeChannel/chanx"

func main() {
	ch := chanx.Make(int, 1)
}
`,
	},
	{
		name: "replace chanx.Make(), size 1",
		in: `package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 1)
	ch <- 1
	fmt.Println(<-ch)
	close(ch)
}
`,
		out: `package main

import (
	"fmt"
	"github.com/m4oyu/visualizeChannel/chanx"
)

func main() {
	ch := chanx.Make(int, 1)
	chanx.Send(1)
	fmt.Println(<-ch)
	close(ch)
}
`,
	},
}

func Test_injection(t *testing.T) {
	for _, test := range tests {
		file := parse(t, test.name, test.in)
		var before bytes.Buffer
		ast.Fprint(&before, fset, file, nil)

		injection(fset, file)
		if got := print(t, test.name, file); got != test.out {
			t.Errorf("first run: %s:\ngot: %s\nwant: %s", test.name, got, test.out)
			var after bytes.Buffer
			ast.Fprint(&after, fset, file, nil)
			t.Logf("AST before:\n%s\nAST after:\n%s\n", before.String(), after.String())
		}
	}
}
