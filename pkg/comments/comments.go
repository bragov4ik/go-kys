package comments

import (
	"go/ast"
	"strings"
)

type Weights struct {
	Word float64 `xml:"word"`
}

type Metric struct {
	Config Weights
	Score  float64
}

func (m *Metric) ParseNode(n ast.Node) {
	v, ok := n.(*ast.Comment)
	if ok {
		m.Score += getCommentComp(v, &m.Config)
	}
}

func (m Metric) Finish() float64 {
	return m.Score
}

func getCommentComp(comment *ast.Comment, config *Weights) float64 {
	return float64(len(strings.Fields(comment.Text))) * config.Word
}
