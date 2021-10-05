package comments

import (
	"go/ast"
	"strings"
)

type Weights struct {
	Word uint `xml:"word"`
}

func GetCommentComp(comment *ast.Comment, config *Weights) uint {
	return uint(len(strings.Fields(comment.Text))) * config.Word
}
