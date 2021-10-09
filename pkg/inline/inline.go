package inline

import (
	"go/ast"
	"go/token"
)

type Weights struct {
	Int    uint
	Float  uint
	Imag   uint
	Char   uint
	String uint
}

type Metric struct {
	Config Weights
	Comp   uint
}

func (m *Metric) ParseNode(n ast.Node) {
	v, ok := n.(*ast.BasicLit)
	if ok {
		m.Comp += getInlineComp(v, &m.Config)
	}
}

func (m Metric) Finish() float64 {
	return float64(m.Comp)
}

func getInlineComp(literal *ast.BasicLit, config *Weights) uint {
	var comp uint
	switch literal.Kind {
	case token.INT:
		comp = config.Int
	case token.FLOAT:
		comp = config.Float
	case token.IMAG:
		comp = config.Imag
	case token.CHAR:
		comp = config.Char
	case token.STRING:
		comp = config.String
	}
	return comp
}
