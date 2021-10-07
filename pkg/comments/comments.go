package comments

import (
	"go/ast"
	"strings"
)

type Weights struct {
	Word uint `xml:"word"`
}

type Metric struct {
	Config Weights
	Score  uint
}

func (m *Metric) ParseNode(n ast.Node) {
	v, ok := n.(*ast.Comment)
	if ok {
		m.Score += GetCommentComp(v, &m.Config)
	}
}

func (m Metric) Finish() float64 {
	return float64(m.Score)
}

func GetCommentComp(comment *ast.Comment, config *Weights) uint {
	return uint(len(strings.Fields(comment.Text))) * config.Word
}
