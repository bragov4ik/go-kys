package comments

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

/// Test with a predefined string
func TestLorem(t *testing.T) {
	src := `package main

			// Lorem ipsum dolor sit amet, consectetur adipiscing elit,
			// sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
			// Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris
			// nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in
			// reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
			// Excepteur sint occaecat cupidatat non proident,
			// sunt in culpa qui officia deserunt mollit anim id est laborum.
			func main(){}`

	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	nwords := uint(76)
	weight := 2.
	want := uint(float64(nwords) * weight)

	m := Metric{Config: Weights{Word: weight}}

	ast.Inspect(file, func(n ast.Node) bool {
		m.ParseNode(n)
		return true
	})

	got := uint(m.Finish())

	if want != got {
		t.Fatalf(`GetCommentComp("Lorem ipsum...") = %v, Wanted %v`, got, want)
	}
}
