package cyclo

import (
	"go/ast"
	"go/token"
)

type Weights struct {
	If   float64 `xml:"if"`
	For  float64 `xml:"for"`
	Rng  float64 `xml:"rng"`
	Case float64 `xml:"case"`
	And  float64 `xml:"and"`
	Or   float64 `xml:"or"`
}

type Metric struct {
	Config Weights
	Comp   float64
}

func (m *Metric) ParseNode(n ast.Node) {
	v, ok := n.(*ast.FuncDecl)
	if ok {
		m.Comp += getCycloComp(v, &m.Config)
	}
}

func (m Metric) Finish() float64 {
	return m.Comp
}

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
