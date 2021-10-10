package cyclo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestCycloComp1(t *testing.T) {
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
			got := uint(getCycloComp(v, &Weights{If: 1}))

			if got != want {
				t.Fatalf(`GetCycloComp("package main...") = %v, Wanted %v`, got, want)
			}
		}
		return true
	})
}

func TestCycloComp2(t *testing.T) {
	fset := token.NewFileSet()

	srcCode := `package main
				func main() {
					for {
						for _, _ := range []int{1, 2, 3} {
							if (f1() && f2()) || f3() && f4() {
								switch v {
								case 0:
								case 1:
								default:
								}
							}
						}
					}
				}`
	expr, err := parser.ParseFile(fset, "", srcCode, parser.ParseComments)

	if err != nil {
		t.Fatal(err)
	}

	m := Metric{Config: Weights{If: 1, For: 1, Case: 1, Rng: 1, And: 1, Or: 1}}

	ast.Inspect(expr, func(n ast.Node) bool {
		m.ParseNode(n)
		return true
	})

	want := uint(10)
	got := uint(m.Finish())

	if got != want {
		t.Fatalf(`GetCycloComp("package main...") = %v, Wanted %v`, got, want)
	}
}
