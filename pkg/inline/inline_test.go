package inline

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestInlineData(t *testing.T) {
	src := `package main

import (
	"fmt"
)
	

func foo(t int) float64 {
	a := 123
	b := 1.23
	return float64(a+t) + b
}
		
func main(){
	fmt.Printf("Result: %v", foo(1))
}`

	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	m := Metric{Config: Weights{Int: 2, Float: 2, Imag: 1, Char: 2, String: 1}}

	ast.Inspect(file, func(n ast.Node) bool {
		m.ParseNode(n)
		return true
	})
	got := m.Comp
	want := uint(33)
	if want != got {
		t.Errorf("TestCodeStructComp got %v, want %v", got, want)
	}
}
func TestInlineData2(t *testing.T) {
	src := `package main

	import (
		"fmt"
	)
	
	func foo(t int) complex128 {
		a := 'a'
		b := 32i
		result := (complex(float64(t), 0) + b)
		fmt.Printf("%s = %v\n", string(a), result)
		return result
	}
	
	func main() {
		fmt.Printf("Result: %v\n", foo(1))
	}
	`

	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	m := Metric{Config: Weights{Int: 2, Float: 2, Imag: 2, Char: 1, String: 1}}

	ast.Inspect(file, func(n ast.Node) bool {
		m.ParseNode(n)
		return true
	})
	got := m.Comp
	want := uint(43)
	if want != got {
		t.Errorf("TestCodeStructComp got %v, want %v", got, want)
	}
}
