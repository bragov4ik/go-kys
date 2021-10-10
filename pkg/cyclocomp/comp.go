// Checks cyclo complexity of code using metric
package cyclo

import (
	"go/ast"
	"go/token"
)

// Config for metric with various weights for syntactical structures
type Weights struct {
	// If weight
	If float64 `xml:"if"`
	// For weight
	For float64 `xml:"for"`
	// Range weight
	Rng float64 `xml:"rng"`
	// Case weight
	Case float64 `xml:"case"`
	// Boolean and weight
	And float64 `xml:"and"`
	// Boolean or weight
	Or float64 `xml:"or"`
}

// Intermidiate state of metric
type Metric struct {
	// Config with weights
	Config Weights
	comp   float64
}

// Parses ast node and collects all of its metrics
func (m *Metric) ParseNode(n ast.Node) {
	if v, ok := n.(*ast.FuncDecl); ok {
		m.comp += getCycloComp(v, &m.Config)
	}
}

// Returns final score
func (m Metric) Finish() float64 { return m.comp }

type branchVisitor func(n ast.Node) (w ast.Visitor)

func (v branchVisitor) Visit(n ast.Node) (w ast.Visitor) {
	return v(n)
}

func getCycloComp(fd *ast.FuncDecl, config *Weights) float64 {
	var comp float64 = 1
	var v ast.Visitor
	v = branchVisitor(func(n ast.Node) (w ast.Visitor) {
		switch n := n.(type) {
		case *ast.IfStmt:
			comp += config.If
		case *ast.ForStmt:
			comp += config.For
		case *ast.RangeStmt:
			comp += config.Rng
		case *ast.CaseClause, *ast.CommClause:
			comp += config.Case
		case *ast.BinaryExpr:
			if n.Op == token.LAND {
				comp += config.And
			} else if n.Op == token.LOR {
				comp += config.Or
			}
		}
		return v
	})
	ast.Walk(v, fd)

	return comp
}
