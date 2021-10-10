// Package with metric which checks comment complexity
package comments

import (
	"go/ast"
	"strings"
)

// Config for metric
type Weights struct {
	// Weight of each word in every comment
	Word float64 `xml:"word"`
}

// Intermidiate state of metric
type Metric struct {
	// Config with weights
	Config Weights
	score  float64
}

// Parses ast node and collects metric result
func (m *Metric) ParseNode(n ast.Node) {
	if v, ok := n.(*ast.Comment); ok {
		m.score += getCommentComp(v, &m.Config)
	}
}

// Returns metric result
func (m Metric) Finish() float64 { return m.score }

func getCommentComp(comment *ast.Comment, config *Weights) float64 {
	return float64(len(strings.Fields(comment.Text))) * config.Word
}
