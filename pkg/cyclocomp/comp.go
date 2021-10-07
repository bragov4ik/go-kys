package cyclo

import (
	"go/ast"
	"go/token"
)

type Weights struct {
	If   uint `xml:"if"`
	For  uint `xml:"for"`
	Rng  uint `xml:"rng"`
	Case uint `xml:"case"`
	And  uint `xml:"and"`
	Or   uint `xml:"or"`
}

type branchVisitor func(n ast.Node) (w ast.Visitor)

func (v branchVisitor) Visit(n ast.Node) (w ast.Visitor) {
	return v(n)
}

type Metric struct {
	Config Weights
	Comp   uint
}

func (m *Metric) ParseNode(n ast.Node) {
	v, ok := n.(*ast.FuncDecl)
	if ok {
		m.Comp += GetCycloComp(v, &m.Config)
	}
}

func (m Metric) Finish() float64 {
	return float64(m.Comp)
}

func GetCycloComp(fd *ast.FuncDecl, config *Weights) uint {
	var comp uint = 1
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
