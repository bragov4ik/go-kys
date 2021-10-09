package inline

import (
	"go/ast"
	"go/token"
)

type Weights struct {
	Int          uint `xml:"int"`
	Float        uint `xml:"float"`
	Imag         uint `xml:"imag"`
	Char         uint `xml:"char"`
	String       uint `xml:"string"`
	CompositeLit uint `xml:"composite"`
}

type Metric struct {
	Config Weights
	Comp   uint
}

func (m *Metric) ParseNode(n ast.Node) {
	switch v := n.(type) {
	case *ast.BasicLit:
		m.Comp += getBasicLitComp(v, &m.Config)
	case *ast.CompositeLit:
		m.Comp += m.Config.CompositeLit
	}
}

func (m Metric) Finish() float64 {
	return float64(m.Comp)
}

func getBasicLitComp(literal *ast.BasicLit, config *Weights) uint {
	var comp uint
	len := uint(len(literal.Value))
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
		comp = config.String * len
	}
	return comp
}
