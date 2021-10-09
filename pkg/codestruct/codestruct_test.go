package codestruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestCodeStructComp(t *testing.T) {
	src := `package main

			type foo struct {}
			type bar struct {}
			
			type foobar interface {}
			
			func f1(){}
			func f2(){}
			
			func main(){}`

	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	got := uint(0)

	ast.Inspect(file, func(n ast.Node) bool {
		got += getCodeStructComp(n, &Weights{Struct: 1, Func: 1, Interface: 1})
		return true
	})

	want := uint(6)
	if want != got {
		t.Errorf("TestCodeStructComp got %v, want %v", got, want)
	}
}
