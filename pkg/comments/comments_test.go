package comments

import (
	"go/ast"
	"testing"
)

/// Test with predefined string
func TestLoremIpsum(t *testing.T) {
	comment := ast.Comment{}
	comment.Text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, "
	comment.Text += "sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. "
	comment.Text += "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris "
	comment.Text += "nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in "
	comment.Text += "reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. "
	comment.Text += "Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	nwords := uint(69)
	weight := uint(2)
	want := nwords * weight

	got := GetCommentComp(&comment, &Weights{weight})
	if got != want {
		t.Fatalf(`GetCommentComp("Lorem ipsum...") = %v, Wanted %v`, got, want)
	}
}
