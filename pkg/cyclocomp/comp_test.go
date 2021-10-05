package cyclo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestGetCycloComp(t *testing.T) {
	fset := token.NewFileSet()

	srcCode := `package main
				func main() {
					if c1()  { 
						f1() 
					} else { 
						f2() 
					} 
					if c2() { 
						f3() 
					} else { 
						f4() 
					}
				}`
	expr, err := parser.ParseFile(fset, "", srcCode, parser.ParseComments)

	if err != nil {
		t.Fatal(err)
	}

	ast.Inspect(expr, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.FuncDecl:
			want := uint(3)
			got := GetCycloComp(v, &Weights{If: 1})

			if got != want {
				t.Fatalf(`GetCycloComp("package main...") = %v, Wanted %v`, got, want)
			}
		}
		return true
	})
}
